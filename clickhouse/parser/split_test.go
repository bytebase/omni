package parser

import (
	"errors"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		empty   []bool
		offsets [][2]int
	}{
		{
			name:   "empty script",
			script: "",
			want:   nil,
		},
		{
			name:    "single statement with semicolon",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 9}},
		},
		{
			name:    "single statement without semicolon",
			script:  "SELECT 1",
			want:    []string{"SELECT 1"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 8}},
		},
		{
			name:    "two statements on one line",
			script:  "SELECT 1; SELECT 2;",
			want:    []string{"SELECT 1;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:    "same-line trailing -- comment stays attached to statement",
			script:  "SELECT 1; -- note\nSELECT 2;",
			want:    []string{"SELECT 1; -- note", "\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 17}, {17, 27}},
		},
		{
			name:    "same-line trailing # comment stays attached to statement",
			script:  "SELECT 1; # note\nSELECT 2;",
			want:    []string{"SELECT 1; # note", "\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 16}, {16, 26}},
		},
		{
			name:    "same-line trailing comment with CRLF",
			script:  "SELECT 1; -- note\r\nSELECT 2;",
			want:    []string{"SELECT 1; -- note", "\r\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 17}, {17, 28}},
		},
		{
			name:    "next-line comment belongs to following statement",
			script:  "SELECT 1;\n-- note\nSELECT 2;",
			want:    []string{"SELECT 1;", "\n-- note\nSELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 27}},
		},
		{
			name:    "semicolons inside single-quoted string and escapes",
			script:  `SELECT 'a;b', 'it\'s;ok', 'x'';y'; SELECT 2;`,
			want:    []string{`SELECT 'a;b', 'it\'s;ok', 'x'';y';`, ` SELECT 2;`},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 34}, {34, 44}},
		},
		{
			name:    "semicolons inside dollar-quoted strings",
			script:  "SELECT $$a;b$$, $tag$c;d\ne;f$tag$; SELECT 2;",
			want:    []string{"SELECT $$a;b$$, $tag$c;d\ne;f$tag$;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 34}, {34, 44}},
		},
		{
			name:    "semicolons inside double-quoted and backtick identifiers",
			script:  "SELECT `a;b`, `c``;d`, \"e;f\" FROM t; SELECT 2;",
			want:    []string{"SELECT `a;b`, `c``;d`, \"e;f\" FROM t;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 36}, {36, 46}},
		},
		{
			name:    "semicolons inside nested block comments",
			script:  "/* outer /* inner ; */ still comment ; */ SELECT 1; /* tail /* ; */ */",
			want:    []string{"/* outer /* inner ; */ still comment ; */ SELECT 1;", " /* tail /* ; */ */"},
			empty:   []bool{false, true},
			offsets: [][2]int{{0, 51}, {51, 70}},
		},
		{
			name:    "empty semicolons and whitespace preserved as empty",
			script:  ";SELECT 1;;SELECT 2;;",
			want:    []string{";", "SELECT 1;", ";", "SELECT 2;", ";"},
			empty:   []bool{true, false, true, false, true},
			offsets: [][2]int{{0, 1}, {1, 10}, {10, 11}, {11, 20}, {20, 21}},
		},
		{
			name:    "comment only script is marked empty",
			script:  "-- line\n# hash\n/* outer /* inner */ outer */\n",
			want:    []string{"-- line\n# hash\n/* outer /* inner */ outer */\n"},
			empty:   []bool{true},
			offsets: [][2]int{{0, 45}},
		},
		{
			name:    "unclosed string consumes to EOF and stays non-empty",
			script:  "SELECT 1;\nSELECT 'unterminated; SELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT 'unterminated; SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 41}},
		},
		{
			name:    "unclosed dollar-quoted string consumes to EOF and stays non-empty",
			script:  "SELECT 1;\nSELECT $doc$unterminated; SELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT $doc$unterminated; SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 45}},
		},
		{
			name:    "unclosed nested block comment stays non-empty",
			script:  "SELECT 1; /* outer /* inner */ ; SELECT 2;",
			want:    []string{"SELECT 1;", " /* outer /* inner */ ; SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 42}},
		},
		{
			name:    "trailing escape in unterminated string terminates cleanly",
			script:  `SELECT 'abc\`,
			want:    []string{`SELECT 'abc\`},
			empty:   []bool{false},
			offsets: [][2]int{{0, 12}},
		},
		{
			name:    "invalid hex literal and unexpected dollar terminate cleanly",
			script:  "SELECT X'zz', $1; SELECT 2;",
			want:    []string{"SELECT X'zz', $1;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 17}, {17, 27}},
		},
		{
			name:    "insert format csv preserves semicolons in inline payload",
			script:  "INSERT INTO t FORMAT CSV\nhello;world\nfoo;bar\n",
			want:    []string{"INSERT INTO t FORMAT CSV\nhello;world\nfoo;bar\n"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 45}},
		},
		{
			name:    "insert format csv after earlier statement preserves payload",
			script:  "SELECT 1;\nINSERT INTO db.t (a, b) FORMAT CSV\na@b;1\n\"unclosed;2",
			want:    []string{"SELECT 1;", "\nINSERT INTO db.t (a, b) FORMAT CSV\na@b;1\n\"unclosed;2"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 9}, {9, 62}},
		},
		{
			name:    "insert format with settings including setting named format",
			script:  "INSERT INTO TABLE t (a) SETTINGS format = 'CSV', async_insert = 1 FORMAT JSONEachRow\n{\"a\":\"x;y\"}\n",
			want:    []string{"INSERT INTO TABLE t (a) SETTINGS format = 'CSV', async_insert = 1 FORMAT JSONEachRow\n{\"a\":\"x;y\"}\n"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 97}},
		},
		{
			name:    "insert into function format preserves inline payload",
			script:  "INSERT INTO FUNCTION s3('https://example.com/a;b', 'CSV') FORMAT CSV\n1;2\n",
			want:    []string{"INSERT INTO FUNCTION s3('https://example.com/a;b', 'CSV') FORMAT CSV\n1;2\n"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 73}},
		},
		{
			name:    "explain ast insert format preserves inline payload",
			script:  "EXPLAIN AST INSERT INTO t FORMAT CSV\nhello;world\n",
			want:    []string{"EXPLAIN AST INSERT INTO t FORMAT CSV\nhello;world\n"},
			empty:   []bool{false},
			offsets: [][2]int{{0, 49}},
		},
		{
			name:    "insert values still splits across statements",
			script:  "INSERT INTO t VALUES (1, 'a');\nINSERT INTO t VALUES (2, 'b');",
			want:    []string{"INSERT INTO t VALUES (1, 'a');", "\nINSERT INTO t VALUES (2, 'b');"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 30}, {30, 61}},
		},
		{
			name:    "insert parenthesized select still splits across statements",
			script:  "INSERT INTO t (SELECT 1); SELECT 2;",
			want:    []string{"INSERT INTO t (SELECT 1);", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 25}, {25, 35}},
		},
		{
			name:    "truncated insert before format still splits at semicolon",
			script:  "INSERT INTO t FORMAT; SELECT 2;",
			want:    []string{"INSERT INTO t FORMAT;", " SELECT 2;"},
			empty:   []bool{false, false},
			offsets: [][2]int{{0, 21}, {21, 31}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			segs := Split(tc.script)
			var texts []string
			prevEnd := 0
			for i, seg := range segs {
				texts = append(texts, seg.Text)
				if seg.ByteStart != prevEnd {
					t.Errorf("segment %d ByteStart = %d, want contiguous %d", i, seg.ByteStart, prevEnd)
				}
				if seg.Text != tc.script[seg.ByteStart:seg.ByteEnd] {
					t.Errorf("segment %d Text %q != script[%d:%d] (%q)", i, seg.Text, seg.ByteStart, seg.ByteEnd, tc.script[seg.ByteStart:seg.ByteEnd])
				}
				prevEnd = seg.ByteEnd
			}
			if prevEnd != len(tc.script) {
				t.Errorf("final ByteEnd = %d, want len(script) = %d", prevEnd, len(tc.script))
			}
			if !reflect.DeepEqual(texts, tc.want) {
				t.Fatalf("segments = %q, want %q", texts, tc.want)
			}
			for i, wantEmpty := range tc.empty {
				if got := segs[i].Empty(); got != wantEmpty {
					t.Errorf("segment %d Empty() = %v, want %v (text=%q)", i, got, wantEmpty, segs[i].Text)
				}
			}
			for i, off := range tc.offsets {
				if segs[i].ByteStart != off[0] || segs[i].ByteEnd != off[1] {
					t.Errorf("segment %d offsets = [%d,%d), want [%d,%d)", i, segs[i].ByteStart, segs[i].ByteEnd, off[0], off[1])
				}
			}
		})
	}
}

