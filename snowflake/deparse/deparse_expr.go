package deparse

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bytebase/omni/snowflake/ast"
)

// writeExpr dispatches to the appropriate expression writer.
func (w *writer) writeExpr(node ast.Node) error {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.Literal:
		return w.writeLiteral(n)
	case *ast.ColumnRef:
		return w.writeColumnRef(n)
	case *ast.StarExpr:
		return w.writeStarExpr(n)
	case *ast.BinaryExpr:
		return w.writeBinaryExpr(n)
	case *ast.UnaryExpr:
		return w.writeUnaryExpr(n)
	case *ast.ParenExpr:
		return w.writeParenExpr(n)
	case *ast.CastExpr:
		return w.writeCastExpr(n)
	case *ast.CaseExpr:
		return w.writeCaseExpr(n)
	case *ast.FuncCallExpr:
		return w.writeFuncCallExpr(n)
	case *ast.IffExpr:
		return w.writeIffExpr(n)
	case *ast.CollateExpr:
		return w.writeCollateExpr(n)
	case *ast.IsExpr:
		return w.writeIsExpr(n)
	case *ast.BetweenExpr:
		return w.writeBetweenExpr(n)
	case *ast.InExpr:
		return w.writeInExpr(n)
	case *ast.LikeExpr:
		return w.writeLikeExpr(n)
	case *ast.AccessExpr:
		return w.writeAccessExpr(n)
	case *ast.ArrayLiteralExpr:
		return w.writeArrayLiteralExpr(n)
	case *ast.JsonLiteralExpr:
		return w.writeJsonLiteralExpr(n)
	case *ast.LambdaExpr:
		return w.writeLambdaExpr(n)
	case *ast.SubqueryExpr:
		return w.writeSubqueryExpr(n)
	case *ast.ExistsExpr:
		return w.writeExistsExpr(n)
	default:
		return fmt.Errorf("deparse: unsupported expression node type %T", node)
	}
}

func (w *writer) writeLiteral(n *ast.Literal) error {
	w.ensureSpace()
	switch n.Kind {
	case ast.LitNull:
		w.buf.WriteString("NULL")
	case ast.LitBool:
		if n.Bval {
			w.buf.WriteString("TRUE")
		} else {
			w.buf.WriteString("FALSE")
		}
	case ast.LitInt:
		w.buf.WriteString(strconv.FormatInt(n.Ival, 10))
	case ast.LitFloat:
		w.buf.WriteString(n.Value)
	case ast.LitString:
		w.buf.WriteString(quoteString(n.Value))
	default:
		// fallback: use raw Value
		w.buf.WriteString(n.Value)
	}
	return nil
}

