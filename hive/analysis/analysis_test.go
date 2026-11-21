package analysis

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
		err     string
	}{
		{
			name:    "one statement",
			script:  "SELECT 1;",
			want:    []string{"SELECT 1;"},
			offsets: [][2]int{{0, 9}},
		},
		{
			name:    "two statements",
			script:  "SELECT 1;\nSELECT 2;",
			want:    []string{"SELECT 1;", "\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 19}},
		},
		{
			name:   "no terminator",
			script: "SELECT 1",
			want:   []string{"SELECT 1"},
		},
		{
			name:   "a comment line belongs to the statement after it",
			script: "-- a comment\nSELECT 1;",
			want:   []string{"-- a comment\nSELECT 1;"},
		},
		{
			name:   "a comment after the last statement belongs to none",
			script: "SELECT 1;\n-- trailing\n",
			want:   []string{"SELECT 1;"},
		},
		{
			name:   "an empty script",
			script: "\n\n",
			want:   []string{},
		},
		{
			name:   "a block comment that ends its line",
			script: "/* a comment */\nSELECT 1;",
			want:   []string{"/* a comment */\nSELECT 1;"},
		},
		{
			name:   "a block comment whose closing does not end its line",
			script: "/* a comment */ SELECT 1;",
			err:    "`*/` must be the end of the line",
		},
		{
			name:    "CRLF multiple statements",
			script:  "SELECT 1;\r\nSELECT 2;",
			want:    []string{"SELECT 1;", "\r\nSELECT 2;"},
			offsets: [][2]int{{0, 9}, {9, 20}},
		},
		{
			name:    "multi-byte UTF-8",
			script:  "SELECT '中';\nSELECT '文';",
			want:    []string{"SELECT '中';", "\nSELECT '文';"},
			offsets: [][2]int{{0, 13}, {13, 27}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			statements, err := Split(tc.script)
			if tc.err != "" {
				if err == nil || !strings.Contains(err.Error(), tc.err) {
					t.Fatalf("Split error = %v, want %q", err, tc.err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Split: %v", err)
			}
			texts := make([]string, 0, len(statements))
			for _, stmt := range statements {
				texts = append(texts, stmt.Text)
				if stmt.Text != tc.script[stmt.Start:stmt.End] {
					t.Errorf("statement text %q does not match its offsets", stmt.Text)
				}
			}
			if len(texts) != len(tc.want) {
				t.Fatalf("got %d statements %q, want %d", len(texts), texts, len(tc.want))
			}
			for i := range texts {
				if texts[i] != tc.want[i] {
					t.Errorf("statement %d = %q, want %q", i, texts[i], tc.want[i])
				}
			}
			for i, offsets := range tc.offsets {
				if statements[i].Start != offsets[0] || statements[i].End != offsets[1] {
					t.Errorf("statement %d offsets = [%d,%d), want [%d,%d)",
						i, statements[i].Start, statements[i].End, offsets[0], offsets[1])
				}
			}
		})
	}
}

func TestExplainStatement(t *testing.T) {
	cases := []struct {
		statement string
		format    string
		want      string
	}{
		{"SELECT 1", "", "EXPLAIN SELECT 1"},
		{"SELECT 1;", "", "EXPLAIN SELECT 1;"},
		{"SELECT 1", "TEXT", "EXPLAIN SELECT 1"},
		{"SELECT 1", "XML", "EXPLAIN SELECT 1"},
		{"SELECT 1", "JSON", "EXPLAIN FORMATTED SELECT 1"},
		{"select 1", "json", "EXPLAIN FORMATTED select 1"},
		{"EXPLAIN SELECT 1", "JSON", "EXPLAIN SELECT 1"},
		{"EXPLAIN EXTENDED SELECT 1", "JSON", "EXPLAIN EXTENDED SELECT 1"},
		{"explain formatted select 1", "", "explain formatted select 1"},
		{"/* a comment */\nEXPLAIN SELECT 1", "JSON", "/* a comment */\nEXPLAIN SELECT 1"},
		{"-- plan it\nEXPLAIN SELECT 1", "", "-- plan it\nEXPLAIN SELECT 1"},
		{"EXPLAINED", "", "EXPLAIN EXPLAINED"},
		{"explained_view_query()", "", "EXPLAIN explained_view_query()"},
		{"/* /* */ EXPLAIN */ SELECT 1", "JSON", "EXPLAIN FORMATTED /* /* */ EXPLAIN */ SELECT 1"},
	}
	for _, tc := range cases {
		if got := ExplainStatement(tc.statement, tc.format); got != tc.want {
			t.Errorf("ExplainStatement(%q, %q) = %q, want %q", tc.statement, tc.format, got, tc.want)
		}
	}
}

func TestPlanFormat(t *testing.T) {
	cases := []struct {
		statement string
		format    string
		ok        bool
	}{
		{"EXPLAIN SELECT 1", "TEXT", true},
		{"/* c */\nEXPLAIN FORMATTED SELECT 1", "JSON", true},
		{"EXPLAIN /* c */ FORMATTED SELECT 1", "JSON", true},
		{"EXPLAIN ANALYZE SELECT 1", "TEXT", true},
		{"EXPLAINED", "", false},
		{"SELECT 1", "", false},
		{"-- only a comment", "", false},
	}
	for _, tc := range cases {
		format, ok := PlanFormat(tc.statement)
		if ok != tc.ok || format != tc.format {
			t.Errorf("PlanFormat(%q) = (%q, %v), want (%q, %v)", tc.statement, format, ok, tc.format, tc.ok)
		}
	}
}
