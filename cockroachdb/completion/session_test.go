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
