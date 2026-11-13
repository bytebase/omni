package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestExplainStatement(t *testing.T) {
	tests := []struct {
		statement string
		json      bool
		want      string
	}{
		{"SELECT 1", false, "EXPLAIN SELECT 1"},
		{"SELECT 1;", false, "EXPLAIN SELECT 1;"},
		{"SELECT 1", true, "EXPLAIN USING JSON SELECT 1"},
		{"WITH x AS (SELECT 1) SELECT * FROM x", false, "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
		// An EXPLAIN is rebuilt around the statement it plans.
		{"EXPLAIN SELECT 1", false, "EXPLAIN SELECT 1"},
		{"EXPLAIN USING JSON SELECT 1", false, "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1", true, "EXPLAIN USING JSON SELECT 1"},
		{"  explain select 1", false, "EXPLAIN select 1"},
	}
	for _, tc := range tests {
		got, err := analysis.ExplainStatement(tc.statement, tc.json)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
	for statement, want := range map[string]string{
		"SELECT 1; SELECT 2": "EXPLAIN plans a single statement",
		"":                   "EXPLAIN needs a statement to plan",
		"-- comment":         "EXPLAIN needs a statement to plan",
	} {
		_, err := analysis.ExplainStatement(statement, false)
		require.EqualErrorf(t, err, want, "%q", statement)
	}
}

func TestDescribeExplain(t *testing.T) {
	for _, tc := range []struct {
		statement string
		format    string
		isExplain bool
	}{
		{statement: "SELECT 1"},
		{statement: "SELECT 1; EXPLAIN SELECT 1"},
		{statement: "EXPLAIN SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN USING TABULAR SELECT 1", format: "text", isExplain: true},
		{statement: "EXPLAIN USING TEXT SELECT 1", format: "text", isExplain: true},
		{statement: "explain using json select 1", format: "json", isExplain: true},
	} {
		format, executed, isExplain := analysis.DescribeExplain(tc.statement)
		require.Equalf(t, tc.format, format, "%q", tc.statement)
		require.Falsef(t, executed, "%q", tc.statement)
		require.Equalf(t, tc.isExplain, isExplain, "%q", tc.statement)
	}
}
