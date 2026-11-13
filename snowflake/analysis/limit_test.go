package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestRewriteResultLimit(t *testing.T) {
	tests := []struct {
		sql   string
		limit int
		want  string
	}{
		{"SELECT * FROM T LIMIT 15;", 10, "SELECT * FROM T LIMIT 10\n;"},
		{"SELECT * FROM T LIMIT 5;", 10, "SELECT * FROM T LIMIT 5\n;"},
		{"SELECT * FROM T LIMIT 0", 10, "SELECT * FROM T LIMIT 10\n;"},
		{"SELECT * FROM T OFFSET 12 FETCH 123;", 10, "SELECT * FROM T OFFSET 12 FETCH 10\n;"},
		{"SELECT TOP 125 * FROM T;", 10, "SELECT TOP 10 * FROM T\n;"},
		{"SELECT * FROM T;  \n", 10, "SELECT * FROM T LIMIT 10\n;"},
		// Only the right-most SELECT of a set operation is limited.
		{"SELECT A FROM T UNION SELECT B FROM U", 10, "SELECT A FROM T UNION SELECT B FROM U LIMIT 10\n;"},
		{"SHOW TABLES ->> SELECT * FROM $1;", 100, "SHOW TABLES ->> SELECT * FROM $1 LIMIT 100\n;"},
		// A statement that is not a query is returned unchanged.
		{"SHOW TABLES", 10, "SHOW TABLES\n;"},
	}
	for _, tc := range tests {
		got, err := analysis.RewriteResultLimit(tc.sql, tc.limit)
		require.NoErrorf(t, err, "%q", tc.sql)
		require.Equalf(t, tc.want, got, "%q", tc.sql)
	}
	for _, sql := range []string{
		"SELECT * FROM T OFFSET 2",
		"SELECT * FROM T LIMIT $x",
		"SELECT 1; SELECT 2",
		"SELEC 1",
	} {
		_, err := analysis.RewriteResultLimit(sql, 10)
		require.Errorf(t, err, "%q", sql)
	}
}
