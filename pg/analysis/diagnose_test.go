package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	require.Nil(t, Diagnose("SELECT 1"))
	require.Nil(t, Diagnose("SELECT 1;\n  "))
	require.Nil(t, Diagnose(""))

	err := Diagnose("SELECT * FROM")
	require.NotNil(t, err)
	require.Equal(t, "SELECT * FROM;", err.Text, "the last statement may omit its ';'")
	require.True(t, err.ParseError)
	require.Equal(t, 13, err.Offset)

	// The parser reports some errors without a position.
	err = Diagnose("SELECT 1 FETCH FIRST 1 ROWS WITH TIES")
	require.NotNil(t, err)
	require.True(t, err.ParseError)
	require.Equal(t, -1, err.Offset)
	require.Equal(t, "WITH TIES cannot be specified without ORDER BY clause", err.Message)
}
