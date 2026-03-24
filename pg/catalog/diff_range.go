package catalog

import "sort"

// diffRanges compares range types between two catalogs and returns diff entries.
func diffRanges(from, to *Catalog) []RangeDiffEntry {
	type rangeKey struct {
		Schema string
		Name   string
	}

	type rangeInfo struct {
		key     rangeKey
		subtype string
		rt      *RangeType
	}

	// resolveRanges builds a map of schema.name → rangeInfo for all range types in a catalog.
	resolveRanges := func(c *Catalog) map[rangeKey]*rangeInfo {
		m := make(map[rangeKey]*rangeInfo, len(c.rangeTypes))
		for typeOID, rt := range c.rangeTypes {
			bt := c.typeByOID[typeOID]
			if bt == nil {
				continue
			}
			s := c.schemas[bt.Namespace]
			if s == nil {
				continue
			}
			k := rangeKey{Schema: s.Name, Name: bt.TypeName}
			m[k] = &rangeInfo{
				key:     k,
				subtype: c.FormatType(rt.SubTypeOID, -1),
				rt:      rt,
			}
		}
		return m
	}

	fromMap := resolveRanges(from)
	toMap := resolveRanges(to)

	var entries []RangeDiffEntry

	// Dropped: in from but not in to.
	for k, fi := range fromMap {
		if _, ok := toMap[k]; !ok {
			entries = append(entries, RangeDiffEntry{
				Action:     DiffDrop,
				SchemaName: k.Schema,
				Name:       k.Name,
				From:       fi.rt,
			})
		}
	}

	// Added: in to but not in from.
	for k, ti := range toMap {
		if _, ok := fromMap[k]; !ok {
			entries = append(entries, RangeDiffEntry{
				Action:     DiffAdd,
				SchemaName: k.Schema,
				Name:       k.Name,
				To:         ti.rt,
			})
		}
	}

	// Modified: in both, compare subtype.
	for k, fi := range fromMap {
		ti, ok := toMap[k]
		if !ok {
			continue
		}
		if fi.subtype != ti.subtype {
			entries = append(entries, RangeDiffEntry{
				Action:     DiffModify,
				SchemaName: k.Schema,
				Name:       k.Name,
				From:       fi.rt,
				To:         ti.rt,
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
