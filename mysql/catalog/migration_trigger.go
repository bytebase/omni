package catalog

// MySQL SDL generate — triggers (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so trigger diffs
// (SchemaDiff.Triggers) become DDL. Inert no-op placeholder returning nil until the trigger
// breadth node renders the CREATE/DROP TRIGGER ops (OpCreateTrigger/OpDropTrigger,
// priorityTrigger). Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:trigger breadth node
func generateTriggerDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no trigger ops (empty diff -> empty plan).
	return nil
}
