package catalog

// MySQL SDL diff — views (scaffold stub).
//
// Wired into DiffWithNormalizer (diff.go) so the SchemaDiff reports view changes via
// SchemaDiff.Views. Inert no-op placeholder returning nil until the view breadth node
// implements the real comparison. Signature mirrors diffTables (diff_table.go): whole-
// catalog, taking from/to *Catalog and the version-fixing Normalizer (views are database-
// level objects keyed by (database, name), not per-table).
//
// implemented by omni:view breadth node
func diffViews(_, _ *Catalog, _ *Normalizer) []ViewDiffEntry {
	// Scaffold: returns nil so the view diff stays empty (self-diff inert).
	return nil
}
