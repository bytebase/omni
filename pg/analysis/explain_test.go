package analysis

import (
	"strings"
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

	for _, statement := range []string{
		"ANALYZE DELETE FROM t",
		"EXPLAIN ANALYZE SELECT 1",
		"EXPLAIN ANALYZE DELETE FROM t",
		"EXPLAIN (ANALYZE) DELETE FROM t",
		"EXPLAIN (ANALYZE, VERBOSE) DELETE FROM t",
		"(ANALYZE) DELETE FROM t",
		"VERBOSE SELECT 1",
		"",
		"-- nothing here",
	} {
		_, err := ExplainStatement(statement, "")
		require.Errorf(t, err, "%q", statement)
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

func TestExplainStatementNeverPlansAPlan(t *testing.T) {
	for _, statement := range []string{
		"EXPLAIN SELECT 1",
		"EXPLAIN INSERT INTO t VALUES (1)",
		"EXPLAIN EXECUTE p",
		"EXPLAIN DECLARE c CURSOR FOR SELECT 1",
		"EXPLAIN CREATE TABLE t AS SELECT 1",
		"EXPLAIN (COSTS OFF) VALUES (1)",
		"EXPLAIN TABLE t",
	} {
		for _, format := range []string{"", "json"} {
			got, err := ExplainStatement(statement, format)
			require.NoErrorf(t, err, "%q", statement)
			require.Equalf(t, 1, strings.Count(strings.ToUpper(got), "EXPLAIN"), "%q became %q", statement, got)
		}
	}
}

func TestCockroachExplainStatement(t *testing.T) {
	tests := []struct {
		statement string
		distsql   bool
		want      string
	}{
		{"SELECT 1", true, "EXPLAIN (DISTSQL) SELECT 1"},
		{"SELECT 1", false, "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1", true, "EXPLAIN (DISTSQL) SELECT 1"},
		{"EXPLAIN SELECT 1", false, "EXPLAIN SELECT 1"},
		{"EXPLAIN (VERBOSE) SELECT 1", true, "EXPLAIN (DISTSQL) SELECT 1"},
		{"EXPLAIN (DISTSQL) SELECT 1", true, "EXPLAIN (DISTSQL) SELECT 1"},
	}
	for _, tc := range tests {
		got, err := CockroachExplainStatement(tc.statement, tc.distsql)
		require.NoError(t, err, tc.statement)
		require.Equal(t, tc.want, got, tc.statement)
	}
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
