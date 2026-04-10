package ast

// NodeLoc returns the source location of n, or NoLoc() if n is nil or its
// concrete type carries no Loc field. The function is a type-switch dispatcher
// — every concrete node type added under snowflake/ast must add a case here.
//
// The pattern matches pg/ast.NodeLoc.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	case *File:
		return v.Loc
	case *ObjectName:
		return v.Loc
	default:
		return NoLoc()
	}
}

// SpanNodes returns the smallest Loc that covers every node in nodes.
// Nil entries and nodes whose Loc is invalid are skipped. Returns NoLoc()
// when no node has a valid Loc (including the empty-args case).
func SpanNodes(nodes ...Node) Loc {
	out := NoLoc()
	for _, n := range nodes {
		if n == nil {
			continue
		}
		out = out.Merge(NodeLoc(n))
	}
	return out
}
