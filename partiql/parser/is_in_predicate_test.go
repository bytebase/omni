package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestIsPredicate_Accept covers the `expr IS [NOT] <type>` predicate
// (PartiQLParser.g4:486 — `lhs=exprPredicate IS NOT? type`). The RHS is the
// full `type` production (g4:674-686), which crucially includes NULL and
// MISSING as atomic types — so `x IS NULL` / `x IS MISSING` are parsed as
// `IS <type>`, not as a separate keyword form.
//
// Oracle (antlr_fallback): the executable generated ANTLR parser at
// github.com/bytebase/parser/partiql ACCEPTS every one of these via the
// top-level Script() rule (verified differentially as `SELECT <expr> FROM t;`,
// parsed to EOF with zero lexer/parser syntax errors). omni previously
// REJECTED all of them with "IS predicate requires NULL, MISSING, TRUE, or
// FALSE" — that was a migration bug.
func TestIsPredicate_Accept(t *testing.T) {
	cases := []struct {
		name  string
		input string
		// want is the ast.NodeToString rendering of the parsed expression.
		want string
	}{
		{
			name:  "is_null",
			input: "x IS NULL",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:NULL} Not:false}",
		},
		{
			name:  "is_missing",
			input: "x IS MISSING",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:MISSING} Not:false}",
		},
		{
			name:  "is_not_null",
			input: "x IS NOT NULL",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:NULL} Not:true}",
		},
		{
			name:  "is_not_missing",
			input: "x IS NOT MISSING",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:MISSING} Not:true}",
		},
		{
			name:  "is_int",
			input: "x IS INT",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:INT} Not:false}",
		},
		{
			name:  "is_struct",
			input: "x IS STRUCT",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:STRUCT} Not:false}",
		},
		{
			name:  "is_decimal_args",
			input: "x IS DECIMAL(10, 2)",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:DECIMAL Args:[10,2]} Not:false}",
		},
		{
			name:  "is_not_int",
			input: "a IS NOT INT",
			want:  "IsExpr{Expr:VarRef{Name:a} Type:TypeRef{Name:INT} Not:true}",
		},
		{
			name:  "is_varchar_arg",
			input: "x IS VARCHAR(255)",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:VARCHAR Args:[255]} Not:false}",
		},
		{
			name:  "is_double_precision",
			input: "x IS DOUBLE PRECISION",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:DOUBLE PRECISION} Not:false}",
		},
		{
			name:  "is_character_varying",
			input: "x IS CHARACTER VARYING",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:CHARACTER VARYING} Not:false}",
		},
		{
			name:  "is_time_with_time_zone",
			input: "x IS TIME WITH TIME ZONE",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:TIME WithTimeZone:true} Not:false}",
		},
		{
			name:  "is_bag",
			input: "x IS BAG",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:BAG} Not:false}",
		},
		{
			name:  "is_custom_type_ident",
			input: "x IS foobar",
			want:  "IsExpr{Expr:VarRef{Name:x} Type:TypeRef{Name:foobar} Not:false}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("ParseExpr(%q) unexpected error: %v", tc.input, err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("ParseExpr(%q)\n got: %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestIsPredicate_Reject covers the IS forms the executable ANTLR parser
// REJECTS. TRUE / FALSE are NOT in the `type` production, and the legacy
// PartiQLParser.g4 precedence comment (line 437-438) marks `IS TRUE/FALSE`
// as "Not yet implemented in PartiQL". UNKNOWN is likewise not a type.
//
// Oracle: via Script(), each yields a syntax error of the form
// `mismatched input 'TRUE' expecting { <type first-set> }` — an unambiguous
// reject. omni previously ACCEPTED `IS TRUE`/`IS FALSE` (the removed
// IsTypeTrue/IsTypeFalse enum values); this test pins the corrected reject.
func TestIsPredicate_Reject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{name: "is_true", input: "x IS TRUE", wantErrIn: "expected type"},
		{name: "is_false", input: "x IS FALSE", wantErrIn: "expected type"},
		{name: "is_not_true", input: "x IS NOT TRUE", wantErrIn: "expected type"},
		{name: "is_not_false", input: "x IS NOT FALSE", wantErrIn: "expected type"},
		{name: "is_unknown", input: "x IS UNKNOWN", wantErrIn: "expected type"},
		{name: "is_missing_type", input: "x IS", wantErrIn: "expected type"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) = nil error, want error containing %q", tc.input, tc.wantErrIn)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("ParseExpr(%q) error = %q, want to contain %q", tc.input, err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestInPredicate_Accept covers the non-empty `expr [NOT] IN (...)` list form
// and the bare-expression `IN <mathOp00>` form (e.g. an array literal), both
// confirmed ACCEPT by the ANTLR oracle.
//
// Grammar (g4:487-488):
//
//	lhs NOT? IN PAREN_LEFT expr PAREN_RIGHT   # PredicateIn  (>=1 expr)
//	lhs NOT? IN rhs=mathOp00                   # PredicateIn  (array/subquery)
func TestInPredicate_Accept(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "in_single",
			input: "a IN (1)",
			want:  "InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1}] Not:false}",
		},
		{
			name:  "in_multi",
			input: "a IN (1, 2, 3)",
			want:  "InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}] Not:false}",
		},
		{
			name:  "not_in_multi",
			input: "a NOT IN (1, 2)",
			want:  "InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} NumberLit{Val:2}] Not:true}",
		},
		{
			name:  "in_array_form",
			input: "a IN [1, 2]",
			want:  "InExpr{Expr:VarRef{Name:a} List:[ListLit{Items:[NumberLit{Val:1} NumberLit{Val:2}]}] Not:false}",
		},
		// The parenthesized items are full exprs (g4:487, and the valueList
		// of g4:566-567), not mathOp00s: omni rejected every one of these
		// at the first token a mathOp00 cannot hold.
		{
			name:  "in_subquery",
			input: "a IN (SELECT x FROM t)",
			want:  "InExpr{Expr:VarRef{Name:a} Subquery:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}} Not:false}",
		},
		{
			name:  "not_in_subquery",
			input: "a NOT IN (SELECT x FROM t WHERE y = 1)",
			want:  "InExpr{Expr:VarRef{Name:a} Subquery:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t} Where:BinaryExpr{Op:= Left:VarRef{Name:y} Right:NumberLit{Val:1}}} Not:true}",
		},
		{
			name:  "in_set_operation",
			input: "a IN (SELECT x FROM t UNION SELECT y FROM u)",
			want:  "InExpr{Expr:VarRef{Name:a} Subquery:SetOpStmt{Op:UNION Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:y}}] From:VarRef{Name:u}}} Not:false}",
		},
		{
			name:  "in_set_operation_of_parenthesized_queries",
			input: "a IN ((SELECT x FROM t) UNION (SELECT y FROM u))",
			want:  "InExpr{Expr:VarRef{Name:a} Subquery:SetOpStmt{Op:UNION Left:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}} Right:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:y}}] From:VarRef{Name:u}}} Not:false}",
		},
		// A query in parentheses of its own is a parenthesized expression,
		// one item of the list, as it was before.
		{
			name:  "in_parenthesized_subquery",
			input: "a IN ((SELECT x FROM t))",
			want:  "InExpr{Expr:VarRef{Name:a} List:[SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}}] Not:false}",
		},
		{
			name:  "not_in_parenthesized_subquery",
			input: "a NOT IN ((SELECT x FROM t))",
			want:  "InExpr{Expr:VarRef{Name:a} List:[SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}}] Not:true}",
		},
		{
			name:  "in_mixed_list",
			input: "a IN (1, b + 1, c = 2, d OR e, NOT f, (SELECT x FROM t))",
			want:  "InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} BinaryExpr{Op:+ Left:VarRef{Name:b} Right:NumberLit{Val:1}} BinaryExpr{Op:= Left:VarRef{Name:c} Right:NumberLit{Val:2}} BinaryExpr{Op:OR Left:VarRef{Name:d} Right:VarRef{Name:e}} UnaryExpr{Op:NOT Operand:VarRef{Name:f}} SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}}] Not:false}",
		},
		{
			name:  "in_list_ending_in_query",
			input: "a IN (1, SELECT x FROM t)",
			want:  "InExpr{Expr:VarRef{Name:a} List:[NumberLit{Val:1} SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}}] Not:false}",
		},
		{
			name:  "in_single_predicate",
			input: "a IN (b IS NULL)",
			want:  "InExpr{Expr:VarRef{Name:a} List:[IsExpr{Expr:VarRef{Name:b} Type:TypeRef{Name:NULL} Not:false}] Not:false}",
		},
		// rhs=mathOp00 (g4:488) whose first operand is parenthesized:
		// g4:487 cannot read past its ')', so the rhs is the whole mathOp00,
		// and List holds it alone, as for `a IN [1, 2]`. omni rejected
		// every one of these at the token after the ')'.
		{
			name:  "in_parenthesized_operand_plus",
			input: "a IN (1) + 1",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:+ Left:NumberLit{Val:1} Right:NumberLit{Val:1}}] Not:false}",
		},
		{
			name:  "not_in_parenthesized_operand_concat",
			input: "a NOT IN (b) || 'x'",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:|| Left:VarRef{Name:b} Right:StringLit{Val:\"x\"}}] Not:true}",
		},
		{
			name:  "in_parenthesized_operand_precedence",
			input: "a IN ((1)) * 2 - 3 || 'x'",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:|| Left:BinaryExpr{Op:- Left:BinaryExpr{Op:* Left:NumberLit{Val:1} Right:NumberLit{Val:2}} Right:NumberLit{Val:3}} Right:StringLit{Val:\"x\"}}] Not:false}",
		},
		{
			name:  "in_parenthesized_operand_index",
			input: "a IN (x)[0]",
			want:  "InExpr{Expr:VarRef{Name:a} List:[PathExpr{Root:VarRef{Name:x} Steps:[IndexStep{Index:NumberLit{Val:0}}]}] Not:false}",
		},
		{
			name:  "in_parenthesized_operand_field",
			input: "a IN (x).y + 1",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:+ Left:PathExpr{Root:VarRef{Name:x} Steps:[DotStep{Field:y}]} Right:NumberLit{Val:1}}] Not:false}",
		},
		{
			name:  "in_subquery_operand",
			input: "a IN (SELECT x FROM t) + 1",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:+ Left:SubLink{Stmt:SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:x}}] From:VarRef{Name:t}}} Right:NumberLit{Val:1}}] Not:false}",
		},
		{
			name:  "in_parenthesized_operand_then_comparison",
			input: "a IN (1) - 1 = b",
			want:  "BinaryExpr{Op:= Left:InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:- Left:NumberLit{Val:1} Right:NumberLit{Val:1}}] Not:false} Right:VarRef{Name:b}}",
		},
		{
			name:  "in_nested_parenthesized_operands",
			input: "a IN (b IN (1) + 1) + 1",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:+ Left:InExpr{Expr:VarRef{Name:b} List:[BinaryExpr{Op:+ Left:NumberLit{Val:1} Right:NumberLit{Val:1}}] Not:false} Right:NumberLit{Val:1}}] Not:false}",
		},
		// A graph match in the parentheses is exprGraphMatchMany
		// (g4:625-626), which only g4:488 reads.
		{
			name:  "in_graph_match",
			input: "a IN (g MATCH (n))",
			want:  "InExpr{Expr:VarRef{Name:a} List:[MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:n}}]}]}] Not:false}",
		},
		{
			name:  "in_graph_match_operand",
			input: "a IN (g MATCH (n)) || b",
			want:  "InExpr{Expr:VarRef{Name:a} List:[BinaryExpr{Op:|| Left:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:n}}]}]} Right:VarRef{Name:b}}] Not:false}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("ParseExpr(%q) unexpected error: %v", tc.input, err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("ParseExpr(%q)\n got: %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestInPredicate_Reject covers the empty parenthesized IN list, which the
// ANTLR oracle REJECTS (g4:487 requires `PAREN_LEFT expr PAREN_RIGHT`, i.e.
// at least one expr — `a IN ()` yields `no viable alternative at input
// 'IN ()'`). omni previously ACCEPTED the empty list (the `if p.cur.Type !=
// tokPAREN_RIGHT` guard let it slip through with a nil List); this pins the
// corrected reject. The trailing-comma case is also confirmed reject.
func TestInPredicate_Reject(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		{name: "in_empty", input: "a IN ()", wantErrIn: "IN list requires at least one expression"},
		{name: "not_in_empty", input: "a NOT IN ()", wantErrIn: "IN list requires at least one expression"},
		{name: "in_trailing_comma", input: "a IN (1,)", wantErrIn: "unexpected token"},
		{name: "in_unclosed", input: "a IN (1, 2", wantErrIn: "expected PAREN_RIGHT"},
		{name: "in_unclosed_subquery", input: "a IN (SELECT x FROM t", wantErrIn: "expected PAREN_RIGHT"},
		{name: "in_subquery_missing_from_source", input: "a IN (SELECT x FROM)", wantErrIn: "unexpected token"},
		{name: "in_incomplete_comparison", input: "a IN (b =)", wantErrIn: "unexpected token"},
		{name: "in_parenthesized_operand_incomplete", input: "a IN (1) +", wantErrIn: "unexpected token"},
		{name: "in_parenthesized_operand_unclosed_index", input: "a IN (x)[0", wantErrIn: "expected BRACKET_RIGHT"},
		// A list followed by an operator makes a valueList the operand of
		// the rhs (g4:488, 566-567), which omni has no node for.
		{name: "in_list_operand", input: "a IN (1, 2) + 1", wantErrIn: "unexpected token"},
		{name: "in_graph_match_in_list", input: "a IN (1, g MATCH (n))", wantErrIn: "expected PAREN_RIGHT"},
		{name: "in_graph_match_unclosed", input: "a IN (g MATCH (n), x)", wantErrIn: "expected ')' to close graph MATCH expression"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("ParseExpr(%q) = nil error, want error containing %q", tc.input, tc.wantErrIn)
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("ParseExpr(%q) error = %q, want to contain %q", tc.input, err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestInSubqueryLoc: in a statement, the IN predicate spans its closing
// parenthesis and the subquery spans the query inside the parentheses.
func TestInSubqueryLoc(t *testing.T) {
	const sql = "SELECT * FROM t WHERE a NOT IN (SELECT b FROM u)"
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	where := list.Items[0].(*ast.SelectStmt).Where
	in, ok := where.(*ast.InExpr)
	if !ok {
		t.Fatalf("WHERE is %T, want *ast.InExpr", where)
	}
	if got := sql[in.Loc.Start:in.Loc.End]; got != "a NOT IN (SELECT b FROM u)" {
		t.Errorf("InExpr spans %q", got)
	}
	if in.Subquery == nil {
		t.Fatalf("Subquery is nil, List = %v", in.List)
	}
	if loc := in.Subquery.GetLoc(); sql[loc.Start:loc.End] != "SELECT b FROM u" {
		t.Errorf("Subquery spans %q", sql[loc.Start:loc.End])
	}
}

// TestInParenthesizedOperandLoc: when the rhs of IN is a mathOp00 that
// begins with a parenthesized operand, the IN predicate ends at the rhs's
// last token, and a subquery operand spans its parentheses, as it does
// anywhere else.
func TestInParenthesizedOperandLoc(t *testing.T) {
	for _, tc := range []struct {
		sql, in, operand string
	}{
		{"SELECT * FROM t WHERE a IN (SELECT b FROM u) + 1", "a IN (SELECT b FROM u) + 1", "(SELECT b FROM u)"},
		{"SELECT * FROM t WHERE a IN (b) + (2) AND c", "a IN (b) + (2)", "b"},
		{"SELECT * FROM t WHERE a NOT IN (g MATCH (n)) AND c", "a NOT IN (g MATCH (n))", ""},
		{"SELECT * FROM t WHERE a IN (x)[0] = c", "a IN (x)[0]", "x"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", tc.sql, err)
		}
		var in *ast.InExpr
		switch w := list.Items[0].(*ast.SelectStmt).Where.(type) {
		case *ast.InExpr:
			in = w
		case *ast.BinaryExpr:
			in = w.Left.(*ast.InExpr)
		}
		if got := tc.sql[in.Loc.Start:in.Loc.End]; got != tc.in {
			t.Errorf("%q: InExpr spans %q, want %q", tc.sql, got, tc.in)
		}
		if tc.operand == "" {
			continue
		}
		var operand ast.Node
		switch rhs := in.List[0].(type) {
		case *ast.BinaryExpr:
			operand = rhs.Left
		case *ast.PathExpr:
			operand = rhs.Root
		}
		if loc := operand.GetLoc(); tc.sql[loc.Start:loc.End] != tc.operand {
			t.Errorf("%q: operand spans %q, want %q", tc.sql, tc.sql[loc.Start:loc.End], tc.operand)
		}
	}
}
