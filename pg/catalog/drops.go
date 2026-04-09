package catalog

import "sort"

// DropsFromDiff returns a MigrationPlan containing only the destructive
// (DROP*) operations that GenerateMigration would produce for the same
// inputs, in the same form, with three differences:
//
//   - SQL is always empty. Callers needing DDL text must call GenerateMigration.
//   - Warning is always empty. GenerateMigration populates warnings on some
//     destructive ops (e.g., table recreates); DropsFromDiff omits them.
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
	ops = append(ops, dropsForExtensions(from, to, diff)...)
	ops = append(ops, dropsForEnums(from, to, diff)...)
	ops = append(ops, dropsForDomains(from, to, diff)...)
	ops = append(ops, dropsForRanges(from, to, diff)...)
	ops = append(ops, dropsForComposites(from, to, diff)...)
	ops = append(ops, dropsForSequences(from, to, diff)...)
	ops = append(ops, dropsForFunctions(from, to, diff)...)
	ops = append(ops, dropsForTables(from, to, diff)...)
	ops = append(ops, dropsForTableRecreates(from, to, diff)...)
	ops = append(ops, dropsForColumns(from, to, diff)...)
	ops = append(ops, dropsForCheckCascades(from, to, diff)...)
	ops = append(ops, dropsForConstraints(from, to, diff)...)
	ops = append(ops, dropsForViews(from, to, diff)...)
	ops = append(ops, dropsForIndexes(from, to, diff)...)
	ops = append(ops, dropsForTriggers(from, to, diff)...)
	ops = append(ops, dropsForPolicies(from, to, diff)...)
	ops = append(ops, dropsForDependentViews(from, to, diff, ops)...)

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