func TestLexerAndTokenize(t *testing.T) {
	toks, err := Tokenize("/* outer /* inner */ */ SELECT $tag$hi$tag$, 1;")
	if err != nil {
		t.Fatalf("Tokenize unexpected error: %v", err)
	}
	wantTypes := []TokenType{KwSelect, TokString, TokComma, TokInt, TokSemicolon, TokEOF}
	if len(toks) != len(wantTypes) {
		t.Fatalf("got %d tokens, want %d: %+v", len(toks), len(wantTypes), toks)
	}
	for i, wt := range wantTypes {
		if toks[i].Type != wt {
			t.Errorf("token %d type = %v, want %v", i, toks[i].Type, wt)
		}
	}

	// Streaming Lexer records first lexical error and terminates.
	l := NewLexer("SELECT 'unterminated")
	if tok := l.NextToken(); tok.Type != KwSelect {
		t.Fatalf("first token = %+v, want KwSelect", tok)
	}
	if tok := l.NextToken(); tok.Type != TokIllegal || tok.Pos != 7 {
		t.Fatalf("second token = %+v, want TokIllegal at 7", tok)
	}
	if tok := l.NextToken(); tok.Type != TokEOF {
		t.Fatalf("third token = %+v, want TokEOF", tok)
	}
	var pe *ParseError
	if !errors.As(l.Err(), &pe) || pe.Position != 7 {
		t.Fatalf("Err() = %v, want *ParseError at 7", l.Err())
	}
}

