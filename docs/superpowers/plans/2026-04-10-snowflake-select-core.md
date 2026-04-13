# Snowflake SELECT Core (T1.4) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace F4's `unsupported("SELECT")` and `unsupported("WITH")` dispatch stubs with real parsers that produce `*SelectStmt` AST nodes, and fill T1.3's subquery placeholders (SubqueryExpr, ExistsExpr, IN-subquery) so expressions can contain subqueries.

**Architecture:** `parseSelectStmt()` parses each SELECT clause in order (DISTINCT/ALL, TOP, target list, FROM, WHERE, GROUP BY, HAVING, QUALIFY, ORDER BY, LIMIT/OFFSET/FETCH). `parseWithSelect()` handles `WITH [RECURSIVE] name AS (SELECT ...) SELECT ...`. Subquery fills modify `expr.go`'s `parseParenExpr`, `parsePrefixExpr`, and `parseInExpr` to call `parseSelectStmt` when `(SELECT` or `EXISTS(SELECT` is encountered.

**Tech Stack:** Go 1.25, stdlib only.

**Spec:** `docs/superpowers/specs/2026-04-10-snowflake-select-core-design.md` (commit `a659ca6`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/select-core` on branch `feat/snowflake/select-core`

**Commit policy:** No commits during implementation. User reviews the full diff at the end.

---

## File Structure

### Modified

| File | Changes |
|------|---------|
| `snowflake/ast/parsenodes.go` | Append SelectStmt + SelectTarget + TableRef + CTE + GroupByClause + GroupByKind + FetchClause (~100 LOC) |
| `snowflake/ast/nodetags.go` | Append T_SelectStmt + String() case |
| `snowflake/ast/loc.go` | Append *SelectStmt case |
| `snowflake/ast/walk_generated.go` | Regenerated |
| `snowflake/parser/parser.go` | Replace SELECT/WITH dispatch stubs (~4 lines) |
| `snowflake/parser/expr.go` | Fill SubqueryExpr/ExistsExpr/IN-subquery placeholders (~50 LOC) |

### Created

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/parser/select.go` | SELECT parser | 500 |
| `snowflake/parser/select_test.go` | SELECT tests | 600 |

Total: ~1,250 LOC across 2 new + 6 modified files.

---

## Task 1: Add AST types + nodetag + NodeLoc + walker regen

**Files:**
- Modify: `snowflake/ast/parsenodes.go`, `snowflake/ast/nodetags.go`, `snowflake/ast/loc.go`
- Regenerate: `snowflake/ast/walk_generated.go`

- [ ] **Step 1: Confirm worktree**

Run: `pwd && git rev-parse --abbrev-ref HEAD`

- [ ] **Step 2: Append AST types to parsenodes.go**

Append after the existing expression compile-time assertions at the end of `snowflake/ast/parsenodes.go`:

```go

// ---------------------------------------------------------------------------
// Statement nodes
// ---------------------------------------------------------------------------

// SelectStmt represents a SELECT statement.
type SelectStmt struct {
	With     []*CTE          // WITH clause CTEs; nil if absent
	Distinct bool            // SELECT DISTINCT
	All      bool            // SELECT ALL
	Top      Node            // TOP n expression; nil if absent
	Targets  []*SelectTarget // SELECT list items
	From     []*TableRef     // FROM table references; nil if absent
	Where    Node            // WHERE condition; nil if absent
	GroupBy  *GroupByClause  // GROUP BY; nil if absent
	Having   Node            // HAVING condition; nil if absent
	Qualify  Node            // QUALIFY condition; nil if absent (Snowflake-specific)
	OrderBy  []*OrderItem    // ORDER BY; nil if absent
	Limit    Node            // LIMIT n; nil if absent
	Offset   Node            // OFFSET n; nil if absent
	Fetch    *FetchClause    // FETCH FIRST/NEXT; nil if absent
	Loc      Loc
}

func (n *SelectStmt) Tag() NodeTag { return T_SelectStmt }

var _ Node = (*SelectStmt)(nil)

// SelectTarget is one item in a SELECT list.
// For expressions: Expr is set, Star is false.
// For star: Star is true, Expr may be a qualifier (table.*) or nil (bare *).
type SelectTarget struct {
	Expr    Node    // expression; nil for bare *
	Alias   Ident   // AS alias; zero Ident if absent
	Star    bool    // true for * or qualifier.*
	Exclude []Ident // EXCLUDE columns; nil if absent
	Loc     Loc
}

// TableRef is a table reference in the FROM clause.
// T1.5 extends this with join syntax.
type TableRef struct {
	Name  *ObjectName // table name
	Alias Ident       // AS alias; zero if absent
	Loc   Loc
}

// CTE represents a Common Table Expression in a WITH clause.
type CTE struct {
	Name      Ident   // CTE name
	Columns   []Ident // optional column aliases
	Query     Node    // the SELECT body (*SelectStmt)
	Recursive bool    // WITH RECURSIVE flag
	Loc       Loc
}

// GroupByClause represents a GROUP BY clause with optional variant.
type GroupByClause struct {
	Kind  GroupByKind
	Items []Node // group-by expressions
	Loc   Loc
}

// GroupByKind enumerates GROUP BY variants.
type GroupByKind int

const (
	GroupByNormal      GroupByKind = iota // GROUP BY a, b
	GroupByCube                           // GROUP BY CUBE (a, b)
	GroupByRollup                         // GROUP BY ROLLUP (a, b)
	GroupByGroupingSets                   // GROUP BY GROUPING SETS ((a), (b))
	GroupByAll                            // GROUP BY ALL
)

// FetchClause represents FETCH FIRST/NEXT n ROWS ONLY.
type FetchClause struct {
	Count Node // the count expression
	Loc   Loc
}
```

- [ ] **Step 3: Add T_SelectStmt to nodetags.go + String() case**

- [ ] **Step 4: Add *SelectStmt to NodeLoc in loc.go**

- [ ] **Step 5: Regenerate walker**

Run: `go generate ./snowflake/ast/...`

- [ ] **Step 6: Verify build + tests**

Run: `go build ./snowflake/... && go vet ./snowflake/... && go test ./snowflake/ast/...`

---

## Task 2: Replace dispatch stubs + fill subquery placeholders

**Files:**
- Modify: `snowflake/parser/parser.go`
- Modify: `snowflake/parser/expr.go`

- [ ] **Step 1: Replace SELECT/WITH dispatch stubs in parser.go**

In `snowflake/parser/parser.go`, find the `parseStmt()` dispatch switch. Replace the SELECT and WITH cases:

```go
	case kwSELECT:
		return p.parseSelectStmt()
	case kwWITH:
		return p.parseWithSelect()
```

This will cause a compile error since `parseSelectStmt` and `parseWithSelect` don't exist yet. That's expected — Task 3 creates them.

- [ ] **Step 2: Fill SubqueryExpr placeholder in expr.go**

In `snowflake/parser/expr.go`, find `parseParenExpr()`. Currently it returns an error when it sees `kwSELECT` or `kwWITH` after `(`. Replace that error path with real subquery parsing:

After consuming `(`, check:
```go
if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
    var query ast.Node
    var err error
    if p.cur.Type == kwWITH {
        query, err = p.parseWithSelect()
    } else {
        query, err = p.parseSelectStmt()
    }
    if err != nil {
        return nil, err
    }
    closeTok, err := p.expect(')')
    if err != nil {
        return nil, err
    }
    return &ast.SubqueryExpr{Query: query, Loc: ast.Loc{Start: startLoc, End: closeTok.Loc.End}}, nil
}
```

Note: `SubqueryExpr` in T1.3 was defined as a placeholder with just `Loc`. Now it needs a `Query Node` field. **Add this field to the SubqueryExpr struct** in `parsenodes.go`:

```go
type SubqueryExpr struct {
    Query Node // the SELECT statement
    Loc   Loc
}
```

Similarly update `ExistsExpr` to add a `Query Node` field:

```go
type ExistsExpr struct {
    Query Node
    Loc   Loc
}
```

- [ ] **Step 3: Fill ExistsExpr placeholder in expr.go**

In `parsePrefixExpr()`, find the EXISTS handling (currently returns an error). Replace with:

```go
case kwEXISTS:
    existsTok := p.advance()
    if _, err := p.expect('('); err != nil {
        return nil, err
    }
    var query ast.Node
    var err error
    if p.cur.Type == kwWITH {
        query, err = p.parseWithSelect()
    } else {
        query, err = p.parseSelectStmt()
    }
    if err != nil {
        return nil, err
    }
    closeTok, err := p.expect(')')
    if err != nil {
        return nil, err
    }
    return &ast.ExistsExpr{Query: query, Loc: ast.Loc{Start: existsTok.Loc.Start, End: closeTok.Loc.End}}, nil
```

- [ ] **Step 4: Fill IN-subquery in parseInExpr**

In `parseInExpr()`, after consuming `(`, check if the next token is SELECT/WITH. If so, parse as SubqueryExpr:

```go
// Inside parseInExpr, after consuming '(':
if p.cur.Type == kwSELECT || p.cur.Type == kwWITH {
    // Subquery IN
    var query ast.Node
    var err error
    if p.cur.Type == kwWITH {
        query, err = p.parseWithSelect()
    } else {
        query, err = p.parseSelectStmt()
    }
    if err != nil {
        return nil, err
    }
    closeTok, err := p.expect(')')
    if err != nil {
        return nil, err
    }
    subq := &ast.SubqueryExpr{Query: query, Loc: ast.Loc{Start: openLoc, End: closeTok.Loc.End}}
    return &ast.InExpr{Expr: left, Values: []ast.Node{subq}, Not: not, Loc: ...}, nil
}
```

- [ ] **Step 5: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: compile error — `parseSelectStmt` and `parseWithSelect` are undefined. This is expected; Task 3 creates them.

---

## Task 3: Write the SELECT parser (select.go)

**Files:**
- Create: `snowflake/parser/select.go`

- [ ] **Step 1: Write the full select.go**

Create `snowflake/parser/select.go` with all parse functions. The implementer should:

1. Read the spec's "parseSelectStmt flow" and "parseWithSelect flow" sections
2. Read `mysql/parser/select.go` as the structural reference
3. Implement all functions listed in the spec

**Key functions:**

```go
// parseSelectStmt parses SELECT [DISTINCT|ALL] [TOP n] target_list
//   [FROM table_refs] [WHERE expr] [GROUP BY ...] [HAVING expr]
//   [QUALIFY expr] [ORDER BY ...] [LIMIT n [OFFSET n]] [FETCH ...]
func (p *Parser) parseSelectStmt() (*ast.SelectStmt, error)

// parseWithSelect parses WITH [RECURSIVE] cte_list SELECT ...
func (p *Parser) parseWithSelect() (ast.Node, error)

// parseCTEList parses a comma-separated list of CTEs
func (p *Parser) parseCTEList(recursive bool) ([]*ast.CTE, error)

// parseCTE parses one CTE: name [(columns)] AS (SELECT ...)
func (p *Parser) parseCTE(recursive bool) (*ast.CTE, error)

// parseSelectList parses comma-separated SELECT targets
func (p *Parser) parseSelectList() ([]*ast.SelectTarget, error)

// parseSelectTarget parses one target: expr [AS alias] or * [EXCLUDE ...]
func (p *Parser) parseSelectTarget() (*ast.SelectTarget, error)

// parseFromClause parses FROM table_ref [, table_ref ...]
func (p *Parser) parseFromClause() ([]*ast.TableRef, error)

// parseTableRef parses one table: object_name [AS alias]
func (p *Parser) parseTableRef() (*ast.TableRef, error)

// parseGroupByClause parses GROUP BY [CUBE|ROLLUP|GROUPING SETS|ALL] (exprs)
func (p *Parser) parseGroupByClause() (*ast.GroupByClause, error)

// parseLimitOffsetFetch fills stmt.Limit, stmt.Offset, stmt.Fetch
func (p *Parser) parseLimitOffsetFetch(stmt *ast.SelectStmt) error

// parseOptionalAlias returns (alias, true) if an alias follows, or (zero, false)
func (p *Parser) parseOptionalAlias() (ast.Ident, bool)
```

**parseOptionalAlias is the trickiest helper.** It must distinguish `SELECT a b FROM t` (b is alias) from `SELECT a FROM t` (FROM is clause keyword, not alias). Implementation:

```go
func (p *Parser) parseOptionalAlias() (ast.Ident, bool) {
    if p.cur.Type == kwAS {
        p.advance()
        id, err := p.parseIdent()
        if err != nil {
            return ast.Ident{}, false
        }
        return id, true
    }
    // Check if current token can be an alias (ident or non-reserved keyword
    // that is NOT a clause-starting keyword)
    if p.cur.Type == tokIdent || p.cur.Type == tokQuotedIdent ||
       (p.cur.Type >= 700 && !keywordReserved[p.cur.Type] && !isClauseKeyword(p.cur.Type)) {
        id, err := p.parseIdent()
        if err != nil {
            return ast.Ident{}, false
        }
        return id, true
    }
    return ast.Ident{}, false
}

// isClauseKeyword returns true for keywords that start SQL clauses and
// should NOT be consumed as implicit aliases.
func isClauseKeyword(t int) bool {
    switch t {
    case kwFROM, kwWHERE, kwGROUP, kwHAVING, kwQUALIFY, kwORDER,
         kwLIMIT, kwFETCH, kwUNION, kwEXCEPT, kwMINUS, kwINTERSECT,
         kwINTO, kwON, kwJOIN, kwINNER, kwLEFT, kwRIGHT, kwFULL,
         kwCROSS, kwNATURAL, kwWITH, kwSELECT, kwSET, kwWHEN,
         kwTHEN, kwELSE, kwEND, kwCASE, kwWINDOW:
        return true
    }
    return false
}
```

**parseSelectTarget** handles three forms:
1. `*` → `SelectTarget{Star: true}`
2. `qualifier.*` → detect by checking if current is an ident followed by `.` then `*`
3. `expr [AS alias]` → parse expression, check for EXCLUDE, check for alias

**parseGroupByClause** handles five variants:
- `GROUP BY expr, expr` → `GroupByNormal`
- `GROUP BY CUBE (expr, expr)` → `GroupByCube`
- `GROUP BY ROLLUP (expr, expr)` → `GroupByRollup`
- `GROUP BY GROUPING SETS ((expr), (expr))` → `GroupByGroupingSets`
- `GROUP BY ALL` → `GroupByAll`

Check if keywords `kwCUBE`, `kwROLLUP`, `kwGROUPING` exist in F2. If not, use tokIdent with string comparison.

**parseLimitOffsetFetch** handles:
- `LIMIT n [OFFSET n]`
- `[OFFSET n] FETCH FIRST|NEXT n ROWS ONLY`
- `OFFSET n` alone (without FETCH)

- [ ] **Step 2: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0.

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

- [ ] **Step 3: Smoke test**

Run a quick smoke test:
```bash
cat > /tmp/select-smoke.go <<'EOF'
package main

import (
    "fmt"
    "github.com/bytebase/omni/snowflake/parser"
    "github.com/bytebase/omni/snowflake/ast"
)

func main() {
    cases := []string{
        "SELECT 1;",
        "SELECT a, b FROM t;",
        "WITH cte AS (SELECT 1) SELECT * FROM cte;",
    }
    for _, c := range cases {
        result := parser.ParseBestEffort(c)
        fmt.Printf("%q → %d stmts, %d errors\n", c, len(result.File.Stmts), len(result.Errors))
        for _, s := range result.File.Stmts {
            fmt.Printf("  type: %s\n", s.Tag())
            if sel, ok := s.(*ast.SelectStmt); ok {
                fmt.Printf("  targets: %d, from: %d\n", len(sel.Targets), len(sel.From))
            }
        }
    }
}
EOF
go run /tmp/select-smoke.go
rm /tmp/select-smoke.go
```

Expected: each SELECT produces a `*SelectStmt` with correct target/from counts. The key signal: `"SELECT 1;" → 1 stmts, 0 errors` (NOT the old "1 errors" from the unsupported stub).

---

## Task 4: Write SELECT tests (select_test.go)

**Files:**
- Create: `snowflake/parser/select_test.go`

- [ ] **Step 1: Write select_test.go**

Create `snowflake/parser/select_test.go` with table-driven tests covering all 22 categories from the spec. Use a helper:

```go
func testParseSelectStmt(input string) (*ast.SelectStmt, []ParseError) {
    result := ParseBestEffort(input)
    if len(result.File.Stmts) == 0 {
        return nil, result.Errors
    }
    sel, ok := result.File.Stmts[0].(*ast.SelectStmt)
    if !ok {
        return nil, append(result.Errors, ParseError{Msg: "not a SelectStmt"})
    }
    return sel, result.Errors
}
```

Test categories (22):
1. `SELECT 1` — 1 target, no FROM/WHERE/etc
2. `SELECT 1, 2, 3` — 3 targets
3. `SELECT a AS x, b AS y FROM t` — aliases with AS
4. `SELECT a x FROM t` — alias without AS
5. `SELECT * FROM t` — star target
6. `SELECT * EXCLUDE (a, b) FROM t` — EXCLUDE
7. `SELECT t.* FROM t` — qualified star
8. `SELECT DISTINCT a FROM t`
9. `SELECT TOP 10 a FROM t`
10. `SELECT a FROM t1, t2 AS x` — comma FROM
11. `SELECT a FROM t WHERE a > 0`
12. `SELECT a, COUNT(*) FROM t GROUP BY a` — normal GROUP BY
13. GROUP BY CUBE / ROLLUP / GROUPING SETS / ALL variants
14. `SELECT ... HAVING COUNT(*) > 1`
15. `SELECT ... QUALIFY ROW_NUMBER() OVER (ORDER BY a) = 1`
16. `SELECT ... ORDER BY a DESC NULLS LAST`
17. `SELECT ... LIMIT 10 OFFSET 5`
18. `SELECT ... FETCH FIRST 10 ROWS ONLY`
19. `WITH cte AS (SELECT 1 AS x) SELECT * FROM cte`
20. `SELECT (SELECT 1)` — SubqueryExpr
21. `SELECT * FROM t WHERE EXISTS (SELECT 1)` — ExistsExpr
22. `SELECT * FROM t WHERE a IN (SELECT b FROM t2)` — IN-subquery

Each test should:
- Call `testParseSelectStmt(input)`
- Assert zero errors (or expected errors for error cases)
- Type-assert and verify key fields (target count, from count, distinct flag, etc.)

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/...`
Expected: all pass.

---

## Task 5: Final acceptance sweep

- [ ] **Step 1: Build** — `go build ./snowflake/...`
- [ ] **Step 2: Vet** — `go vet ./snowflake/...`
- [ ] **Step 3: Gofmt** — `gofmt -l snowflake/` (fix if needed)
- [ ] **Step 4: Test** — `go test ./snowflake/...`
- [ ] **Step 5: Walker regen** — verify byte-identical
- [ ] **Step 6: Verify Parse("SELECT 1;") no longer returns unsupported error** — critical acceptance criterion
- [ ] **Step 7: STOP and present for review**

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| SelectStmt Node + helpers + GroupByKind | Task 1 |
| T_SelectStmt tag | Task 1 |
| Replace dispatch stubs | Task 2 |
| Fill SubqueryExpr/ExistsExpr/IN-subquery | Task 2 |
| parseSelectStmt + all helpers | Task 3 |
| parseWithSelect + CTE parsing | Task 3 |
| parseOptionalAlias (clause keyword disambiguation) | Task 3 |
| 22 test categories | Task 4 |
| Acceptance criteria | Task 5 |
