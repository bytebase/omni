package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
)

// TestMain removes any CockroachDB container started during a recording or
// verification run.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}
