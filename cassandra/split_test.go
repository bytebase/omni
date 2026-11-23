package cassandra

import (
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"SELECT * FROM users", 1},
		{"SELECT * FROM users;", 1},
		{"SELECT * FROM users; INSERT INTO t (id) VALUES (1)", 2},
		{"SELECT * FROM users; ; INSERT INTO t (id) VALUES (1)", 3}, // the empty segment is kept
		{"SELECT * FROM users WHERE name = 'hello;world'", 1},       // semicolon inside string
		{`SELECT * FROM users WHERE name = "test;col"`, 1},          // semicolon inside quoted ident
		{"SELECT * FROM users; -- comment with ; inside\nINSERT INTO t (id) VALUES (1)", 2},
		{"SELECT * FROM users /* ; */ ; SELECT 1", 2},
		{"CREATE FUNCTION f() RETURNS NULL ON NULL INPUT RETURNS text LANGUAGE java AS $$return \";\";$$", 1}, // semicolon in code block
		{"BEGIN BATCH INSERT INTO t (id) VALUES (1); INSERT INTO t (id) VALUES (2); APPLY BATCH;", 1},         // a batch is one statement
		{"BEGIN UNLOGGED BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH; SELECT 1", 2},
		{"BEGIN COUNTER BATCH UPDATE t SET c = c + 1 WHERE id = 1; APPLY BATCH", 1},
		{"BEGIN BATCH INSERT INTO t (id) VALUES (1); -- one; two\nAPPLY BATCH;", 1}, // a semicolon in a comment does not close a batch
		{"BEGIN BATCH INSERT INTO t (id) VALUES (1); INSERT INTO t (id) VALUES (2); INSERT INTO t (id) VALUES (3); APPLY BATCH; SELECT 1", 2},
		{"BEGIN BATCH APPLY BATCH; SELECT 1", 2}, // a batch that closes at once is one statement too
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			segs := Split(tt.input)
			if len(segs) != tt.expected {
				t.Errorf("Split(%q): got %d segments, want %d", tt.input, len(segs), tt.expected)
				for i, s := range segs {
					t.Logf("  segment %d: %q (empty=%v)", i, s.Text, s.Empty)
				}
			}
		})
	}
}

func TestSplitPositions(t *testing.T) {
	input := "SELECT 1; SELECT 2"
	segs := Split(input)
	if len(segs) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segs))
	}
	if segs[0].ByteStart != 0 || segs[0].ByteEnd != 8 {
		t.Errorf("seg 0: ByteStart=%d ByteEnd=%d, want 0..8", segs[0].ByteStart, segs[0].ByteEnd)
	}
	if segs[0].Text != "SELECT 1" {
		t.Errorf("seg 0: Text=%q, want %q", segs[0].Text, "SELECT 1")
	}
	if segs[1].ByteStart != 10 || segs[1].ByteEnd != 18 {
		t.Errorf("seg 1: ByteStart=%d ByteEnd=%d, want 10..18", segs[1].ByteStart, segs[1].ByteEnd)
	}
}

// TestSplitMarksEmptySegments holds the wrapper to the parser's segments: a
// segment holding no statement is kept and marked Empty.
func TestSplitMarksEmptySegments(t *testing.T) {
	input := "SELECT 1;;-- c\n;SELECT 2;\n-- tail"
	want := []Segment{
		{Text: "SELECT 1", ByteStart: 0, ByteEnd: 8},
		{Text: "", ByteStart: 9, ByteEnd: 9, Empty: true},
		{Text: "-- c", ByteStart: 10, ByteEnd: 14, Empty: true},
		{Text: "SELECT 2", ByteStart: 16, ByteEnd: 24},
		{Text: "-- tail", ByteStart: 26, ByteEnd: 33, Empty: true},
	}
	got := Split(input)
	if len(got) != len(want) {
		t.Fatalf("Split(%q): got %d segments %+v, want %d", input, len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("segment %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
