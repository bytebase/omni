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

func has(cands []Candidate, text string, typ CandidateType) bool {
	for _, c := range cands {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

func hasText(cands []Candidate, text string) bool {
	for _, c := range cands {
		if c.Text == text {
			return true
		}
	}
	return false
}

func testMetadata() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "orders",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "INT8"},
							{Name: "amount", Type: "DECIMAL"},
						},
						Indexes: []*metadatapb.IndexMetadata{
							{Name: "orders_pkey", Primary: true},
							{Name: "orders_amount_idx"},
						},
					},
					{
						Name: "order",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "INT8"},
						},
					},
					{
						Name: "Order Items",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "Item ID", Type: "INT8"},
							{Name: "Order ID", Type: "INT8"},
						},
					},
					{
						Name: "my-table",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "col1", Type: "INT8"},
						},
					},
					{
						Name: "events_log",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "message", Type: "STRING"},
							{Name: "rowid", Type: "INT8", IsInvisible: true},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name: "active_orders",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "INT8"},
							{Name: "amount", Type: "DECIMAL"},
						},
					},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					{Name: "orders_summary", Definition: "SELECT id, amount FROM orders"},
				},
				Sequences: []*metadatapb.SequenceMetadata{
					{Name: "order_seq"},
				},
				EnumTypes: []*metadatapb.EnumTypeMetadata{
					{Name: "status", Values: []string{"open", "closed"}},
					{Name: "OrderStatus", Values: []string{"new"}},
				},
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "CalculateTotal", Definition: `CREATE FUNCTION public."CalculateTotal"() RETURNS INT8 LANGUAGE SQL AS $$ SELECT 1 $$`},
				},
			},
			{
				Name: "analytics",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "events",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "event_id", Type: "INT8"},
						},
					},
				},
			},
		},
	}
}

func testConfig() SessionConfig {
	meta := testMetadata()
	return SessionConfig{
		DefaultDatabase: "db",
		DefaultSchema:   "public",
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return meta, nil
		},
	}
}

func completeAt(t *testing.T, cfg SessionConfig, input string) []Candidate {
	t.Helper()
	stmt, line, col := caret(t, input)
	return CompleteSession(context.Background(), cfg, stmt, line, col)
}

func TestQuoteIdentifierIfNeeded(t *testing.T) {
	require.Equal(t, "orders", QuoteIdentifierIfNeeded("orders", false))
	require.Equal(t, "_orders_1", QuoteIdentifierIfNeeded("_orders_1", false))
	require.Equal(t, `"order"`, QuoteIdentifierIfNeeded("order", false))
	require.Equal(t, `"Order Items"`, QuoteIdentifierIfNeeded("Order Items", false))
	require.Equal(t, `"my-table"`, QuoteIdentifierIfNeeded("my-table", false))
	require.Equal(t, `"1abc"`, QuoteIdentifierIfNeeded("1abc", false))
	require.Equal(t, `"a""b"`, QuoteIdentifierIfNeeded(`a"b`, false))
	require.Equal(t, `""`, QuoteIdentifierIfNeeded("", false))
	// An unreserved CockroachDB keyword reads back bare.
	require.Equal(t, "zone", QuoteIdentifierIfNeeded("zone", false))

	require.Equal(t, "order", QuoteIdentifierIfNeeded("order", true))
	require.Equal(t, `a""b`, QuoteIdentifierIfNeeded(`a"b`, true))
}

func TestLineOffsetToBytePos(t *testing.T) {
	input := "SELECT '😀', | FROM orders"
	stmt, line, col := caret(t, input)
	require.Equal(t, strings.Index(input, "|"), lineOffsetToBytePos(stmt, line, col))
	require.Equal(t, 3, lineOffsetToBytePos("abc\ndef", 1, 99))
	require.Equal(t, 7, lineOffsetToBytePos("abc\ndef", 9, 0))
}

