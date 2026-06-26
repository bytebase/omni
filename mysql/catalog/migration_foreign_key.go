package catalog

// MySQL SDL generate — foreign keys (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so FK diffs become DDL.
// Inert no-op placeholder returning nil until the FK breadth node renders the ALTER
// TABLE ... ADD/DROP FOREIGN KEY ops from TableDiffEntry.ForeignKeys. FK creation is
// deferred to PhasePost (priorityForeignKey) so every referenced table exists first —
// mirroring PG. Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:foreign-key breadth node
func generateForeignKeyDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no FK ops (empty diff -> empty plan).
	return nil
}
