package driver

import (
	"os"
	"testing"

	"github.com/bytebase/omni/enginetest/tape"
)

func TestMain(m *testing.M) {
	os.Exit(tape.Main(m))
}
