package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestMain accounts for the package's tapes and removes the Trino server a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
