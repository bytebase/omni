package catalog

import (
	"fmt"
	"strings"
	"unicode"
)

// DeparseAnalyzedExpr converts an analyzed expression back to SQL text.
func DeparseAnalyzedExpr(expr AnalyzedExpr, q *Query) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case *VarExprQ:
		return deparseVarExprQ(e, q)
	case *ConstExprQ:
		return deparseConstExprQ(e)
	case *OpExprQ:
		return deparseOpExprQ(e, q)
	case *BoolExprQ:
		return deparseBoolExprQ(e, q)
	case *FuncCallExprQ:
		return deparseFuncCallExprQ(e, q)
	case *CaseExprQ:
		return deparseCaseExprQ(e, q)
	case *CoalesceExprQ:
		return deparseCoalesceExprQ(e, q)
	case *NullTestExprQ:
		return deparseNullTestExprQ(e, q)
	case *InListExprQ:
		return deparseInListExprQ(e, q)
	case *BetweenExprQ:
		return deparseBetweenExprQ(e, q)
	case *SubLinkExprQ:
		return deparseSubLinkExprQ(e, q)
	case *CastExprQ:
		return deparseCastExprQ(e, q)
	case *RowExprQ:
		return deparseRowExprQ(e, q)
	default:
		return fmt.Sprintf("/* unknown expr %T */", expr)
	}
}

func deparseVarExprQ(v *VarExprQ, q *Query) string {
	if q == nil || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		return "?"
	}
	rte := q.RangeTable[v.RangeIdx]
	colName := "?"
	if v.AttNum >= 1 && v.AttNum <= len(rte.ColNames) {
		colName = rte.ColNames[v.AttNum-1]
	}
	return backtickIDQ(rte.ERef) + "." + backtickIDQ(colName)
}

func deparseConstExprQ(c *ConstExprQ) string {
	if c.IsNull {
		return "NULL"
	}
	if isNumericLiteralQ(c.Value) || isBoolLiteralQ(c.Value) {
		return c.Value
	}
	return "'" + strings.ReplaceAll(c.Value, "'", "''") + "'"
}

func deparseOpExprQ(o *OpExprQ, q *Query) string {
	// Postfix operators: IS TRUE, IS FALSE, IS UNKNOWN, etc.
	if o.Left != nil && o.Right == nil {
		return "(" + DeparseAnalyzedExpr(o.Left, q) + " " + o.Op + ")"
	}
	// Prefix unary operators: -x, ~x, BINARY x
	if o.Left == nil && o.Right != nil {
		if o.Op == "-" || o.Op == "~" {
			return "(" + o.Op + DeparseAnalyzedExpr(o.Right, q) + ")"
		}
		return "(" + o.Op + " " + DeparseAnalyzedExpr(o.Right, q) + ")"
	}
	return "(" + DeparseAnalyzedExpr(o.Left, q) + " " + o.Op + " " + DeparseAnalyzedExpr(o.Right, q) + ")"
}

func deparseBoolExprQ(b *BoolExprQ, q *Query) string {
	switch b.Op {
	case BoolAnd:
		parts := make([]string, len(b.Args))
		for i, arg := range b.Args {
			parts[i] = DeparseAnalyzedExpr(arg, q)
		}
		return "(" + strings.Join(parts, " and ") + ")"
	case BoolOr:
		parts := make([]string, len(b.Args))
		for i, arg := range b.Args {
			parts[i] = DeparseAnalyzedExpr(arg, q)
		}
		return "(" + strings.Join(parts, " or ") + ")"
	case BoolNot:
		if len(b.Args) > 0 {
			return "(not " + DeparseAnalyzedExpr(b.Args[0], q) + ")"
		}
		return "(not ?)"
	default:
		return "?"
	}
}

func deparseFuncCallExprQ(f *FuncCallExprQ, q *Query) string {
	if f.IsAggregate && len(f.Args) == 0 && strings.ToLower(f.Name) == "count" {
		return "count(*)"
	}

	var sb strings.Builder
	sb.WriteString(f.Name)
	sb.WriteByte('(')

	if f.Distinct {
		sb.WriteString("distinct ")
	}

	for i, arg := range f.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(DeparseAnalyzedExpr(arg, q))
	}
	sb.WriteByte(')')

	if f.Over != nil {
		sb.WriteString(" over ")
		sb.WriteString(deparseWindowDefQ(f.Over, q))
	}

	return sb.String()
}

func deparseWindowDefQ(w *WindowDefQ, q *Query) string {
	if w.Name != "" && len(w.PartitionBy) == 0 && len(w.OrderBy) == 0 && w.FrameClause == "" {
		return backtickIDQ(w.Name)
	}

	var sb strings.Builder
	sb.WriteByte('(')
	if w.Name != "" {
		sb.WriteString(backtickIDQ(w.Name) + " ")
	}

	needSpace := false
	if len(w.PartitionBy) > 0 {
		sb.WriteString("partition by ")
		for i, expr := range w.PartitionBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(DeparseAnalyzedExpr(expr, q))
		}
		needSpace = true
	}
	if len(w.OrderBy) > 0 {
		if needSpace {
			sb.WriteByte(' ')
		}
		sb.WriteString("order by ")
		for i, sc := range w.OrderBy {
			if i > 0 {
				sb.WriteString(", ")
			}
			if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
				sb.WriteString(DeparseAnalyzedExpr(q.TargetList[sc.TargetIdx-1].Expr, q))
			}
			if sc.Descending {
				sb.WriteString(" desc")
			}
		}
		needSpace = true
	}
	if w.FrameClause != "" {
		if needSpace {
			sb.WriteByte(' ')
		}
		sb.WriteString(strings.ToLower(w.FrameClause))
	}
	sb.WriteByte(')')
	return sb.String()
}

