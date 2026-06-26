package catalog

// MySQL SDL generate — table partitioning (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so a partitioning change
// (TableDiffEntry.PartitionChanged) becomes DDL. Inert no-op placeholder returning nil
// until the partition breadth node renders the ALTER TABLE ... PARTITION BY / REMOVE
// PARTITIONING op. Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:partition breadth node
func generatePartitionDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no partition ops (empty diff -> empty plan).
	return nil
}
