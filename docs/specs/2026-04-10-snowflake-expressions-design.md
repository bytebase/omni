# Snowflake Expressions (T1.3) — Design

**DAG node:** T1.3 — expressions
**Branch:** `feat/snowflake/expressions`
**Depends on:** T1.2 (data types), T1.1 (identifiers).
**Unblocks:** T1.4 (SELECT core), T1.5 (joins), T2.x (DDL), T5.x (DML).

## Purpose

T1.3 adds the expression parser — the core of any SQL parser. Every SQL clause that contains a value (SELECT list, WHERE, GROUP BY, HAVING, ORDER BY, CAST arguments, function arguments, DEFAULT values, CHECK constraints) calls `parseExpr()`.

T1.3 ships:
1. **22 AST node types + 5 helper types + 5 enums** in `snowflake/ast`
2. **A Pratt parser** (`parseExpr` + ~25 helper functions) in `snowflake/parser/expr.go`
3. **Comprehensive tests** (~800 LOC) in `snowflake/parser/expr_test.go`

## Architecture — Pratt Parsing

Mirrors `mysql/parser/expr.go` (1,932 LOC). The Pratt pattern:

```
parseExpr()
  └── parseExprPrec(minBP)
        ├── parsePrefixExpr()           // prefix: unary -, +, NOT, primary
        │     └── parsePrimaryExpr()    // literals, idents, parens, CASE, CAST, IFF, etc.
        └── infix loop                  // binary ops + postfix + special forms
              ├── infixBindingPower()   // maps token → (bp, BinaryOp)
              ├── BinaryExpr           // regular binary
              └── special handlers      // IS, IN, BETWEEN, LIKE, ::, [], :, ., OVER, COLLATE
```

### Binding power table

```go
const (
    bpNone       = 0
    bpOr         = 10  // OR
    bpAnd        = 20  // AND
    bpNot        = 30  // NOT (prefix)
    bpComparison = 40  // =, <>, <, >, <=, >=, != + IS, IN, BETWEEN, LIKE
    bpAdd        = 50  // +, -, ||
    bpMul        = 60  // *, /, %
    bpUnary      = 70  // unary -, +
    bpPostfix    = 80  // ::, COLLATE, OVER
    bpAccess     = 90  // :, [], .
)
```

## AST Types

### Enums

```go
type BinaryOp int
const (
    BinAdd BinaryOp = iota // +
    BinSub                  // -
    BinMul                  // *
    BinDiv                  // /
    BinMod                  // %
    BinConcat               // ||
    BinEq                   // =
    BinNe                   // <> or !=
    BinLt                   // <
    BinGt                   // >
    BinLe                   // <=
    BinGe                   // >=
    BinAnd                  // AND
    BinOr                   // OR
)

type UnaryOp int
const ( UnaryMinus UnaryOp = iota; UnaryPlus; UnaryNot )

type LiteralKind int
const ( LitNull LiteralKind = iota; LitBool; LitInt; LitFloat; LitString )

type LikeOp int
const ( LikeOpLike LikeOp = iota; LikeOpILike; LikeOpRLike; LikeOpRegexp )

type AccessKind int
const ( AccessColon AccessKind = iota; AccessBracket; AccessDot )
```

### Expression Nodes (22 types)

Each is a Node with `Tag() NodeTag` and `Loc Loc`.

**Core (11):**

| Node | Fields | Grammar alternative |
|---|---|---|
| `Literal` | Kind LiteralKind, Value string, Ival int64, Bval bool | `literal` rule |
| `ColumnRef` | Parts []Ident | `full_column_name` (1-4 parts) |
| `StarExpr` | Qualifier *ObjectName (optional) | `id_ '.' STAR` or bare `*` |
| `BinaryExpr` | Op BinaryOp, Left/Right Node | arithmetic/comparison/logical/concat |
| `UnaryExpr` | Op UnaryOp, Expr Node | prefix `-, +, NOT` |
| `ParenExpr` | Expr Node | `bracket_expression` (non-subquery) |
| `CastExpr` | Expr Node, TypeName *TypeName, TryCast bool, ColonColon bool | `cast_expr`, `try_cast_expr`, `expr :: data_type` |
| `CaseExpr` | Operand Node (nil for searched), Whens []*WhenClause, Else Node | `case_expression` |
| `FuncCallExpr` | Name ObjectName, Args []Node, Star bool, Distinct bool, OrderBy []*OrderItem, Over *WindowSpec | `function_call`, `aggregate_function`, `ranking_windowed_function`, `trim_expression` |
| `IffExpr` | Cond/Then/Else Node | `iff_expr` |
| `CollateExpr` | Expr Node, Collation string | `expr COLLATE string` |

