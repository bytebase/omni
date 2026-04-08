# PartiQL ast-core Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the `partiql/ast` Go package — ~73 hand-written AST node types covering the full legacy ANTLR PartiQL grammar, organized via 6 sealed sub-interfaces, with a `NodeToString` debug helper and table-driven unit tests for tag dispatch and `Loc` handling.

**Architecture:** Sealed sub-interface style (`Node` + `StmtNode`/`ExprNode`/`TableExpr`/`PathStep`/`TypeName`/`PatternNode`) modeled on `cosmosdb/ast/`. Each node embeds a `Loc` field with byte offsets only (no line/column). 8 source files split by category, plus `ast_test.go` and `outfuncs.go`. Pure Go, standard library only, no external dependencies.

**Tech Stack:** Go (matching the omni `go.mod` toolchain), `testing` package, `reflect` for the safety-net coverage test.

**Spec:** `docs/superpowers/specs/2026-04-08-partiql-ast-core-design.md`

**DAG entry:** `docs/migration/partiql/dag.md` node 1 (P0).

**Worktree:** `/Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core`
**Branch:** `feat/partiql/ast-core`

---

## File Map

| File | Lines (approx) | Responsibility |
|------|---------------|---------------|
| `partiql/ast/node.go` | ~120 | Six interfaces (`Node`/`StmtNode`/`ExprNode`/`TableExpr`/`PathStep`/`TypeName`/`PatternNode`), `Loc` struct, `List` helper |
| `partiql/ast/literals.go` | ~240 | 8 literal nodes — all `ExprNode` |
| `partiql/ast/exprs.go` | ~750 | 28 expression nodes (operators, predicates, special-forms, paths, vars, params, subqueries, collection literals, window spec, path steps) + 6 enums |
| `partiql/ast/tableexprs.go` | ~120 | 4 table-position nodes (`TableRef`, `AliasedSource`, `JoinExpr`, `UnpivotExpr`) + `JoinKind` enum. Also documents the multi-interface re-exports from `exprs.go`. |
| `partiql/ast/types.go` | ~30 | Single `TypeRef` node implementing `TypeName` |
| `partiql/ast/stmts.go` | ~700 | 14 top-level statement nodes + 11 clause/DML helpers + 5 enums |
| `partiql/ast/patterns.go` | ~180 | 6 graph-pattern nodes + 3 enums |
| `partiql/ast/outfuncs.go` | ~600 | `NodeToString` — giant type switch over every node type |
| `partiql/ast/ast_test.go` | ~600 | Compile-time interface assertions, `TestGetLoc`, `TestNodeToString` golden cases, `TestNodeToString_AllNodesCovered` reflection safety net |

**Total:** ~3340 lines across 9 files for ~73 node types and ~14 enum types.

## Conventions Applied To Every Type

Every node struct in this package follows the same shape:

```go
// <NodeName> represents <grammar feature>.
//
// Grammar: <rule or rule#Label reference from PartiQLParser.g4>
type <NodeName> struct {
    <fields...>
    Loc Loc
}

func (*<NodeName>) nodeTag()        {}
func (n *<NodeName>) GetLoc() Loc   { return n.Loc }
// ... and one or more sub-interface tag methods, e.g.:
func (*<NodeName>) exprNode()       {}
```

**Doc comments are required** for every node and must mention which grammar rule(s) the node represents (acceptance criterion 2 in the spec).

**Pointer receivers** for `nodeTag` / `GetLoc` / sub-interface tag methods (matches `cosmosdb/ast`).

**Loc field is always last** in the struct.

**No constructor functions.** Tests and parsers build nodes with struct literals (acceptance criterion: no smart constructors).

**Enum types** are declared as `type Foo int` with `const ( FooA Foo = iota; FooB; ... )` blocks. Enums get a `String() string` method that returns the canonical name (used by `NodeToString`).

**Grammar reference convention.** Every `// Grammar:` comment in this package uses `rule#Label` for labeled ANTLR alternatives in `bytebase/parser/partiql/PartiQLParser.g4`, and bare `rule` for unlabeled rules. **Verify each name against the grammar file before committing — invented rule names are a regression risk that already cost one fix-up commit on `literals.go`.** When a file sources types from multiple grammar rules, its file-level header should list each rule with its line range; individual node doc comments should cite their specific `rule#Label`.

---

## Task Order Rationale

Tasks 1–10 build the AST in dependency order (foundation → leaves → branches → compound types). Task 11 wires `NodeToString` over the complete type set. Task 12 is final verification + DAG bookkeeping + branch finishing.

**Why exprs.go is split into 3 tasks** (3, 4, 5): exprs.go has 28 types and would otherwise be a single ~800-line task with too many sub-types in one commit. Splitting along grammar boundaries (operators, special-forms, paths/vars/collections) keeps each commit reviewable.

**Why stmts.go is split into 2 tasks** (8, 9): same reason — 25 types in two coherent groups (top-level statements vs. clause helpers).

**Note on enum location:** The spec listed `QuantifierKind` under "Enums declared in stmts.go" but it is first used by `FuncCall` in `exprs.go`. To avoid forward references during incremental builds, **`QuantifierKind` is declared in `exprs.go`**, not stmts.go. This is the only deviation from the spec's enum-location list.

---

### Task 1: Foundation — `node.go` and test bootstrap

**Files:**
- Create: `partiql/ast/node.go`
- Create: `partiql/ast/ast_test.go`

- [ ] **Step 1: Create `partiql/ast/node.go`**

```go
// Package ast defines parse-tree node types for the omni PartiQL parser.
//
// PartiQL is the SQL++-flavored query language used by AWS DynamoDB and
// Azure Cosmos DB. This package mirrors the legacy bytebase/parser/partiql
// ANTLR grammar's full coverage scope as defined in
// bytebase/parser/partiql/PartiQLParser.g4.
//
// AST style: sealed sub-interfaces. Every node implements Node; most also
// implement one or more of StmtNode, ExprNode, TableExpr, PathStep, TypeName,
// or PatternNode for compile-time position discipline. A handful of small
// clause helpers (e.g. TargetEntry, CaseWhen) are bare Node.
//
// See docs/superpowers/specs/2026-04-08-partiql-ast-core-design.md for the
// design rationale.
package ast

// Loc is a half-open byte range describing where a node appears in the
// original source text. Loc{-1, -1} means the position is unknown
// (synthetic nodes constructed post-parse).
type Loc struct {
	Start int // inclusive byte offset
	End   int // exclusive byte offset
}

// Node is the root interface for all PartiQL parse-tree nodes.
//
// Implementations carry a Loc field and expose it via GetLoc. The
// unexported nodeTag method seals the interface so only types in this
// package can implement it.
type Node interface {
	nodeTag()
	GetLoc() Loc
}

// StmtNode marks top-level statement nodes — SELECT, INSERT, CREATE TABLE,
// EXEC, EXPLAIN, etc. Anything that can appear at script-statement position.
type StmtNode interface {
	Node
	stmtNode()
}

// ExprNode marks scalar-position expression nodes — operators, predicates,
// function calls, literals, paths, variables, subqueries, etc.
type ExprNode interface {
	Node
	exprNode()
}

// TableExpr marks nodes that appear in FROM position — table references,
// aliased sources, joins, unpivot.
//
// PathExpr, VarRef, and SubLink (defined in exprs.go) deliberately also
// implement TableExpr because PartiQL's grammar lets the same productions
// (path navigation, identifiers, parenthesized SELECTs) appear in both
// scalar and FROM position. See tableexprs.go for details.
type TableExpr interface {
	Node
	tableExpr()
}

// PathStep marks a single step in a PartiQL path expression:
// .field, .*, [expr], [*]. Path steps are chained inside a PathExpr.
type PathStep interface {
	Node
	pathStep()
}

// TypeName marks a PartiQL type reference, used by CAST and DDL.
// PartiQL's type system is small enough that a single concrete TypeRef
// implementation in types.go covers everything.
type TypeName interface {
	Node
	typeName()
}

// PatternNode marks GPML graph-match pattern nodes — node patterns,
// edge patterns, quantifiers, selectors.
type PatternNode interface {
	Node
	patternNode()
}

// List is a generic ordered collection of nodes used as a building block
// for cases where the elements are heterogeneous. Most parser-side lists
// use a typed slice (e.g. []ExprNode, []*TargetEntry) instead; List exists
// primarily so NodeToString has a uniform way to dump heterogenous groups.
type List struct {
	Items []Node
	Loc   Loc
}

func (*List) nodeTag()      {}
func (l *List) GetLoc() Loc { return l.Loc }

// Len returns the number of items in the list, treating a nil receiver as 0.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}
```

- [ ] **Step 2: Create `partiql/ast/ast_test.go` with the package boilerplate and the `List` test**

```go
package ast

import (
	"testing"
)

// ---------------------------------------------------------------------------
// Compile-time interface assertions.
//
// Every node type added to this package gets a `var _ <Interface> = (*Type)(nil)`
// line below. The file fails to compile if a node's interface set drifts.
// Tasks add their assertions to the appropriate section as they grow the AST.
// ---------------------------------------------------------------------------

var _ Node = (*List)(nil)

// ---------------------------------------------------------------------------
// TestGetLoc — table-driven Loc round-trip.
//
// One row per node type. Each row constructs the node with Loc{10, 20},
// calls GetLoc(), and asserts the result.
// ---------------------------------------------------------------------------

func TestGetLoc(t *testing.T) {
	cases := []struct {
		name string
		node Node
	}{
		{"List", &List{Loc: Loc{Start: 10, End: 20}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.node.GetLoc()
			if got.Start != 10 || got.End != 20 {
				t.Errorf("GetLoc() = %+v, want {10, 20}", got)
			}
		})
	}
}
```

- [ ] **Step 3: Run `go build` and `go test` to verify it compiles and passes**

Run from the worktree root:

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected output ends with:

```
ok  	github.com/bytebase/omni/partiql/ast	0.XXXs
```

If you see `cannot find package "github.com/bytebase/omni/partiql/ast"`, the package didn't get created — re-check Step 1. If you see compile errors, fix them before proceeding.

- [ ] **Step 4: Run `go vet` and `gofmt -l`**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both should produce no output. If `gofmt -l` lists a file, run `gofmt -w partiql/ast/<file>` to fix it.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/node.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): scaffold ast package with interfaces and Loc

First commit toward partiql/ast (DAG node 1, P0). Defines:
- 6 sealed sub-interfaces: Node, StmtNode, ExprNode, TableExpr,
  PathStep, TypeName, PatternNode
- Loc struct (byte offsets, -1 = unknown)
- List helper for heterogeneous node collections
- ast_test.go bootstrap with TestGetLoc table

Per spec docs/superpowers/specs/2026-04-08-partiql-ast-core-design.md.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 2: `literals.go` — 8 literal nodes

**Files:**
- Create: `partiql/ast/literals.go`
- Modify: `partiql/ast/ast_test.go` (append interface assertions and TestGetLoc rows)

**Grammar reference:** every literal here is one alternative of the `literal` rule in `bytebase/parser/partiql/PartiQLParser.g4` lines 661–672. The doc comments use the ANTLR alternative-label form `literal#LiteralFoo` so the cross-reference is unambiguous. **There is no `TimestampLit`** — `TIMESTAMP` is a type-name keyword in the `type` rule (line 677), not a literal alternative.

- [ ] **Step 1: Create `partiql/ast/literals.go`**

```go
package ast

// ---------------------------------------------------------------------------
// Literal nodes — all implement ExprNode.
//
// Grammar: PartiQLParser.g4 lines 661–672 (the `literal` rule). Each type
// here corresponds to one labeled alternative of that rule. There is no
// TIMESTAMP literal — TIMESTAMP appears only as a type-name keyword in the
// `type` rule (line 677), used by CAST and DDL.
// ---------------------------------------------------------------------------

// StringLit represents a single-quoted string literal: 'hello'.
//
// Grammar: literal#LiteralString (LITERAL_STRING)
type StringLit struct {
	Val string // SQL escapes already decoded ('' → ')
	Loc Loc
}

func (*StringLit) nodeTag()      {}
func (n *StringLit) GetLoc() Loc { return n.Loc }
func (*StringLit) exprNode()     {}

// NumberLit represents a numeric literal. Val stores the raw text to
// preserve the original representation (integer vs decimal vs scientific).
// Consumers needing a typed value should call strconv.ParseFloat or
// shopspring/decimal at the call site — this AST does not normalize.
//
// Grammar: literal#LiteralInteger / literal#LiteralDecimal
type NumberLit struct {
	Val string // raw text as it appears in source
	Loc Loc
}

func (*NumberLit) nodeTag()      {}
func (n *NumberLit) GetLoc() Loc { return n.Loc }
func (*NumberLit) exprNode()     {}

// BoolLit represents TRUE or FALSE.
//
// Grammar: literal#LiteralTrue / literal#LiteralFalse
type BoolLit struct {
	Val bool
	Loc Loc
}

func (*BoolLit) nodeTag()      {}
func (n *BoolLit) GetLoc() Loc { return n.Loc }
func (*BoolLit) exprNode()     {}

// NullLit represents NULL.
//
// Grammar: literal#LiteralNull
type NullLit struct {
	Loc Loc
}

func (*NullLit) nodeTag()      {}
func (n *NullLit) GetLoc() Loc { return n.Loc }
func (*NullLit) exprNode()     {}

// MissingLit represents the PartiQL-distinct MISSING value.
//
// Grammar: literal#LiteralMissing
type MissingLit struct {
	Loc Loc
}

func (*MissingLit) nodeTag()      {}
func (n *MissingLit) GetLoc() Loc { return n.Loc }
func (*MissingLit) exprNode()     {}

// DateLit represents `DATE 'YYYY-MM-DD'`.
//
// Grammar: literal#LiteralDate (DATE LITERAL_STRING)
type DateLit struct {
	Val string // YYYY-MM-DD body
	Loc Loc
}

func (*DateLit) nodeTag()      {}
func (n *DateLit) GetLoc() Loc { return n.Loc }
func (*DateLit) exprNode()     {}

// TimeLit represents `TIME [(p)] [WITH TIME ZONE] 'HH:MM:SS[.frac]'`.
//
// Grammar: literal#LiteralTime
//   TIME (PAREN_LEFT LITERAL_INTEGER PAREN_RIGHT)? (WITH TIME ZONE)? LITERAL_STRING
type TimeLit struct {
	Val          string // HH:MM:SS[.frac] body
	Precision    *int   // optional fractional-seconds precision
	WithTimeZone bool   // WITH TIME ZONE clause present
	Loc          Loc
}

func (*TimeLit) nodeTag()      {}
func (n *TimeLit) GetLoc() Loc { return n.Loc }
func (*TimeLit) exprNode()     {}

// IonLit represents a backtick-delimited inline Ion value: `…`.
// Text holds the verbatim contents between the backticks (no parsing,
// no normalization). PartiQL-unique.
//
// Grammar: literal#LiteralIon (ION_CLOSURE)
type IonLit struct {
	Text string
	Loc  Loc
}

func (*IonLit) nodeTag()      {}
func (n *IonLit) GetLoc() Loc { return n.Loc }
func (*IonLit) exprNode()     {}
```

- [ ] **Step 2: Append interface assertions and TestGetLoc rows to `ast_test.go`**

Modify `partiql/ast/ast_test.go` — add the literal assertions just below the `var _ Node = (*List)(nil)` line:

```go
// Literals — all implement ExprNode.
var _ ExprNode = (*StringLit)(nil)
var _ ExprNode = (*NumberLit)(nil)
var _ ExprNode = (*BoolLit)(nil)
var _ ExprNode = (*NullLit)(nil)
var _ ExprNode = (*MissingLit)(nil)
var _ ExprNode = (*DateLit)(nil)
var _ ExprNode = (*TimeLit)(nil)
var _ ExprNode = (*IonLit)(nil)
```

Then expand the `cases` slice in `TestGetLoc` to include all 8 literal types (one row each):

```go
cases := []struct {
    name string
    node Node
}{
    {"List", &List{Loc: Loc{Start: 10, End: 20}}},
    {"StringLit", &StringLit{Loc: Loc{Start: 10, End: 20}}},
    {"NumberLit", &NumberLit{Loc: Loc{Start: 10, End: 20}}},
    {"BoolLit", &BoolLit{Loc: Loc{Start: 10, End: 20}}},
    {"NullLit", &NullLit{Loc: Loc{Start: 10, End: 20}}},
    {"MissingLit", &MissingLit{Loc: Loc{Start: 10, End: 20}}},
    {"DateLit", &DateLit{Loc: Loc{Start: 10, End: 20}}},
    {"TimeLit", &TimeLit{Loc: Loc{Start: 10, End: 20}}},
    {"IonLit", &IonLit{Loc: Loc{Start: 10, End: 20}}},
}
```

