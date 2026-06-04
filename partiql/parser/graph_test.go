package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/partiql/ast"
)

// TestParseGraphMatch_Accept drives the GPML MATCH expression parser with
// inputs derived from PartiQLParser.g4 (rules gpmlPattern, matchPattern,
// graphPart, node, edge, edgeWSpec, edgeAbbrev, patternQuantifier,
// matchSelector, patternRestrictor, patternPathVariable) and the official
// PartiQL graph-query docs, asserting the pretty-printed AST shape.
//
// The MATCH form is context-sensitive INSIDE parens: (lhs MATCH pattern).
func TestParseGraphMatch_Accept(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		// ---- nodes -------------------------------------------------------
		{
			name:  "single_node_var",
			input: "(g MATCH (a))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]}",
		},
		{
			name:  "empty_node",
			input: "(g MATCH ())",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{}]}]}",
		},
		{
			name:  "node_label_only",
			input: "(g MATCH (:Account))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Labels:[Account]}]}]}",
		},
		{
			name:  "node_var_label",
			input: "(g MATCH (a:Account))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a} Labels:[Account]}]}]}",
		},
		{
			name:  "node_var_where",
			input: "(g MATCH (a WHERE a.owner = 'Bob'))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a} Where:BinaryExpr{Op:= Left:PathExpr{Root:VarRef{Name:a} Steps:[DotStep{Field:owner}]} Right:StringLit{Val:\"Bob\"}}}]}]}",
		},
		{
			name:  "node_var_label_where",
			input: "(g MATCH (a:Account WHERE a.balance > 100))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a} Labels:[Account] Where:BinaryExpr{Op:> Left:PathExpr{Root:VarRef{Name:a} Steps:[DotStep{Field:balance}]} Right:NumberLit{Val:100}}}]}]}",
		},
		// ---- edges: edgeWSpec (with bracket spec) ------------------------
		{
			name:  "edge_right_spec",
			input: "(g MATCH (a)-[e:Knows]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Labels:[Knows]} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_left_spec",
			input: "(g MATCH (a)<-[e]-(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<- Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_undirected_tilde_spec",
			input: "(g MATCH (a)~[e]~(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:~ Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_bidirectional_spec",
			input: "(g MATCH (a)<-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_undirected_right_spec",
			input: "(g MATCH (a)~[e]~>(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:~> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_undirected_left_spec",
			input: "(g MATCH (a)<~[e]~(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<~ Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_undirected_bidirectional_spec",
			input: "(g MATCH (a)-[e]-(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:- Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_spec_where",
			input: "(g MATCH (a)-[e:Knows WHERE e.since > 2000]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Labels:[Knows] Where:BinaryExpr{Op:> Left:PathExpr{Root:VarRef{Name:e} Steps:[DotStep{Field:since}]} Right:NumberLit{Val:2000}}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- edges: edgeAbbrev (no bracket spec) -------------------------
		{
			name:  "edge_abbrev_right",
			input: "(g MATCH (a)->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:->} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_left",
			input: "(g MATCH (a)<-(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<-} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_undirected_bidirectional",
			input: "(g MATCH (a)-(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_bidirectional",
			input: "(g MATCH (a)<->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<->} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_tilde",
			input: "(g MATCH (a)~(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:~} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_tilde_right",
			input: "(g MATCH (a)~>(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:~>} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "edge_abbrev_left_tilde",
			input: "(g MATCH (a)<~(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:<~} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- quantifiers -------------------------------------------------
		{
			name:  "quantifier_plus",
			input: "(g MATCH (a)-[e]->+(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Quantifier:PatternQuantifier{Min:1 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "quantifier_star",
			input: "(g MATCH (a)-[e]->*(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Quantifier:PatternQuantifier{Min:0 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "quantifier_range",
			input: "(g MATCH (a)-[e]->{1,3}(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Quantifier:PatternQuantifier{Min:1 Max:3}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "quantifier_range_open",
			input: "(g MATCH (a)-[e]->{2,}(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e} Quantifier:PatternQuantifier{Min:2 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "quantifier_abbrev_edge",
			input: "(g MATCH (a)->*(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Quantifier:PatternQuantifier{Min:0 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- selectors ---------------------------------------------------
		{
			name:  "selector_any",
			input: "(g MATCH ANY (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ANY} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "selector_any_k",
			input: "(g MATCH ANY 3 (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ANY} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "selector_all_shortest",
			input: "(g MATCH ALL SHORTEST (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ALL_SHORTEST} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "selector_any_shortest",
			input: "(g MATCH ANY SHORTEST (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ALL_SHORTEST} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "selector_shortest_k",
			input: "(g MATCH SHORTEST 5 (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:SHORTEST_K K:5} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "selector_shortest_k_group",
			input: "(g MATCH SHORTEST 2 GROUP (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:SHORTEST_K K:2} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- restrictors -------------------------------------------------
		{
			name:  "restrictor_trail",
			input: "(g MATCH TRAIL (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Restrictor:TRAIL Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "restrictor_acyclic",
			input: "(g MATCH ACYCLIC (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Restrictor:ACYCLIC Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "restrictor_simple",
			input: "(g MATCH SIMPLE (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Restrictor:SIMPLE Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- path variable ----------------------------------------------
		{
			name:  "path_variable",
			input: "(g MATCH p = (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Variable:VarRef{Name:p} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "path_variable_with_restrictor",
			input: "(g MATCH TRAIL p = (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Restrictor:TRAIL Variable:VarRef{Name:p} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- multi-hop & comma list -------------------------------------
		{
			name:  "multi_hop",
			input: "(g MATCH (a)-[e1]->(b)-[e2]->(c))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e1}} NodePattern{Variable:VarRef{Name:b}} EdgePattern{Direction:-> Variable:VarRef{Name:e2}} NodePattern{Variable:VarRef{Name:c}}]}]}",
		},
		{
			name:  "pattern_list",
			input: "(g MATCH (a)-[e]->(b), (c)-[f]->(d))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]} GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:c}} EdgePattern{Direction:-> Variable:VarRef{Name:f}} NodePattern{Variable:VarRef{Name:d}}]}]}",
		},
		// ---- grouped sub-pattern ----------------------------------------
		{
			name:  "grouped_subpattern_paren",
			input: "(g MATCH (a)((x)-[e]->(y))*(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:x}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:y}}] Quantifier:PatternQuantifier{Min:0 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "grouped_subpattern_bracket",
			input: "(g MATCH (a)[(x)-[e]->(y)]+(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:x}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:y}}] Quantifier:PatternQuantifier{Min:1 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "grouped_subpattern_where_and_quantifier",
			input: "(g MATCH (a)((x)-[e]->(y) WHERE x.v > 0){2,5}(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:x}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:y}}] Where:BinaryExpr{Op:> Left:PathExpr{Root:VarRef{Name:x} Steps:[DotStep{Field:v}]} Right:NumberLit{Val:0}} Quantifier:PatternQuantifier{Min:2 Max:5}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		{
			name:  "grouped_subpattern_with_restrictor_and_pathvar",
			input: "(g MATCH (a)(TRAIL q = (x)-[e]->(y))+(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} GraphPattern{Restrictor:TRAIL Variable:VarRef{Name:q} Parts:[NodePattern{Variable:VarRef{Name:x}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:y}}] Quantifier:PatternQuantifier{Min:1 Max:-1}} NodePattern{Variable:VarRef{Name:b}}]}]}",
		},
		// ---- lhs is an expression ---------------------------------------
		{
			name:  "lhs_path_expr",
			input: "(mygraph.field MATCH (a))",
			want:  "MatchExpr{Expr:PathExpr{Root:VarRef{Name:mygraph} Steps:[DotStep{Field:field}]} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]}",
		},
		// ---- quoted label ------------------------------------------------
		{
			name:  "node_quoted_label",
			input: `(g MATCH (a:"Weird Label"))`,
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a} Labels:[Weird Label]}]}]}",
		},
		// ---- ANTLR-permissive forms (graphPart* allows zero parts) -------
		// matchPattern is `restrictor? variable? graphPart*`, so the legacy
		// ANTLR grammar (the antlr_fallback oracle) accepts an empty pattern,
		// a trailing/dangling edge, and a path-variable with no body. These
		// are syntactically valid; any "must have endpoints" rule is a
		// semantic concern, not a parser one. We follow ANTLR and accept.
		{
			name:  "empty_pattern",
			input: "(g MATCH )",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[]}]}",
		},
		{
			name:  "trailing_edge",
			input: "(g MATCH (a)-[e]->)",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}}]}]}",
		},
		{
			name:  "path_var_empty_body",
			input: "(g MATCH p = )",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Variable:VarRef{Name:p} Parts:[]}]}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			expr, err := p.ParseExpr()
			if err != nil {
				t.Fatalf("parse error for %q: %v", tc.input, err)
			}
			got := ast.NodeToString(expr)
			if got != tc.want {
				t.Errorf("AST mismatch for %q\n got: %s\nwant: %s", tc.input, got, tc.want)
			}
		})
	}
}

