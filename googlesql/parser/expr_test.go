package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// expr_test.go is the hand-written unit suite for the `expressions` node. It
// complements expr_oracle_test.go (the live Spanner-emulator differential):
//
//   - the differential proves accept/reject parity for SHARED GoogleSQL-core
//     forms against the live oracle (both polarities);
//   - this file proves (a) STRUCTURAL correctness — precedence/associativity
//     grouping, which accept/reject cannot catch (1+2*3 "accepts" however it
//     groups), via a deparse-style sexpr round-trip; (b) coverage of the
//     BigQuery-only forms the Spanner emulator cannot authoritatively verdict
//     (inline WITH(...), REPLACE_FIELDS, NEW proto / braced constructors); and
//     (c) the node-internal invariants (Loc spans, walker traversal, the
//     notify-style accept alternatives).

// ---------------------------------------------------------------------------
// Structural precedence / associativity (the structural gate)
// ---------------------------------------------------------------------------

// sexpr renders an expression tree to a compact prefix form so precedence and
// associativity are asserted exactly, not just "it parsed".
func sexpr(n ast.Node) string {
	switch v := n.(type) {
	case *ast.BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", v.Op, sexpr(v.Left), sexpr(v.Right))
	case *ast.CompareExpr:
		if v.Quantifier != "" {
			var rhs string
			switch {
			case v.QuantUnnest != nil:
				rhs = "unnest"
			case v.QuantSubquery != nil:
				rhs = "subq"
			default:
				rhs = fmt.Sprintf("%d-vals", len(v.QuantValues))
			}
			return fmt.Sprintf("(%s-%s %s %s)", v.Op, v.Quantifier, sexpr(v.Left), rhs)
		}
		return fmt.Sprintf("(%s %s %s)", v.Op, sexpr(v.Left), sexpr(v.Right))
	case *ast.UnaryExpr:
		return fmt.Sprintf("(%s %s)", v.Op, sexpr(v.Expr))
	case *ast.IsExpr:
		if v.DistinctFrom != nil {
			pfx := "ISDISTINCT"
			if v.Not {
				pfx = "ISNOTDISTINCT"
			}
			return fmt.Sprintf("(%s %s %s)", pfx, sexpr(v.Expr), sexpr(v.DistinctFrom))
		}
		pfx := "IS"
		if v.Not {
			pfx = "ISNOT"
		}
		return fmt.Sprintf("(%s %s %s)", pfx, sexpr(v.Expr), v.Pred)
	case *ast.BetweenExpr:
		pfx := "BETWEEN"
		if v.Not {
			pfx = "NOTBETWEEN"
		}
		return fmt.Sprintf("(%s %s %s %s)", pfx, sexpr(v.Expr), sexpr(v.Low), sexpr(v.High))
	case *ast.InExpr:
		pfx := "IN"
		if v.Not {
			pfx = "NOTIN"
		}
		var rhs string
		switch {
		case v.Unnest != nil:
			rhs = "unnest"
		case v.Subquery != nil:
			rhs = "subq"
		default:
			rhs = fmt.Sprintf("%d-vals", len(v.Values))
		}
		return fmt.Sprintf("(%s %s %s)", pfx, sexpr(v.Expr), rhs)
	case *ast.LikeExpr:
		pfx := "LIKE"
		if v.Not {
			pfx = "NOTLIKE"
		}
		if v.Quantifier != "" {
			return fmt.Sprintf("(%s-%s %s)", pfx, v.Quantifier, sexpr(v.Expr))
		}
		return fmt.Sprintf("(%s %s %s)", pfx, sexpr(v.Expr), sexpr(v.Pattern))
	case *ast.Identifier:
		return v.Name
	case *ast.PathExpr:
		return "path:" + v.String()
	case *ast.Literal:
		return v.Value
	case *ast.TypedLiteral:
		return v.TypeKeyword + ":" + v.Value
	case *ast.Parameter:
		if v.Positional {
			return "?"
		}
		return "@" + v.Name
	case *ast.SystemVariable:
		return "@@" + strings.Join(v.Parts, ".")
	case *ast.ParenExpr:
		return "(paren " + sexpr(v.Expr) + ")"
	case *ast.FieldAccess:
		return fmt.Sprintf("(. %s %s)", sexpr(v.Expr), v.Field)
	case *ast.IndexAccess:
		return fmt.Sprintf("([] %s %s)", sexpr(v.Expr), sexpr(v.Index))
	case *ast.ExtensionAccess:
		return fmt.Sprintf("(.ext %s %s)", sexpr(v.Expr), v.Path.String())
	case *ast.FuncCall:
		parts := []string{"call:" + v.Name.String()}
		if v.Distinct {
			parts = append(parts, "DISTINCT")
		}
		for _, a := range v.Args {
			parts = append(parts, sexpr(a))
		}
		if v.NullHandling != "" {
			parts = append(parts, v.NullHandling)
		}
		if v.Over != nil {
			parts = append(parts, "OVER")
		}
		return "(" + strings.Join(parts, " ") + ")"
	case *ast.StarExpr:
		return "*"
	case *ast.CaseExpr:
		pfx := "case"
		if v.Operand != nil {
			pfx = "case:" + sexpr(v.Operand)
		}
		s := pfx
		for _, w := range v.Whens {
			s += fmt.Sprintf(" (when %s %s)", sexpr(w.Cond), sexpr(w.Result))
		}
		if v.Else != nil {
			s += " (else " + sexpr(v.Else) + ")"
		}
		return "(" + s + ")"
	case *ast.CastExpr:
		k := "cast"
		if v.Safe {
			k = "safecast"
		}
		return fmt.Sprintf("(%s %s %s)", k, sexpr(v.Expr), v.Type.Text)
	case *ast.ExtractExpr:
		return fmt.Sprintf("(extract %s %s)", sexpr(v.Part), sexpr(v.From))
	case *ast.ArrayExpr:
		return fmt.Sprintf("(array %d)", len(v.Elements))
	case *ast.StructExpr:
		return fmt.Sprintf("(struct %d)", len(v.Fields))
	case *ast.IntervalExpr:
		return fmt.Sprintf("(interval %s %s)", sexpr(v.Value), v.From)
	case *ast.LambdaExpr:
		return fmt.Sprintf("(lambda %s %s)", strings.Join(v.Params, ","), sexpr(v.Body))
	case *ast.NamedArg:
		return fmt.Sprintf("(=> %s %s)", v.Name, sexpr(v.Value))
	case *ast.SubqueryExpr:
		return "(subquery " + v.RawText + ")"
	case *ast.ExistsExpr:
		return "(exists " + v.RawText + ")"
	case *ast.ArraySubqueryExpr:
		return "(arraysubq " + v.RawText + ")"
	case *ast.NewConstructor:
		return fmt.Sprintf("(new %s %d)", v.Type.Text, len(v.Args))
	case *ast.WithExpr:
		return fmt.Sprintf("(with %d %s)", len(v.Vars), sexpr(v.Body))
	case *ast.ReplaceFieldsExpr:
		return fmt.Sprintf("(replacefields %s %d)", sexpr(v.Expr), len(v.Items))
	case *ast.BracedConstructor:
		return fmt.Sprintf("(braced %d)", len(v.Fields))
	case nil:
		return "<nil>"
	default:
		return "<" + n.Tag().String() + ">"
	}
}

