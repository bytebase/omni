package completion

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeywordCandidates(t *testing.T) {
	tests := []struct {
		prefix string
		want   []string
		absent []string
	}{
		{prefix: "", want: []string{"SELECT", "UPSERT", "SHOW", "IMPORT", "BACKUP"}, absent: []string{"FROM", "WHERE"}},
		{prefix: "SELECT * FROM t", want: []string{"WHERE", "JOIN", "AS"}, absent: []string{"SELECT", "ABORT"}},
		{prefix: "SELECT * FROM t AS OF", want: []string{"SYSTEM"}},
		{prefix: "SELECT * FROM t AS OF SYSTEM", want: []string{"TIME"}},
		{prefix: "SELECT * FROM t WHERE a", want: []string{"AND", "NOT", "IN", "BETWEEN"}},
		{prefix: "CREATE TABLE", want: []string{"IF"}, absent: []string{"INDEX", "ABORT"}},
		// INT and STRING read here as names of types, which the name
		// context offers; DECIMAL takes a precision a type name cannot.
		{prefix: "CREATE TABLE t (a", want: []string{"DECIMAL", "VARCHAR"}, absent: []string{"ABORT"}},
		{prefix: "CREATE TABLE t (a INT,", want: []string{"CONSTRAINT", "FAMILY", "INDEX", "PRIMARY", "UNIQUE"}},
		{prefix: "ALTER TABLE t", want: []string{"CONFIGURE", "SCATTER", "SPLIT", "RELOCATE"}},
		{prefix: "ALTER TABLE t CONFIGURE", want: []string{"ZONE"}},
		{prefix: "SHOW", want: []string{"RANGES", "JOBS", "ZONE", "CREATE"}},
		{prefix: "INSERT INTO t", want: []string{"VALUES", "SELECT", "DEFAULT"}},
	}
	for _, tc := range tests {
		got := keywordCandidates(context.Background(), tc.prefix, false)
		for _, w := range tc.want {
			require.Containsf(t, got, w, "prefix %q", tc.prefix)
		}
		for _, a := range tc.absent {
			require.NotContainsf(t, got, a, "prefix %q", tc.prefix)
		}
		for _, internal := range []string{"INDEX_BEFORE_PAREN", "NOTHING_AFTER_RETURNING"} {
			require.NotContains(t, got, internal)
		}
	}
}

func TestKeywordCandidatesCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.Empty(t, keywordCandidates(ctx, "SELECT * FROM canceled_prefix", false))
}

func TestHelpFollowers(t *testing.T) {
	require.Contains(t, helpFollowers()["if"], "not")
	require.Contains(t, helpFollowers()["configure"], "zone")
}

func TestPrefixCache(t *testing.T) {
	c := &prefixCache{limit: 2}
	c.put("a", []string{"A"})
	c.put("b", []string{"B"})
	c.put("c", []string{"C"})
	_, ok := c.get("a")
	require.False(t, ok)
	got, ok := c.get("c")
	require.True(t, ok)
	require.Equal(t, []string{"C"}, got)
}