**Predicates (5):**

| Node | Fields | Grammar alternative |
|---|---|---|
| `IsExpr` | Expr Node, Not bool, Null bool, DistinctFrom Node | `expr IS [NOT] NULL` / `expr IS [NOT] DISTINCT FROM expr` |
| `BetweenExpr` | Expr/Low/High Node, Not bool | `expr [NOT] BETWEEN low AND high` |
| `InExpr` | Expr Node, Values []Node, Not bool | `expr [NOT] IN (expr_list)` |
| `LikeExpr` | Expr/Pattern Node, Escape Node, Op LikeOp, Not bool, Any bool, AnyValues []Node | `expr [NOT] LIKE/ILIKE/RLIKE pattern` + `LIKE ANY (...)` |
| `ExistsExpr` | Query Node | placeholder for T1.4 |

**Semi-structured + special (5):**

| Node | Fields | Grammar alternative |
|---|---|---|
| `AccessExpr` | Expr Node, Kind AccessKind, Field Ident (colon/dot), Index Node (bracket) | `expr:field`, `expr[idx]`, `expr.field` |
| `ArrayLiteralExpr` | Elements []Node | `arr_literal` |
| `JsonLiteralExpr` | Pairs []KeyValuePair | `json_literal` |
| `LambdaExpr` | Params []Ident, Body Node | `lambda_params '->' expr` |
| `SubqueryExpr` | Query Node | placeholder for T1.4 |

**Window (1):**

| Node | Fields |
|---|---|
| `WindowSpec` | PartitionBy []Node, OrderBy []*OrderItem, Frame *WindowFrame |

### Helper Types (5, NOT Nodes)

```go
type WhenClause struct { Cond, Result Node; Loc Loc }
type KeyValuePair struct { Key string; Value Node; Loc Loc }
type OrderItem struct { Expr Node; Desc bool; NullsFirst *bool; Loc Loc }
type WindowFrame struct { Kind WindowFrameKind; Start, End WindowBound; Loc Loc }
type WindowBound struct { Kind WindowBoundKind; Offset Node; Loc Loc }
```

Window frame enums:

```go
type WindowFrameKind int
const ( FrameRows WindowFrameKind = iota; FrameRange; FrameGroups )

type WindowBoundKind int
const (
    BoundUnboundedPreceding WindowBoundKind = iota
    BoundPreceding      // N PRECEDING
    BoundCurrentRow
    BoundFollowing      // N FOLLOWING
    BoundUnboundedFollowing
)
```

### NodeTag additions

17 new tags: `T_Literal`, `T_ColumnRef`, `T_StarExpr`, `T_BinaryExpr`, `T_UnaryExpr`, `T_ParenExpr`, `T_CastExpr`, `T_CaseExpr`, `T_FuncCallExpr`, `T_IffExpr`, `T_CollateExpr`, `T_IsExpr`, `T_BetweenExpr`, `T_InExpr`, `T_LikeExpr`, `T_AccessExpr`, `T_ArrayLiteralExpr`, `T_JsonLiteralExpr`, `T_LambdaExpr`.

(`ExistsExpr`, `SubqueryExpr`, `WindowSpec` are either T1.4 stubs or helper types without their own tags.)

## Parser Functions

All in `snowflake/parser/expr.go`:

### Core Pratt loop

```go
func (p *Parser) parseExpr() (ast.Node, error)                    // entry: calls parseExprPrec(bpNone + 1)
func (p *Parser) parseExprPrec(minBP int) (ast.Node, error)        // Pratt loop
func (p *Parser) parsePrefixExpr() (ast.Node, error)               // unary + primary dispatch
func (p *Parser) parsePrimaryExpr() (ast.Node, error)              // big switch: literals, idents, parens, CASE, CAST, IFF, etc.
func (p *Parser) infixBindingPower() (int, BinaryOp, bool)         // maps current token to (bp, op, ok)
```

