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
	}
}