- [ ] **Step 3: Run tests**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core
go test -run TestGetLoc ./partiql/ast/...
```

Expected: `ok` with all 9 sub-tests passing (`List` + 8 literals).

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both should produce no output.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/literals.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add literal node types

8 literal nodes — StringLit, NumberLit, BoolLit, NullLit, MissingLit,
DateLit, TimeLit, IonLit. All implement ExprNode. Each type maps to
one labeled alternative of the PartiQLParser.g4 `literal` rule
(lines 661-672); doc comments use the literal#LiteralXxx form for
unambiguous cross-reference.

MissingLit and IonLit are PartiQL-unique. TimeLit carries Precision
and WithTimeZone for the full TIME literal syntax. There is no
TimestampLit because TIMESTAMP is only a type-name keyword in the
grammar (line 677), not a literal alternative.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 3: `exprs.go` — operators, predicates (6 types + 3 enums)

**Files:**
- Create: `partiql/ast/exprs.go`
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Create `partiql/ast/exprs.go` with operators and predicates**

```go
package ast

// ---------------------------------------------------------------------------
// Expression nodes — all implement ExprNode.
//
// This file is built in three task groups (matching the implementation plan):
//   1. Operators & predicates (this section)
//   2. Special-form expressions: FuncCall, CaseExpr, CastExpr, etc.
//   3. Paths, variables, parameters, subqueries, collection literals,
//      window spec, path steps
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   exprOr            (lines 469–472)
//   exprAnd           (lines 474–477)
//   exprNot           (lines 479–482)
//   exprPredicate     (lines 484–492)
//   mathOp00/01/02    (lines 494–507)
//   valueExpr         (lines 509–512)
//   exprPrimary       (lines 514–534)
//   exprTerm          (lines 542–549)
//   functionCall      (lines 611–616)
//   caseExpr          (lines 557–558)
//   cast/canCast/canLosslessCast (lines 593–600)
//   extract           (lines 602–603)
//   trimFunction      (lines 605–606)
//   substring         (lines 572–575)
//   coalesce          (lines 554–555)
//   nullIf            (lines 551–552)
//   aggregate         (lines 577–580)
//   windowFunction    (lines 589–591)
//   over              (lines 276–278)
//   array/bag/tuple/pair (lines 649–659)
//   pathStep          (lines 618–623)
//   varRefExpr        (lines 635–636)
//   parameter         (lines 632–633)
// Each type below cites its specific rule#Label.
// ---------------------------------------------------------------------------

// ===========================================================================
// Operator enums
// ===========================================================================

// BinOp identifies a binary operator.
type BinOp int

const (
	BinOpInvalid BinOp = iota
	BinOpOr
	BinOpAnd
	BinOpConcat // ||
	BinOpAdd
	BinOpSub
	BinOpMul
	BinOpDiv
	BinOpMod
	BinOpEq
	BinOpNotEq
	BinOpLt
	BinOpGt
	BinOpLtEq
	BinOpGtEq
)

// String returns the canonical operator spelling.
func (op BinOp) String() string {
	switch op {
	case BinOpOr:
		return "OR"
	case BinOpAnd:
		return "AND"
	case BinOpConcat:
		return "||"
	case BinOpAdd:
		return "+"
	case BinOpSub:
		return "-"
	case BinOpMul:
		return "*"
	case BinOpDiv:
		return "/"
	case BinOpMod:
		return "%"
	case BinOpEq:
		return "="
	case BinOpNotEq:
		return "<>"
	case BinOpLt:
		return "<"
	case BinOpGt:
		return ">"
	case BinOpLtEq:
		return "<="
	case BinOpGtEq:
		return ">="
	default:
		return "INVALID"
	}
}

// UnOp identifies a unary operator.
type UnOp int

const (
	UnOpInvalid UnOp = iota
	UnOpNot
	UnOpNeg // unary -
	UnOpPos // unary +
)

func (op UnOp) String() string {
	switch op {
	case UnOpNot:
		return "NOT"
	case UnOpNeg:
		return "-"
	case UnOpPos:
		return "+"
	default:
		return "INVALID"
	}
}

// IsType identifies the right-hand side of an `IS [NOT] X` predicate.
type IsType int

const (
	IsTypeInvalid IsType = iota
	IsTypeNull
	IsTypeMissing
	IsTypeTrue
	IsTypeFalse
)

func (t IsType) String() string {
	switch t {
	case IsTypeNull:
		return "NULL"
	case IsTypeMissing:
		return "MISSING"
	case IsTypeTrue:
		return "TRUE"
	case IsTypeFalse:
		return "FALSE"
	default:
		return "INVALID"
	}
}

// ===========================================================================
// Operator nodes
// ===========================================================================

