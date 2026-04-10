# PartiQL Parser-DML — Design Spec

**DAG node:** 6 (parser-dml), P0
**Depends on:** parser-foundation (node 4, done)
**Unblocks:** parse-entry (8)
**Package:** `partiql/parser`
**Files added:** `dml.go`
**Files modified:** `parser.go` (extend ParseStatement with DML cases, replacing deferred-feature stubs from node 7), `exprprimary.go` (replace VALUES/valueList stubs), `parser_test.go`
**AST additions:** may need new types in `partiql/ast/stmts.go` for ReplaceStmt, UpsertStmt, RemoveStmt, SetCmd if not already present

---

## 1. Goal

Replace the INSERT/UPDATE/DELETE/VALUES/valueList deferred-feature stubs from parser-foundation with real implementations. After this node, the parser handles all PartiQL DML operations. The AWS corpus smoke test should see ~28 more files fully parse (those starting with INSERT/UPDATE/DELETE).

## 2. Scope

### In scope
- `insertCommand` — legacy form (`INSERT INTO path VALUE expr [AT expr]`) and new RFC form (`INSERT INTO symbol [AS alias] value=expr`)
- `deleteCommand` — `DELETE fromClauseSimple [WHERE expr] [RETURNING ...]`
- `updateClause` + `dmlBaseCommand` (compound DML: UPDATE table SET/REMOVE/INSERT)
- `setCommand` — `SET path = expr, ...`
- `replaceCommand` — `REPLACE INTO symbol [AS alias] value=expr`
- `upsertCommand` — `UPSERT INTO symbol [AS alias] value=expr`
- `removeCommand` — `REMOVE pathSimple`
- `onConflictClause` — ON CONFLICT [target] DO NOTHING/REPLACE/UPDATE
- `returningClause` — RETURNING (MODIFIED|ALL) (OLD|NEW) (* | expr)
- `pathSimple` — simplified path for DML targets (symbol + literal/symbol bracket steps + dot steps)
- `fromClauseSimple` — FROM pathSimple [AS alias] [AT alias] [BY alias]

### Out of scope
- `VALUES (row), (row)` row-set constructor — requires a ValuesExpr AST node. If ast-core doesn't have it, this node adds it to `partiql/ast/exprs.go`.

## 3. Architecture

### 3.1 File layout

| File | Responsibility | Est. lines |
|------|---------------|------------|
| `dml.go` | parseInsertStmt, parseDeleteStmt, parseUpdateStmt, parseReplaceStmt, parseUpsertStmt, parseRemoveStmt, parseSetCommand, parseOnConflict, parseReturningClause, parsePathSimple, parseFromClauseSimple | ~350 |

### 3.2 Dispatch integration

**parser.go changes:**
- DML dispatches from `ParseStatement`, NOT from `parseSelectExpr`. Node 6 adds `tokINSERT`, `tokUPDATE`, `tokDELETE`, `tokREPLACE`, `tokUPSERT`, `tokREMOVE` cases to `ParseStatement` in `parser.go`, replacing the deferred-feature stubs that node 7 (parser-ddl) initially placed there.
- `ParseStatement` is introduced by node 7 (parser-ddl). Node 6 extends it with DML cases.
- `parseSelectExpr` is NOT modified by this node. The INSERT/UPDATE/DELETE stubs previously in `parseSelectExpr` are removed as part of node 7's work (those stubs served only as markers during foundation; real dispatch is at statement level).
- DQL (SELECT) in statement position continues to work via `ParseStatement`'s default fallback: it calls `ParseExpr`, which descends through the expression ladder and reaches `parseSelectExpr` → `parseSFWQuery`.

**exprprimary.go changes:**
- Replace `tokVALUES` stub with a real `parseValues` implementation (or a better error if a ValuesExpr AST node isn't added).
- Replace the valueList `(expr, expr, ...)` stub — if a ValuesExpr exists, parse it; otherwise keep the deferred error but improve the message.

### 3.3 Stub replacement audit

Foundation stubs being replaced:
- `ParseStatement` INSERT/UPDATE/DELETE/REPLACE/UPSERT/REMOVE deferred-feature stubs (parser.go) — these were placed by node 7 and are now filled in by node 6
- `parsePrimaryBase` VALUES stub (exprprimary.go)
- `parseParenExpr` valueList stub (exprprimary.go)

### 3.4 AST nodes

Existing in ast-core stmts.go: `InsertStmt`, `UpdateStmt`, `DeleteStmt`, `ExecStmt`

May need to add: `ReplaceStmt`, `UpsertStmt`, `RemoveStmt`, `SetCmd`, `OnConflict`, `ReturningColumn`. Check `partiql/ast/stmts.go` before implementation.

## 4. Testing

- `testdata/parser-dml/` — ~15 golden pairs
- `TestParser_AWSCorpus` — ~28 more files fully parse
- Error tests for missing VALUE after INSERT, missing SET after UPDATE, etc.

## 5. Design decisions

- **D1** DML parsing dispatches from `ParseStatement` at the statement level, not from `parseSelectExpr`. This cleanly separates statement-level dispatch from expression-level parsing. `ParseStatement` is introduced by node 7 (parser-ddl); node 6 extends it with DML cases.
- **D2** pathSimple is a separate parser from the expression-level pathStep (simpler: only symbol + literal/symbol bracket/dot steps, no expression indices or wildcards)
- **D3** ON CONFLICT DO REPLACE/UPDATE bodies are minimal stubs (grammar's doReplace/doUpdate rules have TODO comments in the original g4 file)
- **D4** Compound DML (UPDATE ... SET ... REMOVE ...) parsed iteratively: updateClause first, then loop of dmlBaseCommand
