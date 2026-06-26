package catalog

// MySQL SDL generate — secondary indexes (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so index diffs become DDL.
// Inert no-op placeholder returning nil until the index breadth node renders the ALTER
// TABLE ... ADD/DROP INDEX ops from TableDiffEntry.Indexes (OpAddIndex/OpDropIndex,
// priorityIndex). Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:index breadth node
func generateIndexDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no index ops (empty diff -> empty plan).
	return nil
}
