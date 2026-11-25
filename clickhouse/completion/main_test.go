package completion

import (
	"os"
	"testing"
	"time"
)

// TestMain lifts the completion budget so that a loaded machine cannot cut
// a test's answer short; budget tests set it.
func TestMain(m *testing.M) {
	completionBudget = time.Minute
	os.Exit(m.Run())
}