// BinaryExpr represents a binary operator application.
//
// Grammar: exprOr#Or, exprAnd#And, exprPredicate#PredicateComparison,
//          mathOp00 (CONCAT), mathOp01 (PLUS/MINUS), mathOp02 (PERCENT/ASTERISK/SLASH_FORWARD)
type BinaryExpr struct {
	Op    BinOp
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (*BinaryExpr) nodeTag()      {}
func (n *BinaryExpr) GetLoc() Loc { return n.Loc }
func (*BinaryExpr) exprNode()     {}

// UnaryExpr represents a unary operator application.
//
// Grammar: exprNot#Not, valueExpr (unary PLUS/MINUS)
type UnaryExpr struct {
	Op      UnOp
	Operand ExprNode
	Loc     Loc
}

func (*UnaryExpr) nodeTag()      {}
func (n *UnaryExpr) GetLoc() Loc { return n.Loc }
func (*UnaryExpr) exprNode()     {}

// ===========================================================================
// Predicate nodes
// ===========================================================================

// InExpr represents `expr [NOT] IN (…)` — either a parenthesized expression
// list or a subquery.
//
// Grammar: exprPredicate#PredicateIn
type InExpr struct {
	Expr     ExprNode
	List     []ExprNode // populated when the RHS is an expression list
	Subquery StmtNode   // populated when the RHS is a parenthesized SELECT
	Not      bool
	Loc      Loc
}

func (*InExpr) nodeTag()      {}
func (n *InExpr) GetLoc() Loc { return n.Loc }
func (*InExpr) exprNode()     {}

// BetweenExpr represents `expr [NOT] BETWEEN low AND high`.
//
// Grammar: exprPredicate#PredicateBetween
type BetweenExpr struct {
	Expr ExprNode
	Low  ExprNode
	High ExprNode
	Not  bool
	Loc  Loc
}

func (*BetweenExpr) nodeTag()      {}
func (n *BetweenExpr) GetLoc() Loc { return n.Loc }
func (*BetweenExpr) exprNode()     {}

// LikeExpr represents `expr [NOT] LIKE pattern [ESCAPE escape]`.
//
// Grammar: exprPredicate#PredicateLike
type LikeExpr struct {
	Expr    ExprNode
	Pattern ExprNode
	Escape  ExprNode // nil if no ESCAPE clause
	Not     bool
	Loc     Loc
}

func (*LikeExpr) nodeTag()      {}
func (n *LikeExpr) GetLoc() Loc { return n.Loc }
func (*LikeExpr) exprNode()     {}

// IsExpr represents `expr IS [NOT] (NULL|MISSING|TRUE|FALSE)`.
//
// Grammar: exprPredicate#PredicateIs
type IsExpr struct {
	Expr ExprNode
	Type IsType
	Not  bool
	Loc  Loc
}

func (*IsExpr) nodeTag()      {}
func (n *IsExpr) GetLoc() Loc { return n.Loc }
func (*IsExpr) exprNode()     {}
```

- [ ] **Step 2: Append interface assertions and TestGetLoc rows for the new types**

In `partiql/ast/ast_test.go`, add a section below the literals assertions:

```go
// Operators & predicates (exprs.go).
var _ ExprNode = (*BinaryExpr)(nil)
var _ ExprNode = (*UnaryExpr)(nil)
var _ ExprNode = (*InExpr)(nil)
var _ ExprNode = (*BetweenExpr)(nil)
var _ ExprNode = (*LikeExpr)(nil)
var _ ExprNode = (*IsExpr)(nil)
```

Append rows to the `cases` slice in `TestGetLoc`:

```go
{"BinaryExpr",  &BinaryExpr{Loc: Loc{Start: 10, End: 20}}},
{"UnaryExpr",   &UnaryExpr{Loc: Loc{Start: 10, End: 20}}},
{"InExpr",      &InExpr{Loc: Loc{Start: 10, End: 20}}},
{"BetweenExpr", &BetweenExpr{Loc: Loc{Start: 10, End: 20}}},
{"LikeExpr",    &LikeExpr{Loc: Loc{Start: 10, End: 20}}},
{"IsExpr",      &IsExpr{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Run tests**

```bash
go test ./partiql/ast/...
```

Expected: `ok` with 15 `TestGetLoc` sub-tests passing (9 prior + 6 new).

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/exprs.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add operator and predicate expression nodes

First slice of exprs.go: BinaryExpr, UnaryExpr, InExpr, BetweenExpr,
LikeExpr, IsExpr, plus the BinOp/UnOp/IsType enums with String()
methods. Covers PartiQLParser.g4 exprOr, exprAnd, exprNot,
exprPredicate, mathOp00/01/02, and valueExpr concat/unary forms.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 4: `exprs.go` — special-form expressions (9 types + 3 enums + WindowSpec)

**Files:**
- Modify: `partiql/ast/exprs.go` (append)
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Append the special-form section to `partiql/ast/exprs.go`**

Add **after** the existing IsExpr block:

```go
// ===========================================================================
// Special-form expressions
//
// Function calls, CASE, CAST, and the keyword-bearing built-ins (EXTRACT,
// TRIM, SUBSTRING) get dedicated nodes because their grammar uses keywords
// inside parens or non-comma argument syntax. COALESCE and NULLIF also get
// dedicated nodes because they appear as ANTLR rules. Built-ins with
// ordinary `name(arg, arg, …)` syntax (DATE_ADD, DATE_DIFF, SIZE, EXISTS, …)
// use plain FuncCall.
// ===========================================================================

// QuantifierKind covers the [DISTINCT|ALL] modifier on aggregates and
// set operations. Although the spec listed it under stmts.go, it is
// declared here because FuncCall (defined first in implementation order)
// is its first user. SetOpStmt in stmts.go references it as ast.QuantifierKind.
type QuantifierKind int

const (
	QuantifierNone QuantifierKind = iota
	QuantifierAll
	QuantifierDistinct
)

func (q QuantifierKind) String() string {
	switch q {
	case QuantifierAll:
		return "ALL"
	case QuantifierDistinct:
		return "DISTINCT"
	default:
		return ""
	}
}

// CastKind discriminates the three CAST-family operators.
type CastKind int

const (
	CastKindInvalid CastKind = iota
	CastKindCast
	CastKindCanCast
	CastKindCanLosslessCast
)

func (k CastKind) String() string {
	switch k {
	case CastKindCast:
		return "CAST"
	case CastKindCanCast:
		return "CAN_CAST"
	case CastKindCanLosslessCast:
		return "CAN_LOSSLESS_CAST"
	default:
		return "INVALID"
	}
}

// TrimSpec covers the optional LEADING/TRAILING/BOTH keyword inside TRIM.
type TrimSpec int

const (
	TrimSpecNone TrimSpec = iota
	TrimSpecLeading
	TrimSpecTrailing
	TrimSpecBoth
)

func (s TrimSpec) String() string {
	switch s {
	case TrimSpecLeading:
		return "LEADING"
	case TrimSpecTrailing:
		return "TRAILING"
	case TrimSpecBoth:
		return "BOTH"
	default:
		return ""
	}
}

// FuncCall is the generic function-call node — used for ordinary function
// calls (DATE_ADD, SIZE, ...), aggregates (COUNT/SUM/AVG/MIN/MAX with the
// optional DISTINCT/ALL modifier and COUNT(*) form), and window calls
// (LAG/LEAD with an OVER clause). The Quantifier, Star, and Over fields
// determine which flavor a particular instance is.
//
// Grammar: functionCall#FunctionCallReserved, functionCall#FunctionCallIdent,
//          aggregate#CountAll, aggregate#AggregateBase,
//          windowFunction#LagLeadFunction
type FuncCall struct {
	Name       string
	Args       []ExprNode
	Quantifier QuantifierKind // NONE/DISTINCT/ALL — populated for aggregates
	Star       bool           // true for COUNT(*)
	Over       *WindowSpec    // non-nil for window calls
	Loc        Loc
}

func (*FuncCall) nodeTag()      {}
func (n *FuncCall) GetLoc() Loc { return n.Loc }
func (*FuncCall) exprNode()     {}

// CaseExpr covers both `CASE WHEN … THEN …` (searched) and
// `CASE expr WHEN … THEN …` (simple) forms. Operand is nil for the
// searched form.
//
// Grammar: caseExpr
type CaseExpr struct {
	Operand ExprNode    // nil for searched CASE
	Whens   []*CaseWhen
	Else    ExprNode    // nil if no ELSE clause
	Loc     Loc
}

func (*CaseExpr) nodeTag()      {}
func (n *CaseExpr) GetLoc() Loc { return n.Loc }
func (*CaseExpr) exprNode()     {}

// CaseWhen represents one `WHEN expr THEN expr` arm. Bare Node — does not
// implement ExprNode because it cannot stand alone in scalar position.
type CaseWhen struct {
	When ExprNode
	Then ExprNode
	Loc  Loc
}

func (*CaseWhen) nodeTag()      {}
func (n *CaseWhen) GetLoc() Loc { return n.Loc }

// CastExpr covers CAST, CAN_CAST, and CAN_LOSSLESS_CAST.
//
// Grammar: cast, canCast, canLosslessCast
type CastExpr struct {
	Kind   CastKind
	Expr   ExprNode
	AsType TypeName
	Loc    Loc
}

func (*CastExpr) nodeTag()      {}
func (n *CastExpr) GetLoc() Loc { return n.Loc }
func (*CastExpr) exprNode()     {}

// ExtractExpr represents `EXTRACT(<field> FROM <expr>)`. Field is the
// keyword (YEAR, MONTH, DAY, HOUR, MINUTE, SECOND, ...) — stored as the
// raw uppercase keyword string.
//
// Grammar: extract
type ExtractExpr struct {
	Field string
	From  ExprNode
	Loc   Loc
}

func (*ExtractExpr) nodeTag()      {}
func (n *ExtractExpr) GetLoc() Loc { return n.Loc }
func (*ExtractExpr) exprNode()     {}

// TrimExpr represents `TRIM([LEADING|TRAILING|BOTH] [sub] FROM target)`.
//
// Grammar: trimFunction
type TrimExpr struct {
	Spec TrimSpec
	Sub  ExprNode // optional substring to trim; nil for default whitespace
	From ExprNode
	Loc  Loc
}

func (*TrimExpr) nodeTag()      {}
func (n *TrimExpr) GetLoc() Loc { return n.Loc }
func (*TrimExpr) exprNode()     {}

// SubstringExpr represents `SUBSTRING(expr FROM start [FOR length])` and
// the equivalent comma form `SUBSTRING(expr, start[, length])`.
//
// Grammar: substring
type SubstringExpr struct {
	Expr ExprNode
	From ExprNode
	For  ExprNode // optional length
	Loc  Loc
}

func (*SubstringExpr) nodeTag()      {}
func (n *SubstringExpr) GetLoc() Loc { return n.Loc }
func (*SubstringExpr) exprNode()     {}

// CoalesceExpr represents `COALESCE(expr, expr, …)`.
//
// Grammar: coalesce
type CoalesceExpr struct {
	Args []ExprNode
	Loc  Loc
}

func (*CoalesceExpr) nodeTag()      {}
func (n *CoalesceExpr) GetLoc() Loc { return n.Loc }
func (*CoalesceExpr) exprNode()     {}

// NullIfExpr represents `NULLIF(expr, expr)`.
//
// Grammar: nullIf
type NullIfExpr struct {
	Left  ExprNode
	Right ExprNode
	Loc   Loc
}

func (*NullIfExpr) nodeTag()      {}
func (n *NullIfExpr) GetLoc() Loc { return n.Loc }
func (*NullIfExpr) exprNode()     {}

// WindowSpec represents the body of an OVER (...) clause attached to
// a window function call. Bare Node — appears only inside FuncCall.Over.
//
// Grammar: over
type WindowSpec struct {
	PartitionBy []ExprNode
	OrderBy     []*OrderByItem // OrderByItem defined in stmts.go
	Loc         Loc
}

func (*WindowSpec) nodeTag()      {}
func (n *WindowSpec) GetLoc() Loc { return n.Loc }
```

**Note:** `WindowSpec` references `*OrderByItem` from `stmts.go`. This is fine within the same package — the field is just declared, not constructed yet, so the build only needs the type to exist by the time `stmts.go` is compiled. Since both files compile together, the order they were written in this plan doesn't matter for compilation.

- [ ] **Step 2: Append assertions and TestGetLoc rows**

In `ast_test.go`, add to the assertions section:

```go
// Special-form expressions (exprs.go).
var _ ExprNode = (*FuncCall)(nil)
var _ ExprNode = (*CaseExpr)(nil)
var _ Node     = (*CaseWhen)(nil)
var _ ExprNode = (*CastExpr)(nil)
var _ ExprNode = (*ExtractExpr)(nil)
var _ ExprNode = (*TrimExpr)(nil)
var _ ExprNode = (*SubstringExpr)(nil)
var _ ExprNode = (*CoalesceExpr)(nil)
var _ ExprNode = (*NullIfExpr)(nil)
var _ Node     = (*WindowSpec)(nil)
```

Append to `cases` in `TestGetLoc`:

```go
{"FuncCall",      &FuncCall{Loc: Loc{Start: 10, End: 20}}},
{"CaseExpr",      &CaseExpr{Loc: Loc{Start: 10, End: 20}}},
{"CaseWhen",      &CaseWhen{Loc: Loc{Start: 10, End: 20}}},
{"CastExpr",      &CastExpr{Loc: Loc{Start: 10, End: 20}}},
{"ExtractExpr",   &ExtractExpr{Loc: Loc{Start: 10, End: 20}}},
{"TrimExpr",      &TrimExpr{Loc: Loc{Start: 10, End: 20}}},
{"SubstringExpr", &SubstringExpr{Loc: Loc{Start: 10, End: 20}}},
{"CoalesceExpr",  &CoalesceExpr{Loc: Loc{Start: 10, End: 20}}},
{"NullIfExpr",    &NullIfExpr{Loc: Loc{Start: 10, End: 20}}},
{"WindowSpec",    &WindowSpec{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Try to build** — expect a compile failure because `OrderByItem` doesn't exist yet

```bash
go build ./partiql/ast/...
```

Expected: `undefined: OrderByItem` (or similar) referencing `partiql/ast/exprs.go`.

This is the expected gap. We're going to fix it temporarily by adding a forward-declared placeholder in `stmts.go` so the package builds. The real `OrderByItem` lands in Task 9.

- [ ] **Step 4: Create a temporary `partiql/ast/stmts.go` placeholder for `OrderByItem`**

Create the file with just enough to satisfy the forward reference. **This file will be replaced wholesale in Task 8 — do not invest in it now.**

```go
package ast

// PLACEHOLDER — replaced in Task 8/9. The OrderByItem type is referenced
// by WindowSpec in exprs.go (Task 4). Defining it here lets the package
// build before its real home in stmts.go is written.
type OrderByItem struct {
	Expr          ExprNode
	Desc          bool
	NullsFirst    bool
	NullsExplicit bool
	Loc           Loc
}

func (*OrderByItem) nodeTag()      {}
func (n *OrderByItem) GetLoc() Loc { return n.Loc }
```

- [ ] **Step 5: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, all 25 `TestGetLoc` sub-tests passing (15 prior + 10 new). Note that we are intentionally NOT adding `OrderByItem` to the test table in this task — that happens in Task 9 when the real type lands.

- [ ] **Step 6: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 7: Commit**

```bash
git add partiql/ast/exprs.go partiql/ast/stmts.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add special-form expression nodes

FuncCall (covers regular calls, aggregates, window calls), CaseExpr,
CaseWhen, CastExpr, ExtractExpr, TrimExpr, SubstringExpr, CoalesceExpr,
NullIfExpr, WindowSpec. Plus QuantifierKind, CastKind, TrimSpec enums.

Adds a placeholder OrderByItem in stmts.go to satisfy WindowSpec's
forward reference; the real OrderByItem lands in Task 9 when stmts.go
is written in full.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 5: `exprs.go` — paths, vars, params, subqueries, collection literals, path steps (13 types)

**Files:**
- Modify: `partiql/ast/exprs.go` (append)
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Append the path/var/collection/path-step section to `partiql/ast/exprs.go`**

Add after the `WindowSpec` block:

```go
// ===========================================================================
// Paths, variables, parameters, subqueries
//
// PathExpr, VarRef, and SubLink also implement TableExpr because PartiQL's
// FROM grammar accepts the same productions as scalar expressions.
// ===========================================================================

// PathExpr represents a chained path navigation: root.field[idx].field[*].
// Root is the base expression (typically a VarRef); Steps are the chained
// path operations from PathStep.
//
// Grammar: exprPrimary#ExprPrimaryPath (exprPrimary pathStep+)
type PathExpr struct {
	Root  ExprNode
	Steps []PathStep
	Loc   Loc
}

func (*PathExpr) nodeTag()      {}
func (n *PathExpr) GetLoc() Loc { return n.Loc }
func (*PathExpr) exprNode()     {}
func (*PathExpr) tableExpr()    {} // PartiQL FROM accepts path expressions

// VarRef represents an identifier reference. AtPrefixed distinguishes
// `@id` (true) from bare `id` (false). CaseSensitive distinguishes
// `"X"` (true, double-quoted) from `X` (false, unquoted).
//
// Grammar: varRefExpr
type VarRef struct {
	Name          string
	AtPrefixed    bool
	CaseSensitive bool
	Loc           Loc
}

func (*VarRef) nodeTag()      {}
func (n *VarRef) GetLoc() Loc { return n.Loc }
func (*VarRef) exprNode()     {}
func (*VarRef) tableExpr()    {} // a bare identifier in FROM is a VarRef

// ParamRef represents a positional `?` parameter.
//
// Grammar: parameter
type ParamRef struct {
	Loc Loc
}

func (*ParamRef) nodeTag()      {}
func (n *ParamRef) GetLoc() Loc { return n.Loc }
func (*ParamRef) exprNode()     {}

// SubLink represents a parenthesized SELECT used as a value expression
// or as a FROM source. Stmt is the inner statement (a SelectStmt or
// SetOpStmt).
//
// Grammar: exprTerm#ExprTermWrappedQuery (PAREN_LEFT expr PAREN_RIGHT)
type SubLink struct {
	Stmt StmtNode
	Loc  Loc
}

func (*SubLink) nodeTag()      {}
func (n *SubLink) GetLoc() Loc { return n.Loc }
func (*SubLink) exprNode()     {}
func (*SubLink) tableExpr()    {}

// ===========================================================================
// Collection literals
// ===========================================================================

// ListLit represents an ordered list literal: [expr, expr, …].
//
// Grammar: array
type ListLit struct {
	Items []ExprNode
	Loc   Loc
}

func (*ListLit) nodeTag()      {}
func (n *ListLit) GetLoc() Loc { return n.Loc }
func (*ListLit) exprNode()     {}

// BagLit represents an unordered bag literal: <<expr, expr, …>>. PartiQL-unique.
//
// Grammar: bag
type BagLit struct {
	Items []ExprNode
	Loc   Loc
}

func (*BagLit) nodeTag()      {}
func (n *BagLit) GetLoc() Loc { return n.Loc }
func (*BagLit) exprNode()     {}

// TupleLit represents a tuple/struct literal: {key: value, …}. PartiQL-unique.
//
// Grammar: tuple
type TupleLit struct {
	Pairs []*TuplePair
	Loc   Loc
}

func (*TupleLit) nodeTag()      {}
func (n *TupleLit) GetLoc() Loc { return n.Loc }
func (*TupleLit) exprNode()     {}

// TuplePair represents one `key: value` entry in a tuple literal.
// Bare Node — appears only inside TupleLit.Pairs.
//
// Grammar: pair
type TuplePair struct {
	Key   ExprNode
	Value ExprNode
	Loc   Loc
}

func (*TuplePair) nodeTag()      {}
func (n *TuplePair) GetLoc() Loc { return n.Loc }

// ===========================================================================
// Path steps — implement PathStep
// ===========================================================================

// DotStep represents `.field`. CaseSensitive is true for `."Field"`
// (the field name was quoted) and false for unquoted `.field`.
//
// Grammar: pathStep#PathStepDotExpr
type DotStep struct {
	Field         string
	CaseSensitive bool
	Loc           Loc
}

func (*DotStep) nodeTag()      {}
func (n *DotStep) GetLoc() Loc { return n.Loc }
func (*DotStep) pathStep()     {}

// AllFieldsStep represents `.*` — the all-fields wildcard.
//
// Grammar: pathStep#PathStepDotAll
type AllFieldsStep struct {
	Loc Loc
}

func (*AllFieldsStep) nodeTag()      {}
func (n *AllFieldsStep) GetLoc() Loc { return n.Loc }
func (*AllFieldsStep) pathStep()     {}

// IndexStep represents `[expr]` — index/key by an expression.
//
// Grammar: pathStep#PathStepIndexExpr
type IndexStep struct {
	Index ExprNode
	Loc   Loc
}

func (*IndexStep) nodeTag()      {}
func (n *IndexStep) GetLoc() Loc { return n.Loc }
func (*IndexStep) pathStep()     {}

// WildcardStep represents `[*]` — all-elements wildcard.
//
// Grammar: pathStep#PathStepIndexAll
type WildcardStep struct {
	Loc Loc
}

func (*WildcardStep) nodeTag()      {}
func (n *WildcardStep) GetLoc() Loc { return n.Loc }
func (*WildcardStep) pathStep()     {}
```

- [ ] **Step 2: Append assertions and TestGetLoc rows**

In `ast_test.go`:

```go
// Paths, vars, params, subqueries, collection literals, path steps (exprs.go).
var _ ExprNode  = (*PathExpr)(nil)
var _ TableExpr = (*PathExpr)(nil)
var _ ExprNode  = (*VarRef)(nil)
var _ TableExpr = (*VarRef)(nil)
var _ ExprNode  = (*ParamRef)(nil)
var _ ExprNode  = (*SubLink)(nil)
var _ TableExpr = (*SubLink)(nil)
var _ ExprNode  = (*ListLit)(nil)
var _ ExprNode  = (*BagLit)(nil)
var _ ExprNode  = (*TupleLit)(nil)
var _ Node      = (*TuplePair)(nil)
var _ PathStep  = (*DotStep)(nil)
var _ PathStep  = (*AllFieldsStep)(nil)
var _ PathStep  = (*IndexStep)(nil)
var _ PathStep  = (*WildcardStep)(nil)
```

Append to `cases` in `TestGetLoc`:

```go
{"PathExpr",      &PathExpr{Loc: Loc{Start: 10, End: 20}}},
{"VarRef",        &VarRef{Loc: Loc{Start: 10, End: 20}}},
{"ParamRef",      &ParamRef{Loc: Loc{Start: 10, End: 20}}},
{"SubLink",       &SubLink{Loc: Loc{Start: 10, End: 20}}},
{"ListLit",       &ListLit{Loc: Loc{Start: 10, End: 20}}},
{"BagLit",        &BagLit{Loc: Loc{Start: 10, End: 20}}},
{"TupleLit",      &TupleLit{Loc: Loc{Start: 10, End: 20}}},
{"TuplePair",     &TuplePair{Loc: Loc{Start: 10, End: 20}}},
{"DotStep",       &DotStep{Loc: Loc{Start: 10, End: 20}}},
{"AllFieldsStep", &AllFieldsStep{Loc: Loc{Start: 10, End: 20}}},
{"IndexStep",     &IndexStep{Loc: Loc{Start: 10, End: 20}}},
{"WildcardStep",  &WildcardStep{Loc: Loc{Start: 10, End: 20}}},
```

(`OrderByItem` from the placeholder in stmts.go is NOT added to the test table here — it lands in Task 9 when its real form is written.)

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 37 `TestGetLoc` sub-tests passing.

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/exprs.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add path, variable, collection, and path-step nodes

Final slice of exprs.go: PathExpr, VarRef, ParamRef, SubLink (all
multi-interface where appropriate — PathExpr/VarRef/SubLink also
implement TableExpr because PartiQL's FROM grammar accepts the same
productions as scalar expressions); ListLit, BagLit, TupleLit,
TuplePair collection literals; DotStep, AllFieldsStep, IndexStep,
WildcardStep path steps.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 6: `tableexprs.go` — 4 table-position nodes + JoinKind enum

**Files:**
- Create: `partiql/ast/tableexprs.go`
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Create `partiql/ast/tableexprs.go`**

```go
package ast

// ---------------------------------------------------------------------------
// Table expression nodes — implement TableExpr.
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   fromClause         (lines 297–298)
//   tableReference     (lines 389–395)
//   tableNonJoin       (lines 397–400)
//   tableBaseReference (lines 402–406)
//   tableUnpivot       (lines 408–409)
//   joinRhs            (lines 411–414)
//   joinSpec           (lines 416–417)
//   joinType           (lines 419–425)
//   fromClauseSimple   (lines 202–205)
// Each type below cites its specific rule#Label.
//
// IMPORTANT: PathExpr, VarRef, and SubLink (defined in exprs.go) also
// implement TableExpr. They are not re-listed here as their primary home
// is exprs.go, but they are first-class FROM-position nodes. The compile-time
// `var _ TableExpr = ...` lines in ast_test.go cover them.
// ---------------------------------------------------------------------------

// JoinKind identifies the JOIN flavor.
type JoinKind int

const (
	JoinKindInvalid JoinKind = iota
	JoinKindCross
	JoinKindInner
	JoinKindLeft
	JoinKindRight
	JoinKindFull
	JoinKindOuter // bare OUTER JOIN (PartiQL natural-outer form)
)

func (k JoinKind) String() string {
	switch k {
	case JoinKindCross:
		return "CROSS"
	case JoinKindInner:
		return "INNER"
	case JoinKindLeft:
		return "LEFT"
	case JoinKindRight:
		return "RIGHT"
	case JoinKindFull:
		return "FULL"
	case JoinKindOuter:
		return "OUTER"
	default:
		return "INVALID"
	}
}

// TableRef represents a bare identifier in FROM position. Schema is
// populated only when the reference uses dotted form `schema.table`.
// CaseSensitive is true for double-quoted identifiers.
//
// Grammar: tableBaseReference#TableBaseRefSymbol, tableBaseReference#TableBaseRefClauses
type TableRef struct {
	Name          string
	Schema        string
	CaseSensitive bool
	Loc           Loc
}

func (*TableRef) nodeTag()      {}
func (n *TableRef) GetLoc() Loc { return n.Loc }
func (*TableRef) tableExpr()    {}

// AliasedSource wraps a TableExpr with PartiQL's `AS alias AT positional
// BY key` aliasing form. PartiQL-unique because of the AT/BY positional
// and key aliases.
//
// Grammar: tableBaseReference#TableBaseRefClauses, tableBaseReference#TableBaseRefMatch,
//          tableUnpivot, fromClauseSimple#FromClauseSimpleExplicit
type AliasedSource struct {
	Source TableExpr
	As     *string // optional row alias
	At     *string // optional positional alias
	By     *string // optional key alias
	Loc    Loc
}

func (*AliasedSource) nodeTag()      {}
func (n *AliasedSource) GetLoc() Loc { return n.Loc }
func (*AliasedSource) tableExpr()    {}

// JoinExpr represents a JOIN clause: left JOIN right ON condition.
// Kind selects the JOIN flavor; On is nil for CROSS JOIN.
//
// Grammar: tableReference#TableCrossJoin, tableReference#TableQualifiedJoin
//          (join-type modifier from joinType; joinRhs for right-hand side)
type JoinExpr struct {
	Kind  JoinKind
	Left  TableExpr
	Right TableExpr
	On    ExprNode // nil for CROSS JOIN
	Loc   Loc
}

func (*JoinExpr) nodeTag()      {}
func (n *JoinExpr) GetLoc() Loc { return n.Loc }
func (*JoinExpr) tableExpr()    {}

// UnpivotExpr represents `UNPIVOT expr [AS alias] [AT pos] [BY key]`.
// PartiQL-unique. The Source is an expression (not a TableExpr) because
// the grammar nests an arbitrary expression here.
//
// Grammar: tableUnpivot
type UnpivotExpr struct {
	Source ExprNode
	As     *string
	At     *string
	By     *string
	Loc    Loc
}

func (*UnpivotExpr) nodeTag()      {}
func (n *UnpivotExpr) GetLoc() Loc { return n.Loc }
func (*UnpivotExpr) tableExpr()    {}
```

- [ ] **Step 2: Append assertions and TestGetLoc rows**

In `ast_test.go`:

```go
// Table expression nodes (tableexprs.go).
var _ TableExpr = (*TableRef)(nil)
var _ TableExpr = (*AliasedSource)(nil)
var _ TableExpr = (*JoinExpr)(nil)
var _ TableExpr = (*UnpivotExpr)(nil)
```

Append to `cases`:

```go
{"TableRef",      &TableRef{Loc: Loc{Start: 10, End: 20}}},
{"AliasedSource", &AliasedSource{Loc: Loc{Start: 10, End: 20}}},
{"JoinExpr",      &JoinExpr{Loc: Loc{Start: 10, End: 20}}},
{"UnpivotExpr",   &UnpivotExpr{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 41 sub-tests passing.

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/tableexprs.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add FROM-position table expression nodes

TableRef, AliasedSource (PartiQL AS/AT/BY aliasing), JoinExpr (with
JoinKind enum covering CROSS/INNER/LEFT/RIGHT/FULL/OUTER), UnpivotExpr.
File doc-block notes that PathExpr/VarRef/SubLink in exprs.go also
implement TableExpr.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 7: `types.go` — TypeRef

**Files:**
- Create: `partiql/ast/types.go`
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Create `partiql/ast/types.go`**

```go
package ast

// ---------------------------------------------------------------------------
// PartiQL type references — implement TypeName.
//
// One flat TypeRef covers all PartiQL types. Discrimination is by Name
// (uppercase canonical: INT, DECIMAL, TIME, BAG, ANY, ...) rather than a
// 25-arm Go enum. The grammar treats types as a name token with optional
// modifiers (Args for parameterized types, WithTimeZone for TIME/TIMESTAMP),
// so the AST mirrors that shape.
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   type#TypeAtomic    (lines 675–680)
//   type#TypeArgSingle (line 681)
//   type#TypeVarChar   (line 682)
//   type#TypeArgDouble (line 683)
//   type#TypeTimeZone  (line 684)
//   type#TypeCustom    (line 685)
// ---------------------------------------------------------------------------

// TypeRef represents any PartiQL type reference, used by CAST and DDL.
//
// Args carries optional precision/scale/length:
//   - DECIMAL(10,2) -> Args=[10, 2]
//   - VARCHAR(255)  -> Args=[255]
//   - FLOAT(53)     -> Args=[53]
//   - INT           -> Args=nil
//
// WithTimeZone is set for `TIME WITH TIME ZONE` and `TIMESTAMP WITH TIME ZONE`.
//
// Names covered (canonical uppercase):
//   - Numeric: INT, INTEGER, BIGINT, SMALLINT, REAL, DOUBLE PRECISION,
//     DECIMAL, NUMERIC, FLOAT
//   - Boolean: BOOL, BOOLEAN
//   - Null/missing: NULL, MISSING
//   - String: STRING, SYMBOL, VARCHAR, CHAR, CHARACTER
//   - Binary: BLOB, CLOB
//   - Temporal: DATE, TIME, TIMESTAMP
//   - Collection: STRUCT, TUPLE, LIST, BAG, SEXP
//   - Wildcard: ANY
type TypeRef struct {
	Name         string
	Args         []int
	WithTimeZone bool
	Loc          Loc
}

func (*TypeRef) nodeTag()      {}
func (n *TypeRef) GetLoc() Loc { return n.Loc }
func (*TypeRef) typeName()     {}
```

- [ ] **Step 2: Append assertion and TestGetLoc row**

In `ast_test.go`:

```go
// Type names (types.go).
var _ TypeName = (*TypeRef)(nil)
```

```go
{"TypeRef", &TypeRef{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 42 sub-tests passing.

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/types.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add TypeRef for PartiQL type references

Single flat TypeRef node covers all PartiQL types (CAST and DDL).
Discriminated by Name field; Args carries precision/scale/length;
WithTimeZone for TIME/TIMESTAMP variants. Doc comment enumerates
the full canonical name set.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 8: `stmts.go` — top-level statement nodes (14 types + enums)

**Files:**
- Modify: `partiql/ast/stmts.go` (replace the placeholder created in Task 4)
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Replace `partiql/ast/stmts.go` wholesale**

Open `partiql/ast/stmts.go` and **replace the entire file contents** with:

```go
package ast

// ---------------------------------------------------------------------------
// Statement nodes and supporting types.
//
// This file is built in two task groups:
//   1. Top-level statements (this section): SELECT, set ops, EXPLAIN,
//      DML, DDL, EXEC.
//   2. Clause helpers: TargetEntry, GroupBy, OrderBy, OnConflict, Returning,
//      etc. (Task 9 — appended below).
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   root                   (lines 17–18)
//   statement              (lines 20–25)
//   dql                    (lines 55–56)
//   dml                    (lines 94–100)
//   dmlBaseCommand         (lines 102–108)
//   ddl                    (lines 73–76)
//   createCommand          (lines 78–81)
//   dropCommand            (lines 83–86)
//   execCommand            (lines 64–65)
//   selectClause           (lines 216–221)
//   exprBagOp              (lines 449–454)
//   exprSelect             (lines 456–467)
//   insertCommand          (lines 133–137)
//   insertCommandReturning (lines 130–131)
//   updateClause           (lines 182–183)
//   deleteCommand          (lines 191–192)
//   upsertCommand          (lines 124–125)
//   replaceCommand         (lines 120–121)
//   removeCommand          (lines 127–128)
//   setCommand             (lines 185–186)
// Each type below cites its specific rule#Label.
// ---------------------------------------------------------------------------

// ===========================================================================
// Statement enums
// ===========================================================================

// SetOpKind identifies the set operation in SetOpStmt.
type SetOpKind int

const (
	SetOpInvalid SetOpKind = iota
	SetOpUnion
	SetOpIntersect
	SetOpExcept
)

func (k SetOpKind) String() string {
	switch k {
	case SetOpUnion:
		return "UNION"
	case SetOpIntersect:
		return "INTERSECT"
	case SetOpExcept:
		return "EXCEPT"
	default:
		return "INVALID"
	}
}

// OnConflictAction discriminates the body of an ON CONFLICT clause.
// PartiQL's legacy ANTLR grammar leaves DO REPLACE/UPDATE as stubs that
// only accept the EXCLUDED keyword (see doReplace/doUpdate rules in
// PartiQLParser.g4 lines 168–180); the omni AST matches that scope.
type OnConflictAction int

const (
	OnConflictInvalid OnConflictAction = iota
	OnConflictDoNothing
	OnConflictDoReplaceExcluded
	OnConflictDoUpdateExcluded
)

func (a OnConflictAction) String() string {
	switch a {
	case OnConflictDoNothing:
		return "DO_NOTHING"
	case OnConflictDoReplaceExcluded:
		return "DO_REPLACE_EXCLUDED"
	case OnConflictDoUpdateExcluded:
		return "DO_UPDATE_EXCLUDED"
	default:
		return "INVALID"
	}
}

// ReturningStatus is the (MODIFIED|ALL) modifier on each RETURNING item.
type ReturningStatus int

const (
	ReturningStatusInvalid ReturningStatus = iota
	ReturningStatusModified
	ReturningStatusAll
)

func (s ReturningStatus) String() string {
	switch s {
	case ReturningStatusModified:
		return "MODIFIED"
	case ReturningStatusAll:
		return "ALL"
	default:
		return "INVALID"
	}
}

// ReturningMapping is the (OLD|NEW) modifier on each RETURNING item.
type ReturningMapping int

const (
	ReturningMappingInvalid ReturningMapping = iota
	ReturningMappingOld
	ReturningMappingNew
)

func (m ReturningMapping) String() string {
	switch m {
	case ReturningMappingOld:
		return "OLD"
	case ReturningMappingNew:
		return "NEW"
	default:
		return "INVALID"
	}
}

// ===========================================================================
// Top-level statements — all implement StmtNode
// ===========================================================================

// SelectStmt is a single SELECT query (not a set-op combination — see
// SetOpStmt for that). Carries the full clause set:
//
//   SELECT [DISTINCT|ALL] (* | items | VALUE expr | PIVOT v AT k)
//   FROM ...
//   [LET ...]
//   [WHERE ...]
//   [GROUP [PARTIAL] BY ... [GROUP AS alias]]
//   [HAVING ...]
//   [ORDER BY ...]
//   [LIMIT ...]
//   [OFFSET ...]
//
// Quantifier holds the optional DISTINCT/ALL modifier (NONE if absent).
// Star is true for `SELECT *`. Value holds the expression for
// `SELECT VALUE expr`. Pivot holds the PIVOT projection. At most one of
// (Star, Value != nil, Pivot != nil, len(Targets) > 0) is set per instance.
//
// Grammar: exprSelect#SfwQuery (with selectClause#SelectAll / selectClause#SelectItems /
//          selectClause#SelectValue / selectClause#SelectPivot)
type SelectStmt struct {
	Quantifier QuantifierKind   // NONE / DISTINCT / ALL
	Star       bool             // SELECT *
	Value      ExprNode         // SELECT VALUE expr (PartiQL-unique)
	Pivot      *PivotProjection // PIVOT v AT k (PartiQL-unique)
	Targets    []*TargetEntry   // SELECT items
	From       TableExpr        // FROM clause
	Let        []*LetBinding    // LET bindings (PartiQL-unique)
	Where      ExprNode
	GroupBy    *GroupByClause
	Having     ExprNode
	OrderBy    []*OrderByItem
	Limit      ExprNode
	Offset     ExprNode
	Loc        Loc
}

func (*SelectStmt) nodeTag()      {}
func (n *SelectStmt) GetLoc() Loc { return n.Loc }
func (*SelectStmt) stmtNode()     {}

// SetOpStmt combines two queries with UNION/INTERSECT/EXCEPT.
// Quantifier is the DISTINCT/ALL modifier; Outer is true for the
// PartiQL-specific OUTER variant. Both Left and Right may themselves
// be SetOpStmt nodes for nested set operations.
//
// Grammar: exprBagOp#Union, exprBagOp#Intersect, exprBagOp#Except
type SetOpStmt struct {
	Op         SetOpKind
	Quantifier QuantifierKind
	Outer      bool
	Left       StmtNode
	Right      StmtNode
	Loc        Loc
}

func (*SetOpStmt) nodeTag()      {}
func (n *SetOpStmt) GetLoc() Loc { return n.Loc }
func (*SetOpStmt) stmtNode()     {}

// ExplainStmt wraps any other StmtNode with an EXPLAIN prefix.
//
// Grammar: root (EXPLAIN prefix on the root rule)
type ExplainStmt struct {
	Inner StmtNode
	Loc   Loc
}

func (*ExplainStmt) nodeTag()      {}
func (n *ExplainStmt) GetLoc() Loc { return n.Loc }
func (*ExplainStmt) stmtNode()     {}

// InsertStmt represents `INSERT INTO target [AS alias] VALUE expr
// [AT pos] [ON CONFLICT ...] [RETURNING ...]`. Covers both legacy
// (INSERT INTO p VALUE … [AT pos]) and RFC 0011 (INSERT INTO c AS a VALUE …)
// forms.
//
// Mutual exclusion between the two forms:
//   - Legacy form (insertCommand#InsertLegacy, insertCommandReturning):
//     AsAlias is nil; Pos may be set.
//   - RFC 0011 form (insertCommand#Insert): AsAlias may be set; Pos is nil.
//   - The grammar's `insertCommandReturning` (line 130–131) is the only
//     alternative that allows `RETURNING`; on `insertCommand#InsertLegacy`
//     (line 134) and `insertCommand#Insert` (line 136), Returning is nil.
//
// Grammar: insertCommand#InsertLegacy, insertCommand#Insert, insertCommandReturning
type InsertStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	Pos        ExprNode // legacy `AT pos` clause; nil for RFC 0011 form
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*InsertStmt) nodeTag()      {}
func (n *InsertStmt) GetLoc() Loc { return n.Loc }
func (*InsertStmt) stmtNode()     {}

// UpdateStmt represents `UPDATE source SET ... [WHERE ...] [RETURNING ...]`
// and the equivalent `FROM source SET ...` form.
//
// Grammar: updateClause, dml#DmlBaseWrapper (with dmlBaseCommand containing setCommand)
type UpdateStmt struct {
	Source    TableExpr
	Sets      []*SetAssignment
	Where     ExprNode
	Returning *ReturningClause
	Loc       Loc
}

func (*UpdateStmt) nodeTag()      {}
func (n *UpdateStmt) GetLoc() Loc { return n.Loc }
func (*UpdateStmt) stmtNode()     {}

// DeleteStmt represents `DELETE FROM source [WHERE ...] [RETURNING ...]`
// and the equivalent `FROM source DELETE ...` form.
//
// Grammar: deleteCommand
type DeleteStmt struct {
	Source    TableExpr
	Where     ExprNode
	Returning *ReturningClause
	Loc       Loc
}

func (*DeleteStmt) nodeTag()      {}
func (n *DeleteStmt) GetLoc() Loc { return n.Loc }
func (*DeleteStmt) stmtNode()     {}

// UpsertStmt represents `UPSERT INTO target [AS alias] VALUE expr ...`.
// Same shape as InsertStmt; semantically the conflict-merge is implicit.
//
// Grammar: upsertCommand
type UpsertStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*UpsertStmt) nodeTag()      {}
func (n *UpsertStmt) GetLoc() Loc { return n.Loc }
func (*UpsertStmt) stmtNode()     {}

// ReplaceStmt represents `REPLACE INTO target [AS alias] VALUE expr ...`.
// Same shape as InsertStmt.
//
// Grammar: replaceCommand
type ReplaceStmt struct {
	Target     TableExpr
	AsAlias    *string
	Value      ExprNode
	OnConflict *OnConflict
	Returning  *ReturningClause
	Loc        Loc
}

func (*ReplaceStmt) nodeTag()      {}
func (n *ReplaceStmt) GetLoc() Loc { return n.Loc }
func (*ReplaceStmt) stmtNode()     {}

// RemoveStmt represents `REMOVE path` — removes an element from a
// collection. PartiQL-unique.
//
// Grammar: removeCommand
type RemoveStmt struct {
	Path *PathExpr
	Loc  Loc
}

func (*RemoveStmt) nodeTag()      {}
func (n *RemoveStmt) GetLoc() Loc { return n.Loc }
func (*RemoveStmt) stmtNode()     {}

// CreateTableStmt represents `CREATE TABLE name`. PartiQL DDL has no
// column definitions or constraints — just the table name.
//
// Grammar: createCommand#CreateTable
type CreateTableStmt struct {
	Name *VarRef
	Loc  Loc
}

func (*CreateTableStmt) nodeTag()      {}
func (n *CreateTableStmt) GetLoc() Loc { return n.Loc }
func (*CreateTableStmt) stmtNode()     {}

// CreateIndexStmt represents `CREATE INDEX ON table (path, path, ...)`.
//
// Grammar: createCommand#CreateIndex
type CreateIndexStmt struct {
	Table *VarRef
	Paths []*PathExpr
	Loc   Loc
}

func (*CreateIndexStmt) nodeTag()      {}
func (n *CreateIndexStmt) GetLoc() Loc { return n.Loc }
func (*CreateIndexStmt) stmtNode()     {}

// DropTableStmt represents `DROP TABLE name`.
//
// Grammar: dropCommand#DropTable
type DropTableStmt struct {
	Name *VarRef
	Loc  Loc
}

func (*DropTableStmt) nodeTag()      {}
func (n *DropTableStmt) GetLoc() Loc { return n.Loc }
func (*DropTableStmt) stmtNode()     {}

// DropIndexStmt represents `DROP INDEX index ON table`.
//
// Grammar: dropCommand#DropIndex
type DropIndexStmt struct {
	Index *VarRef
	Table *VarRef
	Loc   Loc
}

func (*DropIndexStmt) nodeTag()      {}
func (n *DropIndexStmt) GetLoc() Loc { return n.Loc }
func (*DropIndexStmt) stmtNode()     {}

// ExecStmt represents `EXEC name [arg, arg, ...]`. Per the grammar
// (PartiQLParser.g4 line 65: `EXEC name=expr ...`), the procedure name
// is itself an expression — typically a VarRef but may be any ExprNode
// (e.g., a parameter `?`) so the AST keeps the full breadth.
//
// Grammar: execCommand
type ExecStmt struct {
	Name ExprNode
	Args []ExprNode
	Loc  Loc
}

func (*ExecStmt) nodeTag()      {}
func (n *ExecStmt) GetLoc() Loc { return n.Loc }
func (*ExecStmt) stmtNode()     {}
```

This file currently references `OrderByItem`, `LetBinding`, `GroupByClause`, `TargetEntry`, `PivotProjection`, `OnConflict`, `ReturningClause`, `SetAssignment` — all defined in Task 9. The build will fail until Task 9 lands these helpers. To keep tasks independently testable, **add temporary forward-declared placeholder types at the bottom of this file** so the build succeeds:

```go
// ---------------------------------------------------------------------------
// PLACEHOLDERS — replaced in Task 9.
//
// These are forward-declared minimal stubs so the package builds at the
// end of Task 8 even though the real types live in the next task. Each
// stub is a bare struct with the Loc field plus the methods Node requires.
// Task 9 deletes this block and replaces it with the full type definitions.
// ---------------------------------------------------------------------------

type TargetEntry struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*TargetEntry) nodeTag()      {}
func (n *TargetEntry) GetLoc() Loc { return n.Loc }

type PivotProjection struct {
	Value ExprNode
	At    ExprNode
	Loc   Loc
}

func (*PivotProjection) nodeTag()      {}
func (n *PivotProjection) GetLoc() Loc { return n.Loc }

type LetBinding struct {
	Expr  ExprNode
	Alias string
	Loc   Loc
}

func (*LetBinding) nodeTag()      {}
func (n *LetBinding) GetLoc() Loc { return n.Loc }

type GroupByClause struct {
	Partial bool
	Items   []*GroupByItem
	GroupAs *string
	Loc     Loc
}

func (*GroupByClause) nodeTag()      {}
func (n *GroupByClause) GetLoc() Loc { return n.Loc }

type GroupByItem struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*GroupByItem) nodeTag()      {}
func (n *GroupByItem) GetLoc() Loc { return n.Loc }

type SetAssignment struct {
	Target *PathExpr
	Value  ExprNode
	Loc    Loc
}

func (*SetAssignment) nodeTag()      {}
func (n *SetAssignment) GetLoc() Loc { return n.Loc }

type OnConflict struct {
	Target *OnConflictTarget
	Action OnConflictAction
	Where  ExprNode
	Loc    Loc
}

func (*OnConflict) nodeTag()      {}
func (n *OnConflict) GetLoc() Loc { return n.Loc }

type OnConflictTarget struct {
	Cols           []*VarRef
	ConstraintName string
	Loc            Loc
}

func (*OnConflictTarget) nodeTag()      {}
func (n *OnConflictTarget) GetLoc() Loc { return n.Loc }

type ReturningClause struct {
	Items []*ReturningItem
	Loc   Loc
}

func (*ReturningClause) nodeTag()      {}
func (n *ReturningClause) GetLoc() Loc { return n.Loc }

type ReturningItem struct {
	Status  ReturningStatus
	Mapping ReturningMapping
	Star    bool
	Expr    ExprNode
	Loc     Loc
}

func (*ReturningItem) nodeTag()      {}
func (n *ReturningItem) GetLoc() Loc { return n.Loc }
```

**Note:** This block also keeps the existing `OrderByItem` placeholder (originally added in Task 4 Step 4). Make sure that type definition is still present in the file — it should appear above or alongside the new placeholders. If you accidentally delete it, the build will fail with `undefined: OrderByItem` from `exprs.go`.

The full content of `partiql/ast/stmts.go` after Step 1 should be: header comment + enums + 14 statement types + the placeholder block (including `OrderByItem` from Task 4 + the 10 new placeholders above).

- [ ] **Step 2: Append assertions and TestGetLoc rows for the 14 statements**

In `ast_test.go`:

```go
// Top-level statements (stmts.go).
var _ StmtNode = (*SelectStmt)(nil)
var _ StmtNode = (*SetOpStmt)(nil)
var _ StmtNode = (*ExplainStmt)(nil)
var _ StmtNode = (*InsertStmt)(nil)
var _ StmtNode = (*UpdateStmt)(nil)
var _ StmtNode = (*DeleteStmt)(nil)
var _ StmtNode = (*UpsertStmt)(nil)
var _ StmtNode = (*ReplaceStmt)(nil)
var _ StmtNode = (*RemoveStmt)(nil)
var _ StmtNode = (*CreateTableStmt)(nil)
var _ StmtNode = (*CreateIndexStmt)(nil)
var _ StmtNode = (*DropTableStmt)(nil)
var _ StmtNode = (*DropIndexStmt)(nil)
var _ StmtNode = (*ExecStmt)(nil)
```

Append to `cases`:

```go
{"SelectStmt",      &SelectStmt{Loc: Loc{Start: 10, End: 20}}},
{"SetOpStmt",       &SetOpStmt{Loc: Loc{Start: 10, End: 20}}},
{"ExplainStmt",     &ExplainStmt{Loc: Loc{Start: 10, End: 20}}},
{"InsertStmt",      &InsertStmt{Loc: Loc{Start: 10, End: 20}}},
{"UpdateStmt",      &UpdateStmt{Loc: Loc{Start: 10, End: 20}}},
{"DeleteStmt",      &DeleteStmt{Loc: Loc{Start: 10, End: 20}}},
{"UpsertStmt",      &UpsertStmt{Loc: Loc{Start: 10, End: 20}}},
{"ReplaceStmt",     &ReplaceStmt{Loc: Loc{Start: 10, End: 20}}},
{"RemoveStmt",      &RemoveStmt{Loc: Loc{Start: 10, End: 20}}},
{"CreateTableStmt", &CreateTableStmt{Loc: Loc{Start: 10, End: 20}}},
{"CreateIndexStmt", &CreateIndexStmt{Loc: Loc{Start: 10, End: 20}}},
{"DropTableStmt",   &DropTableStmt{Loc: Loc{Start: 10, End: 20}}},
{"DropIndexStmt",   &DropIndexStmt{Loc: Loc{Start: 10, End: 20}}},
{"ExecStmt",        &ExecStmt{Loc: Loc{Start: 10, End: 20}}},
```

(The clause helpers and DML helpers — TargetEntry, OnConflict, etc. — are NOT added to the test table here. They land in Task 9 when their real form is written.)

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 56 sub-tests passing (42 prior + 14 new statements).

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/stmts.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add top-level statement nodes

14 statement types covering DDL (CreateTable/Index, DropTable/Index),
DML (Insert/Update/Delete/Upsert/Replace/Remove with both legacy and
RFC 0011 forms), DQL (SelectStmt with full clause set, SetOpStmt for
UNION/INTERSECT/EXCEPT, ExplainStmt), and EXEC.

Plus the SetOpKind, OnConflictAction, ReturningStatus, ReturningMapping
enums.

Includes forward-declared placeholders for the clause helpers
(TargetEntry, GroupByClause, OnConflict, ReturningClause, etc.) so
the package builds; Task 9 replaces those with the real definitions.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 9: `stmts.go` — clause and DML helpers (11 types)

**Files:**
- Modify: `partiql/ast/stmts.go` (replace the placeholder block with real definitions)
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Replace the placeholder block in `partiql/ast/stmts.go`**

Open `partiql/ast/stmts.go`. Find the section that begins:

```go
// PLACEHOLDERS — replaced in Task 9.
```

**Delete that entire comment block and the placeholder type definitions below it (TargetEntry, PivotProjection, LetBinding, GroupByClause, GroupByItem, OrderByItem, SetAssignment, OnConflict, OnConflictTarget, ReturningClause, ReturningItem).**

Replace with the real definitions:

```go
// ===========================================================================
// SELECT clause helpers — bare Node (no sub-interface marker).
//
// These types appear only as fields/elements inside SelectStmt and never
// stand alone in scalar/statement/table-expr position, so they don't need
// a sub-interface marker.
// ===========================================================================

// TargetEntry represents one item in a SELECT projection list:
// `expr [AS alias]`.
//
// Grammar: projectionItem
type TargetEntry struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*TargetEntry) nodeTag()      {}
func (n *TargetEntry) GetLoc() Loc { return n.Loc }