func mustParse(t *testing.T, in string) ast.Node {
	t.Helper()
	n, errs := ParseExpression(in)
	if len(errs) != 0 {
		t.Fatalf("ParseExpression(%q) unexpected errors: %v", in, errs)
	}
	if n == nil {
		t.Fatalf("ParseExpression(%q) returned nil node", in)
	}
	return n
}

// TestExpr_Precedence asserts the grouping of mixed-operator expressions — the
// structural gate accept/reject cannot cover. Each grouping was cross-checked
// against the GoogleSQL precedence table and the live oracle's accept of the
// equivalent semantically-revealing form (see expr_oracle_test.go probes).
func TestExpr_Precedence(t *testing.T) {
	cases := []struct{ in, want string }{
		// arithmetic
		{"1 + 2 * 3", "(+ 1 (* 2 3))"},
		{"1 * 2 + 3", "(+ (* 1 2) 3)"},
		{"1 - 2 - 3", "(- (- 1 2) 3)"}, // left-assoc
		{"1 / 2 / 3", "(/ (/ 1 2) 3)"},
		{"-a * b", "(* (- a) b)"}, // unary tighter than *
		{"- - 1", "(- (- 1))"},
		{"~ 1 + 2", "(+ (~ 1) 2)"}, // unary tighter than +
		// bitwise / shift precedence ladder: | < ^ < & < shift < +
		{"1 | 2 ^ 3", "(| 1 (^ 2 3))"},
		{"1 ^ 2 & 3", "(^ 1 (& 2 3))"},
		{"1 & 2 << 3", "(& 1 (<< 2 3))"},
		{"1 << 2 + 3", "(<< 1 (+ 2 3))"},
		{"1 & 2 | 3 ^ 4", "(| (& 1 2) (^ 3 4))"},
		// concat precedence (|| below shift, above comparison). The sexpr renderer
		// shows the UNQUOTED literal body (Literal.Value), so 'a' renders as a.
		{"'a' || 'b' || 'c'", "(|| (|| a b) c)"},
		// comparison below bitwise
		{"5 & 3 = 1", "(= (& 5 3) 1)"},
		{"1 = 5 & 3", "(= 1 (& 5 3))"},
		// logical: AND tighter than OR, both left-assoc
		{"a OR b AND c", "(OR a (AND b c))"},
		{"a AND b OR c", "(OR (AND a b) c)"},
		{"a AND b AND c", "(AND (AND a b) c)"},
		{"a OR b OR c", "(OR (OR a b) c)"},
		{"TRUE AND FALSE OR TRUE", "(OR (AND TRUE FALSE) TRUE)"},
		// NOT below comparison (P2), above AND
		{"NOT a = b", "(NOT (= a b))"},
		{"NOT a IS NULL", "(NOT (IS a NULL))"},
		{"NOT a AND b", "(AND (NOT a) b)"},
		{"NOT NOT a", "(NOT (NOT a))"},
		// access tighter than everything
		{"a.b.c", "path:a.b.c"},
		{"f(x).y", "(. (call:f x) y)"},
		{"a[0].b", "(. ([] a 0) b)"},
		{"a.b[0]", "([] path:a.b 0)"},
		// `a.b` as a primary is a single dotted PathExpr (path_expression), so
		// `-a.b` is `-(a.b)`; the unary still binds looser than the path.
		{"-a.b", "(- path:a.b)"},
	}
	for _, c := range cases {
		got := sexpr(mustParse(t, c.in))
		if got != c.want {
			t.Errorf("%q: got %s, want %s", c.in, got, c.want)
		}
	}
}

