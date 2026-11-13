package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/trino/internal/testcontainer"
)

// TestMain removes the Trino server a test started.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}
