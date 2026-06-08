package catalog

import "sort"

// diffEnums compares enum types between two catalogs and returns diff entries.
func diffEnums(from, to *Catalog) []EnumDiffEntry {
	type enumKey struct {
		Schema string
		Name   string
	}

	// resolveEnums builds a map of schema.name → label list for all enum types in a catalog.
	resolveEnums := func(c *Catalog) map[enumKey][]string {
		m := make(map[enumKey][]string, len(c.enumTypes))
		for typeOID, et := range c.enumTypes {
			bt := c.typeByOID[typeOID]
			if bt == nil {
				continue
			}
			s := c.schemas[bt.Namespace]
			if s == nil {
				continue
			}
			labels := make([]string, len(et.Values))
			for i, ev := range et.Values {
				labels[i] = ev.Label
			}
			m[enumKey{Schema: s.Name, Name: bt.TypeName}] = labels
		}
		return m
	}

	fromMap := resolveEnums(from)
	toMap := resolveEnums(to)

	var entries []EnumDiffEntry

	// Dropped: in from but not in to.
	for k, fromVals := range fromMap {
		if _, ok := toMap[k]; !ok {
			entries = append(entries, EnumDiffEntry{
				Action:     DiffDrop,
				SchemaName: k.Schema,
				Name:       k.Name,
				FromValues: fromVals,
			})
		}
	}

	// Added: in to but not in from.
	for k, toVals := range toMap {
		if _, ok := fromMap[k]; !ok {
			entries = append(entries, EnumDiffEntry{
				Action:     DiffAdd,
				SchemaName: k.Schema,
				Name:       k.Name,
				ToValues:   toVals,
			})
		}
	}

	// Modified: in both, compare label slices.
	for k, fromVals := range fromMap {
		toVals, ok := toMap[k]
		if !ok {
			continue
		}
		if !stringSliceEqual(fromVals, toVals) {
			entries = append(entries, EnumDiffEntry{
				Action:     DiffModify,
				SchemaName: k.Schema,
				Name:       k.Name,
				FromValues: fromVals,
				ToValues:   toVals,
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

// stringSliceEqual is defined in diff_constraint.go.
