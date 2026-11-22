package mssql

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func assertSplitInvariants(t *testing.T, sql string, segs []Segment) {
	t.Helper()
	if len(segs) == 0 {
		return
	}
	require.Equal(t, 0, segs[0].ByteStart, "first segment must start at byte 0")
	require.Equal(t, Position{Line: 1, Column: 1}, segs[0].Start, "first segment must start at 1:1")
	for i, seg := range segs {
		require.GreaterOrEqual(t, seg.ByteEnd, seg.ByteStart, "segment %d negative range", i)
		require.LessOrEqual(t, seg.ByteEnd, len(sql), "segment %d out of bounds", i)
		require.Equal(t, sql[seg.ByteStart:seg.ByteEnd], seg.Text, "segment %d text mismatch", i)
		if i > 0 {
			require.Equal(t, segs[i-1].ByteEnd, seg.ByteStart, "segment %d byte contiguity broken", i)
			require.Equal(t, segs[i-1].End, seg.Start, "segment %d position adjacency broken", i)
		}
	}
}

func TestSplitValidScripts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Segment
	}{
		{
			name:  "empty and whitespace only",
			input: "   \n\t  ",
			want:  nil,
		},
		{
			name:  "GO batch separator with multiple batches",
			input: "-- first statement\ndeclare @temp table(a int)\n-- second statement\ninsert into @temp values(1)\n-- third statement\nselect * from @temp\n-- go statement\ngo\ngo\ngo\ngo\n\t\t\t",
			want: []Segment{
				{Text: "-- first statement\ndeclare @temp table(a int)", ByteStart: 0, ByteEnd: 45, Start: Position{1, 1}, End: Position{2, 27}, Empty: false},
				{Text: "\n-- second statement\ninsert into @temp values(1)", ByteStart: 45, ByteEnd: 93, Start: Position{2, 27}, End: Position{4, 28}, Empty: false},
				{Text: "\n-- third statement\nselect * from @temp", ByteStart: 93, ByteEnd: 132, Start: Position{4, 28}, End: Position{6, 20}, Empty: false},
				{Text: "\n-- go statement\ngo", ByteStart: 132, ByteEnd: 151, Start: Position{6, 20}, End: Position{8, 3}, Empty: true},
				{Text: "\ngo", ByteStart: 151, ByteEnd: 154, Start: Position{8, 3}, End: Position{9, 3}, Empty: true},
				{Text: "\ngo", ByteStart: 154, ByteEnd: 157, Start: Position{9, 3}, End: Position{10, 3}, Empty: true},
				{Text: "\ngo", ByteStart: 157, ByteEnd: 160, Start: Position{10, 3}, End: Position{11, 3}, Empty: true},
			},
		},
		{
			name:  "multiple UPDATE statements with blank lines",
			input: "\n\n\n\n\n\n\n\nUPDATE SalesLT.Address SET City = \"Shanghai\";\n\nUPDATE SalesLT.Address SET PostalCode = 0;\n\n\nUPDATE SalesLT.ProductModelProductDescription SET Culture = \"zh-cn\";\n\n",
			want: []Segment{
				{Text: "\n\n\n\n\n\n\n\nUPDATE SalesLT.Address SET City = \"Shanghai\";", ByteStart: 0, ByteEnd: 53, Start: Position{1, 1}, End: Position{9, 46}, Empty: false},
				{Text: "\n\nUPDATE SalesLT.Address SET PostalCode = 0;", ByteStart: 53, ByteEnd: 97, Start: Position{9, 46}, End: Position{11, 43}, Empty: false},
				{Text: "\n\n\nUPDATE SalesLT.ProductModelProductDescription SET Culture = \"zh-cn\";", ByteStart: 97, ByteEnd: 168, Start: Position{11, 43}, End: Position{14, 69}, Empty: false},
			},
		},
		{
			name:  "unicode in table names uses rune columns",
			input: "SELECT * FROM 表名; INSERT INTO 表 VALUES (1);",
			want: []Segment{
				{Text: "SELECT * FROM 表名;", ByteStart: 0, ByteEnd: 21, Start: Position{1, 1}, End: Position{1, 18}, Empty: false},
				{Text: " INSERT INTO 表 VALUES (1);", ByteStart: 21, ByteEnd: 49, Start: Position{1, 18}, End: Position{1, 44}, Empty: false},
			},
		},
		{
			name:  "leading newlines",
			input: "\n\nSELECT 1;",
			want: []Segment{
				{Text: "\n\nSELECT 1;", ByteStart: 0, ByteEnd: 11, Start: Position{1, 1}, End: Position{3, 10}, Empty: false},
			},
		},
		{
			name:  "leading spaces and newlines",
			input: "\n  \n  SELECT 1;",
			want: []Segment{
				{Text: "\n  \n  SELECT 1;", ByteStart: 0, ByteEnd: 15, Start: Position{1, 1}, End: Position{3, 12}, Empty: false},
			},
		},
		{
			name:  "multi-statement with leading whitespace",
			input: "SELECT 1;\n\n  SELECT 2;",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: "\n\n  SELECT 2;", ByteStart: 9, ByteEnd: 22, Start: Position{1, 10}, End: Position{3, 12}, Empty: false},
			},
		},
		{
			name:  "lone semicolon preserved as empty segment",
			input: ";",
			want: []Segment{
				{Text: ";", ByteStart: 0, ByteEnd: 1, Start: Position{1, 1}, End: Position{1, 2}, Empty: true},
			},
		},
		{
			name:  "double semicolon between statements preserved as empty segment",
			input: "SELECT 1;;SELECT 2;",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: ";", ByteStart: 9, ByteEnd: 10, Start: Position{1, 10}, End: Position{1, 11}, Empty: true},
				{Text: "SELECT 2;", ByteStart: 10, ByteEnd: 19, Start: Position{1, 11}, End: Position{1, 20}, Empty: false},
			},
		},
		{
			name:  "leading semicolon before WITH CTE preserved as empty segment",
			input: ";WITH cte AS (SELECT 1 AS a) SELECT a FROM cte;",
			want: []Segment{
				{Text: ";", ByteStart: 0, ByteEnd: 1, Start: Position{1, 1}, End: Position{1, 2}, Empty: true},
				{Text: "WITH cte AS (SELECT 1 AS a) SELECT a FROM cte;", ByteStart: 1, ByteEnd: 47, Start: Position{1, 2}, End: Position{1, 48}, Empty: false},
			},
		},
		{
			name:  "comment-only script preserved as empty segment",
			input: "-- just a comment",
			want: []Segment{
				{Text: "-- just a comment", ByteStart: 0, ByteEnd: 17, Start: Position{1, 1}, End: Position{1, 18}, Empty: true},
			},
		},
		{
			name:  "trailing comment after valid statement preserved as empty segment",
			input: "SELECT 1; -- trailing",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " -- trailing", ByteStart: 9, ByteEnd: 21, Start: Position{1, 10}, End: Position{1, 22}, Empty: true},
			},
		},
		{
			name:  "semicolon inside block comment between statements is not split",
			input: "SELECT 1; /* ; */ SELECT 2;",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " /* ; */ SELECT 2;", ByteStart: 9, ByteEnd: 27, Start: Position{1, 10}, End: Position{1, 28}, Empty: false},
			},
		},
		{
			name:  "unterminated trailing block comment after valid statement stays non-empty",
			input: "SELECT 1; /* unterminated",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " /* unterminated", ByteStart: 9, ByteEnd: 25, Start: Position{1, 10}, End: Position{1, 26}, Empty: false},
			},
		},
		{
			name:  "unterminated trailing string after valid statement stays non-empty",
			input: "SELECT 1; 'unterminated",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " 'unterminated", ByteStart: 9, ByteEnd: 23, Start: Position{1, 10}, End: Position{1, 24}, Empty: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.input)
			assertSplitInvariants(t, tt.input, got)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSplitFallbackUnparseableScripts(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Segment
	}{
		{
			name:  "syntax error in second statement preserves positions",
			input: "SELECT 1;\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: "\nSELECT FROM", ByteStart: 9, ByteEnd: 21, Start: Position{1, 10}, End: Position{2, 12}, Empty: false},
			},
		},
		{
			name:  "syntax error after non-ascii string",
			input: "SELECT N'你好';\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT N'你好';", ByteStart: 0, ByteEnd: 17, Start: Position{1, 1}, End: Position{1, 14}, Empty: false},
				{Text: "\nSELECT FROM", ByteStart: 17, ByteEnd: 29, Start: Position{1, 14}, End: Position{2, 12}, Empty: false},
			},
		},
		{
			name:  "semicolons inside single-quoted string with doubled quote and trailing backslash",
			input: `SELECT 'a;b''c\'; SELECT FROM`,
			want: []Segment{
				{Text: `SELECT 'a;b''c\';`, ByteStart: 0, ByteEnd: 17, Start: Position{1, 1}, End: Position{1, 18}, Empty: false},
				{Text: " SELECT FROM", ByteStart: 17, ByteEnd: 29, Start: Position{1, 18}, End: Position{1, 30}, Empty: false},
			},
		},
		{
			name:  "semicolons inside double-quoted and bracketed identifiers",
			input: `SELECT "a;""b", [c;]]d]; SELECT FROM`,
			want: []Segment{
				{Text: `SELECT "a;""b", [c;]]d];`, ByteStart: 0, ByteEnd: 24, Start: Position{1, 1}, End: Position{1, 25}, Empty: false},
				{Text: " SELECT FROM", ByteStart: 24, ByteEnd: 36, Start: Position{1, 25}, End: Position{1, 37}, Empty: false},
			},
		},
		{
			name:  "semicolons inside line and nested block comments",
			input: "-- comment ; ignored\nSELECT /* outer /* inner ; */ ; */ 1; SELECT FROM",
			want: []Segment{
				{Text: "-- comment ; ignored\nSELECT /* outer /* inner ; */ ; */ 1;", ByteStart: 0, ByteEnd: 58, Start: Position{1, 1}, End: Position{2, 38}, Empty: false},
				{Text: " SELECT FROM", ByteStart: 58, ByteEnd: 70, Start: Position{2, 38}, End: Position{2, 50}, Empty: false},
			},
		},
		{
			name:  "lone semicolons and trailing comment marked empty",
			input: "SELECT FROM;; -- trailing",
			want: []Segment{
				{Text: "SELECT FROM;", ByteStart: 0, ByteEnd: 12, Start: Position{1, 1}, End: Position{1, 13}, Empty: false},
				{Text: ";", ByteStart: 12, ByteEnd: 13, Start: Position{1, 13}, End: Position{1, 14}, Empty: true},
				{Text: " -- trailing", ByteStart: 13, ByteEnd: 25, Start: Position{1, 14}, End: Position{1, 26}, Empty: true},
			},
		},
		{
			name:  "GO batch separators in unparseable script",
			input: "SELECT 1\n-- sep\nGO\nGO 2\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\n-- sep\nGO", ByteStart: 8, ByteEnd: 18, Start: Position{1, 9}, End: Position{3, 3}, Empty: true},
				{Text: "\nGO 2", ByteStart: 18, ByteEnd: 23, Start: Position{3, 3}, End: Position{4, 5}, Empty: true},
				{Text: "\nSELECT FROM", ByteStart: 23, ByteEnd: 35, Start: Position{4, 5}, End: Position{5, 12}, Empty: false},
			},
		},
		{
			name:  "GO immediately followed by line comment without whitespace",
			input: "SELECT 1\nGO-- comment\nGO 2-- count comment\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\nGO", ByteStart: 8, ByteEnd: 11, Start: Position{1, 9}, End: Position{2, 3}, Empty: true},
				{Text: "-- comment\nGO 2", ByteStart: 11, ByteEnd: 26, Start: Position{2, 3}, End: Position{3, 5}, Empty: true},
				{Text: "-- count comment\nSELECT FROM", ByteStart: 26, ByteEnd: 54, Start: Position{3, 5}, End: Position{4, 12}, Empty: false},
			},
		},
		{
			name:  "GO preceded by leading block comment on its line or closing a multiline comment",
			input: "SELECT 1\n/* batch */ GO\nSELECT 2 /* multi\n*/ GO 2\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\n/* batch */ GO", ByteStart: 8, ByteEnd: 23, Start: Position{1, 9}, End: Position{2, 15}, Empty: true},
				{Text: "\nSELECT 2", ByteStart: 23, ByteEnd: 32, Start: Position{2, 15}, End: Position{3, 9}, Empty: false},
				{Text: " /* multi\n*/ GO 2", ByteStart: 32, ByteEnd: 49, Start: Position{3, 9}, End: Position{4, 8}, Empty: true},
				{Text: "\nSELECT FROM", ByteStart: 49, ByteEnd: 61, Start: Position{4, 8}, End: Position{5, 12}, Empty: false},
			},
		},
		{
			name:  "CR-only line endings recognize GO separators",
			input: "SELECT 1\rGO\rGO 2\rSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\rGO", ByteStart: 8, ByteEnd: 11, Start: Position{1, 9}, End: Position{1, 12}, Empty: true},
				{Text: "\rGO 2", ByteStart: 11, ByteEnd: 16, Start: Position{1, 12}, End: Position{1, 17}, Empty: true},
				{Text: "\rSELECT FROM", ByteStart: 16, ByteEnd: 28, Start: Position{1, 17}, End: Position{1, 29}, Empty: false},
			},
		},
		{
			name:  "GO followed by multiline trailing block comment before SQL on closing line",
			input: "SELECT 1\nGO /*\n*/ SELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\nGO", ByteStart: 8, ByteEnd: 11, Start: Position{1, 9}, End: Position{2, 3}, Empty: true},
				{Text: " /*\n*/ SELECT FROM", ByteStart: 11, ByteEnd: 29, Start: Position{2, 3}, End: Position{3, 15}, Empty: false},
			},
		},
		{
			name:  "GO followed by closed block comment with or without whitespace",
			input: "SELECT 1\nGO /* comment */\nGO/* tight */\nGO 2/* count */\nSELECT FROM",
			want: []Segment{
				{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
				{Text: "\nGO", ByteStart: 8, ByteEnd: 11, Start: Position{1, 9}, End: Position{2, 3}, Empty: true},
				{Text: " /* comment */\nGO", ByteStart: 11, ByteEnd: 28, Start: Position{2, 3}, End: Position{3, 3}, Empty: true},
				{Text: "/* tight */\nGO 2", ByteStart: 28, ByteEnd: 44, Start: Position{3, 3}, End: Position{4, 5}, Empty: true},
				{Text: "/* count */\nSELECT FROM", ByteStart: 44, ByteEnd: 67, Start: Position{4, 5}, End: Position{5, 12}, Empty: false},
			},
		},
		{
			name:  "rejected ODBC comment opener --(* stays non-empty",
			input: "SELECT 1; --(* invalid",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " --(* invalid", ByteStart: 9, ByteEnd: 22, Start: Position{1, 10}, End: Position{1, 23}, Empty: false},
			},
		},
		{
			name:  "GO followed by rejected --(* or unclosed block comment is not a batch separator",
			input: "GO--(* invalid\nGO /* unclosed",
			want: []Segment{
				{Text: "GO--(* invalid\nGO /* unclosed", ByteStart: 0, ByteEnd: 29, Start: Position{1, 1}, End: Position{2, 15}, Empty: false},
			},
		},
		{
			name:  "GO inside string or block comment does not split",
			input: "SELECT 'hello\nGO\nworld', /*\nGO\n*/ FROM",
			want: []Segment{
				{Text: "SELECT 'hello\nGO\nworld', /*\nGO\n*/ FROM", ByteStart: 0, ByteEnd: 38, Start: Position{1, 1}, End: Position{5, 8}, Empty: false},
			},
		},
		{
			name:  "GO with zero count or trailing identifier is not a batch separator",
			input: "GO 0\nGO foo\nSELECT FROM",
			want: []Segment{
				{Text: "GO 0\nGO foo\nSELECT FROM", ByteStart: 0, ByteEnd: 23, Start: Position{1, 1}, End: Position{3, 12}, Empty: false},
			},
		},
		{
			name:  "unterminated string stays non-empty",
			input: "SELECT 1; SELECT 'unterminated",
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: " SELECT 'unterminated", ByteStart: 9, ByteEnd: 30, Start: Position{1, 10}, End: Position{1, 31}, Empty: false},
			},
		},
		{
			name:  "unterminated block comment stays non-empty",
			input: "SELECT FROM; /* unterminated",
			want: []Segment{
				{Text: "SELECT FROM;", ByteStart: 0, ByteEnd: 12, Start: Position{1, 1}, End: Position{1, 13}, Empty: false},
				{Text: " /* unterminated", ByteStart: 12, ByteEnd: 28, Start: Position{1, 13}, End: Position{1, 29}, Empty: false},
			},
		},
		{
			name:  "unterminated double-quoted identifier stays non-empty",
			input: `SELECT 1; SELECT "unclosed`,
			want: []Segment{
				{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
				{Text: ` SELECT "unclosed`, ByteStart: 9, ByteEnd: 26, Start: Position{1, 10}, End: Position{1, 27}, Empty: false},
			},
		},
		{
			name:  "trailing whitespace after semicolon is omitted",
			input: "SELECT FROM;\n  \t\n",
			want: []Segment{
				{Text: "SELECT FROM;", ByteStart: 0, ByteEnd: 12, Start: Position{1, 1}, End: Position{1, 13}, Empty: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify Parse fails so Split exercises the fallback path.
			_, err := Parse(tt.input)
			require.Error(t, err, "expected Parse to fail for fallback test case")

			got := Split(tt.input)
			assertSplitInvariants(t, tt.input, got)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSplitBySemicolonDirect(t *testing.T) {
	require.Nil(t, SplitBySemicolon(""))
	require.Nil(t, SplitBySemicolon(" \t\r\n "))

	// Procedure body with internal semicolons stays intact under Split (parsed).
	procSQL := procBatchSQL(true, true)
	parsedSegs := Split(procSQL)
	assertSplitInvariants(t, procSQL, parsedSegs)
	require.Len(t, parsedSegs, 3)
	require.False(t, parsedSegs[0].Empty)
	require.True(t, parsedSegs[1].Empty)
	require.False(t, parsedSegs[2].Empty)

	// CRLF around GO in lexical fallback.
	crlf := "SELECT 1\r\nGO -- sep\r\nSELECT FROM"
	got := SplitBySemicolon(crlf)
	assertSplitInvariants(t, crlf, got)
	require.Equal(t, []Segment{
		{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
		{Text: "\r\nGO", ByteStart: 8, ByteEnd: 12, Start: Position{1, 9}, End: Position{2, 3}, Empty: true},
		{Text: " -- sep\r\nSELECT FROM", ByteStart: 12, ByteEnd: 32, Start: Position{2, 3}, End: Position{3, 12}, Empty: false},
	}, got)

	// Unterminated block comment in lexical fallback stays non-empty.
	unclosedComment := "SELECT 1; /* unterminated"
	gotComment := SplitBySemicolon(unclosedComment)
	assertSplitInvariants(t, unclosedComment, gotComment)
	require.Equal(t, []Segment{
		{Text: "SELECT 1;", ByteStart: 0, ByteEnd: 9, Start: Position{1, 1}, End: Position{1, 10}, Empty: false},
		{Text: " /* unterminated", ByteStart: 9, ByteEnd: 25, Start: Position{1, 10}, End: Position{1, 26}, Empty: false},
	}, gotComment)

	// CR-only line endings after line comments also restore GO eligibility.
	crComment := "SELECT 1\rGO\r-- sep\rGO 2\rSELECT FROM"
	gotCR := SplitBySemicolon(crComment)
	assertSplitInvariants(t, crComment, gotCR)
	require.Equal(t, []Segment{
		{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8, Start: Position{1, 1}, End: Position{1, 9}, Empty: false},
		{Text: "\rGO", ByteStart: 8, ByteEnd: 11, Start: Position{1, 9}, End: Position{1, 12}, Empty: true},
		{Text: "\r-- sep\rGO 2", ByteStart: 11, ByteEnd: 23, Start: Position{1, 12}, End: Position{1, 24}, Empty: true},
		{Text: "\rSELECT FROM", ByteStart: 23, ByteEnd: 35, Start: Position{1, 24}, End: Position{1, 36}, Empty: false},
	}, gotCR)

	// Repeated multiline block comments stay linear-time without rescanning.
	manyComments := strings.Repeat("/*\n*/", 16000) + "SELECT FROM"
	gotMany := SplitBySemicolon(manyComments)
	assertSplitInvariants(t, manyComments, gotMany)
	require.Len(t, gotMany, 1)
	require.False(t, gotMany[0].Empty)
}
