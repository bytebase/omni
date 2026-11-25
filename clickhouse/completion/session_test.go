package completion

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

// caret removes the "|" marker from input and returns the text with the
// marker's 1-based line and 0-based UTF-16 column.
func caret(t *testing.T, input string) (string, int, int) {
	t.Helper()
	line, col := 1, 0
	for i, r := range input {
		if r == '|' {
			return input[:i] + input[i+1:], line, col
		}
		if r == '\n' {
			line++
			col = 0
			continue
		}
		if r > 0xFFFF {
			col += 2
		} else {
			col++
		}
	}
	t.Fatalf("missing caret marker in %q", input)
	return "", 0, 0
}

func testConfig() SessionConfig {
	dbs := map[string]*metadatapb.DatabaseSchemaMetadata{
		"db": {Name: "db", Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{
				{Name: "orders", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "UInt64"}, {Name: "user_id", Type: "UInt64"}, {Name: "amount", Type: "Decimal(18, 2)"}}},
				{Name: "users", Columns: []*metadatapb.ColumnMetadata{{Name: "id", Type: "UInt64"}, {Name: "name", Type: "String"}}},
				{Name: "my-table", Columns: []*metadatapb.ColumnMetadata{{Name: "select", Type: "String"}}},
			},
			Views:             []*metadatapb.ViewMetadata{{Name: "order_totals", Columns: []*metadatapb.ColumnMetadata{{Name: "total", Type: "Decimal(38, 2)"}}}},
			MaterializedViews: []*metadatapb.MaterializedViewMetadata{{Name: "daily_orders"}},
		}}},
		"other": {Name: "other", Schemas: []*metadatapb.SchemaMetadata{{
			Tables: []*metadatapb.TableMetadata{{Name: "events", Columns: []*metadatapb.ColumnMetadata{{Name: "ts", Type: "DateTime"}}}},
		}}},
	}
	return SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			if db, ok := dbs[name]; ok {
				return db, nil
			}
			return nil, fmt.Errorf("database %q not found", name)
		},
		ListDatabaseNames: func(context.Context) ([]string, error) {
			return []string{"db", "other"}, nil
		},
	}
}

func completeAt(t *testing.T, cfg SessionConfig, input string) []Candidate {
	t.Helper()
	text, line, col := caret(t, input)
	return CompleteSession(context.Background(), cfg, text, line, col)
}

// texts returns the candidates of type typ, in order.
func texts(cands []Candidate, typ CandidateType) []string {
	var out []string
	for _, c := range cands {
		if c.Type == typ {
			out = append(out, c.Text)
		}
	}
	return out
}

func TestCompleteSessionKeywords(t *testing.T) {
	cfg := testConfig()
	for _, tc := range []struct {
		input   string
		has     []string
		hasNone []string
	}{
		{input: "|", has: []string{"SELECT", "INSERT", "CREATE", "ALTER", "SYSTEM", "OPTIMIZE"}, hasNone: []string{"WHERE", "LIMIT"}},
		{input: "SEL|", has: []string{"SELECT"}, hasNone: []string{"INSERT"}},
		{input: "SELECT * FROM orders |", has: []string{"WHERE", "PREWHERE", "FINAL", "SAMPLE", "GROUP", "ORDER", "LIMIT", "SETTINGS", "FORMAT", "LEFT", "ARRAY", "JOIN"}, hasNone: []string{"SELECT", "VALUES"}},
		{input: "SELECT * FROM orders o LEFT |", has: []string{"JOIN", "OUTER", "ANY", "SEMI"}, hasNone: []string{"WHERE"}},
		{input: "SELECT a FROM t GROUP |", has: []string{"BY"}},
		{input: "CREATE |", has: []string{"TABLE", "VIEW", "MATERIALIZED", "DATABASE", "DICTIONARY", "FUNCTION", "USER"}},
		{input: "CREATE TABLE t (a UInt64) |", has: []string{"ENGINE", "AS", "COMMENT"}},
		{input: "ALTER TABLE t |", has: []string{"ADD", "DROP", "MODIFY", "UPDATE", "DELETE", "FREEZE"}},
		{input: "SYSTEM |", has: []string{"FLUSH", "RELOAD", "START", "STOP"}},
		{input: "SELECT count() OVER (ORDER BY a |", has: []string{"ROWS", "RANGE", "GROUPS"}},
		{input: "INSERT INTO t |", has: []string{"VALUES", "FORMAT", "SELECT", "SETTINGS"}},
		{input: "SELECT 1;\nSHOW |", has: []string{"TABLES", "DATABASES", "PROCESSLIST", "CREATE"}},
		{input: "SELECT 1; SEL|", has: []string{"SELECT"}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			got := texts(completeAt(t, cfg, tc.input), CandidateKeyword)
			for _, kw := range tc.has {
				require.Contains(t, got, kw)
			}
			for _, kw := range tc.hasNone {
				require.NotContains(t, got, kw)
			}
		})
	}
}

