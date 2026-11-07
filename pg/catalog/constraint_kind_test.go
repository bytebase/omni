package catalog

// nonNullabilityConstraints keeps key/check/FK tests focused on those kinds.
// PostgreSQL 18 also stores NOT NULL constraints, covered by the native oracle.
func nonNullabilityConstraints(cons []*Constraint) []*Constraint {
	var result []*Constraint
	for _, con := range cons {
		if con.Type != ConstraintNotNull {
			result = append(result, con)
		}
	}
	return result
}
