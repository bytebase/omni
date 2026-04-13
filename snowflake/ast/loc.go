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
	case *TypeName:
		return v.Loc
	case *Literal:
		return v.Loc
	case *ColumnRef:
		return v.Loc
	case *StarExpr:
		return v.Loc
	case *BinaryExpr:
		return v.Loc
	case *UnaryExpr:
		return v.Loc
	case *ParenExpr:
		return v.Loc
	case *CastExpr:
		return v.Loc
	case *CaseExpr:
		return v.Loc
	case *FuncCallExpr:
		return v.Loc
	case *IffExpr:
		return v.Loc
	case *CollateExpr:
		return v.Loc
	case *IsExpr:
		return v.Loc
	case *BetweenExpr:
		return v.Loc
	case *InExpr:
		return v.Loc
	case *LikeExpr:
		return v.Loc
	case *AccessExpr:
		return v.Loc
	case *ArrayLiteralExpr:
		return v.Loc
	case *JsonLiteralExpr:
		return v.Loc
	case *LambdaExpr:
		return v.Loc
	case *SubqueryExpr:
		return v.Loc
	case *ExistsExpr:
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
