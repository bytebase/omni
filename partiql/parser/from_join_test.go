package parser

import (
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_DanglingJoinKeyword covers bug B5: a join-type modifier
// (LEFT/RIGHT/FULL/INNER/OUTER, with optional trailing OUTER, plus CROSS)
// that is NOT followed by JOIN must be rejected, not silently dropped.
//
// Oracle: the generated ANTLR PartiQL parser (bytebase/parser/partiql)
// REJECTS every one of these inputs ("no viable alternative" / "mismatched
// input ... expecting JOIN" / "extraneous input"). Before the fix,
// tryParseJoinType consumed the modifier token(s) and returned (Invalid,
// false) without restoring the lexer position, so the keyword vanished and
// the statement wrongly parsed (e.g. "FROM t1 LEFT" parsed as "FROM t1").
//
// We assert via ParseStatement (which requires EOF after the statement):
// because the swallowed keyword now stays in the token stream, parsing
// fails — either inside the FROM-clause join loop or at the trailing-token
// (EOF) check. We only assert rejection (the oracle decides accept/reject);
// the exact message is implementation-defined.
func TestParser_DanglingJoinKeyword(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"left", "SELECT * FROM t1 LEFT"},
		{"right", "SELECT * FROM t1 RIGHT"},
		{"inner", "SELECT * FROM t1 INNER"},
		{"outer", "SELECT * FROM t1 OUTER"},
		{"full", "SELECT * FROM t1 FULL"},
		{"cross", "SELECT * FROM t1 CROSS"},
		{"left_outer", "SELECT * FROM t1 LEFT OUTER"},
		{"right_outer", "SELECT * FROM t1 RIGHT OUTER"},
		{"full_outer", "SELECT * FROM t1 FULL OUTER"},
		// The modifier is followed by a *different* clause keyword, not
		// JOIN: the keyword must not be swallowed into the preceding table
		// reference (pre-fix this parsed as "FROM t1 WHERE a = 1").
		{"left_then_where", "SELECT * FROM t1 LEFT WHERE a = 1"},
		{"right_then_where", "SELECT * FROM t1 RIGHT WHERE a = 1"},
		{"left_outer_then_where", "SELECT * FROM t1 LEFT OUTER WHERE a = 1"},
		// Dangling modifier before ORDER BY (another trailing clause).
		{"inner_then_order", "SELECT * FROM t1 INNER ORDER BY a"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err == nil {
				t.Fatalf("expected rejection, but parsed OK: %s", ast.NodeToString(stmt))
			}
		})
	}
}

