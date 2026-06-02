package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_BuiltinsTyped_Goldens iterates every .partiql file under
// testdata/parser-builtins-typed/ and compares the parser's pretty-printed
// output (via ast.NodeToString) against the matching .golden file.
//
// These golden inputs exercise the keyword-bearing typed builtins owned by
// DAG node 15b (parser-builtins-typed): CAST / CAN_CAST / CAN_LOSSLESS_CAST,
// EXTRACT, TRIM, SUBSTRING, DATE_ADD, DATE_DIFF, COALESCE, NULLIF, and the
// LIST / SEXP sequence constructors.
//
// Run with:
//
//	go test -update -run TestParser_BuiltinsTyped_Goldens ./partiql/parser/...
//
// to regenerate goldens after intentional AST shape changes.
func TestParser_BuiltinsTyped_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-builtins-typed/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/parser-builtins-typed/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(string(input))
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			goldenPath := strings.TrimSuffix(inPath, ".partiql") + ".golden"
			if *update {
				if err := os.WriteFile(goldenPath, []byte(got+"\n"), 0o644); err != nil {
					t.Fatal(err)
				}
				return
			}
			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("golden file missing: %s (run with -update to create)", goldenPath)
			}
			if got+"\n" != string(want) {
				t.Errorf("AST mismatch\ngot:\n%s\nwant:\n%s", got, string(want))
			}
		})
	}
}