// dropsForExtensions mirrors the DiffDrop arm of generateExtensionDDL in
// migration_extension.go. Extensions are top-level objects: SchemaName is
// empty, ObjectName holds the extension name.
//
// from and to are unused (the OID lives on entry.From) but kept for
// signature uniformity across all dropsForX helpers — see DropsFromDiff.
func dropsForExtensions(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Extensions {
		if entry.Action != DiffDrop {
			continue
		}
		var extOID uint32
		if entry.From != nil {
			extOID = entry.From.OID
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropExtension,
			ObjectName:    entry.Name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'e',
			ObjOID:        extOID,
			Priority:      PriorityExtension,
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

// dropsForTableRecreates mirrors the DiffModify arm of generatePartitionDDL
// in migration_partition.go where a table must be recreated via DROP+CREATE
// because PostgreSQL does not support in-place RelKind changes (e.g.,
// regular → partitioned) or inheritance changes.
//
// Only the OpDropTable half of buildTableRecreateOps is emitted; the
// OpCreateTable and OpCreateIndex halves are not drops.
//
// Both from and to catalogs are required for inhParentsEqual.
func dropsForTableRecreates(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Relations {
		if entry.Action != DiffModify {
			continue
		}
		if entry.From == nil || entry.To == nil {
			continue
		}
		fromIsTable := entry.From.RelKind == 'r' || entry.From.RelKind == 'p'
		toIsTable := entry.To.RelKind == 'r' || entry.To.RelKind == 'p'
		if !fromIsTable || !toIsTable {
			continue
		}

		needsRecreate := false
		if entry.From.RelKind != entry.To.RelKind {
			needsRecreate = true
		}
		if !needsRecreate && !inhParentsEqual(from, to, entry.From.InhParents, entry.To.InhParents) {
			needsRecreate = true
		}
		if !needsRecreate {
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

// dropsForColumns mirrors the DiffDrop arm of generateColumnDDL in
// migration_column.go. Column drops live inside DiffModify relation entries
// (not as top-level diff entries). Views and matviews are skipped because
// column changes there are handled by view DDL.
//
// ObjectName is the TABLE name (not the column name) — this is omni's
// established convention for OpDropColumn.
func dropsForColumns(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		if len(rel.Columns) == 0 {
			continue
		}
		if rel.To != nil && (rel.To.RelKind == 'v' || rel.To.RelKind == 'm') {
			continue
		}

		var relOID uint32
		if rel.To != nil {
			relOID = rel.To.OID
		}

		for _, col := range rel.Columns {
			if col.Action != DiffDrop {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropColumn,
				SchemaName:    rel.SchemaName,
				ObjectName:    rel.Name,
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'r',
				ObjOID:        relOID,
				Priority:      PriorityColumn,
			})
		}
	}
	return ops
}

// dropsForCheckCascades mirrors the CHECK-constraint cascade logic inside
// columnModifyOps in migration_column.go (lines 147-174). When a column's
// type changes, CHECK constraints referencing that column must be dropped
// first (PG requirement). The outer loop in generateColumnDDL overrides
// Phase to PhaseMain and Priority to PriorityColumn.
func dropsForCheckCascades(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		if len(rel.Columns) == 0 {
			continue
		}
		if rel.To != nil && (rel.To.RelKind == 'v' || rel.To.RelKind == 'm') {
			continue
		}

		var relOID uint32
		if rel.To != nil {
			relOID = rel.To.OID
		}

		for _, col := range rel.Columns {
			if col.Action != DiffModify {
				continue
			}
			if col.From == nil || col.To == nil {
				continue
			}
			typeChanged := from.FormatType(col.From.TypeOID, col.From.TypeMod) !=
				to.FormatType(col.To.TypeOID, col.To.TypeMod)
			if !typeChanged {
				continue
			}
			fromRel := from.GetRelation(rel.SchemaName, rel.Name)
			if fromRel == nil {
				continue
			}
			for _, con := range from.ConstraintsOf(fromRel.OID) {
				if con.Type != ConstraintCheck {
					continue
				}
				if !containsColumnRef(con.CheckExpr, col.From.Name) {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropConstraint,
					SchemaName:    rel.SchemaName,
					ObjectName:    rel.Name,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'r',
					ObjOID:        relOID,
					Priority:      PriorityColumn,
				})
			}
		}
	}
	return ops
}

// dropsForConstraints mirrors the DiffDrop and DiffModify arms of
// constraintOpsForRelation in migration_constraint.go. Both DiffDrop and
// DiffModify emit drops (modify = DROP old + ADD new). Constraint-trigger
// types are skipped (handled by trigger DDL).
//
// ObjectName is the CONSTRAINT name; ParentObject is the TABLE name.
// ObjOID is NOT set (zero), matching buildDropConstraintOp.
func dropsForConstraints(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		for _, ce := range rel.Constraints {
			if (ce.To != nil && ce.To.Type == ConstraintTrigger) ||
				(ce.From != nil && ce.From.Type == ConstraintTrigger) {
				continue
			}
			switch ce.Action {
			case DiffDrop:
				if ce.From == nil {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:         OpDropConstraint,
					SchemaName:   rel.SchemaName,
					ObjectName:   ce.From.Name,
					ParentObject: rel.Name,
					Phase:        PhasePre,
					ObjType:      'c',
					Priority:     PriorityConstraint,
				})
			case DiffModify:
				if ce.From == nil {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:         OpDropConstraint,
					SchemaName:   rel.SchemaName,
					ObjectName:   ce.From.Name,
					ParentObject: rel.Name,
					Phase:        PhasePre,
					ObjType:      'c',
					Priority:     PriorityConstraint,
				})
			}
		}
	}
	return ops
}

// dropsForViews mirrors the drop-emission sites in generateViewDDL in
// migration_view.go. There are five sites:
//
//  1. DiffDrop view (RelKind 'v')
//  2. DiffDrop matview (RelKind 'm')
//  3. DiffModify RelKind flip (view↔matview)
//  4. DiffModify regular view with viewColumnsChanged
//  5. DiffModify matview (all modifications — matviews don't support ALTER)
func dropsForViews(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Relations {
		switch entry.Action {
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			if entry.From.RelKind != 'v' && entry.From.RelKind != 'm' {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropView,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'r',
				ObjOID:        entry.From.OID,
				Priority:      PriorityView,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			fromIsView := entry.From.RelKind == 'v' || entry.From.RelKind == 'm'
			toIsView := entry.To.RelKind == 'v' || entry.To.RelKind == 'm'

			// Site 3: RelKind flip (view↔matview).
			if (fromIsView || toIsView) && entry.From.RelKind != entry.To.RelKind {
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityView,
				})
				continue
			}

			switch entry.To.RelKind {
			case 'v':
				// Site 4: view with incompatible column changes.
				if viewColumnsChanged(entry) {
					ops = append(ops, MigrationOp{
						Type:          OpDropView,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'r',
						ObjOID:        entry.From.OID,
						Priority:      PriorityView,
					})
				}
			case 'm':
				// Site 5: ALL matview modifications emit drop.
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityView,
				})
			}
		}
	}
	return ops
}

