package catalog

// MySQL SDL generate — scheduled events (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so event diffs (SchemaDiff.Events)
// become DDL. Inert no-op placeholder returning nil until the event breadth node renders the
// CREATE/DROP EVENT ops (OpCreateEvent/OpDropEvent, priorityEvent). Signature mirrors
// generateTableDDL (migration_table.go).
//
// implemented by omni:event breadth node
func generateEventDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no event ops (empty diff -> empty plan).
	return nil
}
