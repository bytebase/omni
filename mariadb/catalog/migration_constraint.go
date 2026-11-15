package catalog

// generateConstraintDDL returns no ops: primary key and unique constraint changes are
// rendered as index ops by generateIndexDDL.
func generateConstraintDDL(_, _ *Catalog, _ *SchemaDiff, _ *Normalizer) []MigrationOp {
	return nil
}
