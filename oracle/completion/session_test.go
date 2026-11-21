package completion

import (
	"context"
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

func oracleTestSchemas() map[string]*metadatapb.DatabaseSchemaMetadata {
	return map[string]*metadatapb.DatabaseSchemaMetadata{
		"SCHEMA1": {
			Name: "SCHEMA1",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "T1",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "C1", Type: "NUMBER"},
						},
					},
					{
						Name: "T2",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "C1", Type: "NUMBER"},
							{Name: "C2", Type: "VARCHAR2(100)"},
						},
					},
					{
						Name: "lower_tbl",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "MixedCol", Type: "NUMBER"},
							{Name: "SELECT", Type: "NUMBER"},
							{Name: `has"quote`, Type: "NUMBER"},
						},
					},
				},
				Views: []*metadatapb.ViewMetadata{
					{Name: "V1", Definition: "CREATE VIEW V1 AS SELECT * FROM T1"},
				},
				Sequences: []*metadatapb.SequenceMetadata{
					{Name: "SEQ1"},
				},
				Functions: []*metadatapb.FunctionMetadata{
					{Name: "GET_TOTAL"},
					{Name: "CalculateTotal"},
				},
			}},
		},
		"SCHEMA2": {
			Name: "SCHEMA2",
			Schemas: []*metadatapb.SchemaMetadata{{
				Name: "",
				Tables: []*metadatapb.TableMetadata{
					{
						Name: "REMOTE_T",
						Columns: []*metadatapb.ColumnMetadata{
							{Name: "R1", Type: "NUMBER"},
						},
					},
				},
			}},
		},
	}
}

func TestLineOffsetToBytePos(t *testing.T) {
	input := "SELECT '😀' FROM t1 WHERE | AND c1 = 1"
	stmt, line, col := catchCaretUTF16(t, input)
	require.Equal(t, strings.Index(input, "|"), LineOffsetToBytePos(stmt, line, col))

	multiline := "SELECT 1;\nSELECT '😀', | FROM t1"
	stmt2, line2, col2 := catchCaretUTF16(t, multiline)
	require.Equal(t, strings.Index(multiline, "|"), LineOffsetToBytePos(stmt2, line2, col2))

	// Column past end of line clamps to newline.
	require.Equal(t, 3, LineOffsetToBytePos("abc\ndef", 1, 99))
	// Line past end of input clamps to len(s).
	require.Equal(t, 7, LineOffsetToBytePos("abc\ndef", 9, 0))
}

func TestCaretInsideDoubleQuotes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "q-quote with apostrophe and double quote",
			input: `SELECT q'[a'b"c]' FROM |`,
			want:  false,
		},
		{
			name:  "nq-quote with apostrophe and double quote",
			input: `SELECT nq'{a'b"c}' FROM |`,
			want:  false,
		},
		{
			name:  "q-quote with angle brackets",
			input: `SELECT Q'<a"b>' FROM |`,
			want:  false,
		},
		{
			name:  "q-quote followed by open quoted identifier",
			input: `SELECT q'[a'b"c]' FROM "my|`,
			want:  true,
		},
		{
			name:  "escaped double quote inside open quoted identifier",
			input: `SELECT "a""b| FROM t1`,
			want:  true,
		},
		{
			name:  "closed double-quoted identifier",
			input: `SELECT "a""b"| FROM t1`,
			want:  false,
		},
		{
			name:  "double quote inside single-quoted string",
			input: `SELECT 'a"b' FROM |`,
			want:  false,
		},
		{
			name:  "double quote inside line comment",
			input: "SELECT -- \"\n| FROM t1",
			want:  false,
		},
		{
			name:  "double quote inside block comment",
			input: `SELECT /* " */ | FROM t1`,
			want:  false,
		},
		{
			name:  "double quote inside unclosed block comment",
			input: `SELECT /* unclosed " |`,
			want:  false,
		},
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
	require.Equal(t, "T1", QuoteIdentifierIfNeeded("T1", false))
	require.Equal(t, "A$B#1", QuoteIdentifierIfNeeded("A$B#1", false))
	require.Equal(t, `"lower_tbl"`, QuoteIdentifierIfNeeded("lower_tbl", false))
	require.Equal(t, `"MixedCol"`, QuoteIdentifierIfNeeded("MixedCol", false))
	require.Equal(t, `"SELECT"`, QuoteIdentifierIfNeeded("SELECT", false))
	require.Equal(t, `"A B"`, QuoteIdentifierIfNeeded("A B", false))
	require.Equal(t, `"1ABC"`, QuoteIdentifierIfNeeded("1ABC", false))
	require.Equal(t, `"has""quote"`, QuoteIdentifierIfNeeded(`has"quote`, false))
	require.Equal(t, `""`, QuoteIdentifierIfNeeded("", false))

	// Inside double quotes: omit outer quotes, double embedded quotes.
	require.Equal(t, "lower_tbl", QuoteIdentifierIfNeeded("lower_tbl", true))
	require.Equal(t, `has""quote`, QuoteIdentifierIfNeeded(`has"quote`, true))
}

