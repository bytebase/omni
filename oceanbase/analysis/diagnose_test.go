package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	for _, s := range []string{"", "  \n", "SELECT 1", "SELECT 1;", "SELECT 1;\nSELECT 2", "SELECT 1; -- c"} {
		require.Nilf(t, Diagnose(s), "%q", s)
	}

	got := Diagnose("SELECT 1;\nSELEC 2")
	require.NotNil(t, got)
	require.Equal(t, "SELECT 1;\nSELEC 2;", got.Text)
	require.Equal(t, 10, got.Offset)
	require.NotEmpty(t, got.Message)

	// The offset counts bytes, so multibyte text before the error moves it.
	got = Diagnose("SELECT 'é😀';\nSELEC 2;")
	require.NotNil(t, got)
	require.Equal(t, "SELECT 'é😀';\nSELEC 2;", got.Text)
	require.Equal(t, 17, got.Offset)
}
