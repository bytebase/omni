package ast

// walkChildren walks the child nodes of node, calling Walk(v, child)
// for each child. Maintained manually until the node count warrants a
// code generator.
func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *File:
		walkNodes(v, n.Stmts)
	case *ObjectName:
		// leaf node, no children
	case *TypeName:
		if n.ElementType != nil {
			Walk(v, n.ElementType)
		}
		if n.ValueType != nil {
			Walk(v, n.ValueType)
		}
		for _, f := range n.Fields {
			if f.Type != nil {
				Walk(v, f.Type)
			}
		}
	case *CreateIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *DropIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *BuildIndexStmt:
		Walk(v, n.Name)
		Walk(v, n.Table)
	case *CreateDatabaseStmt:
		Walk(v, n.Name)
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *AlterDatabaseStmt:
		Walk(v, n.Name)
		if n.NewName != nil {
			Walk(v, n.NewName)
		}
		for _, prop := range n.Properties {
			Walk(v, prop)
		}
	case *DropDatabaseStmt:
		Walk(v, n.Name)
	case *Property:
		// leaf node, no children

	// Expression nodes (T1.3).
	case *BinaryExpr:
		Walk(v, n.Left)
		Walk(v, n.Right)
	case *UnaryExpr:
		Walk(v, n.Expr)
	case *IsExpr:
		Walk(v, n.Expr)
	case *BetweenExpr:
		Walk(v, n.Expr)
		Walk(v, n.Low)
		Walk(v, n.High)
	case *InExpr:
		Walk(v, n.Expr)
		walkNodes(v, n.Values)
	case *LikeExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
		if n.Escape != nil {
			Walk(v, n.Escape)
		}
	case *RegexpExpr:
		Walk(v, n.Expr)
		Walk(v, n.Pattern)
	case *FuncCallExpr:
		Walk(v, n.Name)
		walkNodes(v, n.Args)
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
	case *CastExpr:
		Walk(v, n.Expr)
		Walk(v, n.TypeName)
	case *CaseExpr:
		if n.Operand != nil {
			Walk(v, n.Operand)
		}
		for _, w := range n.Whens {
			Walk(v, w)
		}
		if n.Else != nil {
			Walk(v, n.Else)
		}
	case *WhenClause:
		Walk(v, n.Cond)
		Walk(v, n.Result)
	case *SubqueryExpr:
		// placeholder leaf — no parsed children yet
	case *ColumnRef:
		Walk(v, n.Name)
	case *Literal:
		// leaf node, no children
	case *ParenExpr:
		Walk(v, n.Expr)
	case *ExistsExpr:
		Walk(v, n.Subquery)
	case *IntervalExpr:
		Walk(v, n.Value)
	case *OrderByItem:
		Walk(v, n.Expr)

	// SELECT statement nodes (T1.4).
	case *SelectStmt:
		for _, item := range n.Items {
			Walk(v, item)
		}
		walkNodes(v, n.From)
		if n.Where != nil {
			Walk(v, n.Where)
		}
		walkNodes(v, n.GroupBy)
		if n.Having != nil {
			Walk(v, n.Having)
		}
		if n.Qualify != nil {
			Walk(v, n.Qualify)
		}
		for _, item := range n.OrderBy {
			Walk(v, item)
		}
		if n.Limit != nil {
			Walk(v, n.Limit)
		}
		if n.Offset != nil {
			Walk(v, n.Offset)
		}
	case *SelectItem:
		if n.Expr != nil {
			Walk(v, n.Expr)
		}
		if n.TableName != nil {
			Walk(v, n.TableName)
		}
	case *TableRef:
		if n.Name != nil {
			Walk(v, n.Name)
		}
	case *JoinClause:
		Walk(v, n.Left)
		Walk(v, n.Right)
		if n.On != nil {
			Walk(v, n.On)
		}
	}
}
