package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/starrocks/internal/testcontainer"
)

// TestMain removes the server a recording started.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}
