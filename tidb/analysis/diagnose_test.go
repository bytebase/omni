package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiagnose(t *testing.T) {
	for _, s := range []string{"", "  \n", "SELECT 1", "SELECT 1;", "SELECT 1;\nSELECT 2", "SELECT 1; -- c"} {
		got, err := Diagnose(s)
		require.NoErrorf(t, err, "%q", s)
		require.Nilf(t, got, "%q", s)
	}

	script := "SELECT 1;\nSELECT 'é😀' FROM t WHERE;\nSELECT 3;"
	got, err := Diagnose(script)
	require.NoError(t, err)
	require.NotNil(t, got)
	// The range is the failing statement's, without its ";".
	require.Equal(t, "\nSELECT 'é😀' FROM t WHERE", script[got.ByteStart:got.ByteEnd])
	require.Equal(t, ';', rune(script[got.ByteEnd]))
	// The offset is in the statement's text and counts bytes.
	require.Equal(t, len("\nSELECT 'é😀' FROM t WHERE"), got.Offset)
	require.NotEmpty(t, got.Message)
}
