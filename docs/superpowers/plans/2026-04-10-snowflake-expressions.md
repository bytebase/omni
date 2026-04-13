# Snowflake Expressions (T1.3) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Pratt expression parser for `snowflake/parser` — 22 AST node types + 5 helpers + 5 enums in `snowflake/ast`, plus a complete `parseExpr` with ~25 parse functions in `snowflake/parser/expr.go`, covering every expression form from the legacy grammar's `expr` rule (25 alternatives).

**Architecture:** Pratt parsing (precedence climbing) mirroring `mysql/parser/expr.go`. The `parseExprPrec(minBP)` loop calls `parsePrefixExpr` for the left-hand side then iterates infix/postfix operators. Special forms (IS, IN, BETWEEN, LIKE, ::, OVER, access) are dispatched from the infix loop to dedicated parse helpers. Window frame specification (ROWS/RANGE/GROUPS BETWEEN ... AND ...) is included.

**Tech Stack:** Go 1.25, stdlib only (`strings` for case comparisons in multi-word lookahead).

**Spec:** `docs/superpowers/specs/2026-04-10-snowflake-expressions-design.md` (commit `a4ab008`)

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/expressions` on branch `feat/snowflake/expressions`

**Commit policy:** No commits during implementation. The user reviews the full diff at the end.

---

## File Structure

### Modified

| File | Changes |
|------|---------|
| `snowflake/ast/parsenodes.go` | Append 5 enums + 22 node structs + 5 helper structs + all Tag/String methods (~400 LOC) |
| `snowflake/ast/nodetags.go` | Append 17 T_* constants + String() cases |
| `snowflake/ast/loc.go` | Append 17 NodeLoc cases |
| `snowflake/ast/walk_generated.go` | Regenerated — many new cases with child walks |

### Created

| File | Purpose | Approx LOC |
|------|---------|-----------|
| `snowflake/parser/expr.go` | Pratt parser: parseExpr, all helpers, ParseExpr freestanding | 900 |
| `snowflake/parser/expr_test.go` | 24 test categories | 800 |

**Total: ~2,200 LOC** across 2 new + 4 modified files.

---

## Task 1: Add all AST enum types and expression node structs

This is the largest single task — it adds all the type infrastructure to `snowflake/ast/parsenodes.go`. The node types are pure struct definitions with no complex logic.

**Files:**
- Modify: `snowflake/ast/parsenodes.go`

- [ ] **Step 1: Confirm worktree state**

Run: `pwd && git rev-parse --abbrev-ref HEAD`
Expected:
```
/Users/h3n4l/OpenSource/omni/.worktrees/expressions
feat/snowflake/expressions
```

- [ ] **Step 2: Append all expression enums + node types to parsenodes.go**

Use Edit to append the following after the existing TypeName section at the end of `snowflake/ast/parsenodes.go`. This is a large block — the full expression AST:

```go

// ---------------------------------------------------------------------------
// Expression enums
// ---------------------------------------------------------------------------

// BinaryOp enumerates binary operator types.
type BinaryOp int

const (
	BinAdd    BinaryOp = iota // +
	BinSub                    // -
	BinMul                    // *
	BinDiv                    // /
	BinMod                    // %
	BinConcat                 // ||
	BinEq                     // =
	BinNe                     // <> or !=
	BinLt                     // <
	BinGt                     // >
	BinLe                     // <=
	BinGe                     // >=
	BinAnd                    // AND
	BinOr                     // OR
)

// String returns the operator symbol.
func (op BinaryOp) String() string {
	switch op {
	case BinAdd:
		return "+"
	case BinSub:
		return "-"
	case BinMul:
		return "*"
	case BinDiv:
		return "/"
	case BinMod:
		return "%"
	case BinConcat:
		return "||"
	case BinEq:
		return "="
	case BinNe:
		return "!="
	case BinLt:
		return "<"
	case BinGt:
		return ">"
	case BinLe:
		return "<="
	case BinGe:
		return ">="
	case BinAnd:
		return "AND"
	case BinOr:
		return "OR"
	default:
		return "?"
	}
}

// UnaryOp enumerates unary operator types.
type UnaryOp int

const (
	UnaryMinus UnaryOp = iota // -
	UnaryPlus                 // +
	UnaryNot                  // NOT
)