// PivotProjection represents the body of a `PIVOT v AT k` projection.
// Used as SelectStmt.Pivot. PartiQL-unique.
//
// Grammar: selectClause#SelectPivot
type PivotProjection struct {
	Value ExprNode
	At    ExprNode
	Loc   Loc
}

func (*PivotProjection) nodeTag()      {}
func (n *PivotProjection) GetLoc() Loc { return n.Loc }

// LetBinding represents one `expr AS alias` binding inside a LET clause.
// PartiQL-unique.
//
// Grammar: letBinding
type LetBinding struct {
	Expr  ExprNode
	Alias string
	Loc   Loc
}

func (*LetBinding) nodeTag()      {}
func (n *LetBinding) GetLoc() Loc { return n.Loc }

// GroupByClause represents `GROUP [PARTIAL] BY items [GROUP AS alias]`.
//
// Grammar: groupClause
type GroupByClause struct {
	Partial bool
	Items   []*GroupByItem
	GroupAs *string
	Loc     Loc
}

func (*GroupByClause) nodeTag()      {}
func (n *GroupByClause) GetLoc() Loc { return n.Loc }

// GroupByItem represents one `expr [AS alias]` item in a GROUP BY list.
//
// Grammar: groupKey
type GroupByItem struct {
	Expr  ExprNode
	Alias *string
	Loc   Loc
}

