package catalog

// MySQL SDL generate — CHECK constraints (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so check diffs become DDL.
// Inert no-op placeholder returning nil until the check breadth node renders the ALTER
// TABLE ... ADD/DROP CHECK ops from TableDiffEntry.Checks (OpAddCheck/OpDropCheck). CHECK
// is 8.0-only (see normalize.go CheckSupported); the breadth node gates on the Normalizer's
// version. Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:check breadth node
func generateCheckDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no check ops (empty diff -> empty plan).
	return nil
}
