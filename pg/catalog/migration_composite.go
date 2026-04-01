package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateCompositeDDL produces CREATE TYPE ... AS, DROP TYPE, and ALTER TYPE
// (ADD/DROP ATTRIBUTE) operations for composite type changes.
func generateCompositeDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, entry := range diff.CompositeTypes {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpCreateType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           formatCreateCompositeType(to, entry.SchemaName, entry.Name, entry.To),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       't',
				ObjOID:        entry.To.RowTypeOID,
				Priority:      PriorityType,
			})

		case DiffDrop:
			if entry.From == nil {
				continue
			}
			qn := migrationQualifiedName(entry.SchemaName, entry.Name)
			ops = append(ops, MigrationOp{
				Type:          OpDropType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP TYPE %s", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       't',
				ObjOID:        entry.From.RowTypeOID,
				Priority:      PriorityType,
			})

		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			ops = append(ops, generateCompositeAlterOps(from, to, entry)...)
		}
	}

	// Deterministic ordering: drops first, then creates, then alters.
	sort.Slice(ops, func(i, j int) bool {
		oi, oj := compositeOpOrder(ops[i].Type), compositeOpOrder(ops[j].Type)
		if oi != oj {
			return oi < oj
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	return ops
}

// compositeOpOrder returns a sort key for composite type operation ordering.
func compositeOpOrder(t MigrationOpType) int {
	switch t {
	case OpDropType:
		return 0
	case OpCreateType:
		return 1
	default:
		return 2
	}
}

// formatCreateCompositeType renders a CREATE TYPE ... AS (...) statement.
func formatCreateCompositeType(cat *Catalog, schemaName, name string, rel *Relation) string {
	qn := migrationQualifiedName(schemaName, name)
	var fields []string
	for _, col := range rel.Columns {
		typeName := cat.FormatType(col.TypeOID, col.TypeMod)
		fields = append(fields, fmt.Sprintf("%s %s", quoteIdentAlways(col.Name), typeName))
	}
	return fmt.Sprintf("CREATE TYPE %s AS (%s)", qn, strings.Join(fields, ", "))
}

// generateCompositeAlterOps produces ALTER TYPE ADD/DROP ATTRIBUTE operations
// for a modified composite type.
func generateCompositeAlterOps(fromCat, toCat *Catalog, entry CompositeTypeDiffEntry) []MigrationOp {
	var ops []MigrationOp
	qn := migrationQualifiedName(entry.SchemaName, entry.Name)
	typeOID := entry.To.RowTypeOID

	// Build maps of column name -> Column for both sides.
	fromCols := make(map[string]*Column, len(entry.From.Columns))
	for _, c := range entry.From.Columns {
		fromCols[c.Name] = c
	}
	toCols := make(map[string]*Column, len(entry.To.Columns))
	for _, c := range entry.To.Columns {
		toCols[c.Name] = c
	}

	// Dropped attributes.
	var droppedNames []string
	for _, c := range entry.From.Columns {
		if _, ok := toCols[c.Name]; !ok {
			droppedNames = append(droppedNames, c.Name)
		}
	}
	sort.Strings(droppedNames)
	for _, name := range droppedNames {
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER TYPE %s DROP ATTRIBUTE %s", qn, quoteIdentAlways(name)),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       't',
			ObjOID:        typeOID,
			Priority:      PriorityType,
		})
	}

	// Added attributes.
	var addedNames []string
	for _, c := range entry.To.Columns {
		if _, ok := fromCols[c.Name]; !ok {
			addedNames = append(addedNames, c.Name)
		}
	}
	// Preserve declaration order from To.
	for _, c := range entry.To.Columns {
		found := false
		for _, name := range addedNames {
			if c.Name == name {
				found = true
				break
			}
		}
		if !found {
			continue
		}
		typeName := toCat.FormatType(c.TypeOID, c.TypeMod)
		ops = append(ops, MigrationOp{
			Type:          OpAlterType,
			SchemaName:    entry.SchemaName,
			ObjectName:    entry.Name,
			SQL:           fmt.Sprintf("ALTER TYPE %s ADD ATTRIBUTE %s %s", qn, quoteIdentAlways(c.Name), typeName),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       't',
			ObjOID:        typeOID,
			Priority:      PriorityType,
		})
	}

	// Modified attributes (type changed): ALTER TYPE ... ALTER ATTRIBUTE ... TYPE ...
	for _, tc := range entry.To.Columns {
		fc, ok := fromCols[tc.Name]
		if !ok {
			continue
		}
		fromType := fromCat.FormatType(fc.TypeOID, fc.TypeMod)
		toType := toCat.FormatType(tc.TypeOID, tc.TypeMod)
		if fromType != toType {
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("ALTER TYPE %s ALTER ATTRIBUTE %s TYPE %s", qn, quoteIdentAlways(tc.Name), toType),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       't',
				ObjOID:        typeOID,
				Priority:      PriorityType,
			})
		}
	}

	return ops
}
