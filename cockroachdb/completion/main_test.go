package completion

import (
	"os"
	"testing"
	"time"

	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

// TestMain accounts for the package's tapes and removes the CockroachDB
// server a recording started. The completion budget is lifted so that a
// loaded machine cannot cut a test's answer short; budget tests set it.
func TestMain(m *testing.M) {
	completionBudget = time.Minute
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
