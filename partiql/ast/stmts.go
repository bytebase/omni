package ast

// PLACEHOLDER FILE — will be REPLACED WHOLESALE in Task 8/9. Do NOT
// append new types here; instead, delete this file in Task 8 and author
// the real stmts.go from scratch. This file exists only because
// WindowSpec (in exprs.go, Task 4) has a forward reference to
// OrderByItem, and the package must build between Task 4 and Task 8.
type OrderByItem struct {
	Expr          ExprNode
	Desc          bool
	NullsFirst    bool
	NullsExplicit bool
	Loc           Loc
}

func (*OrderByItem) nodeTag()      {}
func (n *OrderByItem) GetLoc() Loc { return n.Loc }