// TestExpr_NonAssociative asserts the comparison family is non-associative (P1):
// chaining two comparison-family operators at the same level is a syntax error.
// (Cross-checked against the live oracle, which rejects each of these.)
func TestExpr_NonAssociative(t *testing.T) {
	rejects := []string{
		"a = b = c",
		"1 < 2 < 3",
		"a = b IS NULL",
		"a IN (1) IN (2)",
		"1 IS NULL IS NULL",
		"1 LIKE 'a' LIKE 'b'",
		"1 BETWEEN 0 AND 2 BETWEEN 0 AND 1",
		"a < b > c",
		"1 = 2 != 3",
		// A quantified comparison is itself a comparison-family operator: chaining
		// it (or a plain comparison) onto its result is non-associative. (Oracle:
		// "Expression to the left of comparison must be parenthesized".)
		"x = ANY (SELECT v FROM s) = y",
		"x = ANY (1, 2) = z",
	}
	for _, in := range rejects {
		if _, errs := ParseExpression(in); len(errs) == 0 {
			t.Errorf("ParseExpression(%q): expected a syntax error (non-associative), got none", in)
		}
	}
	// But a comparison nested in parens or under AND/OR is fine.
	accepts := []string{
		"(a = b) = c",
		"a = b AND c = d",
		"a < b OR c > d",
		"1 BETWEEN 0 AND 2 AND TRUE",
	}
	for _, in := range accepts {
		if _, errs := ParseExpression(in); len(errs) != 0 {
			t.Errorf("ParseExpression(%q): expected accept, got %v", in, errs)
		}
	}
}

// ---------------------------------------------------------------------------
// Primary forms — structural shape assertions
// ---------------------------------------------------------------------------

func TestExpr_Literals(t *testing.T) {
	cases := []struct {
		in   string
		kind ast.LiteralKind
		val  string
	}{
		{"42", ast.LitInt, "42"},
		{"0x10", ast.LitInt, "0x10"},
		{"3.14", ast.LitFloat, "3.14"},
		{"NULL", ast.LitNull, "NULL"},
		{"TRUE", ast.LitBool, "TRUE"},
		{"FALSE", ast.LitBool, "FALSE"},
		{"'str'", ast.LitString, "str"},
		{"b'by'", ast.LitBytes, "by"},
	}
	for _, c := range cases {
		n := mustParse(t, c.in)
		lit, ok := n.(*ast.Literal)
		if !ok {
			t.Errorf("%q: got %T, want *ast.Literal", c.in, n)
			continue
		}
		if lit.Kind != c.kind || lit.Value != c.val {
			t.Errorf("%q: got {%s %q}, want {%s %q}", c.in, lit.Kind, lit.Value, c.kind, c.val)
		}
	}
	// Integer value is captured.
	if lit := mustParse(t, "42").(*ast.Literal); lit.Ival != 42 {
		t.Errorf("integer Ival: got %d, want 42", lit.Ival)
	}
	// Adjacent string concatenation merges into one value.
	if lit := mustParse(t, "'ab' 'cd'").(*ast.Literal); lit.Value != "abcd" {
		t.Errorf("string concat: got %q, want %q", lit.Value, "abcd")
	}
}

func TestExpr_TypedLiterals(t *testing.T) {
	cases := []struct{ in, kw, val string }{
		{"DATE '2020-01-01'", "DATE", "2020-01-01"},
		{"TIMESTAMP '2020-01-01'", "TIMESTAMP", "2020-01-01"},
		{"NUMERIC '1.5'", "NUMERIC", "1.5"},
		{"BIGNUMERIC '9'", "BIGNUMERIC", "9"},
		{"JSON '{}'", "JSON", "{}"},
		{"DATETIME '2020-01-01'", "DATETIME", "2020-01-01"},
		{"TIME '00:00'", "TIME", "00:00"},
	}
	for _, c := range cases {
		n := mustParse(t, c.in)
		tl, ok := n.(*ast.TypedLiteral)
		if !ok {
			t.Errorf("%q: got %T, want *ast.TypedLiteral", c.in, n)
			continue
		}
		if tl.TypeKeyword != c.kw || tl.Value != c.val {
			t.Errorf("%q: got {%s %q}, want {%s %q}", c.in, tl.TypeKeyword, tl.Value, c.kw, c.val)
		}
	}
	// RANGE<T> '…' renders the element type into the keyword.
	rl := mustParse(t, "RANGE<DATE> '[2020-01-01, 2020-12-31)'").(*ast.TypedLiteral)
	if !strings.HasPrefix(rl.TypeKeyword, "RANGE<") {
		t.Errorf("range literal keyword = %q, want RANGE<…>", rl.TypeKeyword)
	}
	// A bare type keyword NOT followed by a string is NOT a typed literal.
	if _, ok := mustParse(t, "DATE").(*ast.Identifier); !ok {
		t.Errorf("bare DATE: want *ast.Identifier")
	}
	if fc, ok := mustParse(t, "DATE(x)").(*ast.FuncCall); !ok || fc.Name.String() != "DATE" {
		t.Errorf("DATE(x): want a FuncCall named DATE")
	}
}

func TestExpr_Parameters(t *testing.T) {
	if p := mustParse(t, "@foo").(*ast.Parameter); p.Name != "foo" || p.Positional {
		t.Errorf("@foo = %+v", p)
	}
	if p := mustParse(t, "?").(*ast.Parameter); !p.Positional {
		t.Errorf("? should be positional")
	}
	if sv := mustParse(t, "@@a.b.c").(*ast.SystemVariable); strings.Join(sv.Parts, ".") != "a.b.c" {
		t.Errorf("@@a.b.c parts = %v", sv.Parts)
	}
}

