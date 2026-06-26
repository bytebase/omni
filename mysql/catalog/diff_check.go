package catalog

// MySQL SDL diff — CHECK constraints (scaffold stub).
//
// Wired into compareTable (diff_table.go) so a modified table reports its CHECK-constraint
// changes via TableDiffEntry.Checks. Inert no-op placeholder returning nil until the check
// breadth node implements the real comparison. CHECK is 8.0-only (unrepresentable on 5.7;
// see normalize.go CheckSupported) — the breadth node gates on the Normalizer's version.
// Signature mirrors diffColumns: per-table, taking from/to *Table and the Normalizer.
//
// implemented by omni:check breadth node
func diffChecks(_, _ *Table, _ *Normalizer) []CheckDiffEntry {
	// Scaffold: returns nil so the check sub-diff stays empty (self-diff inert).
	return nil
}
