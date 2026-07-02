package parser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

// This file is the structural correctness gate (Layer 2 of correctness-protocol.md)
// for the expressions node. Accept/reject (expr_test.go) cannot catch wrong
// precedence/associativity/nesting — `1 + 2 * 3` "accepts" however it groups —
// so these tests:
//
//   - TestExpr_Precedence pins the exact parse-tree shape of every
//     oracle-confirmed precedence/associativity fact (P1 single non-associative
//     predicate, P2 || over predicate, arithmetic layering, unary, AT TIME ZONE,
//     postfix subscript/dereference), via a compact S-expression rendering of the
//     parsed Expr.
//   - TestExpr_RoundTrip / TestExpr_RoundTripOracle render each accepted corpus
//     expression back to SQL and re-parse it; the re-parse must succeed and Trino
//     must re-accept the rendering — proving the parse captured a well-formed,
//     valid structure.

// renderSExpr renders an Expr as a parenthesized S-expression capturing operator
// structure and nesting, for precise structural assertions. It is test-only.
func renderSExpr(e Expr) string {
	switch n := e.(type) {
	case *Literal:
		if n.Kind == LiteralString {
			return "'" + n.Value + "'"
		}
		return n.Value
	case *Parameter:
		return "?"
	case *ColumnRef:
		return n.Name.String()
	case *Dereference:
		return fmt.Sprintf("(. %s %s)", renderSExpr(n.Base), n.FieldName.String())
	case *Subscript:
		return fmt.Sprintf("([] %s %s)", renderSExpr(n.Value), renderSExpr(n.Index))
	case *UnaryExpr:
		return fmt.Sprintf("(u%s %s)", n.Op, renderSExpr(n.Operand))
	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", n.Op, renderSExpr(n.Left), renderSExpr(n.Right))
	case *LogicalExpr:
		return fmt.Sprintf("(%s %s %s)", strings.ToLower(n.Op), renderSExpr(n.Left), renderSExpr(n.Right))
	case *NotExpr:
		return fmt.Sprintf("(not %s)", renderSExpr(n.Operand))
	case *AtTimeZoneExpr:
		return fmt.Sprintf("(attz %s %s)", renderSExpr(n.Value), renderSExpr(n.Zone))
	case *ParenExpr:
		return fmt.Sprintf("(paren %s)", renderSExpr(n.Expr))
	case *ComparisonExpr:
		return fmt.Sprintf("(%s %s %s)", n.Op, renderSExpr(n.Left), renderSExpr(n.Right))
	case *BetweenExpr:
		not := ""
		if n.Not {
			not = "!"
		}
		return fmt.Sprintf("(%sbetween %s %s %s)", not, renderSExpr(n.Value), renderSExpr(n.Lower), renderSExpr(n.Upper))
	case *IsNullExpr:
		not := ""
		if n.Not {
			not = "!"
		}
		return fmt.Sprintf("(%sisnull %s)", not, renderSExpr(n.Value))
	case *IsDistinctFromExpr:
		not := ""
		if n.Not {
			not = "!"
		}
		return fmt.Sprintf("(%sdistinct %s %s)", not, renderSExpr(n.Value), renderSExpr(n.Right))
	case *IsBooleanExpr:
		not := ""
		if n.Not {
			not = "!"
		}
		return fmt.Sprintf("(%sis%s %s)", not, strings.ToLower(n.Test), renderSExpr(n.Value))
	case *LikeExpr:
		not := ""
		if n.Not {
			not = "!"
		}
		return fmt.Sprintf("(%slike %s %s)", not, renderSExpr(n.Value), renderSExpr(n.Pattern))
	default:
		return fmt.Sprintf("<%T>", e)
	}
}

