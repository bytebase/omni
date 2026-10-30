package catalog

// MySQL SDL generate — PRIMARY KEY / UNIQUE constraints (scaffold stub).
//
// Wired into GenerateMigrationWithNormalizer (migration.go) so constraint diffs become DDL.
// Inert no-op placeholder returning nil until the constraint breadth node renders the
// ALTER TABLE ... ADD/DROP PRIMARY KEY|UNIQUE ops from TableDiffEntry.Constraints
// (OpAddConstraint/OpDropConstraint, priorityConstraint). Signature mirrors generateTableDDL.
//
// implemented by omni:constraint breadth node
func generateConstraintDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	// Scaffold: returns nil so the plan gains no constraint ops (empty diff -> empty plan).
	return nil
}
