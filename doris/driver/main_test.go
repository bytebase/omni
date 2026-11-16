package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/doris/internal/testcontainer"
	"github.com/bytebase/omni/enginetest/tape"
)

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.Terminate()
	os.Exit(code)
}
