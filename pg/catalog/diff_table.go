package catalog

import "sort"

// relKey is the identity key for a relation: qualified name.
type relKey struct {
	schema string
	name   string
}

// diffRelations compares relations (tables only, RelKind='r') between two catalogs.
// Views and materialized views are handled separately.
func diffRelations(from, to *Catalog) []RelationDiffEntry {
	// Build maps of (schema, name) -> *Relation for both catalogs.
	fromMap := buildRelationMap(from)
	toMap := buildRelationMap(to)

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
			// Added.
			result = append(result, RelationDiffEntry{
				Action:     DiffAdd,
				SchemaName: key.schema,
				Name:       key.name,
				To:         toRel,
			})
			continue
		}

		// Both exist — check for modifications.
		if entry, changed := compareRelation(key, fromRel, toRel); changed {
			result = append(result, entry)
		}
	}

	// Sort for determinism: by schema name, then relation name, then action.
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

// buildRelationMap builds a map of (schema, name) -> *Relation for tables (RelKind='r')
// and partitioned tables (RelKind='p').
func buildRelationMap(c *Catalog) map[relKey]*Relation {
	m := make(map[relKey]*Relation)
	for _, s := range c.UserSchemas() {
		for _, rel := range s.Relations {
			if rel.RelKind != 'r' && rel.RelKind != 'p' {
				continue
			}
			m[relKey{schema: s.Name, name: rel.Name}] = rel
		}
	}
	return m
}

// compareRelation checks whether two relations with the same identity differ
// in table-level properties. Column/constraint/index sub-diffs are not yet implemented.
func compareRelation(key relKey, from, to *Relation) (RelationDiffEntry, bool) {
	changed := false

	if from.Persistence != to.Persistence {
		changed = true
	}
	if from.ReplicaIdentity != to.ReplicaIdentity {
		changed = true
	}
	if from.RowSecurity != to.RowSecurity {
		changed = true
	}
	if from.ForceRowSecurity != to.ForceRowSecurity {
		changed = true
	}
	if from.Owner != to.Owner {
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
	}, true
}
