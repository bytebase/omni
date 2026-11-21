package completion

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

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

func redshiftTestMetadata() *metadatapb.DatabaseSchemaMetadata {
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
				},
				ExternalTables: []*metadatapb.ExternalTableMetadata{
					{
						Name: "spectrum_orders",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "external_id", Type: "varchar"},
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

func redshiftSessionConfig() SessionConfig {
	meta := redshiftTestMetadata()
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
	require.Equal(t, `"order"`, QuoteIdentifierIfNeeded("order", false))
	require.Equal(t, `"Order Items"`, QuoteIdentifierIfNeeded("Order Items", false))
	require.Equal(t, `"1abc"`, QuoteIdentifierIfNeeded("1abc", false))
	require.Equal(t, `"a""b"`, QuoteIdentifierIfNeeded(`a"b`, false))
	require.Equal(t, `""`, QuoteIdentifierIfNeeded("", false))

	require.Equal(t, "order", QuoteIdentifierIfNeeded("order", true))
	require.Equal(t, `a""b`, QuoteIdentifierIfNeeded(`a"b`, true))
}

func TestCompleteSession_QuotingSlotsAndMetadata(t *testing.T) {
	cfg := redshiftSessionConfig()

	stmt, line, col := catchCaretUTF16(t, "SELECT 1 FROM |")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "public", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "analytics", CandidateSchema), "got %v", got)
	require.True(t, hasSessionCandidate(got, "orders", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "spectrum_orders", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"order"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"Order Items"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "active_orders", CandidateView), "got %v", got)
	require.True(t, hasSessionCandidate(got, "orders_summary", CandidateMaterializedView), "got %v", got)
	require.True(t, hasSessionCandidate(got, "order_seq", CandidateSequence), "got %v", got)
	require.True(t, slices.IsSortedFunc(got, compareSessionCandidates))

	// Quoted columns are surfaced quoted.
	stmtCol, lineCol, colCol := catchCaretUTF16(t, `SELECT | FROM "Order Items"`)
	gotCol := CompleteSession(context.Background(), cfg, stmtCol, lineCol, colCol)
	require.True(t, hasSessionCandidate(gotCol, `"Item ID"`, CandidateColumn), "got %v", gotCol)
	require.True(t, hasSessionCandidate(gotCol, `"Order ID"`, CandidateColumn), "got %v", gotCol)

	// Inside double quotes, outer quotes are omitted.
	stmtQ, lineQ, colQ := catchCaretUTF16(t, `SELECT * FROM "ord|"`)
	gotQ := CompleteSession(context.Background(), cfg, stmtQ, lineQ, colQ)
	require.True(t, hasSessionCandidate(gotQ, "order", CandidateTable), "got %v", gotQ)
	require.True(t, hasSessionCandidate(gotQ, "orders", CandidateTable), "got %v", gotQ)

	// Redshift-specific slot keywords remain available.
	stmtCopy, lineCopy, colCopy := catchCaretUTF16(t, "COPY orders FROM 's3://bucket/file' |")
	gotCopy := CompleteSession(context.Background(), cfg, stmtCopy, lineCopy, colCopy)
	require.True(t, hasSessionCandidate(gotCopy, "IAM_ROLE", CandidateKeyword), "got %v", gotCopy)
	require.True(t, hasSessionCandidate(gotCopy, "FORMAT", CandidateKeyword), "got %v", gotCopy)
}

func TestCompleteSession_StatementIsolationAndUTF16(t *testing.T) {
	cfg := redshiftSessionConfig()

	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM orders; SELECT | FROM analytics.events")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "event_id", CandidateColumn), "got %v", got)
	require.False(t, hasSessionCandidate(got, "amount", CandidateColumn), "got %v", got)

	stmtEmoji, lineE, colE := catchCaretUTF16(t, "SELECT '😀', | FROM orders")
	gotE := CompleteSession(context.Background(), cfg, stmtEmoji, lineE, colE)
	require.True(t, hasSessionCandidate(gotE, "id", CandidateColumn), "got %v", gotE)
	require.True(t, hasSessionCandidate(gotE, "amount", CandidateColumn), "got %v", gotE)

	// Non-public DefaultSchema resolves analytics tables unqualified.
	analyticsCfg := redshiftSessionConfig()
	analyticsCfg.DefaultSchema = "analytics"
	stmtA, lineA, colA := catchCaretUTF16(t, "SELECT | FROM events")
	gotA := CompleteSession(context.Background(), analyticsCfg, stmtA, lineA, colA)
	require.True(t, hasSessionCandidate(gotA, "event_id", CandidateColumn), "got %v", gotA)
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
	}
}
