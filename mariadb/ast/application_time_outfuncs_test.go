package ast

import (
	"strings"
	"testing"
)

// TestCreateTableAppPeriodOutfuncs locks outfuncs coverage of the application-
// time PERIOD FOR <name> fields on CreateTableStmt.
func TestCreateTableAppPeriodOutfuncs(t *testing.T) {
	n := &CreateTableStmt{
		Table:             &TableRef{Name: "t"},
		AppPeriodName:     "app_time",
		AppPeriodStartCol: "s",
		AppPeriodEndCol:   "e",
	}
	if got := NodeToString(n); !strings.Contains(got, ":period_for app_time (s, e)") {
		t.Errorf("NodeToString missing app-time period:\n%s", got)
	}
}