func TestParseExplainAnalyze(t *testing.T) {
	for _, tc := range []struct {
		sql          string
		wantKind     string
		wantSettings int
	}{
		{"EXPLAIN ANALYZE SELECT 1", "ANALYZE", 0},
		{"EXPLAIN ANALYZE SELECT count() FROM q", "ANALYZE", 0},
		{"EXPLAIN ANALYZE json = 1 SELECT 1", "ANALYZE", 1},
		{"EXPLAIN ANALYZE json = 1, indexes = 1 SELECT 1", "ANALYZE", 2},
		{"EXPLAIN analyze = 1 SELECT 1", "", 1},
	} {
		stmts, err := Parse(tc.sql)
		if err != nil {
			t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
		}
		if len(stmts) != 1 {
			t.Fatalf("Parse(%q) returned %d stmts, want 1", tc.sql, len(stmts))
		}
		ex, ok := stmts[0].(*ExplainStmt)
		if !ok {
			t.Fatalf("Parse(%q) returned %T, want *ExplainStmt", tc.sql, stmts[0])
		}
		if ex.Kind != tc.wantKind {
			t.Errorf("Parse(%q) Kind = %q, want %q", tc.sql, ex.Kind, tc.wantKind)
		}
		if len(ex.Settings) != tc.wantSettings {
			t.Errorf("Parse(%q) len(Settings) = %d, want %d", tc.sql, len(ex.Settings), tc.wantSettings)
		}
	}
}
