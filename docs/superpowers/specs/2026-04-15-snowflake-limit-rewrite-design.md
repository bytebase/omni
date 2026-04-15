# T3.3 — Snowflake LIMIT Injection Rewrite Design

## Scope

T3.3 adds `InjectLimit(sql string, maxRows int) (string, error)` to the
`snowflake/deparse` package. It is called by the Bytebase read-only SQL editor
to cap result-set size before a query is sent to Snowflake.

## Public API

```go
// InjectLimit rewrites the given SQL to cap the result-set size to maxRows.
func InjectLimit(sql string, maxRows int) (string, error)
```

## File placement

New file `snowflake/deparse/rewrite.go` in the existing `deparse` package.
No new packages are created.

## Rewrite rules

| Input state | Action |
|---|---|
| No LIMIT, no FETCH | Add `LIMIT maxRows` |
| `LIMIT n` with n <= maxRows | Leave unchanged |
| `LIMIT n` with n > maxRows | Replace with `LIMIT maxRows` |
| `LIMIT expr` (non-literal) | Wrap: `SELECT * FROM (...) LIMIT maxRows` |
| `FETCH FIRST n ROWS ONLY` with n <= maxRows | Leave unchanged |
| `FETCH FIRST n ROWS ONLY` with n > maxRows | Replace with `FETCH FIRST maxRows ROWS ONLY` |
| `FETCH FIRST expr ROWS ONLY` (non-literal) | Wrap: `SELECT * FROM (...) LIMIT maxRows` |
| `SetOperationStmt` (UNION/INTERSECT/EXCEPT) | Always wrap |
| Non-SELECT (INSERT/UPDATE/DELETE/MERGE/DDL) | Return unchanged |
| Multi-statement | Apply per-statement independently |

### FETCH vs LIMIT normalisation

When a FETCH clause must be lowered, it is lowered in-place (the FETCH clause
is kept as FETCH). If the FETCH clause uses a non-literal expression the
statement is wrapped so semantics are preserved.

When a statement has no LIMIT and no FETCH, we add a LIMIT clause (not FETCH),
keeping the output uniform.

### SetOperationStmt wrapping

Set operations may legally carry a trailing LIMIT only at the outermost level
in Snowflake (applied to the entire result). Because the AST node for a set
operation is a `SetOperationStmt` rather than a `SelectStmt`, we always wrap:

```sql
SELECT * FROM (<original>) LIMIT maxRows
```

A subquery alias (`_q`) is applied so the generated SQL is valid.

### Wrap semantics

```sql
SELECT * FROM (<inner query>) AS _q LIMIT maxRows
```

The outer SELECT uses a bare `*` target (represented in the AST as a
`SelectTarget{Star: true, Expr: &ast.StarExpr{}}`).

## Error handling

- `maxRows <= 0` → return `"", fmt.Errorf("maxRows must be positive")`
- Parse error → forward the error from `parser.Parse` immediately; do not
  attempt partial rewrites.

## Dependencies

- `github.com/bytebase/omni/snowflake/ast` — AST node types
- `github.com/bytebase/omni/snowflake/parser` — `Parse`
- `github.com/bytebase/omni/snowflake/deparse` — `DeparseFile` (same package)
