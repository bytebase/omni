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
