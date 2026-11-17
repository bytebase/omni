package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain accounts for the package's tapes. A replay starts nothing; a
// recording reaches the cluster REDSHIFT_DSN names.
func TestMain(m *testing.M) {
	os.Exit(tape.Main(m))
}
