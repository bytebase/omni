package completion

import (
	"testing"

	"github.com/bytebase/omni/partiql/catalog"
)

func TestComplete(t *testing.T) {
	cat := catalog.New()
	cat.AddTable("Music")
	cat.AddTable("Albums")
	cat.AddTable("Artists")

	tests := []struct {
		name       string
		input      string
		pos        int
		wantTables []string // expected table candidates
		wantMinKW  int      // minimum keyword candidates expected
	}{
		{
			name:       "after_from_empty",
			input:      "SELECT * FROM ",
			pos:        14,
			wantTables: []string{"Albums", "Artists", "Music"},
		},
		{
			name:       "after_from_prefix",
			input:      "SELECT * FROM M",
			pos:        15,
			wantTables: []string{"Music"},
		},
		{
			name:       "after_select_keyword",
			input:      "SEL",
			pos:        3,
			wantTables: nil,
			wantMinKW:  1, // at least SELECT
		},
		{
			name:      "empty_input",
			input:     "",
			pos:       0,
			wantMinKW: 10, // many keywords
		},
		{
			name:       "after_join",
			input:      "SELECT * FROM Music JOIN ",
			pos:        25,
			wantTables: []string{"Albums", "Artists", "Music"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			candidates := Complete(tc.input, tc.pos, cat)

			// Check table candidates
			if tc.wantTables != nil {
				var gotTables []string
				for _, c := range candidates {
					if c.Kind == "table" {
						gotTables = append(gotTables, c.Text)
					}
				}
				if len(gotTables) != len(tc.wantTables) {
					t.Errorf("table candidates = %v, want %v", gotTables, tc.wantTables)
				}
			}

			// Check minimum keyword count
			if tc.wantMinKW > 0 {
				kwCount := 0
				for _, c := range candidates {
					if c.Kind == "keyword" {
						kwCount++
					}
				}
				if kwCount < tc.wantMinKW {
					t.Errorf("keyword count = %d, want >= %d", kwCount, tc.wantMinKW)
				}
			}
		})
	}
}

// TestCompletePositionOutOfRange is a negative test guarding against a slice
// panic when the caller passes a cursor position past the end of the input.
// Complete() is exported, so it must clamp internally rather than rely on a
// caller-side clamp. extractPrefix() already clamps pos, but Complete() used
// the original unclamped pos when slicing input[:pos-len(prefix)], which
// panicked with "slice bounds out of range" for any pos > len(input).
//
// Spec note (oracle = generated ANTLR PartiQL lexer, truth2): tokenizing
// "SELECT * FROM " yields "... FROM WS EOF" — no partial identifier follows
// FROM — so a cursor at/after the end is equivalent to a cursor at the end:
// prefix is empty and the context is FROM, so tables are suggested.
func TestCompletePositionOutOfRange(t *testing.T) {
	cat := catalog.New()
	cat.AddTable("Music")
	cat.AddTable("Albums")
	cat.AddTable("Artists")

	cases := []struct {
		name  string
		input string
		pos   int
	}{
		{name: "from_context_far_past_end", input: "SELECT * FROM ", pos: 100},
		{name: "from_context_one_past_end", input: "SELECT * FROM ", pos: len("SELECT * FROM ") + 1},
		{name: "with_prefix_past_end", input: "SELECT * FROM M", pos: 100},
		{name: "empty_input_past_end", input: "", pos: 5},
		{name: "huge_positive_offset", input: "SELECT", pos: 1 << 30},
		{name: "negative_offset", input: "SELECT", pos: -3},
		{name: "negative_offset_empty_input", input: "", pos: -1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Complete(%q, %d) panicked: %v", tc.input, tc.pos, r)
				}
			}()
			_ = Complete(tc.input, tc.pos, cat)
		})
	}

	// Beyond "does not panic": a cursor past the end of "SELECT * FROM "
	// must behave exactly like a cursor at the end — FROM context, all tables
	// suggested, empty prefix.
	atEnd := Complete("SELECT * FROM ", len("SELECT * FROM "), cat)
	pastEnd := Complete("SELECT * FROM ", 100, cat)
	if !equalCandidates(atEnd, pastEnd) {
		t.Fatalf("past-end completion differs from at-end completion:\n at-end = %v\n past-end = %v", atEnd, pastEnd)
	}
	if countKind(pastEnd, "table") != 3 {
		t.Errorf("past-end FROM context: table candidates = %d, want 3", countKind(pastEnd, "table"))
	}
}

// TestKeywordsNoDuplicates guards against duplicate entries in the keyword
// candidate list. "NOT" was listed twice; it is a single PartiQL keyword
// (oracle: PartiQLLexer.g4 -> `NOT: 'NOT';`, truth2), so it must appear in
// completion output exactly once.
func TestKeywordsNoDuplicates(t *testing.T) {
	seen := map[string]int{}
	for _, kw := range keywords {
		seen[kw]++
	}
	for kw, n := range seen {
		if n > 1 {
			t.Errorf("keyword %q appears %d times in keywords list, want exactly 1", kw, n)
		}
	}
	if seen["NOT"] != 1 {
		t.Errorf("keyword %q appears %d times, want 1", "NOT", seen["NOT"])
	}

	// And the dedup must surface through Complete(): "NOT" must be emitted
	// once for an input where it matches as a prefix.
	cat := catalog.New()
	got := Complete("NO", len("NO"), cat)
	notCount := 0
	for _, c := range got {
		if c.Kind == "keyword" && c.Text == "NOT" {
			notCount++
		}
	}
	if notCount != 1 {
		t.Errorf("Complete(\"NO\") emitted NOT %d times, want 1", notCount)
	}
}

func equalCandidates(a, b []Candidate) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func countKind(cs []Candidate, kind string) int {
	n := 0
	for _, c := range cs {
		if c.Kind == kind {
			n++
		}
	}
	return n
}
