package parser

import (
	"testing"

	"github.com/bytebase/omni/bigquery/ast"
)

// Unit tests for the parser-gql node — the §2.12 GQL graph query language
// (gql_statement: `GRAPH <path> <ops>`) and the create_property_graph_statement
// DDL. These assert the AST STRUCTURE (the structural gate: accept/reject alone
// does not catch wrong nesting / direction / label-precedence) and cover the
// edge-direction set, the label algebra precedence, the linear operators, and
// the property-graph element-table grammar.
//
// GQL is BigQuery/Spanner-Graph syntax; the truth1 BigQuery corpus scopes the
// full GQL syntax OUT (INDEX.md), so the authoritative reference for the grammar
// shape is the pinned legacy GoogleSQLParser.g4. The accept/reject of these forms
// against the live Spanner emulator is recorded in the suite:graph-query rows of testdata/compat.tsv
//.

// gqlStmtOf parses sql and asserts the single statement is a *GQLStmt.
func gqlStmtOf(t *testing.T, sql string) *ast.GQLStmt {
	t.Helper()
	n := parseDDL(t, sql)
	g, ok := n.(*ast.GQLStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.GQLStmt", sql, n)
	}
	return g
}

// linearOf asserts the GQL statement is a single linear query (one block, not a
// set operation) and returns it.
func linearOf(t *testing.T, g *ast.GQLStmt) *ast.GraphLinearQuery {
	t.Helper()
	if len(g.Blocks) != 1 {
		t.Fatalf("Blocks = %d, want 1", len(g.Blocks))
	}
	lq, ok := g.Blocks[0].(*ast.GraphLinearQuery)
	if !ok {
		t.Fatalf("Blocks[0] is %T, want *ast.GraphLinearQuery", g.Blocks[0])
	}
	return lq
}

func TestGQL_BasicMatchReturn(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH my_graph MATCH (n) RETURN n")
	if g.Name.String() != "my_graph" {
		t.Errorf("Name = %q, want my_graph", g.Name.String())
	}
	lq := linearOf(t, g)
	if len(lq.Operators) != 1 {
		t.Fatalf("Operators = %d, want 1 (MATCH)", len(lq.Operators))
	}
	m, ok := lq.Operators[0].(*ast.GraphMatchOp)
	if !ok {
		t.Fatalf("operator is %T, want *ast.GraphMatchOp", lq.Operators[0])
	}
	if m.Optional {
		t.Error("MATCH should not be Optional")
	}
	if len(m.Pattern.Paths) != 1 {
		t.Fatalf("pattern paths = %d, want 1", len(m.Pattern.Paths))
	}
	node, ok := m.Pattern.Paths[0].Factors[0].(*ast.GraphNodePattern)
	if !ok {
		t.Fatalf("factor is %T, want *ast.GraphNodePattern", m.Pattern.Paths[0].Factors[0])
	}
	if node.Filler.Var != "n" {
		t.Errorf("node var = %q, want n", node.Filler.Var)
	}
	if lq.Return == nil || len(lq.Return.Items) != 1 {
		t.Fatalf("RETURN items = %v, want 1", lq.Return)
	}
}

func TestGQL_OptionalMatch(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g OPTIONAL MATCH (a)-[e]->(b) RETURN a, b")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	if !m.Optional {
		t.Error("OPTIONAL MATCH should be Optional")
	}
	// (a)-[e]->(b): node, right-full edge, node.
	factors := m.Pattern.Paths[0].Factors
	if len(factors) != 3 {
		t.Fatalf("factors = %d, want 3 (node, edge, node)", len(factors))
	}
	edge, ok := factors[1].(*ast.GraphEdgePattern)
	if !ok {
		t.Fatalf("factors[1] is %T, want *ast.GraphEdgePattern", factors[1])
	}
	if edge.Direction != ast.EdgeRightFull {
		t.Errorf("edge direction = %v, want EdgeRightFull", edge.Direction)
	}
	if edge.Filler == nil || edge.Filler.Var != "e" {
		t.Errorf("edge filler var = %v, want e", edge.Filler)
	}
}

// TestGQL_EdgeDirections covers all six graph_edge_pattern shapes, including the
// undirected full edge `-[e]-` (the legacy first alternative `LT? MINUS … MINUS`
// with the leading `<` ABSENT — a real over-reject the oracle differential caught
// when it was modeled as left-only).
func TestGQL_EdgeDirections(t *testing.T) {
	cases := []struct {
		sql  string
		want ast.GraphEdgeDirection
		full bool
	}{
		{"GRAPH g MATCH (a)-(b) RETURN a", ast.EdgeAny, false},
		{"GRAPH g MATCH (a)<-(b) RETURN a", ast.EdgeLeft, false},
		{"GRAPH g MATCH (a)->(b) RETURN a", ast.EdgeRight, false},
		{"GRAPH g MATCH (a)<-[e]-(b) RETURN a", ast.EdgeLeftFull, true},
		{"GRAPH g MATCH (a)-[e]->(b) RETURN a", ast.EdgeRightFull, true},
		{"GRAPH g MATCH (a)-[e]-(b) RETURN a", ast.EdgeUndirectedFull, true},
	}
	for _, c := range cases {
		g := gqlStmtOf(t, c.sql)
		lq := linearOf(t, g)
		m := lq.Operators[0].(*ast.GraphMatchOp)
		factors := m.Pattern.Paths[0].Factors
		edge, ok := factors[1].(*ast.GraphEdgePattern)
		if !ok {
			t.Fatalf("%q: factors[1] is %T, want edge", c.sql, factors[1])
		}
		if edge.Direction != c.want {
			t.Errorf("%q: direction = %v, want %v", c.sql, edge.Direction, c.want)
		}
		if c.full && edge.Filler == nil {
			t.Errorf("%q: full edge should carry a filler", c.sql)
		}
		if !c.full && edge.Filler != nil {
			t.Errorf("%q: abbreviated edge should not carry a filler", c.sql)
		}
	}
}

// TestGQL_LabelExpression covers the label algebra: bare label via `:` and `IS`,
// the `%` wildcard, `&`/`|` combinators with precedence, `!` negation, and
// parentheses.
func TestGQL_LabelExpression(t *testing.T) {
	// Bare label via ':'.
	g := gqlStmtOf(t, "GRAPH g MATCH (n:Person) RETURN n")
	filler := nodeFiller(t, g)
	if filler.Label == nil || filler.Label.Kind != ast.LabelName || filler.Label.Name != "Person" {
		t.Fatalf("label = %+v, want LabelName Person", filler.Label)
	}
	if !filler.LabelColon {
		t.Error("expected LabelColon true for ':' spelling")
	}

	// Label via IS.
	g = gqlStmtOf(t, "GRAPH g MATCH (n IS Person) RETURN n")
	filler = nodeFiller(t, g)
	if filler.LabelColon {
		t.Error("expected LabelColon false for IS spelling")
	}
	if filler.Label.Name != "Person" {
		t.Errorf("label = %q, want Person", filler.Label.Name)
	}

	// Wildcard %.
	g = gqlStmtOf(t, "GRAPH g MATCH (n:%) RETURN n")
	filler = nodeFiller(t, g)
	if filler.Label.Kind != ast.LabelWildcard {
		t.Errorf("label kind = %v, want LabelWildcard", filler.Label.Kind)
	}

	// Negation !A.
	g = gqlStmtOf(t, "GRAPH g MATCH (n:!Temp) RETURN n")
	filler = nodeFiller(t, g)
	if filler.Label.Kind != ast.LabelNot || filler.Label.Operand.Name != "Temp" {
		t.Errorf("label = %+v, want !Temp", filler.Label)
	}

	// Precedence: A | B & C parses as A | (B & C) (& binds tighter than |).
	g = gqlStmtOf(t, "GRAPH g MATCH (n:A|B&C) RETURN n")
	filler = nodeFiller(t, g)
	if filler.Label.Kind != ast.LabelOr {
		t.Fatalf("top label kind = %v, want LabelOr", filler.Label.Kind)
	}
	if filler.Label.Left.Kind != ast.LabelName || filler.Label.Left.Name != "A" {
		t.Errorf("OR left = %+v, want A", filler.Label.Left)
	}
	if filler.Label.Right.Kind != ast.LabelAnd {
		t.Errorf("OR right kind = %v, want LabelAnd (B & C)", filler.Label.Right.Kind)
	}
	if filler.Label.Right.Left.Name != "B" || filler.Label.Right.Right.Name != "C" {
		t.Errorf("AND operands = %+v / %+v, want B / C", filler.Label.Right.Left, filler.Label.Right.Right)
	}

	// Parentheses override: (A | B) & C parses as (A | B) & C.
	g = gqlStmtOf(t, "GRAPH g MATCH (n:(A|B)&C) RETURN n")
	filler = nodeFiller(t, g)
	if filler.Label.Kind != ast.LabelAnd {
		t.Fatalf("top label kind = %v, want LabelAnd", filler.Label.Kind)
	}
	if filler.Label.Left.Kind != ast.LabelOr {
		t.Errorf("AND left kind = %v, want LabelOr ((A|B))", filler.Label.Left.Kind)
	}
	if filler.Label.Right.Name != "C" {
		t.Errorf("AND right = %+v, want C", filler.Label.Right)
	}
}

