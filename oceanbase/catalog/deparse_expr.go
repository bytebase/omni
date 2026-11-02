package catalog

import (
	"fmt"
	"strings"
	"unicode"
)

// DeparseAnalyzedExpr converts an analyzed expression back to SQL text.
func DeparseAnalyzedExpr(expr AnalyzedExpr, q *Query) string {
	var b strings.Builder
	writeAnalyzedExpr(&b, expr, q)
	return b.String()
}

// writeAnalyzedExpr writes an analyzed expression into the one builder of the
// call, like every function below: returning each level as a string of its
// own copied the text of a long chain (1 + 1 + ..., a AND b AND ...) once per
// level.
func writeAnalyzedExpr(b *strings.Builder, expr AnalyzedExpr, q *Query) {
	if expr == nil {
		return
	}
	switch e := expr.(type) {
	case *VarExprQ:
		writeVarExprQ(b, e, q)
	case *ConstExprQ:
		writeConstExprQ(b, e)
	case *OpExprQ:
		writeOpExprQ(b, e, q)
	case *BoolExprQ:
		writeBoolExprQ(b, e, q)
	case *FuncCallExprQ:
		writeFuncCallExprQ(b, e, q)
	case *CaseExprQ:
		writeCaseExprQ(b, e, q)
	case *CoalesceExprQ:
		writeCoalesceExprQ(b, e, q)
	case *NullTestExprQ:
		writeNullTestExprQ(b, e, q)
	case *InListExprQ:
		writeInListExprQ(b, e, q)
	case *BetweenExprQ:
		writeBetweenExprQ(b, e, q)
	case *SubLinkExprQ:
		writeSubLinkExprQ(b, e, q)
	case *CastExprQ:
		writeCastExprQ(b, e, q)
	case *RowExprQ:
		writeRowExprQ(b, e, q)
	default:
		fmt.Fprintf(b, "/* unknown expr %T */", expr)
	}
}

func writeVarExprQ(b *strings.Builder, v *VarExprQ, q *Query) {
	if q == nil || v.RangeIdx < 0 || v.RangeIdx >= len(q.RangeTable) {
		b.WriteString("?")
		return
	}
	rte := q.RangeTable[v.RangeIdx]
	colName := "?"
	if v.AttNum >= 1 && v.AttNum <= len(rte.ColNames) {
		colName = rte.ColNames[v.AttNum-1]
	}
	b.WriteString(backtickIDQ(rte.ERef))
	b.WriteString(".")
	b.WriteString(backtickIDQ(colName))
}

func writeConstExprQ(b *strings.Builder, c *ConstExprQ) {
	if c.IsNull {
		b.WriteString("NULL")
		return
	}
	if isNumericLiteralQ(c.Value) || isBoolLiteralQ(c.Value) {
		b.WriteString(c.Value)
		return
	}
	b.WriteString("'")
	b.WriteString(strings.ReplaceAll(c.Value, "'", "''"))
	b.WriteString("'")
}

func writeOpExprQ(b *strings.Builder, o *OpExprQ, q *Query) {
	b.WriteString("(")
	switch {
	case o.Left != nil && o.Right == nil:
		// Postfix operators: IS TRUE, IS FALSE, IS UNKNOWN, etc.
		writeAnalyzedExpr(b, o.Left, q)
		b.WriteString(" ")
		b.WriteString(o.Op)
	case o.Left == nil && o.Right != nil:
		// Prefix unary operators: -x, ~x, BINARY x
		b.WriteString(o.Op)
		if o.Op != "-" && o.Op != "~" {
			b.WriteString(" ")
		}
		writeAnalyzedExpr(b, o.Right, q)
	default:
		writeAnalyzedExpr(b, o.Left, q)
		b.WriteString(" ")
		b.WriteString(o.Op)
		b.WriteString(" ")
		writeAnalyzedExpr(b, o.Right, q)
	}
	b.WriteString(")")
}

func writeBoolExprQ(b *strings.Builder, e *BoolExprQ, q *Query) {
	var sep string
	switch e.Op {
	case BoolAnd:
		sep = " and "
	case BoolOr:
		sep = " or "
	case BoolNot:
		if len(e.Args) > 0 {
			b.WriteString("(not ")
			writeAnalyzedExpr(b, e.Args[0], q)
			b.WriteString(")")
			return
		}
		b.WriteString("(not ?)")
		return
	default:
		b.WriteString("?")
		return
	}
	b.WriteString("(")
	for i, arg := range e.Args {
		if i > 0 {
			b.WriteString(sep)
		}
		writeAnalyzedExpr(b, arg, q)
	}
	b.WriteString(")")
}

func writeFuncCallExprQ(b *strings.Builder, f *FuncCallExprQ, q *Query) {
	if f.IsAggregate && len(f.Args) == 0 && strings.ToLower(f.Name) == "count" {
		b.WriteString("count(*)")
		return
	}

	b.WriteString(f.Name)
	b.WriteByte('(')

	if f.Distinct {
		b.WriteString("distinct ")
	}

	for i, arg := range f.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		writeAnalyzedExpr(b, arg, q)
	}
	b.WriteByte(')')

	if f.Over != nil {
		b.WriteString(" over ")
		writeWindowDefQ(b, f.Over, q)
	}
}