// TestParser_WrappedTableReference covers bug B6: the parenthesized
// table-reference form `( tableReference )` (grammar tableReference#
// TableWrapped, PartiQLParser.g4:394) must be accepted when the wrapped
// content is a join (or comma cross-join, or an inner-aliased table) that
// the value-expression grammar cannot consume.
//
// Oracle: the generated ANTLR PartiQL parser ACCEPTS each input. Before the
// fix, parseTablePrimary routed every leading '(' through the
// value-expression paren path, which failed at the first join keyword
// ("expected PAREN_RIGHT, got CROSS/INNER/...").
//
// The parentheses around a table reference are purely structural, so the
// resulting AST is the inner table expression itself (no wrapper node) —
// matching how parseJoinRhs already unwraps `( tableReference )`.
func TestParser_WrappedTableReference(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "cross_join",
			input: "SELECT * FROM (t1 CROSS JOIN t2)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
		{
			name:  "inner_join_on",
			input: "SELECT * FROM (t1 INNER JOIN t2 ON a = b)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:INNER Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "bare_join_on",
			input: "SELECT * FROM (t1 JOIN t2 ON a = b)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:INNER Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "left_join_on",
			input: "SELECT * FROM (t1 LEFT JOIN t2 ON a = b)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:LEFT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "left_outer_join_on",
			input: "SELECT * FROM (t1 LEFT OUTER JOIN t2 ON a = b)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:LEFT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			// Double-wrapped join: the outer parens wrap a wrapped join.
			name:  "double_wrapped_join",
			input: "SELECT * FROM ((t1 CROSS JOIN t2))",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
		{
			// Left-associative chain inside the parens.
			name:  "three_way_cross_join",
			input: "SELECT * FROM (t1 CROSS JOIN t2 CROSS JOIN t3)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}} Right:VarRef{Name:t3}}}",
		},
		{
			// Comma cross-join inside the parens (grammar: COMMA rhs).
			name:  "comma_cross_join",
			input: "SELECT * FROM (t1, t2)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
		{
			// Wrapped join as the LEFT operand of an outer join.
			name:  "wrapped_left_then_cross",
			input: "SELECT * FROM (t1 CROSS JOIN t2) CROSS JOIN t3",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}} Right:VarRef{Name:t3}}}",
		},
		{
			// Wrapped table reference with an inner alias `t1 AS a`. The
			// value-expression grammar cannot consume "t1 AS a", so this is
			// the TableWrapped reading.
			name:  "wrapped_inner_alias",
			input: "SELECT * FROM (t1 AS a)",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:VarRef{Name:t1} As:a}}",
		},
		{
			// Nested: inner `(t1)` is a value-expr paren base reference,
			// the outer parens wrap the cross join.
			name:  "nested_paren_left_cross",
			input: "SELECT * FROM ((t1) CROSS JOIN t2)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("expected accept, got error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tc.want {
				t.Errorf("AST mismatch\ninput: %s\ngot:  %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestParser_WrappedTableReferencePreserved pins the cases the B6 fix must
// NOT regress: a leading '(' whose content is a plain value expression or a
// subquery is read as a base table reference (grammar tableReference#
// TableRefBase whose source is `(expr)`), and that reading DOES accept a
// trailing AS/AT/BY alias chain. Oracle: ANTLR accepts all of these.
func TestParser_WrappedTableReferencePreserved(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "paren_single_table",
			input: "SELECT * FROM (t1)",
			want:  "SelectStmt{Star:true From:VarRef{Name:t1}}",
		},
		{
			name:  "paren_single_table_alias",
			input: "SELECT * FROM (t1) AS x",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:VarRef{Name:t1} As:x}}",
		},
		{
			name:  "paren_single_table_bare_alias",
			input: "SELECT * FROM (t1) t2",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:VarRef{Name:t1} As:t2}}",
		},
		{
			name:  "double_paren_single_table",
			input: "SELECT * FROM ((t1))",
			want:  "SelectStmt{Star:true From:VarRef{Name:t1}}",
		},
		{
			name:  "subquery_alias",
			input: "SELECT * FROM (SELECT * FROM t1) AS s",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:SubLink{Stmt:SelectStmt{Star:true From:VarRef{Name:t1}}} As:s}}",
		},
		{
			name:  "subquery_no_alias",
			input: "SELECT * FROM (SELECT * FROM t1)",
			want:  "SelectStmt{Star:true From:SubLink{Stmt:SelectStmt{Star:true From:VarRef{Name:t1}}}}",
		},
		{
			// Wrapped join on the JOIN right-hand side (parseJoinRhs path).
			// (A *subquery* with a trailing alias on the JOIN RHS, e.g.
			// "... CROSS JOIN (SELECT ...) AS s", is a separate pre-existing
			// parseJoinRhs divergence and is intentionally out of scope for
			// this B5/B6 fix, which only touches tryParseJoinType and
			// parseTablePrimary.)
			name:  "join_rhs_wrapped_join",
			input: "SELECT * FROM t3 CROSS JOIN (t1 CROSS JOIN t2)",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t3} Right:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("expected accept, got error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tc.want {
				t.Errorf("AST mismatch\ninput: %s\ngot:  %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestParser_WrappedTableReferenceNoOuterAlias covers the negative half of
// B6: the TableWrapped form `( tableReference )` does NOT accept a trailing
// (outer) alias. Oracle: ANTLR REJECTS each of these ("mismatched input
// 'AS'", "extraneous input 't3'"). The wrapped reading must therefore leave
// the trailing alias token unconsumed so the statement fails at EOF.
func TestParser_WrappedTableReferenceNoOuterAlias(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"join_as_alias", "SELECT * FROM (t1 CROSS JOIN t2) AS x"},
		{"join_bare_alias", "SELECT * FROM (t1 CROSS JOIN t2) t3"},
		{"inner_join_as_then_join", "SELECT * FROM (t1 INNER JOIN t2 ON a=b) AS j CROSS JOIN t3"},
		{"inner_alias_then_outer_alias", "SELECT * FROM (t1 AS a) AS b"},
		// Empty parens are never a valid table reference.
		{"empty_parens", "SELECT * FROM ()"},
		// UNPIVOT is a primary, not a postfix table operator: `t1 UNPIVOT t2`
		// is not a valid wrapped tableReference.
		{"wrapped_unpivot_postfix", "SELECT * FROM (t1 UNPIVOT t2)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err == nil {
				t.Fatalf("expected rejection, but parsed OK: %s", ast.NodeToString(stmt))
			}
		})
	}
}

