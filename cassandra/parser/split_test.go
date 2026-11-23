package parser

import "testing"

// TestSplitKeepsEmptySegments pins that Split returns a segment holding no
// statement, marked Empty, so a caller can account for every semicolon. A
// segment's text and range exclude the whitespace around it.
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
		{"empty input", "", nil},
		{"lone semicolon", ";", []seg{{"", 0, 0, true}}},
		{"two semicolons", ";;", []seg{{"", 0, 0, true}, {"", 1, 1, true}}},
		{"double semicolon after a statement", "SELECT 1;;", []seg{{"SELECT 1", 0, 8, false}, {"", 9, 9, true}}},
		{"whitespace between statements", "SELECT 1; ; SELECT 2", []seg{
			{"SELECT 1", 0, 8, false}, {"", 10, 10, true}, {"SELECT 2", 12, 20, false},
		}},
		{"nothing after the last semicolon", "SELECT 1;", []seg{{"SELECT 1", 0, 8, false}}},
		{"whitespace-only tail", "SELECT 1;\n", []seg{{"SELECT 1", 0, 8, false}, {"", 10, 10, true}}},
		{"whitespace only", "  ", []seg{{"", 2, 2, true}}},
		{"dash comment before a semicolon", "-- c\n;SELECT 1", []seg{{"-- c", 0, 4, true}, {"SELECT 1", 6, 14, false}}},
		{"slash comment before a semicolon", "// c\n;SELECT 1", []seg{{"// c", 0, 4, true}, {"SELECT 1", 6, 14, false}}},
		{"block comment before a semicolon", "/* c */;SELECT 1", []seg{{"/* c */", 0, 7, true}, {"SELECT 1", 8, 16, false}}},
		{"comments of every style", " /* a */ -- b\n// c\n ;", []seg{{"/* a */ -- b\n// c", 1, 18, true}}},
		{"dash comment tail", "SELECT 1; -- done", []seg{{"SELECT 1", 0, 8, false}, {"-- done", 10, 17, true}}},
		{"slash comment tail", "SELECT 1; // done", []seg{{"SELECT 1", 0, 8, false}, {"// done", 10, 17, true}}},
		{"block comment tail", "SELECT 1;\n/* done */\n", []seg{{"SELECT 1", 0, 8, false}, {"/* done */", 10, 20, true}}},
		{"comment only", "/* c */", []seg{{"/* c */", 0, 7, true}}},
		{"a comment before a statement belongs to it", "-- c\nSELECT 1", []seg{{"-- c\nSELECT 1", 0, 13, false}}},
		{"a semicolon in a slash comment ends nothing", "SELECT 1 // a;b\n;SELECT 2", []seg{
			{"SELECT 1 // a;b", 0, 15, false}, {"SELECT 2", 17, 25, false},
		}},
		{"unterminated comment is not empty", "; /* open", []seg{{"", 0, 0, true}, {"/* open", 2, 9, false}}},
		{"semicolon in a string", "SELECT ';' FROM t;;", []seg{{"SELECT ';' FROM t", 0, 17, false}, {"", 18, 18, true}}},
		{"batch then empty segments", "BEGIN BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH;; -- end", []seg{
			{"BEGIN BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH", 0, 54, false},
			{"", 55, 55, true},
			{"-- end", 57, 63, true},
		}},
		{"empty segment before a batch", ";BEGIN BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH", []seg{
			{"", 0, 0, true},
			{"BEGIN BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH", 1, 55, false},
		}},
		{"comments inside a batch stay in it", "BEGIN BATCH // a;\nINSERT INTO t (id) VALUES (1); -- b;\nAPPLY BATCH;", []seg{
			{"BEGIN BATCH // a;\nINSERT INTO t (id) VALUES (1); -- b;\nAPPLY BATCH", 0, 66, false},
		}},
		{"batch that closes at once", "BEGIN BATCH APPLY BATCH;;", []seg{{"BEGIN BATCH APPLY BATCH", 0, 23, false}, {"", 24, 24, true}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Split(tt.sql)
			if len(got) != len(tt.want) {
				t.Fatalf("Split(%q): got %d segments %+v, want %d", tt.sql, len(got), got, len(tt.want))
			}
			for i, w := range tt.want {
				g := got[i]
				if g.Text != w.text || g.ByteStart != w.start || g.ByteEnd != w.end || g.Empty != w.empty {
					t.Errorf("segment %d: got {%q %d %d %v}, want {%q %d %d %v}",
						i, g.Text, g.ByteStart, g.ByteEnd, g.Empty, w.text, w.start, w.end, w.empty)
				}
				if g.Text != tt.sql[g.ByteStart:g.ByteEnd] {
					t.Errorf("segment %d: Text %q is not the input's range %d..%d", i, g.Text, g.ByteStart, g.ByteEnd)
				}
			}
		})
	}
}