func TestCompleteSessionQuotingAndSearchPath(t *testing.T) {
	cfg := testConfig()

	got := completeAt(t, cfg, "SELECT 1 FROM |")
	for _, want := range []Candidate{
		{Text: "public", Type: CandidateSchema},
		{Text: "analytics", Type: CandidateSchema},
		{Text: "orders", Type: CandidateTable},
		{Text: `"order"`, Type: CandidateTable},
		{Text: `"Order Items"`, Type: CandidateTable},
		{Text: `"my-table"`, Type: CandidateTable},
		{Text: "active_orders", Type: CandidateView},
		{Text: "orders_summary", Type: CandidateMaterializedView},
		{Text: "order_seq", Type: CandidateSequence},
	} {
		require.Truef(t, has(got, want.Text, want.Type), "missing %v in %v", want, got)
	}

	// Inside double quotes the outer quotes are left out.
	got = completeAt(t, cfg, `SELECT * FROM "ord|"`)
	require.True(t, has(got, "order", CandidateTable), "got %v", got)
	require.True(t, has(got, "orders", CandidateTable), "got %v", got)
	require.False(t, has(got, "SELECT", CandidateKeyword), "got %v", got)

	got = completeAt(t, cfg, `SELECT * FROM "Ord|`)
	require.True(t, has(got, "Order Items", CandidateTable), "got %v", got)

	// A quote in a comment or a dollar-quoted string does not open one.
	got = completeAt(t, cfg, `SELECT * FROM /* " */ |`)
	require.True(t, has(got, `"Order Items"`, CandidateTable), "got %v", got)
	got = completeAt(t, cfg, `SELECT $tag$"$tag$, | FROM "Order Items"`)
	require.True(t, has(got, `"Item ID"`, CandidateColumn), "got %v", got)

	// A non-public DefaultSchema resolves unqualified names first.
	analyticsCfg := testConfig()
	analyticsCfg.DefaultSchema = "analytics"
	got = completeAt(t, analyticsCfg, "SELECT | FROM events")
	require.True(t, has(got, "event_id", CandidateColumn), "got %v", got)
}

func TestCompleteSessionStatementIsolation(t *testing.T) {
	cfg := testConfig()

	got := completeAt(t, cfg, "SELECT * FROM orders; SELECT | FROM analytics.events")
	require.True(t, has(got, "event_id", CandidateColumn), "got %v", got)
	require.False(t, has(got, "amount", CandidateColumn), "got %v", got)

	got = completeAt(t, cfg, "SELECT * FROM orders;| SELECT 1")
	require.True(t, has(got, "SELECT", CandidateKeyword), "got %v", got)

	got = completeAt(t, cfg, "SELECT 1;|")
	require.True(t, has(got, "SELECT", CandidateKeyword), "got %v", got)
	require.True(t, has(got, "UPSERT", CandidateKeyword), "got %v", got)

	got = completeAt(t, cfg, "SELECT '😀', | FROM orders")
	require.True(t, has(got, "id", CandidateColumn), "got %v", got)
	require.True(t, has(got, "amount", CandidateColumn), "got %v", got)

	// Broken later statements and a long script do not stall completion.
	var sheet strings.Builder
	sheet.WriteString("SELECT  FROM orders JOIN analytics.events a2 ON orders.id = a2.event_id;\n")
	sheet.WriteString("SELEC broken FROM oops;\n")
	for i := range 2000 {
		fmt.Fprintf(&sheet, "SELECT col_a, col_b FROM table_%04d WHERE col_a = %d;\n", i, i)
	}
	started := time.Now()
	got = CompleteSession(context.Background(), cfg, sheet.String(), 1, 7)
	require.Less(t, time.Since(started), 2*time.Second)
	require.True(t, has(got, "id", CandidateColumn), "got %v", got)
	require.True(t, has(got, "event_id", CandidateColumn), "got %v", got)
}

func TestCompleteSessionFunctions(t *testing.T) {
	cfg := testConfig()

	got := completeAt(t, cfg, "SELECT | FROM orders")
	require.True(t, has(got, `"CalculateTotal"()`, CandidateFunction), "got %v", got)

	got = completeAt(t, cfg, "SELECT Calc|(1) FROM orders")
	require.True(t, has(got, `"CalculateTotal"`, CandidateFunction), "got %v", got)

	got = completeAt(t, cfg, `SELECT "Calc| FROM orders`)
	require.True(t, has(got, "CalculateTotal", CandidateFunction), "got %v", got)
}

