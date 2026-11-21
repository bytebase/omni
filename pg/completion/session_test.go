package completion

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	metadatapb "github.com/bytebase/omni/metadata"
	"github.com/stretchr/testify/require"
)

func catchCaretUTF16(t *testing.T, input string) (string, int, int) {
	t.Helper()
	line := 1
	col := 0
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

func hasSessionCandidate(cands []Candidate, text string, typ CandidateType) bool {
	for _, c := range cands {
		if c.Text == text && c.Type == typ {
			return true
		}
	}
	return false
}

func pgTestMetadata() *metadatapb.DatabaseSchemaMetadata {
	return &metadatapb.DatabaseSchemaMetadata{
		Name: "db",
		Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "orders",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
							{Name: "amount", Type: "numeric"},
						},
					},
					{
						Name: "order",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
						},
					},
					{
						Name: "Order Items",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "Item ID", Type: "bigint"},
							{Name: "Order ID", Type: "integer"},
						},
					},
					{
						Name: "my-table",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "col1", Type: "integer"},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{
						Name: "active_orders",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "id", Type: "integer"},
							{Name: "amount", Type: "numeric"},
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
					{Name: "CalculateTotal", Definition: "CREATE FUNCTION public.\"CalculateTotal\"() RETURNS integer LANGUAGE sql AS $$ SELECT 1 $$"},
				},
			},
			{
				Name: "analytics",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "events",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "event_id", Type: "bigint"},
						},
					},
				},
			},
		},
	}
}

func pgSessionConfig() SessionConfig {
	meta := pgTestMetadata()
	return SessionConfig{
		DefaultDatabase: "db",
		DefaultSchema:   "public",
		LoadDatabase: func(context.Context, string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return meta, nil
		},
	}
}

func TestLineOffsetToBytePos(t *testing.T) {
	input := "SELECT '😀', | FROM orders"
	stmt, line, col := catchCaretUTF16(t, input)
	require.Equal(t, strings.Index(input, "|"), LineOffsetToBytePos(stmt, line, col))
	require.Equal(t, 3, LineOffsetToBytePos("abc\ndef", 1, 99))
	require.Equal(t, 7, LineOffsetToBytePos("abc\ndef", 9, 0))
}

func TestCaretInsideDoubleQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "unquoted", input: "SELECT * FROM ord|", want: false},
		{name: "inside double quotes", input: `SELECT * FROM "ord|"`, want: true},
		{name: "after closed double quotes", input: `SELECT * FROM "ord"|`, want: false},
		{name: "quote inside comment", input: `SELECT * FROM /* " */ |`, want: false},
		{name: "quote inside dollar-quoted string", input: `SELECT $tag$"$tag$, | FROM orders`, want: false},
		{name: "quote inside single-quoted string", input: `SELECT '"' , | FROM orders`, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmt, line, col := catchCaretUTF16(t, tc.input)
			pos := LineOffsetToBytePos(stmt, line, col)
			require.Equal(t, tc.want, CaretInsideDoubleQuotes(stmt, pos))
		})
	}
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

	require.Equal(t, "order", QuoteIdentifierIfNeeded("order", true))
	require.Equal(t, `a""b`, QuoteIdentifierIfNeeded(`a"b`, true))
}

func TestCompleteSession_QuotingSequencesAndSearchPath(t *testing.T) {
	cfg := pgSessionConfig()

	stmt, line, col := catchCaretUTF16(t, "SELECT 1 FROM |")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "public", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "analytics", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "orders", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"order"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"Order Items"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"my-table"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "active_orders", CandidateView), "got %v", got)
	require.True(t, hasSessionCandidate(got, "orders_summary", CandidateMaterializedView), "got %v", got)
	require.True(t, hasSessionCandidate(got, "order_seq", CandidateSequence), "got %v", got)
	require.True(t, slices.IsSortedFunc(got, compareSessionCandidates))

	// Inside double quotes, outer quotes are omitted.
	stmtQ, lineQ, colQ := catchCaretUTF16(t, `SELECT * FROM "ord|"`)
	gotQ := CompleteSession(context.Background(), cfg, stmtQ, lineQ, colQ)
	require.True(t, hasSessionCandidate(gotQ, "order", CandidateTable), "got %v", gotQ)
	require.True(t, hasSessionCandidate(gotQ, "orders", CandidateTable), "got %v", gotQ)

	// Quote in comment or dollar-quoted string does not suppress quoting.
	stmtComment, lineCm, colCm := catchCaretUTF16(t, `SELECT * FROM /* " */ |`)
	gotCm := CompleteSession(context.Background(), cfg, stmtComment, lineCm, colCm)
	require.True(t, hasSessionCandidate(gotCm, `"Order Items"`, CandidateTable), "got %v", gotCm)

	stmtDollar, lineD, colD := catchCaretUTF16(t, `SELECT $tag$"$tag$, | FROM "Order Items"`)
	gotD := CompleteSession(context.Background(), cfg, stmtDollar, lineD, colD)
	require.True(t, hasSessionCandidate(gotD, `"Item ID"`, CandidateColumn), "got %v", gotD)

	// Non-public DefaultSchema is prepended to search_path.
	analyticsCfg := pgSessionConfig()
	analyticsCfg.DefaultSchema = "analytics"
	stmtA, lineA, colA := catchCaretUTF16(t, "SELECT | FROM events")
	gotA := CompleteSession(context.Background(), analyticsCfg, stmtA, lineA, colA)
	require.True(t, hasSessionCandidate(gotA, "event_id", CandidateColumn), "got %v", gotA)
}

