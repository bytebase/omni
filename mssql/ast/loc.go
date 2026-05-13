package ast

import "reflect"

// NodeLoc returns the source location of n, or NoLoc() if n is nil or its
// concrete type carries no Loc field.
func NodeLoc(n Node) Loc {
	if n == nil {
		return NoLoc()
	}
	switch v := n.(type) {
	case *List:
		return SpanNodes(v.Items...)
	case *SelectStmt:
		return v.Loc
	case *ResTarget:
		return v.Loc
	case *BinaryExpr:
		return v.Loc
	case *UnaryExpr:
		return v.Loc
	case *FuncCallExpr:
		return v.Loc
	case *DatePart:
		return v.Loc
	case *NextValueForExpr:
		return v.Loc
	case *ParseExpr:
		return v.Loc
	case *JsonKeyValueExpr:
		return v.Loc
	case *CaseExpr:
		return v.Loc
	case *CaseWhen:
		return v.Loc
	case *BetweenExpr:
		return v.Loc
	case *InExpr:
		return v.Loc
	case *LikeExpr:
		return v.Loc
	case *IsExpr:
		return v.Loc
	case *ExistsExpr:
		return v.Loc
	case *FullTextPredicate:
		return v.Loc
	case *CastExpr:
		return v.Loc
	case *ConvertExpr:
		return v.Loc
	case *TryCastExpr:
		return v.Loc
	case *TryConvertExpr:
		return v.Loc
	case *CoalesceExpr:
		return v.Loc
	case *NullifExpr:
		return v.Loc
	case *IifExpr:
		return v.Loc
	case *ColumnRef:
		return v.Loc
	case *VariableRef:
		return v.Loc
	case *StarExpr:
		return v.Loc
	case *Literal:
		return v.Loc
	case *SubqueryExpr:
		return v.Loc
	case *SubqueryComparisonExpr:
		return v.Loc
	case *CollateExpr:
		return v.Loc
	case *AtTimeZoneExpr:
		return v.Loc
	case *ParenExpr:
		return v.Loc
	case *TableRef:
		return v.Loc
	case *DataType:
		return v.Loc
	case *OrderByItem:
		return v.Loc
	case *OverClause:
		return v.Loc
	case *CurrentOfExpr:
		return v.Loc
	case *MethodCallExpr:
		return v.Loc
	case *GroupingSetsExpr:
		return v.Loc
	case *RollupExpr:
		return v.Loc
	case *CubeExpr:
		return v.Loc
	default:
		return nodeLocFromField(n)
	}
}

func nodeLocFromField(n Node) Loc {
	v := reflect.ValueOf(n)
	if !v.IsValid() {
		return NoLoc()
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NoLoc()
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return NoLoc()
	}
	field := v.FieldByName("Loc")
	if !field.IsValid() || field.Type() != reflect.TypeOf(Loc{}) {
		return NoLoc()
	}
	loc, ok := field.Interface().(Loc)
	if !ok {
		return NoLoc()
	}
	return loc
}

// SpanNodes returns the smallest Loc that covers every node in nodes.
// Nil entries and nodes whose Loc is invalid are skipped. Returns NoLoc()
// when no node has a valid Loc.
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