func TestCompleteSessionWithoutMetadata(t *testing.T) {
	for _, cfg := range []SessionConfig{
		{DefaultDatabase: "db", DefaultSchema: "public"},
		{
			DefaultDatabase: "db",
			DefaultSchema:   "public",
			LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
				return nil, errors.New("boom")
			},
		},
	} {
		for _, c := range completeAt(t, cfg, "SELECT 1 FROM |") {
			require.NotContains(t, []CandidateType{
				CandidateSchema, CandidateTable, CandidateView,
				CandidateMaterializedView, CandidateColumn, CandidateSequence,
			}, c.Type)
		}

		// Statement-local names stay available.
		got := completeAt(t, cfg, "WITH c(x) AS (SELECT 1) SELECT | FROM c")
		require.True(t, has(got, "x", CandidateColumn), "got %v", got)

		got = completeAt(t, cfg, "WITH c(x) AS (SELECT 1) SELECT * FROM |")
		require.True(t, has(got, "c", CandidateTable), "got %v", got)
	}
}

func TestCompleteSessionMetadataMatrix(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []Candidate
		absent []string
	}{
		{
			name:   "schema-qualified FROM offers that schema's tables",
			input:  "SELECT 1 FROM analytics.|",
			want:   []Candidate{{Text: "events", Type: CandidateTable}},
			absent: []string{"orders", "active_orders", "SELECT"},
		},
		{
			name:  "database-qualified FROM reads the public schema",
			input: "SELECT 1 FROM db.|",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "projection offers table columns",
			input: "SELECT | FROM orders",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "WHERE offers table columns",
			input: "SELECT * FROM orders WHERE |",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "JOIN ON offers both sides' columns",
			input: "SELECT * FROM orders o JOIN analytics.events e ON |",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "event_id", Type: CandidateColumn}},
		},
		{
			name:   "alias-qualified columns",
			input:  "SELECT o.| FROM orders AS o",
			want:   []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
			absent: []string{"event_id", "SELECT"},
		},
		{
			name:   "schema table-qualified columns",
			input:  "SELECT events.| FROM analytics.events",
			want:   []Candidate{{Text: "event_id", Type: CandidateColumn}},
			absent: []string{"amount"},
		},
		{
			name:   "unknown qualifier does not broaden column scope",
			input:  "SELECT missing.| FROM orders",
			absent: []string{"id", "amount", "event_id"},
		},
		{
			name:  "view columns",
			input: "SELECT | FROM active_orders",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "materialized view columns come from its definition",
			input: "SELECT | FROM orders_summary",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "CTE relation after FROM",
			input: "WITH recent_orders AS (SELECT * FROM orders) SELECT * FROM |",
			want:  []Candidate{{Text: "recent_orders", Type: CandidateTable}},
		},
		{
			name:  "CTE columns from its column list",
			input: "WITH recent_orders(id_alias, amount_alias) AS (SELECT id, amount FROM orders) SELECT | FROM recent_orders",
			want:  []Candidate{{Text: "id_alias", Type: CandidateColumn}, {Text: "amount_alias", Type: CandidateColumn}},
		},
		{
			name:  "CTE columns from a star",
			input: "WITH r AS (SELECT * FROM orders) SELECT r.| FROM r",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "set-operation CTE takes its left branch's columns",
			input: "WITH combined AS (SELECT id, amount FROM orders UNION SELECT event_id, event_id FROM analytics.events) SELECT combined.| FROM combined",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "subquery alias columns",
			input: "SELECT s.| FROM (SELECT id AS order_id FROM orders) AS s",
			want:  []Candidate{{Text: "order_id", Type: CandidateColumn}},
		},
		{
			name:  "correlated subquery sees the outer table",
			input: "SELECT * FROM orders o WHERE EXISTS (SELECT 1 FROM analytics.events e WHERE e.event_id = o.|)",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "INSERT INTO offers tables",
			input: "INSERT INTO |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}, {Text: "events", Type: CandidateTable}},
		},
		{
			name:   "INSERT column list offers the target's columns",
			input:  "INSERT INTO orders (|",
			want:   []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
			absent: []string{"event_id"},
		},
		{
			name:  "UPDATE offers tables",
			input: "UPDATE |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "UPDATE SET offers the target's columns",
			input: "UPDATE orders SET |",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "UPDATE WHERE offers the target's columns",
			input: "UPDATE orders SET amount = 1 WHERE |",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}},
		},
		{
			name:  "DELETE FROM offers tables",
			input: "DELETE FROM |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "DELETE WHERE offers the target's columns",
			input: "DELETE FROM orders WHERE |",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "quoted column names are quoted",
			input: `SELECT | FROM "Order Items"`,
			want:  []Candidate{{Text: `"Item ID"`, Type: CandidateColumn}, {Text: `"Order ID"`, Type: CandidateColumn}},
		},
		{
			name:   "typed prefix filters objects",
			input:  "SELECT 1 FROM ord|",
			want:   []Candidate{{Text: "orders", Type: CandidateTable}, {Text: "orders_summary", Type: CandidateMaterializedView}},
			absent: []string{"events", "active_orders"},
		},
		{
			name:  "incomplete trailing clause still sees the FROM table",
			input: "SELECT | FROM orders WHERE",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}},
		},
		{
			name:   "hidden rowid is not offered unprompted",
			input:  "SELECT | FROM events_log",
			want:   []Candidate{{Text: "message", Type: CandidateColumn}},
			absent: []string{"rowid"},
		},
		{
			name:  "hidden rowid is offered once named",
			input: "SELECT row| FROM events_log",
			want:  []Candidate{{Text: "rowid", Type: CandidateColumn}},
		},
		{
			name:   "star in a CTE leaves out the hidden rowid",
			input:  "WITH l AS (SELECT * FROM events_log) SELECT l.r| FROM l",
			absent: []string{"rowid"},
		},
		{
			name:  "schema-qualified function",
			input: "SELECT public.Calc|() FROM orders",
			want:  []Candidate{{Text: `"CalculateTotal"`, Type: CandidateFunction}},
		},
		{
			name:  "pg_catalog-qualified built-in function",
			input: "SELECT pg_catalog.now| FROM orders",
			want:  []Candidate{{Text: "now()", Type: CandidateFunction}},
		},
		{
			name:  "schema-qualified type",
			input: "SELECT amount::public.st| FROM orders",
			want:  []Candidate{{Text: "status", Type: CandidateType_}},
		},
		{
			name:  "LATERAL subquery sees the FROM items before it",
			input: "SELECT * FROM orders o, LATERAL (SELECT o.| FROM analytics.events e) x",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "LATERAL join sees its left side",
			input: "SELECT * FROM orders o JOIN LATERAL (SELECT o.|) x ON true",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:   "a plain FROM subquery does not see the items before it",
			input:  "SELECT * FROM orders o, (SELECT o.|) x",
			absent: []string{"id", "amount"},
		},
		{
			name:  "a FROM subquery sees enclosing queries",
			input: "SELECT * FROM orders o WHERE EXISTS (SELECT * FROM (SELECT o.|) s)",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "a partial alias list keeps the columns it does not rename",
			input: "SELECT t.| FROM orders AS t(order_id)",
			want:  []Candidate{{Text: "order_id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "a partial CTE column list keeps the columns it does not rename",
			input: "WITH c(x) AS (SELECT id, amount FROM orders) SELECT c.| FROM c",
			want:  []Candidate{{Text: "x", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "ORDER BY offers output aliases",
			input: "SELECT amount AS total FROM orders ORDER BY tot|",
			want:  []Candidate{{Text: "total", Type: CandidateColumn}},
		},
		{
			name:   "WHERE does not offer output aliases",
			input:  "SELECT amount AS total FROM orders WHERE tot|",
			absent: []string{"total"},
		},
		{
			name:   "a nested CTE is not visible outside its query",
			input:  "WITH q AS (WITH hidden AS (SELECT 1) SELECT * FROM hidden) SELECT * FROM hi|",
			absent: []string{"hidden"},
		},
		{
			name:  "a nested CTE is visible inside its query",
			input: "WITH q AS (WITH hidden AS (SELECT 1) SELECT * FROM hi|) SELECT 1",
			want:  []Candidate{{Text: "hidden", Type: CandidateTable}},
		},
		{
			name:  "unaliased expressions take CockroachDB's column names",
			input: "WITH c AS (SELECT count(*), 1 + 1, amount::INT8 FROM orders) SELECT c.| FROM c",
			want:  []Candidate{{Text: "count", Type: CandidateColumn}, {Text: `"?column?"`, Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "an aliased parenthesized join offers both tables' columns",
			input: "SELECT j.| FROM (orders o JOIN analytics.events e ON true) AS j",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}, {Text: "event_id", Type: CandidateColumn}},
		},
		{
			name:  "an unaliased parenthesized join keeps its tables' names",
			input: "SELECT e.| FROM (orders o JOIN analytics.events e ON o.id = e.event_id)",
			want:  []Candidate{{Text: "event_id", Type: CandidateColumn}},
		},
		{
			name:  "a parenthesized join's ON sees both tables",
			input: "SELECT * FROM (orders o JOIN analytics.events e ON e.|) AS j",
			want:  []Candidate{{Text: "event_id", Type: CandidateColumn}},
		},
		{
			name:  "a DELETE CTE takes its RETURNING columns",
			input: "WITH d AS (DELETE FROM orders RETURNING id, amount AS total) SELECT d.| FROM d",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "total", Type: CandidateColumn}},
		},
		{
			name:  "an INSERT CTE expands RETURNING *",
			input: "WITH i AS (INSERT INTO orders VALUES (1, 2) RETURNING *) SELECT i.| FROM i",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}, {Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "an UPDATE statement source takes its RETURNING columns",
			input: "SELECT u.| FROM [UPDATE orders SET amount = 1 RETURNING id] AS u",
			want:  []Candidate{{Text: "id", Type: CandidateColumn}},
		},
		{
			name:  "a named window sees the query's sources",
			input: "WITH c AS (SELECT amount AS total FROM orders) SELECT sum(total) OVER w FROM c WINDOW w AS (ORDER BY c.to|)",
			want:  []Candidate{{Text: "total", Type: CandidateColumn}},
		},
		{
			name:  "a named window's partition sees the query's sources",
			input: "SELECT rank() OVER w FROM orders WINDOW w AS (PARTITION BY am|)",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "an inline window sees the query's sources",
			input: "SELECT rank() OVER (ORDER BY am|) FROM orders",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "an ON CONFLICT arbiter predicate sees the target",
			input: "INSERT INTO orders VALUES (1) ON CONFLICT (id) WHERE am| DO NOTHING",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "user-defined types are quoted",
			input: "SELECT 1::Ord|",
			want:  []Candidate{{Text: `"OrderStatus"`, Type: CandidateType_}},
		},
		{
			name:   "a new alias offers no objects",
			input:  "SELECT * FROM orders AS |",
			absent: []string{"orders", "id"},
		},
	}

	cfg := testConfig()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := completeAt(t, cfg, tc.input)
			for _, w := range tc.want {
				require.Truef(t, has(got, w.Text, w.Type), "missing %v in %v", w, got)
			}
			for _, a := range tc.absent {
				require.Falsef(t, hasText(got, a), "unexpected %q in %v", a, got)
			}
		})
	}
}

