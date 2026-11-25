package completion

import (
	"os"
	"testing"

	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain accounts for the package's tapes and removes the CockroachDB
// server a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
