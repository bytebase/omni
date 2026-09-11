# Snowflake Set Operators (T1.7) — Design

**DAG node:** T1.7 — set operators
**Branch:** `feat/snowflake/set-ops`
**Depends on:** T1.4 (SELECT core).
**Unblocks:** T1.8 (statement classification), Tier 2+ nodes.

## Purpose

T1.7 adds UNION/UNION ALL/EXCEPT/MINUS/INTERSECT operators that chain multiple SELECT statements. Currently `parseSelectStmt()` returns one SelectStmt. T1.7 wraps chained SELECTs in a `SetOperationStmt` node.

## Scope

From the legacy grammar's `set_operators` rule:
- `UNION [ALL] [BY NAME]` — combines results, optional dedup, optional column-name matching
- `EXCEPT` — set difference
- `MINUS` — alias for EXCEPT (Snowflake treats them identically)
- `INTERSECT` — set intersection

Each is followed by another SELECT statement (which may itself be wrapped in parens or have its own set operators).

## AST Types

```go
type SetOperationStmt struct {
    Op    SetOp
    All   bool    // UNION ALL
    ByName bool   // UNION [ALL] BY NAME (Snowflake-specific)
    Left  Node    // *SelectStmt or nested *SetOperationStmt
    Right Node    // *SelectStmt or nested *SetOperationStmt
    Loc   Loc
}
func (n *SetOperationStmt) Tag() NodeTag { return T_SetOperationStmt }

type SetOp int
const (
    SetOpUnion     SetOp = iota
    SetOpExcept
    SetOpIntersect
)
```

MINUS maps to `SetOpExcept` (they're semantically identical in Snowflake).

## Parser Changes

The key insight: set operators are parsed OUTSIDE `parseSelectStmt()`. The dispatch flow:

```
parseStmt dispatch:
  case kwSELECT: return p.parseQueryExpr()   // NEW: wraps SELECT + set ops
  case kwWITH:   return p.parseWithQueryExpr() // NEW: wraps WITH + SELECT + set ops

parseQueryExpr:
  left = parseSelectStmt()
  while current is UNION/EXCEPT/MINUS/INTERSECT:
    op, all, byName = parseSetOp()
    right = parseSelectStmt()  // or (SELECT ...) in parens
    left = &SetOperationStmt{Op: op, Left: left, Right: right, All: all, ByName: byName}
  return left  // either bare SelectStmt or SetOperationStmt
```

This means:
- `parseStmt` dispatch calls `parseQueryExpr()` (not `parseSelectStmt()` directly)
- `parseQueryExpr()` is a thin wrapper that parses one SELECT then loops on set operators
- If no set operators follow, returns the bare `*SelectStmt`
- If set operators follow, returns a `*SetOperationStmt` wrapping left + right

### Precedence

Standard SQL: INTERSECT binds tighter than UNION/EXCEPT. So `SELECT 1 UNION SELECT 2 INTERSECT SELECT 3` = `SELECT 1 UNION (SELECT 2 INTERSECT SELECT 3)`.

For simplicity in T1.7, treat all set operators as left-associative at the same precedence. This is what the legacy grammar does (no precedence distinction). Can refine later if needed.

### Parenthesized SELECT in set operators

`(SELECT ...) UNION (SELECT ...)` — the parens are optional around each SELECT. The legacy grammar's `select_statement_in_parentheses` allows recursive parens. T1.7 handles this: if `(` is current before a set operator's right side, check for SELECT inside and parse accordingly.

## File Layout

| File | Change |
|------|--------|
| `snowflake/ast/parsenodes.go` | +SetOperationStmt, SetOp enum |
| `snowflake/ast/nodetags.go` | +T_SetOperationStmt |
| `snowflake/ast/loc.go` | +*SetOperationStmt case |
| `snowflake/ast/walk_generated.go` | Regenerated |
| `snowflake/parser/select.go` | Add parseQueryExpr, parseSetOp; update dispatch |
| `snowflake/parser/parser.go` | Update SELECT/WITH dispatch to call parseQueryExpr/parseWithQueryExpr |
| `snowflake/parser/select_test.go` | Add set operator tests |

**Estimated: ~300 LOC** across 7 modified files.

## Testing (8 categories)

1. `SELECT 1 UNION SELECT 2`
2. `SELECT 1 UNION ALL SELECT 2`
3. `SELECT 1 UNION ALL BY NAME SELECT 2`
4. `SELECT 1 EXCEPT SELECT 2`
5. `SELECT 1 MINUS SELECT 2` (same as EXCEPT)
6. `SELECT 1 INTERSECT SELECT 2`
7. Chained: `SELECT 1 UNION SELECT 2 UNION SELECT 3` — left-associative
8. `(SELECT 1) UNION (SELECT 2)` — parenthesized
9. With CTE: `WITH cte AS (SELECT 1) SELECT * FROM cte UNION SELECT 2`

## Acceptance Criteria

1. Build/vet/gofmt/test all clean
2. `Parse("SELECT 1 UNION ALL SELECT 2;")` returns *SetOperationStmt
3. MINUS maps to SetOpExcept
4. Walker regen correct
