package catalog

import "sort"

// diffIndexes compares indexes between two versions of a relation.
// Identity key is the index Name.
// Indexes with ConstraintOID != 0 (PK/UNIQUE backing indexes) are skipped
// because they are managed by constraints, not as standalone indexes. So is a
// partition's index attached to its parent's that from has too: PostgreSQL
// drops it with the parent's index or the partition. One from has unattached
// is modified, to be attached, and one only to has is added and attached.
func diffIndexes(from, to *Catalog, fromRelOID, toRelOID uint32) []IndexDiffEntry {
	// Build name→*Index maps, filtering out constraint-backing indexes.
	fromMap := make(map[string]*Index)
	fromConIdxMap := make(map[string]*Index)
	for _, idx := range from.IndexesOf(fromRelOID) {
		if idx.ConstraintOID != 0 {
			if idx.ParentIndexOID == 0 {
				fromConIdxMap[idx.Name] = idx
			}
			continue
		}
		fromMap[idx.Name] = idx
	}
	toMap := make(map[string]*Index)
	toConIdxMap := make(map[string]*Index)
	for _, idx := range to.IndexesOf(toRelOID) {
		if idx.ConstraintOID != 0 {
			if idx.ParentIndexOID == 0 {
				toConIdxMap[idx.Name] = idx
			}
			continue
		}
		toMap[idx.Name] = idx
	}
	for name, fromIdx := range fromMap {
		toIdx := toMap[name]
		if fromIdx.ParentIndexOID != 0 && (toIdx == nil || toIdx.ParentIndexOID != 0) {
			delete(fromMap, name)
			delete(toMap, name)
		}
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
		if toIdx.ParentIndexOID != 0 || indexesChanged(from, to, fromIdx, toIdx) {
			result = append(result, IndexDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fromIdx,
				To:     toIdx,
			})
		}
	}

	// Constraint-backed indexes: surface an in-place column rename when the
	// backing index itself was not structurally changed.
	for name, toIdx := range toConIdxMap {
		fromIdx, ok := fromConIdxMap[name]
		if !ok {
			continue
		}
		if !indexesStructurallyChanged(from, to, fromIdx, toIdx) &&
			!stringSliceEqual(effectiveIndexColumnNames(from, fromIdx), effectiveIndexColumnNames(to, toIdx)) {
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
func indexesChanged(fromCat, toCat *Catalog, a, b *Index) bool {
	if indexesStructurallyChanged(fromCat, toCat, a, b) {
		return true
	}
	if !stringSliceEqual(effectiveIndexColumnNames(fromCat, a), effectiveIndexColumnNames(toCat, b)) {
		return true
	}
	return false
}

func indexesStructurallyChanged(fromCat, toCat *Catalog, a, b *Index) bool {
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

	if !sameColumns(fromCat, a.RelOID, a.Columns, toCat, b.RelOID, b.Columns) {
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
	if !stringSliceEqual(a.OpClasses, b.OpClasses) {
		return true
	}
	if !stringSliceEqual(a.OpClassOptions, b.OpClassOptions) {
		return true
	}
	if !stringSlicesEqual(a.Collations, b.Collations) {
		return true
	}

	return false
}

func cloneStringSlices(in [][]string) [][]string {
	if in == nil {
		return nil
	}
	out := make([][]string, len(in))
	for i, s := range in {
		if s != nil {
			out[i] = append([]string(nil), s...)
		}
	}
	return out
}

func stringSlicesEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !stringSliceEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func defaultIndexColumnNames(cat *Catalog, idx *Index) []string {
	if idx == nil || len(idx.Columns) == 0 {
		return nil
	}
	rel := cat.GetRelationByOID(idx.RelOID)
	raw := make([]string, 0, len(idx.Columns))
	exprIdx := 0
	for _, attnum := range idx.Columns {
		if attnum == 0 {
			name := "expr"
			if exprIdx < len(idx.Exprs) {
				if parsed, err := parseScalar("SELECT " + idx.Exprs[exprIdx]); err == nil {
					name = figureIndexColname(parsed)
				}
				exprIdx++
			}
			raw = append(raw, name)
		} else if rel != nil {
			raw = append(raw, resolveAttNum(rel, attnum))
		} else {
			raw = append(raw, "expr")
		}
	}
	return chooseIndexColumnNames(raw)
}

func effectiveIndexColumnNames(cat *Catalog, idx *Index) []string {
	if idx == nil {
		return nil
	}
	if len(idx.ColumnNames) == len(idx.Columns) && len(idx.ColumnNames) > 0 {
		return idx.ColumnNames
	}
	return defaultIndexColumnNames(cat, idx)
}

// int16SliceEqual and stringSliceEqual are defined in diff_constraint.go.