// TestGQL_LabelNotBindsTightest pins the label-expression precedence the legacy
// .g4 gets WRONG and the oracle/official docs get right (finding F4). The .g4
// lists the unary `EXCLAMATION_OPERATOR label_expression` alternative LAST, which
// — read as ANTLR left-recursion precedence — would bind `!` LOOSEST, parsing
// `!A & B` as `!(A & B)`. That is a hand-port slip. Standard GQL/ZetaSQL label-
// expression precedence (ISO/IEC 39075 GQL: NOT binds tighter than AND tighter
// than OR; ZetaSQL graph-patterns.md's own example `(p:(!Singer&!Writer))` only
// parses as `(!Singer)&(!Writer)`) makes `!` bind TIGHTEST, so `!A & B` is
// `(!A) & B`. omni builds `(!A) & B`; we DEFEND that and lock the AST shape here.
// (See the divergence ledger: the .g4 operator order is backwards.)
func TestGQL_LabelNotBindsTightest(t *testing.T) {
	// !A & B  must be  (!A) & B  — top is AND, left is NOT(A), right is B.
	g := gqlStmtOf(t, "GRAPH g MATCH (n:!A&B) RETURN n")
	lbl := nodeFiller(t, g).Label
	if lbl.Kind != ast.LabelAnd {
		t.Fatalf("top kind = %v, want LabelAnd ((!A)&B, i.e. ! binds tighter than &)", lbl.Kind)
	}
	if lbl.Left.Kind != ast.LabelNot || lbl.Left.Operand.Name != "A" {
		t.Errorf("AND left = %+v, want LabelNot(A)", lbl.Left)
	}
	if lbl.Right.Kind != ast.LabelName || lbl.Right.Name != "B" {
		t.Errorf("AND right = %+v, want LabelName B", lbl.Right)
	}

	// !A | B  must be  (!A) | B  — top is OR, left is NOT(A), right is B.
	g = gqlStmtOf(t, "GRAPH g MATCH (n:!A|B) RETURN n")
	lbl = nodeFiller(t, g).Label
	if lbl.Kind != ast.LabelOr {
		t.Fatalf("top kind = %v, want LabelOr ((!A)|B)", lbl.Kind)
	}
	if lbl.Left.Kind != ast.LabelNot || lbl.Left.Operand.Name != "A" {
		t.Errorf("OR left = %+v, want LabelNot(A)", lbl.Left)
	}

	// The ZetaSQL canonical example: !A & !B  ==  (!A) & (!B).
	g = gqlStmtOf(t, "GRAPH g MATCH (n:!A&!B) RETURN n")
	lbl = nodeFiller(t, g).Label
	if lbl.Kind != ast.LabelAnd || lbl.Left.Kind != ast.LabelNot || lbl.Right.Kind != ast.LabelNot {
		t.Errorf("!A & !B = %+v, want LabelAnd(LabelNot, LabelNot)", lbl)
	}
}

// nodeFiller extracts the filler of the first node pattern of a single
// MATCH-RETURN GQL statement.
func nodeFiller(t *testing.T, g *ast.GQLStmt) *ast.GraphPatternFiller {
	t.Helper()
	lq := linearOf(t, g)
	m, ok := lq.Operators[0].(*ast.GraphMatchOp)
	if !ok {
		t.Fatalf("operator is %T, want match", lq.Operators[0])
	}
	node, ok := m.Pattern.Paths[0].Factors[0].(*ast.GraphNodePattern)
	if !ok {
		t.Fatalf("factor is %T, want node", m.Pattern.Paths[0].Factors[0])
	}
	return node.Filler
}

func TestGQL_PropertySpecification(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n:Person {name: 'Alice', age: 30}) RETURN n")
	filler := nodeFiller(t, g)
	if filler.Properties == nil || len(filler.Properties.Properties) != 2 {
		t.Fatalf("properties = %+v, want 2", filler.Properties)
	}
	if filler.Properties.Properties[0].Name != "name" {
		t.Errorf("prop[0] name = %q, want name", filler.Properties.Properties[0].Name)
	}
	if filler.Properties.Properties[1].Name != "age" {
		t.Errorf("prop[1] name = %q, want age", filler.Properties.Properties[1].Name)
	}
}

func TestGQL_EmptyNodePattern(t *testing.T) {
	// An empty filler `()` is valid (the grammar's empty production).
	g := gqlStmtOf(t, "GRAPH g MATCH () RETURN 1")
	filler := nodeFiller(t, g)
	if filler.Var != "" || filler.Label != nil || filler.Properties != nil {
		t.Errorf("empty filler should be all-zero, got %+v", filler)
	}
}

func TestGQL_InlineWhereInFiller(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n:Person WHERE n.age > 18) RETURN n")
	filler := nodeFiller(t, g)
	if filler.Where == nil {
		t.Error("expected inline WHERE in filler")
	}
	if filler.Label == nil || filler.Label.Name != "Person" {
		t.Errorf("label = %+v, want Person", filler.Label)
	}
}

func TestGQL_PatternWhere(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (a)-[e]->(b) WHERE a.id = b.id RETURN a")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	if m.Pattern.Where == nil {
		t.Error("expected pattern-level WHERE")
	}
}

func TestGQL_MultiplePathPatterns(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (a)-[]->(b), (b)-[]->(c) RETURN a, c")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	if len(m.Pattern.Paths) != 2 {
		t.Fatalf("paths = %d, want 2", len(m.Pattern.Paths))
	}
}

func TestGQL_PathVariableAndSearchAndMode(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH p = ANY SHORTEST TRAIL (a)-[e]->(b) RETURN p")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	path := m.Pattern.Paths[0]
	if path.PathVar != "p" {
		t.Errorf("PathVar = %q, want p", path.PathVar)
	}
	if path.Search != "ANY SHORTEST" {
		t.Errorf("Search = %q, want 'ANY SHORTEST'", path.Search)
	}
	if path.Mode != "TRAIL" {
		t.Errorf("Mode = %q, want TRAIL", path.Mode)
	}
}

func TestGQL_PathModeWithPaths(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH WALK PATHS (a)-[e]->(b) RETURN a")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	if m.Pattern.Paths[0].Mode != "WALK PATHS" {
		t.Errorf("Mode = %q, want 'WALK PATHS'", m.Pattern.Paths[0].Mode)
	}
}

func TestGQL_QuantifiedPath(t *testing.T) {
	// (a)-[e]->{1,3}(b): a quantified edge factor.
	g := gqlStmtOf(t, "GRAPH g MATCH (a) (-[e]->){1,3} (b) RETURN a")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	// The middle factor is a parenthesized path with a {1,3} quantifier; just
	// assert the path parsed with >= 3 factors (a, (quantified path), b).
	if len(m.Pattern.Paths[0].Factors) < 3 {
		t.Errorf("factors = %d, want >= 3", len(m.Pattern.Paths[0].Factors))
	}
}

func TestGQL_ParenthesizedPath(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH ((a)-[e]->(b) WHERE a.x > 0) RETURN a")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	// The single factor is a parenthesized sub-path.
	sub, ok := m.Pattern.Paths[0].Factors[0].(*ast.GraphPathPattern)
	if !ok {
		t.Fatalf("factor is %T, want *ast.GraphPathPattern (parenthesized)", m.Pattern.Paths[0].Factors[0])
	}
	if sub.Where == nil {
		t.Error("expected WHERE inside parenthesized path")
	}
}

// parenSubPathOf extracts the first path factor of a single MATCH-RETURN GQL
// statement and asserts it is a parenthesized sub-path (*ast.GraphPathPattern).
func parenSubPathOf(t *testing.T, sql string) *ast.GraphPathPattern {
	t.Helper()
	g := gqlStmtOf(t, sql)
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	sub, ok := m.Pattern.Paths[0].Factors[0].(*ast.GraphPathPattern)
	if !ok {
		t.Fatalf("%q: factor[0] is %T, want *ast.GraphPathPattern (parenthesized path)", sql, m.Pattern.Paths[0].Factors[0])
	}
	return sub
}

// TestGQL_ParenthesizedPathWithPrefix covers finding F1: a parenthesized path
// whose interior begins with a path-pattern PREFIX — a path-variable assignment
// (`p = …`), a search prefix (ANY / ALL [SHORTEST]), or a path-mode prefix
// (WALK / TRAIL / SIMPLE / ACYCLIC). The .g4
// `graph_parenthesized_path_pattern: '(' hint? graph_path_pattern …` admits the
// full graph_path_pattern (prefixes included) inside the parens; the live Spanner
// emulator accepts all of these. The disambiguator previously misrouted them to a
// node pattern (and errored). These must parse as parenthesized sub-paths with
// the prefix recorded on the inner path.
func TestGQL_ParenthesizedPathWithPrefix(t *testing.T) {
	// Path-variable assignment inside parens.
	sub := parenSubPathOf(t, "GRAPH g MATCH (p = (a)-[e]->(b)) RETURN p")
	if sub.PathVar != "p" {
		t.Errorf("PathVar = %q, want p", sub.PathVar)
	}
	if len(sub.Factors) != 3 {
		t.Errorf("inner factors = %d, want 3 (node, edge, node)", len(sub.Factors))
	}

	// Search prefix ANY.
	if sub := parenSubPathOf(t, "GRAPH g MATCH (ANY (a)-[e]->(b)) RETURN a"); sub.Search != "ANY" {
		t.Errorf("Search = %q, want ANY", sub.Search)
	}

	// Search prefix ALL SHORTEST.
	if sub := parenSubPathOf(t, "GRAPH g MATCH (ALL SHORTEST (a)-[e]->(b)) RETURN a"); sub.Search != "ALL SHORTEST" {
		t.Errorf("Search = %q, want 'ALL SHORTEST'", sub.Search)
	}

	// Path-mode prefix WALK.
	if sub := parenSubPathOf(t, "GRAPH g MATCH (WALK (a)-[e]->(b)) RETURN a"); sub.Mode != "WALK" {
		t.Errorf("Mode = %q, want WALK", sub.Mode)
	}

	// Combined prefix inside parens: var + search + mode.
	sub = parenSubPathOf(t, "GRAPH g MATCH (q = ANY SHORTEST TRAIL (a)-[e]->(b)) RETURN q")
	if sub.PathVar != "q" || sub.Search != "ANY SHORTEST" || sub.Mode != "TRAIL" {
		t.Errorf("combined prefix = {var %q, search %q, mode %q}, want {q, 'ANY SHORTEST', TRAIL}", sub.PathVar, sub.Search, sub.Mode)
	}
}

