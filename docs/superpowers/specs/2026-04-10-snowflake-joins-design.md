# Snowflake Joins (T1.5) — Design

**DAG node:** T1.5 — joins
**Branch:** `feat/snowflake/joins`
**Depends on:** T1.4 (SELECT core).
**Unblocks:** T1.7 (set operators), T1.8 (statement classification), T2.x (DDL with subqueries).

## Purpose

T1.5 extends T1.4's FROM clause from simple comma-separated table references to full JOIN syntax. After T1.5, the parser handles every FROM-clause form in the legacy grammar.

## AST Types

### JoinExpr (Node)

```go
type JoinExpr struct {
    Type           JoinType
    Left           Node    // TableRef or nested JoinExpr
    Right          Node    // TableRef or nested JoinExpr
    On             Node    // ON condition; nil for CROSS/NATURAL/USING-only
    Using          []Ident // USING columns; nil if ON or NATURAL
    Natural        bool    // NATURAL JOIN
    Directed       bool    // DIRECTED hint (Snowflake-specific)
    MatchCondition Node    // ASOF MATCH_CONDITION(expr); nil for non-ASOF
    Loc            Loc
}
func (n *JoinExpr) Tag() NodeTag { return T_JoinExpr }

type JoinType int
const (
    JoinInner JoinType = iota // [INNER] JOIN
    JoinLeft                   // LEFT [OUTER] JOIN
    JoinRight                  // RIGHT [OUTER] JOIN
    JoinFull                   // FULL [OUTER] JOIN
    JoinCross                  // CROSS JOIN
    JoinAsof                   // ASOF JOIN (Snowflake-specific)
)
```

### Extended TableRef (3 new fields)

```go
type TableRef struct {
    Name     *ObjectName    // table name; nil for subquery/func sources
    Alias    Ident          // AS alias; zero if absent
    Subquery Node           // (SELECT ...) in FROM; nil for table refs
    FuncCall *FuncCallExpr  // TABLE(func(...)); nil for table refs
    Lateral  bool           // LATERAL prefix
    Loc      Loc
}
```

A TableRef is now polymorphic:
- **Table**: Name is set, others nil
- **Subquery**: Subquery is set, Name is nil
- **Table function**: FuncCall is set, Name is nil
- Any of the above can have `Lateral = true`

### SelectStmt.From type change

```go
// Was: From []*TableRef
From []Node // mixed TableRef and JoinExpr
```

This is a breaking change to T1.4's type. Since we control all consumers (only select_test.go references `From`), the fix is to update those test assertions.

### NodeTag addition

One new tag: `T_JoinExpr`.

## Parser Changes

### Revised `parseFromClause()`

Returns `[]ast.Node` instead of `[]*ast.TableRef`. Parses comma-separated "from items" where each item is a primary source optionally followed by a chain of JOINs.

```go
func (p *Parser) parseFromClause() ([]ast.Node, error) {
    // consume FROM
    var items []ast.Node
    item, err := p.parseFromItem()
    items = append(items, item)
    for p.cur.Type == ',' {
        p.advance()
        item, err = p.parseFromItem()
        items = append(items, item)
    }
    return items, nil
}
```

### New: `parseFromItem()`

Parses one comma-separated FROM item: a primary source + optional JOIN chain.

```go
func (p *Parser) parseFromItem() (ast.Node, error) {
    left, err := p.parsePrimarySource()
    return p.parseJoinChain(left)
}
```

### New: `parsePrimarySource()`

Dispatches on the current token:
- `(` → subquery (`(SELECT ...)`) or parenthesized from-item (`(t1 JOIN t2 ON ...)`)
- `kwLATERAL` → consume, then recurse (sets Lateral=true on result)
- `kwTABLE` → consume `TABLE(`, parse function call, wrap in TableRef{FuncCall: ...}
- Otherwise → `parseTableRef()` (existing: ObjectName + alias)

### New: `parseJoinChain(left)`

Left-associative loop that builds a JoinExpr tree:

```go
func (p *Parser) parseJoinChain(left ast.Node) (ast.Node, error) {
    for {
        joinType, natural, directed, ok := p.parseJoinKeywords()
        if !ok { break }
        right, err := p.parsePrimarySource()
        // parse ON / USING / MATCH_CONDITION
        join := &ast.JoinExpr{Type: joinType, Left: left, Right: right, Natural: natural, Directed: directed, ...}
        left = join
    }
    return left, nil
}
```