func TestExpr_Paths(t *testing.T) {
	// Single identifier → Identifier; dotted → PathExpr.
	if id, ok := mustParse(t, "col").(*ast.Identifier); !ok || id.Name != "col" {
		t.Errorf("col: want Identifier{col}")
	}
	if pe, ok := mustParse(t, "a.b.c").(*ast.PathExpr); !ok || pe.String() != "a.b.c" {
		t.Errorf("a.b.c: want PathExpr a.b.c")
	}
	// Backtick identifiers are unquoted by the lexer.
	if id, ok := mustParse(t, "`weird name`").(*ast.Identifier); !ok || id.Name != "weird name" {
		t.Errorf("backtick ident: got %#v", mustParse(t, "`weird name`"))
	}
	// Non-reserved keyword as a bare identifier.
	if _, ok := mustParse(t, "data").(*ast.Identifier); !ok {
		t.Errorf("non-reserved keyword `data` should parse as identifier")
	}
}

// TestExpr_FunctionCalls covers the function-call suffix grammar: positional /
// named / lambda / sequence args, DISTINCT, *, null-handling, ORDER BY, LIMIT,
// and keyword function names.
func TestExpr_FunctionCalls(t *testing.T) {
	cases := []struct{ in, want string }{
		{"f()", "(call:f)"},
		{"f(1, 2)", "(call:f 1 2)"},
		{"COUNT(*)", "(call:COUNT *)"},
		{"COUNT(DISTINCT x)", "(call:COUNT DISTINCT x)"},
		{"f(a => 1, b => 2)", "(call:f (=> a 1) (=> b 2))"},
		{"ARRAY_AGG(x IGNORE NULLS)", "(call:ARRAY_AGG x IGNORE NULLS)"},
		{"ARRAY_AGG(x RESPECT NULLS)", "(call:ARRAY_AGG x RESPECT NULLS)"},
		{"f(x -> x + 1)", "(call:f (lambda x (+ x 1)))"},
		{"f((a, b) -> a)", "(call:f (lambda a,b a))"},
		{"f(() -> 1)", "(call:f (lambda  1))"},
		{"pkg.sub.fn(x)", "(call:pkg.sub.fn x)"},
		{"IF(a, b, c)", "(call:IF a b c)"},
		{"GROUPING(x)", "(call:GROUPING x)"},
		{"LEFT(s, 3)", "(call:LEFT s 3)"},
	}
	for _, c := range cases {
		got := sexpr(mustParse(t, c.in))
		if got != c.want {
			t.Errorf("%q: got %s, want %s", c.in, got, c.want)
		}
	}
	// ORDER BY / LIMIT inside an aggregate.
	fc := mustParse(t, "ARRAY_AGG(x ORDER BY y DESC LIMIT 5)").(*ast.FuncCall)
	if len(fc.OrderBy) != 1 || !fc.OrderBy[0].Desc || fc.Limit == nil {
		t.Errorf("ARRAY_AGG order/limit: orderBy=%v limit=%v", fc.OrderBy, fc.Limit)
	}
}

// TestExpr_Window covers OVER clauses: named window, inline partition/order/frame.
func TestExpr_Window(t *testing.T) {
	// Named window.
	fc := mustParse(t, "SUM(x) OVER w").(*ast.FuncCall)
	if fc.Over == nil || fc.Over.Name != "w" || fc.Over.Inline {
		t.Fatalf("OVER w: %+v", fc.Over)
	}
	// Inline with partition + order + frame.
	fc = mustParse(t, "SUM(x) OVER (PARTITION BY a, b ORDER BY c DESC ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW)").(*ast.FuncCall)
	w := fc.Over
	if w == nil || !w.Inline {
		t.Fatalf("inline window not parsed: %+v", w)
	}
	if len(w.PartitionBy) != 2 {
		t.Errorf("PARTITION BY count = %d, want 2", len(w.PartitionBy))
	}
	if len(w.OrderBy) != 1 || !w.OrderBy[0].Desc {
		t.Errorf("ORDER BY = %v", w.OrderBy)
	}
	if w.Frame == nil || w.Frame.Kind != ast.FrameRows || !w.Frame.Between {
		t.Fatalf("frame = %+v", w.Frame)
	}
	if w.Frame.Start.Kind != ast.BoundUnboundedPreceding || w.Frame.End.Kind != ast.BoundCurrentRow {
		t.Errorf("frame bounds = %+v / %+v", w.Frame.Start, w.Frame.End)
	}
	// Single-bound frame with an offset expr.
	fc = mustParse(t, "AVG(x) OVER (ORDER BY t RANGE 5 PRECEDING)").(*ast.FuncCall)
	if fc.Over.Frame == nil || fc.Over.Frame.Between || fc.Over.Frame.Start.Kind != ast.BoundPreceding || fc.Over.Frame.Start.Offset == nil {
		t.Errorf("single-bound frame = %+v", fc.Over.Frame)
	}
}

