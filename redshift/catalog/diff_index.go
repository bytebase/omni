package catalog

import "sort"

// diffIndexes compares indexes between two versions of a relation.
// Identity key is the index Name.
// Indexes with ConstraintOID != 0 (PK/UNIQUE backing indexes) are skipped
// because they are managed by constraints, not as standalone indexes.
func diffIndexes(from, to *Catalog, fromRelOID, toRelOID uint32) []IndexDiffEntry {
	// Build name→*Index maps, filtering out constraint-backing indexes.
	fromMap := make(map[string]*Index)
	for _, idx := range from.IndexesOf(fromRelOID) {
		if idx.ConstraintOID != 0 {
			continue
		}
		fromMap[idx.Name] = idx
	}
	toMap := make(map[string]*Index)
	for _, idx := range to.IndexesOf(toRelOID) {
		if idx.ConstraintOID != 0 {
			continue
		}
		toMap[idx.Name] = idx
	}

	var result []IndexDiffEntry

	// Dropped: in from but not in to.
	for name, fromIdx := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, IndexDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fromIdx,
			})
		}
	}

	// Added or modified: in to.
	for name, toIdx := range toMap {
		fromIdx, ok := fromMap[name]
		if !ok {
			result = append(result, IndexDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     toIdx,
			})
			continue
		}

		// Both exist — compare fields.
		if indexesChanged(fromIdx, toIdx) {
			result = append(result, IndexDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fromIdx,
				To:     toIdx,
			})
		}
	}

	// Sort by name for determinism.
	sort.Slice(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// indexesChanged returns true if any compared property differs between two indexes.
func indexesChanged(a, b *Index) bool {
	if a.IsUnique != b.IsUnique {
		return true
	}
	if a.AccessMethod != b.AccessMethod {
		return true
	}
	if a.NKeyColumns != b.NKeyColumns {
		return true
	}
	if a.WhereClause != b.WhereClause {
		return true
	}
	if a.NullsNotDistinct != b.NullsNotDistinct {
		return true
	}

	// Compare Columns (attnums).
	if !int16SliceEqual(a.Columns, b.Columns) {
		return true
	}

	// Compare IndOption (per-column flags).
	if !int16SliceEqual(a.IndOption, b.IndOption) {
		return true
	}

	// Compare Exprs (deparsed expressions).
	if !stringSliceEqual(a.Exprs, b.Exprs) {
		return true
	}

	return false
}

// int16SliceEqual and stringSliceEqual are defined in diff_constraint.go.
