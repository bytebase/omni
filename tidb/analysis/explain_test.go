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
// the visualizer draws. The format is spelled unquoted because TiDB's grammar
// here takes an identifier, so a quoted one would not parse again.
func TestExplainStatementJSON(t *testing.T) {
	tests := []struct {
		statement string
		want      string
	}{
		{"SELECT 1", "EXPLAIN FORMAT=tidb_json SELECT 1"},
		{"EXPLAIN SELECT 1", "EXPLAIN FORMAT=tidb_json SELECT 1"},
		{"EXPLAIN FORMAT=brief SELECT 1", "EXPLAIN FORMAT=tidb_json SELECT 1"},
		{"EXPLAIN FORMAT=tidb_json SELECT 1", "EXPLAIN FORMAT=tidb_json SELECT 1"},
		{"EXPLAIN t", "EXPLAIN t"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, true)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
	// What this builds has to parse again: the driver validates it before running it.
	for _, tc := range tests {
		_, err := ExplainStatement(tc.want, true)
		require.NoErrorf(t, err, "%q", tc.want)
	}
	got, err := ExplainStatement("SELECT 1", false)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN SELECT 1", got)
}

// TestDescribeExplain pins what a result says about a plan the caller asked for
// itself.
func TestDescribeExplain(t *testing.T) {
	tests := []struct {
		statement string
		format    string
		executed  bool
		isExplain bool
	}{
		{statement: "EXPLAIN FORMAT=tidb_json SELECT 1", format: "json", isExplain: true},
		{statement: "explain format = TiDB_JSON select 1", format: "json", isExplain: true},
		{statement: "EXPLAIN FORMAT=brief SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN ANALYZE SELECT 1", format: "text", executed: true, isExplain: true},
		{statement: "SELECT 1"},
		// EXPLAIN over a table name describes its columns. The result is a
		// table like any other, and calling it a query plan would put it in
		// the viewer instead of the grid.
		{statement: "EXPLAIN t"},
		{statement: "DESCRIBE t"},
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

func TestExplainAlterTable(t *testing.T) {
	got, err := ExplainStatement("EXPLAIN ALTER TABLE t ADD INDEX (a)", true)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN FORMAT=tidb_json ALTER TABLE t ADD INDEX (a)", got)
	got, err = ExplainStatement("ALTER TABLE t ADD INDEX (a)", false)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN ALTER TABLE t ADD INDEX (a)", got)
	_, err = ExplainStatement("EXPLAIN ANALYZE ALTER TABLE t ADD INDEX (a)", false)
	require.ErrorIs(t, err, ErrExplainAnalyzeExecutes)
	format, executed, ok := DescribeExplain("EXPLAIN ANALYZE ALTER TABLE t ADD INDEX (a)")
	require.Equal(t, []any{"text", true, true}, []any{format, executed, ok})
}