// String returns the operator symbol.
func (op UnaryOp) String() string {
	switch op {
	case UnaryMinus:
		return "-"
	case UnaryPlus:
		return "+"
	case UnaryNot:
		return "NOT"
	default:
		return "?"
	}
}

// LiteralKind classifies literal value types.
type LiteralKind int

const (
	LitNull   LiteralKind = iota
	LitBool
	LitInt
	LitFloat
	LitString
)

// LikeOp enumerates LIKE operator variants.
type LikeOp int

const (
	LikeOpLike   LikeOp = iota // LIKE
	LikeOpILike                // ILIKE
	LikeOpRLike                // RLIKE
	LikeOpRegexp               // REGEXP
)

// AccessKind enumerates semi-structured access operator types.
type AccessKind int

const (
	AccessColon   AccessKind = iota // expr:field (JSON path)
	AccessBracket                   // expr[idx] (array subscript)
	AccessDot                       // expr.field (dot path)
)

// WindowFrameKind enumerates window frame types.
type WindowFrameKind int

const (
	FrameRows   WindowFrameKind = iota
	FrameRange
	FrameGroups
)

// WindowBoundKind enumerates window frame bound types.
type WindowBoundKind int

const (
	BoundUnboundedPreceding WindowBoundKind = iota
	BoundPreceding
	BoundCurrentRow
	BoundFollowing
	BoundUnboundedFollowing
)

// CaseKind discriminates simple vs searched CASE expressions.
type CaseKind int

const (
	CaseSimple  CaseKind = iota // CASE expr WHEN val THEN ...
	CaseSearched                // CASE WHEN cond THEN ...
)

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

// Literal represents an integer, float, string, boolean, or NULL literal.
type Literal struct {
	Kind  LiteralKind
	Value string // raw source text for all kinds
	Ival  int64  // for LitInt
	Bval  bool   // for LitBool
	Loc   Loc
}

func (n *Literal) Tag() NodeTag { return T_Literal }

// ColumnRef represents a 1-4 part column reference (col, t.col, s.t.col, db.s.t.col).
type ColumnRef struct {
	Parts []Ident
	Loc   Loc
}

func (n *ColumnRef) Tag() NodeTag { return T_ColumnRef }

// StarExpr represents * or qualifier.* in SELECT lists and COUNT(*).
type StarExpr struct {
	Qualifier *ObjectName // optional table qualifier; nil for bare *
	Loc       Loc
}

func (n *StarExpr) Tag() NodeTag { return T_StarExpr }

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Loc   Loc
}

func (n *BinaryExpr) Tag() NodeTag { return T_BinaryExpr }

// UnaryExpr represents a prefix unary operation: op expr.
type UnaryExpr struct {
	Op   UnaryOp
	Expr Node
	Loc  Loc
}

func (n *UnaryExpr) Tag() NodeTag { return T_UnaryExpr }

// ParenExpr represents a parenthesized expression: ( expr ).
type ParenExpr struct {
	Expr Node
	Loc  Loc
}

func (n *ParenExpr) Tag() NodeTag { return T_ParenExpr }

// CastExpr represents CAST(expr AS type), TRY_CAST(expr AS type), or expr::type.
type CastExpr struct {
	Expr       Node
	TypeName   *TypeName
	TryCast    bool // true for TRY_CAST
	ColonColon bool // true for expr::type syntax
	Loc        Loc
}

func (n *CastExpr) Tag() NodeTag { return T_CastExpr }

// CaseExpr represents a CASE expression (simple or searched).
type CaseExpr struct {
	Kind    CaseKind
	Operand Node          // non-nil for CaseSimple; nil for CaseSearched
	Whens   []*WhenClause // one or more WHEN clauses
	Else    Node          // optional ELSE; nil if absent
	Loc     Loc
}

func (n *CaseExpr) Tag() NodeTag { return T_CaseExpr }

// WhenClause is a single WHEN...THEN branch inside a CaseExpr.
// Not a top-level Node — embedded in CaseExpr.
type WhenClause struct {
	Cond   Node // the WHEN condition (or value for simple CASE)
	Result Node // the THEN result
	Loc    Loc
}