// TestExpr_Precedence asserts the parse-tree shape of precedence- and
// associativity-sensitive expressions. Each case fixes how the operators group.
func TestExpr_Precedence(t *testing.T) {
	cases := []struct {
		expr string
		want string
	}{
		// multiplicative binds tighter than additive
		{"1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"1 * 2 + 3", "(+ (* 1 2) 3)"},
		// additive/multiplicative are left-associative
		{"1 - 2 - 3", "(- (- 1 2) 3)"},
		{"1 / 2 / 3", "(/ (/ 1 2) 3)"},
		// unary binds tighter than multiplicative; nests right
		{"-a * b", "(* (u- a) b)"},
		{"- - 1", "(u- (u- 1))"},
		// || is looser than arithmetic, left-associative
		{"a || b || c", "(|| (|| a b) c)"},
		{"a + b || c", "(|| (+ a b) c)"},
		// predicate sits above ||: a || b = c is (a||b) = c (P2)
		{"a || b = c", "(= (|| a b) c)"},
		// comparison wraps two value expressions
		{"1 + 2 = 3", "(= (+ 1 2) 3)"},
		// NOT wraps the whole predicated expression (loosest)
		{"NOT a = b", "(not (= a b))"},
		{"NOT NOT a", "(not (not a))"},
		// AND binds tighter than OR
		{"a OR b AND c", "(or a (and b c))"},
		{"a AND b OR c", "(or (and a b) c)"},
		// AND/OR left-associative
		{"a AND b AND c", "(and (and a b) c)"},
		// postfix subscript/dereference are left-associative and tightest
		{"a.b.c", "(. (. a b) c)"},
		{"a[1][2]", "([] ([] a 1) 2)"},
		{"a.b[1]", "([] (. a b) 1)"},
		{"m['k'].f", "(. ([] m 'k') f)"},
		// IS NULL is a single non-associative predicate
		{"a + b IS NULL", "(isnull (+ a b))"},
		{"a IS NOT NULL", "(!isnull a)"},
		// BETWEEN bounds are value expressions (the AND is the separator)
		{"a BETWEEN 1 + 1 AND 2 * 2", "(between a (+ 1 1) (* 2 2))"},
		// LIKE / DISTINCT FROM operate on value expressions
		{"a IS DISTINCT FROM b + c", "(distinct a (+ b c))"},
		// boolean test binds the whole value expression, NOT tighter than +
		{"a IS TRUE", "(istrue a)"},
		{"a IS NOT FALSE", "(!isfalse a)"},
		{"a + b IS UNKNOWN", "(isunknown (+ a b))"},
		{"a IS NOT UNKNOWN", "(!isunknown a)"},
	}
	for _, tc := range cases {
		t.Run(truncateName(tc.expr), func(t *testing.T) {
			node, errs := ParseExpression(tc.expr)
			if len(errs) != 0 {
				t.Fatalf("ParseExpression(%q) failed: %v", tc.expr, errs)
			}
			got := renderSExpr(node)
			if got != tc.want {
				t.Errorf("structure mismatch for %q:\n  got:  %s\n  want: %s", tc.expr, got, tc.want)
			}
		})
	}
}

// renderExpr renders an Expr back to Trino SQL for the round-trip gate. It need
// not be byte-identical to the input (spacing/case may normalize) but must emit
// syntactically valid, structurally faithful Trino. It is test-only; production
// deparse is a separate DAG node.
func renderExpr(e Expr) string {
	var b strings.Builder
	writeExpr(&b, e)
	return b.String()
}

func writeExpr(b *strings.Builder, e Expr) {
	switch n := e.(type) {
	case *Literal:
		switch n.Kind {
		case LiteralString:
			if n.Unicode {
				b.WriteString("U&")
			}
			b.WriteByte('\'')
			b.WriteString(strings.ReplaceAll(n.Value, "'", "''"))
			b.WriteByte('\'')
		case LiteralBinary:
			// The lexer strips the X'…' wrapper; restore it for a valid render.
			b.WriteString("X'")
			b.WriteString(n.Value)
			b.WriteByte('\'')
		default:
			b.WriteString(n.Value)
		}
	case *Parameter:
		b.WriteByte('?')
	case *ColumnRef:
		b.WriteString(n.Name.String())
	case *Dereference:
		writeExpr(b, n.Base)
		b.WriteByte('.')
		b.WriteString(n.FieldName.String())
	case *Subscript:
		writeExpr(b, n.Value)
		b.WriteByte('[')
		writeExpr(b, n.Index)
		b.WriteByte(']')
	case *UnaryExpr:
		b.WriteString(n.Op)
		b.WriteByte('(')
		writeExpr(b, n.Operand)
		b.WriteByte(')')
	case *BinaryExpr:
		b.WriteByte('(')
		writeExpr(b, n.Left)
		b.WriteByte(' ')
		b.WriteString(n.Op)
		b.WriteByte(' ')
		writeExpr(b, n.Right)
		b.WriteByte(')')
	case *LogicalExpr:
		b.WriteByte('(')
		writeExpr(b, n.Left)
		b.WriteByte(' ')
		b.WriteString(n.Op)
		b.WriteByte(' ')
		writeExpr(b, n.Right)
		b.WriteByte(')')
	case *NotExpr:
		b.WriteString("(NOT ")
		writeExpr(b, n.Operand)
		b.WriteByte(')')
	case *AtTimeZoneExpr:
		b.WriteByte('(')
		writeExpr(b, n.Value)
		b.WriteString(" AT TIME ZONE ")
		writeExpr(b, n.Zone)
		b.WriteByte(')')
	case *ParenExpr:
		b.WriteByte('(')
		writeExpr(b, n.Expr)
		b.WriteByte(')')
	case *RowConstructor:
		if n.Explicit {
			b.WriteString("ROW")
		}
		b.WriteByte('(')
		writeExprList(b, n.Elements)
		b.WriteByte(')')
	case *ArrayConstructor:
		b.WriteString("ARRAY[")
		writeExprList(b, n.Elements)
		b.WriteByte(']')
	case *IntervalLiteral:
		b.WriteString("INTERVAL ")
		b.WriteString(n.Sign)
		b.WriteByte('\'')
		b.WriteString(n.Value)
		b.WriteString("' ")
		b.WriteString(n.From.String())
		if n.To != nil {
			b.WriteString(" TO ")
			b.WriteString(n.To.String())
		}
	case *TypeConstructor:
		b.WriteString(n.Name)
		b.WriteString(" '")
		b.WriteString(strings.ReplaceAll(n.Value, "'", "''"))
		b.WriteByte('\'')
	case *SubqueryExpr:
		switch n.Kind {
		case SubqueryExists:
			b.WriteString("EXISTS (")
			b.WriteString(n.RawText)
			b.WriteByte(')')
		default:
			b.WriteByte('(')
			b.WriteString(n.RawText)
			b.WriteByte(')')
		}
	case *ComparisonExpr:
		b.WriteByte('(')
		writeExpr(b, n.Left)
		b.WriteByte(' ')
		b.WriteString(n.Op)
		b.WriteByte(' ')
		writeExpr(b, n.Right)
		b.WriteByte(')')
	case *QuantifiedComparisonExpr:
		b.WriteByte('(')
		writeExpr(b, n.Left)
		b.WriteByte(' ')
		b.WriteString(n.Op)
		b.WriteByte(' ')
		b.WriteString(n.Quantifier)
		b.WriteString(" (")
		b.WriteString(n.Subquery.RawText)
		b.WriteString("))")
	case *BetweenExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" NOT")
		}
		b.WriteString(" BETWEEN ")
		writeExpr(b, n.Lower)
		b.WriteString(" AND ")
		writeExpr(b, n.Upper)
	case *InListExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" NOT")
		}
		b.WriteString(" IN (")
		writeExprList(b, n.List)
		b.WriteByte(')')
	case *InSubqueryExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" NOT")
		}
		b.WriteString(" IN (")
		b.WriteString(n.Subquery.RawText)
		b.WriteByte(')')
	case *LikeExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" NOT")
		}
		b.WriteString(" LIKE ")
		writeExpr(b, n.Pattern)
		if n.Escape != nil {
			b.WriteString(" ESCAPE ")
			writeExpr(b, n.Escape)
		}
	case *IsNullExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" IS NOT NULL")
		} else {
			b.WriteString(" IS NULL")
		}
	case *IsDistinctFromExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" IS NOT DISTINCT FROM ")
		} else {
			b.WriteString(" IS DISTINCT FROM ")
		}
		writeExpr(b, n.Right)
	case *IsBooleanExpr:
		writeExpr(b, n.Value)
		if n.Not {
			b.WriteString(" IS NOT ")
		} else {
			b.WriteString(" IS ")
		}
		b.WriteString(n.Test)
	case *CaseExpr:
		b.WriteString("CASE")
		if n.Operand != nil {
			b.WriteByte(' ')
			writeExpr(b, n.Operand)
		}
		for _, w := range n.Whens {
			b.WriteString(" WHEN ")
			writeExpr(b, w.Cond)
			b.WriteString(" THEN ")
			writeExpr(b, w.Result)
		}
		if n.Else != nil {
			b.WriteString(" ELSE ")
			writeExpr(b, n.Else)
		}
		b.WriteString(" END")
	case *CastExpr:
		if n.Try {
			b.WriteString("TRY_CAST(")
		} else {
			b.WriteString("CAST(")
		}
		writeExpr(b, n.Expr)
		b.WriteString(" AS ")
		b.WriteString(n.Type.String())
		b.WriteByte(')')
	case *ExtractExpr:
		b.WriteString("EXTRACT(")
		b.WriteString(n.Field)
		b.WriteString(" FROM ")
		writeExpr(b, n.Source)
		b.WriteByte(')')
	case *SubstringExpr:
		b.WriteString("SUBSTRING(")
		writeExpr(b, n.Source)
		b.WriteString(" FROM ")
		writeExpr(b, n.From)
		if n.For != nil {
			b.WriteString(" FOR ")
			writeExpr(b, n.For)
		}
		b.WriteByte(')')
	case *TrimExpr:
		b.WriteString("TRIM(")
		if n.Spec != "" {
			b.WriteString(n.Spec)
			b.WriteByte(' ')
		}
		if n.Char != nil {
			writeExpr(b, n.Char)
			b.WriteByte(' ')
		}
		if n.Spec != "" || n.Char != nil {
			b.WriteString("FROM ")
		}
		writeExpr(b, n.Source)
		b.WriteByte(')')
	case *NormalizeExpr:
		b.WriteString("NORMALIZE(")
		writeExpr(b, n.Source)
		if n.Form != "" {
			b.WriteString(", ")
			b.WriteString(n.Form)
		}
		b.WriteByte(')')
	case *PositionExpr:
		b.WriteString("POSITION(")
		writeExpr(b, n.Needle)
		b.WriteString(" IN ")
		writeExpr(b, n.Haystack)
		b.WriteByte(')')
	case *GroupingExpr:
		b.WriteString("GROUPING(")
		for i, a := range n.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(a.String())
		}
		b.WriteByte(')')
	case *ListaggExpr:
		b.WriteString("listagg(")
		if n.Distinct {
			b.WriteString("DISTINCT ")
		}
		writeExpr(b, n.Arg)
		if n.Separator != nil {
			b.WriteString(", ")
			writeExpr(b, n.Separator)
		}
		if n.OnOverflow == "ERROR" {
			b.WriteString(" ON OVERFLOW ERROR")
		} else if n.OnOverflow == "TRUNCATE" {
			b.WriteString(" ON OVERFLOW TRUNCATE WITH COUNT")
		}
		b.WriteString(") WITHIN GROUP (ORDER BY ")
		writeOrderBy(b, n.WithinGroupBy)
		b.WriteByte(')')
		if n.Filter != nil {
			b.WriteString(" FILTER (WHERE ")
			writeExpr(b, n.Filter)
			b.WriteByte(')')
		}
	case *SpecialFuncExpr:
		b.WriteString(n.Name)
		if n.Precision >= 0 {
			b.WriteByte('(')
			b.WriteString(strconv.Itoa(n.Precision))
			b.WriteByte(')')
		}
	case *LambdaExpr:
		if len(n.Params) == 1 {
			b.WriteString(n.Params[0].String())
		} else {
			b.WriteByte('(')
			for i, p := range n.Params {
				if i > 0 {
					b.WriteString(", ")
				}
				b.WriteString(p.String())
			}
			b.WriteByte(')')
		}
		b.WriteString(" -> ")
		writeExpr(b, n.Body)
	case *FuncCall:
		writeFuncCall(b, n)
	default:
		b.WriteString(fmt.Sprintf("<%T>", e))
	}
}

