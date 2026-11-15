package catalog

// diffConstraints reports nothing: primary and unique keys are compared as
// indexes by diffIndexes, foreign keys by diffForeignKeys, and checks by
// diffChecks.
func diffConstraints(_, _ *Table, _ *Normalizer) []ConstraintDiffEntry {
	return nil
}
