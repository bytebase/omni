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
		if entry, changed := compareRelation(from, to, key, fromRel, toRel); changed {
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
// in table-level properties or sub-objects (columns).
func compareRelation(fromCat, toCat *Catalog, key relKey, from, to *Relation) (RelationDiffEntry, bool) {
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
	if from.RelKind != to.RelKind {
		changed = true
	}

	// Partition info comparison.
	if !partitionInfoEqual(from.PartitionInfo, to.PartitionInfo) {
		changed = true
	}

	// Partition bound comparison.
	if !partitionBoundEqual(from.PartitionBound, to.PartitionBound) {
		changed = true
	}

	// Inheritance parents comparison.
	if !inhParentsEqual(fromCat, toCat, from.InhParents, to.InhParents) {
		changed = true
	}

	// Column sub-diff.
	cols := diffColumns(fromCat, toCat, from, to)
	if len(cols) > 0 {
		changed = true
	}

	// Constraint sub-diff.
	conDiffs := diffConstraints(fromCat, toCat, from.OID, to.OID)
	if len(conDiffs) > 0 {
		changed = true
	}

	// Index sub-diff.
	idxDiffs := diffIndexes(fromCat, toCat, from.OID, to.OID)
	if len(idxDiffs) > 0 {
		changed = true
	}

	// Trigger sub-diff.
	trigDiffs := diffTriggers(fromCat, toCat, from.OID, to.OID)
	if len(trigDiffs) > 0 {
		changed = true
	}

	// Policy sub-diff.
	polDiffs := diffPolicies(fromCat, toCat, from.OID, to.OID)
	if len(polDiffs) > 0 {
		changed = true
	}

	// RLS flag changes.
	var rlsChanged bool
	if from.RowSecurity != to.RowSecurity || from.ForceRowSecurity != to.ForceRowSecurity {
		rlsChanged = true
		changed = true
	}

	if !changed {
		return RelationDiffEntry{}, false
	}

	entry := RelationDiffEntry{
		Action:      DiffModify,
		SchemaName:  key.schema,
		Name:        key.name,
		From:        from,
		To:          to,
		Columns:     cols,
		Constraints: conDiffs,
		Indexes:     idxDiffs,
		Triggers:    trigDiffs,
		Policies:    polDiffs,
	}
	if rlsChanged {
		entry.RLSChanged = true
		entry.RLSEnabled = to.RowSecurity
		entry.ForceRLSEnabled = to.ForceRowSecurity
	}
	return entry, true
}

// partitionInfoEqual returns true if two PartitionInfo values are equivalent.
func partitionInfoEqual(a, b *PartitionInfo) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Strategy != b.Strategy || a.NKeyAttrs != b.NKeyAttrs {
		return false
	}
	if len(a.KeyAttNums) != len(b.KeyAttNums) {
		return false
	}
	for i := range a.KeyAttNums {
		if a.KeyAttNums[i] != b.KeyAttNums[i] {
			return false
		}
	}
	return true
}

// partitionBoundEqual returns true if two PartitionBound values are equivalent.
func partitionBoundEqual(a, b *PartitionBound) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Strategy != b.Strategy || a.IsDefault != b.IsDefault {
		return false
	}
	if !stringSliceEqual(a.ListValues, b.ListValues) {
		return false
	}
	if !stringSliceEqual(a.LowerBound, b.LowerBound) {
		return false
	}
	if !stringSliceEqual(a.UpperBound, b.UpperBound) {
		return false
	}
	if a.Modulus != b.Modulus || a.Remainder != b.Remainder {
		return false
	}
	return true
}

// inhParentsEqual compares inheritance parent lists by resolving OIDs to names.
func inhParentsEqual(fromCat, toCat *Catalog, fromOIDs, toOIDs []uint32) bool {
	fromNames := resolveRelNames(fromCat, fromOIDs)
	toNames := resolveRelNames(toCat, toOIDs)
	return stringSliceEqual(fromNames, toNames)
}

// resolveRelNames resolves a list of relation OIDs to qualified names.
func resolveRelNames(c *Catalog, oids []uint32) []string {
	names := make([]string, 0, len(oids))
	for _, oid := range oids {
		if rel := c.GetRelationByOID(oid); rel != nil && rel.Schema != nil {
			names = append(names, rel.Schema.Name+"."+rel.Name)
		}
	}
	return names
}
