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
		{"SELECT * FROM users; ; INSERT INTO t (id) VALUES (1)", 2}, // empty segment filtered
		{"SELECT * FROM users WHERE name = 'hello;world'", 1},       // semicolon inside string
		{`SELECT * FROM users WHERE name = "test;col"`, 1},          // semicolon inside quoted ident
		{"SELECT * FROM users; -- comment with ; inside\nINSERT INTO t (id) VALUES (1)", 2},
		{"SELECT * FROM users /* ; */ ; SELECT 1", 2},
		{"CREATE FUNCTION f() RETURNS NULL ON NULL INPUT RETURNS text LANGUAGE java AS $$return \";\";$$", 1}, // semicolon in code block
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
