package catalog

// MySQL SDL diff — table partitioning (scaffold stub).
//
// Wired into compareTable (diff_table.go) so a modified table records a partitioning
// change via the TableDiffEntry.PartitionChanged flag (a bool, not a slice — MySQL
// partitioning is a single per-table clause, compared as one PARTITION BY ... spec via
// Table.Partitioning). Inert no-op placeholder returning false until the partition breadth
// node implements the real comparison. Signature mirrors diffColumns: per-table, taking
// from/to *Table and the version-fixing Normalizer.
//
// implemented by omni:partition breadth node
func diffPartitions(_, _ *Table, _ *Normalizer) bool {
	// Scaffold: returns false so PartitionChanged stays unset (self-diff inert).
	return false
}
