package catalog

// MySQL SDL diff — secondary indexes (scaffold stub).
//
// Wired into compareTable (diff_table.go) so a modified table reports its index
// changes via TableDiffEntry.Indexes. This is an inert no-op placeholder: it returns
// nil, so the table sub-diff is unchanged until the index breadth node implements the
// real comparison. Signature mirrors diffColumns (diff_column.go): per-table, taking
// the from/to *Table and the version-fixing Normalizer.
//
// implemented by omni:index breadth node
func diffIndexes(_, _ *Table, _ *Normalizer) []IndexDiffEntry {
	// Scaffold: returns nil so the index sub-diff stays empty (self-diff inert).
	return nil
}
