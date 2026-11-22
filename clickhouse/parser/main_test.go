package parser

import (
	"os"
	"testing"

	"github.com/bytebase/omni/clickhouse/internal/testcontainer"
)

// TestMain removes any ClickHouse container started during a recording or
// verification run.
func TestMain(m *testing.M) {
	code := m.Run()
	testcontainer.Terminate()
	os.Exit(code)
}