// TestParser_ValidJoinsStillParse is a positive control: the fixes must not
// regress ordinary (unwrapped) joins. Oracle: ANTLR accepts all.
func TestParser_ValidJoinsStillParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "left_join",
			input: "SELECT * FROM t1 LEFT JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:LEFT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "left_outer_join",
			input: "SELECT * FROM t1 LEFT OUTER JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:LEFT Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "cross_join",
			input: "SELECT * FROM t1 CROSS JOIN t2",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
		{
			name:  "bare_join",
			input: "SELECT * FROM t1 JOIN t2 ON a = b",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:INNER Left:VarRef{Name:t1} Right:VarRef{Name:t2} On:BinaryExpr{Op:= Left:VarRef{Name:a} Right:VarRef{Name:b}}}}",
		},
		{
			name:  "comma_cross_join_unwrapped",
			input: "SELECT * FROM t1, t2",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t1} Right:VarRef{Name:t2}}}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmt, err := p.ParseStatement()
			if err != nil {
				t.Fatalf("expected accept, got error: %v", err)
			}
			got := ast.NodeToString(stmt)
			if got != tc.want {
				t.Errorf("AST mismatch\ninput: %s\ngot:  %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestParser_FromExpressionSource: a FROM source may be any exprSelect
// (tableBaseReference, PartiQLParser.g4:403-404). One that is not a path, a
// name, a subquery or a graph match, such as a function call, a literal or
// an operator, is an ExprSource holding the expression. The parser dropped
// the expression before, leaving an AliasedSource whose Source was nil.
func TestParser_FromExpressionSource(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  string
		// The text the ExprSource spans, if there is one.
		source string
	}{
		{
			input:  "SELECT * FROM f(x) AS y",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:FuncCall{Name:f Args:[VarRef{Name:x}]}} As:y}}",
			source: "f(x)",
		},
		{
			input:  "SELECT * FROM f(x)",
			want:   "SelectStmt{Star:true From:ExprSource{Expr:FuncCall{Name:f Args:[VarRef{Name:x}]}}}",
			source: "f(x)",
		},
		{
			input:  "SELECT * FROM [1, 2] AS y AT p",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:ListLit{Items:[NumberLit{Val:1} NumberLit{Val:2}]}} As:y At:p}}",
			source: "[1, 2]",
		},
		{
			input:  "SELECT * FROM 'abc' x",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:StringLit{Val:\"abc\"}} As:x}}",
			source: "'abc'",
		},
		{
			input:  "SELECT * FROM a || b AS c",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:BinaryExpr{Op:|| Left:VarRef{Name:a} Right:VarRef{Name:b}}} As:c}}",
			source: "a || b",
		},
		{
			// The parentheses add no node, so the source spans the call.
			input:  "SELECT * FROM (f(x)) AS y",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:FuncCall{Name:f Args:[VarRef{Name:x}]}} As:y}}",
			source: "f(x)",
		},
		{
			input:  "SELECT * FROM t, f(t.x) AS y",
			want:   "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:VarRef{Name:t} Right:AliasedSource{Source:ExprSource{Expr:FuncCall{Name:f Args:[PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:x}]}]}} As:y}}}",
			source: "f(t.x)",
		},
		{
			input:  "SELECT * FROM x IN (1) AS y",
			want:   "SelectStmt{Star:true From:AliasedSource{Source:ExprSource{Expr:InExpr{Expr:VarRef{Name:x} List:[NumberLit{Val:1}] Not:false}} As:y}}",
			source: "x IN (1)",
		},
		{
			input:  "FROM f(x) AS y WHERE a = 1 SET b = 2",
			want:   "UpdateStmt{From:true Source:AliasedSource{Source:ExprSource{Expr:FuncCall{Name:f Args:[VarRef{Name:x}]}} As:y} Sets:[SetAssignment{Target:PathExpr{Root:VarRef{Name:b} Steps:[]} Value:NumberLit{Val:2}}] Where:BinaryExpr{Op:= Left:VarRef{Name:a} Right:NumberLit{Val:1}}}",
			source: "f(x)",
		},
		// Sources that are table expressions themselves stay unwrapped.
		{input: "SELECT * FROM t AS y", want: "SelectStmt{Star:true From:AliasedSource{Source:VarRef{Name:t} As:y}}"},
		{input: "SELECT * FROM f(x).a", want: "SelectStmt{Star:true From:PathExpr{Root:FuncCall{Name:f Args:[VarRef{Name:x}]} Steps:[DotStep{Field:a}]}}"},
		{input: "SELECT * FROM (SELECT * FROM t) s", want: "SelectStmt{Star:true From:AliasedSource{Source:SubLink{Stmt:SelectStmt{Star:true From:VarRef{Name:t}}} As:s}}"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			stmt, err := NewParser(tc.input).ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement(%q): %v", tc.input, err)
			}
			if got := ast.NodeToString(stmt); got != tc.want {
				t.Errorf("AST mismatch\ninput: %s\ngot:  %s\nwant: %s", tc.input, got, tc.want)
			}
			var sources []*ast.ExprSource
			collectExprSources(stmt, &sources)
			if tc.source == "" {
				if len(sources) != 0 {
					t.Errorf("got %d ExprSources, want none", len(sources))
				}
				return
			}
			if len(sources) != 1 {
				t.Fatalf("got %d ExprSources, want 1", len(sources))
			}
			es := sources[0]
			if got := tc.input[es.Loc.Start:es.Loc.End]; got != tc.source {
				t.Errorf("ExprSource spans %q, want %q", got, tc.source)
			}
			if es.Expr.GetLoc() != es.Loc {
				t.Errorf("ExprSource Loc %v, its expression's %v", es.Loc, es.Expr.GetLoc())
			}
		})
	}
}

// collectExprSources appends every ExprSource under n.
func collectExprSources(n ast.Node, out *[]*ast.ExprSource) {
	if es, ok := n.(*ast.ExprSource); ok {
		*out = append(*out, es)
	}
	for _, c := range appendChildren(nil, n) {
		collectExprSources(c, out)
	}
}
