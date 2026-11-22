package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// TestMain removes the Spanner emulator container and harness binary a
// recording started, if any.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}
