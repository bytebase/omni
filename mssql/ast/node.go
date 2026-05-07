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
// -1 means "unknown" for either field.
type Loc struct {
	Start int // inclusive start byte offset (-1 if unknown)
	End   int // exclusive end byte offset (-1 if unknown)
}

// NoLoc returns a Loc with both Start and End set to -1 (unknown).
func NoLoc() Loc {
	return Loc{Start: -1, End: -1}
}

// IsValid reports whether both Start and End are non-negative.
func (l Loc) IsValid() bool {
	return l.Start >= 0 && l.End >= 0
}

// Merge returns the smallest Loc that contains both l and other.
// If either side is invalid, returns the other side. If both are invalid,
// returns NoLoc().
func (l Loc) Merge(other Loc) Loc {
	if !l.IsValid() {
		return other
	}
	if !other.IsValid() {
		return l
	}
	out := l
	if other.Start < out.Start {
		out.Start = other.Start
	}
	if other.End > out.End {
		out.End = other.End
	}
	return out
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
