package catalog

import "sort"

// diffSchemas compares schemas between two catalogs and returns diff entries.
// Only user schemas are compared (pg_catalog and pg_toast are excluded).
func diffSchemas(from, to *Catalog) []SchemaDiffEntry {
	fromSchemas := from.UserSchemas()
	toSchemas := to.UserSchemas()

	// Build name-based maps.
	fromMap := make(map[string]*Schema, len(fromSchemas))
	for _, s := range fromSchemas {
		fromMap[s.Name] = s
	}
	toMap := make(map[string]*Schema, len(toSchemas))
	for _, s := range toSchemas {
		toMap[s.Name] = s
	}

	var entries []SchemaDiffEntry

	// Iterate from: key not in to → DiffDrop
	for _, s := range fromSchemas {
		if _, ok := toMap[s.Name]; !ok {
			entries = append(entries, SchemaDiffEntry{
				Action: DiffDrop,
				Name:   s.Name,
				From:   s,
			})
		}
	}

	// Iterate to: key not in from → DiffAdd
	for _, s := range toSchemas {
		if _, ok := fromMap[s.Name]; !ok {
			entries = append(entries, SchemaDiffEntry{
				Action: DiffAdd,
				Name:   s.Name,
				To:     s,
			})
		}
	}

	// Both exist → compare fields
	for _, fs := range fromSchemas {
		ts, ok := toMap[fs.Name]
		if !ok {
			continue
		}
		if fs.Owner != ts.Owner {
			entries = append(entries, SchemaDiffEntry{
				Action: DiffModify,
				Name:   fs.Name,
				From:   fs,
				To:     ts,
			})
		}
	}

	// Sort by name for determinism.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries
}
