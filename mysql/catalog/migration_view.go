package catalog

// MySQL SDL generate — views (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so view diffs (SchemaDiff.Views)
// become DDL. Inert no-op placeholder returning nil until the view breadth node renders the
// CREATE/DROP VIEW ops (OpCreateView/OpDropView, priorityView). Signature mirrors
// generateTableDDL (migration_table.go).
//
// implemented by omni:view breadth node
func generateViewDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no view ops (empty diff -> empty plan).
	return nil
}