func TestCompleteSessionCockroachDBSyntax(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []Candidate
		absent []string
	}{
		{
			name:  "UPSERT starts a statement",
			input: "UPS|",
			want:  []Candidate{{Text: "UPSERT", Type: CandidateKeyword}},
		},
		{
			name:  "UPSERT INTO offers tables",
			input: "UPSERT INTO |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "UPSERT column list offers the target's columns",
			input: "UPSERT INTO orders (|) VALUES (1)",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "AS OF SYSTEM TIME",
			input: "SELECT * FROM orders AS OF |",
			want:  []Candidate{{Text: "SYSTEM", Type: CandidateKeyword}},
		},
		{
			name:  "index hint offers the table's indexes",
			input: "SELECT * FROM orders@|",
			want:  []Candidate{{Text: "orders_pkey", Type: CandidateIndex}, {Text: "orders_amount_idx", Type: CandidateIndex}},
		},
		{
			name:  "aliased index hint",
			input: "SELECT o.id FROM orders@ord| AS o",
			want:  []Candidate{{Text: "orders_amount_idx", Type: CandidateIndex}},
		},
		{
			name:  "ALTER TABLE offers CockroachDB's commands",
			input: "ALTER TABLE orders |",
			want: []Candidate{
				{Text: "CONFIGURE", Type: CandidateKeyword},
				{Text: "SCATTER", Type: CandidateKeyword},
				{Text: "SPLIT", Type: CandidateKeyword},
				{Text: "UNSPLIT", Type: CandidateKeyword},
			},
		},
		{
			name:  "CONFIGURE ZONE",
			input: "ALTER TABLE orders CONFIGURE |",
			want:  []Candidate{{Text: "ZONE", Type: CandidateKeyword}},
		},
		{
			name:  "SHOW RANGES",
			input: "SHOW RAN|",
			want:  []Candidate{{Text: "RANGES", Type: CandidateKeyword}},
		},
		{
			name:  "SHOW RANGES FROM TABLE offers tables",
			input: "SHOW RANGES FROM TABLE |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "SHOW COLUMNS FROM offers tables",
			input: "SHOW COLUMNS FROM |",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
		{
			name:  "CREATE TABLE offers IF",
			input: "CREATE TABLE |",
			want:  []Candidate{{Text: "IF", Type: CandidateKeyword}},
		},
		{
			name:   "keywords used only as names are not offered",
			input:  "SELECT * FROM orders WHERE |",
			want:   []Candidate{{Text: "NOT", Type: CandidateKeyword}, {Text: "EXISTS", Type: CandidateKeyword}},
			absent: []string{"ABORT", "ZONE"},
		},
		{
			name:  "lookahead keywords: NOT before IN",
			input: "SELECT * FROM orders WHERE id |",
			want:  []Candidate{{Text: "NOT", Type: CandidateKeyword}, {Text: "IN", Type: CandidateKeyword}},
		},
		{
			name:  "CREATE INDEX column list offers the table's columns",
			input: "CREATE INDEX ON orders (|)",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "ALTER TABLE DROP COLUMN offers the table's columns",
			input: "ALTER TABLE orders DROP COLUMN |",
			want:  []Candidate{{Text: "amount", Type: CandidateColumn}},
		},
		{
			name:  "statement source",
			input: "SELECT * FROM [SHOW COLUMNS FROM |]",
			want:  []Candidate{{Text: "orders", Type: CandidateTable}},
		},
	}

	cfg := testConfig()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := completeAt(t, cfg, tc.input)
			for _, w := range tc.want {
				require.Truef(t, has(got, w.Text, w.Type), "missing %v in %v", w, got)
			}
			for _, a := range tc.absent {
				require.Falsef(t, hasText(got, a), "unexpected %q in %v", a, got)
			}
		})
	}
}

