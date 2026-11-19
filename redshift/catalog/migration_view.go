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
				ops = append(ops, buildModifyMatViewOps(from, to, entry)...)
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
// Materialized views don't support CREATE OR REPLACE, so the drop takes the
// views that depend on it: they are dropped first and created again from the
// definitions they had.
func buildModifyMatViewOps(from, to *Catalog, entry RelationDiffEntry) []MigrationOp {
	var ops []MigrationOp

	// Deepest dependents first, so each is dropped before the view it reads.
	// Only the ones the target keeps unchanged are ours: a removed or changed
	// dependent has a diff of its own, which would duplicate this work.
	var dependents []viewRef
	for _, dep := range dependentViewsInDropOrder(from, entry.From) {
		fromRel, toRel := from.GetRelation(dep.schema, dep.name), to.GetRelation(dep.schema, dep.name)
		if fromRel == nil || toRel == nil || toRel.RelKind != 'v' || fromRel.CheckOption != toRel.CheckOption {
			continue
		}
		fromDef, _ := from.GetViewDefinition(dep.schema, dep.name)
		toDef, _ := to.GetViewDefinition(dep.schema, dep.name)
		if fromDef != toDef {
			continue
		}
		dependents = append(dependents, dep)
	}
	for _, dep := range dependents {
		op := MigrationOp{
			Type:          OpDropView,
			SchemaName:    dep.schema,
			ObjectName:    dep.name,
			SQL:           fmt.Sprintf("DROP VIEW %s", migrationQualifiedName(dep.schema, dep.name)),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			Priority:      PriorityView,
		}
		// The plan's dependency sort orders the drops by these OIDs.
		if rel := from.GetRelation(dep.schema, dep.name); rel != nil {
			op.ObjOID = rel.OID
		}
		ops = append(ops, op)
	}

	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetMatViewDefinition(entry.SchemaName, entry.Name)

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

	// The view's storage travels with the statement that names it, as it does
	// when the view is added.
	options := ""
	if rel := to.GetRelation(entry.SchemaName, entry.Name); rel != nil {
		options = deparseMatViewOptions(rel)
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

	ops = append(ops, dropOp, createOp)

	// The dependents read the view again, shallowest first, with the grants
	// the drop took with it.
	for i := len(dependents) - 1; i >= 0; i-- {
		dep := dependents[i]
		ops = append(ops, buildRecreateViewOp(from, dep))
		rel := to.GetRelation(dep.schema, dep.name)
		if rel == nil {
			continue
		}
		for _, g := range to.grants {
			if g.ObjType != 'r' || g.ObjOID != rel.OID {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpGrant,
				SchemaName:    dep.schema,
				ObjectName:    dep.name,
				SQL:           formatGrantSQL(g, rel.Name),
				Transactional: grantTransactional(g),
				Phase:         PhaseMain,
				ObjType:       g.ObjType,
				ObjOID:        g.ObjOID,
				Priority:      PriorityMetadata,
			})
		}
	}
	return ops
}

// buildRecreateViewOp creates a CREATE VIEW op for a view that a materialized
// view's re-creation dropped, from the definition it had.
func buildRecreateViewOp(from *Catalog, ref viewRef) MigrationOp {
	op := MigrationOp{
		Type:          OpCreateView,
		SchemaName:    ref.schema,
		ObjectName:    ref.name,
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'r',
		Priority:      PriorityView,
	}
	rel := from.GetRelation(ref.schema, ref.name)
	if rel != nil {
		op.ObjOID = rel.OID
	}
	def, _ := from.GetViewDefinition(ref.schema, ref.name)
	var b strings.Builder
	b.WriteString("CREATE VIEW ")
	b.WriteString(migrationQualifiedName(ref.schema, ref.name))
	b.WriteString(" AS ")
	b.WriteString(strings.TrimRight(def, " \t\n\r;"))
	if rel != nil {
		switch rel.CheckOption {
		case 'l':
			b.WriteString("\n WITH LOCAL CHECK OPTION")
		case 'c':
			b.WriteString("\n WITH CASCADED CHECK OPTION")
		}
	}
	op.SQL = b.String()
	return op
}

// dependentViewsInDropOrder returns the views that depend on a relation,
// transitively and each once, deepest first: a view that depends on another
// comes before it. A dependency is the one the catalog records for the view's
// definition, not a name in its text, which would take a view reading a column
// such as customer_id for one reading a relation named id.
func dependentViewsInDropOrder(c *Catalog, rel *Relation) []viewRef {
	seen := make(map[viewRef]bool)
	var ordered []viewRef
	var visit func(rel *Relation)
	visit = func(rel *Relation) {
		if rel == nil {
			return
		}
		for _, dep := range c.findNormalDependents('r', rel.OID) {
			if dep.ObjType != 'r' {
				continue
			}
			child := c.relationByOID[dep.ObjOID]
			if child == nil || child.RelKind != 'v' || child.Schema == nil {
				continue
			}
			ref := viewRef{schema: child.Schema.Name, name: child.Name}
			if seen[ref] {
				continue
			}
			seen[ref] = true
			visit(child)
			ordered = append(ordered, ref)
		}
	}
	visit(rel)
	return ordered
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
