package analysis

import (
	"os"
	"testing"

	spannertest "github.com/bytebase/omni/googlesql/internal/testcontainer"
)

func TestMain(m *testing.M) {
	code := m.Run()
	spannertest.RemoveHarness()
	os.Exit(code)
}
