package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/trino/ast"
	"github.com/bytebase/omni/trino/parser"
)

func TestExplainStatement(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"SELECT 1", "EXPLAIN SELECT 1"},
		{"SELECT 1;", "EXPLAIN SELECT 1;"},
		{"DELETE FROM t", "EXPLAIN DELETE FROM t"},
		{"WITH x AS (SELECT 1) SELECT * FROM x", "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
		{"WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", "EXPLAIN WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)"},
		// A statement that already asks for a plan is rebuilt around the statement
		// it plans, so no request plans a plan.
		{"EXPLAIN SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1;", "EXPLAIN SELECT 1"},
		{"EXPLAIN (TYPE LOGICAL) SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1 ORDER BY 1", "EXPLAIN SELECT 1 ORDER BY 1"},
		{"EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x", "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
		{"EXPLAIN WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", "EXPLAIN WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)"},
		// Empty statements around the one planned are left out.
		{"; SELECT 1", "EXPLAIN SELECT 1"},
		{"; SELECT 1;", "EXPLAIN SELECT 1;"},
		{";; EXPLAIN (FORMAT TEXT) SELECT 1;", "EXPLAIN SELECT 1"},
		// An EXPLAIN whose inner statement span cannot be located runs as written
		// rather than planning a plan.
		{"EXPLAIN EXECUTE p", "EXPLAIN EXECUTE p"},
		{"; EXPLAIN EXECUTE p", "EXPLAIN EXECUTE p"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, false)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
}

func TestExplainStatementJSON(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"SELECT 1", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN SELECT 1", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN (FORMAT TEXT) SELECT 1", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN (FORMAT JSON) SELECT 1", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN (TYPE LOGICAL, FORMAT GRAPHVIZ) SELECT 1", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)", "EXPLAIN (FORMAT JSON) WITH FUNCTION f(x bigint) RETURNS bigint RETURN x + 1 SELECT f(1)"},
		{"EXPLAIN EXECUTE p", "EXPLAIN EXECUTE p"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, true)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
}

// TestExplainStatementRejects locks what never becomes an EXPLAIN: a write with
// EXPLAIN options smuggled in front of it, EXPLAIN ANALYZE (which executes what
// it plans), empty input, and multi-statement scripts.
func TestExplainStatementRejects(t *testing.T) {
	for _, statement := range []string{
		"EXPLAIN ANALYZE SELECT 1",
		"EXPLAIN ANALYZE VERBOSE SELECT 1",
		"EXPLAIN ANALYZE DELETE FROM t",
	} {
		_, err := ExplainStatement(statement, false)
		require.ErrorIsf(t, err, ErrExplainAnalyzeExecutes, "%q", statement)
		require.ErrorIsf(t, err, analysis.ErrExplainAnalyzeExecutes, "%q", statement)

		_, err = ExplainStatement(statement, true)
		require.ErrorIsf(t, err, ErrExplainAnalyzeExecutes, "%q", statement)
	}

	_, err := ExplainStatement("ANALYZE DELETE FROM t", false)
	require.Error(t, err)
	var pe *parser.ParseError
	require.ErrorAs(t, err, &pe)

	for _, empty := range []string{"", "   ", ";", "-- comment"} {
		_, err := ExplainStatement(empty, false)
		require.EqualErrorf(t, err, "EXPLAIN needs a statement to plan", "%q", empty)
	}

	_, err = ExplainStatement("SELECT 1; SELECT 2", false)
	require.EqualError(t, err, "EXPLAIN plans a single statement")
}

func TestExplainedTextUnlocatableFallback(t *testing.T) {
	cases := []struct {
		name    string
		explain *parser.ExplainStmt
		text    string
	}{
		{
			name:    "nil inner statement",
			explain: &parser.ExplainStmt{Loc: ast.Loc{Start: 0, End: 7}},
			text:    "EXPLAIN",
		},
		{
			name: "negative inner start",
			explain: &parser.ExplainStmt{
				Statement: &parser.QueryStmt{Loc: ast.Loc{Start: -1, End: 16}},
				Loc:       ast.Loc{Start: 0, End: 16},
			},
			text: "EXPLAIN SELECT 1",
		},
		{
			name: "end exceeds text length",
			explain: &parser.ExplainStmt{
				Statement: &parser.QueryStmt{Loc: ast.Loc{Start: 8, End: 16}},
				Loc:       ast.Loc{Start: 0, End: 100},
			},
			text: "EXPLAIN SELECT 1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := explainedText(tc.explain, tc.text)
			require.False(t, ok)
			require.Empty(t, got)
		})
	}
}

func TestDescribeExplain(t *testing.T) {
	tests := []struct {
		statement string
		format    string
		executed  bool
		ok        bool
	}{
		{statement: "SELECT 1"},
		{statement: "SELECT 1; EXPLAIN SELECT 1"},
		{statement: "EXPLAIN SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN (FORMAT JSON) SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN (FORMAT GRAPHVIZ) SELECT 1", format: "GRAPHVIZ", ok: true},
		{statement: "EXPLAIN (TYPE DISTRIBUTED) SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN (TYPE IO) SELECT 1", format: "JSON", ok: true},
		{statement: "EXPLAIN (TYPE IO, FORMAT TEXT) SELECT 1", format: "TEXT", ok: true},
		{statement: "EXPLAIN (TYPE VALIDATE) SELECT 1"},
		{statement: "EXPLAIN ANALYZE SELECT 1", format: "TEXT", executed: true, ok: true},
		{statement: "explain analyze verbose select 1", format: "TEXT", executed: true, ok: true},
	}
	for _, tc := range tests {
		format, executed, ok := DescribeExplain(tc.statement)
		require.Equalf(t, tc.ok, ok, "%q", tc.statement)
		require.Equalf(t, tc.format, format, "%q", tc.statement)
		require.Equalf(t, tc.executed, executed, "%q", tc.statement)
	}
}
