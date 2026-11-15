package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
	"github.com/bytebase/omni/oceanbase/internal/testcontainer"
)

func TestMain(m *testing.M) {
	code := tape.Main(m)
	testcontainer.OceanBase.Terminate()
	os.Exit(code)
}