// dropsForIndexes mirrors the DiffDrop and DiffModify arms of
// generateIndexDDL in migration_index.go. Walks diff.Relations for
// DiffModify entries, then walks relEntry.Indexes. Both DiffDrop and
// DiffModify on indexes emit drops (modified indexes = DROP old + CREATE
// new). Indexes where ConstraintOID != 0 are skipped (managed by
// constraint DDL).
//
// from and to are unused but kept for signature uniformity — see DropsFromDiff.
func dropsForIndexes(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, relEntry := range diff.Relations {
		if relEntry.Action != DiffModify {
			continue
		}
		for _, idxEntry := range relEntry.Indexes {
			switch idxEntry.Action {
			case DiffDrop:
				idx := idxEntry.From
				if idx.ConstraintOID != 0 {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropIndex,
					SchemaName:    relEntry.SchemaName,
					ObjectName:    idx.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'i',
					ObjOID:        idx.OID,
					Priority:      PriorityIndex,
				})
			case DiffModify:
				idxFrom := idxEntry.From
				idxTo := idxEntry.To
				if idxFrom.ConstraintOID != 0 || idxTo.ConstraintOID != 0 {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropIndex,
					SchemaName:    relEntry.SchemaName,
					ObjectName:    idxFrom.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'i',
					ObjOID:        idxFrom.OID,
					Priority:      PriorityIndex,
				})
			}
		}
	}
	return ops
}

// dropsForTriggers mirrors the DiffDrop and DiffModify arms of
// triggerOpsForRelation in migration_trigger.go. Both DiffDrop and DiffModify
// emit drops (modify = DROP old + CREATE new).
//
// ObjectName is the TRIGGER name; ParentObject is the TABLE name.
// ObjType is 'T' (capital T).
func dropsForTriggers(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		for _, te := range rel.Triggers {
			switch te.Action {
			case DiffDrop:
				if te.From == nil {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropTrigger,
					SchemaName:    rel.SchemaName,
					ObjectName:    te.From.Name,
					ParentObject:  rel.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'T',
					ObjOID:        te.From.OID,
					Priority:      PriorityTrigger,
				})
			case DiffModify:
				if te.From == nil {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropTrigger,
					SchemaName:    rel.SchemaName,
					ObjectName:    te.From.Name,
					ParentObject:  rel.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'T',
					ObjOID:        te.From.OID,
					Priority:      PriorityTrigger,
				})
			}
		}
	}
	return ops
}

// dropsForPolicies mirrors the DiffDrop and DiffModify arms of
// policyOpsForRelation in migration_policy.go. Walks diff.Relations for
// DiffModify entries, then walks rel.Policies. Emits drops for DiffDrop
// (always) and DiffModify (only when CmdType or Permissive changed,
// which forces DROP + CREATE).
//
// ObjectName is the POLICY name; ParentObject is the TABLE name (unqualified).
// from and to are unused but kept for signature uniformity — see DropsFromDiff.
func dropsForPolicies(_, _ *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		for _, pe := range rel.Policies {
			switch pe.Action {
			case DiffDrop:
				if pe.From == nil {
					continue
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropPolicy,
					SchemaName:    rel.SchemaName,
					ObjectName:    pe.From.Name,
					ParentObject:  rel.Name,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'p',
					ObjOID:        pe.From.OID,
					Priority:      PriorityPolicy,
				})
			case DiffModify:
				if pe.From == nil || pe.To == nil {
					continue
				}
				if pe.From.CmdType != pe.To.CmdType || pe.From.Permissive != pe.To.Permissive {
					ops = append(ops, MigrationOp{
						Type:          OpDropPolicy,
						SchemaName:    rel.SchemaName,
						ObjectName:    pe.From.Name,
						ParentObject:  rel.Name,
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'p',
						ObjOID:        pe.From.OID,
						Priority:      PriorityPolicy,
					})
				}
			}
		}
	}
	return ops
}

