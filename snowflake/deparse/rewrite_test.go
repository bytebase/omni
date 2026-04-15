package deparse_test

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/stretchr/testify/require"
)

// injectLimitCase is one test case for InjectLimit.
type injectLimitCase struct {
	name    string
	sql     string
	maxRows int
	want    string // expected output SQL; empty means check wantContains / wantPrefix
	// wantContains are substrings that must all appear in the output.
	wantContains []string
	wantErr      bool
}

func TestInjectLimit(t *testing.T) {
	cases := []injectLimitCase{
		// 1. No LIMIT — add LIMIT.
		{
			name:    "no limit adds limit",
			sql:     "SELECT * FROM t",
			maxRows: 100,
			want:    "SELECT * FROM t LIMIT 100",
		},
		// 2. Existing LIMIT n <= maxRows — leave unchanged.
		{
			name:    "existing limit smaller unchanged",
			sql:     "SELECT * FROM t LIMIT 50",
			maxRows: 100,
			want:    "SELECT * FROM t LIMIT 50",
		},
		// 3. Existing LIMIT n > maxRows — lower it.
		{
			name:    "existing limit larger lowered",
			sql:     "SELECT * FROM t LIMIT 200",
			maxRows: 100,
			want:    "SELECT * FROM t LIMIT 100",
		},
		// 4. Existing LIMIT with non-integer expression — wrap.
		{
			name:    "non-literal limit wraps",
			sql:     "SELECT * FROM t LIMIT 1+1",
			maxRows: 100,
			wantContains: []string{
				"SELECT * FROM (SELECT * FROM t LIMIT 1 + 1)",
				"LIMIT 100",
			},
		},
		// 5. FETCH FIRST n > maxRows — lower in-place.
		{
			name:    "fetch first lowered",
			sql:     "SELECT * FROM t FETCH FIRST 200 ROWS ONLY",
			maxRows: 100,
			want:    "SELECT * FROM t FETCH FIRST 100 ROWS ONLY",
		},
		// 5b. FETCH FIRST n <= maxRows — unchanged.
		{
			name:    "fetch first smaller unchanged",
			sql:     "SELECT * FROM t FETCH FIRST 50 ROWS ONLY",
			maxRows: 100,
			want:    "SELECT * FROM t FETCH FIRST 50 ROWS ONLY",
		},
		// 6. UNION — wrap entire set operation.
		{
			name:    "union wrapped",
			sql:     "SELECT a FROM t UNION SELECT b FROM u",
			maxRows: 100,
			wantContains: []string{
				"SELECT * FROM (SELECT a FROM t UNION SELECT b FROM u)",
				"LIMIT 100",
			},
		},
		// 7. WITH CTE + SELECT — add LIMIT at outer SELECT level.
		{
			name:    "cte select adds limit",
			sql:     "WITH cte AS (SELECT id FROM src) SELECT * FROM cte",
			maxRows: 100,
			wantContains: []string{
				"WITH cte AS (SELECT id FROM src)",
				"SELECT * FROM cte",
				"LIMIT 100",
			},
		},
		// 8. INSERT — unchanged.
		{
			name:    "insert unchanged",
			sql:     "INSERT INTO t VALUES (1)",
			maxRows: 100,
			want:    "INSERT INTO t VALUES (1)",
		},
		// 9. CREATE TABLE — unchanged.
		{
			name:    "create table unchanged",
			sql:     "CREATE TABLE t (id INT)",
			maxRows: 100,
			want:    "CREATE TABLE t (id INT)",
		},
		// 10. Multi-statement — apply to SELECTs only.
		{
			name:    "multi-statement select and non-select",
			sql:     "SELECT 1; INSERT INTO t VALUES (1); SELECT 2",
			maxRows: 100,
			wantContains: []string{
				"SELECT 1 LIMIT 100",
				"INSERT INTO t VALUES (1)",
				"SELECT 2 LIMIT 100",
			},
		},
		// 11. Invalid SQL → error.
		{
			name:    "invalid sql returns error",
			sql:     "SELECT FROM WHERE",
			maxRows: 100,
			wantErr: true,
		},
		// 12. maxRows <= 0 → error.
		{
			name:    "zero maxRows error",
			sql:     "SELECT * FROM t",
			maxRows: 0,
			wantErr: true,
		},
		{
			name:    "negative maxRows error",
			sql:     "SELECT * FROM t",
			maxRows: -1,
			wantErr: true,
		},
		// Edge: LIMIT exactly equal to maxRows — unchanged.
		{
			name:    "limit equal to maxRows unchanged",
			sql:     "SELECT * FROM t LIMIT 100",
			maxRows: 100,
			want:    "SELECT * FROM t LIMIT 100",
		},
		// Edge: FETCH FIRST via non-literal expression → wrap.
		{
			name:    "non-literal fetch wraps",
			sql:     "SELECT * FROM t FETCH FIRST (SELECT 10) ROWS ONLY",
			maxRows: 100,
			wantContains: []string{
				"SELECT * FROM (SELECT * FROM t FETCH FIRST (SELECT 10) ROWS ONLY)",
				"LIMIT 100",
			},
		},
		// Edge: INTERSECT — wrap.
		{
			name:    "intersect wrapped",
			sql:     "SELECT a FROM t INTERSECT SELECT b FROM u",
			maxRows: 50,
			wantContains: []string{
				"SELECT * FROM (SELECT a FROM t INTERSECT SELECT b FROM u)",
				"LIMIT 50",
			},
		},
		// Edge: SELECT with ORDER BY, no LIMIT — add LIMIT after ORDER BY.
		{
			name:    "select with order by adds limit",
			sql:     "SELECT id FROM t ORDER BY id DESC",
			maxRows: 10,
			want:    "SELECT id FROM t ORDER BY id DESC LIMIT 10",
		},
		// Edge: SELECT with WHERE, no LIMIT.
		{
			name:    "select with where adds limit",
			sql:     "SELECT id, name FROM employees WHERE active = TRUE",
			maxRows: 25,
			want:    "SELECT id, name FROM employees WHERE active = TRUE LIMIT 25",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := deparse.InjectLimit(tc.sql, tc.maxRows)
			if tc.wantErr {
				require.Error(t, err, "expected an error for input %q", tc.sql)
				return
			}
			require.NoError(t, err)
			if tc.want != "" {
				require.Equal(t, tc.want, got)
			}
			for _, sub := range tc.wantContains {
				require.True(t,
					strings.Contains(got, sub),
					"output %q does not contain %q", got, sub)
			}
		})
	}
}
