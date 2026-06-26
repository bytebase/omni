package catalog

// MySQL SDL diff — stored routines (functions + procedures) (scaffold stub).
//
// Two hooks, because the dispatcher (DiffWithNormalizer in diff.go) packs functions and
// procedures into separate SchemaDiff slices (SchemaDiff.Functions / SchemaDiff.Procedures),
// both element type RoutineDiffEntry (distinguished by RoutineDiffEntry.IsProcedure). The
// routine breadth node fills BOTH here. Inert no-op placeholders returning nil until then.
// Signatures mirror diffTables (diff_table.go): whole-catalog, taking from/to *Catalog and
// the version-fixing Normalizer (routines are database-level objects keyed by (database, name)).
//
// implemented by omni:routine breadth node
func diffFunctions(_, _ *Catalog, _ *Normalizer) []RoutineDiffEntry {
	// Scaffold: returns nil so the function diff stays empty (self-diff inert).
	return nil
}

// diffProcedures is the stored-procedure half of the routine breadth node; see diffFunctions.
//
// implemented by omni:routine breadth node
func diffProcedures(_, _ *Catalog, _ *Normalizer) []RoutineDiffEntry {
	// Scaffold: returns nil so the procedure diff stays empty (self-diff inert).
	return nil
}
