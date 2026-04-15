package ast

// This file holds expression AST node types for Doris SQL (T1.3).

// ---------------------------------------------------------------------------
// Operator types
// ---------------------------------------------------------------------------

// BinaryOp identifies a binary operator kind.
type BinaryOp int

const (
	// Logical operators.
	BinOr  BinaryOp = iota // OR, ||
	BinAnd                  // AND, &&
	BinXor                  // XOR

	// Comparison operators.
	BinEq        // =
	BinNe        // <> or !=
	BinLt        // <
	BinGt        // >
	BinLe        // <=
	BinGe        // >=
	BinNullSafeEq // <=>

	// Arithmetic operators.
	BinAdd // +
	BinSub // -
	BinMul // *
	BinDiv // /
	BinMod // % or MOD
	BinIntDiv // DIV

	// Bitwise operators.
	BinBitOr  // |
	BinBitAnd // &
	BinBitXor // ^
	BinShiftLeft  // <<
	BinShiftRight // >>
)

// String returns the SQL text form of the operator.
func (op BinaryOp) String() string {
	switch op {
	case BinOr:
		return "OR"
	case BinAnd:
		return "AND"
	case BinXor:
		return "XOR"
	case BinEq:
		return "="
	case BinNe:
		return "<>"
	case BinLt:
		return "<"
	case BinGt:
		return ">"
	case BinLe:
		return "<="
	case BinGe:
		return ">="
	case BinNullSafeEq:
		return "<=>"
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
	case BinIntDiv:
		return "DIV"
	case BinBitOr:
		return "|"
	case BinBitAnd:
		return "&"
	case BinBitXor:
		return "^"
	case BinShiftLeft:
		return "<<"
	case BinShiftRight:
		return ">>"
	default:
		return "?"
	}
}

// UnaryOp identifies a unary operator kind.
type UnaryOp int

const (
	UnaryMinus  UnaryOp = iota // -
	UnaryPlus                   // +
	UnaryBitNot                 // ~
	UnaryNot                    // NOT
)

// String returns the SQL text form of the unary operator.
func (op UnaryOp) String() string {
	switch op {
	case UnaryMinus:
		return "-"
	case UnaryPlus:
		return "+"
	case UnaryBitNot:
		return "~"
	case UnaryNot:
		return "NOT"
	default:
		return "?"
	}
}

// LitKind classifies a literal value.
type LitKind int

const (
	LitInt     LitKind = iota // integer literal
	LitFloat                  // decimal/float literal
	LitString                 // string literal
	LitBool                   // TRUE or FALSE
	LitNull                   // NULL
	LitKeyword                // keyword used as literal value, e.g. DEFAULT
)

// CaseKind classifies a CASE expression as simple or searched.
type CaseKind int

const (
	CaseSearched CaseKind = iota // CASE WHEN ...
	CaseSimple                    // CASE operand WHEN ...
)

// ---------------------------------------------------------------------------
// Expression nodes
// ---------------------------------------------------------------------------

// BinaryExpr represents a binary operation: left op right.
type BinaryExpr struct {
	Op    BinaryOp
	Left  Node
	Right Node
	Loc   Loc
}

func (n *BinaryExpr) Tag() NodeTag { return T_BinaryExpr }

var _ Node = (*BinaryExpr)(nil)

// UnaryExpr represents a unary operation: op expr.
type UnaryExpr struct {
	Op   UnaryOp
	Expr Node
	Loc  Loc
}

func (n *UnaryExpr) Tag() NodeTag { return T_UnaryExpr }

var _ Node = (*UnaryExpr)(nil)

// IsExpr represents expr IS [NOT] NULL / TRUE / FALSE.
type IsExpr struct {
	Expr    Node
	Not     bool
	IsWhat  string // "NULL", "TRUE", "FALSE"
	Loc     Loc
}

func (n *IsExpr) Tag() NodeTag { return T_IsExpr }

var _ Node = (*IsExpr)(nil)

// BetweenExpr represents expr [NOT] BETWEEN low AND high.
type BetweenExpr struct {
	Expr Node
	Low  Node
	High Node
	Not  bool
	Loc  Loc
}

func (n *BetweenExpr) Tag() NodeTag { return T_BetweenExpr }

var _ Node = (*BetweenExpr)(nil)

// InExpr represents expr [NOT] IN (values...) or expr [NOT] IN (subquery).
type InExpr struct {
	Expr     Node
	Values   []Node // list of values; for subquery, contains one SubqueryExpr
	Not      bool
	Loc      Loc
}

func (n *InExpr) Tag() NodeTag { return T_InExpr }

var _ Node = (*InExpr)(nil)

