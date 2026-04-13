package partiql

import (
	"testing"
)

func TestSplit(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		wantN     int
		wantText  []string
		wantEmpty []bool
	}{
		{
			name:      "single_statement",
			input:     "SELECT * FROM t",
			wantN:     1,
			wantText:  []string{"SELECT * FROM t"},
			wantEmpty: []bool{false},
		},
		{
			name:      "two_statements",
			input:     "SELECT 1; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT 1;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "trailing_semicolon",
			input:     "SELECT 1;",
			wantN:     1,
			wantText:  []string{"SELECT 1;"},
			wantEmpty: []bool{false},
		},
		{
			name:      "multiple_trailing_semicolons",
			input:     "SELECT 1;;",
			wantN:     2,
			wantText:  []string{"SELECT 1;", ";"},
			wantEmpty: []bool{false, true},
		},
		{
			name:      "semicolon_in_string",
			input:     "SELECT 'a;b'",
			wantN:     1,
			wantText:  []string{"SELECT 'a;b'"},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_quoted_ident",
			input:     `SELECT "a;b"`,
			wantN:     1,
			wantText:  []string{`SELECT "a;b"`},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_ion_literal",
			input:     "SELECT `{a;b}`",
			wantN:     1,
			wantText:  []string{"SELECT `{a;b}`"},
			wantEmpty: []bool{false},
		},
		{
			name:      "semicolon_in_line_comment",
			input:     "SELECT 1 -- comment; not a split\n; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT 1 -- comment; not a split\n;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "semicolon_in_block_comment",
			input:     "SELECT /* ; */ 1; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT /* ; */ 1;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:  "empty_input",
			input: "",
			wantN: 0,
		},
		{
			name:      "whitespace_only",
			input:     "   \n\t  ",
			wantN:     1,
			wantText:  []string{"   \n\t  "},
			wantEmpty: []bool{true},
		},
		{
			name:      "comment_only",
			input:     "-- just a comment",
			wantN:     1,
			wantText:  []string{"-- just a comment"},
			wantEmpty: []bool{true},
		},
		{
			name:      "doubled_quote_escape",
			input:     "SELECT 'it''s'; SELECT 1",
			wantN:     2,
			wantText:  []string{"SELECT 'it''s';", " SELECT 1"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "nested_block_comment",
			input:     "SELECT /* outer /* inner */ still comment */ 1",
			wantN:     1,
			wantText:  []string{"SELECT /* outer /* inner */ still comment */ 1"},
			wantEmpty: []bool{false},
		},
		{
			name:      "ion_with_backtick_in_short_string",
			input:     "SELECT `\"a`b\"`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `\"a`b\"`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_backtick_in_comment",
			input:     "SELECT `/* not end ` still ion */`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `/* not end ` still ion */`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_backtick_in_symbol",
			input:     "SELECT `'a`b'`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `'a`b'`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_with_triple_quoted_long_string",
			input:     "SELECT `'''contains ` backtick'''`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `'''contains ` backtick'''`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
		{
			name:      "ion_simple_no_inner_strings",
			input:     "SELECT `{a: 1}`; SELECT 2",
			wantN:     2,
			wantText:  []string{"SELECT `{a: 1}`;", " SELECT 2"},
			wantEmpty: []bool{false, false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			segs := Split(tc.input)
			if len(segs) != tc.wantN {
				t.Fatalf("Split(%q) returned %d segments, want %d", tc.input, len(segs), tc.wantN)
			}
			for i, seg := range segs {
				if i < len(tc.wantText) && seg.Text != tc.wantText[i] {
					t.Errorf("seg[%d].Text = %q, want %q", i, seg.Text, tc.wantText[i])
				}
				if i < len(tc.wantEmpty) && seg.Empty() != tc.wantEmpty[i] {
					t.Errorf("seg[%d].Empty() = %v, want %v", i, seg.Empty(), tc.wantEmpty[i])
				}
			}
		})
	}
}
