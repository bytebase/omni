# PartiQL Parser-DDL — Design Spec

**DAG node:** 7 (parser-ddl), P0
**Depends on:** parser-foundation (node 4, done)
**Unblocks:** parse-entry (8)
**Package:** `partiql/parser`
**Files added:** `ddl.go`
**Files modified:** `parser.go` (add ParseStatement — the new statement-level entry point), `expr.go` (remove INSERT/UPDATE/DELETE stubs from parseSelectExpr), `parser_test.go`

---

## 1. Goal

Add CREATE TABLE, CREATE INDEX, DROP TABLE, and DROP INDEX parsing. PartiQL's DDL grammar is intentionally minimal (DynamoDB's DDL is done via API, not SQL). This is the smallest of the 3 parallel parser nodes.

## 2. Scope

### In scope
- `createCommand` — CREATE TABLE symbolPrimitive (no column definitions in PartiQL DDL)
- `createCommand` — CREATE INDEX ON symbolPrimitive ( pathSimple, pathSimple, ... )
- `dropCommand` — DROP TABLE symbolPrimitive
- `dropCommand` — DROP INDEX symbolPrimitive ON symbolPrimitive
- `pathSimple` for index key paths (shared with DML if node 6 adds it; otherwise implemented locally)

### Out of scope
- Column definitions in CREATE TABLE (PartiQL is schema-on-read)
- ALTER TABLE (not in the grammar)

## 3. Architecture

### 3.1 File layout

| File | Responsibility | Est. lines |
|------|---------------|------------|
| `ddl.go` | parseCreateCommand, parseDropCommand, parsePathSimple (if not already added by node 6) | ~120 |

### 3.2 Dispatch integration

**parser.go changes — introduce ParseStatement:**
Node 7 is the first node to introduce `ParseStatement`, the public statement-level entry point. It dispatches as follows:

```
tokCREATE → parseCreateCommand (ddl.go)
tokDROP   → parseDropCommand   (ddl.go)
tokINSERT / tokUPDATE / tokDELETE / tokREPLACE / tokUPSERT / tokREMOVE
          → deferredFeature("...", "parser-dml (DAG node 6)")
tokEXEC / tokEXECUTE
          → deferredFeature("EXEC", "parse-entry (DAG node 8)")
default   → ParseExpr() DQL fallback:
              call parseExprTop(); if result implements ast.StmtNode, return it;
              otherwise deferredFeature("bare expression as statement", "parse-entry (DAG node 8)")
```

The DQL fallback works cleanly because `SelectStmt` and `SetOpStmt` both implement both `ast.ExprNode` and `ast.StmtNode`. When the input starts with SELECT, `parseExprTop` → `parseSelectExpr` → `parseSFWQuery` returns a `*SelectStmt` which satisfies `ast.StmtNode`. For bare expressions that are not statements (e.g., `1 + 2`), the deferred error is returned and resolved by node 8 (parse-entry).

**expr.go changes:**
- Remove the INSERT/UPDATE/DELETE stubs from `parseSelectExpr`. Those stubs served as deferred-feature markers during parser-foundation; now that `ParseStatement` exists for statement-level dispatch, they are no longer needed in the expression ladder. CREATE/DROP were never stubbed in `parseSelectExpr` — they route exclusively through `ParseStatement`.

### 3.3 AST nodes

All exist in ast-core stmts.go:
- `CreateTableStmt` — Name field
- `CreateIndexStmt` — Table, Keys fields
- `DropTableStmt` — Name field
- `DropIndexStmt` — Index, Table fields

### 3.4 Stub replacement audit

Stubs replaced by this node:
- INSERT/UPDATE/DELETE deferred-feature stubs in `parseSelectExpr` (expr.go) — removed entirely; DML now dispatches from `ParseStatement` (node 6 will fill those cases in).

Stubs added by this node (to be filled by later nodes):
- `tokINSERT/UPDATE/DELETE/REPLACE/UPSERT/REMOVE` → `deferredFeature(...)` in ParseStatement — filled by node 6 (parser-dml)
- `tokEXEC/EXECUTE` → `deferredFeature(...)` in ParseStatement — filled by node 8 (parse-entry)
- Default DQL bare-expression path → `deferredFeature(...)` in ParseStatement — filled by node 8 (parse-entry)

## 4. Testing

- `testdata/parser-ddl/` — ~8 golden pairs (CREATE TABLE, CREATE INDEX, DROP TABLE, DROP INDEX variants)
- `TestParser_StmtGoldens` — a new test harness (parallel to `TestParser_Goldens`) that calls `ParseStatement()` instead of `ParseExpr()`. Golden files live under `testdata/parser-ddl/`. The existing `TestParser_Goldens` harness (calling `ParseExpr`) is unchanged and continues to test expression-level parsing.
- No AWS corpus impact (corpus has no DDL files)
- Error tests (table-driven in `parser_test.go`) for: CREATE without TABLE/INDEX, DROP without TABLE/INDEX, missing ON in DROP INDEX, missing identifier after TABLE/INDEX keyword, and unknown sub-keywords after CREATE/DROP

## 5. Design decisions

- **D1** DDL dispatches from `ParseStatement` (introduced by this node), NOT from `parseSelectExpr`. Statement-level dispatch is cleanly separated from expression-level parsing. `parseSelectExpr` remains purely for expressions (SELECT, set operations, and value expressions).
- **D2** CREATE TABLE has no column definitions (PartiQL DDL is minimal; schema-on-read)
- **D3** pathSimple reuses the same implementation if node 6 lands first; otherwise implements locally with a note that node 6 may deduplicate
- **D4** No EXEC command in this node — it is stubbed in ParseStatement with a deferredFeature error pointing to node 8 (parse-entry)
- **D5** The DQL fallback in ParseStatement uses a type assertion (`expr.(ast.StmtNode)`) rather than a type switch. Both `SelectStmt` and `SetOpStmt` implement `ast.StmtNode`, so the assertion covers all select/bag-op outputs cleanly. Bare expressions that are not statements are deferred to node 8.