func TestExpr_CaseCastExtract(t *testing.T) {
	// Searched CASE.
	ce := mustParse(t, "CASE WHEN a THEN 1 WHEN b THEN 2 ELSE 3 END").(*ast.CaseExpr)
	if ce.Operand != nil || len(ce.Whens) != 2 || ce.Else == nil {
		t.Errorf("searched CASE: operand=%v whens=%d else=%v", ce.Operand, len(ce.Whens), ce.Else)
	}
	// Simple CASE.
	ce = mustParse(t, "CASE x WHEN 1 THEN 'a' END").(*ast.CaseExpr)
	if ce.Operand == nil || len(ce.Whens) != 1 || ce.Else != nil {
		t.Errorf("simple CASE: operand=%v whens=%d else=%v", ce.Operand, len(ce.Whens), ce.Else)
	}
	// CAST / SAFE_CAST with the rendered type.
	cast := mustParse(t, "CAST(x AS ARRAY<INT64>)").(*ast.CastExpr)
	if cast.Safe || cast.Type.Text != "ARRAY<INT64>" {
		t.Errorf("CAST type = %q safe=%v", cast.Type.Text, cast.Safe)
	}
	if !mustParse(t, "SAFE_CAST(x AS STRING)").(*ast.CastExpr).Safe {
		t.Errorf("SAFE_CAST: Safe should be true")
	}
	// CAST with FORMAT + AT TIME ZONE.
	cast = mustParse(t, "CAST(x AS STRING FORMAT 'f' AT TIME ZONE 'UTC')").(*ast.CastExpr)
	if cast.Format == nil || cast.TimeZone == nil {
		t.Errorf("CAST format/tz: format=%v tz=%v", cast.Format, cast.TimeZone)
	}
	// EXTRACT with AT TIME ZONE.
	ee := mustParse(t, "EXTRACT(HOUR FROM ts AT TIME ZONE 'UTC')").(*ast.ExtractExpr)
	if ee.Part == nil || ee.From == nil || ee.TimeZone == nil {
		t.Errorf("EXTRACT: %+v", ee)
	}
}

func TestExpr_ArrayStruct(t *testing.T) {
	// Bare / ARRAY / typed array.
	ae := mustParse(t, "[1, 2, 3]").(*ast.ArrayExpr)
	if ae.HasArrayKeyword || ae.ElemType != nil || len(ae.Elements) != 3 {
		t.Errorf("[1,2,3]: %+v", ae)
	}
	ae = mustParse(t, "ARRAY[1, 2]").(*ast.ArrayExpr)
	if !ae.HasArrayKeyword || len(ae.Elements) != 2 {
		t.Errorf("ARRAY[1,2]: %+v", ae)
	}
	ae = mustParse(t, "ARRAY<INT64>[1]").(*ast.ArrayExpr)
	if !ae.HasArrayKeyword || ae.ElemType == nil || ae.ElemType.Text != "ARRAY<INT64>" {
		t.Errorf("ARRAY<INT64>[1]: elemType=%v", ae.ElemType)
	}
	if len(mustParse(t, "[]").(*ast.ArrayExpr).Elements) != 0 {
		t.Errorf("empty array should have 0 elements")
	}
	// STRUCT forms.
	se := mustParse(t, "STRUCT(1 AS x, 2 AS y)").(*ast.StructExpr)
	if !se.HasStruct || se.Type != nil || len(se.Fields) != 2 || se.Fields[0].Alias != "x" {
		t.Errorf("STRUCT(...): %+v", se)
	}
	se = mustParse(t, "STRUCT<x INT64>(1)").(*ast.StructExpr)
	if !se.HasStruct || se.Type == nil || se.Type.Text != "STRUCT<x INT64>" || len(se.Fields) != 1 {
		t.Errorf("STRUCT<x INT64>(1): type=%v fields=%d", se.Type, len(se.Fields))
	}
	// Bare tuple (>= 2 elements) is a struct.
	se = mustParse(t, "(1, 2, 3)").(*ast.StructExpr)
	if se.HasStruct || len(se.Fields) != 3 {
		t.Errorf("(1,2,3): hasStruct=%v fields=%d", se.HasStruct, len(se.Fields))
	}
	// A single parenthesized expr is NOT a struct.
	if _, ok := mustParse(t, "(1)").(*ast.ParenExpr); !ok {
		t.Errorf("(1): want ParenExpr, got %T", mustParse(t, "(1)"))
	}
}

func TestExpr_Access(t *testing.T) {
	// OFFSET/ORDINAL inside a subscript are ordinary function calls.
	ia := mustParse(t, "a[OFFSET(0)]").(*ast.IndexAccess)
	fc, ok := ia.Index.(*ast.FuncCall)
	if !ok || fc.Name.String() != "OFFSET" {
		t.Errorf("a[OFFSET(0)] index = %T", ia.Index)
	}
	// Extension access `. ( path )`.
	ext := mustParse(t, "a.(pkg.field)").(*ast.ExtensionAccess)
	if ext.Path.String() != "pkg.field" {
		t.Errorf("extension path = %q", ext.Path.String())
	}
	// Chained access folds left.
	if got := sexpr(mustParse(t, "a.b[0].c")); got != "(. ([] path:a.b 0) c)" {
		t.Errorf("a.b[0].c = %s", got)
	}
}

func TestExpr_Interval(t *testing.T) {
	ie := mustParse(t, "INTERVAL 5 DAY").(*ast.IntervalExpr)
	if ie.From != "DAY" || ie.To != "" {
		t.Errorf("INTERVAL 5 DAY: from=%q to=%q", ie.From, ie.To)
	}
	ie = mustParse(t, "INTERVAL x YEAR TO MONTH").(*ast.IntervalExpr)
	if ie.From != "YEAR" || ie.To != "MONTH" {
		t.Errorf("INTERVAL TO: from=%q to=%q", ie.From, ie.To)
	}
}