func deparseCaseExprQ(c *CaseExprQ, q *Query) string {
	var sb strings.Builder
	sb.WriteString("case")
	if c.TestExpr != nil {
		sb.WriteByte(' ')
		sb.WriteString(DeparseAnalyzedExpr(c.TestExpr, q))
	}
	for _, w := range c.Args {
		sb.WriteString(" when ")
		sb.WriteString(DeparseAnalyzedExpr(w.Cond, q))
		sb.WriteString(" then ")
		sb.WriteString(DeparseAnalyzedExpr(w.Then, q))
	}
	if c.Default != nil {
		sb.WriteString(" else ")
		sb.WriteString(DeparseAnalyzedExpr(c.Default, q))
	}
	sb.WriteString(" end")
	return sb.String()
}

func deparseCoalesceExprQ(c *CoalesceExprQ, q *Query) string {
	var sb strings.Builder
	sb.WriteString("coalesce(")
	for i, arg := range c.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(DeparseAnalyzedExpr(arg, q))
	}
	sb.WriteByte(')')
	return sb.String()
}

func deparseNullTestExprQ(n *NullTestExprQ, q *Query) string {
	arg := DeparseAnalyzedExpr(n.Arg, q)
	if n.IsNull {
		return "(" + arg + " is null)"
	}
	return "(" + arg + " is not null)"
}

func deparseInListExprQ(in *InListExprQ, q *Query) string {
	var sb strings.Builder
	sb.WriteString(DeparseAnalyzedExpr(in.Arg, q))
	if in.Negated {
		sb.WriteString(" not in (")
	} else {
		sb.WriteString(" in (")
	}
	for i, item := range in.List {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(DeparseAnalyzedExpr(item, q))
	}
	sb.WriteByte(')')
	return sb.String()
}

func deparseBetweenExprQ(be *BetweenExprQ, q *Query) string {
	arg := DeparseAnalyzedExpr(be.Arg, q)
	lower := DeparseAnalyzedExpr(be.Lower, q)
	upper := DeparseAnalyzedExpr(be.Upper, q)
	if be.Negated {
		return "(" + arg + " not between " + lower + " and " + upper + ")"
	}
	return "(" + arg + " between " + lower + " and " + upper + ")"
}

func deparseSubLinkExprQ(s *SubLinkExprQ, q *Query) string {
	inner := DeparseQuery(s.Subquery)
	switch s.Kind {
	case SubLinkExists:
		return "exists (" + inner + ")"
	case SubLinkScalar:
		return "(" + inner + ")"
	case SubLinkIn:
		return DeparseAnalyzedExpr(s.TestExpr, q) + " in (" + inner + ")"
	case SubLinkAny:
		return DeparseAnalyzedExpr(s.TestExpr, q) + " " + s.Op + " any (" + inner + ")"
	case SubLinkAll:
		return DeparseAnalyzedExpr(s.TestExpr, q) + " " + s.Op + " all (" + inner + ")"
	default:
		return "(" + inner + ")"
	}
}

func deparseCastExprQ(c *CastExprQ, q *Query) string {
	arg := DeparseAnalyzedExpr(c.Arg, q)
	typeName := deparseResolvedTypeQ(c.TargetType)
	return "cast(" + arg + " as " + typeName + ")"
}

func deparseRowExprQ(r *RowExprQ, q *Query) string {
	var sb strings.Builder
	sb.WriteString("row(")
	for i, arg := range r.Args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(DeparseAnalyzedExpr(arg, q))
	}
	sb.WriteByte(')')
	return sb.String()
}

func deparseResolvedTypeQ(rt *ResolvedType) string {
	if rt == nil {
		return "char"
	}
	switch rt.BaseType {
	case BaseTypeBigInt:
		if rt.Unsigned {
			return "unsigned"
		}
		return "signed"
	case BaseTypeChar:
		if rt.Length > 0 {
			return fmt.Sprintf("char(%d)", rt.Length)
		}
		return "char"
	case BaseTypeBinary:
		if rt.Length > 0 {
			return fmt.Sprintf("binary(%d)", rt.Length)
		}
		return "binary"
	case BaseTypeDecimal:
		if rt.Precision > 0 && rt.Scale > 0 {
			return fmt.Sprintf("decimal(%d, %d)", rt.Precision, rt.Scale)
		}
		if rt.Precision > 0 {
			return fmt.Sprintf("decimal(%d)", rt.Precision)
		}
		return "decimal"
	case BaseTypeDate:
		return "date"
	case BaseTypeDateTime:
		return "datetime"
	case BaseTypeTime:
		return "time"
	case BaseTypeJSON:
		return "json"
	case BaseTypeFloat:
		return "float"
	case BaseTypeDouble:
		return "double"
	default:
		return "char"
	}
}

func isNumericLiteralQ(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' || s[0] == '+' {
		start = 1
	}
	if start >= len(s) {
		return false
	}
	hasDot := false
	for i := start; i < len(s); i++ {
		c := rune(s[i])
		if c == '.' {
			if hasDot {
				return false
			}
			hasDot = true
		} else if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func isBoolLiteralQ(s string) bool {
	upper := strings.ToUpper(s)
	return upper == "TRUE" || upper == "FALSE"
}
