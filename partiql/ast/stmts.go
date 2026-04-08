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