// TestParser_BuiltinsTyped is the inline AST-shape test for the typed
// builtins. It mirrors the table-driven style of TestParser_FuncCall and
// asserts the exact ast.NodeToString output. These complement the golden
// files by locking the shape inline (easier to review the parse contract
// at a glance) and by covering structural edge cases.
func TestParser_BuiltinsTyped(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		// ----------------------------------------------------------------
		// CAST family (cast / canCast / canLosslessCast grammar rules).
		// ----------------------------------------------------------------
		{
			name:  "cast_int",
			input: "CAST(x AS INT)",
			want:  "CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:INT}}",
		},
		{
			name:  "cast_decimal_args",
			input: "CAST(Price AS DECIMAL(10,2))",
			want:  "CastExpr{Kind:CAST Expr:VarRef{Name:Price} AsType:TypeRef{Name:DECIMAL Args:[10,2]}}",
		},
		{
			name:  "cast_lowercase_keyword",
			input: "cast(x as int)",
			want:  "CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:INT}}",
		},
		{
			name:  "cast_expr_arg",
			input: "CAST(a + b AS FLOAT)",
			want:  "CastExpr{Kind:CAST Expr:BinaryExpr{Op:+ Left:VarRef{Name:a} Right:VarRef{Name:b}} AsType:TypeRef{Name:FLOAT}}",
		},
		{
			name:  "can_cast",
			input: "CAN_CAST(x AS INT)",
			want:  "CastExpr{Kind:CAN_CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:INT}}",
		},
		{
			name:  "can_lossless_cast",
			input: "CAN_LOSSLESS_CAST(x AS DECIMAL(10,2))",
			want:  "CastExpr{Kind:CAN_LOSSLESS_CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:DECIMAL Args:[10,2]}}",
		},
		{
			name:  "cast_custom_type",
			input: `CAST(x AS "MyType")`,
			want:  "CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:MyType}}",
		},
		{
			name:  "cast_double_precision",
			input: "CAST(x AS DOUBLE PRECISION)",
			want:  "CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:DOUBLE PRECISION}}",
		},
		// CAST result can itself carry path steps: CAST(...).field.
		{
			name:  "cast_with_path_step",
			input: "CAST(x AS STRUCT).a",
			want:  "PathExpr{Root:CastExpr{Kind:CAST Expr:VarRef{Name:x} AsType:TypeRef{Name:STRUCT}} Steps:[DotStep{Field:a}]}",
		},

		// ----------------------------------------------------------------
		// EXTRACT (field FROM expr). Field is an IDENTIFIER per the
		// grammar (extract: EXTRACT '(' IDENTIFIER FROM expr ')'), stored
		// as the raw token text.
		// ----------------------------------------------------------------
		{
			name:  "extract_year",
			input: "EXTRACT(YEAR FROM ts)",
			want:  "ExtractExpr{Field:YEAR From:VarRef{Name:ts}}",
		},
		{
			name:  "extract_month_lower",
			input: "EXTRACT(month FROM ts)",
			want:  "ExtractExpr{Field:month From:VarRef{Name:ts}}",
		},
		{
			name:  "extract_timezone_hour",
			input: "EXTRACT(timezone_hour FROM ts)",
			want:  "ExtractExpr{Field:timezone_hour From:VarRef{Name:ts}}",
		},
		{
			name:  "extract_path_from",
			input: "EXTRACT(SECOND FROM t.created)",
			want:  "ExtractExpr{Field:SECOND From:PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:created}]}}",
		},

		// ----------------------------------------------------------------
		// TRIM ([ [LEADING|TRAILING|BOTH] [sub] FROM ] target). The
		// LEADING/TRAILING/BOTH modifier and the removal char only apply
		// when a FROM follows; otherwise a bare LEADING/BOTH is the target.
		// ----------------------------------------------------------------
		{
			name:  "trim_bare",
			input: "TRIM(s)",
			want:  "TrimExpr{From:VarRef{Name:s}}",
		},
		{
			name:  "trim_from_only",
			input: "TRIM(FROM s)",
			want:  "TrimExpr{From:VarRef{Name:s}}",
		},
		{
			name:  "trim_leading",
			input: "TRIM(LEADING FROM s)",
			want:  "TrimExpr{Spec:LEADING From:VarRef{Name:s}}",
		},
		{
			name:  "trim_trailing",
			input: "TRIM(TRAILING FROM s)",
			want:  "TrimExpr{Spec:TRAILING From:VarRef{Name:s}}",
		},
		{
			name:  "trim_both",
			input: "TRIM(BOTH FROM s)",
			want:  "TrimExpr{Spec:BOTH From:VarRef{Name:s}}",
		},
		{
			name:  "trim_lowercase_mod",
			input: "TRIM(both FROM s)",
			want:  "TrimExpr{Spec:BOTH From:VarRef{Name:s}}",
		},
		{
			name:  "trim_sub_only",
			input: "TRIM(' ' FROM s)",
			want:  `TrimExpr{Sub:StringLit{Val:" "} From:VarRef{Name:s}}`,
		},
		{
			name:  "trim_both_sub",
			input: "TRIM(BOTH ' ' FROM s)",
			want:  `TrimExpr{Spec:BOTH Sub:StringLit{Val:" "} From:VarRef{Name:s}}`,
		},
		{
			name:  "trim_leading_sub",
			input: "TRIM(LEADING '0' FROM acct)",
			want:  `TrimExpr{Spec:LEADING Sub:StringLit{Val:"0"} From:VarRef{Name:acct}}`,
		},
		// A bare LEADING/BOTH with no FROM is the target, not a modifier.
		{
			name:  "trim_keyword_as_target",
			input: "TRIM(both)",
			want:  "TrimExpr{From:VarRef{Name:both}}",
		},
		// LEADING followed by an infix operator is a target expression,
		// not a modifier (the modifier only precedes a removal char/FROM).
		{
			name:  "trim_keyword_in_concat",
			input: "TRIM(leading || x)",
			want:  "TrimExpr{From:BinaryExpr{Op:|| Left:VarRef{Name:leading} Right:VarRef{Name:x}}}",
		},
		// A non-LEADING/TRAILING/BOTH identifier before a FROM-less close
		// is just the target.
		{
			name:  "trim_other_ident_target",
			input: "TRIM(col)",
			want:  "TrimExpr{From:VarRef{Name:col}}",
		},

		// ----------------------------------------------------------------
		// SUBSTRING — both the FROM/FOR keyword form and the comma form.
		// ----------------------------------------------------------------
		{
			name:  "substring_from",
			input: "SUBSTRING(s FROM 2)",
			want:  "SubstringExpr{Expr:VarRef{Name:s} From:NumberLit{Val:2}}",
		},
		{
			name:  "substring_from_for",
			input: "SUBSTRING(s FROM 2 FOR 3)",
			want:  "SubstringExpr{Expr:VarRef{Name:s} From:NumberLit{Val:2} For:NumberLit{Val:3}}",
		},
		{
			name:  "substring_comma",
			input: "SUBSTRING(s, 2)",
			want:  "SubstringExpr{Expr:VarRef{Name:s} From:NumberLit{Val:2}}",
		},
		{
			name:  "substring_comma_len",
			input: "SUBSTRING(s, 2, 3)",
			want:  "SubstringExpr{Expr:VarRef{Name:s} From:NumberLit{Val:2} For:NumberLit{Val:3}}",
		},
		{
			name:  "substring_expr_args",
			input: "SUBSTRING(name FROM start FOR len)",
			want:  "SubstringExpr{Expr:VarRef{Name:name} From:VarRef{Name:start} For:VarRef{Name:len}}",
		},

		// ----------------------------------------------------------------
		// DATE_ADD / DATE_DIFF — generic FuncCall, the date-part is an
		// IDENTIFIER argument (parses as a VarRef).
		// ----------------------------------------------------------------
		{
			name:  "date_add_year",
			input: "DATE_ADD(year, 5, ts)",
			want:  "FuncCall{Name:DATE_ADD Args:[VarRef{Name:year} NumberLit{Val:5} VarRef{Name:ts}]}",
		},
		{
			name:  "date_add_lowercase_keyword",
			input: "date_add(day, n, ts)",
			want:  "FuncCall{Name:DATE_ADD Args:[VarRef{Name:day} VarRef{Name:n} VarRef{Name:ts}]}",
		},
		{
			name:  "date_diff_month",
			input: "DATE_DIFF(month, t1, t2)",
			want:  "FuncCall{Name:DATE_DIFF Args:[VarRef{Name:month} VarRef{Name:t1} VarRef{Name:t2}]}",
		},
		// Only the date-part (arg 1) is constrained; args 2 and 3 are
		// general expressions (here an arithmetic expr and a path).
		{
			name:  "date_add_expr_value_args",
			input: "DATE_ADD(day, n + 1, t.created)",
			want:  "FuncCall{Name:DATE_ADD Args:[VarRef{Name:day} BinaryExpr{Op:+ Left:VarRef{Name:n} Right:NumberLit{Val:1}} PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:created}]}]}",
		},

		// ----------------------------------------------------------------
		// COALESCE (1+ args) and NULLIF (exactly 2).
		// ----------------------------------------------------------------
		{
			name:  "coalesce_two",
			input: "COALESCE(a, b)",
			want:  "CoalesceExpr{Args:[VarRef{Name:a} VarRef{Name:b}]}",
		},
		{
			name:  "coalesce_one",
			input: "COALESCE(x)",
			want:  "CoalesceExpr{Args:[VarRef{Name:x}]}",
		},
		{
			name:  "coalesce_many",
			input: "COALESCE(a, b, c, d)",
			want:  "CoalesceExpr{Args:[VarRef{Name:a} VarRef{Name:b} VarRef{Name:c} VarRef{Name:d}]}",
		},
		{
			name:  "nullif",
			input: "NULLIF(a, b)",
			want:  "NullIfExpr{Left:VarRef{Name:a} Right:VarRef{Name:b}}",
		},
		{
			name:  "nullif_string",
			input: "NULLIF(status, 'active')",
			want:  `NullIfExpr{Left:VarRef{Name:status} Right:StringLit{Val:"active"}}`,
		},

		// ----------------------------------------------------------------
		// LIST / SEXP sequence constructors — generic FuncCall.
		// ----------------------------------------------------------------
		{
			name:  "list_three",
			input: "LIST(1, 2, 3)",
			want:  "FuncCall{Name:LIST Args:[NumberLit{Val:1} NumberLit{Val:2} NumberLit{Val:3}]}",
		},
		{
			name:  "list_empty",
			input: "LIST()",
			want:  "FuncCall{Name:LIST Args:[]}",
		},
		{
			name:  "sexp_two",
			input: "SEXP(1, 2)",
			want:  "FuncCall{Name:SEXP Args:[NumberLit{Val:1} NumberLit{Val:2}]}",
		},
		{
			name:  "sexp_empty",
			input: "SEXP()",
			want:  "FuncCall{Name:SEXP Args:[]}",
		},
		{
			name:  "list_lowercase_keyword",
			input: "list('a', x)",
			want:  `FuncCall{Name:LIST Args:[StringLit{Val:"a"} VarRef{Name:x}]}`,
		},

		// ----------------------------------------------------------------
		// Nesting — the typed builtins compose with each other.
		// ----------------------------------------------------------------
		{
			name:  "nested_cast_extract",
			input: "CAST(EXTRACT(YEAR FROM ts) AS INT)",
			want:  "CastExpr{Kind:CAST Expr:ExtractExpr{Field:YEAR From:VarRef{Name:ts}} AsType:TypeRef{Name:INT}}",
		},
		{
			name:  "nested_coalesce_nullif",
			input: "COALESCE(NULLIF(a, b), c)",
			want:  "CoalesceExpr{Args:[NullIfExpr{Left:VarRef{Name:a} Right:VarRef{Name:b}} VarRef{Name:c}]}",
		},
		{
			name:  "trim_inside_upper",
			input: "UPPER(TRIM(s))",
			want:  "FuncCall{Name:UPPER Args:[TrimExpr{From:VarRef{Name:s}}]}",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("AST mismatch\n got: %s\nwant: %s", got, tc.want)
			}
		})
	}
}

