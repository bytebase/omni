# PartiQL Parser-DDL — Design Spec

**DAG node:** 7 (parser-ddl), P0
**Depends on:** parser-foundation (node 4, done)
**Unblocks:** parse-entry (8)
**Package:** `partiql/parser`
**Files added:** `ddl.go`
**Files modified:** `expr.go` (add CREATE/DROP dispatch in parseSelectExpr), `parser_test.go`

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

**expr.go changes:**
- `parseSelectExpr`: add `tokCREATE` and `tokDROP` detection with calls to parseCreateCommand / parseDropCommand in ddl.go.

### 3.3 AST nodes

All exist in ast-core stmts.go:
- `CreateTableStmt` — Name field
- `CreateIndexStmt` — Table, Keys fields
- `DropTableStmt` — Name field
- `DropIndexStmt` — Index, Table fields

### 3.4 Stub replacement audit

Foundation stubs being replaced:
- None from foundation (CREATE/DROP weren't stubbed; they hit the default "unexpected token" error). This node ADDS dispatch for these tokens.

## 4. Testing

- `testdata/parser-ddl/` — ~8 golden pairs
- No AWS corpus impact (corpus has no DDL files)
- Error tests for CREATE without TABLE/INDEX, DROP without TABLE/INDEX, missing ON in DROP INDEX

## 5. Design decisions

- **D1** DDL dispatch in parseSelectExpr alongside SELECT/DML, matching the grammar's statement-level position
- **D2** CREATE TABLE has no column definitions (PartiQL DDL is minimal)
- **D3** pathSimple reuses the same implementation if node 6 lands first; otherwise implements locally with a note that node 6 may deduplicate
- **D4** No EXEC command in this node — that's a separate concern potentially for node 8 or a P1 follow-up