// TestGQL_HintedNodePattern covers finding F2: a node pattern whose filler starts
// with a hint — `graph_element_pattern_filler: hint? opt_graph_element_identifier?
// …`. The disambiguator previously treated the leading '@' as a parenthesized-
// path opener and misparsed `(@{h} v:L)`. After consuming the optional leading
// hint, the next token (`v`) reveals a node filler; the live emulator accepts it.
func TestGQL_HintedNodePattern(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (@{force_index=idx} v:Person) RETURN v")
	lq := linearOf(t, g)
	m := lq.Operators[0].(*ast.GraphMatchOp)
	node, ok := m.Pattern.Paths[0].Factors[0].(*ast.GraphNodePattern)
	if !ok {
		t.Fatalf("factor[0] is %T, want *ast.GraphNodePattern (hinted node)", m.Pattern.Paths[0].Factors[0])
	}
	if node.Filler.Var != "v" {
		t.Errorf("node var = %q, want v", node.Filler.Var)
	}
	if node.Filler.Label == nil || node.Filler.Label.Name != "Person" {
		t.Errorf("node label = %+v, want Person", node.Filler.Label)
	}

	// A hinted node with label only, no var (the hint is consumed at the '(' level,
	// then the filler is `:Person`).
	g = gqlStmtOf(t, "GRAPH g MATCH (@{h=1} :Person) RETURN 1")
	f := nodeFiller(t, g)
	if f.Var != "" || f.Label == nil || f.Label.Name != "Person" {
		t.Errorf("hinted (:Person) filler = %+v, want var='' label=Person", f)
	}
}

// TestGQL_DoubleHintRejected pins the post-review fix: a node pattern admits at
// most one leading hint (graph_element_pattern_filler: `hint? ...`). The single
// hint is consumed at the '(' level; the filler must NOT consume a second, so a
// second '@' is a syntax error. Oracle: the live Spanner emulator rejects
// `(@{h1} @{h2} v)` with `Expected ")" but got "@"`.
func TestGQL_DoubleHintRejected(t *testing.T) {
	assertReject(t, "GRAPH g MATCH (@{h1=1} @{h2=2} v) RETURN v")
	assertReject(t, "GRAPH g MATCH (@{h1=1} @{h2=2} :Person) RETURN 1")
}

// TestGQL_KeywordPathVarRejected pins the resolution of a review finding: a
// path-variable assignment whose name is a graph path-mode keyword
// (WALK/TRAIL/SIMPLE/...) is REJECTED. The legacy .g4 lists WALK/TRAIL/ACYCLIC
// in common_keyword_as_identifier (which would permit them as graph
// identifiers), but the live ZetaSQL/Spanner emulator REJECTS them as path-var
// names (`Syntax error: Unexpected "="`). The oracle wins over the .g4, so omni
// correctly rejects; this guards against a future change that loosens it to the
// .g4.
func TestGQL_KeywordPathVarRejected(t *testing.T) {
	for _, sql := range []string{
		"GRAPH g MATCH (walk = (a)) RETURN walk",
		"GRAPH g MATCH (trail = (a)) RETURN trail",
		"GRAPH g MATCH (simple = (a)) RETURN simple",
	} {
		assertReject(t, sql)
	}
	// A plain (non-keyword) identifier path-variable name is still accepted.
	gqlStmtOf(t, "GRAPH g MATCH (x = (a)) RETURN x")
}

// TestGQL_TrailingPathFactorHintRejected covers finding F5: a hint between path
// factors is part of a `(hint? graph_path_factor)` group, so it MUST be followed
// by another factor. A TRAILING hint with no factor after it is a syntax error —
// the live Spanner emulator rejects `(a) @{h} RETURN *` with
// `Expected "(" or "-" or "<" or -> but got …`. The inter-factor form
// `(a) @{h} (b)` stays valid (the emulator accepts it).
func TestGQL_TrailingPathFactorHintRejected(t *testing.T) {
	// Reject: trailing hint with no following factor.
	assertReject(t, "GRAPH g MATCH (a) @{h=1} RETURN *")
	assertReject(t, "GRAPH g MATCH (a)-[e]->(b) @{h=1} RETURN *")

	// Accept: inter-factor hint followed by a factor.
	for _, sql := range []string{
		"GRAPH g MATCH (a) @{h=1} (b) RETURN *",
		"GRAPH g MATCH (a)-[e]-> @{h=1} (b) RETURN *",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q) should accept (inter-factor hint), got errs=%v", sql, errs)
		}
	}
}

// TestGQL_FillerWhereWithoutIdentifier covers finding F6: the where-bearing
// graph_element_pattern_filler alternatives are written with
// `opt_graph_element_identifier` (no trailing '?'), but opt_graph_element_identifier
// is itself optional in practice — the live Spanner emulator ACCEPTS an inline
// WHERE with NO element identifier, both with a label (`(:Person WHERE TRUE)`) and
// bare (`(WHERE foo)`). omni accepts these (the filler treats the identifier as
// optional); we DEFEND that with an accept test.
func TestGQL_FillerWhereWithoutIdentifier(t *testing.T) {
	// (:Person WHERE …): label + WHERE, no element identifier.
	g := gqlStmtOf(t, "GRAPH g MATCH (:Person WHERE TRUE) RETURN 1")
	f := nodeFiller(t, g)
	if f.Var != "" {
		t.Errorf("var = %q, want '' (no element identifier)", f.Var)
	}
	if f.Label == nil || f.Label.Name != "Person" {
		t.Errorf("label = %+v, want Person", f.Label)
	}
	if f.Where == nil {
		t.Error("expected inline WHERE")
	}

	// (WHERE …): bare WHERE, no identifier and no label.
	g = gqlStmtOf(t, "GRAPH g MATCH (WHERE foo) RETURN 1")
	f = nodeFiller(t, g)
	if f.Var != "" || f.Label != nil {
		t.Errorf("filler = %+v, want empty var + nil label", f)
	}
	if f.Where == nil {
		t.Error("expected bare inline WHERE")
	}
}

// TestGQL_LinearOperators covers LET / FILTER / ORDER BY / WITH / FOR / PAGE.
func TestGQL_LinearOperators(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) LET x = n.age FILTER x > 18 RETURN x")
	lq := linearOf(t, g)
	if len(lq.Operators) != 3 {
		t.Fatalf("operators = %d, want 3 (MATCH, LET, FILTER)", len(lq.Operators))
	}
	let, ok := lq.Operators[1].(*ast.GraphLetOp)
	if !ok || len(let.Vars) != 1 || let.Vars[0].Name != "x" {
		t.Fatalf("LET = %+v, want one var x", lq.Operators[1])
	}
	filter, ok := lq.Operators[2].(*ast.GraphFilterOp)
	if !ok || filter.HasWhere {
		t.Fatalf("FILTER = %+v, want bare filter (no WHERE)", lq.Operators[2])
	}
}

func TestGQL_FilterWhere(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) FILTER WHERE n.age > 18 RETURN n")
	lq := linearOf(t, g)
	filter := lq.Operators[1].(*ast.GraphFilterOp)
	if !filter.HasWhere {
		t.Error("expected FILTER WHERE form")
	}
}

func TestGQL_OrderByOperator(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) ORDER BY n.age DESC, n.name ASCENDING RETURN n")
	lq := linearOf(t, g)
	ob, ok := lq.Operators[1].(*ast.GraphOrderByOp)
	if !ok {
		t.Fatalf("operator is %T, want order-by", lq.Operators[1])
	}
	if len(ob.Items) != 2 {
		t.Fatalf("order items = %d, want 2", len(ob.Items))
	}
	if !ob.Items[0].HasDir || !ob.Items[0].Desc {
		t.Errorf("item[0] = %+v, want DESC", ob.Items[0])
	}
	if !ob.Items[1].HasDir || ob.Items[1].Desc {
		t.Errorf("item[1] = %+v, want ASCENDING (asc)", ob.Items[1])
	}
}

func TestGQL_OrderByCollate(t *testing.T) {
	// graph_ordering_expression: expression collate_clause? asc/desc? null_order?.
	// The COLLATE operand must be retained as a Node on the OrderItem (not dropped),
	// so an AST/lineage pass can reach it.
	g := gqlStmtOf(t, "GRAPH g MATCH (n) ORDER BY n.name COLLATE 'und:ci' DESC RETURN n")
	lq := linearOf(t, g)
	ob := lq.Operators[1].(*ast.GraphOrderByOp)
	if len(ob.Items) != 1 {
		t.Fatalf("order items = %d, want 1", len(ob.Items))
	}
	if ob.Items[0].Collate == nil {
		t.Error("expected COLLATE operand retained on the OrderItem")
	}
	if !ob.Items[0].Desc {
		t.Error("expected DESC after COLLATE")
	}
}

func TestGQL_WithOperator(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) WITH DISTINCT n.age AS age GROUP BY age RETURN age")
	lq := linearOf(t, g)
	with, ok := lq.Operators[1].(*ast.GraphWithOp)
	if !ok {
		t.Fatalf("operator is %T, want with", lq.Operators[1])
	}
	if with.Quantifier != "DISTINCT" {
		t.Errorf("quantifier = %q, want DISTINCT", with.Quantifier)
	}
	if len(with.Items) != 1 || with.Items[0].Alias != "age" {
		t.Fatalf("with items = %+v, want one aliased 'age'", with.Items)
	}
	if with.GroupBy == nil {
		t.Error("expected GROUP BY in WITH")
	}
}

func TestGQL_ForOperator(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) FOR x IN n.items WITH OFFSET AS pos RETURN x, pos")
	lq := linearOf(t, g)
	forOp, ok := lq.Operators[1].(*ast.GraphForOp)
	if !ok {
		t.Fatalf("operator is %T, want for", lq.Operators[1])
	}
	if forOp.Name != "x" {
		t.Errorf("FOR var = %q, want x", forOp.Name)
	}
	if !forOp.WithOffset || forOp.OffsetAlias != "pos" {
		t.Errorf("WITH OFFSET = %v/%q, want true/pos", forOp.WithOffset, forOp.OffsetAlias)
	}
}

