package catalog

// MySQL SDL diff — triggers (scaffold stub).
//
// Wired into DiffWithNormalizer (diff.go) so the SchemaDiff reports trigger changes via
// SchemaDiff.Triggers. Inert no-op placeholder returning nil until the trigger breadth node
// implements the real comparison. Signature mirrors diffTables (diff_table.go): whole-
// catalog, taking from/to *Catalog and the version-fixing Normalizer (MySQL triggers are
// database-level objects keyed by (database, name), each bound to a table via Trigger.Table).
//
// implemented by omni:trigger breadth node
func diffTriggers(_, _ *Catalog, _ *Normalizer) []TriggerDiffEntry {
	// Scaffold: returns nil so the trigger diff stays empty (self-diff inert).
	return nil
}