func TestExpr_Subqueries(t *testing.T) {
	// Subquery RawText captures the inner query (Query stays nil for parser-select).
	sq := mustParse(t, "(SELECT 1)").(*ast.SubqueryExpr)
	if sq.Query != nil || sq.RawText != "SELECT 1" {
		t.Errorf("subquery: query=%v raw=%q", sq.Query, sq.RawText)
	}
	ex := mustParse(t, "EXISTS(SELECT x FROM t)").(*ast.ExistsExpr)
	if ex.RawText != "SELECT x FROM t" {
		t.Errorf("exists raw = %q", ex.RawText)
	}
	as := mustParse(t, "ARRAY(SELECT 1)").(*ast.ArraySubqueryExpr)
	if as.RawText != "SELECT 1" {
		t.Errorf("array subquery raw = %q", as.RawText)
	}
	// Balanced nesting inside the subquery body.
	sq = mustParse(t, "(SELECT f(1, (2 + 3)) FROM t WHERE x IN (1, 2))").(*ast.SubqueryExpr)
	if !strings.Contains(sq.RawText, "WHERE x IN (1, 2)") {
		t.Errorf("nested subquery raw = %q", sq.RawText)
	}
	// A subquery participates in an outer expression.
	if got := sexpr(mustParse(t, "(SELECT 1) + 2")); got != "(+ (subquery SELECT 1) 2)" {
		t.Errorf("(SELECT 1)+2 = %s", got)
	}
}

// TestExpr_QuantifiedComparison covers the quantified comparison predicate
// `expr {= != <> < <= > >=} {ANY|SOME|ALL} <rhs>` (the any_some_all production
// on comparative_operator). This is SHARED GoogleSQL core (the live Spanner
// emulator accepts every form here — see expr_oracle_test.go and divergence
// #201). Before this node, omni rejected it with "syntax error at or near
// ANY/ALL"; it now parses, mirroring the long-standing `LIKE ANY|SOME|ALL` form.
func TestExpr_QuantifiedComparison(t *testing.T) {
	// Structural shape: operator, quantifier and RHS kind. RHS kind is "subq"
	// (parenthesized query), "N-vals" (parenthesized list), or "unnest".
	cases := []struct{ in, want string }{
		// Subquery RHS, all six operators × the three quantifiers.
		{"x = ANY (SELECT v FROM s)", "(=-ANY x subq)"},
		// `!=` and `<>` both canonicalize to CmpNe (rendered "!=").
		{"x != ANY (SELECT v FROM s)", "(!=-ANY x subq)"},
		{"x <> ALL (SELECT v FROM s)", "(!=-ALL x subq)"},
		{"x < SOME (SELECT v FROM s)", "(<-SOME x subq)"},
		{"x <= ALL (SELECT v FROM s)", "(<=-ALL x subq)"},
		{"x > ALL (SELECT v FROM s)", "(>-ALL x subq)"},
		{"x >= ANY (SELECT v FROM s)", "(>=-ANY x subq)"},
		{"x = SOME (SELECT v FROM s)", "(=-SOME x subq)"},
		// List RHS (>= 1 element).
		{"x = ANY (1, 2, 3)", "(=-ANY x 3-vals)"},
		{"x > ALL (1, 2)", "(>-ALL x 2-vals)"},
		{"x = ANY (5)", "(=-ANY x 1-vals)"},
		// UNNEST RHS.
		{"x = ANY UNNEST([1, 2, 3])", "(=-ANY x unnest)"},
		{"x > ALL UNNEST([1, 2])", "(>-ALL x unnest)"},
		// LHS may be a higher-precedence expression (arithmetic binds tighter).
		{"x + 1 > ALL (SELECT v FROM s)", "(>-ALL (+ x 1) subq)"},
		// Optional @{...} hint before the RHS (matches the IN/LIKE hint slot).
		{"x = ANY @{a=1} (SELECT v FROM s)", "(=-ANY x subq)"},
		// Quantified comparison feeds into AND (it is below AND in precedence).
		{"x = ALL (SELECT v FROM s) AND y > 0", "(AND (=-ALL x subq) (> y 0))"},
	}
	for _, c := range cases {
		got := sexpr(mustParse(t, c.in))
		if got != c.want {
			t.Errorf("%q: got %s, want %s", c.in, got, c.want)
		}
	}

	// Field-level assertions for one representative of each RHS form.
	sub := mustParse(t, "x = ANY (SELECT v FROM s)").(*ast.CompareExpr)
	if sub.Op != ast.CmpEq || sub.Quantifier != "ANY" || sub.Right != nil {
		t.Errorf("subquery form: op=%v quant=%q right=%v", sub.Op, sub.Quantifier, sub.Right)
	}
	if sq, ok := sub.QuantSubquery.(*ast.SubqueryExpr); !ok || sq.RawText != "SELECT v FROM s" {
		t.Errorf("subquery form: QuantSubquery=%v", sub.QuantSubquery)
	}
	lst := mustParse(t, "x > ALL (1, 2, 3)").(*ast.CompareExpr)
	if lst.Quantifier != "ALL" || len(lst.QuantValues) != 3 || lst.QuantSubquery != nil || lst.QuantUnnest != nil {
		t.Errorf("list form: quant=%q vals=%d", lst.Quantifier, len(lst.QuantValues))
	}
	un := mustParse(t, "x = ANY UNNEST([1, 2])").(*ast.CompareExpr)
	if un.QuantUnnest == nil || un.QuantSubquery != nil || un.QuantValues != nil {
		t.Errorf("unnest form: unnest=%v", un.QuantUnnest)
	}

	// Inspect visits the quantified RHS children (genwalker coverage).
	var sawSubquery bool
	ast.Inspect(sub, func(n ast.Node) bool {
		if _, ok := n.(*ast.SubqueryExpr); ok {
			sawSubquery = true
		}
		return true
	})
	if !sawSubquery {
		t.Errorf("Inspect did not visit the quantified subquery RHS")
	}

	// Rejects: empty list, double quantifier, missing/unparenthesized RHS, and
	// IS does not take a quantifier. (All confirmed reject by the live oracle.)
	rejects := []string{
		"x = ANY ()",                    // empty list
		"x = ANY ANY (SELECT v FROM s)", // double quantifier
		"x = ANY",                       // no RHS
		"x = ANY SELECT v FROM s",       // subquery not parenthesized
		"x IS ANY (SELECT v FROM s)",    // IS has no quantified form
	}
	for _, in := range rejects {
		if _, errs := ParseExpression(in); len(errs) == 0 {
			t.Errorf("ParseExpression(%q): expected a syntax error, got none", in)
		}
	}
}