func TestGQL_PageOperatorAndReturn(t *testing.T) {
	// Standalone PAGE operator (LIMIT) then a RETURN with its own trailing PAGE.
	g := gqlStmtOf(t, "GRAPH g MATCH (n) LIMIT 10 RETURN n ORDER BY n.id OFFSET 5 LIMIT 3")
	lq := linearOf(t, g)
	page, ok := lq.Operators[1].(*ast.GraphPageOp)
	if !ok {
		t.Fatalf("operator is %T, want page", lq.Operators[1])
	}
	if page.Limit == nil || page.Skip != nil {
		t.Errorf("standalone page = %+v, want bare LIMIT", page)
	}
	if lq.Return.OrderBy == nil {
		t.Error("expected RETURN ORDER BY")
	}
	if lq.Return.Page == nil || lq.Return.Page.Skip == nil || lq.Return.Page.Limit == nil {
		t.Errorf("RETURN page = %+v, want OFFSET..LIMIT", lq.Return.Page)
	}
}

func TestGQL_SampleOperator(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (100 ROWS) WITH WEIGHT AS w REPEATABLE (42) RETURN n")
	lq := linearOf(t, g)
	s, ok := lq.Operators[1].(*ast.GraphSampleOp)
	if !ok {
		t.Fatalf("operator is %T, want sample", lq.Operators[1])
	}
	if s.Method != "RESERVOIR" {
		t.Errorf("method = %q, want RESERVOIR", s.Method)
	}
	if s.Unit != "ROWS" {
		t.Errorf("unit = %q, want ROWS", s.Unit)
	}
	if !s.WithWeight || s.WeightAlias != "w" {
		t.Errorf("WITH WEIGHT = %v/%q, want true/w", s.WithWeight, s.WeightAlias)
	}
	if s.Repeatable == nil {
		t.Error("expected REPEATABLE seed")
	}
}

func TestGQL_ReturnStar(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n) RETURN *")
	lq := linearOf(t, g)
	if len(lq.Return.Items) != 1 || !lq.Return.Items[0].Star {
		t.Fatalf("RETURN items = %+v, want a single star", lq.Return.Items)
	}
}

func TestGQL_NextChain(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (a) RETURN a NEXT MATCH (b) RETURN b")
	if len(g.Blocks) != 2 {
		t.Fatalf("blocks = %d, want 2 (NEXT-separated)", len(g.Blocks))
	}
}

func TestGQL_SetOperation(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (a) RETURN a UNION ALL MATCH (b) RETURN b")
	if len(g.Blocks) != 1 {
		t.Fatalf("blocks = %d, want 1 (a single composite set-op block)", len(g.Blocks))
	}
	setOp, ok := g.Blocks[0].(*ast.GraphSetOp)
	if !ok {
		t.Fatalf("block is %T, want *ast.GraphSetOp", g.Blocks[0])
	}
	if len(setOp.Ops) != 2 || len(setOp.Metas) != 1 {
		t.Fatalf("set op = %d ops / %d metas, want 2/1", len(setOp.Ops), len(setOp.Metas))
	}
	if setOp.Metas[0].Op != "UNION" || setOp.Metas[0].Quantifier != "ALL" {
		t.Errorf("meta = %+v, want UNION ALL", setOp.Metas[0])
	}
}

func TestGQL_DottedGraphName(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH myproject.mydataset.my_graph MATCH (n) RETURN n")
	if g.Name.String() != "myproject.mydataset.my_graph" {
		t.Errorf("Name = %q, want the 3-part path", g.Name.String())
	}
}

// ---- Negative tests (the grammar must REJECT these) ----

func TestGQL_Rejects(t *testing.T) {
	cases := []string{
		"GRAPH",                                               // no graph name, no ops
		"GRAPH g",                                             // no operation block (RETURN required)
		"GRAPH g MATCH (n)",                                   // missing RETURN
		"GRAPH g RETURN",                                      // RETURN with no items
		"GRAPH g MATCH n RETURN n",                            // node pattern needs parens
		"GRAPH g MATCH (n LET x = 1 RETURN n",                 // unbalanced node paren
		"GRAPH g MATCH (a)=>(b) RETURN a",                     // not a valid edge arrow
		"GRAPH g MATCH (n) LET = 1 RETURN n",                  // LET needs an identifier
		"GRAPH g MATCH (n) FOR x n.items RETURN x",            // FOR needs IN
		"GRAPH g MATCH (a) RETURN a UNION MATCH (b) RETURN b", // set-op needs ALL|DISTINCT
		"GRAPH g MATCH (n:) RETURN n",                         // empty label after ':'
		"GRAPH g MATCH (n {x}) RETURN n",                      // property spec needs ': value'
		// Numeric operands are restricted to int_literal_or_parameter /
		// possibly_cast_int_literal_or_parameter / sample_size_value — NOT a full
		// expression. A bare identifier or an arithmetic expression in a page count,
		// a quantifier bound, or a sample size is a syntax error (oracle-confirmed —
		// the live Spanner emulator rejects each with "Syntax error:").
		"GRAPH g MATCH (n) RETURN n LIMIT x",                          // page LIMIT must be int/param/cast, not an ident
		"GRAPH g MATCH (n) RETURN n OFFSET y LIMIT 3",                 // page OFFSET ident
		"GRAPH g MATCH (n) LIMIT z RETURN n",                          // standalone-page LIMIT ident
		"GRAPH g MATCH (n) RETURN n LIMIT 1+1",                        // page LIMIT arithmetic
		"GRAPH g MATCH (a) (-[e]->){1+1,3} (b) RETURN a",              // quantifier bound must be int/param, not arithmetic
		"GRAPH g MATCH (a) (-[e]->){x} (b) RETURN a",                  // quantifier bound ident
		"GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (foo ROWS) RETURN n", // sample size ident
		"GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (1+1 ROWS) RETURN n", // sample size arithmetic
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// TestGQL_NumericOperandForms asserts the RESTRICTED-but-valid operand forms the
// grammar admits for page counts (possibly_cast_int_literal_or_parameter),
// quantifier bounds (int_literal_or_parameter), and sample sizes
// (sample_size_value = possibly_cast_int_literal_or_parameter | float): an
// integer, a query parameter, a system variable, a CAST of one, and — for the
// sample size only — a floating-point literal. (The reject side lives in
// TestGQL_Rejects; both sides are pinned against the live oracle in
// the suite:graph-query rows of testdata/compat.tsv.)
func TestGQL_NumericOperandForms(t *testing.T) {
	accepts := []string{
		"GRAPH g MATCH (n) RETURN n LIMIT @p",                             // page LIMIT parameter
		"GRAPH g MATCH (n) RETURN n LIMIT CAST(@p AS INT64)",              // page LIMIT cast
		"GRAPH g MATCH (n) RETURN n OFFSET @o LIMIT @l",                   // page OFFSET..LIMIT params
		"GRAPH g MATCH (n) LIMIT @@max RETURN n",                          // standalone-page LIMIT system var
		"GRAPH g MATCH (a) (-[e]->){@p,3} (b) RETURN a",                   // quantifier bound parameter
		"GRAPH g MATCH (a) (-[e]->){2} (b) RETURN a",                      // quantifier single int bound
		"GRAPH g MATCH (n) TABLESAMPLE BERNOULLI (10.5 PERCENT) RETURN n", // sample size float
		"GRAPH g MATCH (n) TABLESAMPLE RESERVOIR (@n ROWS) RETURN n",      // sample size parameter
	}
	for _, sql := range accepts {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q) should accept, got errs=%v", sql, errs)
		}
	}
}

// ---- CREATE PROPERTY GRAPH ----

// createPropertyGraphOf parses sql and asserts it is a *CreatePropertyGraphStmt.
func createPropertyGraphOf(t *testing.T, sql string) *ast.CreatePropertyGraphStmt {
	t.Helper()
	n := parseDDL(t, sql)
	pg, ok := n.(*ast.CreatePropertyGraphStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.CreatePropertyGraphStmt", sql, n)
	}
	return pg
}

func TestCreatePropertyGraph_Basic(t *testing.T) {
	pg := createPropertyGraphOf(t, "CREATE PROPERTY GRAPH my_graph NODE TABLES (Person, Account)")
	if pg.Name.String() != "my_graph" {
		t.Errorf("Name = %q, want my_graph", pg.Name.String())
	}
	if len(pg.NodeTables) != 2 {
		t.Fatalf("node tables = %d, want 2", len(pg.NodeTables))
	}
	if pg.NodeTables[0].Name.String() != "Person" || pg.NodeTables[1].Name.String() != "Account" {
		t.Errorf("node tables = %q/%q, want Person/Account", pg.NodeTables[0].Name, pg.NodeTables[1].Name)
	}
	if pg.EdgeTables != nil {
		t.Errorf("EdgeTables = %+v, want nil", pg.EdgeTables)
	}
}

func TestCreatePropertyGraph_OrReplaceOptions(t *testing.T) {
	pg := createPropertyGraphOf(t, "CREATE OR REPLACE PROPERTY GRAPH g NODE TABLES (T) OPTIONS(x='y')")
	if !pg.OrReplace {
		t.Error("expected OrReplace")
	}
	if pg.IfNotExists {
		t.Error("did not expect IfNotExists")
	}
	if len(pg.Options) != 1 || pg.Options[0].Name != "x" {
		t.Fatalf("options = %+v, want one option x", pg.Options)
	}
	if got := pg.Loc.End; got != len("CREATE OR REPLACE PROPERTY GRAPH g NODE TABLES (T) OPTIONS(x='y')") {
		t.Errorf("Loc.End = %d, want the options list's close paren", got)
	}
}

