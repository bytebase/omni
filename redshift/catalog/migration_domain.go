package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateDomainDDL produces CREATE DOMAIN, DROP DOMAIN, and ALTER DOMAIN
// operations from the diff.
func generateDomainDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.Domains {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpCreateType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           formatCreateDomain(to, entry.SchemaName, entry.Name, entry.To),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       't',
				ObjOID:        entry.To.TypeOID,
				Priority:      PriorityType,
			})

		case DiffDrop:
			qn := migrationQualifiedName(entry.SchemaName, entry.Name)
			var typeOID uint32
			if entry.From != nil {
				typeOID = entry.From.TypeOID
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP DOMAIN %s", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       't',
				ObjOID:        typeOID,
				Priority:      PriorityType,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			ops = append(ops, generateDomainAlterOps(from, to, entry)...)
		}
	}

	// Deterministic ordering: drops first, then creates, then alters; within each group by schema+name.
	sort.Slice(ops, func(i, j int) bool {
		oi, oj := domainOpOrder(ops[i].Type), domainOpOrder(ops[j].Type)
		if oi != oj {
			return oi < oj
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		if ops[i].ObjectName != ops[j].ObjectName {
			return ops[i].ObjectName < ops[j].ObjectName
		}
		return ops[i].SQL < ops[j].SQL
	})

	return ops
}

// domainOpOrder returns a sort key for domain operation types.
func domainOpOrder(t MigrationOpType) int {
	switch t {
	case OpDropType:
		return 0
	case OpCreateType:
		return 1
	default:
		return 2
	}
}

// formatCreateDomain renders a CREATE DOMAIN statement.
func formatCreateDomain(cat *Catalog, schemaName, name string, d *DomainType) string {
	var buf strings.Builder
	qn := migrationQualifiedName(schemaName, name)
	baseType := cat.FormatType(d.BaseTypeOID, d.BaseTypMod)
	buf.WriteString(fmt.Sprintf("CREATE DOMAIN %s AS %s", qn, baseType))

	if d.Default != "" {
		buf.WriteString(fmt.Sprintf(" DEFAULT %s", d.Default))
	}

	if d.NotNull {
		buf.WriteString(" NOT NULL")
	}

	// Sort constraints by name for determinism.
	constraints := make([]*DomainConstraint, len(d.Constraints))
	copy(constraints, d.Constraints)
	sort.Slice(constraints, func(i, j int) bool {
		return constraints[i].Name < constraints[j].Name
	})

	for _, c := range constraints {
		buf.WriteString(fmt.Sprintf(" CONSTRAINT %s CHECK (%s)",
			quoteIdentAlways(c.Name), c.CheckExpr))
	}

	return buf.String()
}

// generateDomainAlterOps produces ALTER DOMAIN statements for a modified domain.
func generateDomainAlterOps(from, to *Catalog, entry DomainDiffEntry) []MigrationOp {
	var ops []MigrationOp
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	d0 := entry.From
	d1 := entry.To
	typeOID := d1.TypeOID

	// Base type change: cannot ALTER DOMAIN type — emit warning with DROP+CREATE.
	fromType := from.FormatType(d0.BaseTypeOID, d0.BaseTypMod)
	toType := to.FormatType(d1.BaseTypeOID, d1.BaseTypMod)
	if fromType != toType {
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("-- cannot ALTER DOMAIN base type; must DROP and recreate\nDROP DOMAIN %s;\n%s", qn, formatCreateDomain(to, entry.SchemaName, entry.Name, d1)),
			Warning:       fmt.Sprintf("domain %s base type changed from %s to %s; requires DROP and recreate", entry.Name, fromType, toType),
			Transactional: true,
		})
		return ops
	}

	// Default change.
	if d0.Default != d1.Default {
		if d1.Default == "" {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER DOMAIN %s DROP DEFAULT", qn),
				Transactional: true,
			})
		} else {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER DOMAIN %s SET DEFAULT %s", qn, d1.Default),
				Transactional: true,
			})
		}
	}

	// NOT NULL change.
	if d0.NotNull != d1.NotNull {
		if d1.NotNull {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER DOMAIN %s SET NOT NULL", qn),
				Transactional: true,
			})
		} else {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER DOMAIN %s DROP NOT NULL", qn),
				Transactional: true,
			})
		}
	}

	// Constraint changes.
	fromCons := make(map[string]*DomainConstraint, len(d0.Constraints))
	for _, c := range d0.Constraints {
		fromCons[c.Name] = c
	}
	toCons := make(map[string]*DomainConstraint, len(d1.Constraints))
	for _, c := range d1.Constraints {
		toCons[c.Name] = c
	}

	// Dropped constraints (sorted for determinism).
	var droppedNames []string
	for name := range fromCons {
		if _, ok := toCons[name]; !ok {
			droppedNames = append(droppedNames, name)
		}
	}
	sort.Strings(droppedNames)
	for _, name := range droppedNames {
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER DOMAIN %s DROP CONSTRAINT %s", qn, quoteIdentAlways(name)),
			Transactional: true,
		})
	}

	// Added constraints (sorted for determinism).
	var addedNames []string
	for name := range toCons {
		if _, ok := fromCons[name]; !ok {
			addedNames = append(addedNames, name)
		}
	}
	sort.Strings(addedNames)
	for _, name := range addedNames {
		c := toCons[name]
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER DOMAIN %s ADD CONSTRAINT %s CHECK (%s)", qn, quoteIdentAlways(name), c.CheckExpr),
			Transactional: true,
		})
	}

	// Modified constraints (expression changed): drop + add (sorted for determinism).
	var modifiedNames []string
	for name, fc := range fromCons {
		tc, ok := toCons[name]
		if !ok {
			continue
		}
		if fc.CheckExpr != tc.CheckExpr {
			modifiedNames = append(modifiedNames, name)
		}
	}
	sort.Strings(modifiedNames)
	for _, name := range modifiedNames {
		tc := toCons[name]
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER DOMAIN %s DROP CONSTRAINT %s", qn, quoteIdentAlways(name)),
			Transactional: true,
		})
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER DOMAIN %s ADD CONSTRAINT %s CHECK (%s)", qn, quoteIdentAlways(name), tc.CheckExpr),
			Transactional: true,
		})
	}

	// Set metadata on all domain alter ops.
	for i := range ops {
		ops[i].Phase = PhaseMain
		ops[i].ObjType = 't'
		ops[i].ObjOID = typeOID
		ops[i].Priority = PriorityType
	}
	return ops
}
