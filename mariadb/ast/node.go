// Package ast defines MySQL parse tree node types.
package ast

// Node is the interface implemented by all MySQL parse tree nodes.
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
// Zero values mean "not yet set" (parsers must explicitly set Start and End).
type Loc struct {
	Start int // inclusive start byte offset, 0 if not yet set
	End   int // exclusive end byte offset, 0 if not yet set
}

// LocStart returns the start byte offset.
// This method enables field promotion: any struct with a Loc field
// automatically gets LocStart() without explicit implementation.
func (l Loc) LocStart() int { return l.Start }

// LocEnd returns the end byte offset.
func (l Loc) LocEnd() int { return l.End }

// List is a generic list of nodes.
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
