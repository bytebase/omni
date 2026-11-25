package parser

import (
	"slices"
	"testing"
)

func TestExpectedKeywords(t *testing.T) {
	tests := []struct {
		prefix string
		want   []string // the whole list, when set
		has    []string
		not    []string
	}{
		{prefix: "INSERT ", want: []string{"INTO", "OVERWRITE"}},
		{prefix: "SELECT a FROM t GROUP ", want: []string{"BY"}},
		{prefix: "SELECT * FROM t LEFT ", want: []string{"ANTI", "JOIN", "OUTER", "SEMI"}},
		{prefix: "UPDATE t ", want: []string{"AS", "SET"}},
		// A keyword the parser takes only by the token after it.
		{prefix: "SELECT * FROM t ", has: []string{"CROSS", "FULL", "INNER", "JOIN", "LEFT", "RIGHT", "WHERE"}},
		// IF is reserved; the non-reserved keywords there are table names.
		// CREATE TABLE FUNCTION is a function (function.go).
		{prefix: "CREATE TABLE ", want: []string{"FUNCTION", "IF"}},
		// Non-reserved keywords read as clauses rather than as an alias.
		{prefix: "SELECT * FROM t ", has: []string{"TABLET"}, not: []string{"COMMENT", "USER"}},
		{prefix: "SELECT * FROM t WHERE a ", has: []string{"AND", "BETWEEN", "IN", "IS", "LIKE", "NOT"}, not: []string{"FROM", "SELECT"}},
		{prefix: "", has: []string{"ALTER", "CREATE", "SELECT", "SHOW"}, not: []string{"FROM", "JOIN"}},
		// The prefix itself fails: nothing fits after it.
		{prefix: "SELECT FROM FROM ", want: nil},
		// No separator is added after punctuation-free text ending in a word.
		{prefix: "INSERT", want: []string{"INTO", "OVERWRITE"}},
	}
	for _, tc := range tests {
		got := ExpectedKeywords(tc.prefix)
		if tc.has == nil && tc.not == nil {
			if !slices.Equal(got, tc.want) {
				t.Errorf("%q: got %v, want %v", tc.prefix, got, tc.want)
			}
			continue
		}
		for _, w := range tc.has {
			if !slices.Contains(got, w) {
				t.Errorf("%q: missing %s in %v", tc.prefix, w, got)
			}
		}
		for _, w := range tc.not {
			if slices.Contains(got, w) {
				t.Errorf("%q: unexpected %s in %v", tc.prefix, w, got)
			}
		}
	}
}
