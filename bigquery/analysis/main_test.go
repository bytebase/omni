package analysis

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain runs the tests through the tape harness.
func TestMain(m *testing.M) {
	os.Exit(tape.Main(m))
}
