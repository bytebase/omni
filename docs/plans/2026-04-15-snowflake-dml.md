# T5.1: Snowflake DML Implementation Plan

**Goal:** Implement INSERT/UPDATE/DELETE/MERGE parsing for the Snowflake engine.

**Architecture:** Four new parser files (`insert.go`, `update.go`, `delete.go`, `merge.go`) and corresponding test files. New AST node types appended to `parsenodes.go`, tags to `nodetags.go`, walker cases regenerated.

---

## File Map

| File | Action | Responsibility |
|------|--------|----------------|
| `snowflake/ast/parsenodes.go` | Modify | Add InsertStmt, InsertMultiStmt, InsertMultiBranch, UpdateStmt, UpdateSet, DeleteStmt, MergeStmt, MergeWhen, MergeAction |
| `snowflake/ast/nodetags.go` | Modify | Add T_InsertStmt, T_InsertMultiStmt, T_UpdateStmt, T_DeleteStmt, T_MergeStmt tags |
| `snowflake/ast/walk_generated.go` | Regenerate | Auto-generated walker cases |
| `snowflake/parser/parser.go` | Modify | Replace 4 unsupported stubs with real calls |
| `snowflake/parser/insert.go` | Create | parseInsertStmt, parseInsertMultiStmt |
| `snowflake/parser/update.go` | Create | parseUpdateStmt |
| `snowflake/parser/delete.go` | Create | parseDeleteStmt |
| `snowflake/parser/merge.go` | Create | parseMergeStmt |
| `snowflake/parser/insert_test.go` | Create | INSERT tests |
| `snowflake/parser/update_test.go` | Create | UPDATE tests |
| `snowflake/parser/delete_test.go` | Create | DELETE tests |
| `snowflake/parser/merge_test.go` | Create | MERGE tests |

---

## Steps

- [ ] Step 1: Add AST types to parsenodes.go
- [ ] Step 2: Add node tags to nodetags.go
- [ ] Step 3: Regenerate walk_generated.go
- [ ] Step 4: Implement insert.go + insert_test.go
- [ ] Step 5: Implement update.go + update_test.go
- [ ] Step 6: Implement delete.go + delete_test.go
- [ ] Step 7: Implement merge.go + merge_test.go
- [ ] Step 8: Wire parseStmt dispatch in parser.go
- [ ] Step 9: go test ./snowflake/... + gofmt -w
- [ ] Step 10: Commit + PR
