package catalog

import "sort"

// diffViews compares views (RelKind='v') and materialized views (RelKind='m')
// between two catalogs. Returns RelationDiffEntry entries (same type as tables).
func diffViews(from, to *Catalog) []RelationDiffEntry {
	fromMap := buildViewMap(from)
	toMap := buildViewMap(to)

	var result []RelationDiffEntry

	// Dropped: in from but not in to.
	for key, fromRel := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, RelationDiffEntry{
				Action:     DiffDrop,
				SchemaName: key.schema,
				Name:       key.name,
				From:       fromRel,
			})
		}
	}

	// Added or modified: in to.
	for key, toRel := range toMap {
		fromRel, ok := fromMap[key]
		if !ok {
			result = append(result, RelationDiffEntry{
				Action:     DiffAdd,
				SchemaName: key.schema,
				Name:       key.name,
				To:         toRel,
			})
			continue
		}

		// Both exist — check for modifications.
		if entry, changed := compareView(from, to, key, fromRel, toRel); changed {
			result = append(result, entry)
		}
	}

	// Sort for determinism.
	sort.Slice(result, func(i, j int) bool {
		if result[i].SchemaName != result[j].SchemaName {
			return result[i].SchemaName < result[j].SchemaName
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Action < result[j].Action
	})

	return result
}

// buildViewMap builds a map of (schema, name) -> *Relation for views and matviews.
func buildViewMap(c *Catalog) map[relKey]*Relation {
	m := make(map[relKey]*Relation)
	for _, s := range c.UserSchemas() {
		for _, rel := range s.Relations {
			if rel.RelKind != 'v' && rel.RelKind != 'm' {
				continue
			}
			m[relKey{schema: s.Name, name: rel.Name}] = rel
		}
	}
	return m
}

// compareView checks whether two views/matviews with the same identity differ.
func compareView(fromCat, toCat *Catalog, key relKey, from, to *Relation) (RelationDiffEntry, bool) {
	changed := false

	// Compare view definition using deparse.
	fromDef := getViewDef(fromCat, key, from)
	toDef := getViewDef(toCat, key, to)
	if fromDef != toDef {
		changed = true
	}

	// Compare CheckOption (views only).
	if from.CheckOption != to.CheckOption {
		changed = true
	}

	// Column sub-diff (views are relations too).
	cols := diffColumns(fromCat, toCat, from, to)
	if len(cols) > 0 {
		changed = true
	}

	if !changed {
		return RelationDiffEntry{}, false
	}

	return RelationDiffEntry{
		Action:     DiffModify,
		SchemaName: key.schema,
		Name:       key.name,
		From:       from,
		To:         to,
		Columns:    cols,
	}, true
}

// getViewDef returns the deparsed definition of a view or matview.
func getViewDef(c *Catalog, key relKey, rel *Relation) string {
	var def string
	var err error
	if rel.RelKind == 'm' {
		def, err = c.GetMatViewDefinition(key.schema, key.name)
	} else {
		def, err = c.GetViewDefinition(key.schema, key.name)
	}
	if err != nil {
		return ""
	}
	return def
}