### Primary expression helpers

```go
func (p *Parser) parseLiteral() (*ast.Literal, error)
func (p *Parser) parseIdentExpr() (ast.Node, error)                // ident → ColumnRef or FuncCallExpr (if followed by '(')
func (p *Parser) parseFuncCall(name ast.ObjectName) (*ast.FuncCallExpr, error)
func (p *Parser) parseParenExpr() (ast.Node, error)                // (expr) or (expr_list) for IN
func (p *Parser) parseCaseExpr() (*ast.CaseExpr, error)
func (p *Parser) parseCastExpr() (*ast.CastExpr, error)
func (p *Parser) parseTryCastExpr() (*ast.CastExpr, error)
func (p *Parser) parseIffExpr() (*ast.IffExpr, error)
func (p *Parser) parseArrayLiteral() (*ast.ArrayLiteralExpr, error)
func (p *Parser) parseJsonLiteral() (*ast.JsonLiteralExpr, error)
func (p *Parser) parseLambdaExpr() (*ast.LambdaExpr, error)        // heuristic: ident -> or ( ident, ident ) ->
```

### Infix/postfix special handlers

```go
func (p *Parser) parseIsExpr(left ast.Node) (ast.Node, error)
func (p *Parser) parseBetweenExpr(left ast.Node) (ast.Node, error)
func (p *Parser) parseInExpr(left ast.Node) (ast.Node, error)
func (p *Parser) parseLikeExpr(left ast.Node) (ast.Node, error)
func (p *Parser) parseCastPostfix(left ast.Node) (*ast.CastExpr, error)     // :: postfix
func (p *Parser) parseAccessExpr(left ast.Node) (*ast.AccessExpr, error)    // :, [], .
func (p *Parser) parseOverClause() (*ast.WindowSpec, error)
func (p *Parser) parseWindowFrame() (*ast.WindowFrame, error)
func (p *Parser) parseWindowBound() (ast.WindowBound, error)
```

### Utilities

```go
func (p *Parser) parseExprList() ([]ast.Node, error)               // comma-separated expression list
func (p *Parser) parseOrderByList() ([]*ast.OrderItem, error)       // for window ORDER BY and WITHIN GROUP
```

### Public helper

```go
func ParseExpr(input string) (ast.Node, []ParseError)              // freestanding for tests + callers
```

## Subquery Handling (T1.4 placeholder)

When `parsePrimaryExpr` encounters `(` followed by `kwSELECT` or `kwWITH`, it returns a `ParseError` "subquery expressions not yet supported". When `parseExpr` encounters `kwEXISTS` followed by `(`, same error.

T1.4 will fill in the real implementations by:
1. Adding a `kwSELECT` / `kwWITH` branch inside `parseParenExpr`
2. Adding an `kwEXISTS` branch inside `parsePrimaryExpr`

The Pratt framework does NOT need modification — only the primary-expression dispatch gains new cases.

## Lambda Detection Heuristic

Snowflake's `lambda_params '->' expr` is syntactically ambiguous: `x -> x + 1` starts with an identifier, which could also be a column reference. The parser must lookahead to distinguish.

**Heuristic** (same approach as many SQL parsers):
- If the current token is an identifier AND the next is `->`, parse as lambda
- If the current token is `(` AND after matching identifiers and `)` the next is `->`, parse as lambda
- Otherwise, parse as a normal expression

This is handled inside `parsePrimaryExpr`'s identifier branch: peek for `->` before falling through to parseIdentExpr.

## File Layout

| File | Change | Approx LOC |
|------|--------|-----------|
| `snowflake/ast/parsenodes.go` | MODIFY: append 22 node types + 5 helpers + 5 enums | +400 |
| `snowflake/ast/nodetags.go` | MODIFY: +17 T_* tags + String() cases | +40 |
| `snowflake/ast/loc.go` | MODIFY: +17 NodeLoc cases | +35 |
| `snowflake/ast/walk_generated.go` | REGENERATED: many new cases | ~auto |
| `snowflake/parser/expr.go` | NEW: Pratt parser + all helpers | 900 |
| `snowflake/parser/expr_test.go` | NEW: comprehensive tests | 800 |

