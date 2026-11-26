package analysis

import (
	"os"
	"testing"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain removes the server a recording started.
func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
