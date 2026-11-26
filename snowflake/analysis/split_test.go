package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestSplitStatements(t *testing.T) {
	require.Nil(t, analysis.SplitStatements(""))
	require.Equal(t, []analysis.Statement{
		{Text: ";", Empty: true, Start: 0, End: 1},
		{Text: ";", Empty: true, Start: 1, End: 2},
	}, analysis.SplitStatements(";;"))
	require.Equal(t, []analysis.Statement{
		{Text: "SELECT 1;", Start: 0, End: 9},
		{Text: ";", Empty: true, Start: 9, End: 10},
		{Text: "  SELECT 'a;b'", Start: 10, End: 24},
	}, analysis.SplitStatements("SELECT 1;;  SELECT 'a;b'"))
}
