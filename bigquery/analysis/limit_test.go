package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestStatementWithResultLimit holds the dialect dispatch: BigQuery wraps
// unconditionally because its caller already confirmed a SELECT, and Spanner
// keeps RewriteResultLimit's in-place rewrite. The corpus is the BigQuery
// adapter's cases recorded from the package that shipped them.
func TestStatementWithResultLimit(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		limit     int
		want      string
	}{
		{
			name:      "bigquery caps a suffix-only statement with a pipe limit",
			statement: "|> WHERE x > 1",
			limit:     10,
			want:      "|> WHERE x > 1 |> LIMIT 10",
		},
		{
			name:      "bigquery appends a literal cap after a parameterized pipe limit",
			statement: "|> LIMIT @n",
			limit:     10,
			want:      "|> LIMIT @n |> LIMIT 10",
		},
		{
			name:      "bigquery wraps a select",
			statement: "SELECT * FROM t",
			limit:     10,
			want:      "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery trims the trailing semicolon",
			statement: "SELECT * FROM t;",
			limit:     10,
			want:      "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery trims surrounding whitespace",
			statement: "  SELECT * FROM t;  ",
			limit:     10,
			want:      "WITH result AS (SELECT * FROM t) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery keeps the closing parenthesis out of a line comment",
			statement: "SELECT 1 -- note",
			limit:     10,
			want:      "WITH result AS (SELECT 1 -- note\n) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery drops a terminator before a trailing comment",
			statement: "SELECT 1; /* note */",
			limit:     10,
			want:      "WITH result AS (SELECT 1) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery keeps a trailing block comment",
			statement: "SELECT 1 /* note */",
			limit:     10,
			want:      "WITH result AS (SELECT 1 /* note */\n) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery keeps the closing parenthesis out of a hash comment",
			statement: "SELECT 1 # note",
			limit:     10,
			want:      "WITH result AS (SELECT 1 # note\n) SELECT * FROM result LIMIT 10;",
		},
		{
			name:      "bigquery omits a zero limit",
			statement: "SELECT * FROM t",
			limit:     0,
			want:      "WITH result AS (SELECT * FROM t) SELECT * FROM result;",
		},
		{
			name:      "bigquery omits a negative limit",
			statement: "SELECT * FROM t",
			limit:     -1,
			want:      "WITH result AS (SELECT * FROM t) SELECT * FROM result;",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, StatementWithResultLimit(tc.statement, tc.limit))
		})
	}
}