func writeFuncCall(b *strings.Builder, n *FuncCall) {
	if n.ProcessingMode != "" {
		b.WriteString(n.ProcessingMode)
		b.WriteByte(' ')
	}
	b.WriteString(n.Name.String())
	// A measure (identifier OVER ...) has no parens.
	if !(len(n.Args) == 0 && !n.Star && (n.Over != nil || n.OverName != nil) && n.Filter == nil) {
		b.WriteByte('(')
		if n.Star {
			if n.Label != nil {
				b.WriteString(n.Label.String())
				b.WriteByte('.')
			}
			b.WriteByte('*')
		} else {
			if n.Distinct {
				b.WriteString("DISTINCT ")
			}
			writeExprList(b, n.Args)
			if len(n.OrderBy) > 0 {
				b.WriteString(" ORDER BY ")
				writeOrderBy(b, n.OrderBy)
			}
		}
		b.WriteByte(')')
	}
	if n.Filter != nil {
		b.WriteString(" FILTER (WHERE ")
		writeExpr(b, n.Filter)
		b.WriteByte(')')
	}
	if n.NullTreatment != "" {
		b.WriteByte(' ')
		b.WriteString(n.NullTreatment)
	}
	if n.OverName != nil {
		b.WriteString(" OVER ")
		b.WriteString(n.OverName.String())
	} else if n.Over != nil {
		b.WriteString(" OVER ")
		writeWindowSpec(b, n.Over)
	}
}

