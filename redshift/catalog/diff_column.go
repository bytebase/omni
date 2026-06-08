package catalog

import "sort"

// diffColumns compares columns between two versions of a relation.
// Identity key is the column Name within the relation.
func diffColumns(from, to *Catalog, fromRel, toRel *Relation) []ColumnDiffEntry {
	// Build name→*Column maps.
	fromMap := make(map[string]*Column, len(fromRel.Columns))
	for _, col := range fromRel.Columns {
		fromMap[col.Name] = col
	}
	toMap := make(map[string]*Column, len(toRel.Columns))
	for _, col := range toRel.Columns {
		toMap[col.Name] = col
	}

	var result []ColumnDiffEntry

	// Dropped: in from but not in to.
	for name, fromCol := range fromMap {
		if _, ok := toMap[name]; !ok {
			result = append(result, ColumnDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fromCol,
			})
		}
	}

	// Added or modified: in to.
	for name, toCol := range toMap {
		fromCol, ok := fromMap[name]
		if !ok {
			result = append(result, ColumnDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     toCol,
			})
			continue
		}

		// Both exist — compare fields.
		if columnsChanged(from, to, fromCol, toCol) {
			result = append(result, ColumnDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fromCol,
				To:     toCol,
			})
		}
	}

	// Sort by column name for determinism.
	sort.Slice(result, func(i, j int) bool {
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// columnsChanged returns true if any compared property differs between two columns.
func columnsChanged(from, to *Catalog, a, b *Column) bool {
	// Type comparison via FormatType (never compare raw OIDs).
	if from.FormatType(a.TypeOID, a.TypeMod) != to.FormatType(b.TypeOID, b.TypeMod) {
		return true
	}
	if a.NotNull != b.NotNull {
		return true
	}
	if a.HasDefault != b.HasDefault {
		return true
	}
	if a.Default != b.Default {
		return true
	}
	if a.Identity != b.Identity {
		return true
	}
	if a.Generated != b.Generated {
		return true
	}
	if a.GenerationExpr != b.GenerationExpr {
		return true
	}
	if a.CollationName != b.CollationName {
		return true
	}
	return false
}
