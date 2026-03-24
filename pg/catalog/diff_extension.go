package catalog

import "sort"

// diffExtensions compares extensions between two catalogs and returns diff entries.
func diffExtensions(from, to *Catalog) []ExtensionDiffEntry {
	type extInfo struct {
		ext        *Extension
		schemaName string
	}

	// resolveExts builds a name → extInfo map for all extensions in a catalog.
	resolveExts := func(c *Catalog) map[string]*extInfo {
		m := make(map[string]*extInfo, len(c.extensions))
		for _, ext := range c.extensions {
			schemaName := ""
			if s := c.schemas[ext.SchemaOID]; s != nil {
				schemaName = s.Name
			}
			m[ext.Name] = &extInfo{
				ext:        ext,
				schemaName: schemaName,
			}
		}
		return m
	}

	fromMap := resolveExts(from)
	toMap := resolveExts(to)

	var entries []ExtensionDiffEntry

	// Dropped: in from but not in to.
	for name, fi := range fromMap {
		if _, ok := toMap[name]; !ok {
			entries = append(entries, ExtensionDiffEntry{
				Action: DiffDrop,
				Name:   name,
				From:   fi.ext,
			})
		}
	}

	// Added: in to but not in from.
	for name, ti := range toMap {
		if _, ok := fromMap[name]; !ok {
			entries = append(entries, ExtensionDiffEntry{
				Action: DiffAdd,
				Name:   name,
				To:     ti.ext,
			})
		}
	}

	// Modified: in both, compare schema and relocatable.
	for name, fi := range fromMap {
		ti, ok := toMap[name]
		if !ok {
			continue
		}
		if fi.schemaName != ti.schemaName || fi.ext.Relocatable != ti.ext.Relocatable {
			entries = append(entries, ExtensionDiffEntry{
				Action: DiffModify,
				Name:   name,
				From:   fi.ext,
				To:     ti.ext,
			})
		}
	}

	// Sort by name for determinism.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries
}
