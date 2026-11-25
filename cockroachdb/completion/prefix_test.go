package completion

import (
	"context"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCollapseLists(t *testing.T) {
	tests := []struct{ in, want string }{
		{"SELECT a, b, c, d FROM t WHERE", "SELECT a, d FROM t WHERE"},
		{"SELECT a, b, c, ", "SELECT a, "},
		{"SELECT a, b FROM t", "SELECT a, b FROM t"},
		{"SELECT * FROM t, u, v, w", "SELECT * FROM t, w"},
		{"SELECT * FROM t ORDER BY a, b, c", "SELECT * FROM t ORDER BY a, c"},
		{"INSERT INTO t VALUES (1), (2), (3)", "INSERT INTO t VALUES (1), (3)"},
		{"SELECT f(a, b, c, d) FROM t", "SELECT f(a, b, c, d) FROM t"},
		{"SELECT * FROM t WHERE a AND b AND c AND ", "SELECT * FROM t WHERE a AND "},
		{"SELECT * FROM t WHERE a OR b AND c AND d AND ", "SELECT * FROM t WHERE a OR b AND "},
		{"SELECT * FROM t WHERE a AND b OR c AND d OR e AND ", "SELECT * FROM t WHERE a AND b OR e AND "},
		{"SELECT * FROM t WHERE a BETWEEN 1 AND 2 AND b AND c AND ", "SELECT * FROM t WHERE a BETWEEN 1 AND 2 AND "},
		{"SELECT CASE WHEN a AND b THEN c AND d END AND e", "SELECT CASE WHEN a AND b THEN c AND d END AND e"},
		{"SELECT a AND b AND c, d AND e", "SELECT a AND c, d AND e"},
		{"SELECT * FROM t@{FORCE_INDEX=i, NO_FULL_SCAN}, u, v", "SELECT * FROM t@{FORCE_INDEX=i, NO_FULL_SCAN}, v"},
		{"SELECT a, b, c INTO x, y, z", "SELECT a, c INTO x, y, z"},
		{"SELECT 'open", "SELECT 'open"},
	}
	for _, tc := range tests {
		require.Equal(t, tc.want, collapseLists(tc.in), tc.in)
	}
}

func TestInnermostQuery(t *testing.T) {
	require.Equal(t, "SELECT b FROM u WHERE ", innermostQuery("SELECT * FROM t WHERE x IN (SELECT b FROM u WHERE "))
	require.Equal(t, "SELECT f(b, ", innermostQuery("SELECT * FROM t WHERE x IN (SELECT f(b, "))
	require.Equal(t, "SELECT * FROM t WHERE f(", innermostQuery("SELECT * FROM t WHERE f("))
	// VALUES and WITH may continue differently inside parentheses.
	require.Equal(t, "SELECT * FROM (VALUES", innermostQuery("SELECT * FROM (VALUES"))
	require.Equal(t, "SELECT (WITH c AS (SELECT 1) ", innermostQuery("SELECT (WITH c AS (SELECT 1) "))
	require.Equal(t, "SELECT * FROM t WHERE x IN (SELECT 1) AND ", innermostQuery("SELECT * FROM t WHERE x IN (SELECT 1) AND "))
}

func TestProbePrefixKeepsErrors(t *testing.T) {
	// A prefix that fails before its end is not reduced, so the reduction
	// cannot hide the error.
	p := "SELECT a, b, c, d FROM FROM t WHERE"
	require.Equal(t, p, probePrefix(p))
	require.Equal(t, "SELECT a, d FROM t WHERE", probePrefix("SELECT a, b, c, d FROM t WHERE"))
}

// TestProbePrefixKeepsKeywords holds the reductions to their claim on
// statements of PostgreSQL's regression suite that CockroachDB parses, cut
// at random tokens: the reduced prefix admits exactly the keywords the full
// one does.
func TestProbePrefixKeepsKeywords(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "pg", "pgregress", "testdata", "sql", "*.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, files)
	rng := rand.New(rand.NewPCG(1, 2))
	checked := 0
	for _, f := range files {
		b, err := os.ReadFile(f)
		require.NoError(t, err)
		for script := string(b); script != "" && checked < 150; {
			n, ok := splitFirst(script)
			if !ok {
				break
			}
			stmt := script[:n]
			script = script[n:]
			if len(stmt) > 3000 || rng.IntN(4) != 0 {
				continue
			}
			if _, err := parseOne(stmt); err != nil {
				continue
			}
			toks := tokenize(stmt)
			if len(toks) < 4 {
				continue
			}
			prefix := stmt[:toks[1+rng.IntN(len(toks)-1)].start]
			reduced := probePrefix(prefix)
			if reduced == prefix {
				continue
			}
			checked++
			require.Equal(t, probeKeywords(context.Background(), prefix), probeKeywords(context.Background(), reduced), "%q reduced to %q", prefix, reduced)
		}
	}
	require.Equal(t, 150, checked)
}
