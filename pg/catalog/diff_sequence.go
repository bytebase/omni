package catalog

import "sort"

// diffSequences compares standalone sequences between two catalogs.
// Identity key is (schemaName, seqName).
// Sequences with OwnerRelOID != 0 (SERIAL/IDENTITY-owned) are skipped
// because they are managed by the owning column, not as standalone objects.
func diffSequences(from, to *Catalog) []SequenceDiffEntry {
	type seqKey struct {
		schema string
		name   string
	}

	// Build name-based maps from both catalogs, skipping owned sequences.
	fromMap := make(map[seqKey]*Sequence)
	for _, s := range from.UserSchemas() {
		for _, seq := range s.Sequences {
			if seq.OwnerRelOID != 0 {
				continue
			}
			fromMap[seqKey{schema: s.Name, name: seq.Name}] = seq
		}
	}

	toMap := make(map[seqKey]*Sequence)
	for _, s := range to.UserSchemas() {
		for _, seq := range s.Sequences {
			if seq.OwnerRelOID != 0 {
				continue
			}
			toMap[seqKey{schema: s.Name, name: seq.Name}] = seq
		}
	}

	var result []SequenceDiffEntry

	// Dropped: in from but not in to.
	for key, fromSeq := range fromMap {
		if _, ok := toMap[key]; !ok {
			result = append(result, SequenceDiffEntry{
				Action:     DiffDrop,
				SchemaName: key.schema,
				Name:       key.name,
				From:       fromSeq,
			})
		}
	}

	// Added or modified: in to.
	for key, toSeq := range toMap {
		fromSeq, ok := fromMap[key]
		if !ok {
			result = append(result, SequenceDiffEntry{
				Action:     DiffAdd,
				SchemaName: key.schema,
				Name:       key.name,
				To:         toSeq,
			})
			continue
		}

		// Both exist -- compare fields.
		if sequencesChanged(from, to, fromSeq, toSeq) {
			result = append(result, SequenceDiffEntry{
				Action:     DiffModify,
				SchemaName: key.schema,
				Name:       key.name,
				From:       fromSeq,
				To:         toSeq,
			})
		}
	}

	// Sort by (schema, name) for determinism.
	sort.Slice(result, func(i, j int) bool {
		if result[i].SchemaName != result[j].SchemaName {
			return result[i].SchemaName < result[j].SchemaName
		}
		return result[i].Name < result[j].Name
	})

	return result
}

// sequencesChanged returns true if any compared property differs.
// Type is compared via FormatType (never raw OIDs).
func sequencesChanged(fromCat, toCat *Catalog, a, b *Sequence) bool {
	// Compare type via FormatType.
	if fromCat.FormatType(a.TypeOID, -1) != toCat.FormatType(b.TypeOID, -1) {
		return true
	}
	if a.Increment != b.Increment {
		return true
	}
	if a.MinValue != b.MinValue {
		return true
	}
	if a.MaxValue != b.MaxValue {
		return true
	}
	if a.Start != b.Start {
		return true
	}
	if a.CacheValue != b.CacheValue {
		return true
	}
	if a.Cycle != b.Cycle {
		return true
	}
	return false
}
