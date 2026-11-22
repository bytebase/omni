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
				CompositeTypes: []*metadatapb.CompositeTypeMetadata{
					{Name: "AddressType"},
					{Name: "ORDER_TYPE"},
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
	require.Equal(t, "CAST()", formatSessionCandidateText(Candidate{Text: "CAST", Type: CandidateFunction}, false, true))

	// DDL/DCL object-reference contexts omit "()" on function names.
	stmtDropFn, lineDF, colDF := catchCaretUTF16(t, "DROP FUNCTION |")
	gotDropFn := CompleteSession(context.Background(), cfg, stmtDropFn, lineDF, colDF)
	require.True(t, hasSessionCandidate(gotDropFn, "GET_TOTAL", CandidateFunction), "got %v", gotDropFn)
	require.False(t, hasSessionCandidate(gotDropFn, "GET_TOTAL()", CandidateFunction), "got %v", gotDropFn)
	require.True(t, hasSessionCandidate(gotDropFn, `"CalculateTotal"`, CandidateFunction), "got %v", gotDropFn)

	stmtGrantFn, lineGF, colGF := catchCaretUTF16(t, "GRANT EXECUTE ON |")
	gotGrantFn := CompleteSession(context.Background(), cfg, stmtGrantFn, lineGF, colGF)
	require.True(t, hasSessionCandidate(gotGrantFn, "GET_TOTAL", CandidateFunction), "got %v", gotGrantFn)
	require.False(t, hasSessionCandidate(gotGrantFn, "GET_TOTAL()", CandidateFunction), "got %v", gotGrantFn)

	stmtGrantQual, lineGQ, colGQ := catchCaretUTF16(t, "GRANT EXECUTE ON SCHEMA1.|")
	gotGrantQual := CompleteSession(context.Background(), cfg, stmtGrantQual, lineGQ, colGQ)
	require.True(t, hasSessionCandidate(gotGrantQual, "GET_TOTAL", CandidateFunction), "got %v", gotGrantQual)
	require.False(t, hasSessionCandidate(gotGrantQual, "GET_TOTAL()", CandidateFunction), "got %v", gotGrantQual)

	// Metadata-backed composite types are quoted when needed, while built-in data types stay unquoted.
	stmtType, lineT, colT := catchCaretUTF16(t, "DROP TYPE |")
	gotType := CompleteSession(context.Background(), cfg, stmtType, lineT, colT)
	require.True(t, hasSessionCandidate(gotType, `"AddressType"`, CandidateType_), "got %v", gotType)
	require.True(t, hasSessionCandidate(gotType, "ORDER_TYPE", CandidateType_), "got %v", gotType)

	stmtTypeQ, lineTQ, colTQ := catchCaretUTF16(t, `DROP TYPE "Addr|"`)
	gotTypeQ := CompleteSession(context.Background(), cfg, stmtTypeQ, lineTQ, colTQ)
	require.True(t, hasSessionCandidate(gotTypeQ, "AddressType", CandidateType_), "got %v", gotTypeQ)

	stmtColType, lineCT, colCT := catchCaretUTF16(t, "CREATE TABLE t (c |)")
	gotColType := CompleteSession(context.Background(), cfg, stmtColType, lineCT, colCT)
	require.True(t, hasSessionCandidate(gotColType, "NUMBER", CandidateType_), "got %v", gotColType)
	require.True(t, hasSessionCandidate(gotColType, `"AddressType"`, CandidateType_), "got %v", gotColType)
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

	// Plain SELECT table reference is allowed in read-only mode.
	stmtPlainSel, linePS, colPS := catchCaretUTF16(t, "SELECT * FROM |")
	gotPS := CompleteSession(context.Background(), cfg, stmtPlainSel, linePS, colPS)
	require.True(t, hasSessionCandidate(gotPS, "T1", CandidateTable), "got %v", gotPS)
	require.True(t, hasSessionCandidate(gotPS, "T2", CandidateTable), "got %v", gotPS)

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

func oracleSessionCorpusConfig() SessionConfig {
	return SessionConfig{
		DefaultSchema: "SCHEMA1",
		ListSchemaNames: func(context.Context) ([]string, error) {
			return []string{"SCHEMA1", "SCHEMA2", "SCHEMA3"}, nil
		},
		LoadSchema: func(_ context.Context, name string) (*metadatapb.DatabaseSchemaMetadata, error) {
			return loadTestSchema(strings.ToUpper(name)), nil
		},
	}
}

func completeSessionCorpus(t *testing.T, input string) []Candidate {
	t.Helper()
	stmt, line, col := catchCaretUTF16(t, input)
	return CompleteSession(context.Background(), oracleSessionCorpusConfig(), stmt, line, col)
}

func TestCompleteSession_GoldenCorpus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Candidate
	}{
		{
			name: "multiple statements limits to current statement",
			input: `select count(1) from t1 where id; SELECT * FROM |
select * from QRTZ_J0B_DETAILS qjd where J0B_NAME like '%gray';`,
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "derived table star expansion",
			input: "SELECT x.| FROM (SELECT * FROM t2) x",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateColumn, Text: "C2"},
			},
		},
		{
			name:  "cte star expansion",
			input: "WITH x AS (SELECT * FROM t2) SELECT x.| FROM x;",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateColumn, Text: "C2"},
			},
		},
		{
			name:  "cte explicit column list",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT x.| FROM x;",
			want: []Candidate{
				{Type: CandidateColumn, Text: "X1"},
				{Type: CandidateColumn, Text: "X2"},
			},
		},
		{
			name:  "cte explicit columns unqualified select",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT | FROM x",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateTable, Text: "X"},
				{Type: CandidateColumn, Text: "X1"},
				{Type: CandidateColumn, Text: "X2"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "select aliases in group by",
			input: "SELECT c1 as eid, c2 as xid FROM t2 GROUP BY |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "EID"},
				{Type: CandidateColumn, Text: "XID"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "select aliases in order by",
			input: "SELECT c1 as eid, c2 as xid FROM t2 ORDER BY |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "EID"},
				{Type: CandidateColumn, Text: "XID"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "single select alias in order by",
			input: "SELECT c1 as eid FROM t1 ORDER BY |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "EID"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name: "select columns after invalid first statement",
			input: `SELECT FROM basdkfjasldf;
SELECT | FROM t1`,
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "derived table alias and columns",
			input: "SELECT | FROM (SELECT c1 FROM t1) cc1",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "CC1"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "join left alias qualified",
			input: "SELECT cc1.| FROM t1 cc1 JOIN t2 on NOT cc1.c1 = t2.c1",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
			},
		},
		{
			name:  "join right table qualified in on",
			input: "SELECT * from t1 cc1 JOIN t2 on cc1.c1 = t2.|",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
			},
		},
		{
			name:  "join left alias qualified in on",
			input: "SELECT * from t1 cc1 JOIN t2 on cc1.|",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
			},
		},
		{
			name:  "function argument qualified alias",
			input: "SELECT MAX(cc1.|) FROM t2 cc1",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
			},
		},
		{
			name:  "select list qualified alias",
			input: "SELECT cc1.| FROM t2 cc1",
			want: []Candidate{
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T2 | int, NOT NULL"},
				{Type: CandidateColumn, Text: "C2", Definition: "SCHEMA1.T2 | int, NOT NULL"},
			},
		},
		{
			name:  "select from root objects",
			input: "SELECT * FROM |",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "select list physical alias offered as table",
			input: "SELECT | FROM t1 x",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateTable, Text: "X"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
		{
			name:  "schema qualified table in from",
			input: "SELECT | FROM schema1.t1",
			want: []Candidate{
				{Type: CandidateSchema, Text: "SCHEMA1"},
				{Type: CandidateSchema, Text: "SCHEMA2"},
				{Type: CandidateSchema, Text: "SCHEMA3"},
				{Type: CandidateTable, Text: "T1"},
				{Type: CandidateTable, Text: "T2"},
				{Type: CandidateColumn, Text: "C1", Definition: "SCHEMA1.T1 | int, NOT NULL"},
				{Type: CandidateView, Text: "V1"},
				{Type: CandidateSequence, Text: "SEQ1"},
				{Type: CandidateSequence, Text: "USER_ID_SEQ"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := completeSessionCorpus(t, tt.input)
			var filtered []Candidate
			for _, c := range got {
				if c.Type == CandidateKeyword || c.Type == CandidateFunction {
					continue
				}
				filtered = append(filtered, c)
			}
			slices.SortFunc(tt.want, compareSessionCandidates)
			require.Equal(t, tt.want, filtered)
		})
	}
}