func TestCompleteSession_QuotingAndFunctions(t *testing.T) {
	schemas := oracleTestSchemas()
	cfg := SessionConfig{
		DefaultSchema: "SCHEMA1",
		ListSchemaNames: func(context.Context) ([]string, error) {
			return []string{"SCHEMA1", "SCHEMA2"}, nil
		},
		LoadSchema: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return schemas[strings.ToUpper(name)], nil
		},
	}

	// Lowercase table is quoted when caret is outside double quotes.
	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM |")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "T1", CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, `"lower_tbl"`, CandidateTable), "got %v", got)
	require.True(t, hasSessionCandidate(got, "V1", CandidateView), "got %v", got)
	require.True(t, hasSessionCandidate(got, "SEQ1", CandidateSequence), "got %v", got)
	require.True(t, hasSessionCandidate(got, "SCHEMA2", CandidateSchema), "got %v", got)

	// Inside double quotes, outer quotes are omitted.
	stmtQ, lineQ, colQ := catchCaretUTF16(t, `SELECT * FROM "low|"`)
	gotQ := CompleteSession(context.Background(), cfg, stmtQ, lineQ, colQ)
	require.True(t, hasSessionCandidate(gotQ, "lower_tbl", CandidateTable), "got %v", gotQ)
	require.False(t, hasSessionCandidate(gotQ, `"lower_tbl"`, CandidateTable), "got %v", gotQ)

	// Quoted columns with mixed case, reserved keyword, and embedded quote.
	stmtC, lineC, colC := catchCaretUTF16(t, `SELECT | FROM "lower_tbl"`)
	gotC := CompleteSession(context.Background(), cfg, stmtC, lineC, colC)
	require.True(t, hasSessionCandidate(gotC, `"MixedCol"`, CandidateColumn), "got %v", gotC)
	require.True(t, hasSessionCandidate(gotC, `"SELECT"`, CandidateColumn), "got %v", gotC)
	require.True(t, hasSessionCandidate(gotC, `"has""quote"`, CandidateColumn), "got %v", gotC)

	// Functions receive trailing "()"; mixed-case/special metadata functions are quoted first.
	stmtF, lineF, colF := catchCaretUTF16(t, "SELECT GET_T| FROM T1")
	gotF := CompleteSession(context.Background(), cfg, stmtF, lineF, colF)
	require.True(t, hasSessionCandidate(gotF, "GET_TOTAL()", CandidateFunction), "got %v", gotF)

	stmtFQ, lineFQ, colFQ := catchCaretUTF16(t, "SELECT Calc| FROM T1")
	gotFQ := CompleteSession(context.Background(), cfg, stmtFQ, lineFQ, colFQ)
	require.True(t, hasSessionCandidate(gotFQ, `"CalculateTotal"()`, CandidateFunction), "got %v", gotFQ)

	stmtFInQ, lineFIQ, colFIQ := catchCaretUTF16(t, `SELECT "Calc|" FROM T1`)
	gotFInQ := CompleteSession(context.Background(), cfg, stmtFInQ, lineFIQ, colFIQ)
	require.True(t, hasSessionCandidate(gotFInQ, "CalculateTotal", CandidateFunction), "got %v", gotFInQ)

	// Built-in keyword functions remain unquoted.
	require.Equal(t, "CAST()", formatSessionCandidateText(Candidate{Text: "CAST", Type: CandidateFunction}, false))
}