func TestCompleteSessionNames(t *testing.T) {
	cfg := testConfig()
	for _, tc := range []struct {
		input string
		typ   CandidateType
		want  []string
	}{
		{input: "SELECT * FROM |", typ: CandidateTable, want: []string{"`my-table`", "orders", "users"}},
		{input: "SELECT * FROM |", typ: CandidateView, want: []string{"order_totals"}},
		{input: "SELECT * FROM |", typ: CandidateMaterializedView, want: []string{"daily_orders"}},
		{input: "SELECT * FROM |", typ: CandidateDatabase, want: []string{"db", "other"}},
		{input: "SELECT * FROM ord|", typ: CandidateTable, want: []string{"orders"}},
		{input: "SELECT * FROM my|", typ: CandidateTable, want: []string{"`my-table`"}},
		{input: "SELECT * FROM `my|", typ: CandidateTable, want: []string{"my-table"}},
		{input: "SELECT * FROM orders, |", typ: CandidateTable, want: []string{"`my-table`", "orders", "users"}},
		{input: "SELECT * FROM orders o, |", typ: CandidateTable, want: []string{"`my-table`", "orders", "users"}},
		{input: "SELECT * FROM orders ORDER BY id, |", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
		{input: "SELECT * FROM orders ORDER BY id, |", typ: CandidateTable},
		{input: "SELECT user_id FROM orders GROUP BY user_id, |", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
		{input: "SELECT * FROM other.|", typ: CandidateTable, want: []string{"events"}},
		{input: "INSERT INTO |", typ: CandidateTable, want: []string{"`my-table`", "orders", "users"}},
		{input: "OPTIMIZE TABLE u|", typ: CandidateTable, want: []string{"users"}},
		{input: "USE |", typ: CandidateDatabase, want: []string{"db", "other"}},
		{input: "SELECT | FROM orders", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
		{input: "SELECT u| FROM orders", typ: CandidateColumn, want: []string{"user_id"}},
		{input: "SELECT * FROM orders WHERE a|", typ: CandidateColumn, want: []string{"amount"}},
		{input: "SELECT o.| FROM orders AS o JOIN users u ON o.id = u.id", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
		{input: "SELECT u.| FROM orders AS o JOIN users u ON o.id = u.id", typ: CandidateColumn, want: []string{"id", "name"}},
		{input: "SELECT * FROM orders o JOIN users u ON o.id = u.|", typ: CandidateColumn, want: []string{"id", "name"}},
		{input: "SELECT orders.| FROM orders", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
		{input: "SELECT e.| FROM other.events e", typ: CandidateColumn, want: []string{"ts"}},
		{input: "SELECT other.events.| FROM other.events", typ: CandidateColumn, want: []string{"ts"}},
		{input: "SELECT | FROM `my-table`", typ: CandidateColumn, want: []string{"select"}},
		{input: "SELECT * FROM orders WHERE id IN (SELECT | FROM users)", typ: CandidateColumn, want: []string{"amount", "id", "name", "user_id"}},
		{input: "ALTER TABLE orders UPDATE amount = 0 WHERE |", typ: CandidateColumn, want: []string{"amount", "id", "user_id"}},
	} {
		t.Run(tc.input, func(t *testing.T) {
			require.Equal(t, tc.want, texts(completeAt(t, cfg, tc.input), tc.typ))
		})
	}
}

func TestCompleteSessionQualifiedHasNoKeywords(t *testing.T) {
	got := completeAt(t, testConfig(), "SELECT * FROM other.|")
	require.Empty(t, texts(got, CandidateKeyword))
}

func TestCompleteSessionNothingToComplete(t *testing.T) {
	cfg := testConfig()
	for _, input := range []string{
		"SELECT 'ord|",
		"SELECT 1 -- comment |",
		"SELECT 1 # comment |",
		"SELECT 1 /* open |",
	} {
		require.Empty(t, completeAt(t, cfg, input), input)
	}
}

func TestCompleteSessionWithoutMetadata(t *testing.T) {
	failing := SessionConfig{
		DefaultDatabase: "db",
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return nil, errors.New("down")
		},
		ListDatabaseNames: func(context.Context) ([]string, error) { return nil, errors.New("down") },
	}
	for _, cfg := range []SessionConfig{{}, failing} {
		got := completeAt(t, cfg, "SELECT * FROM |")
		require.Empty(t, texts(got, CandidateTable))
		got = completeAt(t, cfg, "SELECT * FROM t |")
		require.Contains(t, texts(got, CandidateKeyword), "WHERE")
	}
}

func TestCompleteSessionMultiLineUTF16(t *testing.T) {
	got := completeAt(t, testConfig(), "SELECT '😀😀' AS x,\n  | FROM orders")
	require.Contains(t, texts(got, CandidateColumn), "amount")
}

// TestCompleteSessionBudget holds a long statement to the completion budget.
func TestCompleteSessionBudget(t *testing.T) {
	saved := completionBudget
	completionBudget = 50 * time.Millisecond
	t.Cleanup(func() { completionBudget = saved })
	var cols []string
	for i := range 3000 {
		cols = append(cols, fmt.Sprintf("coalesce(c%d, 1) + %d AS a%d", i, i, i))
	}
	start := time.Now()
	completeAt(t, testConfig(), "SELECT "+strings.Join(cols, ", ")+" FROM orders |")
	require.Less(t, time.Since(start), 5*time.Second)
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	for name, want := range map[string]string{
		"orders":   "orders",
		"Orders_1": "Orders_1",
		"my-table": "`my-table`",
		"select":   "select",
		"from":     "`from`",
		"id":       "id",
		"null":     "`null`",
		"true":     "`true`",
		"FALSE":    "`FALSE`",
		"nan":      "`nan`",
		"1abc":     "`1abc`",
		"a`b":      "`a``b`",
		`a\b`:      "`a\\\\b`",
	} {
		require.Equal(t, want, QuoteIdentifierIfNeeded(name), name)
	}
}

// TestCompleteSessionEveryCaret completes at every byte of statements that
// end in odd places: an INSERT's data, an open quote, a comment.
func TestCompleteSessionEveryCaret(t *testing.T) {
	cfg := testConfig()
	for _, script := range []string{
		"INSERT INTO orders FORMAT CSV\n1,2,3\n;SELECT o.id FROM orders AS o WHERE o.amount > 1",
		"SELECT `my-table`.`select`, \"x\" FROM `my-table` -- tail\n/* open",
		"WITH t AS (SELECT 1) SELECT * FROM t ARRAY JOIN [1] AS x SETTINGS a = 'b",
		"SELECT '😀', x'00' FROM other.events e GLOBAL ANY LEFT JOIN users USING (id) FORMAT JSON",
	} {
		for i := 0; i <= len(script); i++ {
			before := script[:i]
			line := strings.Count(before, "\n") + 1
			col := 0
			for _, r := range before[strings.LastIndexByte(before, '\n')+1:] {
				col++
				if r > 0xFFFF {
					col++
				}
			}
			CompleteSession(context.Background(), cfg, script, line, col)
		}
	}
}
