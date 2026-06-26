package catalog

// MySQL SDL generate — stored routines (functions + procedures) (scaffold stub).
//
// A single generate hook covers both routine kinds: wired into
// GenerateMigrationWithNormalizer (migration.go), it consumes BOTH SchemaDiff.Functions and
// SchemaDiff.Procedures and emits the CREATE/DROP FUNCTION|PROCEDURE ops
// (OpCreateFunction/OpDropFunction/OpCreateProcedure/OpDropProcedure, priorityRoutine).
// Inert no-op placeholder returning nil until the routine breadth node implements it.
// Signature mirrors generateTableDDL (migration_table.go).
//
// implemented by omni:routine breadth node
func generateRoutineDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no routine ops (empty diff -> empty plan).
	return nil
}
