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
//
//	tableUnpivot, fromClauseSimple#FromClauseSimpleExplicit
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

// ExprSource is a FROM source written as an expression that is not a table
// expression itself: a function call, a literal, a collection, an operator.
// The grammar takes any exprSelect as a FROM source; the paths, variables,
// subqueries and graph matches among them are TableExprs and appear
// unwrapped. `FROM f(x) AS y` is AliasedSource{Source: ExprSource{Expr:
// f(x)}, As: y}. Loc is the expression's.
//
// Grammar: tableBaseReference#TableBaseRefSymbol, tableBaseReference#TableBaseRefClauses
// (source=exprSelect, PartiQLParser.g4 lines 403-404)
type ExprSource struct {
	Expr ExprNode
	Loc  Loc
}

func (*ExprSource) nodeTag()      {}
func (n *ExprSource) GetLoc() Loc { return n.Loc }
func (*ExprSource) tableExpr()    {}

// JoinExpr represents a JOIN clause: left JOIN right ON condition.
// Kind selects the JOIN flavor; On is nil for CROSS JOIN.
//
// Grammar: tableReference#TableCrossJoin, tableReference#TableQualifiedJoin
//
//	(join-type modifier from joinType; joinRhs for right-hand side)
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
