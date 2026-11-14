package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExplainStatement(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"SELECT 1", "EXPLAIN SELECT 1"},
		{"SELECT 1;", "EXPLAIN SELECT 1;"},
		{"DELETE FROM t", "EXPLAIN DELETE FROM t"},
		{"UPDATE t SET a = 1", "EXPLAIN UPDATE t SET a = 1"},
		// A statement that already asks for a plan is rebuilt around the statement
		// it plans, so no request plans a plan.
		{"EXPLAIN SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN FORMAT=JSON SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1 ORDER BY 1", "EXPLAIN SELECT 1 ORDER BY 1"},
		// An EXPLAIN with no statement to plan is left as the user wrote it: it
		// describes a table or another connection's query, and planning a plan
		// would only turn a working statement into a syntax error.
		{"EXPLAIN t", "EXPLAIN t"},
		{"DESCRIBE t", "DESCRIBE t"},
		{"EXPLAIN FOR CONNECTION 5", "EXPLAIN FOR CONNECTION 5"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, false)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
}

// TestExplainStatementJSON pins the statement that asks for the structured plan
// the visualizer draws. A caller's own FORMAT is replaced rather than kept, so
// the request decides the format whatever the statement named.
func TestExplainStatementJSON(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"SELECT 1", "EXPLAIN FORMAT=JSON SELECT 1"},
		{"EXPLAIN SELECT 1", "EXPLAIN FORMAT=JSON SELECT 1"},
		{"EXPLAIN FORMAT=TREE SELECT 1", "EXPLAIN FORMAT=JSON SELECT 1"},
		{"EXPLAIN FORMAT=JSON SELECT 1", "EXPLAIN FORMAT=JSON SELECT 1"},
		// Nothing to plan, so the statement runs as written and the format the
		// request named has nowhere to go.
		{"EXPLAIN t", "EXPLAIN t"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, true)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
	got, err := ExplainStatement("SELECT 1", false)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN SELECT 1", got)
}

// TestDescribeExplain pins what a result says about a plan the caller asked for
// itself. An omitted FORMAT is unknown rather than text, because from MySQL
// 8.0.32 it follows the session's explain_format.
func TestDescribeExplain(t *testing.T) {
	tests := []struct {
		statement string
		format    string
		executed  bool
		isExplain bool
	}{
		{statement: "EXPLAIN FORMAT=JSON SELECT 1", format: "json", isExplain: true},
		{statement: "explain format = json select 1", format: "json", isExplain: true},
		{statement: "EXPLAIN FORMAT=TREE SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN FORMAT=TRADITIONAL SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN SELECT 1", isExplain: true},
		{statement: "EXPLAIN ANALYZE SELECT 1", executed: true, isExplain: true},
		{statement: "SELECT 1"},
		{statement: "SELECT 1; SELECT 2"},
		// EXPLAIN over a table name describes its columns. The result is a
		// table like any other, and calling it a query plan would put it in
		// the viewer instead of the grid.
		{statement: "EXPLAIN t"},
		{statement: "DESCRIBE t"},
		{statement: "EXPLAIN FOR CONNECTION 5"},
	}
	for _, tc := range tests {
		format, executed, isExplain := DescribeExplain(tc.statement)
		require.Equalf(t, tc.isExplain, isExplain, "%q", tc.statement)
		require.Equalf(t, tc.format, format, "%q", tc.statement)
		require.Equalf(t, tc.executed, executed, "%q", tc.statement)
	}
}

// TestExplainStatementRejects locks what never becomes an EXPLAIN: a write with
// EXPLAIN options smuggled in front of it, which anything that simply prefixed
// EXPLAIN would turn into an EXPLAIN ANALYZE DELETE that runs the DELETE, and an
// EXPLAIN ANALYZE the caller wrote, which plans by executing what it plans.
func TestExplainStatementRejects(t *testing.T) {
	for _, statement := range []string{
		"ANALYZE DELETE FROM t",
		"EXPLAIN ANALYZE SELECT * FROM t",
		"EXPLAIN ANALYZE DELETE FROM t",
		"FORMAT=JSON SELECT 1",
		"SELECT 1; SELECT 2",
		"",
	} {
		_, err := ExplainStatement(statement, false)
		require.Errorf(t, err, "%q", statement)
	}
}

// TestExplainStatementExecutableComment keeps the rebuild off text whose locations
// the lexer reports against a spliced buffer. The statement runs as written, which
// still never plans a plan.
func TestExplainStatementExecutableComment(t *testing.T) {
	const statement = "EXPLAIN /*!40001 SELECT 1 */"
	got, err := ExplainStatement(statement, false)
	require.NoError(t, err)
	require.Equal(t, statement, got)
}

func TestExplainStatementNotExplainable(t *testing.T) {
	for _, statement := range []string{"ANALYZE TABLE t", "SHOW TABLES", "CREATE TABLE t (a INT)", "SET @a = 1"} {
		_, err := ExplainStatement(statement, false)
		require.ErrorIsf(t, err, ErrNotExplainable, "%q", statement)
	}
}