// dropsForDependentViews mirrors wrapColumnTypeChangesWithViewOps in
// migration.go:241-401. When a table's columns are dropped, retyped, or
// the table undergoes a RelKind/inheritance change, PostgreSQL requires
// dependent views to be dropped first. This helper walks from.deps via
// BFS to find all transitively dependent views and emits OpDropView ops
// for each, deduplicating against existingOps so views already handled
// by dropsForViews are not double-counted.
//
// Only regular views (RelKind 'v') are considered — matviews have their
// own handling in buildModifyMatViewOps.
func dropsForDependentViews(from, to *Catalog, diff *SchemaDiff, existingOps []MigrationOp) []MigrationOp {
	// Step 1: Find OIDs of tables that need dependent-view wrapping.
	tableOIDs := make(map[uint32]bool)
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		needsViewWrap := false

		if rel.From != nil && rel.To != nil {
			if rel.From.RelKind != rel.To.RelKind {
				needsViewWrap = true
			}
			if !inhParentsEqual(from, to, rel.From.InhParents, rel.To.InhParents) {
				needsViewWrap = true
			}
		}

		if !needsViewWrap {
			for _, col := range rel.Columns {
				if col.Action == DiffDrop {
					needsViewWrap = true
					break
				}
				if col.Action != DiffModify || col.From == nil || col.To == nil {
					continue
				}
				if from.FormatType(col.From.TypeOID, col.From.TypeMod) != to.FormatType(col.To.TypeOID, col.To.TypeMod) {
					needsViewWrap = true
					break
				}
			}
		}
		if needsViewWrap {
			r := from.GetRelation(rel.SchemaName, rel.Name)
			if r != nil {
				tableOIDs[r.OID] = true
			}
		}
	}
	if len(tableOIDs) == 0 {
		return nil
	}

	// Step 2: BFS over from.deps to find dependent views transitively.
	type viewInfo struct {
		schema string
		name   string
		oid    uint32
	}
	seen := make(map[uint32]bool)
	var viewsToDrop []viewInfo

	queue := make([]uint32, 0, len(tableOIDs))
	for oid := range tableOIDs {
		queue = append(queue, oid)
	}
	for len(queue) > 0 {
		refOID := queue[0]
		queue = queue[1:]
		for _, d := range from.deps {
			if d.RefType != 'r' || d.RefOID != refOID || d.ObjType != 'r' {
				continue
			}
			if seen[d.ObjOID] {
				continue
			}
			rel := from.GetRelationByOID(d.ObjOID)
			if rel == nil || rel.RelKind != 'v' {
				continue
			}
			seen[d.ObjOID] = true
			if rel.Schema == nil {
				continue
			}
			viewsToDrop = append(viewsToDrop, viewInfo{schema: rel.Schema.Name, name: rel.Name, oid: rel.OID})
			queue = append(queue, d.ObjOID)
		}
	}

	if len(viewsToDrop) == 0 {
		return nil
	}

	// Step 3: Sort for determinism.
	sort.Slice(viewsToDrop, func(i, j int) bool {
		if viewsToDrop[i].schema != viewsToDrop[j].schema {
			return viewsToDrop[i].schema < viewsToDrop[j].schema
		}
		return viewsToDrop[i].name < viewsToDrop[j].name
	})

	// Step 4: Build dedup set from existing ops.
	existing := make(map[string]bool)
	for _, op := range existingOps {
		if op.Type == OpDropView || op.Type == OpCreateView || op.Type == OpAlterView {
			existing[op.SchemaName+"."+op.ObjectName] = true
		}
	}

	// Step 5: Emit OpDropView for each dependent view not already covered.
	var ops []MigrationOp
	for _, v := range viewsToDrop {
		key := v.schema + "." + v.name
		if existing[key] {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpDropView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        v.oid,
			Priority:      PriorityView,
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
