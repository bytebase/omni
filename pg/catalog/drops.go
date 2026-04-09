package catalog

import "sort"

// DropsFromDiff returns a MigrationPlan containing only the destructive
// (DROP*) operations that GenerateMigration would produce for the same
// inputs, in the same form, with two differences:
//
//   - SQL is always empty. Callers needing DDL text must call GenerateMigration.
//   - Operations are sorted by a deterministic lexicographic key
//     (SchemaName, ObjectName, ParentObject, Type) — see sortDropOps — not
//     topologically.
//
// This is a fast path for callers (e.g., destructive-change advisories) that
// need to know what would be dropped without paying the cost of building DDL
// text or running the dependency-driven topological sort.
//
// All metadata fields (Type, SchemaName, ObjectName, ParentObject, Phase,
// ObjType, ObjOID, Priority) are populated identically to what
// GenerateMigration sets. ParentObject follows omni's existing convention:
// for OpDropConstraint and OpDropTrigger it holds the unqualified parent
// table name; for OpDropColumn the table name is in ObjectName (not
// ParentObject) per migration_column.go.
func DropsFromDiff(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	if diff == nil {
		return &MigrationPlan{}
	}
	var ops []MigrationOp
	// Producers are appended in the same order as GenerateMigration.
	// Each dropsForX helper documents which migration_*.go function it mirrors.
	// All helpers take the same (from, to, diff) signature for call-site
	// uniformity, even when a particular helper does not need from or to.
	ops = append(ops, dropsForSchemas(from, to, diff)...)
	ops = append(ops, dropsForEnums(from, to, diff)...)
	ops = append(ops, dropsForDomains(from, to, diff)...)
	ops = append(ops, dropsForRanges(from, to, diff)...)
	ops = append(ops, dropsForComposites(from, to, diff)...)
	ops = append(ops, dropsForSequences(from, to, diff)...)
	ops = append(ops, dropsForFunctions(from, to, diff)...)
	ops = append(ops, dropsForTables(from, to, diff)...)

	sortDropOps(ops)
	return &MigrationPlan{Ops: ops}
}

// dropsForSchemas mirrors the DiffDrop arm of generateSchemaDDL in
// migration_schema.go. Schemas are top-level objects: SchemaName is empty,
// ObjectName holds the schema name.
//
// from and to are unused (the OID lives on entry.From) but kept for
// signature uniformity across all dropsForX helpers — see DropsFromDiff.
func dropsForSchemas(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Schemas {
		if entry.Action != DiffDrop {
			continue
		}
		var schemaOID uint32
		if entry.From != nil {
			schemaOID = entry.From.OID
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropSchema,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'n',
			ObjOID:        schemaOID,
			Priority:      PrioritySchema,
		})
	}
	return ops
}

// dropsForEnums mirrors the DiffDrop arm of generateEnumDDL in
// migration_enum.go. The OID is resolved against the from catalog by name
// (resolveTypeOIDByName), so this helper genuinely uses from; to is unused
// but kept for signature uniformity — see DropsFromDiff.
func dropsForEnums(from, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Enums {
		if entry.Action != DiffDrop {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       't',
			ObjOID:        resolveTypeOIDByName(from, entry.SchemaName, entry.Name),
			Priority:      PriorityType,
		})
	}
	return ops
}

// dropsForDomains mirrors the DiffDrop arm of generateDomainDDL in
// migration_domain.go. The OID lives on entry.From.TypeOID (NOT entry.From.OID),
// guarded by an entry.From nil check. from and to are unused but kept for
// signature uniformity — see DropsFromDiff.
func dropsForDomains(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Domains {
		if entry.Action != DiffDrop {
			continue
		}
		var typeOID uint32
		if entry.From != nil {
			typeOID = entry.From.TypeOID
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       't',
			ObjOID:        typeOID,
			Priority:      PriorityType,
		})
	}
	return ops
}

// dropsForRanges mirrors the DiffDrop arm of generateRangeDDL in
// migration_range.go. Skips entries with nil From. from and to are unused
// but kept for signature uniformity — see DropsFromDiff.
func dropsForRanges(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Ranges {
		if entry.Action != DiffDrop {
			continue
		}
		if entry.From == nil {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       't',
			ObjOID:        entry.From.OID,
			Priority:      PriorityType,
		})
	}
	return ops
}