// LikeExpr represents expr [NOT] LIKE pattern [ESCAPE esc].
type LikeExpr struct {
	Expr    Node
	Pattern Node
	Escape  Node // nil if no ESCAPE clause
	Not     bool
	Loc     Loc
}

func (n *LikeExpr) Tag() NodeTag { return T_LikeExpr }

var _ Node = (*LikeExpr)(nil)

// RegexpExpr represents expr [NOT] REGEXP|RLIKE pattern.
type RegexpExpr struct {
	Expr    Node
	Pattern Node
	Not     bool
	Loc     Loc
}

func (n *RegexpExpr) Tag() NodeTag { return T_RegexpExpr }

var _ Node = (*RegexpExpr)(nil)

// FuncCallExpr represents a function call: name(args...).
type FuncCallExpr struct {
	Name     *ObjectName
	Args     []Node
	Distinct bool   // COUNT(DISTINCT x)
	Star     bool   // COUNT(*)
	OrderBy  []*OrderByItem // optional ORDER BY in aggregate (GROUP_CONCAT)
	Separator string // optional SEPARATOR value for GROUP_CONCAT
	Loc      Loc
}

func (n *FuncCallExpr) Tag() NodeTag { return T_FuncCallExpr }

var _ Node = (*FuncCallExpr)(nil)

// CastExpr represents CAST(expr AS type) or TRY_CAST(expr AS type).
type CastExpr struct {
	Expr     Node
	TypeName *TypeName
	TryCast  bool // true for TRY_CAST
	Loc      Loc
}

func (n *CastExpr) Tag() NodeTag { return T_CastExpr }

var _ Node = (*CastExpr)(nil)

// CaseExpr represents CASE [operand] WHEN...THEN...ELSE...END.
type CaseExpr struct {
	Kind    CaseKind
	Operand Node          // nil for searched CASE
	Whens   []*WhenClause
	Else    Node          // nil if no ELSE
	Loc     Loc
}

func (n *CaseExpr) Tag() NodeTag { return T_CaseExpr }

var _ Node = (*CaseExpr)(nil)

// WhenClause represents one WHEN cond THEN result arm of a CASE expression.
type WhenClause struct {
	Cond   Node
	Result Node
	Loc    Loc
}

func (n *WhenClause) Tag() NodeTag { return T_WhenClause }

var _ Node = (*WhenClause)(nil)

// SubqueryExpr is a placeholder for a subquery appearing in expression
// position. Real subquery parsing comes in T1.4; for now it stores the
// raw text between parentheses.
type SubqueryExpr struct {
	RawText string // raw SQL text of the subquery (without outer parens)
	Loc     Loc
}

func (n *SubqueryExpr) Tag() NodeTag { return T_SubqueryExpr }

var _ Node = (*SubqueryExpr)(nil)

// ColumnRef represents a qualified column reference used as an expression.
// It reuses ObjectName for the multipart identifier.
type ColumnRef struct {
	Name *ObjectName
	Loc  Loc
}

func (n *ColumnRef) Tag() NodeTag { return T_ColumnRef }

var _ Node = (*ColumnRef)(nil)

// Literal represents a literal value: numeric, string, boolean, or NULL.
type Literal struct {
	Kind  LitKind
	Value string // original source text
	Loc   Loc
}

func (n *Literal) Tag() NodeTag { return T_Literal }

var _ Node = (*Literal)(nil)

// ParenExpr represents a parenthesized expression: (expr).
type ParenExpr struct {
	Expr Node
	Loc  Loc
}

func (n *ParenExpr) Tag() NodeTag { return T_ParenExpr }

var _ Node = (*ParenExpr)(nil)

// ExistsExpr represents EXISTS (subquery).
type ExistsExpr struct {
	Subquery *SubqueryExpr
	Loc      Loc
}

func (n *ExistsExpr) Tag() NodeTag { return T_ExistsExpr }

var _ Node = (*ExistsExpr)(nil)

// IntervalExpr represents INTERVAL expr unit.
type IntervalExpr struct {
	Value Node
	Unit  string // e.g., "DAY", "HOUR", "MINUTE", "SECOND", "MONTH", "YEAR"
	Loc   Loc
}

func (n *IntervalExpr) Tag() NodeTag { return T_IntervalExpr }

var _ Node = (*IntervalExpr)(nil)

// OrderByItem represents an expression with optional ASC/DESC and NULLS FIRST/LAST.
type OrderByItem struct {
	Expr       Node
	Desc       bool   // true for DESC, false for ASC (default)
	NullsFirst *bool  // nil if not specified; true for NULLS FIRST, false for NULLS LAST
	Loc        Loc
}

func (n *OrderByItem) Tag() NodeTag { return T_OrderByItem }

var _ Node = (*OrderByItem)(nil)
