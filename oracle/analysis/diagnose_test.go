package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	require.Nil(t, Diagnose(""))
	require.Nil(t, Diagnose("SELECT 1 FROM DUAL"))
	require.Nil(t, Diagnose("SELECT 1 FROM DUAL;\nSELECT 2 FROM DUAL  \n"))

	got := Diagnose("SELECT 1 FROM DUAL;\nSELECT * FROM")
	require.NotNil(t, got)
	require.Equal(t, "SELECT 1 FROM DUAL;\nSELECT * FROM;", got.Text)
	require.GreaterOrEqual(t, got.Offset, len("SELECT 1 FROM DUAL;\n"))
	require.NotEmpty(t, got.Message)
}