// FuncCallExpr represents a function call including aggregates and window functions.
//
// Star is true for COUNT(*). Distinct is true for COUNT(DISTINCT x).
// OrderBy is used by aggregates with WITHIN GROUP (ORDER BY ...).
// Over is non-nil for window functions (SUM(x) OVER (...)).
type FuncCallExpr struct {
	Name     ObjectName   // function name (may be qualified: schema.func)
	Args     []Node       // argument expressions
	Star     bool         // COUNT(*)
	Distinct bool         // DISTINCT keyword in aggregate
	OrderBy  []*OrderItem // for WITHIN GROUP (ORDER BY ...)
	Over     *WindowSpec  // non-nil for window functions
	Loc      Loc
}

func (n *FuncCallExpr) Tag() NodeTag { return T_FuncCallExpr }

// IffExpr represents Snowflake's IFF(condition, then_expr, else_expr).
type IffExpr struct {
	Cond Node
	Then Node
	Else Node
	Loc  Loc
}

func (n *IffExpr) Tag() NodeTag { return T_IffExpr }

// CollateExpr represents expr COLLATE collation_name.
type CollateExpr struct {
	Expr      Node
	Collation string
	Loc       Loc
}

func (n *CollateExpr) Tag() NodeTag { return T_CollateExpr }

// IsExpr represents expr IS [NOT] NULL or expr IS [NOT] DISTINCT FROM expr.
type IsExpr struct {
	Expr         Node
	Not          bool // IS NOT NULL / IS NOT DISTINCT FROM
	Null         bool // true for IS [NOT] NULL; false for IS [NOT] DISTINCT FROM
	DistinctFrom Node // non-nil for IS [NOT] DISTINCT FROM expr
	Loc          Loc
}

func (n *IsExpr) Tag() NodeTag { return T_IsExpr }

// BetweenExpr represents expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Expr Node
	Low  Node
	High Node
	Not  bool
	Loc  Loc
}

func (n *BetweenExpr) Tag() NodeTag { return T_BetweenExpr }

// InExpr represents expr [NOT] IN (value_list).
// Subquery form (expr IN (SELECT ...)) is handled by T1.4.
type InExpr struct {
	Expr   Node
	Values []Node
	Not    bool
	Loc    Loc
}

func (n *InExpr) Tag() NodeTag { return T_InExpr }

// LikeExpr represents expr [NOT] LIKE/ILIKE/RLIKE/REGEXP pattern [ESCAPE esc].
// Also handles LIKE ANY (...) via the Any + AnyValues fields.
type LikeExpr struct {
	Expr      Node
	Pattern   Node   // the pattern expression
	Escape    Node   // optional ESCAPE clause; nil if absent
	Op        LikeOp // LIKE, ILIKE, RLIKE, REGEXP
	Not       bool
	Any       bool   // true for LIKE ANY (...)
	AnyValues []Node // values for LIKE ANY; nil unless Any is true
	Loc       Loc
}

func (n *LikeExpr) Tag() NodeTag { return T_LikeExpr }

// AccessExpr represents semi-structured data access:
//   - AccessColon:   expr:field (JSON path)
//   - AccessBracket: expr[index] (array subscript)
//   - AccessDot:     expr.field (dot path chaining)
type AccessExpr struct {
	Expr  Node
	Kind  AccessKind
	Field Ident // for Colon and Dot access
	Index Node  // for Bracket access
	Loc   Loc
}

func (n *AccessExpr) Tag() NodeTag { return T_AccessExpr }

// ArrayLiteralExpr represents an array literal: [elem1, elem2, ...].
type ArrayLiteralExpr struct {
	Elements []Node
	Loc      Loc
}

func (n *ArrayLiteralExpr) Tag() NodeTag { return T_ArrayLiteralExpr }

// JsonLiteralExpr represents a JSON object literal: {key: value, ...}.
type JsonLiteralExpr struct {
	Pairs []KeyValuePair
	Loc   Loc
}

func (n *JsonLiteralExpr) Tag() NodeTag { return T_JsonLiteralExpr }

// KeyValuePair is a single key-value pair in a JsonLiteralExpr.
type KeyValuePair struct {
	Key   string
	Value Node
	Loc   Loc
}

// LambdaExpr represents a lambda expression: param -> body or (p1, p2) -> body.
type LambdaExpr struct {
	Params []Ident
	Body   Node
	Loc    Loc
}

func (n *LambdaExpr) Tag() NodeTag { return T_LambdaExpr }

// SubqueryExpr is a placeholder for subquery expressions (SELECT ...).
// T1.4 will fill in the real implementation.
type SubqueryExpr struct {
	Loc Loc
}

