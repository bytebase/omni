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