// TestExpr_BigQueryOnly covers forms the Spanner emulator cannot authoritatively
// verdict (BigQuery-only). They are validated structurally here and triangulated
// against the legacy GoogleSQLParser.g4 (which parses each) + the BigQuery docs.
// (oracle.md routing: BigQuery-only ⇒ triangulation, emulator verdict ignored.)
func TestExpr_BigQueryOnly(t *testing.T) {
	// Inline WITH(...) expression (with_expression).
	we := mustParse(t, "WITH(a AS 1, b AS a + 1, a + b)").(*ast.WithExpr)
	if len(we.Vars) != 2 || we.Vars[0].Alias != "a" || we.Body == nil {
		t.Errorf("WITH expr: vars=%d body=%v", len(we.Vars), we.Body)
	}
	// NEW proto constructor (new_constructor). Use a non-keyword type name so
	// source case is preserved (keyword-named path parts fold to upper case — a
	// pre-existing identifierText behavior; see TestExpr_KeywordPathCaseFolding).
	nc := mustParse(t, "NEW pkg.MyMessage(1 AS x, 2 AS y)").(*ast.NewConstructor)
	if nc.Type.Text != "pkg.MyMessage" || len(nc.Args) != 2 {
		t.Errorf("NEW: type=%q args=%d", nc.Type.Text, len(nc.Args))
	}
	// NEW with an extension-path alias.
	nc = mustParse(t, "NEW pkg.MyMessage(v AS (ext.field))").(*ast.NewConstructor)
	if len(nc.Args) != 1 || nc.Args[0].Alias != "(ext.field)" {
		t.Errorf("NEW ext: %+v", nc.Args)
	}
	// Braced proto constructor (braced_constructor / braced_new_constructor).
	if _, ok := mustParse(t, "NEW pkg.Type {a: 1, b: 2}").(*ast.NewConstructor); !ok {
		t.Errorf("NEW braced: want NewConstructor")
	}
	bc := mustParse(t, "STRUCT<x INT64> {x: 1}").(*ast.BracedConstructor)
	if bc.Type == nil || bc.Type.Text != "STRUCT<x INT64>" {
		t.Errorf("struct braced: type=%v", bc.Type)
	}
	// REPLACE_FIELDS (replace_fields_expression).
	rf := mustParse(t, "REPLACE_FIELDS(s, 1 AS a, 2 AS b.c)").(*ast.ReplaceFieldsExpr)
	if len(rf.Items) != 2 || rf.Items[1].Alias != "b.c" {
		t.Errorf("REPLACE_FIELDS: items=%d", len(rf.Items))
	}
}

// TestExpr_KeywordPathCaseFolding locks in that the shared identifierText helper
// (datatypes.go) PRESERVES the source casing of a dotted path/type component
// that happens to be a GoogleSQL word-keyword (e.g. Type, Path, Value). It used
// to fold such a part to the canonical UPPER-case keyword name (`pkg.Type` →
// `pkg.TYPE`) on the false premise that keyword tokens carry no source spelling;
// in fact the lexer records the verbatim text in Token.Str (lexer.go
// scanIdentOrKeyword), so identifierText now returns it. This matters because
// keyword-named proto/enum path parts ARE case-sensitive for GoogleSQL
// resolution — folding them silently corrupted the name (and PathExpr shares the
// same helper, so the fix lands here too). Accept/reject parity is unaffected
// (the oracle accepts the form regardless of case — see expr_oracle_test.go).
func TestExpr_KeywordPathCaseFolding(t *testing.T) {
	// Non-keyword parts keep their source case (unchanged).
	if pe := mustParse(t, "Foo.Bar.Baz").(*ast.PathExpr); pe.String() != "Foo.Bar.Baz" {
		t.Errorf("non-keyword path case: got %q, want Foo.Bar.Baz", pe.String())
	}
	// A keyword-named part now ALSO keeps its source case (the fix).
	if pe := mustParse(t, "pkg.Type").(*ast.PathExpr); pe.String() != "pkg.Type" {
		t.Errorf("keyword path-part case: got %q, want pkg.Type (source casing must be preserved)", pe.String())
	}
}

// ---------------------------------------------------------------------------
// Node-internal invariants: Loc spans, walker, notify-style accepts
// ---------------------------------------------------------------------------

