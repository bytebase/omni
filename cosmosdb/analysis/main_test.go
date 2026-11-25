package analysis_test

import (
	"os"
	"testing"

	"github.com/bytebase/omni/cosmosdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain runs the tests through the tape harness and removes the emulator
// a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
