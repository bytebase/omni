package catalog

import "sort"

// diffComposites compares composite types (RelKind='c') between two catalogs
// and returns diff entries.
func diffComposites(from, to *Catalog) []CompositeTypeDiffEntry {
	fromMap := buildCompositeMap(from)
	toMap := buildCompositeMap(to)

	var entries []CompositeTypeDiffEntry

	// Dropped: in from but not in to.
	for k, fromRel := range fromMap {
		if _, ok := toMap[k]; !ok {
			entries = append(entries, CompositeTypeDiffEntry{
				Action:     DiffDrop,
				SchemaName: k.schema,
				Name:       k.name,
				From:       fromRel,
			})
		}
	}

	// Added: in to but not in from.
	for k, toRel := range toMap {
		if _, ok := fromMap[k]; !ok {
			entries = append(entries, CompositeTypeDiffEntry{
				Action:     DiffAdd,
				SchemaName: k.schema,
				Name:       k.name,
				To:         toRel,
			})
		}
	}

	// Modified: in both, compare columns.
	for k, fromRel := range fromMap {
		toRel, ok := toMap[k]
		if !ok {
			continue
		}
		if !compositeColumnsEqual(from, to, fromRel, toRel) {
			entries = append(entries, CompositeTypeDiffEntry{
				Action:     DiffModify,
				SchemaName: k.schema,
				Name:       k.name,
				From:       fromRel,
				To:         toRel,
			})
		}
	}

	// Sort by schema + name for determinism.
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SchemaName != entries[j].SchemaName {
			return entries[i].SchemaName < entries[j].SchemaName
		}
		return entries[i].Name < entries[j].Name
	})

	return entries
}

// buildCompositeMap builds a map of (schema, name) -> *Relation for composite
// types (RelKind='c').
func buildCompositeMap(c *Catalog) map[relKey]*Relation {
	m := make(map[relKey]*Relation)
	for _, s := range c.UserSchemas() {
		for _, rel := range s.Relations {
			if rel.RelKind != 'c' {
				continue
			}
			m[relKey{schema: s.Name, name: rel.Name}] = rel
		}
	}
	return m
}

// compositeColumnsEqual returns true if two composite type relations have
// the same columns (by name, type, and order).
func compositeColumnsEqual(fromCat, toCat *Catalog, from, to *Relation) bool {
	if len(from.Columns) != len(to.Columns) {
		return false
	}
	for i := range from.Columns {
		fc := from.Columns[i]
		tc := to.Columns[i]
		if fc.Name != tc.Name {
			return false
		}
		if fromCat.FormatType(fc.TypeOID, fc.TypeMod) != toCat.FormatType(tc.TypeOID, tc.TypeMod) {
			return false
		}
	}
	return true
}