// TestExpr_Loc spot-checks that node Loc spans cover the full source extent of
// the parsed construct.
func TestExpr_Loc(t *testing.T) {
	cases := []string{
		"1 + 2", "f(a, b, c)", "CASE WHEN a THEN 1 END", "CAST(x AS INT64)",
		"a[0]", "STRUCT(1 AS x)", "[1, 2, 3]", "a BETWEEN 1 AND 2",
		"SUM(x) OVER (PARTITION BY a)",
	}
	for _, in := range cases {
		n := mustParse(t, in)
		loc := nodeLoc(n)
		if loc.Start != 0 || loc.End != len(in) {
			t.Errorf("%q: Loc = {%d,%d}, want {0,%d}", in, loc.Start, loc.End, len(in))
		}
	}
}

// TestExpr_Walk verifies the generated walker descends into expression children:
// walking `f(a + b, c)` visits both operands of the inner BinaryExpr.
func TestExpr_Walk(t *testing.T) {
	n := mustParse(t, "f(a + b, c)")
	var idents []string
	ast.Inspect(n, func(node ast.Node) bool {
		if id, ok := node.(*ast.Identifier); ok {
			idents = append(idents, id.Name)
		}
		return true
	})
	// f is a PathExpr (name), so the identifiers reached are a, b, c.
	want := map[string]bool{"a": true, "b": true, "c": true}
	if len(idents) != 3 {
		t.Fatalf("walk reached identifiers %v, want a,b,c", idents)
	}
	for _, id := range idents {
		if !want[id] {
			t.Errorf("unexpected identifier %q reached", id)
		}
	}
}

// TestExpr_NotifyStyleAccept verifies the grammar's notify-style alternatives
// that parse a tree while flagging a diagnostic still PARSE (the oracle
// classifies them ACCEPT). A bare SELECT as a function argument is the canonical
// case (`func(SELECT 1)`): a diagnostic is recorded but the call is accepted.
func TestExpr_NotifyStyleAccept(t *testing.T) {
	node, errs := ParseExpression("func(SELECT 1)")
	if node == nil {
		t.Fatal("func(SELECT 1): expected a parsed node (notify-style accept)")
	}
	if _, ok := node.(*ast.FuncCall); !ok {
		t.Errorf("func(SELECT 1): want FuncCall, got %T", node)
	}
	// A diagnostic is recorded (the "argument is not a query" notify).
	if len(errs) == 0 {
		t.Errorf("func(SELECT 1): expected a recorded diagnostic")
	} else if !strings.Contains(errs[0].Msg, "expression, not a query") {
		t.Errorf("func(SELECT 1): diagnostic = %q", errs[0].Msg)
	}
}

// TestExpr_TrailingTokens verifies ParseExpression reports trailing tokens after
// a complete expression (a standalone-entry-point concern).
func TestExpr_TrailingTokens(t *testing.T) {
	node, errs := ParseExpression("1 + 2 garbage")
	if node == nil {
		t.Fatal("expected the leading expression to parse")
	}
	if len(errs) == 0 || !strings.Contains(errs[0].Msg, "unexpected token after expression") {
		t.Errorf("trailing tokens: errs = %v", errs)
	}
}

// TestExpr_NoPanic feeds pathological / truncated / degenerate input and
// asserts ParseExpression never panics or hangs — it must always return cleanly
// (with an error). Covers unbalanced brackets, dangling operators, deep nesting,
// and truncated keyword forms. Bounded by the test timeout to catch any
// no-forward-progress loop.
func TestExpr_NoPanic(t *testing.T) {
	inputs := []string{
		"", "(", ")", "((((", "))))", "[[[[", "@@@@", "...", "a.", ".a",
		"f(", "f(((", "CASE", "CAST(", "EXTRACT(", "NOT NOT NOT",
		"- - - - -", "a[", "a[[", "STRUCT<", "ARRAY<", "NEW",
		"1 + + + +", "a => => b", "() -> () -> x", "WITH(", "REPLACE_FIELDS(",
		"a IN UNNEST", "OVER", "f() OVER (", strings.Repeat("(", 200),
		strings.Repeat("a.", 200) + "b", strings.Repeat("NOT ", 100) + "x",
		strings.Repeat("-", 100) + "1", "f(" + strings.Repeat("x,", 100) + "y)",
		"@", "@@", "?", "0x", "''", "``", "INTERVAL", "LIKE", "BETWEEN",
		strings.Repeat("[", 100) + "1" + strings.Repeat("]", 100),
		"(SELECT (SELECT (SELECT 1)))", "EXISTS(", "ARRAY(",
	}
	for _, in := range inputs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PANIC on %q: %v", in, r)
				}
			}()
			_, _ = ParseExpression(in) // must not panic or hang
		}()
	}
}

// TestExpr_Rejects covers syntax rejects beyond the non-associative set (these
// are also confirmed by the live oracle in expr_oracle_test.go).
func TestExpr_Rejects(t *testing.T) {
	rejects := []string{
		"1 +",
		"CASE END",
		"CAST(x INT64)",
		"CAST(x AS)",
		"EXTRACT(x y)",
		"a IN ()",
		"a..b",
		"INTERVAL 5",        // missing datepart
		"NEW",               // missing type
		"REPLACE_FIELDS(x)", // needs >= 1 replacement
		"@",                 // bare @ with no name
		"f(,)",              // empty arg before comma
	}
	for _, in := range rejects {
		if _, errs := ParseExpression(in); len(errs) == 0 {
			t.Errorf("ParseExpression(%q): expected a syntax error, got none", in)
		}
	}
}
