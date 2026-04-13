# Snowflake SELECT Core (T1.4) — Design

**DAG node:** T1.4 — SELECT core
**Branch:** `feat/snowflake/select-core`
**Depends on:** T1.3 (expressions), T1.2 (data types), T1.1 (identifiers).
**Unblocks:** T1.5 (joins), T1.7 (set operators). Also effectively completes T1.6 (CTEs) since basic CTE parsing is included.

## Purpose

T1.4 replaces F4's `unsupported("SELECT")` and `unsupported("WITH")` dispatch stubs with real parsers. For the first time, `Parse("SELECT ...")` returns a real `*SelectStmt` AST node instead of an error. T1.4 also fills in the T1.3 subquery placeholders (SubqueryExpr, ExistsExpr, IN-subquery) since SELECT parsing enables them.

## Scope

**Included:**
- SELECT list (expressions + aliases + star + EXCLUDE)
- DISTINCT / ALL
- TOP n
- FROM clause (comma-separated table references with aliases)
- WHERE clause
- GROUP BY (normal, CUBE, ROLLUP, GROUPING SETS, ALL)
- HAVING clause
- QUALIFY clause (Snowflake-specific)
- ORDER BY (reuses T1.3's parseOrderByList)
- LIMIT / OFFSET
- FETCH FIRST/NEXT n ROWS ONLY
- Basic CTEs: WITH [RECURSIVE] name [(columns)] AS (SELECT ...) (merges T1.6 scope)
- Subquery expressions: (SELECT ...), EXISTS (SELECT ...), expr IN (SELECT ...)

**Deferred:**
- JOIN syntax (T1.5) — FROM only handles comma-separated ObjectName table refs
- Set operators UNION/EXCEPT/MINUS/INTERSECT (T1.7)
- Subquery in FROM position (`SELECT * FROM (SELECT ...) AS sub`) — requires T1.5's richer table-source model
- INTO clause (Snowflake Scripting)

## AST Types

### SelectStmt (Node)

```go
type SelectStmt struct {
    With     []*CTE
    Distinct bool
    All      bool
    Top      Node            // TOP n; nil if absent
    Targets  []*SelectTarget // SELECT list
    From     []*TableRef     // FROM; nil if absent
    Where    Node            // WHERE; nil if absent
    GroupBy  *GroupByClause  // GROUP BY; nil if absent
    Having   Node            // HAVING; nil if absent
    Qualify  Node            // QUALIFY; nil if absent (Snowflake-specific)
    OrderBy  []*OrderItem    // ORDER BY; nil if absent (reuses T1.3 type)
    Limit    Node            // LIMIT n; nil if absent
    Offset   Node            // OFFSET n; nil if absent
    Fetch    *FetchClause    // FETCH FIRST/NEXT; nil if absent
    Loc      Loc
}
func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }
```

SelectStmt is the ONLY new Node type. All helper types below are non-Nodes (no Tag method). The walker's generated case for `*SelectStmt` traverses its Node-typed fields (Top, Where, Having, Qualify, Limit, Offset) automatically. Slice fields containing non-Node helpers (Targets, From, With, GroupBy, OrderBy) are not auto-traversed by the walker — bytebase's query_span_extractor can manually iterate them when needed.

### Helper types (5, NOT Nodes)

```go
type SelectTarget struct {
    Expr    Node    // the expression; nil for bare star
    Alias   Ident   // AS alias; zero if absent
    Star    bool    // true for * or qualifier.*
    Exclude []Ident // EXCLUDE columns; nil if absent
    Loc     Loc
}

type TableRef struct {
    Name  *ObjectName
    Alias Ident       // zero if absent
    Loc   Loc
}

type CTE struct {
    Name      Ident
    Columns   []Ident   // optional column aliases
    Query     Node       // *SelectStmt
    Recursive bool
    Loc       Loc
}

type GroupByClause struct {
    Kind  GroupByKind
    Items []Node
    Loc   Loc
}

type GroupByKind int
const (
    GroupByNormal GroupByKind = iota
    GroupByCube
    GroupByRollup
    GroupByGroupingSets
    GroupByAll
)

type FetchClause struct {
    Count Node
    Loc   Loc
}
```

### NodeTag addition

One new tag: `T_SelectStmt`.

## Parser Changes

### Dispatch stub replacements in `parser.go`

```go
case kwSELECT:
    return p.parseSelectStmt()   // replaces p.unsupported("SELECT")
case kwWITH:
    return p.parseWithSelect()   // replaces p.unsupported("WITH")
```

This is the FIRST time a dispatch stub gets a real implementation. The pattern for future Tier 1+/2+ nodes: replace `p.unsupported("X")` with `p.parseXStmt()`.

### Subquery placeholder fills in `expr.go`

1. **`parseParenExpr()`** — when `(` is followed by `kwSELECT` or `kwWITH`, parse a SelectStmt via `p.parseSelectStmt()` (or `p.parseWithSelect()`), expect `)`, wrap in `&SubqueryExpr{Query: stmt}`.

2. **`parsePrefixExpr()`** — when `kwEXISTS` is current, consume it, expect `(`, parse SelectStmt, expect `)`, wrap in `&ExistsExpr{Query: stmt}`.

3. **`parseInExpr()`** — when the `(` is followed by `kwSELECT` or `kwWITH`, parse as SubqueryExpr inside InExpr's Values (single-element containing the SubqueryExpr).

### New file: `snowflake/parser/select.go`

```go
func (p *Parser) parseSelectStmt() (*ast.SelectStmt, error)
func (p *Parser) parseWithSelect() (ast.Node, error)
func (p *Parser) parseCTEList(recursive bool) ([]*ast.CTE, error)
func (p *Parser) parseCTE(recursive bool) (*ast.CTE, error)
func (p *Parser) parseSelectBody() (*ast.SelectStmt, error)  // SELECT keyword already consumed
func (p *Parser) parseSelectList() ([]*ast.SelectTarget, error)
func (p *Parser) parseSelectTarget() (*ast.SelectTarget, error)
func (p *Parser) parseFromClause() ([]*ast.TableRef, error)
func (p *Parser) parseTableRef() (*ast.TableRef, error)
func (p *Parser) parseGroupByClause() (*ast.GroupByClause, error)
func (p *Parser) parseLimitOffsetFetch(stmt *ast.SelectStmt) error
func (p *Parser) parseOptionalAlias() (ast.Ident, bool)
```

### parseSelectStmt flow

```
parseSelectStmt:
    consume kwSELECT
    check DISTINCT / ALL
    check TOP n
    parseSelectList → targets
    if kwFROM: parseFromClause → from
    if kwWHERE: parseExpr → where
    if kwGROUP: parseGroupByClause → groupby
    if kwHAVING: parseExpr → having
    if kwQUALIFY: parseExpr → qualify
    if kwORDER: parseOrderByList → orderby  (reuses T1.3)
    parseLimitOffsetFetch → limit/offset/fetch
    return &SelectStmt{...}
```

### parseWithSelect flow

```
parseWithSelect:
    consume kwWITH
    check kwRECURSIVE
    parseCTEList → ctes
    parseSelectStmt → stmt
    stmt.With = ctes
    return stmt
```

### parseOptionalAlias

Handles `[AS] alias` where AS is optional:
- If current is `kwAS`, consume it, then `parseIdent()` → alias
- If current is an identifier or non-reserved keyword (but NOT a clause keyword like FROM/WHERE/GROUP/etc.), consume as alias
- Otherwise, no alias

This is the trickiest helper — it needs to distinguish `SELECT a b FROM t` (where `b` is the alias) from `SELECT a FROM t` (where `FROM` is NOT an alias). The lookahead checks whether the current token is a clause-starting keyword.

## File Layout

| File | Change | Approx LOC |
|------|--------|-----------|
| `snowflake/ast/parsenodes.go` | MODIFY: +SelectStmt, SelectTarget, TableRef, CTE, GroupByClause, GroupByKind, FetchClause | +100 |
| `snowflake/ast/nodetags.go` | MODIFY: +T_SelectStmt | +3 |
| `snowflake/ast/loc.go` | MODIFY: +*SelectStmt case | +2 |
| `snowflake/ast/walk_generated.go` | REGENERATED | auto |
| `snowflake/parser/parser.go` | MODIFY: replace SELECT/WITH stubs | ~4 lines changed |
| `snowflake/parser/expr.go` | MODIFY: fill SubqueryExpr/ExistsExpr/IN-subquery | +50 |
| `snowflake/parser/select.go` | NEW: SELECT parser | 500 |
| `snowflake/parser/select_test.go` | NEW: SELECT tests | 600 |

**Estimated total: ~1,250 LOC** across 2 new + 5 modified files.

## Testing

Table-driven with `testParseSelect(input) (*ast.SelectStmt, error)` helper.

### Categories (22)

1. `SELECT 1` — simplest
2. `SELECT 1, 2, 3` — multiple targets
3. `SELECT a AS x, b AS y FROM t` — aliases (with AS)
4. `SELECT a x FROM t` — alias without AS
5. `SELECT * FROM t` — star
6. `SELECT * EXCLUDE (a, b) FROM t` — EXCLUDE
7. `SELECT t.* FROM t` — qualified star
8. `SELECT DISTINCT a FROM t`
9. `SELECT TOP 10 a FROM t`
10. `SELECT a FROM t1, t2 AS x` — comma FROM with alias
11. `SELECT a FROM t WHERE a > 0`
12. `SELECT a, COUNT(*) FROM t GROUP BY a`
13. `SELECT a FROM t GROUP BY CUBE (a, b)` + ROLLUP, GROUPING SETS, ALL
14. `SELECT a FROM t GROUP BY a HAVING COUNT(*) > 1`
15. `SELECT a FROM t QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1`
16. `SELECT a FROM t ORDER BY a DESC NULLS LAST`
17. `SELECT a FROM t LIMIT 10 OFFSET 5`
18. `SELECT a FROM t FETCH FIRST 10 ROWS ONLY`
19. `WITH cte AS (SELECT 1 AS x) SELECT * FROM cte` — basic CTE
20. `SELECT (SELECT 1)` — SubqueryExpr
21. `SELECT * FROM t WHERE EXISTS (SELECT 1 FROM t2)` — ExistsExpr
22. `SELECT * FROM t WHERE a IN (SELECT b FROM t2)` — IN-subquery

Run scope: `go test ./snowflake/...`

## Out of Scope

| Feature | Where |
|---|---|
| JOIN syntax | T1.5 |
| Set operators (UNION/EXCEPT/MINUS/INTERSECT) | T1.7 |
| Subquery in FROM position | T1.5 |
| INTO clause | Snowflake Scripting |
| PIVOT/UNPIVOT in FROM | T5.3 |
| SAMPLE/TABLESAMPLE | T5.3 |

## Acceptance Criteria

1. `go build ./snowflake/...` succeeds
2. `go vet ./snowflake/...` clean
3. `gofmt -l snowflake/` clean
4. `go test ./snowflake/...` passes
5. `Parse("SELECT 1;")` returns `*SelectStmt` with one `Literal` target (NOT an unsupported error)
6. `Parse("WITH cte AS (SELECT 1) SELECT * FROM cte;")` returns a `*SelectStmt` with `With` populated
7. `ParseExpr("(SELECT 1)")` returns a `*SubqueryExpr` containing a `*SelectStmt`
8. `ParseExpr("EXISTS (SELECT 1)")` returns an `*ExistsExpr`
9. Walker regeneration produces correct output
10. After merge, dag.md T1.4 (and effectively T1.6) status → `done`
