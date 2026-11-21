package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain accounts for the package's tapes: a replay fails on unused tapes or
// streams, and a recording removes stale ones.
func TestMain(m *testing.M) {
	os.Exit(tape.Main(m))
}