// TestParseGraphMatch_Reject feeds malformed MATCH expressions and asserts a
// parse error (NOT a panic, NOT a silent accept). Each case maps to a
// grammar constraint in PartiQLParser.g4.
func TestParseGraphMatch_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"node_unclosed", "(g MATCH (a)"},
		{"node_paren_unclosed", "(g MATCH (a -[e]->(b))"},
		{"edge_spec_unclosed", "(g MATCH (a)-[e->(b))"},
		{"leading_edge", "(g MATCH -[e]->(b))"},
		{"quantifier_brace_no_comma", "(g MATCH (a)-[e]->{3}(b))"},
		{"quantifier_brace_no_lower", "(g MATCH (a)-[e]->{,3}(b))"},
		{"quantifier_brace_unclosed", "(g MATCH (a)-[e]->{1,3(b))"},
		{"selector_shortest_no_k", "(g MATCH SHORTEST (a))"},
		{"label_missing_name", "(g MATCH (a:))"},
		{"label_not_identifier", "(g MATCH (a:123))"},
		{"trailing_garbage", "(g MATCH (a) garbage)"},
		{"bad_edge_double_bracket", "(g MATCH (a)-[[e]]->(b))"},
		{"group_subpattern_unclosed", "(g MATCH (a)((x)-[e]->(y)(b))"},
		{"all_without_shortest", "(g MATCH ALL (a))"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.ParseExpr()
			if err == nil {
				t.Fatalf("expected parse error for %q, got nil", tc.input)
			}
			// A deferred-feature stub is NOT an acceptable rejection here:
			// graph MATCH is implemented, so any rejection must be a real
			// syntax error.
			if strings.Contains(err.Error(), "deferred to") {
				t.Fatalf("got deferred-feature stub for %q (graph MATCH must be implemented): %v", tc.input, err)
			}
		})
	}
}

// TestParseGraphMatch_Goldens iterates every .partiql file under
// testdata/parser-graph/ and compares the parser's pretty-printed AST
// (ast.NodeToString of ParseExpr) against the matching .golden file. These
// durable golden artifacts capture the exact structural shape (nesting,
// grouping, quantifier placement) of representative GPML inputs.
//
// Run with `go test -update -run TestParseGraphMatch_Goldens ./partiql/parser/...`
// to regenerate after intentional AST shape changes.
func TestParseGraphMatch_Goldens(t *testing.T) {
	files, err := filepath.Glob("testdata/parser-graph/*.partiql")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no golden inputs found under testdata/parser-graph/")
	}
	for _, inPath := range files {
		name := strings.TrimSuffix(filepath.Base(inPath), ".partiql")
		t.Run(name, func(t *testing.T) {
			input, err := os.ReadFile(inPath)
			if err != nil {
				t.Fatal(err)
			}
			p := NewParser(strings.TrimRight(string(input), "\n"))
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
