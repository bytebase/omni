package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateViewDDL produces CREATE VIEW, DROP VIEW, CREATE MATERIALIZED VIEW,
// DROP MATERIALIZED VIEW, and related modification operations from the diff.
func generateViewDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Relations {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			switch entry.To.RelKind {
			case 'v':
				ops = append(ops, buildCreateViewOp(to, entry))
			case 'm':
				ops = append(ops, buildCreateMatViewOp(to, entry))
			}

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			switch entry.From.RelKind {
			case 'v':
				qn := migrationQualifiedName(entry.SchemaName, entry.Name)
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("DROP VIEW %s", qn),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityView,
				})
			case 'm':
				qn := migrationQualifiedName(entry.SchemaName, entry.Name)
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("DROP MATERIALIZED VIEW %s", qn),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityView,
				})
			}

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}

			// RelKind changed between view types (e.g., view → matview or vice versa): DROP old + CREATE new.
			// Only handle view-type RelKinds; table changes are handled by generatePartitionDDL.
			fromIsView := entry.From.RelKind == 'v' || entry.From.RelKind == 'm'
			toIsView := entry.To.RelKind == 'v' || entry.To.RelKind == 'm'
			if (fromIsView || toIsView) && entry.From.RelKind != entry.To.RelKind {
				qn := migrationQualifiedName(entry.SchemaName, entry.Name)
				// Drop the original kind.
				dropSQL := fmt.Sprintf("DROP VIEW %s", qn)
				if entry.From.RelKind == 'm' {
					dropSQL = fmt.Sprintf("DROP MATERIALIZED VIEW %s", qn)
				}
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           dropSQL,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        entry.From.OID,
					Priority:      PriorityView,
				})
				// Create the new kind.
				switch entry.To.RelKind {
				case 'v':
					ops = append(ops, buildCreateViewOp(to, entry))
				case 'm':
					ops = append(ops, buildCreateMatViewOp(to, entry))
				}
				continue
			}

			switch entry.To.RelKind {
			case 'v':
				if viewColumnsChanged(entry) {
					// Column order/names/types changed — can't use CREATE OR REPLACE.
					qn := migrationQualifiedName(entry.SchemaName, entry.Name)
					ops = append(ops, MigrationOp{
						Type:          OpDropView,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("DROP VIEW %s", qn),
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'r',
						ObjOID:        entry.From.OID,
						Priority:      PriorityView,
					})
					ops = append(ops, buildCreateViewOp(to, entry))
				} else {
					ops = append(ops, buildModifyViewOps(from, to, entry)...)
				}
			case 'm':
				ops = append(ops, buildModifyMatViewOps(to, entry)...)
			}
		}
	}

	// Deterministic ordering: drops first, then creates/replaces, by schema + name.
	sort.SliceStable(ops, func(i, j int) bool {
		iDrop := ops[i].Type == OpDropView
		jDrop := ops[j].Type == OpDropView
		if iDrop != jDrop {
			return iDrop
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// buildCreateViewOp creates a CREATE VIEW op for an added view.
func buildCreateViewOp(to *Catalog, entry RelationDiffEntry) MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetViewDefinition(entry.SchemaName, entry.Name)

	var b strings.Builder
	b.WriteString("CREATE VIEW ")
	b.WriteString(qn)
	b.WriteString(" AS ")
	b.WriteString(strings.TrimRight(def, " \t\n\r;"))

	// Check option.
	switch entry.To.CheckOption {
	case 'l':
		b.WriteString("\n WITH LOCAL CHECK OPTION")
	case 'c':
		b.WriteString("\n WITH CASCADED CHECK OPTION")
	}

	return MigrationOp{
		Type:          OpCreateView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityView,
	}
}

// buildCreateMatViewOp creates a CREATE MATERIALIZED VIEW op.
func buildCreateMatViewOp(to *Catalog, entry RelationDiffEntry) MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetMatViewDefinition(entry.SchemaName, entry.Name)

	// The view's storage travels with the statement that names it: the USING
	// access method and the option list, before AS.
	options := ""
	if rel := to.GetRelation(entry.SchemaName, entry.Name); rel != nil {
		options = deparseMatViewOptions(rel)
	}

	sql := fmt.Sprintf("CREATE MATERIALIZED VIEW %s%s AS %s",
		qn, options, strings.TrimRight(def, " \t\n\r;"))

	return MigrationOp{
		Type:          OpCreateView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           sql,
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityView,
	}
}

// buildModifyViewOps creates a CREATE OR REPLACE VIEW op for a modified view.
// Also generates warnings for dependent views that may need recreation.
func buildModifyViewOps(from, to *Catalog, entry RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp

	// Check for dependent views and generate warnings.
	dependents := findDependentViews(from, entry.SchemaName, entry.Name)
	if len(dependents) > 0 {
		names := make([]string, len(dependents))
		for i, d := range dependents {
			names[i] = migrationQualifiedName(d.schema, d.name)
		}
		sort.Strings(names)

		ops = append(ops, MigrationOp{
			Type:          OpAlterView,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           "", // no SQL — warning only
			Warning:       fmt.Sprintf("view %s has dependent views that may need recreation: %s", migrationQualifiedName(entry.SchemaName, entry.Name), strings.Join(names, ", ")),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'r',
			ObjOID:        entry.To.OID,
			Priority:      PriorityView,
		})
	}

	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetViewDefinition(entry.SchemaName, entry.Name)

	var b strings.Builder
	b.WriteString("CREATE OR REPLACE VIEW ")
	b.WriteString(qn)
	b.WriteString(" AS ")
	b.WriteString(strings.TrimRight(def, " \t\n\r;"))

	switch entry.To.CheckOption {
	case 'l':
		b.WriteString("\n WITH LOCAL CHECK OPTION")
	case 'c':
		b.WriteString("\n WITH CASCADED CHECK OPTION")
	}

	ops = append(ops, MigrationOp{
		Type:          OpAlterView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           b.String(),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityView,
	})

	return ops
}

// buildModifyMatViewOps creates DROP + CREATE ops for a modified matview.
// Materialized views don't support CREATE OR REPLACE. The views and
// materialized views that read it are handled beside this by
// wrapMatViewReplacementsWithDependents, which drops them first and creates
// them again from the target.
func buildModifyMatViewOps(to *Catalog, entry RelationDiffEntry) []MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetMatViewDefinition(entry.SchemaName, entry.Name)

	// The view's storage travels with the statement that names it, as it does
	// when the view is added.
	options := ""
	if rel := to.GetRelation(entry.SchemaName, entry.Name); rel != nil {
		options = deparseMatViewOptions(rel)
	}

	dropOp := MigrationOp{
		Type:          OpDropView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           fmt.Sprintf("DROP MATERIALIZED VIEW %s", qn),
		Transactional: true,
		Phase:         PhasePre,
		ObjType:       'r',
		ObjOID:        entry.From.OID,
		Priority:      PriorityView,
	}

	createSQL := fmt.Sprintf("CREATE MATERIALIZED VIEW %s%s AS %s",
		qn, options, strings.TrimRight(def, " \t\n\r;"))

	createOp := MigrationOp{
		Type:          OpCreateView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           createSQL,
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        entry.To.OID,
		Priority:      PriorityView,
	}

	return []MigrationOp{dropOp, createOp}
}

