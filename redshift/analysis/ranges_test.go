package analysis

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStatementRanges(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		want      []Range
	}{
		{
			name:      "multi statement",
			statement: "SELECT 1;\n  SELECT 2;",
			want: []Range{
				{Start: 0, End: 9},
				{Start: 12, End: 21},
			},
		},
		{
			name:      "comment prefix",
			statement: "/*🙂*/ SELECT 1;",
			want: []Range{
				{Start: 9, End: 18},
			},
		},
		{
			name:      "trailing semicolon is included",
			statement: "SELECT 1;SELECT 2",
			want: []Range{
				{Start: 0, End: 9},
				{Start: 9, End: 17},
			},
		},
		{
			name:      "comment between statements",
			statement: "SELECT 1;\n-- note\nSELECT 2;",
			want: []Range{
				{Start: 0, End: 9},
				{Start: 18, End: 27},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, StatementRanges(tc.statement))
		})
	}
}

func TestStatementRangesToleratesIncompleteSQL(t *testing.T) {
	require.Equal(t, []Range{{Start: 0, End: 13}}, StatementRanges("SELECT * FROM"))
	// A script the parser rejects keeps the splitter's boundaries, so an
	// unterminated comment stays inside its range.
	require.Equal(t, []Range{{Start: 0, End: 15}}, StatementRanges("/* unterminated"))
	require.Equal(t, []Range{{Start: 0, End: 18}}, StatementRanges("   /* unterminated"))
}

func TestStatementRangesSkipsEmptySegments(t *testing.T) {
	require.Empty(t, StatementRanges(""))
	require.Empty(t, StatementRanges("   \n\t"))
	require.Empty(t, StatementRanges("-- only a comment"))
	require.Equal(t, []Range{{Start: 0, End: 9}}, StatementRanges("SELECT 1;;"))
}

func TestStatementRangesScalesLinearlyForLargeScript(t *testing.T) {
	smallDuration := timeStatementRanges(t, statementRangeScript(1000))
	largeDuration := timeStatementRanges(t, statementRangeScript(10000))

	require.Lessf(t, largeDuration, smallDuration*40, "10x more statements should not take %s vs %s", largeDuration, smallDuration)
}

func statementRangeScript(count int) string {
	var builder strings.Builder
	for i := range count {
		if i > 0 {
			builder.WriteByte('\n')
		}
		table := fmt.Sprintf("perf_redshift_%d", i)
		switch i % 4 {
		case 0:
			fmt.Fprintf(&builder, "CREATE TABLE %s (id INT, payload VARCHAR(255))", table)
		case 1:
			fmt.Fprintf(&builder, "INSERT INTO %s (id, payload) VALUES (%d, 'payload-%d')", table, i, i)
		case 2:
			fmt.Fprintf(&builder, "UPDATE %s SET payload = 'updated-%d' WHERE id = %d", table, i, i)
		default:
			fmt.Fprintf(&builder, "DROP TABLE IF EXISTS %s", table)
		}
		builder.WriteByte(';')
	}
	return builder.String()
}

func timeStatementRanges(t *testing.T, statement string) time.Duration {
	t.Helper()
	start := time.Now()
	ranges := StatementRanges(statement)
	require.NotEmpty(t, ranges)
	return time.Since(start)
}