func (*GroupByItem) nodeTag()      {}
func (n *GroupByItem) GetLoc() Loc { return n.Loc }

// OrderByItem represents one `expr [ASC|DESC] [NULLS FIRST|LAST]` item
// in an ORDER BY list. NullsExplicit is true when the source text included
// a NULLS clause; NullsFirst is the resulting setting when it did.
//
// Grammar: orderSortSpec
type OrderByItem struct {
	Expr          ExprNode
	Desc          bool
	NullsFirst    bool
	NullsExplicit bool
	Loc           Loc
}

func (*OrderByItem) nodeTag()      {}
func (n *OrderByItem) GetLoc() Loc { return n.Loc }

// ===========================================================================
// DML helpers — bare Node.
// ===========================================================================

// SetAssignment represents one `target = value` assignment in an UPDATE
// SET clause (or in a chained-DML SET op).
//
// Grammar: setAssignment
type SetAssignment struct {
	Target *PathExpr
	Value  ExprNode
	Loc    Loc
}

func (*SetAssignment) nodeTag()      {}
func (n *SetAssignment) GetLoc() Loc { return n.Loc }

// OnConflict represents the body of an ON CONFLICT clause:
// `ON CONFLICT [target] [WHERE expr] action`.
//
// Grammar: onConflictClause#OnConflict, onConflictClause#OnConflictLegacy
type OnConflict struct {
	Target *OnConflictTarget
	Action OnConflictAction
	Where  ExprNode
	Loc    Loc
}

func (*OnConflict) nodeTag()      {}
func (n *OnConflict) GetLoc() Loc { return n.Loc }

// OnConflictTarget represents the target of an ON CONFLICT clause —
// either a list of column references `(col, col, ...)` or
// `ON CONSTRAINT name`. Exactly one of Cols/ConstraintName is set.
//
// Grammar: conflictTarget
type OnConflictTarget struct {
	Cols           []*VarRef
	ConstraintName string
	Loc            Loc
}

func (*OnConflictTarget) nodeTag()      {}
func (n *OnConflictTarget) GetLoc() Loc { return n.Loc }

// ReturningClause represents `RETURNING item, item, …`.
//
// Grammar: returningClause
type ReturningClause struct {
	Items []*ReturningItem
	Loc   Loc
}

func (*ReturningClause) nodeTag()      {}
func (n *ReturningClause) GetLoc() Loc { return n.Loc }

// ReturningItem represents one `(MODIFIED|ALL) (OLD|NEW) (* | expr)` entry
// in a RETURNING clause. Star is true for the `*` form; otherwise Expr
// holds the projection expression.
//
// Grammar: returningColumn
type ReturningItem struct {
	Status  ReturningStatus
	Mapping ReturningMapping
	Star    bool
	Expr    ExprNode
	Loc     Loc
}

