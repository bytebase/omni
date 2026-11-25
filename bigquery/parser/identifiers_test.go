package parser

import "testing"

// TestNormalizeGoogleSQLIdentifier locks the byte-for-byte parity with the
// legacy bytebase unquoteIdentifierByText helper that the bigquery/spanner
// query-span extractors use. The legacy guard is:
//
//	len >= 3 && hasPrefix "`" && hasSuffix "`"  =>  strip the surrounding backticks
//	otherwise                                   =>  return unchanged
func TestNormalizeGoogleSQLIdentifier(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"bare", "users", "users"},
		{"bare_mixed_case_preserved", "MyTable", "MyTable"}, // legacy did NOT lower-case here
		{"backtick_quoted", "`my table`", "my table"},
		{"backtick_with_dot", "`project.dataset.table`", "project.dataset.table"},
		{"backtick_min_len3", "`a`", "a"},
		// len < 3 guard: a lone or doubled backtick is NOT stripped.
		{"single_backtick", "`", "`"},
		{"empty_backticks_len2", "``", "``"},
		{"empty", "", ""},
		// Only SURROUNDING backticks are stripped; an interior backtick stays.
		{"interior_backtick", "a`b", "a`b"},
		// No lower-casing, no escape collapsing (legacy sliced verbatim).
		{"backtick_with_backslash", "`a\\`b`", "a\\`b"},
		{"prefix_only", "`abc", "`abc"},
		{"suffix_only", "abc`", "abc`"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeGoogleSQLIdentifier(c.in); got != c.want {
				t.Errorf("NormalizeGoogleSQLIdentifier(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestIsIdentifierStart(t *testing.T) {
	// A plain identifier token (also covers backtick idents, which the lexer
	// collapses to tokIdentifier).
	if !isIdentifierStart(tokIdentifier) {
		t.Error("isIdentifierStart(tokIdentifier) = false, want true")
	}
	// A non-reserved keyword may stand in as an identifier.
	if !isIdentifierStart(kwOPTIONS) { // OPTIONS is non-reserved
		t.Error("isIdentifierStart(kwOPTIONS) = false, want true (non-reserved keyword)")
	}
	if keywordReserved[kwOPTIONS] {
		t.Fatal("test premise wrong: kwOPTIONS should be non-reserved")
	}
	// A reserved keyword may NOT be a bare identifier.
	if isIdentifierStart(kwSELECT) {
		t.Error("isIdentifierStart(kwSELECT) = true, want false (reserved)")
	}
	if !keywordReserved[kwSELECT] {
		t.Fatal("test premise wrong: kwSELECT should be reserved")
	}
	// Non-identifier tokens.
	for _, tk := range []int{tokEOF, tokInteger, tokString, int('('), int('@')} {
		if isIdentifierStart(tk) {
			t.Errorf("isIdentifierStart(%s) = true, want false", TokenName(tk))
		}
	}
}

func TestIsAnyKeywordIdentifier(t *testing.T) {
	// The permissive variant accepts ANY keyword, reserved or not, plus
	// tokIdentifier.
	if !isAnyKeywordIdentifier(tokIdentifier) {
		t.Error("isAnyKeywordIdentifier(tokIdentifier) = false, want true")
	}
	if !isAnyKeywordIdentifier(kwSELECT) { // reserved
		t.Error("isAnyKeywordIdentifier(kwSELECT) = false, want true (any keyword allowed)")
	}
	if !isAnyKeywordIdentifier(kwOPTIONS) { // non-reserved
		t.Error("isAnyKeywordIdentifier(kwOPTIONS) = false, want true")
	}
	// Non-keyword, non-identifier tokens are still rejected.
	for _, tk := range []int{tokEOF, tokInteger, int('.'), int('@')} {
		if isAnyKeywordIdentifier(tk) {
			t.Errorf("isAnyKeywordIdentifier(%s) = true, want false", TokenName(tk))
		}
	}
}