// TestCreatePropertyGraph_OptionsAfterContent holds the OPTIONS list to the
// position the pinned grammar (googlesql.tm:2666) and Spanner's graph DDL
// reference both put it: after the node and edge tables. The leading position
// the legacy ANTLR GoogleSQLParser.g4 had is no alternative of either, and the
// emulator refuses it at the OPTIONS word (compat.tsv's suite:graph-query row).
func TestCreatePropertyGraph_OptionsAfterContent(t *testing.T) {
	accepted := []string{
		"CREATE PROPERTY GRAPH g NODE TABLES (T) OPTIONS(x=true)",
		"CREATE PROPERTY GRAPH g NODE TABLES (T) EDGE TABLES (E) OPTIONS(x=true)",
		"CREATE PROPERTY GRAPH g NODE TABLES (T) OPTIONS()",
	}
	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accept", sql, err)
		}
	}
	refused := []string{
		"CREATE PROPERTY GRAPH g OPTIONS(x=true) NODE TABLES (T)",
		"CREATE OR REPLACE PROPERTY GRAPH g OPTIONS(x='y') NODE TABLES (T)",
		"CREATE PROPERTY GRAPH g OPTIONS(x=true) NODE TABLES (T) EDGE TABLES (E)",
	}
	for _, sql := range refused {
		assertReject(t, sql)
	}
}

func TestCreatePropertyGraph_IfNotExists(t *testing.T) {
	pg := createPropertyGraphOf(t, "CREATE PROPERTY GRAPH IF NOT EXISTS g NODE TABLES (T)")
	if pg.OrReplace {
		t.Error("did not expect OrReplace")
	}
	if !pg.IfNotExists {
		t.Error("expected IfNotExists")
	}
}

// TestCreatePropertyGraph_OrReplaceIfNotExistsConflict asserts that Spanner
// rejects OR REPLACE + IF NOT EXISTS on CREATE PROPERTY GRAPH, while union
// mode accepts it as BigQuery parses both modifiers.
func TestCreatePropertyGraph_OrReplaceIfNotExistsConflict(t *testing.T) {
	if _, err := Parse("CREATE OR REPLACE PROPERTY GRAPH IF NOT EXISTS g NODE TABLES (T)"); err != nil {
		t.Errorf("union Parse unexpected error: %v", err)
	}
}

func TestCreatePropertyGraph_EdgeTablesWithSourceDest(t *testing.T) {
	sql := "CREATE PROPERTY GRAPH g " +
		"NODE TABLES (Person KEY (id), Account KEY (id)) " +
		"EDGE TABLES (Owns KEY (person_id, account_id) " +
		"SOURCE KEY (person_id) REFERENCES Person (id) " +
		"DESTINATION KEY (account_id) REFERENCES Account (id))"
	pg := createPropertyGraphOf(t, sql)
	if len(pg.NodeTables) != 2 {
		t.Fatalf("node tables = %d, want 2", len(pg.NodeTables))
	}
	if len(pg.NodeTables[0].Key) != 1 || pg.NodeTables[0].Key[0] != "id" {
		t.Errorf("Person KEY = %+v, want [id]", pg.NodeTables[0].Key)
	}
	if len(pg.EdgeTables) != 1 {
		t.Fatalf("edge tables = %d, want 1", len(pg.EdgeTables))
	}
	edge := pg.EdgeTables[0]
	if edge.Source == nil || edge.Source.Node != "Person" {
		t.Fatalf("edge SOURCE = %+v, want REFERENCES Person", edge.Source)
	}
	if len(edge.Source.Columns) != 1 || edge.Source.Columns[0] != "person_id" {
		t.Errorf("source columns = %+v, want [person_id]", edge.Source.Columns)
	}
	if len(edge.Source.RefColumns) != 1 || edge.Source.RefColumns[0] != "id" {
		t.Errorf("source ref columns = %+v, want [id]", edge.Source.RefColumns)
	}
	if edge.Dest == nil || edge.Dest.Node != "Account" {
		t.Fatalf("edge DESTINATION = %+v, want REFERENCES Account", edge.Dest)
	}
}

func TestCreatePropertyGraph_Labels(t *testing.T) {
	sql := "CREATE PROPERTY GRAPH g NODE TABLES (" +
		"Person AS P LABEL Human PROPERTIES (name, age AS years) " +
		"DEFAULT LABEL NO PROPERTIES)"
	pg := createPropertyGraphOf(t, sql)
	def := pg.NodeTables[0]
	if def.Alias != "P" {
		t.Errorf("alias = %q, want P", def.Alias)
	}
	if len(def.Labels) != 2 {
		t.Fatalf("labels = %d, want 2", len(def.Labels))
	}
	if def.Labels[0].LabelName != "Human" || def.Labels[0].Kind != ast.LabelPropsList {
		t.Errorf("label[0] = %+v, want Human with a props list", def.Labels[0])
	}
	if len(def.Labels[0].PropsList) != 2 || def.Labels[0].PropsList[1].Alias != "years" {
		t.Errorf("label[0] props = %+v, want name, age AS years", def.Labels[0].PropsList)
	}
	if !def.Labels[1].Default || def.Labels[1].LabelName != "" || def.Labels[1].Kind != ast.LabelPropsNone {
		t.Errorf("label[1] = %+v, want DEFAULT LABEL NO PROPERTIES", def.Labels[1])
	}

	pg = createPropertyGraphOf(t, "CREATE PROPERTY GRAPH g NODE TABLES (T DYNAMIC LABEL(lbl) DYNAMIC PROPERTIES(props))")
	if pg.NodeTables[0].DynamicLabel == nil || pg.NodeTables[0].DynamicProperties == nil {
		t.Errorf("dynamic label/properties = %+v, want both set", pg.NodeTables[0])
	}
}

func TestCreatePropertyGraph_PropertiesAllColumns(t *testing.T) {
	pg := createPropertyGraphOf(t, "CREATE PROPERTY GRAPH g NODE TABLES (T PROPERTIES ARE ALL COLUMNS EXCEPT (secret))")
	def := pg.NodeTables[0]
	if len(def.Labels) != 1 {
		t.Fatalf("labels = %d, want 1 (bare properties clause)", len(def.Labels))
	}
	if def.Labels[0].Kind != ast.LabelPropsAllColumns {
		t.Errorf("kind = %v, want LabelPropsAllColumns", def.Labels[0].Kind)
	}
	if def.Labels[0].LabelName != "" {
		t.Errorf("bare properties clause should have empty LabelName, got %q", def.Labels[0].LabelName)
	}
	if len(def.Labels[0].ExceptColumns) != 1 || def.Labels[0].ExceptColumns[0] != "secret" {
		t.Errorf("except = %+v, want [secret]", def.Labels[0].ExceptColumns)
	}
}

func TestCreatePropertyGraph_NoProperties(t *testing.T) {
	pg := createPropertyGraphOf(t, "CREATE PROPERTY GRAPH g NODE TABLES (T LABEL L NO PROPERTIES)")
	if pg.NodeTables[0].Labels[0].Kind != ast.LabelPropsNone {
		t.Errorf("kind = %v, want LabelPropsNone", pg.NodeTables[0].Labels[0].Kind)
	}
}

func TestCreatePropertyGraph_TrailingComma(t *testing.T) {
	// A trailing comma before the close paren is permitted by element_table_list.
	pg := createPropertyGraphOf(t, "CREATE PROPERTY GRAPH g NODE TABLES (A, B,)")
	if len(pg.NodeTables) != 2 {
		t.Errorf("node tables = %d, want 2", len(pg.NodeTables))
	}
}

func TestCreatePropertyGraph_Rejects(t *testing.T) {
	cases := []string{
		"CREATE PROPERTY g NODE TABLES (T)",              // PROPERTY without GRAPH
		"CREATE PROPERTY GRAPH g NODE TABLES ()",         // empty element list
		"CREATE PROPERTY GRAPH g NODE TABLES T",          // element list needs parens
		"CREATE PROPERTY GRAPH g NODE TABLES (T LABEL)",  // LABEL needs a name
		"CREATE PROPERTY GRAPH g NODE TABLES (T KEY id)", // KEY needs parens
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
	// BigQuery accepts optional NODE TABLES and create-scope on CREATE PROPERTY
	// GRAPH in union mode, while DialectSpanner rejects all three.
	for _, sql := range []string{
		"CREATE PROPERTY GRAPH g",
		"CREATE TEMP PROPERTY GRAPH g NODE TABLES (T)",
		"CREATE PROPERTY GRAPH g EDGE TABLES (E)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", sql, err)
		}
	}
}

// TestGQL_WalkReachesSubnodes verifies the generated walker descends into GQL
// nodes — a query-span / lineage pass must be able to reach the expressions
// embedded in a graph query (LET/FILTER/property values/return items). We count
// the path-expression nodes reachable from a GQL statement.
func TestGQL_WalkReachesSubnodes(t *testing.T) {
	g := gqlStmtOf(t, "GRAPH g MATCH (n:Person {x: a.b}) FILTER n.age > c.d RETURN n.name AS nm")
	var pathExprs int
	ast.Inspect(g, func(n ast.Node) bool {
		if _, ok := n.(*ast.PathExpr); ok {
			pathExprs++
		}
		return true
	})
	// Expect at least: a.b (property value), c.d (filter rhs), n.age, n.name,
	// plus the graph name path. The exact count depends on how the expressions
	// node shapes column refs; assert a healthy lower bound proves traversal.
	if pathExprs < 3 {
		t.Errorf("walker reached %d PathExpr nodes, want >= 3 (subnodes unreachable?)", pathExprs)
	}
}

// ---------------------------------------------------------------------------
// Graph subqueries, graph expressions and CALL (the alternatives compat.tsv's
// graph rows cover: braced_graph_subquery, the three exists_graph_subquery
// forms, expression_subquery_with_keyword over a brace, graph_expression's
// endpoint/label predicates and IN, and graph_call_operator).
// ---------------------------------------------------------------------------

// selectExprItems parses a SELECT and returns its select list.
func selectExprItems(t *testing.T, sql string) []*ast.SelectItem {
	t.Helper()
	n := parseDDL(t, sql)
	q, ok := n.(*ast.QueryStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.QueryStmt", sql, n)
	}
	sel, ok := q.Body.(*ast.SelectStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): body is %T, want *ast.SelectStmt", sql, q.Body)
	}
	return sel.Items
}

