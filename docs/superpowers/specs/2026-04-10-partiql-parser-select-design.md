# PartiQL Parser-Select — Design Spec

**DAG node:** 5 (parser-select), P0
**Depends on:** parser-foundation (node 4, done)
**Unblocks:** parser-let-pivot (12), parse-entry (8)
**Package:** `partiql/parser`
**Files added:** `select.go`, `from.go`
**Files modified:** `expr.go` (replace SELECT/UNION/INTERSECT/EXCEPT stubs), `exprprimary.go` (SubLink wrapping in parseParenExpr), `parser_test.go`

---

## 1. Goal

Replace the SELECT and UNION/INTERSECT/EXCEPT deferred-feature stubs from parser-foundation with real implementations. After this node, the parser handles the full SFW (SELECT-FROM-WHERE) query shape, set operations, joins, ordering, and grouping. The AWS corpus smoke test should see ~33 files fully parse (those starting with SELECT).

## 2. Scope

### In scope
- `selectClause` (4 variants: SELECT *, SELECT items, SELECT VALUE, PIVOT stub)
- `fromClause` + `tableReference` (tableBaseReference with aliases/AT/BY, JOINs, UNPIVOT, comma-joins, parenthesized table refs)
- `whereClauseSelect`
- `groupClause` (GROUP [PARTIAL] BY with aliases)
- `havingClause`
- `orderByClause` (ASC/DESC, NULLS FIRST/LAST)
- `limitClause` / `offsetByClause`
- `exprBagOp` (UNION/INTERSECT/EXCEPT with OUTER and DISTINCT/ALL modifiers)
- SubLink wrapping: `(SELECT ...)` in expression position wraps in ast.SubLink

### Out of scope (deferred)
- `letClause` (LET bindings) — deferred to parser-let-pivot (node 12)
- `PIVOT expr AT expr` in selectClause — deferred to parser-let-pivot (node 12)
- Graph match patterns in FROM — deferred to parser-graph (node 16)

## 3. Architecture

### 3.1 File layout

| File | Responsibility | Est. lines |
|------|---------------|------------|
| `select.go` | parseSFWQuery (the full SELECT...FROM...WHERE...GROUP...HAVING...ORDER...LIMIT...OFFSET shape), parseSelectClause, parseProjectionItem, parseGroupClause, parseHavingClause, parseOrderByClause, parseLimitClause, parseOffsetClause | ~350 |
| `from.go` | parseFromClause, parseTableReference (join left-recursion loop), parseTableNonJoin, parseTableBaseReference (with alias/AT/BY), parseTableUnpivot, parseJoinType, parseJoinSpec | ~250 |

### 3.2 Dispatch integration

**expr.go changes:**
- `parseSelectExpr`: replace the `if p.cur.Type == tokSELECT` stub with a call to `parseSFWQuery` in select.go. Add `tokPIVOT` detection for the PIVOT form.
- `parseBagOp`: replace the UNION/INTERSECT/EXCEPT stubs with a real left-associative loop that parses set operations with optional OUTER and DISTINCT/ALL modifiers, building `ast.SetOpStmt`.

**exprprimary.go changes:**
- `parseParenExpr`: after parsing the inner expression, if the result is a SelectStmt (detected via type assertion), wrap it in `ast.SubLink` before returning.

**Note on ParseStatement:**
Node 5 does NOT introduce `ParseStatement`. SELECT stays entirely within the expression ladder: `parseSelectExpr` dispatches to `parseSFWQuery` when the current token is `tokSELECT` (or `tokPIVOT`). When `ParseStatement` is eventually called (introduced by node 7, parser-ddl), it dispatches DQL through `ParseExpr`, which descends through the expression ladder and reaches `parseSelectExpr` normally. Node 5's changes are purely within `expr.go`, `select.go`, and `from.go`.

### 3.3 AST nodes produced

All exist in `partiql/ast/stmts.go` and `partiql/ast/tableexprs.go`:
- `SelectStmt` — the SFW query with TargetList, From, Where, GroupBy, Having, OrderBy, Limit, Offset fields
- `SetOpStmt` — UNION/INTERSECT/EXCEPT with Op, Quantifier, Outer, Left, Right
- `TargetEntry` — one projection item (Expr + optional alias)
- `OrderByItem` — one sort key (Expr + Dir + NullsOrder)
- `GroupByItem` — one group key (Expr + optional alias)
- `AliasedSource` — FROM source with As/At/By aliases
- `JoinExpr` — join with Kind, Left, Right, On
- `UnpivotExpr` — UNPIVOT source with aliases

### 3.4 Stub replacement audit

Foundation stubs being replaced:
- `parseSelectExpr` SELECT stub (expr.go) — replaced with parseSFWQuery call
- `parseBagOp` UNION/INTERSECT/EXCEPT stubs (expr.go) — replaced with real set-op loop

New stubs added by this node:
- `tokLET` after FROM → "LET is deferred to parser-let-pivot (DAG node 12)"
- `tokPIVOT` in selectClause → "PIVOT is deferred to parser-let-pivot (DAG node 12)"

## 4. Testing

- `testdata/parser-select/` — ~25 golden pairs covering all SFW clauses, join types, set operations, subquery wrapping
- `TestParser_AWSCorpus` — expect ~33 files to fully parse (those starting with SELECT)
- Error tests added to `TestParser_Errors` for new syntax errors (e.g., missing FROM after SELECT items)

## 5. Design decisions

- **D1** LET/PIVOT deferred to node 12 per DAG
- **D2** Join parsing uses iterative left-recursion (same pattern as binary operators in expr.go)
- **D3** `parseTableBaseReference` uses `exprSelect` for the source (grammar line 403-404), which enables path expressions and subqueries as FROM sources
- **D4** SelectStmt.SelectValue bool field distinguishes SELECT VALUE from SELECT items
- **D5** SetOpStmt carries the source spelling (UNION vs INTERSECT vs EXCEPT) for round-trip fidelity