func writeWindowDefQ(b *strings.Builder, w *WindowDefQ, q *Query) {
	if w.Name != "" && len(w.PartitionBy) == 0 && len(w.OrderBy) == 0 && w.FrameClause == "" {
		b.WriteString(backtickIDQ(w.Name))
		return
	}

	b.WriteByte('(')
	if w.Name != "" {
		b.WriteString(backtickIDQ(w.Name) + " ")
	}

	needSpace := false
	if len(w.PartitionBy) > 0 {
		b.WriteString("partition by ")
		for i, expr := range w.PartitionBy {
			if i > 0 {
				b.WriteString(", ")
			}
			writeAnalyzedExpr(b, expr, q)
		}
		needSpace = true
	}
	if len(w.OrderBy) > 0 {
		if needSpace {
			b.WriteByte(' ')
		}
		b.WriteString("order by ")
		for i, sc := range w.OrderBy {
			if i > 0 {
				b.WriteString(", ")
			}
			if sc.TargetIdx >= 1 && sc.TargetIdx <= len(q.TargetList) {
				writeAnalyzedExpr(b, q.TargetList[sc.TargetIdx-1].Expr, q)
			}
			if sc.Descending {
				b.WriteString(" desc")
			}
		}
		needSpace = true
	}
	if w.FrameClause != "" {
		if needSpace {
			b.WriteByte(' ')
		}
		b.WriteString(strings.ToLower(w.FrameClause))
	}
	b.WriteByte(')')
}

func writeCaseExprQ(b *strings.Builder, c *CaseExprQ, q *Query) {
	b.WriteString("case")
	if c.TestExpr != nil {
		b.WriteByte(' ')
		writeAnalyzedExpr(b, c.TestExpr, q)
	}
	for _, w := range c.Args {
		b.WriteString(" when ")
		writeAnalyzedExpr(b, w.Cond, q)
		b.WriteString(" then ")
		writeAnalyzedExpr(b, w.Then, q)
	}
	if c.Default != nil {
		b.WriteString(" else ")
		writeAnalyzedExpr(b, c.Default, q)
	}
	b.WriteString(" end")
}

func writeCoalesceExprQ(b *strings.Builder, c *CoalesceExprQ, q *Query) {
	b.WriteString("coalesce(")
	for i, arg := range c.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		writeAnalyzedExpr(b, arg, q)
	}
	b.WriteByte(')')
}

func writeNullTestExprQ(b *strings.Builder, n *NullTestExprQ, q *Query) {
	b.WriteString("(")
	writeAnalyzedExpr(b, n.Arg, q)
	if n.IsNull {
		b.WriteString(" is null)")
	} else {
		b.WriteString(" is not null)")
	}
}

func writeInListExprQ(b *strings.Builder, in *InListExprQ, q *Query) {
	writeAnalyzedExpr(b, in.Arg, q)
	if in.Negated {
		b.WriteString(" not in (")
	} else {
		b.WriteString(" in (")
	}
	for i, item := range in.List {
		if i > 0 {
			b.WriteString(", ")
		}
		writeAnalyzedExpr(b, item, q)
	}
	b.WriteByte(')')
}

func writeBetweenExprQ(b *strings.Builder, be *BetweenExprQ, q *Query) {
	b.WriteString("(")
	writeAnalyzedExpr(b, be.Arg, q)
	if be.Negated {
		b.WriteString(" not between ")
	} else {
		b.WriteString(" between ")
	}
	writeAnalyzedExpr(b, be.Lower, q)
	b.WriteString(" and ")
	writeAnalyzedExpr(b, be.Upper, q)
	b.WriteString(")")
}

func writeSubLinkExprQ(b *strings.Builder, s *SubLinkExprQ, q *Query) {
	switch s.Kind {
	case SubLinkExists:
		b.WriteString("exists (")
	case SubLinkIn:
		writeAnalyzedExpr(b, s.TestExpr, q)
		b.WriteString(" in (")
	case SubLinkAny:
		writeAnalyzedExpr(b, s.TestExpr, q)
		b.WriteString(" ")
		b.WriteString(s.Op)
		b.WriteString(" any (")
	case SubLinkAll:
		writeAnalyzedExpr(b, s.TestExpr, q)
		b.WriteString(" ")
		b.WriteString(s.Op)
		b.WriteString(" all (")
	default: // SubLinkScalar
		b.WriteString("(")
	}
	writeQueryQ(b, s.Subquery)
	b.WriteString(")")
}

func writeCastExprQ(b *strings.Builder, c *CastExprQ, q *Query) {
	b.WriteString("cast(")
	writeAnalyzedExpr(b, c.Arg, q)
	b.WriteString(" as ")
	b.WriteString(deparseResolvedTypeQ(c.TargetType))
	b.WriteString(")")
}

func writeRowExprQ(b *strings.Builder, r *RowExprQ, q *Query) {
	b.WriteString("row(")
	for i, arg := range r.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		writeAnalyzedExpr(b, arg, q)
	}
	b.WriteByte(')')
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