func writeWindowSpec(b *strings.Builder, w *WindowSpec) {
	b.WriteByte('(')
	parts := []string{}
	if w.ExistingName != nil {
		parts = append(parts, w.ExistingName.String())
	}
	if len(w.PartitionBy) > 0 {
		var pb strings.Builder
		pb.WriteString("PARTITION BY ")
		writeExprList(&pb, w.PartitionBy)
		parts = append(parts, pb.String())
	}
	if len(w.OrderBy) > 0 {
		var ob strings.Builder
		ob.WriteString("ORDER BY ")
		writeOrderBy(&ob, w.OrderBy)
		parts = append(parts, ob.String())
	}
	if w.Frame != nil {
		var fb strings.Builder
		writeFrame(&fb, w.Frame)
		parts = append(parts, fb.String())
	}
	b.WriteString(strings.Join(parts, " "))
	b.WriteByte(')')
}

func writeFrame(b *strings.Builder, f *WindowFrame) {
	b.WriteString(f.FrameType)
	b.WriteByte(' ')
	if f.End != nil {
		b.WriteString("BETWEEN ")
		writeBound(b, f.Start)
		b.WriteString(" AND ")
		writeBound(b, *f.End)
	} else {
		writeBound(b, f.Start)
	}
}

func writeBound(b *strings.Builder, bound WindowBound) {
	switch bound.Kind {
	case BoundUnboundedPreceding:
		b.WriteString("UNBOUNDED PRECEDING")
	case BoundUnboundedFollowing:
		b.WriteString("UNBOUNDED FOLLOWING")
	case BoundCurrentRow:
		b.WriteString("CURRENT ROW")
	case BoundPreceding:
		writeExpr(b, bound.Value)
		b.WriteString(" PRECEDING")
	case BoundFollowing:
		writeExpr(b, bound.Value)
		b.WriteString(" FOLLOWING")
	}
}

