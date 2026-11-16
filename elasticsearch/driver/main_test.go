package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/elasticsearch/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain accounts for the package's tapes.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