// TestCompleteSessionKeywordSets pins the keywords offered at the table and
// expression positions users type most. A keyword the grammar only reads as
// a name or a function call there is not offered as a keyword.
func TestCompleteSessionKeywordSets(t *testing.T) {
	keywords := func(cands []Candidate) []string {
		var out []string
		for _, c := range cands {
			if c.Type == CandidateKeyword {
				out = append(out, c.Text)
			}
		}
		return out
	}
	cfg := testConfig()
	require.Equal(t, []string{"LATERAL", "ONLY", "ROWS"}, keywords(completeAt(t, cfg, "SELECT * FROM |")))
	require.Equal(t, []string{
		"ARRAY", "CASE", "CURRENT_CATALOG", "CURRENT_ROLE", "DEFAULT", "EXISTS", "FALSE", "NOT", "NULL", "TRUE", "USER",
	}, keywords(completeAt(t, cfg, "SELECT * FROM orders WHERE |")))
	// After a bare SELECT the target list may still be empty, so clause
	// keywords parse too (SELECT FROM t, SELECT WHERE true).
	require.Equal(t, []string{
		"ALL", "ARRAY", "CASE", "CURRENT_CATALOG", "CURRENT_ROLE", "DEFAULT", "DISTINCT", "EXCEPT", "EXISTS",
		"FALSE", "FETCH", "FOR", "FROM", "GROUP", "HAVING", "INTERSECT", "LIMIT", "NOT", "NULL", "OFFSET",
		"ORDER", "TRUE", "UNION", "USER", "WHERE", "WINDOW",
	}, keywords(completeAt(t, cfg, "SELECT |")))

	for _, input := range []string{"SELECT * FROM |", "SELECT |", "SELECT * FROM orders WHERE |", "SELECT * FROM orders o JOIN analytics.events e ON |"} {
		got := keywords(completeAt(t, cfg, input))
		for _, noise := range []string{"JOIN", "CROSS", "LEFT", "INNER", "CAST", "BIGINT", "INT", "STRING", "IS", "LIKE", "BETWEEN", "AUTHORIZATION", "COLLATION", "WORK", "VALUES", "CURRENT_DATE", "SUBSTRING", "COALESCE"} {
			require.NotContainsf(t, got, noise, "%s", input)
		}
	}
}

