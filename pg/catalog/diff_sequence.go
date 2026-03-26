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
			// Check if the sequence exists in `to` as an owned sequence
			// (absorbed into SERIAL/IDENTITY). If so, it wasn't truly dropped.
			if seqExistsAsOwned(to, key.schema, key.name) {
				continue
			}
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
			// Check if the sequence exists in `from` as an owned sequence
			// (was SERIAL/IDENTITY, now standalone). If so, it wasn't truly added.
			if seqExistsAsOwned(from, key.schema, key.name) {
				continue
			}
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

// seqExistsAsOwned checks if a sequence with the given name exists in the
// catalog as an owned sequence (OwnerRelOID != 0). This happens when a SERIAL
// or IDENTITY column absorbs the sequence into the table definition.
func seqExistsAsOwned(c *Catalog, schemaName, seqName string) bool {
	s := c.GetSchema(schemaName)
	if s == nil {
		return false
	}
	if seq, ok := s.Sequences[seqName]; ok && seq.OwnerRelOID != 0 {
		return true
	}
	return false
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
