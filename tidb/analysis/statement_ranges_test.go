package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStatementRanges(t *testing.T) {
	tests := []struct {
		sql  string
		want []ByteRange
	}{
		{"", nil},
		{"   ", nil},
		{"-- just a comment\n", nil},
		{"/* block */", nil},
		{"SELECT 1; SELECT 2;", []ByteRange{{Start: 0, End: 9}, {Start: 10, End: 19}}},
		// The last statement has no ';' to include.
		{"SELECT 1;\n  SELECT 2", []ByteRange{{Start: 0, End: 9}, {Start: 12, End: 20}}},
		// Unicode whitespace is skipped too: U+3000 is three bytes.
		{"SELECT 1;　SELECT 2;", []ByteRange{{Start: 0, End: 9}, {Start: 12, End: 21}}},
		// A syntax error does not drop the ranges around it.
		{"SELECT 1; SELEC 2; SELECT 3;", []ByteRange{{Start: 0, End: 9}, {Start: 10, End: 18}, {Start: 19, End: 28}}},
	}
	for _, tc := range tests {
		require.Equalf(t, tc.want, StatementRanges(tc.sql), "%q", tc.sql)
	}
}

func TestStatementRangesCustomDelimiter(t *testing.T) {
	sql := "DELIMITER $$\nSELECT 1$$\nDELIMITER ;\nSELECT 2;"
	got := StatementRanges(sql)
	var texts []string
	for _, r := range got {
		texts = append(texts, sql[r.Start:r.End])
	}
	require.Equal(t, []string{"SELECT 1", "SELECT 2;"}, texts)
}
