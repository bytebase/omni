package catalog

// MySQL SDL diff — foreign keys (scaffold stub).
//
// Wired into compareTable (diff_table.go) so a modified table reports its FK changes
// via TableDiffEntry.ForeignKeys. Inert no-op placeholder returning nil until the FK
// breadth node implements the real comparison. Signature mirrors diffColumns: per-table,
// taking from/to *Table and the version-fixing Normalizer. FK identity/equality (the
// referenced table/columns and ON DELETE/UPDATE actions) is the breadth node's concern.
//
// implemented by omni:foreign-key breadth node
func diffForeignKeys(_, _ *Table, _ *Normalizer) []ForeignKeyDiffEntry {
	// Scaffold: returns nil so the FK sub-diff stays empty (self-diff inert).
	return nil
}