func (n *SubqueryExpr) Tag() NodeTag { return T_SubqueryExpr }

// ExistsExpr is a placeholder for EXISTS (SELECT ...).
// T1.4 will fill in the real implementation.
type ExistsExpr struct {
	Loc Loc
}

func (n *ExistsExpr) Tag() NodeTag { return T_ExistsExpr }

// WindowSpec describes a window specification for OVER (...).
// Not a top-level Node — embedded in FuncCallExpr.
type WindowSpec struct {
	PartitionBy []Node
	OrderBy     []*OrderItem
	Frame       *WindowFrame // nil if no frame clause
	Loc         Loc
}

// OrderItem represents one element in an ORDER BY clause.
type OrderItem struct {
	Expr       Node
	Desc       bool  // true for DESC
	NullsFirst *bool // nil = unspecified, true = NULLS FIRST, false = NULLS LAST
	Loc        Loc
}

// WindowFrame represents ROWS/RANGE/GROUPS BETWEEN start AND end.
type WindowFrame struct {
	Kind  WindowFrameKind
	Start WindowBound
	End   WindowBound // End.Kind may be zero if single-bound form (not BETWEEN)
	Loc   Loc
}

// WindowBound represents one end of a window frame specification.
type WindowBound struct {
	Kind   WindowBoundKind
	Offset Node // for BoundPreceding/BoundFollowing: the N in "N PRECEDING"; nil otherwise
}

// Compile-time assertions.
var (
	_ Node = (*Literal)(nil)
	_ Node = (*ColumnRef)(nil)
	_ Node = (*StarExpr)(nil)
	_ Node = (*BinaryExpr)(nil)
	_ Node = (*UnaryExpr)(nil)
	_ Node = (*ParenExpr)(nil)
	_ Node = (*CastExpr)(nil)
	_ Node = (*CaseExpr)(nil)
	_ Node = (*FuncCallExpr)(nil)
	_ Node = (*IffExpr)(nil)
	_ Node = (*CollateExpr)(nil)
	_ Node = (*IsExpr)(nil)
	_ Node = (*BetweenExpr)(nil)
	_ Node = (*InExpr)(nil)
	_ Node = (*LikeExpr)(nil)
	_ Node = (*AccessExpr)(nil)
	_ Node = (*ArrayLiteralExpr)(nil)
	_ Node = (*JsonLiteralExpr)(nil)
	_ Node = (*LambdaExpr)(nil)
	_ Node = (*SubqueryExpr)(nil)
	_ Node = (*ExistsExpr)(nil)
)
```

- [ ] **Step 3: Verify compile**

Run: `go build ./snowflake/ast/...`
Expected: no output, exit 0. (The tag values are NOT yet defined — the file will fail because T_Literal etc. don't exist. This is expected. Move to Task 2 immediately.)

Actually — the compile WILL fail because the Tag() methods reference undefined T_* constants. This is expected. Proceed to Task 2 which adds them. Do NOT halt on this expected failure.

---

## Task 2: Add NodeTags + NodeLoc + regenerate walker

**Files:**
- Modify: `snowflake/ast/nodetags.go`
- Modify: `snowflake/ast/loc.go`
- Regenerate: `snowflake/ast/walk_generated.go`

- [ ] **Step 1: Add 17 new T_* constants to nodetags.go**

Append after `T_TypeName`:

```go
	T_Literal
	T_ColumnRef
	T_StarExpr
	T_BinaryExpr
	T_UnaryExpr
	T_ParenExpr
	T_CastExpr
	T_CaseExpr
	T_FuncCallExpr
	T_IffExpr
	T_CollateExpr
	T_IsExpr
	T_BetweenExpr
	T_InExpr
	T_LikeExpr
	T_AccessExpr
	T_ArrayLiteralExpr
	T_JsonLiteralExpr
	T_LambdaExpr
	T_SubqueryExpr
	T_ExistsExpr