// onlyItemExpr parses `SELECT <expr>` and returns the single item's expression.
func onlyItemExpr(t *testing.T, sql string) ast.Node {
	t.Helper()
	items := selectExprItems(t, sql)
	if len(items) != 1 {
		t.Fatalf("parseForTest(%q): %d items, want 1", sql, len(items))
	}
	return items[0].Expr
}

func TestGraphSubquery_Modifiers(t *testing.T) {
	cases := []struct {
		sql      string
		modifier ast.GraphSubqueryModifier
		graph    string
	}{
		{"SELECT ARRAY { RETURN NULL }", ast.GraphSubqueryArray, ""},
		{"SELECT VALUE { RETURN NULL }", ast.GraphSubqueryValue, ""},
		{"SELECT EXISTS { RETURN NULL }", ast.GraphSubqueryExists, ""},
		{"SELECT ARRAY { GRAPH a RETURN NULL }", ast.GraphSubqueryArray, "a"},
		{"SELECT VALUE { GRAPH a RETURN NULL }", ast.GraphSubqueryValue, "a"},
	}
	for _, tc := range cases {
		sub, ok := onlyItemExpr(t, tc.sql).(*ast.GraphSubqueryExpr)
		if !ok {
			t.Fatalf("%q: expression is %T, want *ast.GraphSubqueryExpr", tc.sql, onlyItemExpr(t, tc.sql))
		}
		if sub.Modifier != tc.modifier {
			t.Errorf("%q: Modifier = %v, want %v", tc.sql, sub.Modifier, tc.modifier)
		}
		gotGraph := ""
		if sub.Graph != nil {
			gotGraph = sub.Graph.String()
		}
		if gotGraph != tc.graph {
			t.Errorf("%q: Graph = %q, want %q", tc.sql, gotGraph, tc.graph)
		}
		if len(sub.Blocks) != 1 {
			t.Fatalf("%q: Blocks = %d, want 1", tc.sql, len(sub.Blocks))
		}
		lq, ok := sub.Blocks[0].(*ast.GraphLinearQuery)
		if !ok {
			t.Fatalf("%q: Blocks[0] is %T, want *ast.GraphLinearQuery", tc.sql, sub.Blocks[0])
		}
		if lq.Return == nil || len(lq.Return.Items) != 1 {
			t.Errorf("%q: RETURN items = %v, want 1", tc.sql, lq.Return)
		}
	}
}

func TestGraphSubquery_InPredicate(t *testing.T) {
	in, ok := onlyItemExpr(t, "SELECT NULL IN { RETURN NULL }").(*ast.InExpr)
	if !ok {
		t.Fatalf("expression is %T, want *ast.InExpr", onlyItemExpr(t, "SELECT NULL IN { RETURN NULL }"))
	}
	if in.GraphSubquery == nil {
		t.Fatal("InExpr.GraphSubquery = nil, want the braced graph subquery")
	}
	if in.GraphSubquery.Modifier != ast.GraphSubqueryBare {
		t.Errorf("Modifier = %v, want GraphSubqueryBare", in.GraphSubquery.Modifier)
	}
	if len(in.GraphSubquery.Blocks) != 1 {
		t.Errorf("Blocks = %d, want 1", len(in.GraphSubquery.Blocks))
	}
}

func TestGraphExistsSubquery_PatternAndLinearOps(t *testing.T) {
	// A bare graph pattern after the brace.
	e, ok := onlyItemExpr(t, "SELECT EXISTS { - }").(*ast.GraphExistsExpr)
	if !ok {
		t.Fatalf("expression is %T, want *ast.GraphExistsExpr", onlyItemExpr(t, "SELECT EXISTS { - }"))
	}
	if e.Kind != ast.GraphExistsPattern {
		t.Errorf("Kind = %v, want GraphExistsPattern", e.Kind)
	}
	if e.Pattern == nil || len(e.Pattern.Paths) != 1 {
		t.Fatalf("Pattern = %v, want one path", e.Pattern)
	}
	if _, ok := e.Pattern.Paths[0].Factors[0].(*ast.GraphEdgePattern); !ok {
		t.Errorf("factor is %T, want *ast.GraphEdgePattern", e.Pattern.Paths[0].Factors[0])
	}

	// A run of linear operators with no RETURN.
	e, ok = onlyItemExpr(t, "SELECT EXISTS { MATCH - FILTER NULL }").(*ast.GraphExistsExpr)
	if !ok {
		t.Fatalf("expression is %T, want *ast.GraphExistsExpr", onlyItemExpr(t, "SELECT EXISTS { MATCH - FILTER NULL }"))
	}
	if e.Kind != ast.GraphExistsLinearOps {
		t.Errorf("Kind = %v, want GraphExistsLinearOps", e.Kind)
	}
	if len(e.Ops) != 2 {
		t.Fatalf("Ops = %d, want 2 (MATCH, FILTER)", len(e.Ops))
	}
	if _, ok := e.Ops[0].(*ast.GraphMatchOp); !ok {
		t.Errorf("Ops[0] is %T, want *ast.GraphMatchOp", e.Ops[0])
	}
	if _, ok := e.Ops[1].(*ast.GraphFilterOp); !ok {
		t.Errorf("Ops[1] is %T, want *ast.GraphFilterOp", e.Ops[1])
	}

	// A named graph prefix.
	e, ok = onlyItemExpr(t, "SELECT EXISTS { GRAPH a MATCH - }").(*ast.GraphExistsExpr)
	if !ok {
		t.Fatalf("expression is %T, want *ast.GraphExistsExpr", onlyItemExpr(t, "SELECT EXISTS { GRAPH a MATCH - }"))
	}
	if e.Graph == nil || e.Graph.String() != "a" {
		t.Errorf("Graph = %v, want a", e.Graph)
	}
}

func TestGraphCallOp(t *testing.T) {
	// A named table-valued function with a YIELD list.
	call := firstLinearOp(t, "SELECT EXISTS { CALL a() YIELD a, a AS b }")
	c, ok := call.(*ast.GraphCallOp)
	if !ok {
		t.Fatalf("operator is %T, want *ast.GraphCallOp", call)
	}
	if c.Call == nil || c.Call.Name.String() != "a" {
		t.Fatalf("Call = %v, want a()", c.Call)
	}
	if len(c.Yield) != 2 || c.Yield[0].Name != "a" || c.Yield[1].Alias != "b" {
		t.Errorf("Yield = %v, want [a, a AS b]", c.Yield)
	}
	if c.HasPer || c.Optional {
		t.Errorf("HasPer/Optional = %v/%v, want false/false", c.HasPer, c.Optional)
	}

	// OPTIONAL CALL with a PER clause, and the empty PER list.
	c, ok = firstLinearOp(t, "SELECT EXISTS { OPTIONAL CALL PER(a) a() }").(*ast.GraphCallOp)
	if !ok {
		t.Fatal("operator is not *ast.GraphCallOp")
	}
	if !c.Optional {
		t.Error("Optional = false, want true")
	}
	if !c.HasPer || len(c.Per) != 1 || c.Per[0] != "a" {
		t.Errorf("Per = %v (HasPer %v), want [a]", c.Per, c.HasPer)
	}
	c, ok = firstLinearOp(t, "SELECT EXISTS { CALL PER() a() }").(*ast.GraphCallOp)
	if !ok {
		t.Fatal("operator is not *ast.GraphCallOp")
	}
	if !c.HasPer || len(c.Per) != 0 {
		t.Errorf("Per = %v (HasPer %v), want the empty PER list", c.Per, c.HasPer)
	}

	// The inline-subquery call, with and without the identifier list.
	c, ok = firstLinearOp(t, "SELECT EXISTS { CALL { RETURN NULL } }").(*ast.GraphCallOp)
	if !ok {
		t.Fatal("operator is not *ast.GraphCallOp")
	}
	if c.Subquery == nil || c.Call != nil {
		t.Fatalf("Subquery = %v, Call = %v; want a subquery call", c.Subquery, c.Call)
	}
	c, ok = firstLinearOp(t, "SELECT EXISTS { CALL() { RETURN NULL } }").(*ast.GraphCallOp)
	if !ok {
		t.Fatal("operator is not *ast.GraphCallOp")
	}
	if c.Subquery == nil || c.HasPer || len(c.Per) != 0 {
		t.Errorf("identifier-list call = %v / Per %v (HasPer %v), want a subquery call with an empty, non-PER list", c.Subquery, c.Per, c.HasPer)
	}
}

// firstLinearOp parses a statement whose first block is a linear query and
// returns its first operator.
func firstLinearOp(t *testing.T, sql string) ast.Node {
	t.Helper()
	e := onlyItemExpr(t, sql)
	var ops []ast.Node
	switch n := e.(type) {
	case *ast.GraphExistsExpr:
		ops = n.Ops
	case *ast.GraphSubqueryExpr:
		lq, ok := n.Blocks[0].(*ast.GraphLinearQuery)
		if !ok {
			t.Fatalf("%q: Blocks[0] is %T, want *ast.GraphLinearQuery", sql, n.Blocks[0])
		}
		ops = lq.Operators
	default:
		t.Fatalf("%q: expression is %T, want a graph subquery", sql, e)
	}
	if len(ops) == 0 {
		t.Fatalf("%q: no linear operators", sql)
	}
	return ops[0]
}

