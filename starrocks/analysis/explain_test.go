package analysis

import (
	"strings"
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
		{"WITH x AS (SELECT 1) SELECT * FROM x", "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
		// A statement that already asks for a plan is rebuilt around the statement
		// it plans, so no request plans a plan.
		{"EXPLAIN SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN VERBOSE SELECT 1", "EXPLAIN SELECT 1"},
		{"EXPLAIN SELECT 1 ORDER BY 1", "EXPLAIN SELECT 1 ORDER BY 1"},
		{"EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x", "EXPLAIN WITH x AS (SELECT 1) SELECT * FROM x"},
	}
	for _, tc := range tests {
		got, err := ExplainStatement(tc.statement)
		require.NoErrorf(t, err, "%q", tc.statement)
		require.Equalf(t, tc.want, got, "%q", tc.statement)
	}
}

func TestExplainStatementRejects(t *testing.T) {
	for _, statement := range []string{
		"SELECT 1; SELECT 2",
		"",
	} {
		_, err := ExplainStatement(statement)
		require.Errorf(t, err, "%q", statement)
	}
}

// A plan request over text with several parse failures reports the first one,
// as the engine entry points do, not the ParseErrors summary.
func TestExplainStatementReportsFirstParseError(t *testing.T) {
	_, err := ExplainStatement("EXPLAIN SELECT exposes the SELECT's tables")
	require.Error(t, err)
	require.False(t, strings.Contains(err.Error(), "(and "), "error = %q", err)
}
