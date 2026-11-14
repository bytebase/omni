package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExplainStatement(t *testing.T) {
	tests := []struct {
		statement string
		format    string
		want      string
	}{
		{"SELECT 1", "", "EXPLAIN SELECT 1"},
		{"SELECT 1", "json", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"SELECT 1", "xml", "EXPLAIN (FORMAT XML) SELECT 1"},
		{"SELECT 1", "yaml", "EXPLAIN (FORMAT YAML) SELECT 1"},
		{"DELETE FROM t", "", "EXPLAIN DELETE FROM t"},
		// An EXPLAIN is rebuilt around the statement it plans.
		{"EXPLAIN (FORMAT XML) SELECT 1", "json", "EXPLAIN (FORMAT JSON) SELECT 1"},
		{"EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x", "", "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
		{"EXPLAIN (ANALYZE false) SELECT 1", "", "EXPLAIN SELECT 1"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement, tc.format)
		require.NoError(t, err, tc.statement)
		require.Equal(t, tc.want, got, tc.statement)
	}

	for _, statement := range []string{"EXPLAIN ANALYZE SELECT 1", "EXPLAIN (ANALYZE) DELETE FROM t"} {
		_, err := ExplainStatement(statement, "")
		require.ErrorIs(t, err, ErrExplainAnalyzeExecutes)
		_, err = CockroachExplainStatement(statement, true)
		require.ErrorIs(t, err, ErrExplainAnalyzeExecutes)
	}
	_, err := ExplainStatement("SELECT 1; SELECT 2", "")
	require.EqualError(t, err, "EXPLAIN plans a single statement")
	_, err = ExplainStatement(";", "")
	require.EqualError(t, err, "EXPLAIN needs a statement to plan")
}

func TestCockroachExplainStatement(t *testing.T) {
	got, err := CockroachExplainStatement("SELECT 1", true)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN (DISTSQL) SELECT 1", got)
	got, err = CockroachExplainStatement("EXPLAIN SELECT 1", false)
	require.NoError(t, err)
	require.Equal(t, "EXPLAIN SELECT 1", got)
}

func TestDescribeExplain(t *testing.T) {
	tests := []struct {
		statement string
		format    string
		executed  bool
		ok        bool
	}{
		{"EXPLAIN SELECT 1", "text", false, true},
		{"EXPLAIN (FORMAT JSON) SELECT 1", "json", false, true},
		{"EXPLAIN (ANALYZE, FORMAT YAML) SELECT 1", "yaml", true, true},
		{"EXPLAIN (FORMAT 1) SELECT 1", "", false, true},
		{"SELECT 1", "", false, false},
		{"EXPLAIN", "", false, false},
	}
	for _, tc := range tests {
		format, executed, ok := DescribeExplain(tc.statement)
		require.Equal(t, tc.format, format, tc.statement)
		require.Equal(t, tc.executed, executed, tc.statement)
		require.Equal(t, tc.ok, ok, tc.statement)
	}
}
