package completion

import (
	"testing"

	"github.com/bytebase/omni/metadata"
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

func TestCompleteWithMetadataColumns(t *testing.T) {
	meta := &metadata.DatabaseSchemaMetadata{
		Name: "xx-ap-east-1",
		Schemas: []*metadata.SchemaMetadata{
			{
				Name: "",
				Tables: []*metadata.TableMetadata{
					{
						Name: "Music",
						Columns: []*metadata.ColumnMetadata{
							{Name: "Artist"},
							{Name: "SongTitle"},
						},
					},
					{
						Name: "Album",
						Columns: []*metadata.ColumnMetadata{
							{Name: "AlbumId"},
							{Name: "Title"},
						},
					},
					{
						Name: "order",
						Columns: []*metadata.ColumnMetadata{
							{Name: "OrderId"},
						},
					},
				},
			},
		},
	}

	cat := catalog.New()
	cat.LoadMetadata(meta)

	tests := []struct {
		name     string
		input    string
		cat      *catalog.Catalog
		wantCols []string
		wantTabs []string
	}{
		{
			name:     "select item from table",
			input:    "SELECT | FROM Music",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "select item with prefix",
			input:    "SELECT Ar| FROM Music",
			cat:      cat,
			wantCols: []string{"Artist"},
		},
		{
			name:     "where clause columns",
			input:    "SELECT * FROM Music WHERE |",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "having clause columns",
			input:    "SELECT Artist FROM Music GROUP BY Artist HAVING |",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "order by clause columns",
			input:    "SELECT * FROM Music ORDER BY |",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "update set columns",
			input:    "UPDATE Music SET |",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "comma separated tables",
			input:    "SELECT | FROM Music,Album",
			cat:      cat,
			wantCols: []string{"AlbumId", "Artist", "SongTitle", "Title"},
		},
		{
			name:     "aliased table with comment",
			input:    "SELECT | FROM Music AS m /* music */",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "quoted keyword table name",
			input:    `SELECT | FROM "order"`,
			cat:      cat,
			wantCols: []string{"OrderId"},
		},
		{
			name:     "multi statement semicolon isolation",
			input:    "SELECT * FROM Album; SELECT | FROM Music",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "semicolon inside string literal",
			input:    "SELECT | FROM Music WHERE Artist = 'a;b'",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "no from clause falls back to all tables",
			input:    "SELECT |",
			cat:      cat,
			wantCols: []string{"AlbumId", "Artist", "OrderId", "SongTitle", "Title"},
		},
		{
			name:     "from clause offers tables not columns",
			input:    "SELECT * FROM |",
			cat:      cat,
			wantCols: nil,
			wantTabs: []string{"Album", "Music", "order"},
		},
		{
			name:     "nil catalog derives tables and columns from metadata",
			input:    "SELECT * FROM |",
			cat:      nil,
			wantCols: nil,
			wantTabs: []string{"Album", "Music", "order"},
		},
		{
			name:     "nil catalog resolves columns from metadata",
			input:    "SELECT | FROM Music",
			cat:      nil,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "qualified alias restricts columns to that table",
			input:    "SELECT m.| FROM Music m, Album a",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "nested subquery does not leak table into outer query",
			input:    "SELECT | FROM Album WHERE EXISTS (SELECT 1 FROM Music)",
			cat:      cat,
			wantCols: []string{"AlbumId", "Title"},
		},
		{
			name:     "string literal FROM does not suppress column completion",
			input:    "SELECT * FROM Music WHERE Artist = 'FROM' |",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "cursor immediately after semicolon starts new statement",
			input:    "SELECT * FROM Music WHERE Artist;|",
			cat:      cat,
			wantCols: nil,
		},
		{
			name:     "expression parentheses keep outer FROM scope",
			input:    "SELECT (|) FROM Music",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "string literal FROM Album does not leak Album columns",
			input:    "SELECT | FROM Music WHERE Artist = 'x FROM Album WHERE y'",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "correlated subquery resolves outer qualified alias",
			input:    "SELECT * FROM Album a WHERE EXISTS (SELECT 1 FROM Music m WHERE a.|)",
			cat:      cat,
			wantCols: []string{"AlbumId", "Title"},
		},
		{
			name:     "block comment containing subquery does not hijack scope",
			input:    "SELECT | FROM Music WHERE Artist = 'a' /* (SELECT 1 FROM Album */",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "line comment containing subquery does not hijack scope",
			input:    "SELECT | FROM Music -- (SELECT 1 FROM Album\nWHERE Artist = 'a'",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "parenthesized table in FROM resolves table columns",
			input:    "SELECT m.| FROM (Music) m",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "parenthesized subquery in FROM projects only selected columns",
			input:    "SELECT m.| FROM (SELECT Artist AS performer FROM Music) m",
			cat:      cat,
			wantCols: []string{"performer"},
		},
		{
			name:     "parenthesized star subquery in FROM expands inner table columns",
			input:    "SELECT m.| FROM (SELECT * FROM Music) m",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
		},
		{
			name:     "expression FROM inside TRIM does not suggest tables or suppress columns",
			input:    "SELECT TRIM(' ' FROM |) FROM Music",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
			wantTabs: []string{},
		},
		{
			name:     "expression FROM inside SUBSTRING does not suppress later select-list columns",
			input:    "SELECT SUBSTRING(Artist FROM 1 FOR 3), | FROM Music",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
			wantTabs: []string{},
		},
		{
			name:     "expression FROM inside WHERE expression keeps column context",
			input:    "SELECT * FROM Music WHERE SUBSTRING(Artist FROM |)",
			cat:      cat,
			wantCols: []string{"Artist", "SongTitle"},
			wantTabs: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sql, pos := splitMarker(t, tc.input)
			got := CompleteWithMetadata(sql, pos, tc.cat, meta)

			var gotCols, gotTabs []string
			for _, c := range got {
				switch c.Kind {
				case "column":
					gotCols = append(gotCols, c.Text)
				case "table":
					gotTabs = append(gotTabs, c.Text)
				}
			}
			if !equalStringSlices(gotCols, tc.wantCols) {
				t.Errorf("columns = %v, want %v", gotCols, tc.wantCols)
			}
			if tc.wantTabs != nil && !equalStringSlices(gotTabs, tc.wantTabs) {
				t.Errorf("tables = %v, want %v", gotTabs, tc.wantTabs)
			}
		})
	}
}

func TestCompleteNilCatalogDoesNotPanic(t *testing.T) {
	got := Complete("SELECT * FROM ", 14, nil)
	if countKind(got, "table") != 0 {
		t.Errorf("expected no tables for nil catalog, got %v", got)
	}
	if countKind(got, "keyword") == 0 {
		t.Errorf("expected keyword candidates for nil catalog")
	}
}

func splitMarker(t *testing.T, input string) (string, int) {
	t.Helper()
	for i := 0; i < len(input); i++ {
		if input[i] == '|' {
			return input[:i] + input[i+1:], i
		}
	}
	t.Fatalf("missing '|' marker in %q", input)
	return "", 0
}

func equalStringSlices(a, b []string) bool {
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
