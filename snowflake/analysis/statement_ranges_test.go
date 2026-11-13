package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestStatementRanges(t *testing.T) {
	tests := []struct {
		sql  string
		want []analysis.ByteRange
	}{
		{"", nil},
		{";", nil},
		{"  -- note\n ", []analysis.ByteRange{{Start: 2, End: 9}}},
		{"SELECT 1; SELECT 2;", []analysis.ByteRange{{Start: 0, End: 9}, {Start: 10, End: 19}}},
		// A leading comment stays in the statement's range.
		{"\n/* c */ SELECT 1", []analysis.ByteRange{{Start: 1, End: 17}}},
		// A comment after the last statement gets a range of its own.
		{"SELECT 1; -- tail\n\n", []analysis.ByteRange{{Start: 0, End: 9}, {Start: 10, End: 17}}},
	}
	for _, tc := range tests {
		require.Equalf(t, tc.want, analysis.StatementRanges(tc.sql), "%q", tc.sql)
	}
}