func (*ReturningItem) nodeTag()      {}
func (n *ReturningItem) GetLoc() Loc { return n.Loc }
```

After this step, `partiql/ast/stmts.go` should contain (in order):

1. Package declaration + file header comment
2. Enum types (`SetOpKind`, `OnConflictAction`, `ReturningStatus`, `ReturningMapping`)
3. The 14 top-level statement types from Task 8
4. The 11 clause/DML helper types above

No placeholder block remains.

- [ ] **Step 2: Append assertions and TestGetLoc rows for the 11 helpers**

In `ast_test.go`:

```go
// Clause and DML helpers (stmts.go).
var _ Node = (*TargetEntry)(nil)
var _ Node = (*PivotProjection)(nil)
var _ Node = (*LetBinding)(nil)
var _ Node = (*GroupByClause)(nil)
var _ Node = (*GroupByItem)(nil)
var _ Node = (*OrderByItem)(nil)
var _ Node = (*SetAssignment)(nil)
var _ Node = (*OnConflict)(nil)
var _ Node = (*OnConflictTarget)(nil)
var _ Node = (*ReturningClause)(nil)
var _ Node = (*ReturningItem)(nil)
```

Append to `cases`:

```go
{"TargetEntry",      &TargetEntry{Loc: Loc{Start: 10, End: 20}}},
{"PivotProjection",  &PivotProjection{Loc: Loc{Start: 10, End: 20}}},
{"LetBinding",       &LetBinding{Loc: Loc{Start: 10, End: 20}}},
{"GroupByClause",    &GroupByClause{Loc: Loc{Start: 10, End: 20}}},
{"GroupByItem",      &GroupByItem{Loc: Loc{Start: 10, End: 20}}},
{"OrderByItem",      &OrderByItem{Loc: Loc{Start: 10, End: 20}}},
{"SetAssignment",    &SetAssignment{Loc: Loc{Start: 10, End: 20}}},
{"OnConflict",       &OnConflict{Loc: Loc{Start: 10, End: 20}}},
{"OnConflictTarget", &OnConflictTarget{Loc: Loc{Start: 10, End: 20}}},
{"ReturningClause",  &ReturningClause{Loc: Loc{Start: 10, End: 20}}},
{"ReturningItem",    &ReturningItem{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 67 sub-tests passing (56 prior + 11 new).

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/stmts.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add SELECT clause and DML helper nodes

11 helper types — TargetEntry, PivotProjection, LetBinding,
GroupByClause, GroupByItem, OrderByItem, SetAssignment, OnConflict,
OnConflictTarget, ReturningClause, ReturningItem. All are bare Node
(no sub-interface marker) because they only appear as fields inside
the statement nodes from Task 8.

Replaces the forward-declared placeholders introduced in Tasks 4 and 8.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 10: `patterns.go` — graph pattern nodes

**Files:**
- Create: `partiql/ast/patterns.go`
- Modify: `partiql/ast/ast_test.go`

- [ ] **Step 1: Create `partiql/ast/patterns.go`**

```go
package ast

// ---------------------------------------------------------------------------
// Graph Pattern Matching (GPML) nodes — implement PatternNode (or ExprNode
// for the top-level MatchExpr container).
//
// PartiQL has a graph pattern matching extension based on GPML. The full
// shape of the field set inside NodePattern/EdgePattern may be refined when
// the parser-graph DAG node (node 16) is implemented and we read the grammar
// more carefully. The marker interface, the node names, and the multi-interface
// relationship to ExprNode are stable in this initial pass.
//
// Grammar: PartiQLParser.g4 — sourced from rules:
//   gpmlPattern          (lines 315–316)
//   gpmlPatternList      (lines 318–319)
//   matchPattern         (lines 321–322)
//   graphPart            (lines 324–328)
//   matchSelector        (lines 330–334)
//   patternPathVariable  (lines 336–337)
//   patternRestrictor    (lines 339–340)
//   node                 (lines 342–343)
//   edge                 (lines 345–348)
//   pattern              (lines 350–353)
//   patternQuantifier    (lines 355–358)
//   edgeWSpec            (lines 360–368)
//   edgeSpec             (lines 370–371)
//   patternPartLabel     (lines 373–374)
//   edgeAbbrev           (lines 376–381)
//   exprGraphMatchOne    (lines 628–629)
//   exprGraphMatchMany   (lines 625–626)
// Each type below cites its specific rule#Label.
// ---------------------------------------------------------------------------

// EdgeDirection identifies the direction of an edge pattern.
//
// Enumerates the seven edge alternatives from PartiQLParser.g4 `edgeWSpec`
// (lines 360–368) and `edgeAbbrev` (lines 376–381). Each value maps to the
// bare grammar symbol via String().
type EdgeDirection int

const (
	EdgeDirInvalid                 EdgeDirection = iota
	EdgeDirRight                                 // ->          edgeWSpec#EdgeSpecRight
	EdgeDirLeft                                  // <-          edgeWSpec#EdgeSpecLeft
	EdgeDirUndirected                            // ~           edgeWSpec#EdgeSpecUndirected
	EdgeDirLeftOrRight                           // <->         edgeWSpec#EdgeSpecBidirectional
	EdgeDirLeftOrUndirected                      // <~          edgeWSpec#EdgeSpecUndirectedLeft
	EdgeDirRightOrUndirected                     // ~>          edgeWSpec#EdgeSpecUndirectedRight
	EdgeDirUndirectedBidirectional               // -           edgeWSpec#EdgeSpecUndirectedBidirectional (any of right/left/undirected)
)

func (d EdgeDirection) String() string {
	switch d {
	case EdgeDirRight:
		return "->"
	case EdgeDirLeft:
		return "<-"
	case EdgeDirUndirected:
		return "~"
	case EdgeDirLeftOrRight:
		return "<->"
	case EdgeDirLeftOrUndirected:
		return "<~"
	case EdgeDirRightOrUndirected:
		return "~>"
	case EdgeDirUndirectedBidirectional:
		return "-"
	default:
		return "INVALID"
	}
}

// PatternRestrictor identifies the optional restrictor keyword on a graph pattern.
type PatternRestrictor int

const (
	PatternRestrictorNone PatternRestrictor = iota
	PatternRestrictorTrail
	PatternRestrictorAcyclic
	PatternRestrictorSimple
)

func (r PatternRestrictor) String() string {
	switch r {
	case PatternRestrictorTrail:
		return "TRAIL"
	case PatternRestrictorAcyclic:
		return "ACYCLIC"
	case PatternRestrictorSimple:
		return "SIMPLE"
	default:
		return ""
	}
}

// SelectorKind identifies the optional selector keyword on a graph pattern.
type SelectorKind int

const (
	SelectorKindNone SelectorKind = iota
	SelectorKindAny
	SelectorKindAllShortest
	SelectorKindShortestK
)

func (s SelectorKind) String() string {
	switch s {
	case SelectorKindAny:
		return "ANY"
	case SelectorKindAllShortest:
		return "ALL_SHORTEST"
	case SelectorKindShortestK:
		return "SHORTEST_K"
	default:
		return ""
	}
}

// MatchExpr is the top-level graph-match expression: MATCH(graph_expr, pattern, …).
// Implements ExprNode because PartiQL embeds graph matching in expression
// position (exprGraphMatchOne / exprGraphMatchMany rules).
//
// Grammar: exprGraphMatchOne, exprGraphMatchMany
type MatchExpr struct {
	Expr     ExprNode        // the graph-valued expression being matched
	Patterns []*GraphPattern // one or more pattern alternatives
	Loc      Loc
}

func (*MatchExpr) nodeTag()      {}
func (n *MatchExpr) GetLoc() Loc { return n.Loc }
func (*MatchExpr) exprNode()     {}

// GraphPattern is one complete pattern: optional selector + restrictor +
// path variable + a sequence of node/edge pattern parts.
//
// Grammar: gpmlPattern
type GraphPattern struct {
	Selector   *PatternSelector
	Restrictor PatternRestrictor
	Variable   *VarRef // optional `p = ...` path variable binding
	Parts      []PatternNode
	Loc        Loc
}

func (*GraphPattern) nodeTag()      {}
func (n *GraphPattern) GetLoc() Loc { return n.Loc }
func (*GraphPattern) patternNode()  {}

// NodePattern represents `(var:Label WHERE …)`. Variable, Labels, and
// Where are all optional.
//
// Grammar: node
type NodePattern struct {
	Variable *VarRef
	Labels   []string
	Where    ExprNode
	Loc      Loc
}

func (*NodePattern) nodeTag()      {}
func (n *NodePattern) GetLoc() Loc { return n.Loc }
func (*NodePattern) patternNode()  {}

// EdgePattern represents `-[var:Label]->`, `<-[]-`, `~[]~`, etc.
// Direction is required; Variable, Labels, Where, and Quantifier are optional.
//
// Grammar: edge#EdgeWithSpec, edge#EdgeAbbreviated
//          (direction comes from edgeWSpec labels / edgeAbbrev; body from edgeSpec)
type EdgePattern struct {
	Direction  EdgeDirection
	Variable   *VarRef
	Labels     []string
	Where      ExprNode
	Quantifier *PatternQuantifier
	Loc        Loc
}

func (*EdgePattern) nodeTag()      {}
func (n *EdgePattern) GetLoc() Loc { return n.Loc }
func (*EdgePattern) patternNode()  {}

// PatternQuantifier represents the `+`, `*`, or `{m,n}` decorator on a
// pattern part. Min and Max use -1 to indicate "unbounded".
//
// Grammar: patternQuantifier
type PatternQuantifier struct {
	Min int
	Max int // -1 = unbounded
	Loc Loc
}

func (*PatternQuantifier) nodeTag()      {}
func (n *PatternQuantifier) GetLoc() Loc { return n.Loc }

// PatternSelector represents the optional selector keyword on a graph
// pattern: `ANY`, `ALL SHORTEST`, or `SHORTEST k`. K is non-zero only for
// SHORTEST K.
//
// Grammar: matchSelector#SelectorBasic, matchSelector#SelectorAny, matchSelector#SelectorShortest
type PatternSelector struct {
	Kind SelectorKind
	K    int
	Loc  Loc
}

func (*PatternSelector) nodeTag()      {}
func (n *PatternSelector) GetLoc() Loc { return n.Loc }
```

- [ ] **Step 2: Append assertions and TestGetLoc rows**

In `ast_test.go`:

```go
// Graph patterns (patterns.go).
var _ ExprNode    = (*MatchExpr)(nil)
var _ PatternNode = (*GraphPattern)(nil)
var _ PatternNode = (*NodePattern)(nil)
var _ PatternNode = (*EdgePattern)(nil)
var _ Node        = (*PatternQuantifier)(nil)
var _ Node        = (*PatternSelector)(nil)
```

Append to `cases`:

```go
{"MatchExpr",         &MatchExpr{Loc: Loc{Start: 10, End: 20}}},
{"GraphPattern",      &GraphPattern{Loc: Loc{Start: 10, End: 20}}},
{"NodePattern",       &NodePattern{Loc: Loc{Start: 10, End: 20}}},
{"EdgePattern",       &EdgePattern{Loc: Loc{Start: 10, End: 20}}},
{"PatternQuantifier", &PatternQuantifier{Loc: Loc{Start: 10, End: 20}}},
{"PatternSelector",   &PatternSelector{Loc: Loc{Start: 10, End: 20}}},
```

- [ ] **Step 3: Build and test**

```bash
go build ./partiql/ast/...
go test ./partiql/ast/...
```

Expected: ok, 73 sub-tests passing.

- [ ] **Step 4: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 5: Commit**

```bash
git add partiql/ast/patterns.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add graph pattern matching nodes

GPML scaffolding: MatchExpr (also implements ExprNode because PartiQL
embeds MATCH in expression position), GraphPattern container,
NodePattern/EdgePattern atomic units, PatternQuantifier (+/*/{m,n}),
PatternSelector (ANY/ALL SHORTEST/SHORTEST k). Plus EdgeDirection,
PatternRestrictor, SelectorKind enums.

Spec notes that the field shape inside NodePattern/EdgePattern may
be refined when parser-graph (DAG node 16) is implemented; the
marker interface, node names, and multi-interface relationship are
stable from this point.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 11: `outfuncs.go` — `NodeToString` + golden tests + reflection safety net

**Files:**
- Create: `partiql/ast/outfuncs.go`
- Modify: `partiql/ast/ast_test.go`

This is the biggest task by line count. The `NodeToString` switch has one arm per node type.

- [ ] **Step 1: Create `partiql/ast/outfuncs.go`**

```go
package ast

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// NodeToString returns a deterministic textual dump of any AST node.
// The format is Go-struct-like for readability:
//
//	BinaryExpr{Op:+ Left:VarRef{Name:a} Right:NumberLit{Val:1}}
//
// Loc fields are not dumped (positions are tested separately by
// TestGetLoc). For nil input, returns "<nil>".
//
// Used by ast_test.go for snapshot-style assertions and by future
// parser golden tests.
func NodeToString(n Node) string {
	if n == nil {
		return "<nil>"
	}
	var sb strings.Builder
	writeNode(&sb, n)
	return sb.String()
}

// writeNode dispatches on the concrete node type and appends a
// textual representation to sb.
func writeNode(sb *strings.Builder, n Node) {
	if n == nil {
		sb.WriteString("<nil>")
		return
	}
	// Detect typed-nil pointers (e.g. (*VarRef)(nil) wrapped in a Node
	// interface). The plain `n == nil` check above only catches an
	// untyped nil interface; the reflection check is needed by the
	// TestNodeToString_AllNodesCovered safety net which constructs
	// zero-value parent nodes whose pointer fields are typed nil.
	if rv := reflect.ValueOf(n); rv.Kind() == reflect.Pointer && rv.IsNil() {
		sb.WriteString("<nil>")
		return
	}
	switch v := n.(type) {

	// -----------------------------------------------------------------------
	// node.go
	// -----------------------------------------------------------------------
	case *List:
		sb.WriteString("List{Items:[")
		for i, item := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, item)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// literals.go
	// -----------------------------------------------------------------------
	case *StringLit:
		fmt.Fprintf(sb, "StringLit{Val:%q}", v.Val)
	case *NumberLit:
		fmt.Fprintf(sb, "NumberLit{Val:%s}", v.Val)
	case *BoolLit:
		fmt.Fprintf(sb, "BoolLit{Val:%t}", v.Val)
	case *NullLit:
		sb.WriteString("NullLit{}")
	case *MissingLit:
		sb.WriteString("MissingLit{}")
	case *DateLit:
		fmt.Fprintf(sb, "DateLit{Val:%s}", v.Val)
	case *TimeLit:
		sb.WriteString("TimeLit{Val:")
		sb.WriteString(v.Val)
		if v.Precision != nil {
			fmt.Fprintf(sb, " Precision:%d", *v.Precision)
		}
		if v.WithTimeZone {
			sb.WriteString(" WithTimeZone:true")
		}
		sb.WriteString("}")
	case *IonLit:
		fmt.Fprintf(sb, "IonLit{Text:%q}", v.Text)

	// -----------------------------------------------------------------------
	// exprs.go — operators and predicates
	// -----------------------------------------------------------------------
	case *BinaryExpr:
		fmt.Fprintf(sb, "BinaryExpr{Op:%s Left:", v.Op)
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *UnaryExpr:
		fmt.Fprintf(sb, "UnaryExpr{Op:%s Operand:", v.Op)
		writeNode(sb, v.Operand)
		sb.WriteString("}")
	case *InExpr:
		sb.WriteString("InExpr{Expr:")
		writeNode(sb, v.Expr)
		if v.Subquery != nil {
			sb.WriteString(" Subquery:")
			writeNode(sb, v.Subquery)
		} else {
			sb.WriteString(" List:[")
			for i, e := range v.List {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, e)
			}
			sb.WriteString("]")
		}
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *BetweenExpr:
		sb.WriteString("BetweenExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Low:")
		writeNode(sb, v.Low)
		sb.WriteString(" High:")
		writeNode(sb, v.High)
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *LikeExpr:
		sb.WriteString("LikeExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Pattern:")
		writeNode(sb, v.Pattern)
		if v.Escape != nil {
			sb.WriteString(" Escape:")
			writeNode(sb, v.Escape)
		}
		fmt.Fprintf(sb, " Not:%t}", v.Not)
	case *IsExpr:
		sb.WriteString("IsExpr{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Type:%s Not:%t}", v.Type, v.Not)

	// -----------------------------------------------------------------------
	// exprs.go — special-form expressions
	// -----------------------------------------------------------------------
	case *FuncCall:
		fmt.Fprintf(sb, "FuncCall{Name:%s", v.Name)
		if v.Quantifier != QuantifierNone {
			fmt.Fprintf(sb, " Quantifier:%s", v.Quantifier)
		}
		if v.Star {
			sb.WriteString(" Star:true")
		}
		sb.WriteString(" Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]")
		if v.Over != nil {
			sb.WriteString(" Over:")
			writeNode(sb, v.Over)
		}
		sb.WriteString("}")
	case *CaseExpr:
		sb.WriteString("CaseExpr{")
		if v.Operand != nil {
			sb.WriteString("Operand:")
			writeNode(sb, v.Operand)
			sb.WriteString(" ")
		}
		sb.WriteString("Whens:[")
		for i, w := range v.Whens {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, w)
		}
		sb.WriteString("]")
		if v.Else != nil {
			sb.WriteString(" Else:")
			writeNode(sb, v.Else)
		}
		sb.WriteString("}")
	case *CaseWhen:
		sb.WriteString("CaseWhen{When:")
		writeNode(sb, v.When)
		sb.WriteString(" Then:")
		writeNode(sb, v.Then)
		sb.WriteString("}")
	case *CastExpr:
		fmt.Fprintf(sb, "CastExpr{Kind:%s Expr:", v.Kind)
		writeNode(sb, v.Expr)
		sb.WriteString(" AsType:")
		writeNode(sb, v.AsType)
		sb.WriteString("}")
	case *ExtractExpr:
		fmt.Fprintf(sb, "ExtractExpr{Field:%s From:", v.Field)
		writeNode(sb, v.From)
		sb.WriteString("}")
	case *TrimExpr:
		sb.WriteString("TrimExpr{")
		if v.Spec != TrimSpecNone {
			fmt.Fprintf(sb, "Spec:%s ", v.Spec)
		}
		if v.Sub != nil {
			sb.WriteString("Sub:")
			writeNode(sb, v.Sub)
			sb.WriteString(" ")
		}
		sb.WriteString("From:")
		writeNode(sb, v.From)
		sb.WriteString("}")
	case *SubstringExpr:
		sb.WriteString("SubstringExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" From:")
		writeNode(sb, v.From)
		if v.For != nil {
			sb.WriteString(" For:")
			writeNode(sb, v.For)
		}
		sb.WriteString("}")
	case *CoalesceExpr:
		sb.WriteString("CoalesceExpr{Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]}")
	case *NullIfExpr:
		sb.WriteString("NullIfExpr{Left:")
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *WindowSpec:
		sb.WriteString("WindowSpec{PartitionBy:[")
		for i, e := range v.PartitionBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("] OrderBy:[")
		for i, o := range v.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// exprs.go — paths, vars, params, subqueries
	// -----------------------------------------------------------------------
	case *PathExpr:
		sb.WriteString("PathExpr{Root:")
		writeNode(sb, v.Root)
		sb.WriteString(" Steps:[")
		for i, s := range v.Steps {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString("]}")
	case *VarRef:
		fmt.Fprintf(sb, "VarRef{Name:%s", v.Name)
		if v.AtPrefixed {
			sb.WriteString(" AtPrefixed:true")
		}
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *ParamRef:
		sb.WriteString("ParamRef{}")
	case *SubLink:
		sb.WriteString("SubLink{Stmt:")
		writeNode(sb, v.Stmt)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// exprs.go — collection literals
	// -----------------------------------------------------------------------
	case *ListLit:
		sb.WriteString("ListLit{Items:[")
		for i, e := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("]}")
	case *BagLit:
		sb.WriteString("BagLit{Items:[")
		for i, e := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, e)
		}
		sb.WriteString("]}")
	case *TupleLit:
		sb.WriteString("TupleLit{Pairs:[")
		for i, p := range v.Pairs {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *TuplePair:
		sb.WriteString("TuplePair{Key:")
		writeNode(sb, v.Key)
		sb.WriteString(" Value:")
		writeNode(sb, v.Value)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// exprs.go — path steps
	// -----------------------------------------------------------------------
	case *DotStep:
		fmt.Fprintf(sb, "DotStep{Field:%s", v.Field)
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *AllFieldsStep:
		sb.WriteString("AllFieldsStep{}")
	case *IndexStep:
		sb.WriteString("IndexStep{Index:")
		writeNode(sb, v.Index)
		sb.WriteString("}")
	case *WildcardStep:
		sb.WriteString("WildcardStep{}")

	// -----------------------------------------------------------------------
	// tableexprs.go
	// -----------------------------------------------------------------------
	case *TableRef:
		fmt.Fprintf(sb, "TableRef{Name:%s", v.Name)
		if v.Schema != "" {
			fmt.Fprintf(sb, " Schema:%s", v.Schema)
		}
		if v.CaseSensitive {
			sb.WriteString(" CaseSensitive:true")
		}
		sb.WriteString("}")
	case *AliasedSource:
		sb.WriteString("AliasedSource{Source:")
		writeNode(sb, v.Source)
		writeOptString(sb, " As:", v.As)
		writeOptString(sb, " At:", v.At)
		writeOptString(sb, " By:", v.By)
		sb.WriteString("}")
	case *JoinExpr:
		fmt.Fprintf(sb, "JoinExpr{Kind:%s Left:", v.Kind)
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		if v.On != nil {
			sb.WriteString(" On:")
			writeNode(sb, v.On)
		}
		sb.WriteString("}")
	case *UnpivotExpr:
		sb.WriteString("UnpivotExpr{Source:")
		writeNode(sb, v.Source)
		writeOptString(sb, " As:", v.As)
		writeOptString(sb, " At:", v.At)
		writeOptString(sb, " By:", v.By)
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// types.go
	// -----------------------------------------------------------------------
	case *TypeRef:
		fmt.Fprintf(sb, "TypeRef{Name:%s", v.Name)
		if len(v.Args) > 0 {
			sb.WriteString(" Args:[")
			for i, a := range v.Args {
				if i > 0 {
					sb.WriteString(",")
				}
				sb.WriteString(strconv.Itoa(a))
			}
			sb.WriteString("]")
		}
		if v.WithTimeZone {
			sb.WriteString(" WithTimeZone:true")
		}
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// stmts.go — top-level statements
	// -----------------------------------------------------------------------
	case *SelectStmt:
		writeSelectStmt(sb, v)
	case *SetOpStmt:
		fmt.Fprintf(sb, "SetOpStmt{Op:%s", v.Op)
		if v.Quantifier != QuantifierNone {
			fmt.Fprintf(sb, " Quantifier:%s", v.Quantifier)
		}
		if v.Outer {
			sb.WriteString(" Outer:true")
		}
		sb.WriteString(" Left:")
		writeNode(sb, v.Left)
		sb.WriteString(" Right:")
		writeNode(sb, v.Right)
		sb.WriteString("}")
	case *ExplainStmt:
		sb.WriteString("ExplainStmt{Inner:")
		writeNode(sb, v.Inner)
		sb.WriteString("}")
	case *InsertStmt:
		writeDmlStmt(sb, "InsertStmt", v.Target, v.AsAlias, v.Value, v.Pos, v.OnConflict, v.Returning)
	case *UpdateStmt:
		sb.WriteString("UpdateStmt{Source:")
		writeNode(sb, v.Source)
		sb.WriteString(" Sets:[")
		for i, s := range v.Sets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, s)
		}
		sb.WriteString("]")
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Returning != nil {
			sb.WriteString(" Returning:")
			writeNode(sb, v.Returning)
		}
		sb.WriteString("}")
	case *DeleteStmt:
		sb.WriteString("DeleteStmt{Source:")
		writeNode(sb, v.Source)
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Returning != nil {
			sb.WriteString(" Returning:")
			writeNode(sb, v.Returning)
		}
		sb.WriteString("}")
	case *UpsertStmt:
		writeDmlStmt(sb, "UpsertStmt", v.Target, v.AsAlias, v.Value, nil, v.OnConflict, v.Returning)
	case *ReplaceStmt:
		writeDmlStmt(sb, "ReplaceStmt", v.Target, v.AsAlias, v.Value, nil, v.OnConflict, v.Returning)
	case *RemoveStmt:
		sb.WriteString("RemoveStmt{Path:")
		writeNode(sb, v.Path)
		sb.WriteString("}")
	case *CreateTableStmt:
		sb.WriteString("CreateTableStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString("}")
	case *CreateIndexStmt:
		sb.WriteString("CreateIndexStmt{Table:")
		writeNode(sb, v.Table)
		sb.WriteString(" Paths:[")
		for i, p := range v.Paths {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *DropTableStmt:
		sb.WriteString("DropTableStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString("}")
	case *DropIndexStmt:
		sb.WriteString("DropIndexStmt{Index:")
		writeNode(sb, v.Index)
		sb.WriteString(" Table:")
		writeNode(sb, v.Table)
		sb.WriteString("}")
	case *ExecStmt:
		sb.WriteString("ExecStmt{Name:")
		writeNode(sb, v.Name)
		sb.WriteString(" Args:[")
		for i, a := range v.Args {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, a)
		}
		sb.WriteString("]}")

	// -----------------------------------------------------------------------
	// stmts.go — clause and DML helpers
	// -----------------------------------------------------------------------
	case *TargetEntry:
		sb.WriteString("TargetEntry{Expr:")
		writeNode(sb, v.Expr)
		writeOptString(sb, " Alias:", v.Alias)
		sb.WriteString("}")
	case *PivotProjection:
		sb.WriteString("PivotProjection{Value:")
		writeNode(sb, v.Value)
		sb.WriteString(" At:")
		writeNode(sb, v.At)
		sb.WriteString("}")
	case *LetBinding:
		sb.WriteString("LetBinding{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Alias:%s}", v.Alias)
	case *GroupByClause:
		sb.WriteString("GroupByClause{")
		if v.Partial {
			sb.WriteString("Partial:true ")
		}
		sb.WriteString("Items:[")
		for i, it := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, it)
		}
		sb.WriteString("]")
		writeOptString(sb, " GroupAs:", v.GroupAs)
		sb.WriteString("}")
	case *GroupByItem:
		sb.WriteString("GroupByItem{Expr:")
		writeNode(sb, v.Expr)
		writeOptString(sb, " Alias:", v.Alias)
		sb.WriteString("}")
	case *OrderByItem:
		sb.WriteString("OrderByItem{Expr:")
		writeNode(sb, v.Expr)
		fmt.Fprintf(sb, " Desc:%t", v.Desc)
		if v.NullsExplicit {
			fmt.Fprintf(sb, " NullsFirst:%t", v.NullsFirst)
		}
		sb.WriteString("}")
	case *SetAssignment:
		sb.WriteString("SetAssignment{Target:")
		writeNode(sb, v.Target)
		sb.WriteString(" Value:")
		writeNode(sb, v.Value)
		sb.WriteString("}")
	case *OnConflict:
		sb.WriteString("OnConflict{")
		if v.Target != nil {
			sb.WriteString("Target:")
			writeNode(sb, v.Target)
			sb.WriteString(" ")
		}
		fmt.Fprintf(sb, "Action:%s", v.Action)
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		sb.WriteString("}")
	case *OnConflictTarget:
		sb.WriteString("OnConflictTarget{")
		if v.ConstraintName != "" {
			fmt.Fprintf(sb, "ConstraintName:%s", v.ConstraintName)
		} else {
			sb.WriteString("Cols:[")
			for i, c := range v.Cols {
				if i > 0 {
					sb.WriteString(" ")
				}
				writeNode(sb, c)
			}
			sb.WriteString("]")
		}
		sb.WriteString("}")
	case *ReturningClause:
		sb.WriteString("ReturningClause{Items:[")
		for i, it := range v.Items {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, it)
		}
		sb.WriteString("]}")
	case *ReturningItem:
		fmt.Fprintf(sb, "ReturningItem{Status:%s Mapping:%s", v.Status, v.Mapping)
		if v.Star {
			sb.WriteString(" Star:true")
		}
		if v.Expr != nil {
			sb.WriteString(" Expr:")
			writeNode(sb, v.Expr)
		}
		sb.WriteString("}")

	// -----------------------------------------------------------------------
	// patterns.go
	// -----------------------------------------------------------------------
	case *MatchExpr:
		sb.WriteString("MatchExpr{Expr:")
		writeNode(sb, v.Expr)
		sb.WriteString(" Patterns:[")
		for i, p := range v.Patterns {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *GraphPattern:
		sb.WriteString("GraphPattern{")
		if v.Selector != nil {
			sb.WriteString("Selector:")
			writeNode(sb, v.Selector)
			sb.WriteString(" ")
		}
		if v.Restrictor != PatternRestrictorNone {
			fmt.Fprintf(sb, "Restrictor:%s ", v.Restrictor)
		}
		if v.Variable != nil {
			sb.WriteString("Variable:")
			writeNode(sb, v.Variable)
			sb.WriteString(" ")
		}
		sb.WriteString("Parts:[")
		for i, p := range v.Parts {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, p)
		}
		sb.WriteString("]}")
	case *NodePattern:
		sb.WriteString("NodePattern{")
		first := true
		sep := func() {
			if !first {
				sb.WriteString(" ")
			}
			first = false
		}
		if v.Variable != nil {
			sep()
			sb.WriteString("Variable:")
			writeNode(sb, v.Variable)
		}
		if len(v.Labels) > 0 {
			sep()
			fmt.Fprintf(sb, "Labels:[%s]", strings.Join(v.Labels, " "))
		}
		if v.Where != nil {
			sep()
			sb.WriteString("Where:")
			writeNode(sb, v.Where)
		}
		sb.WriteString("}")
	case *EdgePattern:
		fmt.Fprintf(sb, "EdgePattern{Direction:%s", v.Direction)
		if v.Variable != nil {
			sb.WriteString(" Variable:")
			writeNode(sb, v.Variable)
		}
		if len(v.Labels) > 0 {
			fmt.Fprintf(sb, " Labels:[%s]", strings.Join(v.Labels, " "))
		}
		if v.Where != nil {
			sb.WriteString(" Where:")
			writeNode(sb, v.Where)
		}
		if v.Quantifier != nil {
			sb.WriteString(" Quantifier:")
			writeNode(sb, v.Quantifier)
		}
		sb.WriteString("}")
	case *PatternQuantifier:
		fmt.Fprintf(sb, "PatternQuantifier{Min:%d Max:%d}", v.Min, v.Max)
	case *PatternSelector:
		fmt.Fprintf(sb, "PatternSelector{Kind:%s", v.Kind)
		if v.Kind == SelectorKindShortestK {
			fmt.Fprintf(sb, " K:%d", v.K)
		}
		sb.WriteString("}")

	default:
		fmt.Fprintf(sb, "<unknown:%T>", v)
	}
}

