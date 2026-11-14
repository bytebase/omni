package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatementRanges(t *testing.T) {
	r := func(sl, sc, el, ec uint32) Range {
		return Range{Start: Position{Line: sl, Character: sc}, End: Position{Line: el, Character: ec}}
	}
	tests := []struct {
		statement string
		want      []Range
	}{
		{"", nil},
		{"  \n ", nil},
		{"SELECT 1", []Range{r(0, 0, 1, 0)}},
		// A range that ends at a ';' on its line ends one character past it.
		{"SELECT 1;", []Range{r(0, 0, 0, 10)}},
		{"SELECT 1; SELECT 2;", []Range{r(0, 0, 0, 10), r(0, 10, 0, 20)}},
		{"SELECT 1;\nSELECT 2", []Range{r(0, 0, 1, 0), r(1, 0, 2, 0)}},
		{"  SELECT 1;\n\n  SELECT 2;  ", []Range{r(0, 2, 1, 0), r(2, 2, 2, 12)}},
		// UTF-16: an astral character counts two.
		{"SELECT '😀'; SELECT 2", []Range{r(0, 0, 0, 13), r(0, 13, 1, 0)}},
		// Every ';' splits, even one inside a string.
		{"SELECT ';'", []Range{r(0, 0, 0, 10), r(0, 9, 1, 0)}},
	}
	for _, tc := range tests {
		require.Equal(t, tc.want, StatementRanges(tc.statement), "%q", tc.statement)
	}
}