// dropsForComposites mirrors the DiffDrop arm of generateCompositeDDL in
// migration_composite.go. Note: ObjType is 'r' (relation), not 't' — composite
// types are stored as relations in PostgreSQL. Skips entries with nil From.
// from and to are unused but kept for signature uniformity — see DropsFromDiff.
func dropsForComposites(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.CompositeTypes {
		if entry.Action != DiffDrop {
			continue
		}
		if entry.From == nil {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        entry.From.OID,
			Priority:      PriorityType,
		})
	}
	return ops
}

// dropsForSequences mirrors the DiffDrop arm of generateSequenceDDL in
// migration_sequence.go. Skips identity-backed sequences (DROP IDENTITY on
// the owning column drops them automatically) — this filter requires the
// from catalog. to is unused but kept for signature uniformity — see
// DropsFromDiff.
func dropsForSequences(from, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Sequences {
		if entry.Action != DiffDrop {
			continue
		}
		if entry.From == nil {
			continue
		}
		if isIdentitySequence(from, entry.SchemaName, entry.Name) {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropSequence,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'S',
			ObjOID:        entry.From.OID,
			Priority:      PrioritySequence,
		})
	}
	return ops
}

// dropsForFunctions mirrors the DiffDrop arm of generateFunctionDDL in
// migration_function.go, plus the DiffModify arm where signatureChanged
// returns true (signature changes force DROP + CREATE; body-only changes
// use CREATE OR REPLACE and emit no drop op).
//
// ObjectName is set to entry.Identity (e.g. "add_one(integer)"), not
// entry.Name, so overloaded functions remain distinguishable. Both from
// and to catalogs are required because the DiffModify path must call
// signatureChanged(from, entry.From, to, entry.To).
func dropsForFunctions(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Functions {
		switch entry.Action {
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, dropFunctionOp(entry))
		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			if signatureChanged(from, entry.From, to, entry.To) {
				ops = append(ops, dropFunctionOp(entry))
			}
		}
	}
	return ops
}

// dropFunctionOp builds the OpDropFunction MigrationOp for a function diff
// entry, mirroring buildDropFunctionOp in migration_function.go minus SQL.
func dropFunctionOp(entry FunctionDiffEntry) MigrationOp {
	return MigrationOp{
		Type:          OpDropFunction,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Identity,
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'f',
		ObjOID:        entry.From.OID,
		Priority:      PriorityFunction,
	}
}

// dropsForTables mirrors the DiffDrop arm of generateTableDDL in
// migration_table.go. Only relations with RelKind 'r' (regular table)
// or 'p' (partitioned table) qualify; views, matviews, and composite
// types are handled by their own helpers (dropsForViews,
// dropsForComposites). The DiffModify recreate cases (RelKind flip,
// inheritance change) are handled by dropsForTableRecreates.
//
// from and to are unused (the OID lives on entry.From) but kept for
// signature uniformity — see DropsFromDiff.
func dropsForTables(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Relations {
		if entry.Action != DiffDrop {
			continue
		}
		if entry.From == nil {
			continue
		}
		rk := entry.From.RelKind
		if rk != 'r' && rk != 'p' {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropTable,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        entry.From.OID,
			Priority:      PriorityTable,
		})
	}
	return ops
}

// sortDropOps gives DropsFromDiff output deterministic ordering. We do NOT
// run sortMigrationOps because that performs an expensive topological sort
// over the dep graph that drop advice does not need.
//
// The sort key includes ParentObject because OpDropConstraint and
// OpDropTrigger ops carry the constraint/trigger name in ObjectName and the
// parent table name in ParentObject — without ParentObject in the key, two
// drops of the same-named constraint on different tables would be ordered
// nondeterministically.
func sortDropOps(ops []MigrationOp) {
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		if ops[i].ObjectName != ops[j].ObjectName {
			return ops[i].ObjectName < ops[j].ObjectName
		}
		if ops[i].ParentObject != ops[j].ParentObject {
			return ops[i].ParentObject < ops[j].ParentObject
		}
		return ops[i].Type < ops[j].Type
	})
}