// TestParser_BuiltinsTyped_Errors covers the REJECT cases for the typed
// builtins — malformed forms the grammar does not accept. The parser must
// reject each (NEGATIVE tests are required by the correctness protocol).
func TestParser_BuiltinsTyped_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// --- CAST family: require expr AS type ---
		{
			name:      "cast_missing_as",
			input:     "CAST(x INT)",
			wantErrIn: "expected AS",
		},
		{
			name:      "cast_missing_type",
			input:     "CAST(x AS)",
			wantErrIn: "expected type",
		},
		{
			name:      "cast_missing_paren",
			input:     "CAST x AS INT",
			wantErrIn: "expected PAREN_LEFT",
		},
		{
			name:      "cast_unclosed",
			input:     "CAST(x AS INT",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "cast_empty",
			input:     "CAST()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "can_cast_missing_as",
			input:     "CAN_CAST(x INT)",
			wantErrIn: "expected AS",
		},

		// --- EXTRACT: require IDENTIFIER FROM expr ---
		{
			name:      "extract_missing_from",
			input:     "EXTRACT(YEAR ts)",
			wantErrIn: "expected FROM",
		},
		{
			name:      "extract_missing_field",
			input:     "EXTRACT(FROM ts)",
			wantErrIn: "expected identifier",
		},
		{
			name:      "extract_quoted_field",
			input:     `EXTRACT("YEAR" FROM ts)`,
			wantErrIn: "EXTRACT field must be an unquoted identifier",
		},
		{
			name:      "extract_unclosed",
			input:     "EXTRACT(YEAR FROM ts",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "extract_empty",
			input:     "EXTRACT()",
			wantErrIn: "expected identifier",
		},

		// --- TRIM: when a modifier/sub is given, FROM is required ---
		{
			name:      "trim_unclosed",
			input:     "TRIM(s",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "trim_empty",
			input:     "TRIM()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "trim_mod_sub_no_from",
			input:     "TRIM(LEADING '0' acct)",
			wantErrIn: "expected FROM",
		},

		// --- SUBSTRING: cannot mix comma and FROM/FOR ---
		{
			name:      "substring_unclosed",
			input:     "SUBSTRING(s FROM 2",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "substring_empty",
			input:     "SUBSTRING()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "substring_for_without_from",
			input:     "SUBSTRING(s FOR 3)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "substring_comma_then_for",
			input:     "SUBSTRING(s, 2 FOR 3)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "substring_too_many_commas",
			input:     "SUBSTRING(s, 1, 2, 3)",
			wantErrIn: "expected PAREN_RIGHT",
		},

		// --- DATE_ADD / DATE_DIFF: part IDENTIFIER then 2 exprs ---
		// Too few args: after the second arg the third COMMA is required.
		{
			name:      "date_add_too_few_args",
			input:     "DATE_ADD(year, 5)",
			wantErrIn: "expected COMMA",
		},
		// Too many args: after the third arg the closing paren is required.
		{
			name:      "date_diff_too_many_args",
			input:     "DATE_DIFF(year, a, b, c)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "date_add_unclosed",
			input:     "DATE_ADD(year, 5, ts",
			wantErrIn: "expected PAREN_RIGHT",
		},
		// The date-part (dt=IDENTIFIER in the dateFunction rule) must be a
		// bare unquoted identifier — reject a quoted string, an expression,
		// or a path in the first-argument position.
		{
			name:      "date_add_quoted_part",
			input:     `DATE_ADD('year', 5, ts)`,
			wantErrIn: "date part must be an unquoted identifier",
		},
		{
			name:      "date_add_double_quoted_part",
			input:     `DATE_ADD("year", 5, ts)`,
			wantErrIn: "date part must be an unquoted identifier",
		},
		// The date-part is `dt=IDENTIFIER COMMA` — a single bare identifier
		// immediately followed by a comma. An expression (year + 1) or a
		// path (a.b) starts with an identifier but is then followed by an
		// operator/dot instead of the required comma, so it is rejected
		// there rather than being consumed as a compound part.
		{
			name:      "date_add_expr_part",
			input:     "DATE_ADD(year + 1, 5, ts)",
			wantErrIn: "expected COMMA",
		},
		{
			name:      "date_diff_path_part",
			input:     "DATE_DIFF(a.b, x, y)",
			wantErrIn: "expected COMMA",
		},
		{
			name:      "date_add_numeric_part",
			input:     "DATE_ADD(1, 5, ts)",
			wantErrIn: "date part must be an unquoted identifier",
		},
		{
			name:      "date_add_missing_part",
			input:     "DATE_ADD(, 5, ts)",
			wantErrIn: "date part must be an unquoted identifier",
		},

		// --- NULLIF: exactly 2 args ---
		{
			name:      "nullif_one_arg",
			input:     "NULLIF(a)",
			wantErrIn: "expected COMMA",
		},
		{
			name:      "nullif_three_args",
			input:     "NULLIF(a, b, c)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "nullif_empty",
			input:     "NULLIF()",
			wantErrIn: "unexpected token",
		},

		// --- COALESCE: at least 1 arg ---
		{
			name:      "coalesce_empty",
			input:     "COALESCE()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "coalesce_trailing_comma",
			input:     "COALESCE(a,)",
			wantErrIn: "unexpected token",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantErrIn) {
				t.Errorf("error = %q, want to contain %q", err.Error(), tc.wantErrIn)
			}
		})
	}
}

