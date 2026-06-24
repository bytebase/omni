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

// TestAlterTableAppPeriodOutfuncs locks outfuncs coverage of the period name on
// ADD/DROP PERIOD AlterTableCmd nodes, so application-time and system-time
// periods are distinguishable in NodeToString.
func TestAlterTableAppPeriodOutfuncs(t *testing.T) {
	add := &AlterTableCmd{Type: ATAddPeriod, PeriodName: "app_time", PeriodStartCol: "s", PeriodEndCol: "e"}
	if got := NodeToString(add); !strings.Contains(got, ":period_for app_time (s, e)") {
		t.Errorf("ADD PERIOD NodeToString missing period name:\n%s", got)
	}
	drop := &AlterTableCmd{Type: ATDropPeriod, PeriodName: "app_time"}
	if got := NodeToString(drop); !strings.Contains(got, ":period_for app_time") {
		t.Errorf("DROP PERIOD NodeToString missing period name:\n%s", got)
	}
}
