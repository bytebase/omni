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
				})
			case 'm':
				qn := migrationQualifiedName(entry.SchemaName, entry.Name)
				ops = append(ops, MigrationOp{
					Type:          OpDropView,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("DROP MATERIALIZED VIEW %s", qn),
					Transactional: true,
				})
			}

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			switch entry.To.RelKind {
			case 'v':
				ops = append(ops, buildModifyViewOps(from, to, entry)...)
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
	}
}

// buildCreateMatViewOp creates a CREATE MATERIALIZED VIEW op.
func buildCreateMatViewOp(to *Catalog, entry RelationDiffEntry) MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetMatViewDefinition(entry.SchemaName, entry.Name)

	sql := fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s",
		qn, strings.TrimRight(def, " \t\n\r;"))

	return MigrationOp{
		Type:          OpCreateView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           sql,
		Transactional: true,
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
	})

	return ops
}

// buildModifyMatViewOps creates DROP + CREATE ops for a modified matview.
// Materialized views don't support CREATE OR REPLACE.
func buildModifyMatViewOps(to *Catalog, entry RelationDiffEntry) []MigrationOp {
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	def, _ := to.GetMatViewDefinition(entry.SchemaName, entry.Name)

	dropOp := MigrationOp{
		Type:          OpDropView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           fmt.Sprintf("DROP MATERIALIZED VIEW %s", qn),
		Transactional: true,
	}

	createSQL := fmt.Sprintf("CREATE MATERIALIZED VIEW %s AS %s",
		qn, strings.TrimRight(def, " \t\n\r;"))

	createOp := MigrationOp{
		Type:          OpCreateView,
		SchemaName:    entry.SchemaName,
		ObjectName:    entry.Name,
		SQL:           createSQL,
		Transactional: true,
	}

	return []MigrationOp{dropOp, createOp}
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
