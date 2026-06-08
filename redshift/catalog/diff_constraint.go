package catalog

import "sort"

// diffConstraints compares constraints on two versions of the same relation.
// Identity key is the constraint Name.
func diffConstraints(from, to *Catalog, fromRelOID, toRelOID uint32) []ConstraintDiffEntry {
	fromCons := from.ConstraintsOf(fromRelOID)
	toCons := to.ConstraintsOf(toRelOID)

	// Build name → *Constraint maps.
	fromMap := make(map[string]*Constraint, len(fromCons))
	for _, c := range fromCons {
		fromMap[c.Name] = c
	}
	toMap := make(map[string]*Constraint, len(toCons))
	for _, c := range toCons {
		toMap[c.Name] = c
	}

	var result []ConstraintDiffEntry

	// Dropped: in from but not in to.
	for name, fc := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, ConstraintDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fc,
			})
		}
	}

	// Added or modified: in to.
	for name, tc := range toMap {
		fc, ok := fromMap[name]
		if !ok {
			result = append(result, ConstraintDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     tc,
			})
			continue
		}

		// Both exist — compare fields.
		if constraintsEqual(from, to, fc, tc) {
			continue
		}
		result = append(result, ConstraintDiffEntry{
			Action: DiffModify,
			Name:   name,
			From:   fc,
			To:     tc,
		})
	}

	// Sort by name for determinism.
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// constraintsEqual returns true if two constraints are semantically identical.
func constraintsEqual(from, to *Catalog, a, b *Constraint) bool {
	if a.Type != b.Type {
		return false
	}
	if !int16SliceEqual(a.Columns, b.Columns) {
		return false
	}
	if a.CheckExpr != b.CheckExpr {
		return false
	}
	if a.FKUpdAction != b.FKUpdAction {
		return false
	}
	if a.FKDelAction != b.FKDelAction {
		return false
	}
	if a.FKMatchType != b.FKMatchType {
		return false
	}
	if a.Deferrable != b.Deferrable {
		return false
	}
	if a.Deferred != b.Deferred {
		return false
	}

	// For FK constraints, compare referenced table by resolved name and referenced columns.
	if a.Type == ConstraintFK {
		if !int16SliceEqual(a.FColumns, b.FColumns) {
			return false
		}
		// Resolve referenced table names via GetRelationByOID.
		if !fkRefTableEqual(from, to, a, b) {
			return false
		}
	}

	// EXCLUDE: compare operators.
	if a.Type == ConstraintExclude {
		if !stringSliceEqual(a.ExclOps, b.ExclOps) {
			return false
		}
	}

	return true
}

// fkRefTableEqual compares the referenced table of two FK constraints
// by resolving FRelOID to schema.name in their respective catalogs.
func fkRefTableEqual(from, to *Catalog, a, b *Constraint) bool {
	fromRel := from.GetRelationByOID(a.FRelOID)
	toRel := to.GetRelationByOID(b.FRelOID)
	if fromRel == nil && toRel == nil {
		return true
	}
	if fromRel == nil || toRel == nil {
		return false
	}
	if fromRel.Name != toRel.Name {
		return false
	}
	// Compare schema names.
	fromSchema := ""
	if fromRel.Schema != nil {
		fromSchema = fromRel.Schema.Name
	}
	toSchema := ""
	if toRel.Schema != nil {
		toSchema = toRel.Schema.Name
	}
	return fromSchema == toSchema
}

// int16SliceEqual returns true if two int16 slices are identical.
func int16SliceEqual(a, b []int16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// stringSliceEqual returns true if two string slices are identical.
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
