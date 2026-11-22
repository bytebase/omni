package completion

import (
	"context"
	"errors"
	"fmt"
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

	// Caret immediately after the final semicolon (no trailing space) starts an empty scope.
	stmtFinalSemi, lineFS, colFS := catchCaretUTF16(t, "SELECT 1;|")
	gotFinalSemi := CompleteSession(context.Background(), cfg, stmtFinalSemi, lineFS, colFS)
	require.True(t, hasSessionCandidate(gotFinalSemi, "SELECT", CandidateKeyword), "got %v", gotFinalSemi)

	// Quoted user-defined functions are quoted when needed.
	require.Equal(t, `"CalculateTotal"`, formatSessionCandidateText(Candidate{Text: "CalculateTotal", Type: CandidateFunction}, false))
	require.Equal(t, "CalculateTotal", formatSessionCandidateText(Candidate{Text: "CalculateTotal", Type: CandidateFunction}, true))

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

		// Statement-local scopes (CTEs and alias columns) remain available without metadata.
		stmtCTE, lineC, colC := catchCaretUTF16(t, "WITH c(x) AS (SELECT 1) SELECT | FROM c")
		gotCTE := CompleteSession(context.Background(), cfg, stmtCTE, lineC, colC)
		require.True(t, hasSessionCandidate(gotCTE, "x", CandidateColumn), "got %v", gotCTE)

		stmtCTERel, lineR, colR := catchCaretUTF16(t, "WITH c(x) AS (SELECT 1) SELECT * FROM |")
		gotCTERel := CompleteSession(context.Background(), cfg, stmtCTERel, lineR, colR)
		require.True(t, hasSessionCandidate(gotCTERel, "c", CandidateTable), "got %v", gotCTERel)
	}
}

