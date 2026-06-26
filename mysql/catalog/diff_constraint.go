package catalog

// MySQL SDL diff — PRIMARY KEY / UNIQUE constraints (scaffold stub).
//
// Wired into compareTable (diff_table.go) so a modified table reports its PK/UNIQUE
// constraint changes via TableDiffEntry.Constraints. Inert no-op placeholder returning
// nil until the constraint breadth node implements the real comparison. Signature mirrors
// diffColumns: per-table, taking from/to *Table and the version-fixing Normalizer.
//
// implemented by omni:constraint breadth node
func diffConstraints(_, _ *Table, _ *Normalizer) []ConstraintDiffEntry {
	// Scaffold: returns nil so the constraint sub-diff stays empty (self-diff inert).
	return nil
}
