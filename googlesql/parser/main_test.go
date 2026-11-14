package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/googlesql/internal/testcontainer"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.RemoveHarness()
	testcontainer.Terminate()
	os.Exit(code)
}