func TestCompleteSession_RedshiftSlotsAndCoverageMatrix(t *testing.T) {
	type wantSpec struct {
		text string
		typ  CandidateType
	}
	type matrixCase struct {
		name       string
		input      string
		want       []wantSpec
		absentText []string
	}

	// Redshift grammar slot and prefix cases (no metadata needed).
	slotCases := []matrixCase{
		{name: "statement start", input: "|", want: []wantSpec{{"SELECT", CandidateKeyword}, {"COPY", CandidateKeyword}, {"UNLOAD", CandidateKeyword}, {"SHOW", CandidateKeyword}}},
		{name: "prefix filtering", input: "SH|", want: []wantSpec{{"SHOW", CandidateKeyword}}, absentText: []string{"SELECT", "COPY"}},
		{name: "create dispatch", input: "CREATE |", want: []wantSpec{{"TABLE", CandidateKeyword}, {"VIEW", CandidateKeyword}, {"SCHEMA", CandidateKeyword}, {"DATABASE", CandidateKeyword}}},
		{name: "alter dispatch", input: "ALTER |", want: []wantSpec{{"TABLE", CandidateKeyword}, {"DATABASE", CandidateKeyword}, {"ROLE", CandidateKeyword}, {"SCHEMA", CandidateKeyword}}},
		{name: "drop dispatch", input: "DROP |", want: []wantSpec{{"TABLE", CandidateKeyword}, {"VIEW", CandidateKeyword}, {"INDEX", CandidateKeyword}, {"FUNCTION", CandidateKeyword}}},
		{name: "create table options", input: "CREATE TABLE t (id INT) |", want: []wantSpec{{"DISTSTYLE", CandidateKeyword}, {"DISTKEY", CandidateKeyword}, {"SORTKEY", CandidateKeyword}, {"ENCODE", CandidateKeyword}}},
		{name: "create table option prefix DIS", input: "CREATE TABLE t (id INT) DIS|", want: []wantSpec{{"DISTSTYLE", CandidateKeyword}, {"DISTKEY", CandidateKeyword}}, absentText: []string{"SORTKEY", "ENCODE"}},
		{name: "create table option prefix D", input: "CREATE TABLE t (id INT) D|", want: []wantSpec{{"DISTSTYLE", CandidateKeyword}}, absentText: []string{"ENCODE"}},
		{name: "create table option prefix DI", input: "CREATE TABLE t (id INT) DI|", want: []wantSpec{{"DISTKEY", CandidateKeyword}}, absentText: []string{"ENCODE"}},
		{name: "create table option prefix S", input: "CREATE TABLE t (id INT) S|", want: []wantSpec{{"SORTKEY", CandidateKeyword}}, absentText: []string{"ENCODE"}},
		{name: "create table option prefix E", input: "CREATE TABLE t (id INT) E|", want: []wantSpec{{"ENCODE", CandidateKeyword}}, absentText: []string{"SORTKEY"}},
		{name: "copy options", input: "COPY t FROM 's3://bucket/file' |", want: []wantSpec{{"IAM_ROLE", CandidateKeyword}, {"CREDENTIALS", CandidateKeyword}, {"FORMAT", CandidateKeyword}, {"MANIFEST", CandidateKeyword}}},
		{name: "copy option prefix C", input: "COPY orders FROM 's3://bucket/file' C|", want: []wantSpec{{"CREDENTIALS", CandidateKeyword}}, absentText: []string{"FORMAT"}},
		{name: "copy option prefix F", input: "COPY orders FROM 's3://bucket/file' F|", want: []wantSpec{{"FORMAT", CandidateKeyword}}, absentText: []string{"MANIFEST"}},
		{name: "copy option prefix M", input: "COPY orders FROM 's3://bucket/file' M|", want: []wantSpec{{"MANIFEST", CandidateKeyword}}, absentText: []string{"FORMAT"}},
		{name: "copy option prefix G", input: "COPY orders FROM 's3://bucket/file' G|", want: []wantSpec{{"GZIP", CandidateKeyword}}, absentText: []string{"BZIP2", "MANIFEST"}},
		{name: "copy option prefix B", input: "COPY orders FROM 's3://bucket/file' B|", want: []wantSpec{{"BZIP2", CandidateKeyword}}, absentText: []string{"GZIP"}},
		{name: "unload options", input: "UNLOAD ('SELECT * FROM t') TO 's3://bucket/out' |", want: []wantSpec{{"IAM_ROLE", CandidateKeyword}, {"FORMAT", CandidateKeyword}, {"MANIFEST", CandidateKeyword}, {"HEADER", CandidateKeyword}}},
		{name: "unload option prefix F", input: "UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' F|", want: []wantSpec{{"FORMAT", CandidateKeyword}}, absentText: []string{"MANIFEST"}},
		{name: "unload option prefix M", input: "UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' M|", want: []wantSpec{{"MANIFEST", CandidateKeyword}}, absentText: []string{"FORMAT"}},
		{name: "unload option prefix G", input: "UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' G|", want: []wantSpec{{"GZIP", CandidateKeyword}}, absentText: []string{"BZIP2"}},
		{name: "unload option prefix B", input: "UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' B|", want: []wantSpec{{"BZIP2", CandidateKeyword}}, absentText: []string{"GZIP"}},
		{name: "unload option prefix H", input: "UNLOAD ('SELECT * FROM orders') TO 's3://bucket/out' H|", want: []wantSpec{{"HEADER", CandidateKeyword}}, absentText: []string{"MANIFEST", "FORMAT"}},
		{name: "show subcommands", input: "SHOW |", want: []wantSpec{{"DATABASES", CandidateKeyword}, {"SCHEMAS", CandidateKeyword}, {"TABLES", CandidateKeyword}, {"DATASHARES", CandidateKeyword}}},
		{name: "show subcommand prefix DATA", input: "SHOW DATA|", want: []wantSpec{{"DATABASES", CandidateKeyword}, {"DATASHARES", CandidateKeyword}}, absentText: []string{"TABLES", "COLUMNS"}},
		{name: "show subcommand prefix D", input: "SHOW D|", want: []wantSpec{{"DATASHARES", CandidateKeyword}}, absentText: []string{"TABLES"}},
		{name: "show subcommand prefix S", input: "SHOW S|", want: []wantSpec{{"SCHEMAS", CandidateKeyword}}, absentText: []string{"TABLES"}},
		{name: "show subcommand prefix T", input: "SHOW T|", want: []wantSpec{{"TABLES", CandidateKeyword}}, absentText: []string{"SCHEMAS"}},
		{name: "show subcommand prefix C", input: "SHOW C|", want: []wantSpec{{"COLUMNS", CandidateKeyword}}, absentText: []string{"DATASHARES"}},
		{name: "show subcommand prefix G", input: "SHOW G|", want: []wantSpec{{"GRANTS", CandidateKeyword}}, absentText: []string{"TABLES"}},
	}

	// Metadata relation, schema-qualified, column, qualifier, prefix, and CTE cases.
	metaCases := []matrixCase{
		// Relation contexts
		{name: "select from", input: "SELECT * FROM |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "select join", input: "SELECT * FROM orders JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "select comma join", input: "SELECT * FROM orders, |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "left join", input: "SELECT * FROM orders LEFT JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "right join", input: "SELECT * FROM orders RIGHT JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "full join", input: "SELECT * FROM orders FULL JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "cross join", input: "SELECT * FROM orders CROSS JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "insert into", input: "INSERT INTO |", want: []wantSpec{{"orders", CandidateTable}, {"events", CandidateTable}}},
		{name: "copy relation", input: "COPY | FROM 's3://bucket/file'", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "update relation", input: "UPDATE |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "delete relation", input: "DELETE FROM |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "truncate relation", input: "TRUNCATE TABLE |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "create table as from", input: "CREATE TABLE new_orders AS SELECT * FROM |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "create view as from", input: "CREATE VIEW v AS SELECT * FROM |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "create materialized view as from", input: "CREATE MATERIALIZED VIEW mv AS SELECT * FROM |", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "subquery from", input: "SELECT * FROM (SELECT * FROM |) s", want: []wantSpec{{"orders", CandidateTable}}},
		{name: "exists subquery from", input: "SELECT * FROM orders WHERE EXISTS (SELECT 1 FROM |)", want: []wantSpec{{"events", CandidateTable}}},
		{name: "in subquery from", input: "SELECT * FROM orders WHERE id IN (SELECT event_id FROM |)", want: []wantSpec{{"events", CandidateTable}}},
		{name: "union right from", input: "SELECT * FROM orders UNION ALL SELECT * FROM |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "with body from", input: "WITH x AS (SELECT * FROM |) SELECT * FROM x", want: []wantSpec{{"orders", CandidateTable}}},

		// Column contexts
		{name: "projection", input: "SELECT | FROM orders", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "projection after expression", input: "SELECT count(*), | FROM orders", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "where", input: "SELECT * FROM orders WHERE |", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "where after predicate", input: "SELECT * FROM orders WHERE amount > 10 AND |", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "group by", input: "SELECT id, sum(amount) FROM orders GROUP BY |", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "order by", input: "SELECT id, amount FROM orders ORDER BY |", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "having", input: "SELECT id, sum(amount) FROM orders GROUP BY id HAVING |", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "join on left", input: "SELECT * FROM orders o JOIN analytics.events e ON |", want: []wantSpec{{"id", CandidateColumn}, {"event_id", CandidateColumn}}},
		{name: "join on right predicate", input: "SELECT * FROM orders o JOIN analytics.events e ON o.id = e.event_id AND |", want: []wantSpec{{"id", CandidateColumn}, {"event_id", CandidateColumn}}},
		{name: "subquery projection", input: "SELECT * FROM (SELECT | FROM orders) s", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "exists projection", input: "SELECT * FROM orders WHERE EXISTS (SELECT | FROM analytics.events)", want: []wantSpec{{"event_id", CandidateColumn}}},
		{name: "with body projection", input: "WITH x AS (SELECT | FROM orders) SELECT * FROM x", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "view projection", input: "SELECT | FROM active_orders", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "materialized view projection", input: "SELECT | FROM orders_summary", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}},
		{name: "external table projection", input: "SELECT | FROM spectrum_orders", want: []wantSpec{{"external_id", CandidateColumn}}},

		// Qualifier contexts
		{name: "alias o", input: "SELECT o.| FROM orders AS o", want: []wantSpec{{"id", CandidateColumn}, {"amount", CandidateColumn}}, absentText: []string{"event_id"}},
		{name: "alias e", input: "SELECT e.| FROM orders o JOIN analytics.events e ON o.id = e.event_id", want: []wantSpec{{"event_id", CandidateColumn}}},
		{name: "table alias named orders", input: "SELECT orders.| FROM orders AS orders", want: []wantSpec{{"id", CandidateColumn}}},
		{name: "table events", input: "SELECT events.| FROM analytics.events", want: []wantSpec{{"event_id", CandidateColumn}}, absentText: []string{"amount"}},
		{name: "cte alias", input: "WITH x(a, b) AS (SELECT id, amount FROM orders) SELECT x.| FROM x", want: []wantSpec{{"a", CandidateColumn}, {"b", CandidateColumn}}},
		{name: "unknown qualifier", input: "SELECT nope.| FROM orders", absentText: []string{"id", "amount", "event_id"}},

		// Prefix contexts
		{name: "relation prefix ord", input: "SELECT * FROM ord|", want: []wantSpec{{"orders", CandidateTable}, {"orders_summary", CandidateMaterializedView}}, absentText: []string{"events", "active_orders"}},
		{name: "relation prefix act", input: "SELECT * FROM act|", want: []wantSpec{{"active_orders", CandidateView}}, absentText: []string{"orders"}},
		{name: "relation prefix spec", input: "SELECT * FROM spec|", want: []wantSpec{{"spectrum_orders", CandidateTable}}, absentText: []string{"active_orders"}},
		{name: "schema prefix ana", input: "SELECT * FROM ana|", want: []wantSpec{{"analytics", CandidateSchema}}, absentText: []string{"public"}},
		{name: "column prefix am", input: "SELECT am| FROM orders", want: []wantSpec{{"amount", CandidateColumn}}, absentText: []string{"id"}},
		{name: "column prefix event", input: "SELECT event| FROM analytics.events", want: []wantSpec{{"event_id", CandidateColumn}}, absentText: []string{"amount"}},

		// CTE contexts
		{name: "cte from", input: "WITH x AS (SELECT * FROM orders) SELECT * FROM |", want: []wantSpec{{"x", CandidateTable}}},
		{name: "cte projection alias cols", input: "WITH x(a, b) AS (SELECT id, amount FROM orders) SELECT | FROM x", want: []wantSpec{{"a", CandidateColumn}, {"b", CandidateColumn}}},
		{name: "cte prefix", input: "WITH recent_orders AS (SELECT * FROM orders) SELECT * FROM recent_|", want: []wantSpec{{"recent_orders", CandidateTable}}},
		{name: "cte join relation", input: "WITH x AS (SELECT * FROM orders) SELECT * FROM x JOIN |", want: []wantSpec{{"events", CandidateTable}}},
		{name: "cte join columns", input: "WITH x(a, b) AS (SELECT id, amount FROM orders) SELECT | FROM x JOIN analytics.events e ON x.a = e.event_id", want: []wantSpec{{"a", CandidateColumn}, {"event_id", CandidateColumn}}},
	}

	// Schema-qualified relation templates across both schemas.
	qualifiedTemplates := []string{
		"SELECT * FROM %s.|",
		"SELECT * FROM orders JOIN %s.|",
		"INSERT INTO %s.|",
		"UPDATE %s.|",
		"DELETE FROM %s.|",
		"COPY %s.| FROM 's3://bucket/file'",
		"CREATE TABLE copy AS SELECT * FROM %s.|",
		"CREATE VIEW copy_v AS SELECT * FROM %s.|",
		"SELECT * FROM (SELECT * FROM %s.|) s",
		"WITH x AS (SELECT * FROM %s.|) SELECT * FROM x",
	}
	for _, s := range []struct {
		schema string
		want   string
		absent []string
	}{
		{"public", "orders", []string{"events"}},
		{"analytics", "events", []string{"orders", "active_orders"}},
	} {
		for i, tmpl := range qualifiedTemplates {
			metaCases = append(metaCases, matrixCase{
				name:       fmt.Sprintf("schema qualified %s %02d", s.schema, i),
				input:      fmt.Sprintf(tmpl, s.schema),
				want:       []wantSpec{{s.want, CandidateTable}},
				absentText: s.absent,
			})
		}
	}

	runCases := func(cfg SessionConfig, cases []matrixCase) {
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				stmt, line, col := catchCaretUTF16(t, tc.input)
				got := CompleteSession(context.Background(), cfg, stmt, line, col)
				for _, w := range tc.want {
					require.Truef(t, hasSessionCandidate(got, w.text, w.typ), "missing {%s, %v} in %v", w.text, w.typ, got)
				}
				for _, absent := range tc.absentText {
					for _, c := range got {
						require.NotEqualf(t, absent, c.Text, "unexpected %q in %v", absent, got)
					}
				}
			})
		}
	}

	runCases(SessionConfig{DefaultDatabase: "db", DefaultSchema: "public"}, slotCases)
	runCases(redshiftSessionConfig(), metaCases)
}