// wrapMatViewReplacementsWithDependents adds what a replaced materialized view
// needs: the relations that read it, transitively — views and materialized
// views alike — dropped before it and created again from the target, with the
// grants, comments and indexes their drops take with them. Redshift has no ALTER
// for most of a materialized view's storage, so the catalog re-creates the view;
// the replacement's own ops for those relations are filtered out, since this
// function writes them.
//
// The drops carry the from catalog's OIDs and the creations the target's, which
// is what sortMigrationOps orders each phase by.
func wrapMatViewReplacementsWithDependents(from, to *Catalog, diff *SchemaDiff, ops []MigrationOp) []MigrationOp {
	var replaced []*Relation
	for _, entry := range diff.Relations {
		if entry.Action == DiffModify && entry.To != nil && entry.To.RelKind == 'm' && entry.From != nil {
			replaced = append(replaced, entry.From)
		}
	}
	if len(replaced) == 0 {
		return ops
	}

	// The dependents, transitively, by the catalog's recorded dependencies: a
	// name in a definition would take a view reading a column such as
	// customer_id for one reading a relation named id.
	seen := make(map[uint32]bool)
	var dependents []*Relation
	queue := make([]uint32, 0, len(replaced))
	for _, rel := range replaced {
		queue = append(queue, rel.OID)
	}
	for len(queue) > 0 {
		oid := queue[0]
		queue = queue[1:]
		for _, dep := range from.deps {
			if dep.RefType != 'r' || dep.RefOID != oid || dep.ObjType != 'r' || seen[dep.ObjOID] {
				continue
			}
			rel := from.GetRelationByOID(dep.ObjOID)
			if rel == nil || rel.Schema == nil || (rel.RelKind != 'v' && rel.RelKind != 'm') {
				continue
			}
			seen[dep.ObjOID] = true
			dependents = append(dependents, rel)
			queue = append(queue, rel.OID)
		}
	}
	if len(dependents) == 0 {
		// The view itself is still replaced: its own grants and comments went
		// with the drop, and their replacements are added below.
		var extra []MigrationOp
		for _, rel := range replaced {
			target := to.GetRelation(rel.Schema.Name, rel.Name)
			if target == nil {
				continue
			}
			extra = append(extra, buildRelationMetadataOps(to, target)...)
		}
		return append(ops, extra...)
	}

	// Whatever the diff wrote for those relations is this function's, grants,
	// comments and indexes included. An index op names the index, not the
	// relation, so the indexes of each dependent are matched by OID.
	owned := make(map[string]bool, len(dependents))
	ownedIndexes := make(map[uint32]bool)
	for _, rel := range dependents {
		owned[rel.Schema.Name+"."+rel.Name] = true
		for _, idx := range from.IndexesOf(rel.OID) {
			ownedIndexes[idx.OID] = true
		}
		// An index the target adds is created after the relation this chain
		// creates, not beside it.
		if target := to.GetRelation(rel.Schema.Name, rel.Name); target != nil {
			for _, idx := range to.IndexesOf(target.OID) {
				ownedIndexes[idx.OID] = true
			}
		}
	}
	kept := ops[:0]
	for _, op := range ops {
		if owned[op.SchemaName+"."+op.ObjectName] {
			switch op.Type {
			case OpDropView, OpCreateView, OpAlterView, OpGrant, OpRevoke, OpComment:
				continue
			}
		}
		if (op.Type == OpCreateIndex || op.Type == OpDropIndex) && ownedIndexes[op.ObjOID] {
			continue
		}
		// A trigger op names the trigger and holds the relation it belongs to.
		if owned[op.SchemaName+"."+op.ParentObject] {
			switch op.Type {
			case OpCreateTrigger, OpDropTrigger:
				continue
			}
		}
		kept = append(kept, op)
	}

	var extra []MigrationOp
	for _, rel := range dependents {
		kind := "VIEW"
		if rel.RelKind == 'm' {
			kind = "MATERIALIZED VIEW"
		}
		extra = append(extra, MigrationOp{
			Type:          OpDropView,
			SchemaName:    rel.Schema.Name,
			ObjectName:    rel.Name,
			SQL:           fmt.Sprintf("DROP %s %s", kind, migrationQualifiedName(rel.Schema.Name, rel.Name)),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        rel.OID,
			Priority:      PriorityView,
		})
	}

	for _, rel := range dependents {
		target := to.GetRelation(rel.Schema.Name, rel.Name)
		if target == nil || (target.RelKind != 'v' && target.RelKind != 'm') {
			// The target drops it, or holds a table of that name now: this
			// function's drop above is its whole story, and the table's own
			// diff says the rest.
			continue
		}
		extra = append(extra, buildRecreateRelationOps(to, target)...)
	}

	// The replaced views' own grants and comments went with their drops.
	for _, rel := range replaced {
		// By name: the two catalogs allocate OIDs independently, so a source
		// OID names no relation in the target.
		target := to.GetRelation(rel.Schema.Name, rel.Name)
		if target == nil {
			continue
		}
		extra = append(extra, buildRelationMetadataOps(to, target)...)
	}

	return append(kept, extra...)
}

