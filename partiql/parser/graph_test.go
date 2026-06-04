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
			// matchSelector#SelectorAny: `ANY k=LITERAL_INTEGER?` — the count k
			// is preserved on PatternSelector.K (Codex finding: it was previously
			// consumed and discarded).
			name:  "selector_any_k",
			input: "(g MATCH ANY 3 (a)-[e]->(b))",
			want:  "MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ANY K:3} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}",
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

// TestParseGraphMatch_FromClause_Accept exercises the FROM-position graph match
// (tableBaseReference#TableBaseRefMatch, g4:405 — source=exprGraphMatchOne).
// This is the PRIMARY GPML usage. The full statement is parsed so the asserted
// AST shows where the MatchExpr lands in the SelectStmt.From — the whole point
// of these cases is that the match is REALLY CARRIED in the AST (the prior bug
// silently dropped it, yielding From:AliasedSource{Source:nil}), so each `want`
// pins the MatchExpr (and any alias) present under From.
//
// Two source forms both reach a MatchExpr-in-FROM:
//   - unparenthesised:  FROM g MATCH (a)         (exprGraphMatchOne, this fix)
//   - parenthesised:    FROM (g MATCH (a))       (exprGraphMatchMany via parens)
//
// The grammar makes the unparenthesised form a SINGLE pattern: a COMMA after the
// pattern belongs to the FROM tableReference comma-join, not to the pattern
// list, so `FROM g MATCH (a), h` is a cross-join of the match with `h` — covered
// below to lock that boundary in.
func TestParseGraphMatch_FromClause_Accept(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			// Unparenthesised single-node FROM match: the match must be the
			// From source, NOT a nil-source AliasedSource.
			name:  "from_unparenthesised_node",
			input: "SELECT * FROM g MATCH (a)",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]}}",
		},
		{
			// Parenthesised FROM match (exprGraphMatchMany shape) — also lands
			// as a real MatchExpr in From, alias-free.
			name:  "from_parenthesised_node",
			input: "SELECT * FROM (g MATCH (a))",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]}}",
		},
		{
			name:  "from_unparenthesised_edge",
			input: "SELECT * FROM g MATCH (a)-[e]->(b)",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}}",
		},
		{
			// AS alias on a FROM match must be preserved around the MatchExpr.
			name:  "from_match_as_alias",
			input: "SELECT * FROM g MATCH (a)-[e]->(b) AS r",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]} As:r}}",
		},
		{
			// AS + AT alias chain (TableBaseRefMatch: asIdent? atIdent? byIdent?).
			name:  "from_match_as_at_alias",
			input: "SELECT a FROM g MATCH (x)-[e]->(y) AS r AT pos",
			want:  "SelectStmt{Targets:[TargetEntry{Expr:VarRef{Name:a}}] From:AliasedSource{Source:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:x}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:y}}]}]} As:r At:pos}}",
		},
		{
			// Selector with a count on an unparenthesised FROM match
			// (exprGraphMatchOne -> gpmlPattern -> matchSelector#SelectorAny);
			// K is preserved.
			name:  "from_match_selector_any_k",
			input: "SELECT * FROM g MATCH ANY 3 (a)-[e]->(b)",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Selector:PatternSelector{Kind:ANY K:3} Parts:[NodePattern{Variable:VarRef{Name:a}} EdgePattern{Direction:-> Variable:VarRef{Name:e}} NodePattern{Variable:VarRef{Name:b}}]}]}}",
		},
		{
			// Boundary: a COMMA after the unparenthesised single pattern is the
			// FROM comma-join, NOT a second pattern. `FROM g MATCH (a), h` ==
			// (g MATCH (a)) CROSS JOIN h.
			name:  "from_match_then_comma_join",
			input: "SELECT * FROM g MATCH (a), h",
			want:  "SelectStmt{Star:true From:JoinExpr{Kind:CROSS Left:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]} Right:VarRef{Name:h}}}",
		},
		{
			// A path-expression graph source, unparenthesised.
			name:  "from_match_path_source",
			input: "SELECT * FROM mygraph.field MATCH (a)",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:PathExpr{Root:VarRef{Name:mygraph} Steps:[DotStep{Field:field}]} Patterns:[GraphPattern{Parts:[NodePattern{Variable:VarRef{Name:a}}]}]}}",
		},
		{
			// ANTLR-permissive (oracle = antlr_fallback): matchPattern is
			// `restrictor? variable? graphPart*`, so an empty pattern is valid.
			// This is the FROM-position analog of the `empty_pattern` accept
			// case and the documented divergence #3; `FROM g MATCH` is accepted
			// as an empty-pattern match, NOT a syntax error.
			name:  "from_match_empty_pattern",
			input: "SELECT * FROM g MATCH",
			want:  "SelectStmt{Star:true From:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[]}]}}",
		},
		{
			// The empty-pattern FROM match still accepts an alias chain.
			name:  "from_match_empty_pattern_alias",
			input: "SELECT * FROM g MATCH AS r",
			want:  "SelectStmt{Star:true From:AliasedSource{Source:MatchExpr{Expr:VarRef{Name:g} Patterns:[GraphPattern{Parts:[]}]} As:r}}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			stmts, err := p.parseScript()
			if err != nil {
				t.Fatalf("parse error for %q: %v", tc.input, err)
			}
			if len(stmts.Items) != 1 {
				t.Fatalf("expected 1 statement for %q, got %d", tc.input, len(stmts.Items))
			}
			got := ast.NodeToString(stmts.Items[0])
			if got != tc.want {
				t.Errorf("AST mismatch for %q\n got: %s\nwant: %s", tc.input, got, tc.want)
			}
			// Guard against regressing the silent-drop bug: the rendered AST
			// must contain a MatchExpr and must NOT contain a nil source.
			if !strings.Contains(got, "MatchExpr{") {
				t.Errorf("FROM match dropped from AST for %q: %s", tc.input, got)
			}
			if strings.Contains(got, "Source:}") || strings.Contains(got, "Source:nil") {
				t.Errorf("FROM match produced a nil-source AliasedSource for %q: %s", tc.input, got)
			}
		})
	}
}

// TestParseGraphMatch_FromClause_Reject feeds malformed FROM-position graph
// matches and asserts a real parse error (not a panic, not a silent accept, not
// a deferred-feature stub).
func TestParseGraphMatch_FromClause_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		// NOTE: `FROM g MATCH` (no graphPart) is NOT here — per the
		// antlr_fallback oracle an empty pattern is accepted (divergence #3);
		// it is covered as an accept case above.
		//
		// A leading edge with no source node is rejected in pattern position.
		{"from_match_leading_edge", "SELECT * FROM g MATCH -[e]->(b)"},
		// Unclosed node in the FROM pattern.
		{"from_match_node_unclosed", "SELECT * FROM g MATCH (a"},
		// Unclosed edge spec.
		{"from_match_edge_spec_unclosed", "SELECT * FROM g MATCH (a)-[e->(b)"},
		// Bad quantifier (single value, no comma) — follows ANTLR {m,n}/{m,}.
		{"from_match_quantifier_no_comma", "SELECT * FROM g MATCH (a)-[e]->{3}(b)"},
		// SHORTEST requires a count.
		{"from_match_shortest_no_k", "SELECT * FROM g MATCH SHORTEST (a)"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.parseScript()
			if err == nil {
				t.Fatalf("expected parse error for %q, got nil", tc.input)
			}
			if strings.Contains(err.Error(), "deferred to") {
				t.Fatalf("got deferred-feature stub for %q (FROM graph MATCH must be implemented): %v", tc.input, err)
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