// quoteString wraps a string value in single quotes, escaping internal single quotes.
func quoteString(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func (w *writer) writeColumnRef(n *ast.ColumnRef) error {
	w.ensureSpace()
	for i, part := range n.Parts {
		if i > 0 {
			w.buf.WriteByte('.')
		}
		w.buf.WriteString(part.String())
	}
	return nil
}

func (w *writer) writeStarExpr(n *ast.StarExpr) error {
	w.ensureSpace()
	if n.Qualifier != nil {
		w.writeObjectNameNoSpace(n.Qualifier)
		w.buf.WriteString(".*")
	} else {
		w.buf.WriteByte('*')
	}
	return nil
}

func (w *writer) writeBinaryExpr(n *ast.BinaryExpr) error {
	if err := w.writeExpr(n.Left); err != nil {
		return err
	}
	op := n.Op.String()
	switch n.Op {
	case ast.BinAnd, ast.BinOr:
		w.buf.WriteByte(' ')
		w.buf.WriteString(op)
	default:
		w.buf.WriteByte(' ')
		w.buf.WriteString(op)
	}
	w.buf.WriteByte(' ')
	return w.writeExpr(n.Right)
}

func (w *writer) writeUnaryExpr(n *ast.UnaryExpr) error {
	w.ensureSpace()
	switch n.Op {
	case ast.UnaryNot:
		w.buf.WriteString("NOT ")
	case ast.UnaryMinus:
		w.buf.WriteByte('-')
	case ast.UnaryPlus:
		w.buf.WriteByte('+')
	}
	return w.writeExpr(n.Expr)
}

func (w *writer) writeParenExpr(n *ast.ParenExpr) error {
	w.ensureSpace()
	w.buf.WriteByte('(')
	savedLen := w.buf.Len()
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	// trim leading space that ensureSpace may have added inside the paren
	s := w.buf.String()
	inner := s[savedLen:]
	if strings.HasPrefix(inner, " ") {
		trimmed := s[:savedLen] + inner[1:]
		w.buf.Reset()
		w.buf.WriteString(trimmed)
	}
	w.buf.WriteByte(')')
	return nil
}

func (w *writer) writeCastExpr(n *ast.CastExpr) error {
	if n.ColonColon {
		if err := w.writeExpr(n.Expr); err != nil {
			return err
		}
		w.buf.WriteString("::")
		w.writeTypeName(n.TypeName)
		return nil
	}
	w.ensureSpace()
	if n.TryCast {
		w.buf.WriteString("TRY_CAST(")
	} else {
		w.buf.WriteString("CAST(")
	}
	savedLen := w.buf.Len()
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	s := w.buf.String()
	inner := s[savedLen:]
	if strings.HasPrefix(inner, " ") {
		trimmed := s[:savedLen] + inner[1:]
		w.buf.Reset()
		w.buf.WriteString(trimmed)
	}
	w.buf.WriteString(" AS ")
	w.writeTypeName(n.TypeName)
	w.buf.WriteByte(')')
	return nil
}

func (w *writer) writeCaseExpr(n *ast.CaseExpr) error {
	w.ensureSpace()
	w.buf.WriteString("CASE")
	if n.Kind == ast.CaseSimple && n.Operand != nil {
		if err := w.writeExpr(n.Operand); err != nil {
			return err
		}
	}
	for _, when := range n.Whens {
		w.buf.WriteString(" WHEN ")
		savedLen := w.buf.Len()
		if err := w.writeExpr(when.Cond); err != nil {
			return err
		}
		s := w.buf.String()
		inner := s[savedLen:]
		if strings.HasPrefix(inner, " ") {
			trimmed := s[:savedLen] + inner[1:]
			w.buf.Reset()
			w.buf.WriteString(trimmed)
		}
		w.buf.WriteString(" THEN ")
		savedLen = w.buf.Len()
		if err := w.writeExpr(when.Result); err != nil {
			return err
		}
		s = w.buf.String()
		inner = s[savedLen:]
		if strings.HasPrefix(inner, " ") {
			trimmed := s[:savedLen] + inner[1:]
			w.buf.Reset()
			w.buf.WriteString(trimmed)
		}
	}
	if n.Else != nil {
		w.buf.WriteString(" ELSE ")
		savedLen := w.buf.Len()
		if err := w.writeExpr(n.Else); err != nil {
			return err
		}
		s := w.buf.String()
		inner := s[savedLen:]
		if strings.HasPrefix(inner, " ") {
			trimmed := s[:savedLen] + inner[1:]
			w.buf.Reset()
			w.buf.WriteString(trimmed)
		}
	}
	w.buf.WriteString(" END")
	return nil
}

func (w *writer) writeFuncCallExpr(n *ast.FuncCallExpr) error {
	w.ensureSpace()
	w.writeObjectNameNoSpace(&n.Name)
	w.buf.WriteByte('(')
	if n.Star {
		w.buf.WriteByte('*')
	} else {
		if n.Distinct {
			w.buf.WriteString("DISTINCT ")
		}
		for i, arg := range n.Args {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			savedLen := w.buf.Len()
			if err := w.writeExpr(arg); err != nil {
				return err
			}
			s := w.buf.String()
			inner := s[savedLen:]
			if strings.HasPrefix(inner, " ") {
				trimmed := s[:savedLen] + inner[1:]
				w.buf.Reset()
				w.buf.WriteString(trimmed)
			}
		}
	}
	if len(n.OrderBy) > 0 {
		// WITHIN GROUP (ORDER BY ...) for ordered-set aggregate functions
		w.buf.WriteString(") WITHIN GROUP (ORDER BY ")
		for i, item := range n.OrderBy {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeOrderItem(item); err != nil {
				return err
			}
		}
		w.buf.WriteByte(')')
	} else {
		w.buf.WriteByte(')')
	}
	if n.Over != nil {
		w.buf.WriteString(" OVER (")
		if err := w.writeWindowSpec(n.Over); err != nil {
			return err
		}
		w.buf.WriteByte(')')
	}
	return nil
}

func (w *writer) writeWindowSpec(ws *ast.WindowSpec) error {
	first := true
	if len(ws.PartitionBy) > 0 {
		w.buf.WriteString("PARTITION BY ")
		for i, expr := range ws.PartitionBy {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			savedLen := w.buf.Len()
			if err := w.writeExpr(expr); err != nil {
				return err
			}
			s := w.buf.String()
			inner := s[savedLen:]
			if strings.HasPrefix(inner, " ") {
				trimmed := s[:savedLen] + inner[1:]
				w.buf.Reset()
				w.buf.WriteString(trimmed)
			}
		}
		first = false
	}
	if len(ws.OrderBy) > 0 {
		if !first {
			w.buf.WriteByte(' ')
		}
		w.buf.WriteString("ORDER BY ")
		for i, item := range ws.OrderBy {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			if err := w.writeOrderItem(item); err != nil {
				return err
			}
		}
		first = false
	}
	if ws.Frame != nil {
		if !first {
			w.buf.WriteByte(' ')
		}
		if err := w.writeWindowFrame(ws.Frame); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeWindowFrame(f *ast.WindowFrame) error {
	switch f.Kind {
	case ast.FrameRows:
		w.buf.WriteString("ROWS")
	case ast.FrameRange:
		w.buf.WriteString("RANGE")
	case ast.FrameGroups:
		w.buf.WriteString("GROUPS")
	}
	w.buf.WriteString(" BETWEEN ")
	if err := w.writeWindowBound(f.Start); err != nil {
		return err
	}
	w.buf.WriteString(" AND ")
	return w.writeWindowBound(f.End)
}

func (w *writer) writeWindowBound(b ast.WindowBound) error {
	switch b.Kind {
	case ast.BoundUnboundedPreceding:
		w.buf.WriteString("UNBOUNDED PRECEDING")
	case ast.BoundPreceding:
		if b.Offset != nil {
			savedLen := w.buf.Len()
			if err := w.writeExpr(b.Offset); err != nil {
				return err
			}
			s := w.buf.String()
			inner := s[savedLen:]
			if strings.HasPrefix(inner, " ") {
				trimmed := s[:savedLen] + inner[1:]
				w.buf.Reset()
				w.buf.WriteString(trimmed)
			}
		} else {
			w.buf.WriteByte('1')
		}
		w.buf.WriteString(" PRECEDING")
	case ast.BoundCurrentRow:
		w.buf.WriteString("CURRENT ROW")
	case ast.BoundFollowing:
		if b.Offset != nil {
			savedLen := w.buf.Len()
			if err := w.writeExpr(b.Offset); err != nil {
				return err
			}
			s := w.buf.String()
			inner := s[savedLen:]
			if strings.HasPrefix(inner, " ") {
				trimmed := s[:savedLen] + inner[1:]
				w.buf.Reset()
				w.buf.WriteString(trimmed)
			}
		} else {
			w.buf.WriteByte('1')
		}
		w.buf.WriteString(" FOLLOWING")
	case ast.BoundUnboundedFollowing:
		w.buf.WriteString("UNBOUNDED FOLLOWING")
	}
	return nil
}

func (w *writer) writeIffExpr(n *ast.IffExpr) error {
	w.ensureSpace()
	w.buf.WriteString("IFF(")
	if err := writeExprNoLeadSpace(w, n.Cond); err != nil {
		return err
	}
	w.buf.WriteString(", ")
	if err := writeExprNoLeadSpace(w, n.Then); err != nil {
		return err
	}
	w.buf.WriteString(", ")
	if err := writeExprNoLeadSpace(w, n.Else); err != nil {
		return err
	}
	w.buf.WriteByte(')')
	return nil
}

func (w *writer) writeCollateExpr(n *ast.CollateExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	w.buf.WriteString(" COLLATE ")
	w.buf.WriteString(quoteString(n.Collation))
	return nil
}

func (w *writer) writeIsExpr(n *ast.IsExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	if n.Null {
		if n.Not {
			w.buf.WriteString(" IS NOT NULL")
		} else {
			w.buf.WriteString(" IS NULL")
		}
	} else {
		if n.Not {
			w.buf.WriteString(" IS NOT DISTINCT FROM ")
		} else {
			w.buf.WriteString(" IS DISTINCT FROM ")
		}
		return writeExprNoLeadSpace(w, n.DistinctFrom)
	}
	return nil
}

func (w *writer) writeBetweenExpr(n *ast.BetweenExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	if n.Not {
		w.buf.WriteString(" NOT BETWEEN ")
	} else {
		w.buf.WriteString(" BETWEEN ")
	}
	if err := writeExprNoLeadSpace(w, n.Low); err != nil {
		return err
	}
	w.buf.WriteString(" AND ")
	return writeExprNoLeadSpace(w, n.High)
}

func (w *writer) writeInExpr(n *ast.InExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	if n.Not {
		w.buf.WriteString(" NOT IN (")
	} else {
		w.buf.WriteString(" IN (")
	}
	for i, val := range n.Values {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := writeExprNoLeadSpace(w, val); err != nil {
			return err
		}
	}
	w.buf.WriteByte(')')
	return nil
}

func (w *writer) writeLikeExpr(n *ast.LikeExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	notStr := ""
	if n.Not {
		notStr = "NOT "
	}
	switch n.Op {
	case ast.LikeOpLike:
		if n.Any {
			w.buf.WriteString(" " + notStr + "LIKE ANY (")
			for i, v := range n.AnyValues {
				if i > 0 {
					w.buf.WriteString(", ")
				}
				if err := writeExprNoLeadSpace(w, v); err != nil {
					return err
				}
			}
			w.buf.WriteByte(')')
			return nil
		}
		w.buf.WriteString(" " + notStr + "LIKE ")
	case ast.LikeOpILike:
		w.buf.WriteString(" " + notStr + "ILIKE ")
	case ast.LikeOpRLike:
		w.buf.WriteString(" " + notStr + "RLIKE ")
	case ast.LikeOpRegexp:
		w.buf.WriteString(" " + notStr + "REGEXP ")
	}
	if err := writeExprNoLeadSpace(w, n.Pattern); err != nil {
		return err
	}
	if n.Escape != nil {
		w.buf.WriteString(" ESCAPE ")
		if err := writeExprNoLeadSpace(w, n.Escape); err != nil {
			return err
		}
	}
	return nil
}

func (w *writer) writeAccessExpr(n *ast.AccessExpr) error {
	if err := w.writeExpr(n.Expr); err != nil {
		return err
	}
	switch n.Kind {
	case ast.AccessColon:
		w.buf.WriteByte(':')
		w.buf.WriteString(n.Field.String())
	case ast.AccessDot:
		w.buf.WriteByte('.')
		w.buf.WriteString(n.Field.String())
	case ast.AccessBracket:
		w.buf.WriteByte('[')
		if err := writeExprNoLeadSpace(w, n.Index); err != nil {
			return err
		}
		w.buf.WriteByte(']')
	}
	return nil
}

func (w *writer) writeArrayLiteralExpr(n *ast.ArrayLiteralExpr) error {
	w.ensureSpace()
	w.buf.WriteByte('[')
	for i, elem := range n.Elements {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		if err := writeExprNoLeadSpace(w, elem); err != nil {
			return err
		}
	}
	w.buf.WriteByte(']')
	return nil
}

func (w *writer) writeJsonLiteralExpr(n *ast.JsonLiteralExpr) error {
	w.ensureSpace()
	w.buf.WriteByte('{')
	for i, pair := range n.Pairs {
		if i > 0 {
			w.buf.WriteString(", ")
		}
		w.buf.WriteString(quoteString(pair.Key))
		w.buf.WriteString(": ")
		if err := writeExprNoLeadSpace(w, pair.Value); err != nil {
			return err
		}
	}
	w.buf.WriteByte('}')
	return nil
}

func (w *writer) writeLambdaExpr(n *ast.LambdaExpr) error {
	w.ensureSpace()
	if len(n.Params) == 1 {
		w.buf.WriteString(n.Params[0].String())
	} else {
		w.buf.WriteByte('(')
		for i, p := range n.Params {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(p.String())
		}
		w.buf.WriteByte(')')
	}
	w.buf.WriteString(" -> ")
	return writeExprNoLeadSpace(w, n.Body)
}

func (w *writer) writeSubqueryExpr(n *ast.SubqueryExpr) error {
	w.ensureSpace()
	w.buf.WriteByte('(')
	if err := w.writeQueryExpr(n.Query); err != nil {
		return err
	}
	w.buf.WriteByte(')')
	return nil
}

func (w *writer) writeExistsExpr(n *ast.ExistsExpr) error {
	w.ensureSpace()
	w.buf.WriteString("EXISTS (")
	if err := w.writeQueryExpr(n.Query); err != nil {
		return err
	}
	w.buf.WriteByte(')')
	return nil
}

// writeOrderItem writes a single ORDER BY item.
func (w *writer) writeOrderItem(item *ast.OrderItem) error {
	if err := writeExprNoLeadSpace(w, item.Expr); err != nil {
		return err
	}
	if item.Desc {
		w.buf.WriteString(" DESC")
	} else {
		w.buf.WriteString(" ASC")
	}
	if item.NullsFirst != nil {
		if *item.NullsFirst {
			w.buf.WriteString(" NULLS FIRST")
		} else {
			w.buf.WriteString(" NULLS LAST")
		}
	}
	return nil
}

// writeTypeName writes a data type name.
func (w *writer) writeTypeName(tn *ast.TypeName) {
	if tn == nil {
		return
	}
	w.buf.WriteString(tn.Name)
	if tn.Kind == ast.TypeVector && tn.ElementType != nil {
		w.buf.WriteByte('(')
		w.buf.WriteString(tn.ElementType.Name)
		w.buf.WriteString(", ")
		w.buf.WriteString(strconv.Itoa(tn.VectorDim))
		w.buf.WriteByte(')')
		return
	}
	if tn.Kind == ast.TypeArray && tn.ElementType != nil {
		w.buf.WriteByte('(')
		w.writeTypeName(tn.ElementType)
		w.buf.WriteByte(')')
		return
	}
	if len(tn.Params) > 0 {
		w.buf.WriteByte('(')
		for i, p := range tn.Params {
			if i > 0 {
				w.buf.WriteString(", ")
			}
			w.buf.WriteString(strconv.Itoa(p))
		}
		w.buf.WriteByte(')')
	}
}

// writeExprNoLeadSpace writes expr and strips any leading space that
// ensureSpace may have inserted. This is used inside parentheses, comma
// lists, or after literal keyword+space sequences where the caller has
// already positioned the cursor correctly.
func writeExprNoLeadSpace(w *writer, node ast.Node) error {
	start := w.buf.Len()
	if err := w.writeExpr(node); err != nil {
		return err
	}
	s := w.buf.String()
	if len(s) > start && s[start] == ' ' {
		trimmed := s[:start] + s[start+1:]
		w.buf.Reset()
		w.buf.WriteString(trimmed)
	}
	return nil
}