```

Add matching cases to `String()`:

```go
	case T_Literal:
		return "Literal"
	case T_ColumnRef:
		return "ColumnRef"
	case T_StarExpr:
		return "StarExpr"
	case T_BinaryExpr:
		return "BinaryExpr"
	case T_UnaryExpr:
		return "UnaryExpr"
	case T_ParenExpr:
		return "ParenExpr"
	case T_CastExpr:
		return "CastExpr"
	case T_CaseExpr:
		return "CaseExpr"
	case T_FuncCallExpr:
		return "FuncCallExpr"
	case T_IffExpr:
		return "IffExpr"
	case T_CollateExpr:
		return "CollateExpr"
	case T_IsExpr:
		return "IsExpr"
	case T_BetweenExpr:
		return "BetweenExpr"
	case T_InExpr:
		return "InExpr"
	case T_LikeExpr:
		return "LikeExpr"
	case T_AccessExpr:
		return "AccessExpr"
	case T_ArrayLiteralExpr:
		return "ArrayLiteralExpr"
	case T_JsonLiteralExpr:
		return "JsonLiteralExpr"
	case T_LambdaExpr:
		return "LambdaExpr"
	case T_SubqueryExpr:
		return "SubqueryExpr"
	case T_ExistsExpr:
		return "ExistsExpr"
```

- [ ] **Step 2: Add 21 NodeLoc cases to loc.go**

Add before the `default:` in `NodeLoc`:

```go
	case *Literal:
		return v.Loc
	case *ColumnRef:
		return v.Loc
	case *StarExpr:
		return v.Loc
	case *BinaryExpr:
		return v.Loc
	case *UnaryExpr:
		return v.Loc
	case *ParenExpr:
		return v.Loc
	case *CastExpr:
		return v.Loc
	case *CaseExpr:
		return v.Loc
	case *FuncCallExpr:
		return v.Loc
	case *IffExpr:
		return v.Loc
	case *CollateExpr:
		return v.Loc
	case *IsExpr:
		return v.Loc
	case *BetweenExpr:
		return v.Loc
	case *InExpr:
		return v.Loc
	case *LikeExpr:
		return v.Loc
	case *AccessExpr:
		return v.Loc
	case *ArrayLiteralExpr:
		return v.Loc
	case *JsonLiteralExpr:
		return v.Loc
	case *LambdaExpr:
		return v.Loc
	case *SubqueryExpr:
		return v.Loc
	case *ExistsExpr:
		return v.Loc
