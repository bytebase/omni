// Package ast defines T-SQL parse tree node types.
// These types represent the abstract syntax tree produced by the T-SQL parser.
package ast

// Node is the interface implemented by all parse tree nodes.
type Node interface {
	nodeTag()
}

// ExprNode is the interface for expression nodes.
type ExprNode interface {
	Node
	exprNode()
}

// TableExpr is the interface for table reference nodes in FROM clauses.
type TableExpr interface {
	Node
	tableExpr()
}

// StmtNode is the interface for statement nodes.
type StmtNode interface {
	Node
	stmtNode()
}

// Loc represents a source location range (byte offsets).
// Zero values indicate that the position was not recorded.
type Loc struct {
	Start int // inclusive start byte offset (0 if not set)
	End   int // exclusive end byte offset (0 if not set)
}

// List represents a generic list of nodes.
type List struct {
	Items []Node
}

func (l *List) nodeTag() {}

// Len returns the number of items in the list.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}

// String represents a string value node.
type String struct {
	Str string
}

func (s *String) nodeTag() {}

// Integer represents an integer value node.
type Integer struct {
	Ival int64
}

func (i *Integer) nodeTag() {}

// Float represents a floating-point value node.
// Stored as string to preserve precision.
type Float struct {
	Fval string
}

func (f *Float) nodeTag() {}

// Boolean represents a boolean value node.
type Boolean struct {
	Boolval bool
}

func (b *Boolean) nodeTag() {}
