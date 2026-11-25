package completion

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLexicalSources(t *testing.T) {
	got := lexicalSources(`SELECT x FROM orders o, analytics.events AS e JOIN "Order Items" oi ON a = = b WHERE`)
	require.Equal(t, []source{
		{ref: "o", name: "orders"},
		{ref: "e", schema: "analytics", name: "events"},
		{ref: "oi", name: "Order Items"},
	}, got)

	got = lexicalSources("UPDATE orders SET a = = 1")
	require.Equal(t, []source{{ref: "orders", name: "orders"}}, got)

	require.Empty(t, lexicalSources("SELECT 1"))
}

func TestLexicalContext(t *testing.T) {
	word := func(stmt string) caretWord { return caretWord{start: len(stmt), end: len(stmt)} }

	stmt := "SELECT * FROM orders o WHERE a = = 1 AND FROM "
	nc := lexicalContext(stmt, word(stmt))
	require.Equal(t, kindObject, nc.kind)

	stmt = "SELECT * FROM analytics."
	nc = lexicalContext(stmt, word(stmt))
	require.Equal(t, kindObject, nc.kind)
	require.Equal(t, []string{"analytics"}, nc.qualifier)

	stmt = "SELECT o."
	nc = lexicalContext(stmt, word(stmt))
	require.Equal(t, kindExpr, nc.kind)
	require.Equal(t, []string{"o"}, nc.qualifier)

	stmt = "SELECT 1 1 "
	require.Nil(t, lexicalContext(stmt, word(stmt)))

	stmt = "SELECT 'open"
	require.Nil(t, lexicalContext(stmt, word(stmt)))
}

func TestCompleteSessionLexicalFallback(t *testing.T) {
	got := completeAt(t, testConfig(), "SELECT o.| FROM orders o WHERE id = = 1")
	require.True(t, has(got, "amount", CandidateColumn), "got %v", got)
}