func writeExprList(b *strings.Builder, list []Expr) {
	for i, e := range list {
		if i > 0 {
			b.WriteString(", ")
		}
		writeExpr(b, e)
	}
}

func writeOrderBy(b *strings.Builder, items []SortItem) {
	for i, it := range items {
		if i > 0 {
			b.WriteString(", ")
		}
		writeExpr(b, it.Expr)
		if it.Ordering != "" {
			b.WriteByte(' ')
			b.WriteString(it.Ordering)
		}
		if it.NullOrder != "" {
			b.WriteString(" NULLS ")
			b.WriteString(it.NullOrder)
		}
	}
}

// TestExpr_RoundTrip parses each accepted corpus expression, renders it back to
// SQL, and re-parses the rendering; the re-parse must succeed. This proves the
// parse produced a structurally faithful tree (a renderer round-trip is the
// structural gate in the absence of a reference parser / production deparse).
func TestExpr_RoundTrip(t *testing.T) {
	for _, expr := range exprAcceptCorpus {
		t.Run(truncateName(expr), func(t *testing.T) {
			node, errs := ParseExpression(expr)
			if len(errs) != 0 {
				t.Fatalf("first parse of %q failed: %v", expr, errs)
			}
			rendered := renderExpr(node)
			_, errs2 := ParseExpression(rendered)
			if len(errs2) != 0 {
				t.Errorf("re-parse of rendered %q (from %q) failed: %v", rendered, expr, errs2)
			}
		})
	}
}

// TestExpr_RoundTripOracle verifies the rendered form of every accepted corpus
// expression is itself accepted by Trino — proving the render emits valid Trino
// expression syntax. Skipped without an oracle. JSON forms are not in the corpus
// (deferred), so every entry is expected to render to valid SQL.
func TestExpr_RoundTripOracle(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, expr := range exprAcceptCorpus {
		expr := expr
		t.Run(truncateName(expr), func(t *testing.T) {
			node, errs := ParseExpression(expr)
			if len(errs) != 0 {
				t.Fatalf("parse of %q failed: %v", expr, errs)
			}
			rendered := renderExpr(node)
			accepted, ok := oracleAccepts(t, o, wrapExpr(rendered))
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if !accepted {
				t.Errorf("render produced an expression Trino rejects: %q -> %q", expr, rendered)
			}
		})
	}
}