// TestParser_BuiltinsTyped_RoundTrip is the structural gate substitute for
// the typed builtins. There is no reference parser for PartiQL, so per the
// correctness protocol we use AST-shape stability: parsing a form and then
// re-parsing should yield a structurally identical AST (NodeToString is the
// canonical structural fingerprint). This catches accidental Loc-dependent
// or order-dependent parse instability.
func TestParser_BuiltinsTyped_RoundTrip(t *testing.T) {
	inputs := []string{
		"CAST(x AS INT)",
		"CAN_CAST(x AS DECIMAL(10,2))",
		"CAN_LOSSLESS_CAST(x AS FLOAT)",
		"EXTRACT(YEAR FROM ts)",
		"TRIM(BOTH ' ' FROM s)",
		"TRIM(s)",
		"TRIM(LEADING '0' FROM acct)",
		"SUBSTRING(s FROM 2 FOR 3)",
		"SUBSTRING(s, 2, 3)",
		"DATE_ADD(year, 5, ts)",
		"DATE_DIFF(month, t1, t2)",
		"COALESCE(a, b, c)",
		"NULLIF(a, b)",
		"LIST(1, 2, 3)",
		"SEXP()",
		"CAST(EXTRACT(YEAR FROM ts) AS INT)",
	}
	for _, in := range inputs {
		t.Run(in, func(t *testing.T) {
			p1 := NewParser(in)
			e1, err := p1.ParseExpr()
			if err != nil {
				t.Fatalf("first parse error: %v", err)
			}
			s1 := ast.NodeToString(e1)

			p2 := NewParser(in)
			e2, err := p2.ParseExpr()
			if err != nil {
				t.Fatalf("second parse error: %v", err)
			}
			s2 := ast.NodeToString(e2)

			if s1 != s2 {
				t.Errorf("re-parse instability\nfirst:  %s\nsecond: %s", s1, s2)
			}
		})
	}
}
