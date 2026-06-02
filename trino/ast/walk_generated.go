package ast

// walkChildren walks the child nodes of node, calling Walk(v, child)
// for each child. Maintained manually until the node count warrants a
// code generator (mirroring doris/ast).
//
// Every concrete node type added under trino/ast must add a case here so
// the walker descends into its children (or an explicit leaf comment when
// it has none).
func walkChildren(v Visitor, node Node) {
	switch n := node.(type) {
	case *File:
		walkNodes(v, n.Stmts)
	case *QualifiedName:
		for _, p := range n.Parts {
			if p == nil {
				// Skip nil component parts so the visitor never receives a
				// typed-nil *Identifier (consistent with NormalizedParts /
				// String, which also skip nil parts).
				continue
			}
			Walk(v, p)
		}
	case *Identifier:
		// leaf node, no children
	}
}