func TestCompleteSessionFunctionForms(t *testing.T) {
	cfg := testConfig()
	got := completeAt(t, cfg, "SELECT sub|")
	require.True(t, has(got, "substring()", CandidateFunction), "got %v", got)
	require.False(t, hasText(got, "SUBSTRING"), "got %v", got)
	require.False(t, hasText(got, `"substring"()`), "got %v", got)

	// CockroachDB, unlike PostgreSQL, takes session_user() as a call; the
	// pinned server answers it.
	got = completeAt(t, cfg, "SELECT sess|")
	require.Equal(t, []Candidate{{Text: "session_user()", Type: CandidateFunction}}, got)

	got = completeAt(t, cfg, "SELECT coal|")
	require.Equal(t, []Candidate{{Text: "coalesce()", Type: CandidateFunction}}, got)

	got = completeAt(t, cfg, "SELECT current_d|")
	require.True(t, has(got, "current_date()", CandidateFunction), "got %v", got)
	require.False(t, has(got, "CURRENT_DATE", CandidateKeyword), "got %v", got)

	got = completeAt(t, cfg, "SELECT substr|ing('abc', 2)")
	require.True(t, has(got, "substring", CandidateFunction), "got %v", got)

	// Functions are not offered where a table name goes.
	got = completeAt(t, cfg, "SELECT * FROM cur|")
	require.Empty(t, got)

	// A user function named by a reserved keyword is quoted.
	meta := testMetadata()
	meta.Schemas[0].Functions = append(meta.Schemas[0].Functions, &metadatapb.FunctionMetadata{Name: "order"})
	orderCfg := SessionConfig{DefaultDatabase: "db", DefaultSchema: "public", LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) { return meta, nil }}
	got = completeAt(t, orderCfg, "SELECT ord| FROM orders")
	require.True(t, has(got, `"order"()`, CandidateFunction), "got %v", got)
}

