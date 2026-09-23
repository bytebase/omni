package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/googlesql/internal/spannertest"
)

func TestMain(m *testing.M) {
	code := m.Run()
	spannertest.RemoveHarness()
	os.Exit(code)
}
