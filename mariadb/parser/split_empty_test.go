package parser

import "testing"

// nonEmptySegments returns the segments of Split(sql) that hold a statement,
// for tests that count statements rather than delimiters.
func nonEmptySegments(sql string) []Segment {
	var segs []Segment
	for _, seg := range Split(sql) {
		if !seg.Empty() {
			segs = append(segs, seg)
		}
	}
	return segs
}

// TestSplitKeepsEmptySegments pins that Split returns a segment holding no
// statement, marked Empty, so a caller can account for every delimiter.
func TestSplitKeepsEmptySegments(t *testing.T) {
	type seg struct {
		text       string
		start, end int
		empty      bool
	}
	tests := []struct {
		name string
		sql  string
		want []seg
	}{
		{"lone semicolon", ";", []seg{{"", 0, 0, true}}},
		{"two semicolons", ";;", []seg{{"", 0, 0, true}, {"", 1, 1, true}}},
		{"double semicolon after a statement", "SELECT 1;;", []seg{{"SELECT 1", 0, 8, false}, {"", 9, 9, true}}},
		{"empty between statements", "SELECT 1; ; SELECT 2", []seg{
			{"SELECT 1", 0, 8, false}, {" ", 9, 10, true}, {" SELECT 2", 11, 20, false},
		}},
		{"comments before a semicolon", "/* a */ -- b\n# c\n;SELECT 1", []seg{
			{"/* a */ -- b\n# c\n", 0, 17, true}, {"SELECT 1", 18, 26, false},
		}},
		{"comment-only tail", "SELECT 1; -- done", []seg{{"SELECT 1", 0, 8, false}, {" -- done", 9, 17, true}}},
		{"whitespace-only tail", "SELECT 1;\n", []seg{{"SELECT 1", 0, 8, false}, {"\n", 9, 10, true}}},
		{"nothing after the last semicolon", "SELECT 1;", []seg{{"SELECT 1", 0, 8, false}}},
		{"whitespace only", "  ", []seg{{"  ", 0, 2, true}}},
		{"comment only", "/* c */", []seg{{"/* c */", 0, 7, true}}},
		{"empty input", "", nil},
		{"executable comment is a statement", "/*!40101 SET @a = 1 */;;", []seg{
			{"/*!40101 SET @a = 1 */", 0, 22, false}, {"", 23, 23, true},
		}},
		{"unterminated comment is not empty", "; /* open", []seg{{"", 0, 0, true}, {" /* open", 1, 9, false}}},
		// The whitespace before a DELIMITER directive belongs to no segment,
		// and an empty segment under a custom delimiter is kept like one
		// under ';'.
		{"delimiter block", "DELIMITER $$\nSELECT 1$$$$\n$$\nDELIMITER ;\n;", []seg{
			{"SELECT 1", 13, 21, false}, {"", 23, 23, true}, {"\n", 25, 26, true}, {"", 41, 41, true},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := Split(tt.sql)
			var got []seg
			for _, s := range segs {
				if tt.sql[s.ByteStart:s.ByteEnd] != s.Text {
					t.Errorf("sql[%d:%d] = %q, but Text = %q", s.ByteStart, s.ByteEnd, tt.sql[s.ByteStart:s.ByteEnd], s.Text)
				}
				got = append(got, seg{s.Text, s.ByteStart, s.ByteEnd, s.Empty()})
			}
			if len(got) != len(tt.want) {
				t.Fatalf("Split(%q) = %+v, want %+v", tt.sql, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Split(%q)[%d] = %+v, want %+v", tt.sql, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// TestParseSkipsEmptySegments pins that the empty segments Split keeps add
// no statement to Parse.
func TestParseSkipsEmptySegments(t *testing.T) {
	list, err := Parse(";; SELECT 1;; /* c */ ; SELECT 2; -- end")
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 2 {
		t.Fatalf("Parse returned %d statements, want 2", len(list.Items))
	}
}
