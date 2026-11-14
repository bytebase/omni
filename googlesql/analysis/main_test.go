package analysis

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/googlesql/internal/testcontainer"
)

// TestMain runs the tests through the tape harness and removes the Spanner
// emulator a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