### New: `parseJoinKeywords()`

Recognizes all JOIN keyword combinations:
- `[INNER] JOIN` → JoinInner
- `LEFT [OUTER] JOIN` → JoinLeft
- `RIGHT [OUTER] JOIN` → JoinRight
- `FULL [OUTER] JOIN` → JoinFull
- `CROSS JOIN` → JoinCross
- `NATURAL [LEFT|RIGHT|FULL] JOIN` → sets natural=true + appropriate type
- `ASOF JOIN` → JoinAsof
- `DIRECTED [INNER|LEFT|...] JOIN` → sets directed=true + appropriate type
- `JOIN` alone → JoinInner (implicit INNER)

Returns `(joinType, natural, directed, ok)` where ok=false if current position is not a join keyword.

### Join condition parsing

After the right source:
- If `kwON`: parse expression → `JoinExpr.On`
- If `kwUSING`: expect `(`, parse ident list, expect `)` → `JoinExpr.Using`
- If ASOF: expect `kwMATCH_CONDITION` (or tokIdent "MATCH_CONDITION"), expect `(`, parse expr, expect `)` → `JoinExpr.MatchCondition`
- CROSS and NATURAL joins: no condition required
- All other joins: ON or USING is expected (error if absent)

## File Layout

| File | Change | Approx LOC |
|------|--------|-----------|
| `snowflake/ast/parsenodes.go` | MODIFY: extend TableRef + add JoinExpr + JoinType | +60 |
| `snowflake/ast/nodetags.go` | MODIFY: +T_JoinExpr | +3 |
| `snowflake/ast/loc.go` | MODIFY: +*JoinExpr case | +2 |
| `snowflake/ast/walk_generated.go` | REGENERATED | auto |
| `snowflake/parser/select.go` | MODIFY: replace parseFromClause + add parseFromItem/parsePrimarySource/parseJoinChain/parseJoinKeywords | +250 |
| `snowflake/parser/select_test.go` | MODIFY: update From assertions + add 15 join test categories | +300 |

**Estimated total: ~600 LOC** across 0 new + 6 modified files.

## Testing (15 categories)

1. `FROM t1 JOIN t2 ON t1.id = t2.id` — basic INNER
2. `FROM t1 LEFT JOIN t2 ON ...` — LEFT
3. `FROM t1 RIGHT OUTER JOIN t2 ON ...` — RIGHT OUTER
4. `FROM t1 FULL JOIN t2 ON ...` — FULL
5. `FROM t1 CROSS JOIN t2` — CROSS (no ON)
6. `FROM t1 NATURAL JOIN t2` — NATURAL (no ON/USING)
7. `FROM t1 JOIN t2 USING (id)` — USING
8. `FROM t1 JOIN t2 ON ... JOIN t3 ON ...` — chained left-assoc
9. `FROM t1, t2 JOIN t3 ON ...` — comma + join mixed
10. `FROM (SELECT 1 AS x) AS sub` — subquery in FROM
11. `FROM LATERAL (SELECT ...)` — LATERAL subquery
12. `FROM TABLE(FLATTEN(v:arr))` — table function (basic — may need tokIdent check for FLATTEN)
13. `FROM t1 ASOF JOIN t2 MATCH_CONDITION (t1.ts >= t2.ts)` — Snowflake ASOF
14. `FROM t1 DIRECTED INNER JOIN t2 ON ...` — DIRECTED hint
15. Error cases: JOIN without ON, missing alias for subquery

## Out of Scope

| Feature | Where |
|---|---|
| PIVOT / UNPIVOT | T5.3 |
| SAMPLE / TABLESAMPLE | T5.3 |
| MATCH_RECOGNIZE | T5.3 |
| CHANGES clause | T5.3 |
| AT / BEFORE (time travel) | T5.3 |

## Acceptance Criteria

1. `go build ./snowflake/...` succeeds
2. `go vet ./snowflake/...` clean
3. `gofmt -l snowflake/` clean
4. `go test ./snowflake/...` passes
5. `Parse("SELECT * FROM t1 JOIN t2 ON t1.id = t2.id;")` returns a SelectStmt with `From[0]` being a `*JoinExpr`
6. Existing T1.4 SELECT tests still pass (updated for `From []Node` type change)
7. Walker regen correct
8. After merge, dag.md T1.5 → done
