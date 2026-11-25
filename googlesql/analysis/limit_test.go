package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/googlesql/parser"
)

// TestRewriteResultLimit holds the rewriter to the cases the Spanner driver's
// LIMIT rewrite had before the move: the outermost LIMIT moves, and a shape the
// rewrite cannot trust stays as written.
func TestRewriteResultLimit(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      string
	}{
		{
			name:      "select",
			statement: "SELECT * FROM t",
			want:      "SELECT * FROM t LIMIT 100",
		},
		{
			name:      "with",
			statement: "WITH a AS (SELECT 1 AS x) SELECT x FROM a",
			want:      "WITH a AS (SELECT 1 AS x) SELECT x FROM a LIMIT 100",
		},
		{
			name:      "lowercase with, order by, and trailing semicolon",
			statement: "with a as (select 1 x), b as (select x from a) select x from b order by x desc;",
			want:      "with a as (select 1 x), b as (select x from a) select x from b order by x desc LIMIT 100;",
		},
		{
			name:      "with and set operation",
			statement: "WITH a AS (SELECT 1 AS x) (SELECT x FROM a) UNION ALL (SELECT x FROM a) ORDER BY x",
			want:      "WITH a AS (SELECT 1 AS x) (SELECT x FROM a) UNION ALL (SELECT x FROM a) ORDER BY x LIMIT 100",
		},
		{
			name:      "set operation",
			statement: "SELECT 1 UNION ALL SELECT 2",
			want:      "SELECT 1 UNION ALL SELECT 2 LIMIT 100",
		},
		{
			name:      "inner limit belongs to the parenthesized query",
			statement: "(SELECT x FROM t LIMIT 500)",
			want:      "(SELECT x FROM t LIMIT 500) LIMIT 100",
		},
		{
			name:      "trailing comment",
			statement: "SELECT x FROM t -- note",
			want:      "SELECT x FROM t LIMIT 100 -- note",
		},
		{
			name:      "lower limit is kept",
			statement: "WITH a AS (SELECT x FROM t) SELECT x FROM a LIMIT 10",
			want:      "WITH a AS (SELECT x FROM t) SELECT x FROM a LIMIT 10",
		},
		{
			name:      "higher limit is lowered",
			statement: "WITH a AS (SELECT x FROM t) SELECT x FROM a ORDER BY x LIMIT 5000 OFFSET 10",
			want:      "WITH a AS (SELECT x FROM t) SELECT x FROM a ORDER BY x LIMIT 100 OFFSET 10",
		},
		{
			name:      "non-literal limit is unchanged",
			statement: "WITH a AS (SELECT x FROM t) SELECT x FROM a LIMIT CAST(5000 AS INT64)",
			want:      "WITH a AS (SELECT x FROM t) SELECT x FROM a LIMIT CAST(5000 AS INT64)",
		},
		{
			name:      "for update is unchanged",
			statement: "SELECT x FROM t ORDER BY x FOR UPDATE",
			want:      "SELECT x FROM t ORDER BY x FOR UPDATE",
		},
		{
			// A pipe query's limit is a `|> LIMIT` operator, so the cap is
			// appended as one.
			name:      "pipe query gets a pipe limit",
			statement: "WITH a AS (SELECT x FROM t) SELECT x FROM a |> WHERE x > 1",
			want:      "WITH a AS (SELECT x FROM t) SELECT x FROM a |> WHERE x > 1 |> LIMIT 100",
		},
		{
			name:      "pipe query limit is lowered",
			statement: "SELECT x FROM t |> LIMIT 5000",
			want:      "SELECT x FROM t |> LIMIT 100",
		},
		{
			name:      "pipe query limit all is capped by another operator",
			statement: "SELECT x FROM t |> LIMIT ALL",
			want:      "SELECT x FROM t |> LIMIT ALL |> LIMIT 100",
		},
		{
			name:      "suffix-only statement gets a pipe limit",
			statement: "|> WHERE x > 1",
			want:      "|> WHERE x > 1 |> LIMIT 100",
		},
		{
			name:      "non-literal pipe limit is unchanged",
			statement: "SELECT x FROM t |> LIMIT CAST(5000 AS INT64)",
			want:      "SELECT x FROM t |> LIMIT CAST(5000 AS INT64)",
		},
		{
			// The `TABLE t` query primary is a query, so the in-place LIMIT binds at
			// the clause's end like any other query.
			name:      "table clause",
			statement: "TABLE t",
			want:      "TABLE t LIMIT 100",
		},
		{
			name:      "table clause with filter",
			statement: "TABLE t WHERE a > 1",
			want:      "TABLE t WHERE a > 1 LIMIT 100",
		},
		{
			name:      "dml is unchanged",
			statement: "UPDATE t SET x = 1 WHERE id = 2",
			want:      "UPDATE t SET x = 1 WHERE id = 2",
		},
		{
			name:      "unparsable query is unchanged",
			statement: "SELECT x FROM",
			want:      "SELECT x FROM",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, RewriteResultLimit(tc.statement, 100))
		})
	}
}

func TestRewriteResultLimitParses(t *testing.T) {
	for _, statement := range []string{
		"SELECT x FROM t",
		"WITH a AS (SELECT x FROM t) SELECT x FROM a |> WHERE x > 1",
		"SELECT x FROM t |> WHERE x > 1 |> LIMIT 5000",
		"SELECT x FROM t |> AGGREGATE COUNT(*) AS n",
	} {
		rewritten := RewriteResultLimit(statement, 100)
		_, err := parser.Parse(rewritten)
		require.NoError(t, err, "rewritten: %s", rewritten)
	}
}