func TestCompleteSessionBuiltins(t *testing.T) {
	cfg := testConfig()
	for _, tc := range []struct {
		input  string
		want   []Candidate
		absent []string
	}{
		{input: "SELECT gen_random_u| FROM orders", want: []Candidate{{Text: "gen_random_uuid()", Type: CandidateFunction}}},
		{input: "SELECT * FROM orders WHERE amount > now|", want: []Candidate{{Text: "now()", Type: CandidateFunction}}},
		{input: "SELECT now|() FROM orders", want: []Candidate{{Text: "now", Type: CandidateFunction}}},
		{input: "SELECT o.now| FROM orders o", absent: []string{"now", "now()"}},
		{input: "CREATE TABLE t (a |", want: []Candidate{{Text: "INT8", Type: CandidateType_}, {Text: "JSONB", Type: CandidateType_}, {Text: "UUID", Type: CandidateType_}, {Text: "STRING", Type: CandidateType_}, {Text: "TIMESTAMPTZ", Type: CandidateType_}, {Text: "BYTES", Type: CandidateType_}, {Text: "status", Type: CandidateType_}}},
		{input: "CREATE TABLE t (a JS|", want: []Candidate{{Text: "JSONB", Type: CandidateType_}}, absent: []string{"orders"}},
		{input: "SELECT CAST(amount AS |) FROM orders", want: []Candidate{{Text: "DECIMAL", Type: CandidateType_}}, absent: []string{"orders", "amount"}},
		{input: "SELECT amount::ST| FROM orders", want: []Candidate{{Text: "STRING", Type: CandidateType_}}},
		{input: "ALTER TABLE orders ADD COLUMN note |", want: []Candidate{{Text: "STRING", Type: CandidateType_}}},
		{input: "ALTER TABLE orders ALTER COLUMN amount TYPE |", want: []Candidate{{Text: "FLOAT8", Type: CandidateType_}}},
	} {
		got := completeAt(t, cfg, tc.input)
		for _, w := range tc.want {
			require.Truef(t, has(got, w.Text, w.Type), "%s: missing %v in %v", tc.input, w, got)
		}
		for _, a := range tc.absent {
			require.Falsef(t, hasText(got, a), "%s: unexpected %q in %v", tc.input, a, got)
		}
		// A type keyword is offered once, as a type.
		require.Falsef(t, has(got, "STRING", CandidateKeyword) && has(got, "STRING", CandidateType_), "%s: STRING twice", tc.input)
	}
}

