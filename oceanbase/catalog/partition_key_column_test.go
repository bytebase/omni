package catalog

import "testing"

// partitionKeyColumn tells a partitioning key that is one column, bare or in backticks,
// from an expression (TestOracle_OceanBaseRefusals holds the rule it serves to the engine).
func TestPartitionKeyColumn(t *testing.T) {
	for expr, want := range map[string]string{
		"val":           "val",
		" val ":         "val",
		"`val`":         "val",
		"`val-id`":      "val-id",
		"`part key`":    "part key",
		"`a``b`":        "a`b",
		"`a` + `b`":     "",
		"`a`+1":         "",
		"val + 1":       "",
		"to_days(`d`)":  "",
		"``":            "",
		"`unterminated": "",
		"":              "",
	} {
		got, ok := partitionKeyColumn(expr)
		if got != want || ok != (want != "") {
			t.Errorf("partitionKeyColumn(%q) = %q, %v; want %q", expr, got, ok, want)
		}
	}
}
