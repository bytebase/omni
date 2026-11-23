package parser

import (
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name    string
		script  string
		want    []string
		offsets [][2]int
		empty   []bool
	}{
		{
			name:   "empty string",
			script: "",
			want:   nil,
		},
		{
			name:    "single statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			offsets: [][2]int{{0, 9}},
			empty:   []bool{false},
		},
		{
			name:    "single statement without semicolon",
			script:  "SELECT 1",
			want:    []string{"SELECT 1"},
			offsets: [][2]int{{0, 8}},
			empty:   []bool{false},
		},
		{
			name:    "two statements on one line",
			script:  "SELECT 1; SELECT 2;",
			want:    []string{"SELECT 1;", " SELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 19}},
			empty:   []bool{false, false},
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 20}},
			empty:   []bool{false, false},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			offsets: [][2]int{{0, 13}, {13, 27}},
			empty:   []bool{false, false},
		},
		{
			name:    "comment only script is preserved as empty",
			script:  "\n-- only a comment\n",
			want:    []string{"\n-- only a comment\n"},
			offsets: [][2]int{{0, 19}},
			empty:   []bool{true},
		},
		{
			name:    "lone semicolons are empty",
			script:  ";;\n",
			want:    []string{";", ";", "\n"},
			offsets: [][2]int{{0, 1}, {1, 2}, {2, 3}},
			empty:   []bool{true, true, true},
		},
		{
			name:    "inline block comment before statement",
			script:  "/* a comment */ SELECT 1;",
			want:    []string{"/* a comment */ SELECT 1;"},
			offsets: [][2]int{{0, 25}},
			empty:   []bool{false},
		},
		{
			name:    "leading inline block comment between statements on same line stays with next statement",
			script:  "SELECT 1; /* lead */ SELECT 2;",
			want:    []string{"SELECT 1;", " /* lead */ SELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 30}},
			empty:   []bool{false, false},
		},
		{
			name:    "trailing same-line line comment stays with preceding statement",
			script:  "SELECT 1; -- note\nSELECT 2;",
			want:    []string{"SELECT 1; -- note", "\nSELECT 2;"},
			offsets: [][2]int{{0, 17}, {17, 27}},
			empty:   []bool{false, false},
		},
		{
			name:    "trailing same-line line comment with backslash continuation",
			script:  "SELECT 1; -- line1\\\nline2\nSELECT 2;",
			want:    []string{"SELECT 1; -- line1\\\nline2", "\nSELECT 2;"},
			offsets: [][2]int{{0, 25}, {25, 35}},
			empty:   []bool{false, false},
		},
		{
			name:    "trailing same-line block comment ending its line stays with preceding statement",
			script:  "SELECT 1; /* note */\nSELECT 2;",
			want:    []string{"SELECT 1; /* note */", "\nSELECT 2;"},
			offsets: [][2]int{{0, 20}, {20, 30}},
			empty:   []bool{false, false},
		},
		{
			name:    "trailing multiline block comment after terminator",
			script:  "SELECT 1; /* note\ncontinued */\nSELECT 2;",
			want:    []string{"SELECT 1; /* note\ncontinued */", "\nSELECT 2;"},
			offsets: [][2]int{{0, 30}, {30, 40}},
			empty:   []bool{false, false},
		},
		{
			name:    "nested block comments with semicolons inside",
			script:  "/* outer /* inner ; */ ; */ SELECT 1; SELECT 2;",
			want:    []string{"/* outer /* inner ; */ ; */ SELECT 1;", " SELECT 2;"},
			offsets: [][2]int{{0, 37}, {37, 47}},
			empty:   []bool{false, false},
		},
		{
			name:    "optimizer hint with semicolons in strings",
			script:  "SELECT /*+ BROADCAST(t) */ 'a;b', \"c;d\", r'e;f', R\"g;h\", `i;j`, `k``;l` FROM t; SELECT 2;",
			want:    []string{"SELECT /*+ BROADCAST(t) */ 'a;b', \"c;d\", r'e;f', R\"g;h\", `i;j`, `k``;l` FROM t;", " SELECT 2;"},
			offsets: [][2]int{{0, 79}, {79, 89}},
			empty:   []bool{false, false},
		},
		{
			name:    "escaped quotes and raw string trailing backslash",
			script:  `SELECT 'it\'s;ok', r'foo\' ; SELECT 2;`,
			want:    []string{`SELECT 'it\'s;ok', r'foo\' ;`, ` SELECT 2;`},
			offsets: [][2]int{{0, 28}, {28, 38}},
			empty:   []bool{false, false},
		},
		{
			name:    "compound BEGIN END block stays together",
			script:  "BEGIN SELECT 1; SELECT 2; END; SELECT 3;",
			want:    []string{"BEGIN SELECT 1; SELECT 2; END;", " SELECT 3;"},
			offsets: [][2]int{{0, 30}, {30, 40}},
			empty:   []bool{false, false},
		},
		{
			name:    "compound BEGIN NOT ATOMIC END",
			script:  "BEGIN NOT ATOMIC DECLARE x INT = 1; SET x = 2; SELECT x; END; SELECT 3;",
			want:    []string{"BEGIN NOT ATOMIC DECLARE x INT = 1; SET x = 2; SELECT x; END;", " SELECT 3;"},
			offsets: [][2]int{{0, 61}, {61, 71}},
			empty:   []bool{false, false},
		},
		{
			name:    "empty BEGIN END block",
			script:  "BEGIN END; SELECT 1;",
			want:    []string{"BEGIN END;", " SELECT 1;"},
			offsets: [][2]int{{0, 10}, {10, 20}},
			empty:   []bool{false, false},
		},
		{
			name:    "labeled compound BEGIN END block",
			script:  "outer_blk: BEGIN SELECT 1; SELECT 2; END outer_blk; SELECT 3;",
			want:    []string{"outer_blk: BEGIN SELECT 1; SELECT 2; END outer_blk;", " SELECT 3;"},
			offsets: [][2]int{{0, 51}, {51, 61}},
			empty:   []bool{false, false},
		},
		{
			name: "nested compound blocks with IF CASE WHILE REPEAT LOOP FOR and SQL CASE expression",
			script: strings.Join([]string{
				"BEGIN",
				"  SELECT CASE WHEN x = 1 THEN 'a' ELSE 'b' END FROM t;",
				"  IF x > 0 THEN",
				"    BEGIN",
				"      SELECT 1;",
				"    END;",
				"  ELSEIF CASE WHEN y = 1 THEN true ELSE false END THEN",
				"    SELECT 2;",
				"  ELSE",
				"    CASE x WHEN 1 THEN SELECT 3; ELSE SELECT 4; END CASE;",
				"  END IF;",
				"  lbl: WHILE x < 10 DO",
				"    SET x = x + 1;",
				"  END WHILE lbl;",
				"  REPEAT SELECT 5; UNTIL x = 10 END REPEAT;",
				"  LOOP LEAVE lbl; END LOOP;",
				"  FOR r AS SELECT 1 DO SELECT r; END FOR;",
				"END;",
				"SELECT 99;",
			}, "\n"),
			want: []string{
				strings.Join([]string{
					"BEGIN",
					"  SELECT CASE WHEN x = 1 THEN 'a' ELSE 'b' END FROM t;",
					"  IF x > 0 THEN",
					"    BEGIN",
					"      SELECT 1;",
					"    END;",
					"  ELSEIF CASE WHEN y = 1 THEN true ELSE false END THEN",
					"    SELECT 2;",
					"  ELSE",
					"    CASE x WHEN 1 THEN SELECT 3; ELSE SELECT 4; END CASE;",
					"  END IF;",
					"  lbl: WHILE x < 10 DO",
					"    SET x = x + 1;",
					"  END WHILE lbl;",
					"  REPEAT SELECT 5; UNTIL x = 10 END REPEAT;",
					"  LOOP LEAVE lbl; END LOOP;",
					"  FOR r AS SELECT 1 DO SELECT r; END FOR;",
					"END;",
				}, "\n"),
				"\nSELECT 99;",
			},
			empty: []bool{false, false},
		},
		{
			name:   "DECLARE EXIT HANDLER with BEGIN END inside compound block",
			script: "BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN SELECT 1; END; SELECT 2; END; SELECT 3;",
			want:   []string{"BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION BEGIN SELECT 1; END; SELECT 2; END;", " SELECT 3;"},
			empty:  []bool{false, false},
		},
		{
			name:   "DECLARE EXIT HANDLER with labeled BEGIN END",
			script: "BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION h: BEGIN SELECT 1; END h; SELECT 2; END; SELECT 3;",
			want:   []string{"BEGIN DECLARE EXIT HANDLER FOR SQLEXCEPTION h: BEGIN SELECT 1; END h; SELECT 2; END;", " SELECT 3;"},
			empty:  []bool{false, false},
		},
		{
			name:    "begin and end as non-reserved identifiers do not merge statements",
			script:  "SELECT begin, end FROM t; SELECT 2;",
			want:    []string{"SELECT begin, end FROM t;", " SELECT 2;"},
			offsets: [][2]int{{0, 25}, {25, 35}},
			empty:   []bool{false, false},
		},
		{
			name:   "begin identifier inside compound block does not increase depth",
			script: "BEGIN SELECT begin FROM t; END; SELECT 2;",
			want:   []string{"BEGIN SELECT begin FROM t; END;", " SELECT 2;"},
			empty:  []bool{false, false},
		},
		{
			name:    "unclosed block comment stays non-empty",
			script:  "/* unclosed",
			want:    []string{"/* unclosed"},
			offsets: [][2]int{{0, 11}},
			empty:   []bool{false},
		},
		{
			name:    "unclosed string stays non-empty",
			script:  "SELECT 'unterminated",
			want:    []string{"SELECT 'unterminated"},
			offsets: [][2]int{{0, 20}},
			empty:   []bool{false},
		},
		{
			name:    "unclosed compound block does not loop",
			script:  "BEGIN SELECT 1; SELECT 2;",
			want:    []string{"BEGIN SELECT 1; SELECT 2;"},
			offsets: [][2]int{{0, 25}},
			empty:   []bool{false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			segs := Split(tc.script)
			if len(segs) != len(tc.want) {
				got := make([]string, len(segs))
				for i, s := range segs {
					got[i] = s.Text
				}
				t.Fatalf("Split(%q) got %d segments %q, want %d %q", tc.script, len(segs), got, len(tc.want), tc.want)
			}
			prevEnd := 0
			for i, seg := range segs {
				if seg.Text != tc.script[seg.ByteStart:seg.ByteEnd] {
					t.Errorf("segment %d text %q != script[%d:%d]", i, seg.Text, seg.ByteStart, seg.ByteEnd)
				}
				if seg.ByteStart != prevEnd {
					t.Errorf("segment %d ByteStart = %d, want contiguous %d", i, seg.ByteStart, prevEnd)
				}
				prevEnd = seg.ByteEnd
				if seg.Text != tc.want[i] {
					t.Errorf("segment %d = %q, want %q", i, seg.Text, tc.want[i])
				}
				if i < len(tc.offsets) {
					if seg.ByteStart != tc.offsets[i][0] || seg.ByteEnd != tc.offsets[i][1] {
						t.Errorf("segment %d offsets = [%d,%d), want [%d,%d)",
							i, seg.ByteStart, seg.ByteEnd, tc.offsets[i][0], tc.offsets[i][1])
					}
				}
				if i < len(tc.empty) {
					if got := seg.Empty(); got != tc.empty[i] {
						t.Errorf("segment %d Empty() = %v, want %v", i, got, tc.empty[i])
					}
				}
			}
			if len(tc.script) > 0 && prevEnd != len(tc.script) {
				t.Errorf("segments covered %d bytes, want %d", prevEnd, len(tc.script))
			}
		})
	}
}

func TestTokenizeAndLexer(t *testing.T) {
	toks, err := Tokenize("/* outer /* inner */ */ SELECT /*+ BROADCAST(t) */ `a``b`, r'x\\y' FROM t |> LIMIT 10")
	if err != nil {
		t.Fatalf("Tokenize: %v", err)
	}
	if len(toks) == 0 || toks[0].Kind != KwSelect {
		t.Fatalf("first token = %+v, want KwSelect", toks[0])
	}
	if toks[len(toks)-1].Kind != TokEOF {
		t.Fatalf("last token = %+v, want TokEOF", toks[len(toks)-1])
	}

	if _, err := Tokenize("/* unclosed /* inner */"); err == nil {
		t.Fatal("Tokenize accepted unclosed nested block comment")
	}

	l := NewLexer("/* unclosed")
	if tok := l.NextToken(); tok.Kind != TokEOF {
		t.Errorf("unclosed comment NextToken() = %+v, want TokEOF", tok)
	}
	if l.Err() == nil {
		t.Error("unclosed comment Err() = nil, want error")
	}
}
