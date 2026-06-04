package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParser_Aggregate is the inline AST-shape test for aggregate function
// calls owned by DAG node 14 (parser-aggregates): COUNT(*),
// COUNT|SUM|AVG|MIN|MAX([DISTINCT|ALL] expr), plus COUNT's functionCall
// fallback (empty / multi-arg).
//
// Every accept-case verdict below was confirmed against the legacy ANTLR
// parser (truth2) via a throwaway differential probe wrapping each form as a
// SELECT projection; see the PR body for the oracle table. Where the legacy
// grammar and the docs (truth1) agree the AST shape is straightforward; the
// COUNT()/COUNT(x,y) fallback rows are the only behaviorally subtle ones and
// are exercised here and in the error test below.
func TestParser_Aggregate(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		// ----------------------------------------------------------------
		// COUNT(*) — aggregate#CountAll. Only COUNT admits the star form.
		// ----------------------------------------------------------------
		{
			name:  "count_star",
			input: "COUNT(*)",
			want:  "FuncCall{Name:COUNT Star:true Args:[]}",
		},
		{
			name:  "count_star_lowercase",
			input: "count(*)",
			want:  "FuncCall{Name:COUNT Star:true Args:[]}",
		},
		// ----------------------------------------------------------------
		// COUNT(expr) and COUNT([DISTINCT|ALL] expr) — aggregate#AggregateBase.
		// ----------------------------------------------------------------
		{
			name:  "count_expr",
			input: "COUNT(x)",
			want:  "FuncCall{Name:COUNT Args:[VarRef{Name:x}]}",
		},
		{
			name:  "count_distinct",
			input: "COUNT(DISTINCT x)",
			want:  "FuncCall{Name:COUNT Quantifier:DISTINCT Args:[VarRef{Name:x}]}",
		},
		{
			name:  "count_all",
			input: "COUNT(ALL x)",
			want:  "FuncCall{Name:COUNT Quantifier:ALL Args:[VarRef{Name:x}]}",
		},
		{
			name:  "count_distinct_path",
			input: "COUNT(DISTINCT t.a[0])",
			want:  "FuncCall{Name:COUNT Quantifier:DISTINCT Args:[PathExpr{Root:VarRef{Name:t} Steps:[DotStep{Field:a} IndexStep{Index:NumberLit{Val:0}}]}]}",
		},
		// ----------------------------------------------------------------
		// COUNT fallback — empty / multi-arg route through functionCall
		// (COUNT is also in the FunctionCallReserved list, g4:613). No
		// quantifier, no star; ordinary FuncCall shape.
		// ----------------------------------------------------------------
		{
			name:  "count_empty",
			input: "COUNT()",
			want:  "FuncCall{Name:COUNT Args:[]}",
		},
		{
			name:  "count_multi_arg",
			input: "COUNT(x, y)",
			want:  "FuncCall{Name:COUNT Args:[VarRef{Name:x} VarRef{Name:y}]}",
		},
		// ----------------------------------------------------------------
		// SUM / AVG / MIN / MAX — aggregate#AggregateBase only. One expr,
		// optional DISTINCT/ALL.
		// ----------------------------------------------------------------
		{
			name:  "sum_expr",
			input: "SUM(x)",
			want:  "FuncCall{Name:SUM Args:[VarRef{Name:x}]}",
		},
		{
			name:  "avg_expr",
			input: "AVG(price)",
			want:  "FuncCall{Name:AVG Args:[VarRef{Name:price}]}",
		},
		{
			name:  "min_expr",
			input: "MIN(x)",
			want:  "FuncCall{Name:MIN Args:[VarRef{Name:x}]}",
		},
		{
			name:  "max_expr",
			input: "MAX(x)",
			want:  "FuncCall{Name:MAX Args:[VarRef{Name:x}]}",
		},
		{
			name:  "sum_distinct",
			input: "SUM(DISTINCT x)",
			want:  "FuncCall{Name:SUM Quantifier:DISTINCT Args:[VarRef{Name:x}]}",
		},
		{
			name:  "avg_all",
			input: "AVG(ALL x)",
			want:  "FuncCall{Name:AVG Quantifier:ALL Args:[VarRef{Name:x}]}",
		},
		{
			name:  "sum_distinct_expr",
			input: "SUM(DISTINCT a + b)",
			want:  "FuncCall{Name:SUM Quantifier:DISTINCT Args:[BinaryExpr{Op:+ Left:VarRef{Name:a} Right:VarRef{Name:b}}]}",
		},
		{
			name:  "sum_case_insensitive",
			input: "Sum(x)",
			want:  "FuncCall{Name:SUM Args:[VarRef{Name:x}]}",
		},
		// ----------------------------------------------------------------
		// Nested aggregate as an argument (grammar allows it; semantics are
		// not a parse concern).
		// ----------------------------------------------------------------
		{
			name:  "nested_aggregate",
			input: "SUM(COUNT(x))",
			want:  "FuncCall{Name:SUM Args:[FuncCall{Name:COUNT Args:[VarRef{Name:x}]}]}",
		},
		// ----------------------------------------------------------------
		// Trailing path steps attach via parsePrimary's pathStep+ wrapping.
		// ----------------------------------------------------------------
		{
			name:  "count_expr_trailing_path",
			input: "COUNT(x).a",
			want:  "PathExpr{Root:FuncCall{Name:COUNT Args:[VarRef{Name:x}]} Steps:[DotStep{Field:a}]}",
		},
		{
			name:  "count_star_trailing_path",
			input: "COUNT(*).a",
			want:  "PathExpr{Root:FuncCall{Name:COUNT Star:true Args:[]} Steps:[DotStep{Field:a}]}",
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

// TestParser_Aggregate_Errors covers the REJECT cases (NEGATIVE tests are
// required by the correctness protocol). Every reject verdict below was
// confirmed against the legacy ANTLR parser (truth2): the star form is
// COUNT-only and cannot carry a quantifier; SUM/AVG/MIN/MAX have no star,
// no empty, and no multi-arg form (they are absent from FunctionCallReserved).
func TestParser_Aggregate_Errors(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantErrIn string
	}{
		// --- star is COUNT-only and never combines with a quantifier ---
		{
			name:      "count_distinct_star",
			input:     "COUNT(DISTINCT *)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "count_all_star",
			input:     "COUNT(ALL *)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "count_star_extra_arg",
			input:     "COUNT(*, x)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "sum_star",
			input:     "SUM(*)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "max_distinct_star",
			input:     "MAX(DISTINCT *)",
			wantErrIn: "unexpected token",
		},
		// --- quantifier requires exactly one expr ---
		{
			name:      "count_distinct_empty",
			input:     "COUNT(DISTINCT)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "count_all_empty",
			input:     "COUNT(ALL)",
			wantErrIn: "unexpected token",
		},
		{
			name:      "count_distinct_multi",
			input:     "COUNT(DISTINCT x, y)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "sum_distinct_empty",
			input:     "SUM(DISTINCT)",
			wantErrIn: "unexpected token",
		},
		// --- SUM/AVG/MIN/MAX have no empty or multi-arg form ---
		{
			name:      "sum_empty",
			input:     "SUM()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "min_empty",
			input:     "MIN()",
			wantErrIn: "unexpected token",
		},
		{
			name:      "sum_multi_arg",
			input:     "SUM(x, y)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "max_multi_arg",
			input:     "MAX(x, y)",
			wantErrIn: "expected PAREN_RIGHT",
		},
		// --- missing parens ---
		{
			name:      "count_no_paren",
			input:     "COUNT",
			wantErrIn: "expected PAREN_LEFT",
		},
		{
			name:      "sum_no_paren",
			input:     "SUM x",
			wantErrIn: "expected PAREN_LEFT",
		},
		{
			name:      "count_star_unclosed",
			input:     "COUNT(*",
			wantErrIn: "expected PAREN_RIGHT",
		},
		{
			name:      "sum_unclosed",
			input:     "SUM(x",
			wantErrIn: "expected PAREN_RIGHT",
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

// TestParser_Aggregate_Goldens iterates every .partiql file under
// testdata/parser-aggregates/ and compares the parser's pretty-printed
// output against the matching .golden file. Mirrors the builtins-typed
// golden harness. Run with:
//
//	go test -update -run TestParser_Aggregate_Goldens ./partiql/parser/...
func TestParser_Aggregate_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-aggregates/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/parser-aggregates/")
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

// TestParser_Aggregate_RoundTrip is the structural-gate substitute (no
// reference parser for PartiQL): parsing then re-parsing the canonical
// NodeToString fingerprint must be stable.
func TestParser_Aggregate_RoundTrip(t *testing.T) {
	inputs := []string{
		"COUNT(*)",
		"COUNT(x)",
		"COUNT(DISTINCT x)",
		"COUNT(ALL x)",
		"COUNT()",
		"COUNT(x, y)",
		"SUM(x)",
		"AVG(price)",
		"MIN(x)",
		"MAX(x)",
		"SUM(DISTINCT a + b)",
		"SUM(COUNT(x))",
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
			if s2 := ast.NodeToString(e2); s1 != s2 {
				t.Errorf("unstable parse\nfirst:  %s\nsecond: %s", s1, s2)
			}
		})
	}
}