```

- [ ] **Step 3: Regenerate walker**

Run: `go generate ./snowflake/ast/...`
Expected: `Generated walk_generated.go: N cases, M child fields` where N is significantly more than 2 (the previous count). The exact numbers depend on how many Node-typed fields the generator finds across all structs.

- [ ] **Step 4: Verify build + tests**

Run: `go build ./snowflake/... && go vet ./snowflake/... && go test ./snowflake/ast/...`
Expected: all clean, all existing tests pass.

---

## Task 3: Write the Pratt parser core (expr.go)

This is the central task — the Pratt parser skeleton. It creates `snowflake/parser/expr.go` with the binding power table, the core Pratt loop, prefix/primary dispatch, and the infix binary-operator handler. Special forms (IS, IN, BETWEEN, LIKE, CASE, CAST, etc.) are added in subsequent tasks as stubs initially.

**Files:**
- Create: `snowflake/parser/expr.go`

- [ ] **Step 1: Write the full expr.go file**

Create `snowflake/parser/expr.go` with the complete Pratt parser. This is the biggest single code file in the migration — ~900 LOC.

The implementer should read the spec and the mysql reference (`mysql/parser/expr.go`), then write a Snowflake-adapted version. The file must include:

1. **Binding power constants** (bpNone through bpAccess)
2. **`parseExpr()`** — entry point, calls `parseExprPrec(bpNone + 1)`
3. **`parseExprPrec(minBP)`** — Pratt loop: parsePrefixExpr → infix loop
4. **`parsePrefixExpr()`** — dispatches unary -, +, NOT to `parseUnaryExpr`; delegates everything else to `parsePrimaryExpr`
5. **`parsePrimaryExpr()`** — big switch on `p.cur.Type` dispatching to: literals (tokInt/tokFloat/tokString/kwTRUE/kwFALSE/kwNULL), identifiers → `parseIdentExpr`, `(` → `parseParenExpr`, kwCASE → `parseCaseExpr`, kwCAST → `parseCastExpr`, kwTRY_CAST → `parseTryCastExpr`, kwIFF → `parseIffExpr`, `[` → `parseArrayLiteral`, `{` → `parseJsonLiteral`, `*` → StarExpr, and lambda detection (ident + `->`)
6. **`infixBindingPower()`** — maps current token to (bp, BinaryOp, ok) for standard binary operators; returns ok=true for special forms (IS, IN, BETWEEN, LIKE, ::, :, [, ., OVER, COLLATE) at their respective binding powers
7. **Infix loop handlers** — inside `parseExprPrec`'s loop: regular BinaryExpr for standard ops; dedicated handlers for IS, IN, BETWEEN, LIKE/ILIKE/RLIKE, :: (cast postfix), : (JSON access), [ ] (array access), . (dot access), OVER (window), COLLATE
8. **All helper parse functions**: `parseUnaryExpr`, `parseLiteral`, `parseIdentExpr`, `parseFuncCall`, `parseParenExpr`, `parseCaseExpr`, `parseCastExpr`, `parseTryCastExpr`, `parseIffExpr`, `parseArrayLiteral`, `parseJsonLiteral`, `parseLambdaOrIdent`, `parseIsExpr`, `parseBetweenExpr`, `parseInExpr`, `parseLikeExpr`, `parseCastPostfix`, `parseAccessExpr`, `parseOverClause`, `parseWindowFrame`, `parseWindowBound`, `parseExprList`, `parseOrderByList`, `ParseExpr` (freestanding)

**Key patterns from the spec:**

- The Pratt loop consumes operators greedily: `for { prec, op, ok := p.infixBindingPower(); if !ok || prec < minBP { break } ... }`
- For standard binary ops: consume the operator, parse the right side with `parseExprPrec(prec + 1)` (left-associative), build a `BinaryExpr`
- For postfix ops (::, COLLATE, OVER): consume the operator and its arguments, wrap left in the appropriate node, continue the loop
- For access ops (:, [], .): highest binding power, parse the field/index, wrap left in AccessExpr
- For IS/IN/BETWEEN/LIKE: same precedence as comparison (bpComparison), dedicated parse helpers that consume the full syntax

**Import:** `"strings"` (for multi-word lookups like `strings.ToUpper(p.cur.Str) == "VARYING"` — though expressions don't need this often; the main use is for checking `DISTINCT FROM` after IS NOT).

Also import `"github.com/bytebase/omni/snowflake/ast"`.

- [ ] **Step 2: Verify build**

Run: `go build ./snowflake/parser/...`
Expected: no output, exit 0. (No tests yet — just compile check.)

Run: `go vet ./snowflake/parser/...`
Expected: no output, exit 0.

---

## Task 4: Write expression tests (expr_test.go)

**Files:**
- Create: `snowflake/parser/expr_test.go`

- [ ] **Step 1: Write the test file**

Create `snowflake/parser/expr_test.go` with table-driven tests using a `testParseExpr` helper. The test file should cover all 24 categories from the spec:

1. Literals
2. Column refs
3. Star
4. Arithmetic precedence
5. Parenthesized
6. Comparison operators
7. Logical operators (AND/OR/NOT precedence)
8. Unary
9. Concat
10. CAST / TRY_CAST / ::
11. CASE simple + searched
12. IFF
13. IS NULL / IS NOT NULL / IS DISTINCT FROM / IS NOT DISTINCT FROM
14. BETWEEN / NOT BETWEEN
15. IN / NOT IN
16. LIKE / ILIKE / RLIKE / LIKE ANY
17. Function calls (generic + COUNT(*) + COUNT(DISTINCT) + TRIM)
18. Window functions with OVER (PARTITION BY + ORDER BY + frame)
19. JSON access (:field), array access ([idx]), dot access (.field)
20. Array literal / JSON literal
21. Lambda expressions
22. Nested expressions
23. COLLATE
24. Error cases

Each test should assert the correct AST node type via type assertion and verify key fields (Op, Kind, Loc positions where critical).

The test helper:

```go
func testParseExpr(input string) (ast.Node, error) {
    p := &Parser{lexer: NewLexer(input), input: input}
    p.advance()
    return p.parseExpr()
}
```

- [ ] **Step 2: Run tests**

Run: `go test ./snowflake/...`
Expected: all tests pass.

---

## Task 5: Final acceptance criteria sweep

- [ ] **Step 1: Build**

Run: `go build ./snowflake/...`
Expected: no output, exit 0.

- [ ] **Step 2: Vet**

Run: `go vet ./snowflake/...`
Expected: no output, exit 0.

- [ ] **Step 3: Gofmt**

Run: `gofmt -l snowflake/`
Expected: no output. If any file listed, `gofmt -w snowflake/` and re-check.

- [ ] **Step 4: Test**

Run: `go test ./snowflake/...`
Expected: both ast and parser pass.

- [ ] **Step 5: Walker generation round-trip**

Run: `cp snowflake/ast/walk_generated.go /tmp/wg.go && go generate ./snowflake/ast/... && diff /tmp/wg.go snowflake/ast/walk_generated.go && echo "BYTE_IDENTICAL" && rm /tmp/wg.go`
Expected: `BYTE_IDENTICAL`.

- [ ] **Step 6: STOP and present for review**

Do NOT commit.

---

## Spec Coverage Checklist

| Spec section | Covered by |
|---|---|
| 5 enum types (BinaryOp, UnaryOp, LiteralKind, LikeOp, AccessKind) | Task 1 |
| Window enums (WindowFrameKind, WindowBoundKind, CaseKind) | Task 1 |
| 22 expression node types + 5 helpers | Task 1 |
| 17 NodeTag constants | Task 2 |
| 17 NodeLoc cases | Task 2 |
| Walker regeneration | Task 2 |
| Pratt parser (parseExpr, parseExprPrec, infixBindingPower) | Task 3 |
| Prefix dispatch (parsePrefixExpr, parseUnaryExpr) | Task 3 |
| Primary dispatch (parsePrimaryExpr) | Task 3 |
| Special form helpers (CASE, CAST, IFF, IS, IN, BETWEEN, LIKE) | Task 3 |
| Function calls + aggregates + window | Task 3 |
| Access operators (:, [], .) | Task 3 |
| Lambda expressions | Task 3 |
| Array/JSON literals | Task 3 |
| ParseExpr freestanding | Task 3 |
| 24 test categories | Task 4 |
| Acceptance criteria | Task 5 |

## Implementation Notes for Subagents

**Task 3 is the critical task.** The implementer must:

1. Read `mysql/parser/expr.go` (1932 LOC) as the reference implementation
2. Read the spec's binding power table and node inventory
3. Write a Snowflake-adapted Pratt parser that handles all 25 expr alternatives from the legacy grammar
4. The parser must compile cleanly and integrate with the existing Parser struct from F4

**Key differences from mysql's expr.go:**
- Snowflake has `::` cast postfix (mysql doesn't)
- Snowflake has JSON `:` access (mysql uses `->` / `->>`; snowflake uses `:`)
- Snowflake has IFF() (mysql uses IF())
- Snowflake has ILIKE (mysql doesn't)
- Snowflake has LIKE ANY (...) (mysql doesn't)
- Snowflake has lambda expressions (mysql doesn't)
- Snowflake has array/JSON literals in expression position (mysql doesn't)
- Snowflake does NOT have: `:=` assignment, `^` XOR, `&` bitwise AND, `|` bitwise OR, `<<` `>>` shifts, `<=>` null-safe equal, `SOUNDS LIKE`, `MATCH ... AGAINST`, `CONVERT`, `INTERVAL`

**The Pratt parser structure:**
```
parseExpr → parseExprPrec(1)
  parsePrefixExpr:
    case '-': unary minus at bpUnary
    case '+': skip (unary plus is a no-op)
    case kwNOT: unary NOT at bpNot
    default: parsePrimaryExpr
  parsePrimaryExpr:
    case tokInt/tokFloat/tokString: Literal
    case kwTRUE/kwFALSE: Literal(LitBool)
    case kwNULL: Literal(LitNull)
    case tokIdent/tokQuotedIdent/keyword-as-ident: parseIdentExpr → ColumnRef or FuncCallExpr
    case '(': parseParenExpr (check for SELECT → error placeholder)
    case kwCASE: parseCaseExpr
    case kwCAST: parseCastExpr
    case kwTRY_CAST: parseTryCastExpr
    case kwIFF: parseIffExpr
    case '[': parseArrayLiteral
    case '{': parseJsonLiteral
    case '*': StarExpr
    case kwEXISTS: error "subquery not yet supported"
  infix loop:
    binary ops: +, -, *, /, %, =, <>, <, >, <=, >=, !=, AND, OR, ||
    special: IS (bpComparison), IN (bpComparison), BETWEEN (bpComparison),
             LIKE/ILIKE/RLIKE (bpComparison), :: (bpPostfix),
             COLLATE (bpPostfix), OVER (bpPostfix),
             : (bpAccess), [ (bpAccess), . (bpAccess)
```