func TestCompleteSession_ReadOnlyScene(t *testing.T) {
	schemas := oracleTestSchemas()
	cfg := SessionConfig{
		ReadOnly:      true,
		DefaultSchema: "SCHEMA1",
		ListSchemaNames: func(context.Context) ([]string, error) {
			return []string{"SCHEMA1", "SCHEMA2"}, nil
		},
		LoadSchema: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return schemas[strings.ToUpper(name)], nil
		},
	}

	// Statement start offers only SELECT, WITH, EXPLAIN keywords.
	stmt, line, col := catchCaretUTF16(t, "|")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	for _, kw := range []string{"SELECT", "WITH", "EXPLAIN"} {
		require.True(t, hasSessionCandidate(got, kw, CandidateKeyword), "missing %s in %v", kw, got)
	}
	for _, kw := range []string{"INSERT", "UPDATE", "DELETE", "MERGE", "CREATE", "ALTER", "DROP", "TRUNCATE"} {
		require.False(t, hasSessionCandidate(got, kw, CandidateKeyword), "unexpected %s in %v", kw, got)
	}

	// Prefix at statement start keeps matching query keyword only.
	stmtSel, lineSel, colSel := catchCaretUTF16(t, "SEL|")
	gotSel := CompleteSession(context.Background(), cfg, stmtSel, lineSel, colSel)
	require.True(t, hasSessionCandidate(gotSel, "SELECT", CandidateKeyword))
	require.False(t, hasSessionCandidate(gotSel, "INSERT", CandidateKeyword))

	// Write statements are disallowed completely.
	stmtDrop, lineDrop, colDrop := catchCaretUTF16(t, "DROP TABLE |")
	require.Nil(t, CompleteSession(context.Background(), cfg, stmtDrop, lineDrop, colDrop))

	// WITH-prefixed write statement is disallowed.
	stmtWithIns, lineWI, colWI := catchCaretUTF16(t, "WITH x AS (SELECT c1 FROM t1) INSERT INTO t1 SELECT * FROM |")
	require.Nil(t, CompleteSession(context.Background(), cfg, stmtWithIns, lineWI, colWI))

	// WITH-prefixed SELECT is allowed.
	stmtWithSel, lineWS, colWS := catchCaretUTF16(t, "WITH x AS (SELECT c1 FROM t1) SELECT * FROM |")
	gotWS := CompleteSession(context.Background(), cfg, stmtWithSel, lineWS, colWS)
	require.True(t, hasSessionCandidate(gotWS, "T1", CandidateTable), "got %v", gotWS)
	require.True(t, hasSessionCandidate(gotWS, "X", CandidateTable), "got %v", gotWS)

	// SELECT following a write statement is isolated and allowed.
	stmtMulti, lineM, colM := catchCaretUTF16(t, "DROP TABLE t1; SELECT * FROM |")
	gotM := CompleteSession(context.Background(), cfg, stmtMulti, lineM, colM)
	require.True(t, hasSessionCandidate(gotM, "T1", CandidateTable), "got %v", gotM)

	// Statement start after semicolon filters write keywords.
	stmtAfterSemi, lineAS, colAS := catchCaretUTF16(t, "SELECT 1; |")
	gotAS := CompleteSession(context.Background(), cfg, stmtAfterSemi, lineAS, colAS)
	require.True(t, hasSessionCandidate(gotAS, "SELECT", CandidateKeyword), "got %v", gotAS)
	require.False(t, hasSessionCandidate(gotAS, "DROP", CandidateKeyword), "got %v", gotAS)
}

func TestCompleteSession_LazySchemaLoadingAndCaching(t *testing.T) {
	schemas := oracleTestSchemas()
	loads := map[string]int{}
	cfg := SessionConfig{
		DefaultSchema: "SCHEMA1",
		ListSchemaNames: func(context.Context) ([]string, error) {
			return []string{"SCHEMA1", "SCHEMA2"}, nil
		},
		LoadSchema: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			key := strings.ToUpper(name)
			loads[key]++
			return schemas[key], nil
		},
	}

	stmt, line, col := catchCaretUTF16(t, "SELECT * FROM SCHEMA2.|")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "REMOTE_T", CandidateTable), "got %v", got)
	require.Equal(t, 1, loads["SCHEMA1"])
	require.Equal(t, 1, loads["SCHEMA2"])
}

func TestCompleteSession_MultiStatementAndDeterministicOrder(t *testing.T) {
	schemas := oracleTestSchemas()
	cfg := SessionConfig{
		DefaultSchema: "SCHEMA1",
		LoadSchema: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return schemas[strings.ToUpper(name)], nil
		},
	}

	stmt, line, col := catchCaretUTF16(t, "SELECT FROM broken; SELECT '😀', | FROM T2")
	got := CompleteSession(context.Background(), cfg, stmt, line, col)
	require.True(t, hasSessionCandidate(got, "C1", CandidateColumn), "got %v", got)
	require.True(t, hasSessionCandidate(got, "C2", CandidateColumn), "got %v", got)
	require.True(t, slices.IsSortedFunc(got, compareSessionCandidates))
}
