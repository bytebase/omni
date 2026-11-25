package completion

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// analyze parses input with the placeholder at "|" and returns its context.
func analyze(t *testing.T, input string) *nameContext {
	t.Helper()
	stmt, err := parseOne(strings.Replace(input, "|", placeholder, 1))
	require.NoError(t, err, input)
	return analyzeStatement(context.Background(), stmt, newCatalog(testMetadata(), "public"))
}

func refs(sources []source) []string {
	var out []string
	for _, s := range sources {
		out = append(out, s.ref)
	}
	return out
}

func TestAnalyzeStatementKinds(t *testing.T) {
	tests := []struct {
		input     string
		kind      nameKind
		qualifier []string
		scope     []string
		target    string
	}{
		{input: "SELECT * FROM |", kind: kindObject},
		{input: "SELECT * FROM s.|", kind: kindObject, qualifier: []string{"s"}},
		{input: "SELECT * FROM |.t", kind: kindSchema},
		{input: "DROP TABLE |", kind: kindObject},
		{input: "SHOW CREATE TABLE |", kind: kindObject},
		{input: "SELECT | FROM orders o, analytics.events", kind: kindExpr, scope: []string{"o", "events"}},
		{input: "SELECT o.| FROM orders o", kind: kindExpr, qualifier: []string{"o"}, scope: []string{"o"}},
		{input: "SELECT s.t.| FROM s.t", kind: kindExpr, qualifier: []string{"s", "t"}, scope: []string{"t"}},
		{input: "SELECT |(1) FROM orders", kind: kindFunction, scope: []string{"orders"}},
		{input: "SELECT * FROM orders a JOIN orders b ON |", kind: kindExpr, scope: []string{"a", "b"}},
		{input: "SELECT (SELECT | FROM analytics.events) FROM orders", kind: kindExpr, scope: []string{"events", "orders"}},
		{input: "SELECT * FROM (SELECT | FROM analytics.events) AS s, orders", kind: kindExpr, scope: []string{"events"}},
		{input: "SELECT * FROM orders ORDER BY |", kind: kindExpr, scope: []string{"", "orders"}},
		{input: "SELECT * FROM orders@|", kind: kindIndex, target: "orders"},
		{input: "INSERT INTO orders (|) VALUES (1)", kind: kindTargetColumn, target: "orders"},
		{input: "INSERT INTO orders VALUES (1) ON CONFLICT (id) DO UPDATE SET amount = |", kind: kindExpr, scope: []string{"orders", "excluded"}},
		{input: "INSERT INTO orders VALUES (1) RETURNING |", kind: kindExpr, scope: []string{"orders"}},
		{input: "UPSERT INTO orders (|) VALUES (1)", kind: kindTargetColumn, target: "orders"},
		{input: "UPDATE orders o SET | = 1", kind: kindTargetColumn, target: "o"},
		{input: "UPDATE orders SET amount = 1 FROM analytics.events e WHERE |", kind: kindExpr, scope: []string{"orders", "e"}},
		{input: "DELETE FROM orders USING analytics.events WHERE |", kind: kindExpr, scope: []string{"orders", "events"}},
		{input: "CREATE INDEX ON orders (|)", kind: kindTargetColumn, target: "orders"},
		{input: "CREATE INDEX ON orders (id) STORING (|)", kind: kindTargetColumn, target: "orders"},
		{input: "ALTER TABLE orders DROP COLUMN |", kind: kindTargetColumn, target: "orders"},
		{input: "ALTER TABLE orders RENAME COLUMN | TO x", kind: kindTargetColumn, target: "orders"},
		{input: "EXPLAIN SELECT | FROM orders", kind: kindExpr, scope: []string{"orders"}},
		{input: "SELECT * FROM [SHOW COLUMNS FROM |]", kind: kindObject},
		{input: "CREATE VIEW v AS SELECT | FROM orders", kind: kindExpr, scope: []string{"orders"}},
		{input: "CREATE TABLE t (a |)", kind: kindType},
		{input: "ALTER TABLE orders ADD COLUMN a |", kind: kindType},
		{input: "ALTER TABLE orders ALTER COLUMN a TYPE |", kind: kindType},
		{input: "SELECT CAST(1 AS |)", kind: kindType},
		{input: "SELECT 1::|", kind: kindType},
		{input: "SELECT 1 IS OF (|)", kind: kindType},
		{input: "ALTER TABLE orders ADD CHECK (amount::| > 0)", kind: kindType, scope: []string{"orders"}},
		{input: "ALTER TABLE orders ADD CHECK (| > 0)", kind: kindExpr, scope: []string{"orders"}},
		// A DEFAULT may reference no column.
		{input: "ALTER TABLE orders ALTER COLUMN amount SET DEFAULT |", kind: kindExpr},
		{input: "ALTER TABLE orders ADD COLUMN x INT DEFAULT |", kind: kindExpr},
		{input: "ALTER TABLE orders ADD COLUMN x INT AS (|) STORED", kind: kindExpr, scope: []string{"orders"}},
		{input: "CREATE TABLE t (a INT, b INT AS (|) STORED)", kind: kindExpr, scope: []string{"t"}},
		{input: "CREATE TABLE t (a INT, CHECK (|))", kind: kindExpr, scope: []string{"t"}},
		{input: "CREATE TABLE t (a INT DEFAULT |)", kind: kindExpr},
		{input: "CREATE TABLE t (a INT DEFAULT 1::|)", kind: kindType},
	}
	for _, tc := range tests {
		nc := analyze(t, tc.input)
		require.NotNil(t, nc, tc.input)
		require.Equal(t, tc.kind, nc.kind, tc.input)
		require.Equal(t, tc.qualifier, nc.qualifier, tc.input)
		require.Equal(t, tc.scope, refs(nc.scope), tc.input)
		if tc.target != "" {
			require.Equal(t, tc.target, nc.target.ref, tc.input)
		}
	}
}

func TestAnalyzeStatementNewNames(t *testing.T) {
	for _, input := range []string{
		"SELECT * FROM orders AS |",
		"SELECT 1 AS |",
		"CREATE DATABASE |",
	} {
		require.Nil(t, analyze(t, input), input)
	}
}

func TestAnalyzeStatementCTEs(t *testing.T) {
	nc := analyze(t, "WITH a AS (SELECT 1), b(x) AS (SELECT * FROM a) SELECT * FROM |")
	require.Equal(t, kindObject, nc.kind)
	require.Equal(t, []string{"a", "b"}, nc.ctes)

	nc = analyze(t, "WITH b(x, y) AS (SELECT 1, 2) SELECT b.| FROM b")
	require.Equal(t, []string{"x", "y"}, nc.scope[0].columns)

	nc = analyze(t, "WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n + 1 FROM r WHERE |) SELECT * FROM r")
	require.Equal(t, kindExpr, nc.kind)
	require.Equal(t, []string{"n"}, nc.scope[0].columns)

	nc = analyze(t, "WITH v AS (VALUES (1, 2)) SELECT | FROM v")
	require.Equal(t, []string{"column1", "column2"}, nc.scope[0].columns)

	nc = analyze(t, "WITH l AS (SELECT * FROM events_log) SELECT | FROM l")
	require.Equal(t, []string{"message"}, nc.scope[0].columns)
}