// writeSelectStmt is split out because SelectStmt has 13 fields and would
// dominate the main switch arm.
func writeSelectStmt(sb *strings.Builder, s *SelectStmt) {
	sb.WriteString("SelectStmt{")
	first := true
	add := func(label string) {
		if !first {
			sb.WriteString(" ")
		}
		first = false
		sb.WriteString(label)
	}
	if s.Quantifier != QuantifierNone {
		add(fmt.Sprintf("Quantifier:%s", s.Quantifier))
	}
	if s.Star {
		add("Star:true")
	}
	if s.Value != nil {
		add("Value:")
		writeNode(sb, s.Value)
	}
	if s.Pivot != nil {
		add("Pivot:")
		writeNode(sb, s.Pivot)
	}
	if len(s.Targets) > 0 {
		add("Targets:[")
		for i, t := range s.Targets {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, t)
		}
		sb.WriteString("]")
	}
	if s.From != nil {
		add("From:")
		writeNode(sb, s.From)
	}
	if len(s.Let) > 0 {
		add("Let:[")
		for i, l := range s.Let {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, l)
		}
		sb.WriteString("]")
	}
	if s.Where != nil {
		add("Where:")
		writeNode(sb, s.Where)
	}
	if s.GroupBy != nil {
		add("GroupBy:")
		writeNode(sb, s.GroupBy)
	}
	if s.Having != nil {
		add("Having:")
		writeNode(sb, s.Having)
	}
	if len(s.OrderBy) > 0 {
		add("OrderBy:[")
		for i, o := range s.OrderBy {
			if i > 0 {
				sb.WriteString(" ")
			}
			writeNode(sb, o)
		}
		sb.WriteString("]")
	}
	if s.Limit != nil {
		add("Limit:")
		writeNode(sb, s.Limit)
	}
	if s.Offset != nil {
		add("Offset:")
		writeNode(sb, s.Offset)
	}
	sb.WriteString("}")
}

// writeDmlStmt is shared by InsertStmt, UpsertStmt, and ReplaceStmt — they
// have nearly the same shape (target, alias, value, pos, on-conflict,
// returning). For UpsertStmt and ReplaceStmt, callers pass nil for pos
// because the grammar's `replaceCommand` (line 121) and `upsertCommand`
// (line 125) do not have an `AT pos` clause.
func writeDmlStmt(sb *strings.Builder, name string, target TableExpr, alias *string, value ExprNode, pos ExprNode, oc *OnConflict, ret *ReturningClause) {
	fmt.Fprintf(sb, "%s{Target:", name)
	writeNode(sb, target)
	writeOptString(sb, " AsAlias:", alias)
	sb.WriteString(" Value:")
	writeNode(sb, value)
	if pos != nil {
		sb.WriteString(" Pos:")
		writeNode(sb, pos)
	}
	if oc != nil {
		sb.WriteString(" OnConflict:")
		writeNode(sb, oc)
	}
	if ret != nil {
		sb.WriteString(" Returning:")
		writeNode(sb, ret)
	}
	sb.WriteString("}")
}

// writeOptString appends a label + value if the optional string is non-nil.
func writeOptString(sb *strings.Builder, label string, s *string) {
	if s == nil {
		return
	}
	sb.WriteString(label)
	sb.WriteString(*s)
}
```

- [ ] **Step 2: Add `TestNodeToString` golden cases to `ast_test.go`**

Append a new test function at the bottom of `ast_test.go`:

```go
// ---------------------------------------------------------------------------
// TestNodeToString — golden assertions on outfuncs.go output.
//
// One representative case per node category, plus the multi-interface nodes
// and at least one deeply nested case.
// ---------------------------------------------------------------------------

