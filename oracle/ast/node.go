// Package ast defines Oracle PL/SQL parse tree node types.
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
// Start and End are byte offsets; Start is inclusive and End is exclusive.
// A zero value is a real location at the beginning of the input. Use -1 for
// unknown positions.
type Loc struct {
	Start int // inclusive start byte offset; 0 is a valid position, -1 means unknown
	End   int // exclusive end byte offset; 0 is a valid position, -1 means unknown
}

// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
	return Loc{Start: -1, End: -1}
}

// IsUnknown reports whether both ends of the location use the unknown sentinel.
func (l Loc) IsUnknown() bool {
	return l.Start == -1 && l.End == -1
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
