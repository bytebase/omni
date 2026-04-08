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
//
//	mathOp00 (CONCAT), mathOp01 (PLUS/MINUS), mathOp02 (PERCENT/ASTERISK/SLASH_FORWARD)
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
//
//	aggregate#CountAll, aggregate#AggregateBase,
//	windowFunction#LagLeadFunction
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
	Operand ExprNode // nil for searched CASE
	Whens   []*CaseWhen
	Else    ExprNode // nil if no ELSE clause
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