**Estimated total: ~2,200 LOC** across 2 new + 4 modified files.

## Testing

Table-driven with `testParseExpr(input) (ast.Node, error)`.

### Test categories (22)

1. Literals: `42`, `3.14`, `'hello'`, `TRUE`, `FALSE`, `NULL`
2. Column refs: `col`, `t.col`, `s.t.col`, `db.s.t.col`
3. Star: `*`, `t.*`
4. Arithmetic precedence: `1 + 2 * 3` → BinaryExpr(+, 1, BinaryExpr(*, 2, 3))
5. Parenthesized: `(1 + 2) * 3`
6. Comparison: `a = b`, `a <> b`, `a != b`, `a >= b`
7. Logical: `a AND b OR c` → BinaryExpr(OR, BinaryExpr(AND, a, b), c)
8. NOT prefix: `NOT a`, `NOT NOT a`
9. Unary: `-a`, `+a`
10. Concat: `a || b || c`
11. CAST: `CAST(x AS INT)`, `TRY_CAST(x AS VARCHAR(100))`, `x::INT`
12. CASE simple: `CASE x WHEN 1 THEN 'a' WHEN 2 THEN 'b' ELSE 'c' END`
13. CASE searched: `CASE WHEN x > 0 THEN 'pos' ELSE 'neg' END`
14. IFF: `IFF(x > 0, 'pos', 'neg')`
15. IS NULL: `x IS NULL`, `x IS NOT NULL`, `x IS DISTINCT FROM y`, `x IS NOT DISTINCT FROM y`
16. BETWEEN: `x BETWEEN 1 AND 10`, `x NOT BETWEEN 1 AND 10`
17. IN: `x IN (1, 2, 3)`, `x NOT IN (1, 2, 3)`
18. LIKE: `x LIKE 'pat%'`, `x ILIKE '%pat%'`, `x RLIKE '.*'`, `x LIKE '%' ESCAPE '\'`, `x LIKE ANY ('a%', 'b%')`
19. Function calls: `f()`, `f(1, 2)`, `COUNT(*)`, `COUNT(DISTINCT x)`, `TRIM(x)`
20. Window: `ROW_NUMBER() OVER (PARTITION BY a ORDER BY b)`, `SUM(x) OVER (ORDER BY y ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)`
21. JSON/Array/Access: `v:field`, `v[0]`, `v:field.subfield`, `[1, 2, 3]`, `{'key': 'val'}`
22. Lambda: `x -> x + 1`, `(x, y) -> x + y`
23. Nested: `CAST(IFF(a > b, a, b) AS INT)`, `f(g(x))`, `a + b * c :: INT`
24. Error cases: unclosed paren, missing THEN in CASE, unterminated BETWEEN

## Out of Scope

| Feature | Where |
|---|---|
| Subquery expressions `(SELECT ...)` | T1.4 |
| EXISTS (subquery) | T1.4 |
| IN (subquery) | T1.4 |
| FLATTEN / LATERAL | T5.3 |
| Aggregate WITHIN GROUP (ORDER BY) | Included in FuncCallExpr via OrderBy field |
| LISTAGG / ARRAY_AGG special forms | Handled as FuncCallExpr with OrderBy |

## Acceptance Criteria

1. `go build ./snowflake/...` succeeds
2. `go vet ./snowflake/...` clean
3. `gofmt -l snowflake/` clean
4. `go test ./snowflake/...` passes
5. `go generate ./snowflake/ast/...` produces correct walker with expression node child walks
6. Every expression form from the legacy `expr` rule (25 alternatives) has a corresponding parse path
7. Precedence is correct: `1 + 2 * 3` parses as `1 + (2 * 3)`, not `(1 + 2) * 3`
8. Window frames parse correctly: `ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW`
9. After merge, `docs/migration/snowflake/dag.md` T1.3 status → `done`
