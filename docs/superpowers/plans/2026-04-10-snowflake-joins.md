# Snowflake Joins (T1.5) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend T1.4's FROM clause from simple comma-separated table refs to full JOIN syntax — JoinExpr Node, extended TableRef (subquery/func/lateral), 6 join types, ON/USING/MATCH_CONDITION conditions, left-associative join chain parsing.

**Architecture:** `parseFromClause()` returns `[]ast.Node` (was `[]*ast.TableRef`). Each comma-separated item is parsed by `parseFromItem()` which calls `parsePrimarySource()` for the base table/subquery/function, then `parseJoinChain()` for any following JOINs. JoinExpr nodes form a left-associative binary tree.

**Tech Stack:** Go 1.25, stdlib only (`strings` for MATCH_CONDITION ident check).

**Spec:** `docs/superpowers/specs/2026-04-10-snowflake-joins-design.md` (commit `6cb51b9`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/joins` on branch `feat/snowflake/joins`

**Commit policy:** No commits during implementation. User reviews the full diff at the end.

---

## File Structure

| File | Change | Approx LOC |
|------|--------|-----------|
| `snowflake/ast/parsenodes.go` | MODIFY: extend TableRef + add JoinExpr + JoinType | +60 |
| `snowflake/ast/nodetags.go` | MODIFY: +T_JoinExpr + String() | +3 |
| `snowflake/ast/loc.go` | MODIFY: +*JoinExpr case | +2 |
| `snowflake/ast/walk_generated.go` | REGENERATED | auto |
| `snowflake/parser/select.go` | MODIFY: rewrite parseFromClause + add join functions | +250 |
| `snowflake/parser/select_test.go` | MODIFY: update From assertions + add 15 join tests | +300 |

Total: ~600 LOC across 6 modified files.

---

## Task 1: AST types + nodetag + walker regen

**Files:** Modify parsenodes.go, nodetags.go, loc.go. Regenerate walk_generated.go.

- [ ] **Step 1: Extend TableRef and add JoinExpr + JoinType to parsenodes.go**

Extend the existing `TableRef` struct with 3 new fields and add JoinExpr + JoinType after it:

```go
// Updated TableRef — add Subquery, FuncCall, Lateral fields
type TableRef struct {
    Name     *ObjectName    // table name; nil for subquery/func sources
    Alias    Ident          // AS alias; zero if absent
    Subquery Node           // (SELECT ...) in FROM; nil for table refs
    FuncCall *FuncCallExpr  // TABLE(func(...)); nil for table refs
    Lateral  bool           // LATERAL prefix
    Loc      Loc
}

// JoinExpr represents a JOIN between two FROM sources.
type JoinExpr struct {
    Type           JoinType
    Left           Node    // TableRef or nested JoinExpr
    Right          Node    // TableRef or nested JoinExpr
    On             Node    // ON condition; nil for CROSS/NATURAL/USING-only
    Using          []Ident // USING columns; nil if ON or NATURAL
    Natural        bool
    Directed       bool    // Snowflake DIRECTED hint
    MatchCondition Node    // ASOF MATCH_CONDITION(expr); nil for non-ASOF
    Loc            Loc
}

func (n *JoinExpr) Tag() NodeTag { return T_JoinExpr }
var _ Node = (*JoinExpr)(nil)

type JoinType int
const (
    JoinInner JoinType = iota
    JoinLeft
    JoinRight
    JoinFull
    JoinCross
    JoinAsof
)
```

- [ ] **Step 2: Change SelectStmt.From type**

In the SelectStmt struct, change `From []*TableRef` to `From []Node`.

- [ ] **Step 3: Add T_JoinExpr to nodetags.go + NodeLoc case to loc.go**

- [ ] **Step 4: Regenerate walker**

Run: `go generate ./snowflake/ast/...`

- [ ] **Step 5: Verify AST builds**

Run: `go build ./snowflake/ast/...`

Note: `go build ./snowflake/parser/...` will FAIL because select.go still returns `[]*ast.TableRef` from parseFromClause — Task 2 fixes that.

---

## Task 2: Rewrite parseFromClause + add join parsing functions

**Files:** Modify `snowflake/parser/select.go`

- [ ] **Step 1: Rewrite parseFromClause + add all join functions**

The implementer must:

1. **Read the spec** for the full function list and flow
2. **Read the existing `parseFromClause` and `parseTableRef`** in select.go to understand what to replace
3. **Rewrite `parseFromClause()`** to return `[]ast.Node` and call `parseFromItem()` for each comma-separated item
4. **Add `parseFromItem()`** — calls parsePrimarySource + parseJoinChain
5. **Add `parsePrimarySource()`** — dispatches on `(` (subquery/parens), `kwLATERAL`, `kwTABLE`, otherwise existing parseTableRef
6. **Add `parseJoinChain(left)`** — left-associative loop consuming JOIN keywords + right source + conditions
7. **Add `parseJoinKeywords()`** — recognizes all JOIN keyword combinations and returns (JoinType, natural, directed, ok)
8. **Keep existing `parseTableRef()`** for simple ObjectName + alias cases (it's still called by parsePrimarySource)

Key patterns:

**parseJoinKeywords** must handle these token sequences:
- `kwJOIN` → JoinInner
- `kwINNER kwJOIN` → JoinInner
- `kwLEFT [kwOUTER] kwJOIN` → JoinLeft
- `kwRIGHT [kwOUTER] kwJOIN` → JoinRight
- `kwFULL [kwOUTER] kwJOIN` → JoinFull
- `kwCROSS kwJOIN` → JoinCross
- `kwNATURAL [kwLEFT|kwRIGHT|kwFULL] kwJOIN` → sets natural=true
- `kwASOF kwJOIN` → JoinAsof (check if kwASOF exists; may need tokIdent check)

**ASOF** — check if `kwASOF` exists in F2. If not, use `tokIdent && strings.ToUpper(p.cur.Str) == "ASOF"`.

**DIRECTED** — similarly check kwDIRECTED. If not in F2, use tokIdent check.

**MATCH_CONDITION** — use tokIdent check (`strings.ToUpper(p.cur.Str) == "MATCH_CONDITION"`).

**Subquery in FROM** — in parsePrimarySource, when `(` is current:
- Peek for kwSELECT/kwWITH → parse as subquery, expect `)`, parseOptionalAlias, wrap in TableRef{Subquery: ...}
- Otherwise → parse as parenthesized from-item `(t1 JOIN t2 ON ...)`, expect `)`

**TABLE() function** — when kwTABLE is current:
- Consume kwTABLE, expect `(`, parse function call expression, expect `)`, parseOptionalAlias
- Wrap in TableRef{FuncCall: ...}

**LATERAL** — when kwLATERAL is current:
- Consume kwLATERAL, parse the next source (subquery or table function), set Lateral=true
- Special case: `LATERAL FLATTEN(...)` — FLATTEN is a function call, parse as FuncCallExpr

- [ ] **Step 2: Update select_test.go From assertions**

The existing T1.4 tests that check `From []*ast.TableRef` will break because the type is now `[]ast.Node`. Update them to type-assert `From[i].(*ast.TableRef)`.

- [ ] **Step 3: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: clean.

Run: `go test ./snowflake/...`
Expected: all existing tests pass with the updated assertions.

---

## Task 3: Write join tests

**Files:** Modify `snowflake/parser/select_test.go`

- [ ] **Step 1: Add 15 join test categories**

Add join-specific tests using the existing `testParseSelectStmt` helper. Each test parses a full SELECT statement and asserts the From structure:

1. Basic INNER JOIN — From[0] is *JoinExpr{Type: JoinInner, Left: *TableRef, Right: *TableRef, On: not nil}
2. LEFT JOIN
3. RIGHT OUTER JOIN
4. FULL JOIN
5. CROSS JOIN (no ON)
6. NATURAL JOIN (no ON/USING)
7. USING clause
8. Chained JOINs — From[0] is JoinExpr{Left: JoinExpr{...}, Right: TableRef}
9. Comma + JOIN mixed — From has 2 items: TableRef and JoinExpr
10. Subquery in FROM — From[0] is *TableRef with Subquery set
11. LATERAL subquery
12. TABLE(FLATTEN(...)) — TableRef with FuncCall
13. ASOF JOIN with MATCH_CONDITION
14. DIRECTED JOIN
15. Error cases

- [ ] **Step 2: Run all tests**

Run: `go test ./snowflake/...`
Expected: all pass.

---

## Task 4: Final acceptance sweep

- [ ] **Step 1:** `go build ./snowflake/...` — clean
- [ ] **Step 2:** `go vet ./snowflake/...` — clean
- [ ] **Step 3:** `gofmt -l snowflake/` — clean
- [ ] **Step 4:** `go test ./snowflake/...` — pass
- [ ] **Step 5:** Walker regen byte-identical
- [ ] **Step 6:** STOP and present for review

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| JoinExpr Node + JoinType enum | Task 1 |
| Extended TableRef (Subquery/FuncCall/Lateral) | Task 1 |
| SelectStmt.From type change | Task 1 |
| T_JoinExpr tag | Task 1 |
| parseFromClause rewrite | Task 2 |
| parseFromItem / parsePrimarySource / parseJoinChain / parseJoinKeywords | Task 2 |
| Subquery-in-FROM | Task 2 |
| TABLE() / LATERAL / FLATTEN | Task 2 |
| ASOF / DIRECTED | Task 2 |
| ON / USING / MATCH_CONDITION | Task 2 |
| Existing test updates | Task 2 |
| 15 join test categories | Task 3 |
| Acceptance criteria | Task 4 |