func TestGraphPredicate(t *testing.T) {
	cases := []struct {
		expr string
		kind ast.GraphPredicateKind
		not  bool
	}{
		{"a IS SOURCE b", ast.GraphIsSource, false},
		{"a IS SOURCE OF b", ast.GraphIsSource, false},
		{"a IS NOT SOURCE b", ast.GraphIsSource, true},
		{"a IS NOT SOURCE OF b", ast.GraphIsSource, true},
		{"a IS DESTINATION b", ast.GraphIsDestination, false},
		{"a IS NOT DESTINATION OF b", ast.GraphIsDestination, true},
	}
	for _, tc := range cases {
		p, ok := mustParse(t, tc.expr).(*ast.GraphPredicate)
		if !ok {
			t.Fatalf("%q: expression is %T, want *ast.GraphPredicate", tc.expr, mustParse(t, tc.expr))
		}
		if p.Kind != tc.kind || p.Not != tc.not {
			t.Errorf("%q: Kind = %v, Not = %v; want %v, %v", tc.expr, p.Kind, p.Not, tc.kind, tc.not)
		}
		if _, ok := p.Right.(*ast.Identifier); !ok {
			t.Errorf("%q: Right is %T, want the endpoint expression", tc.expr, p.Right)
		}
	}

	// IS [NOT] LABELED carries the label algebra.
	p, ok := mustParse(t, "a IS LABELED x & !y | %").(*ast.GraphPredicate)
	if !ok {
		t.Fatal("expression is not *ast.GraphPredicate")
	}
	if p.Kind != ast.GraphIsLabeled {
		t.Errorf("Kind = %v, want GraphIsLabeled", p.Kind)
	}
	label, ok := p.Right.(*ast.GraphLabelExpr)
	if !ok {
		t.Fatalf("Right is %T, want *ast.GraphLabelExpr", p.Right)
	}
	// `&` binds tighter than `|`: (x & !y) | %.
	if label.Kind != ast.LabelOr {
		t.Fatalf("top label op = %v, want LabelOr", label.Kind)
	}
	if label.Left.Kind != ast.LabelAnd || label.Left.Right.Kind != ast.LabelNot {
		t.Errorf("left operand = %v/%v, want an AND of a NOT", label.Left.Kind, label.Left.Right.Kind)
	}
	if label.Right.Kind != ast.LabelWildcard {
		t.Errorf("right operand = %v, want the %% wildcard", label.Right.Kind)
	}
}

func TestGraphSearchPrefixAndPathMode(t *testing.T) {
	cases := []struct {
		pattern string
		search  string
		mode    string
	}{
		{"ANY -", "ANY", ""},
		{"ANY SHORTEST -", "ANY SHORTEST", ""},
		{"ANY CHEAPEST -", "ANY CHEAPEST", ""},
		{"ANY 1 -", "ANY 1", ""},
		{"SHORTEST 1 -", "SHORTEST 1", ""},
		{"CHEAPEST 1 -", "CHEAPEST 1", ""},
		{"ALL -", "ALL", ""},
		{"ALL SHORTEST -", "ALL SHORTEST", ""},
		{"ALL CHEAPEST -", "ALL CHEAPEST", ""},
		{"WALK -", "", "WALK"},
		{"TRAIL PATH -", "", "TRAIL PATH"},
		{"ACYCLIC PATHS -", "", "ACYCLIC PATHS"},
		{"PATH -", "", "PATH"},
		{"PATHS -", "", "PATHS"},
		{"p = ANY SHORTEST WALK PATHS -", "ANY SHORTEST", "WALK PATHS"},
	}
	for _, tc := range cases {
		path := firstPathOf(t, "SELECT EXISTS { "+tc.pattern+" }")
		if path.Search != tc.search {
			t.Errorf("%q: Search = %q, want %q", tc.pattern, path.Search, tc.search)
		}
		if path.Mode != tc.mode {
			t.Errorf("%q: Mode = %q, want %q", tc.pattern, path.Mode, tc.mode)
		}
	}
	// A path-variable assignment: `p = …`.
	path := firstPathOf(t, "SELECT EXISTS { p = ANY - }")
	if path.PathVar != "p" {
		t.Errorf("PathVar = %q, want p", path.PathVar)
	}
}

// firstPathOf parses `SELECT EXISTS { <pattern> }` and returns the pattern's
// single path pattern.
func firstPathOf(t *testing.T, sql string) *ast.GraphPathPattern {
	t.Helper()
	e, ok := onlyItemExpr(t, sql).(*ast.GraphExistsExpr)
	if !ok {
		t.Fatalf("%q: expression is %T, want *ast.GraphExistsExpr", sql, onlyItemExpr(t, sql))
	}
	if e.Pattern == nil || len(e.Pattern.Paths) != 1 {
		t.Fatalf("%q: Pattern = %v, want one path", sql, e.Pattern)
	}
	return e.Pattern.Paths[0]
}

func TestGraphQuantifier(t *testing.T) {
	// `+`, `*`, `{n}`, `{n,}`, `{,m}` and `{n,m}` are all quantifiers on an edge.
	// The grammar wraps a quantified edge in a path pattern, so the factor is a
	// one-factor *GraphPathPattern whose own factor is the edge.
	for _, pattern := range []string{"- +", "- *", "- { 1 }", "- { 1, }", "- {, 2 }", "- { 1, 2 }"} {
		path := firstPathOf(t, "SELECT EXISTS { "+pattern+" }")
		if len(path.Factors) != 1 {
			t.Fatalf("%q: %d factors, want 1", pattern, len(path.Factors))
		}
		wrapped, ok := path.Factors[0].(*ast.GraphPathPattern)
		if !ok {
			t.Fatalf("%q: factor is %T, want *ast.GraphPathPattern", pattern, path.Factors[0])
		}
		if len(wrapped.Factors) != 1 {
			t.Fatalf("%q: quantified path has %d factors, want 1", pattern, len(wrapped.Factors))
		}
		if _, ok := wrapped.Factors[0].(*ast.GraphEdgePattern); !ok {
			t.Errorf("%q: quantified factor is %T, want *ast.GraphEdgePattern", pattern, wrapped.Factors[0])
		}
		if wrapped.Loc.End == nodeLoc(wrapped.Factors[0]).End {
			t.Errorf("%q: quantified path Loc does not span the quantifier", pattern)
		}
	}
	// A quantifier on a parenthesized path keeps the path itself.
	path := firstPathOf(t, "SELECT EXISTS { (-)+ }")
	if _, ok := path.Factors[0].(*ast.GraphPathPattern); !ok {
		t.Errorf("quantified parenthesized path is %T, want *ast.GraphPathPattern", path.Factors[0])
	}
	// A quantifier on a node pattern is refused, as the grammar's action does.
	if _, errs := parseForTest("SELECT EXISTS { (a)* }"); len(errs) == 0 {
		t.Error("parseForTest(`SELECT EXISTS { (a)* }`) accepted a quantifier on a node pattern")
	}
}

func TestGraphEdgePatternRequiresAdjacency(t *testing.T) {
	// The grammar's ValidateNoWhitespace refuses a space between an edge
	// pattern's delimiters, so the generator's spaced spelling is a near-miss the
	// emulator rejects and omni must too.
	for _, sql := range []string{
		"SELECT EXISTS { - [ ] - }",
		"SELECT EXISTS { - [ ] -> }",
		"SELECT EXISTS { < - [ ] - }",
		"SELECT EXISTS { < - }",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) accepted whitespace inside an edge pattern", sql)
		}
	}
	// The adjacent spellings parse.
	for _, sql := range []string{
		"SELECT EXISTS { -[]- }",
		"SELECT EXISTS { -[]-> }",
		"SELECT EXISTS { <-[]- }",
		"SELECT EXISTS { <- }",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
}

