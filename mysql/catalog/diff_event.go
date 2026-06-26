package catalog

// MySQL SDL diff — scheduled events (scaffold stub).
//
// Wired into DiffWithNormalizer (diff.go) so the SchemaDiff reports event changes via
// SchemaDiff.Events. Inert no-op placeholder returning nil until the event breadth node
// implements the real comparison. Signature mirrors diffTables (diff_table.go): whole-
// catalog, taking from/to *Catalog and the version-fixing Normalizer (events are database-
// level objects keyed by (database, name)).
//
// implemented by omni:event breadth node
func diffEvents(_, _ *Catalog, _ *Normalizer) []EventDiffEntry {
	// Scaffold: returns nil so the event diff stays empty (self-diff inert).
	return nil
}