// buildRecreateRelationOps creates a view or materialized view again from the
// target, with the metadata its drop took with it.
func buildRecreateRelationOps(to *Catalog, rel *Relation) []MigrationOp {
	qn := migrationQualifiedName(rel.Schema.Name, rel.Name)
	var sql string
	if rel.RelKind == 'm' {
		def, _ := to.GetMatViewDefinition(rel.Schema.Name, rel.Name)
		sql = fmt.Sprintf("CREATE MATERIALIZED VIEW %s%s AS %s",
			qn, deparseMatViewOptions(rel), strings.TrimRight(def, " \t\n\r;"))
	} else {
		def, _ := to.GetViewDefinition(rel.Schema.Name, rel.Name)
		var b strings.Builder
		b.WriteString("CREATE VIEW ")
		b.WriteString(qn)
		b.WriteString(" AS ")
		b.WriteString(strings.TrimRight(def, " \t\n\r;"))
		switch rel.CheckOption {
		case 'l':
			b.WriteString("\n WITH LOCAL CHECK OPTION")
		case 'c':
			b.WriteString("\n WITH CASCADED CHECK OPTION")
		}
		sql = b.String()
	}

	ops := []MigrationOp{{
		Type:          OpCreateView,
		SchemaName:    rel.Schema.Name,
		ObjectName:    rel.Name,
		SQL:           sql,
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		ObjOID:        rel.OID,
		Priority:      PriorityView,
	}}

	// The triggers a dropped view had, and a materialized view's indexes, went
	// with it too.
	for _, trig := range to.TriggersOf(rel.OID) {
		ops = append(ops, buildCreateTriggerOp(to, rel.Schema.Name, rel.Name, trig))
	}
	if rel.RelKind == 'm' {
		for _, idx := range to.IndexesOf(rel.OID) {
			if idx.ConstraintOID != 0 {
				continue // managed by constraint DDL
			}
			ops = append(ops, MigrationOp{
				Type:          OpCreateIndex,
				SchemaName:    rel.Schema.Name,
				ObjectName:    idx.Name,
				SQL:           buildCreateIndexSQL(idx, rel, rel.Schema.Name),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'i',
				ObjOID:        idx.OID,
				Priority:      PriorityIndex,
			})
		}
	}

	return append(ops, buildRelationMetadataOps(to, rel)...)
}