func TestCompleteSession_LongTailTableSources(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantColumns []string
	}{
		{
			name:        "json table explicit columns",
			input:       "SELECT jt.| FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
			wantColumns: []string{"ID", "NAME"},
		},
		{
			name:        "xmltable explicit columns",
			input:       "SELECT xt.| FROM t1, XMLTABLE('/root/row' PASSING t1.c1 COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') xt",
			wantColumns: []string{"ID", "NAME"},
		},
		{
			name:        "inline external declared columns",
			input:       "SELECT ext.| FROM EXTERNAL ((a NUMBER, b VARCHAR2(10)) TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir LOCATION ('x.csv')) ext",
			wantColumns: []string{"A", "B"},
		},
		{
			name:        "lateral inline view columns",
			input:       "SELECT l.| FROM t1, LATERAL (SELECT c1 FROM t2 WHERE t2.c1 = t1.c1) l",
			wantColumns: []string{"C1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := completeSessionCorpus(t, tt.input)
			var columns []string
			for _, c := range got {
				if c.Type == CandidateColumn {
					columns = append(columns, c.Text)
				}
			}
			slices.Sort(columns)
			require.Equal(t, tt.wantColumns, columns)
		})
	}
}

func TestCompleteSession_CoverageMatrix(t *testing.T) {
	schs := func(texts ...string) []candidateSpec { return specs(CandidateSchema, texts...) }

	tests := []struct {
		name    string
		input   string
		want    []candidateSpec
		notWant []candidateSpec
	}{
		{
			name:  "select list uses table columns",
			input: "SELECT | FROM t1",
			want:  cols("C1"),
		},
		{
			name:  "select list uses alias columns",
			input: "SELECT a.| FROM t1 a",
			want:  cols("C1"),
		},
		{
			name:  "select list uses schema qualified table columns",
			input: "SELECT | FROM schema1.t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "column prefix keeps matching column available",
			input: "SELECT C2| FROM t2",
			want:  cols("C2"),
		},
		{
			name:  "qualified column prefix keeps matching column available",
			input: "SELECT t2.C2| FROM t2",
			want:  cols("C2"),
		},
		{
			name:    "unknown qualifier offers no columns",
			input:   "SELECT missing.| FROM t1",
			notWant: cols("C1", "C2"),
		},
		{
			name:  "join on left alias columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON a.|",
			want:  cols("C1"),
		},
		{
			name:  "join on right alias columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON b.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "join where sees both table columns",
			input: "SELECT * FROM t1 a JOIN t2 b ON a.c1 = b.c1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:    "left alias excludes right only column",
			input:   "SELECT a.| FROM t1 a JOIN t2 b ON a.c1 = b.c1",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "join table reference offers tables",
			input: "SELECT * FROM t1 JOIN |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "using clause offers joined columns",
			input: "SELECT * FROM t1 JOIN t2 USING (|)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte explicit columns are available",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT x.| FROM x",
			want:  cols("X1", "X2"),
		},
		{
			name:  "cte star body expands through query span",
			input: "WITH x AS (SELECT * FROM t2) SELECT x.| FROM x",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte table candidate is available without schema qualifier",
			input: "WITH x AS (SELECT * FROM t2) SELECT * FROM |",
			want:  tbls("X", "T1", "T2"),
		},
		{
			name:    "cte table candidate is excluded from schema qualified table reference",
			input:   "WITH x AS (SELECT * FROM t2) SELECT * FROM schema1.|",
			want:    tbls("T1", "T2"),
			notWant: tbls("X"),
		},
		{
			name:  "multiple ctes expose selected cte columns",
			input: "WITH a AS (SELECT c1 FROM t1), b AS (SELECT c1, c2 FROM t2) SELECT b.| FROM a JOIN b ON a.c1 = b.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "derived table explicit projection columns",
			input: "SELECT x.| FROM (SELECT c1 FROM t1) x",
			want:  cols("C1"),
		},
		{
			name:  "derived table star body expands through query span",
			input: "SELECT x.| FROM (SELECT * FROM t2) x",
			want:  cols("C1", "C2"),
		},
		{
			name:    "unqualified select from derived table exposes derived columns",
			input:   "SELECT | FROM (SELECT c1 FROM t2) x",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "lateral alias columns",
			input: "SELECT l.| FROM t1, LATERAL (SELECT c1 FROM t2 WHERE t2.c1 = t1.c1) l",
			want:  cols("C1"),
		},
		{
			name:    "lateral body sees left outer reference",
			input:   "SELECT * FROM t1, LATERAL (SELECT t1.| FROM t2) l",
			want:    cols("C1"),
			notWant: cols("C2"),
		},
		{
			name:  "json table qualified columns",
			input: "SELECT jt.| FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
			want:  cols("ID", "NAME"),
		},
		{
			name:  "json table unqualified columns",
			input: "SELECT | FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
			want:  cols("C1", "ID", "NAME"),
		},
		{
			name:  "xmltable qualified columns",
			input: "SELECT xt.| FROM t1, XMLTABLE('/root/row' PASSING t1.c1 COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') xt",
			want:  cols("ID", "NAME"),
		},
		{
			name:  "inline external qualified columns",
			input: "SELECT ext.| FROM EXTERNAL ((a NUMBER, b VARCHAR2(10)) TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir LOCATION ('x.csv')) ext",
			want:  cols("A", "B"),
		},
		{
			name:  "containers alias resolves base object columns",
			input: "SELECT c.| FROM CONTAINERS(schema1.t2) c",
			want:  cols("C1", "C2"),
		},
		{
			name:    "table function alias does not invent columns",
			input:   "SELECT x.| FROM TABLE(pkg.values()) x",
			notWant: cols("C1", "C2"),
		},
		{
			name:  "insert target table reference",
			input: "INSERT INTO |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "insert column list uses target table columns",
			input: "INSERT INTO t2 (|) VALUES (1)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert select uses source table columns",
			input: "INSERT INTO t1 SELECT | FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update target table reference",
			input: "UPDATE | SET c1 = 1",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "update set uses target table columns",
			input: "UPDATE t2 SET |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update where uses target table columns",
			input: "UPDATE t2 SET c1 = 1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "delete target table reference",
			input: "DELETE FROM |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "delete where uses target table columns",
			input: "DELETE FROM t2 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge target table reference",
			input: "MERGE INTO |",
			want:  tbls("T1", "T2"),
		},
		{
			name:    "merge source table reference offers real tables",
			input:   "MERGE INTO t1 target USING |",
			want:    tbls("T1", "T2"),
			notWant: tbls("TARGET"),
		},
		{
			name:  "merge on uses target and source columns",
			input: "MERGE INTO t1 target USING t2 source ON |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "schema qualified table reference offers schema tables",
			input: "SELECT * FROM schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "table prefix keeps matching table available",
			input: "SELECT * FROM t|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted identifier context preserves candidate text",
			input: `SELECT "T2".| FROM t2`,
			want:  cols("C1", "C2"),
		},
		{
			name:  "inside double-quoted identifier does not double-quote candidates",
			input: `SELECT * FROM "T|`,
			want:  tbls("T1", "T2"),
		},
		{
			name:  "select list uses t2 columns",
			input: "SELECT | FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "where clause uses table columns",
			input: "SELECT * FROM t2 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "order by clause uses table columns",
			input: "SELECT * FROM t2 ORDER BY |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "group by clause uses table columns",
			input: "SELECT * FROM t2 GROUP BY |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "having clause uses grouped table columns",
			input: "SELECT c1, COUNT(*) FROM t2 GROUP BY c1 HAVING |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "distinct select list uses table columns",
			input: "SELECT DISTINCT | FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "table name qualifier uses table columns",
			input: "SELECT t2.| FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "comparison expression uses table columns",
			input: "SELECT * FROM t2 WHERE c1 = |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "arithmetic expression uses table columns",
			input: "SELECT * FROM t2 WHERE c1 + |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "case expression uses table columns",
			input: "SELECT CASE WHEN | THEN c1 ELSE c2 END FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "function argument uses table columns",
			input: "SELECT NVL(|, 0) FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "between lower bound uses table columns",
			input: "SELECT * FROM t2 WHERE c1 BETWEEN | AND 10",
			want:  cols("C1", "C2"),
		},
		{
			name:  "between upper bound uses table columns",
			input: "SELECT * FROM t2 WHERE c1 BETWEEN 1 AND |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "like expression uses table columns",
			input: "SELECT * FROM t2 WHERE c1 LIKE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "is null predicate continuation uses table columns",
			input: "SELECT * FROM t2 WHERE c1 IS NULL OR |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "select alias is available in order by",
			input: "SELECT c1 AS c_alias FROM t1 ORDER BY |",
			want:  cols("C1", "C_ALIAS"),
		},
		{
			name:  "select alias survives nested select in order by",
			input: "SELECT c1 AS c_alias FROM t1 WHERE EXISTS (SELECT c1 FROM t2) ORDER BY |",
			want:  cols("C1", "C_ALIAS"),
		},
		{
			name:    "select alias is not available in where",
			input:   "SELECT c1 AS c_alias FROM t1 WHERE |",
			want:    cols("C1"),
			notWant: cols("C_ALIAS"),
		},
		{
			name:  "select alias is available in group by",
			input: "SELECT c1 AS c_alias FROM t1 GROUP BY |",
			want:  cols("C1", "C_ALIAS"),
		},
		{
			name:  "comma join select list sees both tables",
			input: "SELECT | FROM t1 a, t2 b",
			want:  cols("C1", "C2"),
		},
		{
			name:  "comma join left alias columns",
			input: "SELECT a.| FROM t1 a, t2 b",
			want:  cols("C1"),
		},
		{
			name:  "comma join right alias columns",
			input: "SELECT b.| FROM t1 a, t2 b",
			want:  cols("C1", "C2"),
		},
		{
			name:  "left join select list sees both sides",
			input: "SELECT | FROM t1 LEFT JOIN t2 ON t1.c1 = t2.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "right join select list sees both sides",
			input: "SELECT | FROM t1 RIGHT JOIN t2 ON t1.c1 = t2.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "full outer join select list sees both sides",
			input: "SELECT | FROM t1 FULL OUTER JOIN t2 ON t1.c1 = t2.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cross join select list sees both sides",
			input: "SELECT | FROM t1 CROSS JOIN t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "natural join select list sees both sides",
			input: "SELECT | FROM t1 NATURAL JOIN t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "multi join final alias columns",
			input: "SELECT c.| FROM t1 a JOIN t2 b ON a.c1 = b.c1 JOIN t2 c ON b.c1 = c.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "multi join where sees all aliases",
			input: "SELECT * FROM t1 a JOIN t2 b ON a.c1 = b.c1 JOIN t2 c ON b.c1 = c.c1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "schema qualified join alias columns",
			input: "SELECT s.| FROM t1 a JOIN schema2.t2 s ON a.c1 = s.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:    "schema qualified alias excludes same table from other schema",
			input:   "SELECT a.| FROM schema1.t1 a JOIN schema4.t1 b ON a.c1 = b.c1",
			want:    cols("C1"),
			notWant: cols("SCHEMA4_ONLY"),
		},
		{
			name:  "join table prefix keeps matching tables",
			input: "SELECT * FROM t1 JOIN t|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "left join table reference offers tables",
			input: "SELECT * FROM t1 LEFT JOIN |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "right join table reference offers tables",
			input: "SELECT * FROM t1 RIGHT JOIN |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "cross join table reference offers tables",
			input: "SELECT * FROM t1 CROSS JOIN |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "subquery select list uses inner table columns",
			input: "SELECT * FROM t1 WHERE EXISTS (SELECT | FROM t2)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "subquery where uses inner table columns",
			input: "SELECT * FROM t1 WHERE EXISTS (SELECT * FROM t2 WHERE |)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "correlated subquery can qualify outer table",
			input: "SELECT * FROM t1 outer_t WHERE EXISTS (SELECT * FROM t2 WHERE outer_t.|)",
			want:  cols("C1"),
		},
		{
			name:  "correlated subquery can qualify inner table",
			input: "SELECT * FROM t1 outer_t WHERE EXISTS (SELECT * FROM t2 inner_t WHERE inner_t.|)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "scalar subquery select list uses inner table columns",
			input: "SELECT (SELECT | FROM t2) FROM t1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "in subquery select list uses inner table columns",
			input: "SELECT * FROM t1 WHERE c1 IN (SELECT | FROM t2)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "nested derived table alias columns",
			input: "SELECT y.| FROM (SELECT x.c1 FROM (SELECT c1 FROM t1) x) y",
			want:  cols("C1"),
		},
		{
			name:  "derived table with column alias",
			input: "SELECT x.| FROM (SELECT c1 AS renamed FROM t1) x",
			want:  cols("RENAMED"),
		},
		{
			name:  "derived table join exposes derived alias columns",
			input: "SELECT x.| FROM (SELECT c1 FROM t1) x JOIN t2 ON x.c1 = t2.c1",
			want:  cols("C1"),
		},
		{
			name:  "derived table join where sees derived and physical columns",
			input: "SELECT * FROM (SELECT c1 FROM t1) x JOIN t2 ON x.c1 = t2.c1 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte selected through derived table",
			input: "WITH x AS (SELECT c1 FROM t1) SELECT | FROM (SELECT * FROM x) sub",
			want:  cols("C1"),
		},
		{
			name:  "cte joined with physical table",
			input: "WITH x AS (SELECT c1 FROM t1) SELECT | FROM x JOIN t2 ON x.c1 = t2.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte alias joined with physical table",
			input: "WITH x AS (SELECT c1 FROM t1) SELECT xa.| FROM x xa JOIN t2 ON xa.c1 = t2.c1",
			want:  cols("C1"),
		},
		{
			name:  "cte in insert select",
			input: "WITH x AS (SELECT * FROM t2) INSERT INTO t1 SELECT | FROM x",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte in update scalar subquery",
			input: "WITH x AS (SELECT c1, c2 FROM t2) UPDATE t1 SET c1 = (SELECT | FROM x)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte in delete subquery",
			input: "WITH x AS (SELECT * FROM t2) DELETE FROM t1 WHERE c1 IN (SELECT | FROM x)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte set operation exposes positional columns",
			input: "WITH x AS (SELECT c1 FROM t1 UNION SELECT c1 FROM t2) SELECT x.| FROM x",
			want:  cols("C1"),
		},
		{
			name:  "lateral body uses inner alias columns",
			input: "SELECT * FROM t1, LATERAL (SELECT inner_t.| FROM t2 inner_t WHERE inner_t.c1 = t1.c1) l",
			want:  cols("C1", "C2"),
		},
		{
			name:  "left join lateral alias columns",
			input: "SELECT l.| FROM t1 LEFT JOIN LATERAL (SELECT c1 FROM t2 WHERE t2.c1 = t1.c1) l ON 1 = 1",
			want:  cols("C1"),
		},
		{
			name:  "lateral unqualified select list sees inner columns",
			input: "SELECT * FROM t1, LATERAL (SELECT | FROM t2 WHERE t2.c1 = t1.c1) l",
			want:  cols("C1", "C2"),
		},
		{
			name:  "json table table alias is offered as table candidate",
			input: "SELECT | FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id')) jt",
			want:  tbls("JT"),
		},
		{
			name:  "xmltable unqualified columns",
			input: "SELECT | FROM t1, XMLTABLE('/root/row' PASSING t1.c1 COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') xt",
			want:  cols("C1", "ID", "NAME"),
		},
		{
			name:  "inline external unqualified columns",
			input: "SELECT | FROM EXTERNAL ((a NUMBER, b VARCHAR2(10)) TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir LOCATION ('x.csv')) ext",
			want:  cols("A", "B"),
		},
		{
			name:  "containers without alias resolves base object columns",
			input: "SELECT t2.| FROM CONTAINERS(schema1.t2)",
			want:  cols("C1", "C2"),
		},
		{
			name:    "pivot does not invent transformed columns",
			input:   "SELECT | FROM t2 PIVOT (COUNT(*) FOR c1 IN (1 AS one))",
			notWant: cols("ONE"),
		},
		{
			name:    "unpivot does not invent transformed columns",
			input:   "SELECT | FROM t2 UNPIVOT (val FOR col IN (c1 AS 'C1', c2 AS 'C2'))",
			notWant: cols("VAL", "COL"),
		},
		{
			name:    "model does not invent transformed columns",
			input:   "SELECT | FROM t2 MODEL DIMENSION BY (c1) MEASURES (c2) RULES (c2[1] = 1)",
			notWant: cols("VAL", "COL"),
		},
		{
			name:  "insert schema qualified table reference",
			input: "INSERT INTO schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "insert schema qualified target columns",
			input: "INSERT INTO schema2.t2 (|) VALUES (1)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert values expression uses target columns",
			input: "INSERT INTO t2 VALUES (|)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert second column list position uses target columns",
			input: "INSERT INTO t2 (c1, |) VALUES (1, 2)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert select with schema qualified source",
			input: "INSERT INTO t1 SELECT | FROM schema2.t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update schema qualified target reference",
			input: "UPDATE schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "update schema qualified target set columns",
			input: "UPDATE schema2.t2 SET |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update set expression uses target columns",
			input: "UPDATE t2 SET c1 = |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update where continuation uses target columns",
			input: "UPDATE t2 SET c1 = 1 WHERE c1 > 0 AND |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "delete schema qualified table reference",
			input: "DELETE FROM schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "delete where continuation uses target columns",
			input: "DELETE FROM t2 WHERE c1 > 0 AND |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge target schema qualified table reference",
			input: "MERGE INTO schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "merge source schema qualified table reference",
			input: "MERGE INTO t1 target USING schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "merge qualified target columns",
			input: "MERGE INTO t1 target USING t2 source ON target.|",
			want:  cols("C1"),
		},
		{
			name:  "merge qualified source columns",
			input: "MERGE INTO t1 target USING t2 source ON source.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge update set uses target and source columns",
			input: "MERGE INTO t1 target USING t2 source ON target.c1 = source.c1 WHEN MATCHED THEN UPDATE SET target.c1 = |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge insert values uses target and source columns",
			input: "MERGE INTO t1 target USING t2 source ON target.c1 = source.c1 WHEN NOT MATCHED THEN INSERT (c1) VALUES (|)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "select from root offers schemas",
			input: "SELECT * FROM |",
			want:  schs("SCHEMA1", "SCHEMA2", "SCHEMA3"),
		},
		{
			name:  "select from root offers views",
			input: "SELECT * FROM |",
			want:  vws("V1"),
		},
		{
			name:  "select from root offers sequences",
			input: "SELECT * FROM |",
			want:  seqs("SEQ1", "USER_ID_SEQ"),
		},
		{
			name:  "schema2 table reference offers schema2 tables",
			input: "SELECT * FROM schema2.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "schema3 table reference offers schema3 tables",
			input: "SELECT * FROM schema3.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "schema prefix keeps matching schema available",
			input: "SELECT * FROM schema|",
			want:  schs("SCHEMA1", "SCHEMA2", "SCHEMA3"),
		},
		{
			name:  "view reference offers view candidate",
			input: "SELECT * FROM v|",
			want:  vws("V1"),
		},
		{
			name:  "sequence reference offers sequence candidates",
			input: "DROP SEQUENCE |",
			want:  seqs("SEQ1", "USER_ID_SEQ"),
		},
		{
			name:  "sequence prefix keeps matching sequence",
			input: "DROP SEQUENCE user|",
			want:  seqs("USER_ID_SEQ"),
		},
		{
			name:  "alter table reference offers tables",
			input: "ALTER TABLE |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "alter schema qualified table reference offers tables",
			input: "ALTER TABLE schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "alter drop column offers table columns",
			input: "ALTER TABLE t2 DROP COLUMN |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "drop table reference offers tables",
			input: "DROP TABLE |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "drop schema qualified table reference offers tables",
			input: "DROP TABLE schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "drop view reference offers views",
			input: "DROP VIEW |",
			want:  vws("V1"),
		},
		{
			name:  "truncate table reference offers tables",
			input: "TRUNCATE TABLE |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "comment on table reference offers tables",
			input: "COMMENT ON TABLE |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "comment on column offers table columns",
			input: "COMMENT ON COLUMN t2.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "grant on object offers tables",
			input: "GRANT SELECT ON |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "revoke on object offers tables",
			input: "REVOKE SELECT ON |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "create table datatype offers NUMBER",
			input: "CREATE TABLE t (c |)",
			want:  typs("NUMBER"),
		},
		{
			name:  "create table references offers tables",
			input: "CREATE TABLE child (parent_id NUMBER REFERENCES |)",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "multiple statements limits to statement at caret insert",
			input: "SELECT * FROM t1; INSERT INTO |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "multiple statements limits to statement at caret update",
			input: "SELECT * FROM t1; UPDATE | SET c1 = 1",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "multiple statements limits to statement at caret delete",
			input: "SELECT * FROM t1; DELETE FROM |",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "multiple statements select columns after invalid first statement",
			input: "SELECT FROM broken; SELECT | FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "table prefix after multiple statements",
			input: "SELECT * FROM t1; SELECT * FROM t|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted table alias columns",
			input: `SELECT "Alias".| FROM t2 "Alias"`,
			want:  cols("C1", "C2"),
		},
		{
			name:  "quoted cte alias columns",
			input: `WITH "X"(x1, x2) AS (SELECT * FROM t2) SELECT "X".| FROM "X"`,
			want:  cols("X1", "X2"),
		},
		{
			name:  "quoted derived alias columns",
			input: `SELECT "Sub".| FROM (SELECT c1 FROM t1) "Sub"`,
			want:  cols("C1"),
		},
		{
			name:  "quoted schema qualified table reference",
			input: `SELECT * FROM "SCHEMA1".|`,
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted schema qualified columns",
			input: `SELECT | FROM "SCHEMA1"."T2"`,
			want:  cols("C1", "C2"),
		},
		{
			name:  "reserved-like quoted alias columns",
			input: `SELECT "SELECT".| FROM t2 "SELECT"`,
			want:  cols("C1", "C2"),
		},
		{
			name:    "database link table keeps local table columns unavailable when metadata missing",
			input:   "SELECT remote_t.| FROM t2@remote remote_t",
			notWant: cols("C1", "C2"),
		},
		{
			name:    "schema qualified unknown table offers no known columns",
			input:   "SELECT | FROM schema1.missing_table",
			notWant: cols("C1", "C2"),
		},
		{
			name:    "unknown schema table reference offers no local tables",
			input:   "SELECT * FROM missing_schema.|",
			notWant: tbls("T1", "T2"),
		},
		{
			name:    "unknown cte qualifier offers no columns",
			input:   "WITH x AS (SELECT c1 FROM t1) SELECT y.| FROM x",
			notWant: cols("C1", "C2"),
		},
		{
			name:    "unknown physical alias offers no columns",
			input:   "SELECT z.| FROM t1 a JOIN t2 b ON a.c1 = b.c1",
			notWant: cols("C1", "C2"),
		},
		{
			name:    "unknown lateral alias offers no columns",
			input:   "SELECT missing.| FROM t1, LATERAL (SELECT c1 FROM t2) l",
			notWant: cols("C1", "C2"),
		},
		{
			name:  "hierarchical start with uses from scope",
			input: "SELECT * FROM t2 START WITH | CONNECT BY PRIOR c1 = c2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "hierarchical connect by uses from scope",
			input: "SELECT * FROM t2 START WITH c1 IS NULL CONNECT BY |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "analytic partition by uses table columns",
			input: "SELECT COUNT(*) OVER (PARTITION BY |) FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "analytic order by uses table columns",
			input: "SELECT COUNT(*) OVER (ORDER BY |) FROM t2",
			want:  cols("C1", "C2"),
		},
		{
			name:  "for update select list keeps table columns",
			input: "SELECT | FROM t2 FOR UPDATE",
			want:  cols("C1", "C2"),
		},
		{
			name:  "subquery alias table candidate is offered",
			input: "SELECT | FROM (SELECT c1 FROM t1) sub1",
			want:  tbls("SUB1"),
		},
		{
			name:  "physical alias table candidate is offered",
			input: "SELECT | FROM t1 x",
			want:  tbls("X"),
		},
		{
			name:  "schema2 table prefix keeps matching tables",
			input: "SELECT * FROM schema2.t|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "quoted schema2 table prefix keeps matching tables",
			input: `SELECT * FROM "SCHEMA2".t|`,
			want:  tbls("T1", "T2"),
		},
		{
			name:  "schema3 t1 select list uses columns",
			input: "SELECT | FROM schema3.t1",
			want:  cols("C1"),
		},
		{
			name:  "schema3 alias select list uses columns",
			input: "SELECT s.| FROM schema3.t2 s",
			want:  cols("C1", "C2"),
		},
		{
			name:  "schema3 alias where uses columns",
			input: "SELECT * FROM schema3.t2 s WHERE s.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "schema3 alias order by uses columns",
			input: "SELECT * FROM schema3.t2 s ORDER BY s.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "join schema3 table alias columns",
			input: "SELECT s.| FROM t1 a JOIN schema3.t2 s ON a.c1 = s.c1",
			want:  cols("C1", "C2"),
		},
		{
			name:  "insert schema3 qualified target columns",
			input: "INSERT INTO schema3.t2 (|) VALUES (1)",
			want:  cols("C1", "C2"),
		},
		{
			name:  "update schema3 qualified target set columns",
			input: "UPDATE schema3.t2 SET |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "delete schema3 qualified target where columns",
			input: "DELETE FROM schema3.t2 WHERE |",
			want:  cols("C1", "C2"),
		},
		{
			name:  "merge schema3 target qualified columns",
			input: "MERGE INTO schema3.t1 target USING schema2.t2 source ON source.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "cte table prefix keeps matching cte",
			input: "WITH x AS (SELECT * FROM t2) SELECT * FROM x|",
			want:  tbls("X"),
		},
		{
			name:  "cte explicit column prefix keeps matching column",
			input: "WITH x(x1, x2) AS (SELECT * FROM t2) SELECT x.X1| FROM x",
			want:  cols("X1"),
		},
		{
			name:  "derived table column prefix keeps matching column",
			input: "SELECT x.C1| FROM (SELECT c1 FROM t1) x",
			want:  cols("C1"),
		},
		{
			name:  "json table column prefix keeps matching column",
			input: "SELECT jt.I| FROM t1, JSON_TABLE(t1.c1, '$.rows[*]' COLUMNS (id NUMBER PATH '$.id', name VARCHAR2(100) PATH '$.name')) jt",
			want:  cols("ID"),
		},
		{
			name:  "xmltable column prefix keeps matching column",
			input: "SELECT xt.N| FROM t1, XMLTABLE('/root/row' PASSING t1.c1 COLUMNS id NUMBER PATH 'id', name VARCHAR2(100) PATH 'name') xt",
			want:  cols("NAME"),
		},
		{
			name:  "inline external column prefix keeps matching column",
			input: "SELECT ext.A| FROM EXTERNAL ((a NUMBER, b VARCHAR2(10)) TYPE ORACLE_LOADER DEFAULT DIRECTORY data_dir LOCATION ('x.csv')) ext",
			want:  cols("A"),
		},
		{
			name:  "comment on schema qualified column offers columns",
			input: "COMMENT ON COLUMN schema1.t2.|",
			want:  cols("C1", "C2"),
		},
		{
			name:  "grant on schema qualified object offers tables",
			input: "GRANT SELECT ON schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "revoke on schema qualified object offers tables",
			input: "REVOKE SELECT ON schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "truncate schema qualified table reference offers tables",
			input: "TRUNCATE TABLE schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:  "comment on schema qualified table reference offers tables",
			input: "COMMENT ON TABLE schema1.|",
			want:  tbls("T1", "T2"),
		},
		{
			name:    "drop schema qualified view reference offers only views",
			input:   "DROP VIEW schema1.|",
			want:    vws("V1"),
			notWant: append(tbls("T1", "T2"), seqs("SEQ1", "USER_ID_SEQ")...),
		},
		{
			name:  "view schema qualified table reference offers views",
			input: "SELECT * FROM schema1.v|",
			want:  vws("V1"),
		},
		{
			name:  "sequence schema qualified reference offers sequences",
			input: "DROP SEQUENCE schema1.|",
			want:  seqs("SEQ1", "USER_ID_SEQ"),
		},
		{
			name:    "sequence schema qualified reference uses non default schema",
			input:   "DROP SEQUENCE schema2.|",
			want:    seqs("SCHEMA2_SEQ"),
			notWant: seqs("SEQ1", "USER_ID_SEQ"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := completeSessionCorpus(t, tt.input)
			for _, want := range tt.want {
				require.Truef(t, hasSpec(got, want), "missing %s candidate %q in %v", want.typ, want.text, got)
			}
			for _, notWant := range tt.notWant {
				require.Falsef(t, hasSpec(got, notWant), "unexpected %s candidate %q in %v", notWant.typ, notWant.text, got)
			}
		})
	}
}