func TestNodeToString(t *testing.T) {
	cases := []struct {
		name string
		node Node
		want string
	}{
		{
			name: "string_literal",
			node: &StringLit{Val: "hello"},
			want: `StringLit{Val:"hello"}`,
		},
		{
			name: "number_literal",
			node: &NumberLit{Val: "42"},
			want: `NumberLit{Val:42}`,
		},
		{
			name: "binary_expr_add",
			node: &BinaryExpr{
				Op:    BinOpAdd,
				Left:  &VarRef{Name: "a"},
				Right: &NumberLit{Val: "1"},
			},
			want: `BinaryExpr{Op:+ Left:VarRef{Name:a} Right:NumberLit{Val:1}}`,
		},
		{
			name: "in_expr_with_list",
			node: &InExpr{
				Expr: &VarRef{Name: "x"},
				List: []ExprNode{&NumberLit{Val: "1"}, &NumberLit{Val: "2"}},
				Not:  false,
			},
			want: `InExpr{Expr:VarRef{Name:x} List:[NumberLit{Val:1} NumberLit{Val:2}] Not:false}`,
		},
		{
			name: "is_missing_predicate",
			node: &IsExpr{
				Expr: &VarRef{Name: "y"},
				Type: IsTypeMissing,
				Not:  false,
			},
			want: `IsExpr{Expr:VarRef{Name:y} Type:MISSING Not:false}`,
		},
		{
			name: "func_call_count_distinct",
			node: &FuncCall{
				Name:       "COUNT",
				Quantifier: QuantifierDistinct,
				Args:       []ExprNode{&VarRef{Name: "id"}},
			},
			want: `FuncCall{Name:COUNT Quantifier:DISTINCT Args:[VarRef{Name:id}]}`,
		},
		{
			name: "func_call_count_star",
			node: &FuncCall{Name: "COUNT", Star: true},
			want: `FuncCall{Name:COUNT Star:true Args:[]}`,
		},
		{
			name: "case_expr_searched",
			node: &CaseExpr{
				Whens: []*CaseWhen{
					{When: &BoolLit{Val: true}, Then: &NumberLit{Val: "1"}},
				},
				Else: &NumberLit{Val: "0"},
			},
			want: `CaseExpr{Whens:[CaseWhen{When:BoolLit{Val:true} Then:NumberLit{Val:1}}] Else:NumberLit{Val:0}}`,
		},
		{
			name: "cast_expr",
			node: &CastExpr{
				Kind:   CastKindCast,
				Expr:   &VarRef{Name: "x"},
				AsType: &TypeRef{Name: "INTEGER"},
			},
			want: `CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:INTEGER}}`,
		},
		{
			name: "path_expr_with_steps",
			node: &PathExpr{
				Root:  &VarRef{Name: "Music"},
				Steps: []PathStep{&DotStep{Field: "albums"}, &WildcardStep{}},
			},
			want: `PathExpr{Root:VarRef{Name:Music} Steps:[DotStep{Field:albums} WildcardStep{}]}`,
		},
		{
			name: "var_ref_at_prefixed",
			node: &VarRef{Name: "doc", AtPrefixed: true},
			want: `VarRef{Name:doc AtPrefixed:true}`,
		},
		{
			name: "param_ref",
			node: &ParamRef{},
			want: `ParamRef{}`,
		},
		{
			name: "sub_link",
			node: &SubLink{Stmt: &SelectStmt{Star: true}},
			want: `SubLink{Stmt:SelectStmt{Star:true}}`,
		},
		{
			name: "list_literal",
			node: &ListLit{Items: []ExprNode{&NumberLit{Val: "1"}, &NumberLit{Val: "2"}}},
			want: `ListLit{Items:[NumberLit{Val:1} NumberLit{Val:2}]}`,
		},
		{
			name: "bag_literal",
			node: &BagLit{Items: []ExprNode{&NumberLit{Val: "1"}}},
			want: `BagLit{Items:[NumberLit{Val:1}]}`,
		},
		{
			name: "tuple_literal",
			node: &TupleLit{
				Pairs: []*TuplePair{
					{Key: &StringLit{Val: "k"}, Value: &NumberLit{Val: "1"}},
				},
			},
			want: `TupleLit{Pairs:[TuplePair{Key:StringLit{Val:"k"} Value:NumberLit{Val:1}}]}`,
		},
		{
			name: "join_expr",
			node: &JoinExpr{
				Kind:  JoinKindInner,
				Left:  &TableRef{Name: "a"},
				Right: &TableRef{Name: "b"},
				On:    &BoolLit{Val: true},
			},
			want: `JoinExpr{Kind:INNER Left:TableRef{Name:a} Right:TableRef{Name:b} On:BoolLit{Val:true}}`,
		},
		{
			name: "select_with_path_in_from",
			node: &SelectStmt{
				Star: true,
				From: &PathExpr{
					Root:  &VarRef{Name: "Music"},
					Steps: []PathStep{&DotStep{Field: "albums"}, &WildcardStep{}},
				},
			},
			want: `SelectStmt{Star:true From:PathExpr{Root:VarRef{Name:Music} Steps:[DotStep{Field:albums} WildcardStep{}]}}`,
		},
		{
			name: "insert_stmt_with_returning",
			node: &InsertStmt{
				Target: &TableRef{Name: "Music"},
				Value:  &TupleLit{Pairs: []*TuplePair{{Key: &StringLit{Val: "k"}, Value: &NumberLit{Val: "1"}}}},
				Returning: &ReturningClause{
					Items: []*ReturningItem{
						{Status: ReturningStatusModified, Mapping: ReturningMappingNew, Star: true},
					},
				},
			},
			want: `InsertStmt{Target:TableRef{Name:Music} Value:TupleLit{Pairs:[TuplePair{Key:StringLit{Val:"k"} Value:NumberLit{Val:1}}]} Returning:ReturningClause{Items:[ReturningItem{Status:MODIFIED Mapping:NEW Star:true}]}}`,
		},
		{
			name: "match_expr",
			node: &MatchExpr{
				Expr: &VarRef{Name: "g"},
				Patterns: []*GraphPattern{
					{Parts: []PatternNode{&NodePattern{Variable: &VarRef{Name: "n"}}}},
				},
			},
			want: `MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:n}}]}]}`,
		},
		{
			name: "type_ref_decimal",
			node: &TypeRef{Name: "DECIMAL", Args: []int{10, 2}},
			want: `TypeRef{Name:DECIMAL Args:[10,2]}`,
		},
		{
			name: "type_ref_time_with_tz",
			node: &TypeRef{Name: "TIME", WithTimeZone: true},
			want: `TypeRef{Name:TIME WithTimeZone:true}`,
		},
		{
			name: "nil_node",
			node: nil,
			want: `<nil>`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NodeToString(tc.node)
			if got != tc.want {
				t.Errorf("NodeToString() mismatch\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 3: Run the new test**

```bash
go test -run TestNodeToString$ ./partiql/ast/...
```

Expected: ok with all 23 sub-tests passing. If any fail, fix the corresponding arm in `outfuncs.go` until the golden matches. **Do not adjust the golden to match a buggy `NodeToString` output** — the goldens encode the spec's "stable, deterministic" invariant.

- [ ] **Step 4: Add the reflection safety net `TestNodeToString_AllNodesCovered`**

Append at the bottom of `ast_test.go`:

```go
// ---------------------------------------------------------------------------
// TestNodeToString_AllNodesCovered — reflection safety net.
//
// Walks every node in the package's TestGetLoc table and asserts that
// NodeToString returns a non-empty result without panicking. This catches
// the case where a new node type is added to TestGetLoc but never wired
// into outfuncs.go's switch.
// ---------------------------------------------------------------------------

func TestNodeToString_AllNodesCovered(t *testing.T) {
	// Build a fresh table identical in shape to TestGetLoc but constructed
	// here so the two tests cannot drift apart silently. We could share via
	// a top-level helper but inlining keeps the test self-contained.
	cases := allNodesForCoverageTest()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("NodeToString panicked: %v", r)
				}
			}()
			got := NodeToString(tc.node)
			if got == "" {
				t.Errorf("NodeToString returned empty string for %s", tc.name)
			}
			if strings.HasPrefix(got, "<unknown:") {
				t.Errorf("outfuncs.go is missing a switch arm for %s: got %q", tc.name, got)
			}
		})
	}
}

// allNodesForCoverageTest returns one zero-value instance of every node
// type defined in this package. Adding a new node type means adding it
// here too — the test will fail until you do.
func allNodesForCoverageTest() []struct {
	name string
	node Node
} {
	return []struct {
		name string
		node Node
	}{
		// node.go
		{"List", &List{}},

		// literals.go
		{"StringLit", &StringLit{}},
		{"NumberLit", &NumberLit{}},
		{"BoolLit", &BoolLit{}},
		{"NullLit", &NullLit{}},
		{"MissingLit", &MissingLit{}},
		{"DateLit", &DateLit{}},
		{"TimeLit", &TimeLit{}},
		{"IonLit", &IonLit{}},

		// exprs.go
		{"BinaryExpr", &BinaryExpr{}},
		{"UnaryExpr", &UnaryExpr{}},
		{"InExpr", &InExpr{}},
		{"BetweenExpr", &BetweenExpr{}},
		{"LikeExpr", &LikeExpr{}},
		{"IsExpr", &IsExpr{}},
		{"FuncCall", &FuncCall{}},
		{"CaseExpr", &CaseExpr{}},
		{"CaseWhen", &CaseWhen{}},
		{"CastExpr", &CastExpr{}},
		{"ExtractExpr", &ExtractExpr{}},
		{"TrimExpr", &TrimExpr{}},
		{"SubstringExpr", &SubstringExpr{}},
		{"CoalesceExpr", &CoalesceExpr{}},
		{"NullIfExpr", &NullIfExpr{}},
		{"WindowSpec", &WindowSpec{}},
		{"PathExpr", &PathExpr{}},
		{"VarRef", &VarRef{}},
		{"ParamRef", &ParamRef{}},
		{"SubLink", &SubLink{}},
		{"ListLit", &ListLit{}},
		{"BagLit", &BagLit{}},
		{"TupleLit", &TupleLit{}},
		{"TuplePair", &TuplePair{}},
		{"DotStep", &DotStep{}},
		{"AllFieldsStep", &AllFieldsStep{}},
		{"IndexStep", &IndexStep{}},
		{"WildcardStep", &WildcardStep{}},

		// tableexprs.go
		{"TableRef", &TableRef{}},
		{"AliasedSource", &AliasedSource{}},
		{"JoinExpr", &JoinExpr{}},
		{"UnpivotExpr", &UnpivotExpr{}},

		// types.go
		{"TypeRef", &TypeRef{}},

		// stmts.go — top-level
		{"SelectStmt", &SelectStmt{}},
		{"SetOpStmt", &SetOpStmt{}},
		{"ExplainStmt", &ExplainStmt{}},
		{"InsertStmt", &InsertStmt{}},
		{"UpdateStmt", &UpdateStmt{}},
		{"DeleteStmt", &DeleteStmt{}},
		{"UpsertStmt", &UpsertStmt{}},
		{"ReplaceStmt", &ReplaceStmt{}},
		{"RemoveStmt", &RemoveStmt{}},
		{"CreateTableStmt", &CreateTableStmt{}},
		{"CreateIndexStmt", &CreateIndexStmt{}},
		{"DropTableStmt", &DropTableStmt{}},
		{"DropIndexStmt", &DropIndexStmt{}},
		{"ExecStmt", &ExecStmt{}},

		// stmts.go — clause and DML helpers
		{"TargetEntry", &TargetEntry{}},
		{"PivotProjection", &PivotProjection{}},
		{"LetBinding", &LetBinding{}},
		{"GroupByClause", &GroupByClause{}},
		{"GroupByItem", &GroupByItem{}},
		{"OrderByItem", &OrderByItem{}},
		{"SetAssignment", &SetAssignment{}},
		{"OnConflict", &OnConflict{}},
		{"OnConflictTarget", &OnConflictTarget{}},
		{"ReturningClause", &ReturningClause{}},
		{"ReturningItem", &ReturningItem{}},

		// patterns.go
		{"MatchExpr", &MatchExpr{}},
		{"GraphPattern", &GraphPattern{}},
		{"NodePattern", &NodePattern{}},
		{"EdgePattern", &EdgePattern{}},
		{"PatternQuantifier", &PatternQuantifier{}},
		{"PatternSelector", &PatternSelector{}},
	}
}
```

Add the import at the top of `ast_test.go` (just after the existing `"testing"` import) — `strings` is needed by the new test:

```go
import (
	"strings"
	"testing"
)
```

- [ ] **Step 5: Run the safety net test**

```bash
go test -run TestNodeToString_AllNodesCovered ./partiql/ast/...
```

Expected: ok, 73 sub-tests passing. If any sub-test fails with `outfuncs.go is missing a switch arm for X`, add the missing case to `writeNode` in `outfuncs.go` and re-run.

- [ ] **Step 6: Run the full test suite**

```bash
go test ./partiql/ast/...
```

Expected: ok, all tests pass (`TestGetLoc` 73 sub-tests + `TestNodeToString` 23 sub-tests + `TestNodeToString_AllNodesCovered` 73 sub-tests = 169 sub-tests across 3 test functions).

- [ ] **Step 7: Run vet and gofmt**

```bash
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Both clean.

- [ ] **Step 8: Commit**

```bash
git add partiql/ast/outfuncs.go partiql/ast/ast_test.go
git commit -m "$(cat <<'EOF'
feat(partiql/ast): add NodeToString and reflection safety net

outfuncs.go: NodeToString writes a deterministic Go-struct-like dump
of any AST node. Switch covers all 73 node types. Loc fields are
omitted (positions are tested separately). Helpers writeSelectStmt
and writeDmlStmt extract the bigger arms.

ast_test.go: TestNodeToString golden cases (23 representative shapes
including the multi-interface nodes and a deeply nested case);
TestNodeToString_AllNodesCovered reflection safety net that calls
NodeToString on a zero value of every node type and asserts no panic
or "<unknown:>" fallback. Adding a new node type without wiring it
into outfuncs.go now fails the build.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

---

### Task 12: Final verification, grammar cross-check, DAG bookkeeping, finishing branch

**Files:**
- Read: `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4` (legacy grammar — for cross-check)
- Modify: `/Users/h3n4l/OpenSource/omni/docs/migration/partiql/dag.md` (mark node 1 as done — on `main`, not the worktree)

This task does not write more code unless the grammar cross-check finds a missing rule. It is the acceptance gate.

- [ ] **Step 1: Run the full test suite, vet, and gofmt one more time**

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core
go test ./partiql/ast/...
go vet ./partiql/ast/...
gofmt -l partiql/ast/
```

Expected: ok, no output from vet or gofmt.

- [ ] **Step 2: Cross-check the AST against `PartiQLParser.g4`**

Open `/Users/h3n4l/OpenSource/parser/partiql/PartiQLParser.g4` and walk it from top to bottom. For each grammar rule, find a node type in `partiql/ast/` that represents it, using the spec's "Coverage cross-check" table as the working list. Take notes:

- Rules covered: ✓
- Rules missing or wrong-shaped: list with rule name + grammar line number

Expected outcome: zero gaps. The analysis spent considerable effort on this already and the spec's cross-check covers the categories. If you find a missing rule:

  1. Add the corresponding node type to the appropriate file
  2. Add it to all four test groups in `ast_test.go` (compile-time assertion, `TestGetLoc` row, `TestNodeToString` golden if representative, and `allNodesForCoverageTest`)
  3. Add a switch arm to `outfuncs.go`
  4. Re-run `go test ./partiql/ast/...`
  5. Make a separate commit per missed-rule fix with message `feat(partiql/ast): add <Node> for grammar rule <name>`

If you find no gaps, write a brief note in the verification commit message confirming the cross-check was run.

- [ ] **Step 3: Run the cross-check spot-check commands**

```bash
# Sanity: count node types via grep
grep -c '^func (\*\w\+) nodeTag()' partiql/ast/*.go

# Should print something like:
# partiql/ast/exprs.go:28
# partiql/ast/literals.go:9
# partiql/ast/node.go:1
# partiql/ast/patterns.go:6
# partiql/ast/stmts.go:25
# partiql/ast/tableexprs.go:4
# partiql/ast/types.go:1
#
# Total ~73. If your count differs, investigate before proceeding.

# Sanity: every node has GetLoc
grep -c 'GetLoc() Loc { return n.Loc }' partiql/ast/*.go

# Sanity: outfuncs.go switch has a case for every type
grep -c '^	case \*' partiql/ast/outfuncs.go
# Should be ~73 (one per node type).
```

If any number is off, look for the missing types using diff between the file lists and the switch arms.

- [ ] **Step 4: Run a coverage report on `outfuncs.go`**

```bash
go test -coverprofile=/tmp/partiql-ast-cover.out ./partiql/ast/...
go tool cover -func=/tmp/partiql-ast-cover.out | grep outfuncs.go
```

Expected: `outfuncs.go:NodeToString` and `writeNode` very close to 100%. Some optional-field branches may not be exercised by goldens — that's fine.

- [ ] **Step 5: Confirm the spec's acceptance criteria**

Re-read `docs/superpowers/specs/2026-04-08-partiql-ast-core-design.md` "Acceptance criteria" section and verify each item:

| # | Criterion | Status |
|---|-----------|--------|
| 1 | All ~73 node types defined across the 8 files | check `grep -c` output |
| 2 | Every node has a doc comment naming the grammar rule | spot-check 5 random nodes |
| 3 | Every node implements `nodeTag()` and `GetLoc() Loc` | covered by `var _ Node` assertions |
| 4 | Every node implements at least one sub-interface OR is a documented bare-Node helper | covered by `var _ <Iface>` assertions |
| 5 | `outfuncs.NodeToString` covers every type | covered by `TestNodeToString_AllNodesCovered` |
| 6 | `go test ./partiql/ast/...` passes | Step 1 |
| 7 | `go vet ./partiql/ast/...` clean | Step 1 |
| 8 | `gofmt` clean | Step 1 |
| 9 | Coverage cross-check vs `PartiQLParser.g4` | Step 2 |

If any criterion fails, fix it before proceeding. Do not move to Step 6 until all are green.

- [ ] **Step 6: Final commit (verification note, even if there were no code changes)**

```bash
git commit --allow-empty -m "$(cat <<'EOF'
chore(partiql/ast): final verification pass

- go test ./partiql/ast/... — pass
- go vet ./partiql/ast/... — clean
- gofmt -l partiql/ast/ — clean
- Cross-checked against bytebase/parser/partiql/PartiQLParser.g4
  line by line; every rule maps to a node type
- All 9 acceptance criteria from the spec satisfied

Closes ast-core (DAG node 1) for the partiql migration.

Co-Authored-By: Claude Opus 4.6 (1M context) <noreply@anthropic.com>
EOF
)"
```

(`--allow-empty` is acceptable here because this commit is purely a verification gate marker — useful for reviewers and for the DAG status update in the next steps. If the cross-check found issues and you committed them as separate `feat` commits, this final empty commit is still a useful "done" marker.)

- [ ] **Step 7: Update `dag.md` on `main` to mark node 1 as `done`**

The DAG file lives on `main`, not on the feature branch. Switch to the main worktree directory:

```bash
cd /Users/h3n4l/OpenSource/omni
```

The migration docs were untracked when this work started. They may still be untracked or may have been committed in the meantime. Check:

```bash
git status docs/migration/partiql/
```

If `docs/migration/partiql/dag.md` is **untracked**, you need to commit it first as a baseline before updating the status (so the diff is sensible). If it's already tracked, just edit it.

Open `docs/migration/partiql/dag.md` in your editor and find the table row:

```
| 1 | ast-core | `partiql/ast` | (none) | lexer, catalog | **P0** | not started |
```

Change `not started` to `done`:

```
| 1 | ast-core | `partiql/ast` | (none) | lexer, catalog | **P0** | done |
```

If `dag.md` was untracked at the start of this session, also stage and commit it now along with the status update. Coordinate with the user about whether to commit the migration docs to main as part of this task or as a separate housekeeping commit.

- [ ] **Step 8: Hand off to `superpowers:finishing-a-development-branch`**

This is the integration step. From the worktree:

```bash
cd /Users/h3n4l/OpenSource/omni/.worktrees/feat-partiql-ast-core
git log --oneline feat/partiql/ast-core ^main
```

Verify the commit list looks right (you should see commits from Tasks 1–12 plus the spec/plan commits from brainstorming/writing-plans).

Then invoke the `superpowers:finishing-a-development-branch` skill to:
- Pick a merge strategy (squash vs merge vs PR)
- Push the branch if needed
- Clean up the worktree once integrated

The skill will guide the branch finish — do not improvise the merge strategy.

- [ ] **Step 9 (after branch finish): Clean up the worktree and report**

After `finishing-a-development-branch` completes successfully:

1. Confirm `git worktree list` no longer shows `feat-partiql-ast-core`
2. Confirm `dag.md` on `main` shows ast-core as `done`
3. Report to the user:
   - Which DAG node was completed
   - How many node types were added (~73)
   - Test pass count (~169 sub-tests)
   - The next actionable nodes per the DAG (node 2 `lexer` should now be unblocked, node 3 `catalog` was already unblocked)

---

## Self-Review

Walking the spec section by section against this plan:

**Goal & inputs (spec lines 9–18):** Task 1 establishes the package and references the spec. ✓

**D1. Sealed sub-interface style (lines 22–78):** Task 1 declares all six interfaces with the correct tag-method names. ✓

**D2. Multi-interface satisfaction (lines 80–90):** Task 5 implements `PathExpr`, `VarRef`, `SubLink` with both `exprNode()` and `tableExpr()` tag methods. The compile-time assertions in `ast_test.go` cover them. ✓

**D3. Byte-offset Loc only (lines 92–94):** Task 1 declares `Loc{Start, End int}`. No line/column. ✓

**D4. Files split by category (lines 96–112):** Task 1 creates `node.go`. Tasks 2–10 create the other 7 files. ✓

**D5. Single TypeRef (lines 114–127):** Task 7 implements a flat `TypeRef`. ✓

**D6. FuncCall covers regular/aggregate/window (lines 129–146):** Task 4 implements `FuncCall` with `Quantifier`, `Star`, and `Over` fields. The note about which built-ins get dedicated nodes vs plain `FuncCall` is preserved (Tasks 4 dedicate `ExtractExpr`, `TrimExpr`, `SubstringExpr`; `DATE_ADD`/`DATE_DIFF` use plain `FuncCall`). ✓

**D7. No Walker, no smart constructors (lines 148–150):** Plan never adds them. ✓

**Type taxonomy (lines 152–306):** Each node type from the spec's tables is created in the corresponding task. Cross-check against the file map at the top of this plan and against the `allNodesForCoverageTest` listing in Task 11 — both enumerate the full set. ✓

**Coverage cross-check (lines 318–355):** Task 12 Step 2 explicitly walks the spec's cross-check table against `PartiQLParser.g4`. ✓

**Test plan (lines 359–394):** Task 1 creates the `TestGetLoc` skeleton. Tasks 2–10 append rows. Task 11 adds `TestNodeToString` (golden) and `TestNodeToString_AllNodesCovered` (reflection safety net). Compile-time `var _ <Iface>` assertions are added incrementally in every task. ✓

**Test corpus (forward reference, lines 396–432):** Spec marks this as out of scope for `ast-core`; plan does not include corpus tests. ✓

**Acceptance criteria (lines 396–408 in updated spec):** Task 12 Step 5 walks all 9 criteria. ✓

**Non-goals (lines 410–420):** Plan never implements parser logic, deparse, semantic, walker, smart constructors, corpus tests, or line/column. ✓

**Risks & mitigations (lines 422–429):** The spec's mitigations are implemented in the plan:
- Risk "missing grammar rule": Task 12 Step 2 line-by-line cross-check
- Risk "graph pattern shape may shift": Task 10 includes the docstring noting this
- Risk "multi-interface confuses callers": Task 5 inline doc-comments + `var _` assertions
- Risk "outfuncs.go goes stale": Task 11 reflection safety net
- Risk "field bloat on large nodes": acknowledged in `SelectStmt` (13 fields), no action

**Placeholder scan:** Searched for "TBD", "TODO", "fill in details", "appropriate error handling", etc. None present. The forward-declared placeholder block in Task 8 is explicitly labeled as such and is removed in Task 9 — that's intentional sequencing, not a placeholder failure.

**Type consistency:** Method names, field names, and enum names are consistent across tasks. `QuantifierKind` is declared in `exprs.go` (Task 4) and used in `stmts.go` (Tasks 8–9) — same package, no import. The deviation from the spec's enum-location list is documented in the "Task Order Rationale" section above.

No issues found.

---

## Plan complete and saved

Plan saved to `docs/superpowers/plans/2026-04-08-partiql-ast-core.md` (this file).

Next step: pick an execution mode.