// buildRelationMetadataOps puts back what a relation's grants and comments were,
// which an unchanged grant or comment has no diff entry of its own to say.
func buildRelationMetadataOps(c *Catalog, rel *Relation) []MigrationOp {
	desc := rel.Schema.Name + "." + rel.Name
	var ops []MigrationOp
	for _, g := range c.grants {
		if g.ObjType != 'r' || g.ObjOID != rel.OID {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpGrant,
			SchemaName:    rel.Schema.Name,
			ObjectName:    rel.Name,
			SQL:           formatGrantSQL(g, desc),
			Transactional: grantTransactional(g),
			Phase:         PhaseMain,
			ObjType:       g.ObjType,
			ObjOID:        g.ObjOID,
			Priority:      PriorityMetadata,
		})
	}

	var commentKeys []commentKey
	for key, text := range c.comments {
		if key.ObjType != 'r' || key.ObjOID != rel.OID || text == "" {
			continue
		}
		commentKeys = append(commentKeys, key)
	}
	sort.Slice(commentKeys, func(i, j int) bool { return commentKeys[i].SubID < commentKeys[j].SubID })
	for _, key := range commentKeys {
		sql := formatCommentSQL(c, 'r', desc, key.SubID, c.comments[key])
		if sql == "" {
			continue
		}
		ops = append(ops, MigrationOp{
			Type:          OpComment,
			SchemaName:    rel.Schema.Name,
			ObjectName:    rel.Name,
			SQL:           sql,
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'r',
			Priority:      PriorityMetadata,
		})
	}
	return ops
}

// viewRef identifies a view by schema and name.
type viewRef struct {
	schema string
	name   string
}

// findDependentViews finds views that depend on the given view via the catalog's
// dependency tracking. For now, uses a simple heuristic: scan all views in the
// catalog and check if their definition references this view.
func findDependentViews(c *Catalog, schemaName, viewName string) []viewRef {
	var deps []viewRef

	for _, s := range c.UserSchemas() {
		for _, rel := range s.Relations {
			if rel.RelKind != 'v' {
				continue
			}
			// Skip self.
			if s.Name == schemaName && rel.Name == viewName {
				continue
			}
			// Check if this view references the target view.
			def, err := c.GetViewDefinition(s.Name, rel.Name)
			if err != nil {
				continue
			}
			// Simple heuristic: check if the view definition references the target.
			if strings.Contains(def, quoteIdentAlways(viewName)) || strings.Contains(def, viewName) {
				deps = append(deps, viewRef{schema: s.Name, name: rel.Name})
			}
		}
	}

	return deps
}

// viewColumnsChanged returns true if existing columns were renamed, removed,
// reordered, or had their types changed. Adding new columns at the end is safe
// for CREATE OR REPLACE, but type changes require DROP + CREATE.
func viewColumnsChanged(entry RelationDiffEntry) bool {
	if entry.From == nil || entry.To == nil {
		return false
	}
	fromCols := entry.From.Columns
	toCols := entry.To.Columns
	// Fewer columns in target means columns were removed.
	if len(toCols) < len(fromCols) {
		return true
	}
	// Check that existing columns kept their names, order, and types.
	for i := range fromCols {
		if fromCols[i].Name != toCols[i].Name {
			return true
		}
		// Type change (e.g., function return type changed): PG cannot use
		// CREATE OR REPLACE VIEW if a column's data type changes.
		if fromCols[i].TypeOID != toCols[i].TypeOID {
			return true
		}
	}
	return false
}