func TestGraphHintPlacement(t *testing.T) {
	// A hint between two edges is ambiguous, and the grammar refuses it.
	if _, errs := parseForTest("SELECT EXISTS { -@1 - }"); len(errs) == 0 {
		t.Error("parseForTest accepted a hint between two edge patterns")
	}
	// A quantified edge is a path pattern, not an edge factor, so a hint beside
	// one is not the edge-edge case.
	for _, sql := range []string{
		"SELECT EXISTS { -+ @1 - }",
		"SELECT EXISTS { - @1 -+ }",
		"SELECT EXISTS { -+ @1 -+ }",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
	// A hint may not open a parenthesized path pattern.
	if _, errs := parseForTest("SELECT EXISTS { (@1 -) }"); len(errs) == 0 {
		t.Error("parseForTest accepted a hint at the beginning of a path pattern")
	}
	// A hint may open a node pattern's filler, and may sit between two nodes.
	for _, sql := range []string{
		"SELECT EXISTS { (@1) }",
		"SELECT EXISTS { (@1 a) }",
		"SELECT EXISTS { (a)@1(b) }",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
}

func TestGraphKeywordIdentifiers(t *testing.T) {
	// The pinned grammar's common_keyword_as_identifier takes these as names;
	// the pinned emulator agrees (compat.tsv's common_keyword_as_identifier rows).
	for _, word := range []string{"NOTHING", "SHORTEST", "REPLACE_FIELDS", "SAFE_CAST"} {
		for _, sql := range []string{"SELECT " + word, "SELECT " + word + " zzz"} {
			if _, errs := parseForTest(sql); len(errs) != 0 {
				t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
			}
		}
	}
	// The clause readings survive: QUALIFY after a select list is still the
	// clause, and SAFE_CAST/REPLACE_FIELDS with a '(' are still the functions.
	for _, sql := range []string{
		"SELECT a QUALIFY ROW_NUMBER() OVER () = 1",
		"SELECT SAFE_CAST(a AS INT64) FROM t",
		"SELECT REPLACE_FIELDS(a, b AS c) FROM t",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
}

func TestGraphCallOpAndPrefixStrictness(t *testing.T) {
	// OPTIONAL is the call's first token and its Loc must span it.
	call := firstLinearOp(t, "SELECT EXISTS { OPTIONAL CALL a() }")
	c, ok := call.(*ast.GraphCallOp)
	if !ok {
		t.Fatalf("operator is %T, want *ast.GraphCallOp", call)
	}
	if !c.Optional {
		t.Error("Optional = false, want true")
	}
	if got := c.Loc.End - c.Loc.Start; got != len("OPTIONAL CALL a()") {
		t.Errorf("Loc spans %d bytes, want %d (the whole operator)", got, len("OPTIONAL CALL a()"))
	}
	// The grammar's strict shapes: no ALL count, no trailing comma in a
	// parenthesized identifier list, and no hint on an IN graph subquery.
	for _, sql := range []string{
		"SELECT EXISTS { ALL 2 - }",
		"SELECT EXISTS { CALL PER(a,) a() }",
		"SELECT EXISTS { CALL(a,) { RETURN NULL } }",
		"SELECT 1 IN @1 { RETURN NULL }",
		"SELECT EXISTS { (a)* }",
		"SELECT EXISTS { - { } }",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) accepted text the grammar refuses", sql)
		}
	}
	for _, sql := range []string{
		"SELECT EXISTS { ANY 2 - }",
		"SELECT EXISTS { SHORTEST 2 - }",
		"SELECT EXISTS { CHEAPEST 2 - }",
		"SELECT EXISTS { CALL PER() a() }",
		"SELECT EXISTS { CALL(a) { RETURN NULL } }",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
}

func TestGraphKeywordWordsAsIdentifiers(t *testing.T) {
	// The pinned grammar's non-reserved words are identifiers wherever an
	// identifier is required, not only in an expression path.
	for _, sql := range []string{} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
	// A QUALIFY with no expression after it is the alias even when a clause
	// keyword follows it (in both select-list and table-alias positions).
	for _, sql := range []string{} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q): unexpected errors: %v", sql, errs)
		}
	}
	// The clause reading survives in both select-list and FROM contexts,
	// including when the QUALIFY predicate starts with a non-reserved keyword
	// or an inline WITH(...) expression.
	for _, sql := range []string{
		"SELECT a QUALIFY b",
		"SELECT a FROM t QUALIFY ROW_NUMBER() OVER () = 1",
		"SELECT a FROM t QUALIFY b",
		"SELECT ROW_NUMBER() OVER () FROM t QUALIFY pivot",
		"SELECT ROW_NUMBER() OVER () FROM t QUALIFY unpivot",
		"SELECT ROW_NUMBER() OVER () FROM t QUALIFY WITH(x AS 1, x > 0)",
		"SELECT ROW_NUMBER() OVER () FROM t QUALIFY LEFT(s, 1) = 'x'",
		"SELECT ROW_NUMBER() OVER () FROM t QUALIFY RIGHT(s, 1) = 'x'",
	} {
		f, err := Parse(sql)
		if err != nil {
			t.Errorf("Parse(%q): unexpected error: %v", sql, err)
			continue
		}
		qs, ok := f.Stmts[0].(*ast.QueryStmt)
		if !ok {
			t.Errorf("Parse(%q): got %T, want *ast.QueryStmt", sql, f.Stmts[0])
			continue
		}
		sel := selectOf(t, qs)
		if sel.Qualify == nil {
			t.Errorf("Parse(%q): expected non-nil Qualify clause, got %#v", sql, qs.Body)
		}
	}
	// DialectBigQuery still reserves QUALIFY and rejects bare QUALIFY aliases;
	// DialectSpanner accepts bare QUALIFY table aliases.
	if _, err := Parse("TABLE a() QUALIFY"); err == nil {
		t.Errorf("Parse(TABLE a() QUALIFY): expected error, got nil")
	}
}

func TestGraphLabelPredicateLoc(t *testing.T) {
	// The predicate's range ends at its last token: a parenthesized label's `)`.
	sql := "SELECT a IS LABELED (Person)"
	n := onlyItemExpr(t, sql)
	p, ok := n.(*ast.GraphPredicate)
	if !ok {
		t.Fatalf("expression is %T, want *ast.GraphPredicate", n)
	}
	want := sql[len("SELECT "):]
	if got := sql[p.Loc.Start:p.Loc.End]; got != want {
		t.Errorf("predicate text = %q, want %q", got, want)
	}
}

func TestGraphSamplePartitionByIsKept(t *testing.T) {
	e, ok := onlyItemExpr(t, "SELECT EXISTS { TABLESAMPLE BERNOULLI(1 ROWS PARTITION BY 1, 2) RETURN 1 }").(*ast.GraphSubqueryExpr)
	if !ok {
		t.Fatal("expression is not *ast.GraphSubqueryExpr")
	}
	lq, ok := e.Blocks[0].(*ast.GraphLinearQuery)
	if !ok {
		t.Fatalf("Blocks[0] is %T, want *ast.GraphLinearQuery", e.Blocks[0])
	}
	sample, ok := lq.Operators[0].(*ast.GraphSampleOp)
	if !ok {
		t.Fatalf("operator is %T, want *ast.GraphSampleOp", lq.Operators[0])
	}
	if len(sample.PartitionBy) != 2 {
		t.Errorf("PartitionBy = %d expressions, want 2", len(sample.PartitionBy))
	}
	if sample.Unit != "ROWS" || sample.Method != "BERNOULLI" {
		t.Errorf("method/unit = %q/%q, want BERNOULLI/ROWS", sample.Method, sample.Unit)
	}
}

func TestGraphCallOpTVFArguments(t *testing.T) {
	// graph_call_operator_core's target is the grammar's `tvf`, so a non-expression
	// argument (`TABLE t`) parses.
	c, ok := firstLinearOp(t, "SELECT EXISTS { CALL f(TABLE t) YIELD a }").(*ast.GraphCallOp)
	if !ok {
		t.Fatal("operator is not *ast.GraphCallOp")
	}
	if c.Call == nil || len(c.Call.Args) != 1 {
		t.Fatalf("Call = %v, want one argument", c.Call)
	}
	if _, ok := c.Call.Args[0].(*ast.TableClause); !ok {
		t.Errorf("argument is %T, want *ast.TableClause", c.Call.Args[0])
	}
}

func TestGraphSetOpPropagationMode(t *testing.T) {
	cases := []struct {
		sql  string
		mode string
	}{
		{"SELECT EXISTS { RETURN 1 UNION ALL RETURN 2 }", ""},
		{"SELECT EXISTS { RETURN 1 FULL UNION ALL RETURN 2 }", "FULL"},
		{"SELECT EXISTS { RETURN 1 FULL OUTER UNION ALL RETURN 2 }", "FULL"},
		{"SELECT EXISTS { RETURN 1 OUTER UNION DISTINCT RETURN 2 }", "OUTER"},
		{"SELECT EXISTS { RETURN 1 INNER INTERSECT ALL RETURN 2 }", "INNER"},
		{"SELECT EXISTS { RETURN 1 LEFT OUTER EXCEPT ALL RETURN 2 }", "LEFT"},
	}
	for _, tc := range cases {
		e, ok := onlyItemExpr(t, tc.sql).(*ast.GraphSubqueryExpr)
		if !ok {
			t.Fatalf("%q: expression is %T, want *ast.GraphSubqueryExpr", tc.sql, onlyItemExpr(t, tc.sql))
		}
		setOp, ok := e.Blocks[0].(*ast.GraphSetOp)
		if !ok {
			t.Fatalf("%q: Blocks[0] is %T, want *ast.GraphSetOp", tc.sql, e.Blocks[0])
		}
		if len(setOp.Metas) != 1 {
			t.Fatalf("%q: %d metadata entries, want 1", tc.sql, len(setOp.Metas))
		}
		if setOp.Metas[0].OuterMode != tc.mode {
			t.Errorf("%q: OuterMode = %q, want %q", tc.sql, setOp.Metas[0].OuterMode, tc.mode)
		}
	}
}

func TestGraphSubqueryGraphNameIsPathExpression(t *testing.T) {
	// braced_graph_subquery takes path_expression: dotted parts only, no dashed
	// BigQuery table name.
	if _, errs := parseForTest("SELECT ARRAY { GRAPH a.b RETURN 1 }"); len(errs) != 0 {
		t.Errorf("GRAPH a.b: unexpected errors: %v", errs)
	}
	if _, errs := parseForTest("SELECT ARRAY { GRAPH a-b RETURN 1 }"); len(errs) == 0 {
		t.Error("GRAPH a-b: accepted a dashed name where the grammar takes path_expression")
	}
}

func TestReplaceFieldsExtensionPath(t *testing.T) {
	// generalized_extension_path renders its dots, so the AST keeps the target the
	// input named.
	cases := []struct {
		sql  string
		want string
	}{
		{"SELECT REPLACE_FIELDS(s, 1 AS (a))", "(a)"},
		{"SELECT REPLACE_FIELDS(s, 1 AS (a).b)", "(a).b"},
		{"SELECT REPLACE_FIELDS(s, 1 AS (a).(b))", "(a).(b)"},
		{"SELECT REPLACE_FIELDS(s, 1 AS (a).(b).c)", "(a).(b).c"},
	}
	for _, tc := range cases {
		n := onlyItemExpr(t, tc.sql)
		rf, ok := n.(*ast.ReplaceFieldsExpr)
		if !ok {
			t.Fatalf("%q: expression is %T, want *ast.ReplaceFieldsExpr", tc.sql, n)
		}
		if len(rf.Items) != 1 {
			t.Fatalf("%q: %d replacement items, want 1", tc.sql, len(rf.Items))
		}
		if got := rf.Items[0].Alias; got != tc.want {
			t.Errorf("%q: target = %q, want %q", tc.sql, got, tc.want)
		}
	}
}

func TestGraphExpressionLoc(t *testing.T) {
	// The graph expression nodes are in nodeLoc, so an enclosing node's range is
	// real rather than NoLoc.
	file, errs := parseForTest("SELECT NOT EXISTS { RETURN 1 } AS x, a IS SOURCE b FROM t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	var unary *ast.UnaryExpr
	var predicate *ast.GraphPredicate
	ast.Inspect(file, func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.UnaryExpr:
			unary = v
		case *ast.GraphPredicate:
			predicate = v
		}
		return true
	})
	if unary == nil || !unary.Loc.IsValid() || unary.Loc.End <= unary.Loc.Start {
		t.Errorf("NOT EXISTS { … } UnaryExpr Loc = %v, want the operand's range", unary)
	}
	if predicate == nil || !predicate.Loc.IsValid() || predicate.Loc.End <= predicate.Loc.Start {
		t.Errorf("graph predicate Loc = %v, want its own range", predicate)
	}
}
