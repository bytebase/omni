package cosmosdb

import "testing"

func TestSplit(t *testing.T) {
	cases := []struct {
		name      string
		script    string
		wantText  []string
		wantEmpty []bool
	}{
		{
			name:      "one statement",
			script:    "SELECT * FROM c",
			wantText:  []string{"SELECT * FROM c"},
			wantEmpty: []bool{false},
		},
		{
			name:      "two statements",
			script:    "SELECT * FROM c; SELECT * FROM c WHERE c.country = 'US'",
			wantText:  []string{"SELECT * FROM c", " SELECT * FROM c WHERE c.country = 'US'"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "terminator is not statement text",
			script:    "SELECT * FROM c;",
			wantText:  []string{"SELECT * FROM c"},
			wantEmpty: []bool{false},
		},
		{
			name:      "a lone semicolon is an empty segment",
			script:    "SELECT * FROM c;;",
			wantText:  []string{"SELECT * FROM c", ""},
			wantEmpty: []bool{false, true},
		},
		{
			name:      "semicolon in a double-quoted string",
			script:    `SELECT * FROM c WHERE c.name = "a;b"`,
			wantText:  []string{`SELECT * FROM c WHERE c.name = "a;b"`},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon in a single-quoted string",
			script:    `SELECT * FROM c WHERE c.name = 'a;b'`,
			wantText:  []string{`SELECT * FROM c WHERE c.name = 'a;b'`},
			wantEmpty: []bool{false},
		},
		{
			name:      "quote escaped in a string",
			script:    `SELECT * FROM c WHERE c.name = 'it\'s; here'; SELECT 1`,
			wantText:  []string{`SELECT * FROM c WHERE c.name = 'it\'s; here'`, ` SELECT 1`},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "semicolon in a line comment",
			script:    "SELECT * FROM c -- one; two\n; SELECT 1",
			wantText:  []string{"SELECT * FROM c -- one; two", " SELECT 1"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "comment-only segment is empty",
			script:    "-- nothing here\n;SELECT 1",
			wantText:  []string{"-- nothing here", "SELECT 1"},
			wantEmpty: []bool{true, false},
		},
		{
			name:      "whitespace-only segment is empty",
			script:    "  \n\t ;SELECT 1",
			wantText:  []string{"", "SELECT 1"},
			wantEmpty: []bool{true, false},
		},
		{
			name:      "unterminated string runs to the end",
			script:    "SELECT * FROM c WHERE c.name = 'unclosed; SELECT 2",
			wantText:  []string{"SELECT * FROM c WHERE c.name = 'unclosed; SELECT 2"},
			wantEmpty: []bool{false},
		},
		{
			name:      "empty script",
			script:    "",
			wantText:  nil,
			wantEmpty: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			segments := Split(tc.script)
			if len(segments) != len(tc.wantText) {
				t.Fatalf("Split(%q) returned %d segments, want %d", tc.script, len(segments), len(tc.wantText))
			}
			for i, segment := range segments {
				if segment.Text != tc.wantText[i] {
					t.Errorf("segment %d of %q: Text = %q, want %q", i, tc.script, segment.Text, tc.wantText[i])
				}
				if got := segment.Empty(); got != tc.wantEmpty[i] {
					t.Errorf("segment %d of %q: Empty() = %v, want %v", i, tc.script, got, tc.wantEmpty[i])
				}
				if got := tc.script[segment.ByteStart:segment.ByteEnd]; got != segment.Text {
					t.Errorf("segment %d of %q: range [%d,%d] extracts %q, want Text %q",
						i, tc.script, segment.ByteStart, segment.ByteEnd, got, segment.Text)
				}
			}
		})
	}
}

// TestSplitStatementsParse pins that a script of several statements splits into
// text the parser reads one statement at a time, and that the parser itself
// refuses the whole script: the reason the driver splits before it parses.
func TestSplitStatementsParse(t *testing.T) {
	script := "SELECT c.name FROM c WHERE c.country = 'US'; SELECT TOP 2 c.seq FROM c ORDER BY c.seq"
	if _, err := Parse(script); err == nil {
		t.Fatalf("Parse(%q) accepted a script of two statements", script)
	}

	var parsed int
	for _, segment := range Split(script) {
		if segment.Empty() {
			continue
		}
		statements, err := Parse(segment.Text)
		if err != nil {
			t.Fatalf("Parse(%q): %v", segment.Text, err)
		}
		if len(statements) != 1 {
			t.Fatalf("Parse(%q) returned %d statements, want 1", segment.Text, len(statements))
		}
		parsed++
	}
	if parsed != 2 {
		t.Fatalf("parsed %d statements of %q, want 2", parsed, script)
	}
}
