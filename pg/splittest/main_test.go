package splittest

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestMain accounts for the package's tapes and removes the PostgreSQL a
// recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
