package analysis_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bytebase/omni/snowflake/analysis"
)

func TestSplitStatements(t *testing.T) {
	require.Nil(t, analysis.SplitStatements(""))
	require.Nil(t, analysis.SplitStatements(";;"))
	require.Equal(t, []analysis.Statement{
		{Text: "SELECT 1;", Start: 0, End: 9},
		{Text: "  SELECT 'a;b'", Start: 10, End: 24},
	}, analysis.SplitStatements("SELECT 1;;  SELECT 'a;b'"))
}