func TestCompleteSessionBudget(t *testing.T) {
	// A long sum collapses into nothing shorter, so keyword probing runs into
	// its budget; the columns still come back.
	stmt := "SELECT * FROM orders WHERE amount = 0" + strings.Repeat(" + 1", 4000) + " AND "
	keywordCache.reset()
	started := time.Now()
	got := CompleteSession(context.Background(), testConfig(), stmt, 1, len(stmt))
	require.Less(t, time.Since(started), 2*time.Second)
	require.True(t, has(got, "amount", CandidateColumn), "got %v", got)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got = CompleteSession(ctx, testConfig(), "SELECT * FROM orders WHERE ", 1, 27)
	require.True(t, has(got, "amount", CandidateColumn), "got %v", got)
}

func TestCompleteSessionNothingToComplete(t *testing.T) {
	cfg := testConfig()
	for _, input := range []string{
		"SELECT 'ord|",
		"SELECT 1 -- comment |",
		"SELECT 1 /* open |",
	} {
		require.Empty(t, completeAt(t, cfg, input), input)
	}
}

// benchStatements builds statements of about size bytes, each with a caret
// in the middle and one at the end: a long select list, which keyword
// probing collapses, and a long WHERE conjunction, which it cannot.
func benchStatements(size int) []benchCase {
	var list strings.Builder
	list.WriteString("SELECT ")
	i := 0
	for list.Len() < size/2 {
		fmt.Fprintf(&list, "o.amount + %d AS a%d, ", i, i)
		i++
	}
	listMid := list.Len()
	for list.Len() < size-120 {
		fmt.Fprintf(&list, "o.id * %d AS b%d, ", i, i)
		i++
	}
	list.WriteString("1 FROM orders AS o JOIN analytics.events AS e ON o.id = e.event_id WHERE o.amount > 0 AND ")

	var where strings.Builder
	where.WriteString("SELECT o.id FROM orders AS o JOIN analytics.events AS e ON o.id = e.event_id WHERE ")
	i = 0
	for where.Len() < size/2 {
		fmt.Fprintf(&where, "o.amount > %d AND ", i)
		i++
	}
	whereMid := where.Len()
	for where.Len() < size {
		fmt.Fprintf(&where, "e.event_id <> %d AND ", i)
		i++
	}
	return []benchCase{
		{"list/mid", list.String(), listMid},
		{"list/end", list.String(), list.Len()},
		{"where/mid", where.String(), whereMid},
		{"where/end", where.String(), where.Len()},
	}
}

type benchCase struct {
	name string
	stmt string
	pos  int
}

// BenchmarkCompleteSession measures a first keystroke: the keyword cache is
// emptied before every completion.
func BenchmarkCompleteSession(b *testing.B) {
	cfg := testConfig()
	for _, size := range []int{2 << 10, 11 << 10} {
		for _, bc := range benchStatements(size) {
			b.Run(fmt.Sprintf("%dB/%s", size, bc.name), func(b *testing.B) {
				for b.Loop() {
					keywordCache.reset()
					got := CompleteSession(context.Background(), cfg, bc.stmt, 1, bc.pos)
					if len(got) == 0 {
						b.Fatal("no candidates")
					}
				}
			})
		}
	}
}