func TestCompleteSession_StatementIsolationAndTrailingBrokenSQL(t *testing.T) {
	cfg := pgSessionConfig()

	// Earlier statement does not leak columns.
	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM orders; SELECT | FROM analytics.events")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "event_id", CandidateColumn), "got %v", got)
	require.False(t, hasSessionCandidate(got, "amount", CandidateColumn), "got %v", got)

	// Caret immediately after semicolon begins a fresh statement.
	stmtSemi, lineS, colS := catchCaretUTF16(t, "SELECT * FROM orders;| SELECT 1")
	gotSemi := CompleteSession(context.Background(), cfg, stmtSemi, lineS, colS)
	require.True(t, hasSessionCandidate(gotSemi, "SELECT", CandidateKeyword), "got %v", gotSemi)

	// Caret immediately after the final semicolon (no trailing space) starts an empty scope.
	stmtFinalSemi, lineFS, colFS := catchCaretUTF16(t, "SELECT 1;|")
	gotFinalSemi := CompleteSession(context.Background(), cfg, stmtFinalSemi, lineFS, colFS)
	require.True(t, hasSessionCandidate(gotFinalSemi, "SELECT", CandidateKeyword), "got %v", gotFinalSemi)

	// Quoted user-defined functions are quoted before appending "()", and omit "()" in DDL contexts.
	require.Equal(t, `"CalculateTotal"()`, formatSessionCandidateText(Candidate{Text: "CalculateTotal", Type: CandidateFunction}, false, true))
	require.Equal(t, `"CalculateTotal"`, formatSessionCandidateText(Candidate{Text: "CalculateTotal", Type: CandidateFunction}, false, false))
	require.Equal(t, "CalculateTotal", formatSessionCandidateText(Candidate{Text: "CalculateTotal", Type: CandidateFunction}, true, true))
	require.Equal(t, "my_func()", formatSessionCandidateText(Candidate{Text: "my_func", Type: CandidateFunction}, false, true))

	// UTF-16 caret after emoji resolves columns accurately.
	stmtEmoji, lineE, colE := catchCaretUTF16(t, "SELECT '😀', | FROM orders")
	gotE := CompleteSession(context.Background(), cfg, stmtEmoji, lineE, colE)
	require.True(t, hasSessionCandidate(gotE, "id", CandidateColumn), "got %v", gotE)
	require.True(t, hasSessionCandidate(gotE, "amount", CandidateColumn), "got %v", gotE)

	// Broken trailing statements must not stall completion.
	var sheet strings.Builder
	sheet.WriteString("SELECT  FROM orders JOIN analytics.events a2 ON orders.id = a2.event_id;\n")
	sheet.WriteString("SELEC broken FROM oops;\n")
	for i := range 2000 {
		fmt.Fprintf(&sheet, "SELECT col_a, col_b FROM table_%04d WHERE col_a = %d;\n", i, i)
	}
	started := time.Now()
	gotBig := CompleteSession(context.Background(), cfg, sheet.String(), 1, 7)
	require.Less(t, time.Since(started), 2*time.Second)
	require.True(t, hasSessionCandidate(gotBig, "id", CandidateColumn), "got %v", gotBig)
	require.True(t, hasSessionCandidate(gotBig, "event_id", CandidateColumn), "got %v", gotBig)
}

func TestCompleteSession_NilOrFailedMetadata(t *testing.T) {
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
		stmt, line, col := catchCaretUTF16(t, "SELECT 1 FROM |")
		for _, c := range CompleteSession(context.Background(), cfg, stmt, line, col) {
			require.NotContains(t, []CandidateType{
				CandidateSchema, CandidateTable, CandidateView,
				CandidateMaterializedView, CandidateColumn, CandidateSequence,
			}, c.Type)
		}

		// Statement-local scopes (CTEs and alias columns) remain available without metadata.
		stmtCTE, lineC, colC := catchCaretUTF16(t, "WITH c(x) AS (SELECT 1) SELECT | FROM c")
		gotCTE := CompleteSession(context.Background(), cfg, stmtCTE, lineC, colC)
		require.True(t, hasSessionCandidate(gotCTE, "x", CandidateColumn), "got %v", gotCTE)

		stmtCTERel, lineR, colR := catchCaretUTF16(t, "WITH c(x) AS (SELECT 1) SELECT * FROM |")
		gotCTERel := CompleteSession(context.Background(), cfg, stmtCTERel, lineR, colR)
		require.True(t, hasSessionCandidate(gotCTERel, "c", CandidateTable), "got %v", gotCTERel)
	}
}
