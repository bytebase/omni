// Package ast defines parse tree node types for Azure Cosmos DB NoSQL SQL API.
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
type Loc struct {
	Start int // inclusive start byte offset, -1 = unknown
	End   int // exclusive end byte offset, -1 = unknown
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

// ---------------------------------------------------------------------------
// Literal nodes — all implement ExprNode
// ---------------------------------------------------------------------------

// StringLit represents a string constant ('hello' or "hello").
type StringLit struct {
	Val string
	Loc Loc
}

func (*StringLit) nodeTag()  {}
func (*StringLit) exprNode() {}

// NumberLit represents a numeric constant (integer, float, hex, scientific).
// Val stores the raw text to preserve the original representation.
type NumberLit struct {
	Val string
	Loc Loc
}

func (*NumberLit) nodeTag()  {}
func (*NumberLit) exprNode() {}

// BoolLit represents true or false.
type BoolLit struct {
	Val bool
	Loc Loc
}

func (*BoolLit) nodeTag()  {}
func (*BoolLit) exprNode() {}

// NullLit represents null.
type NullLit struct {
	Loc Loc
}

func (*NullLit) nodeTag()  {}
func (*NullLit) exprNode() {}

// UndefinedLit represents the CosmosDB-specific undefined constant.
type UndefinedLit struct {
	Loc Loc
}

func (*UndefinedLit) nodeTag()  {}
func (*UndefinedLit) exprNode() {}

// InfinityLit represents the Infinity constant (case-sensitive).
type InfinityLit struct {
	Loc Loc
}

func (*InfinityLit) nodeTag()  {}
func (*InfinityLit) exprNode() {}

// NanLit represents the NaN constant (case-sensitive).
type NanLit struct {
	Loc Loc
}

func (*NanLit) nodeTag()  {}
func (*NanLit) exprNode() {}
