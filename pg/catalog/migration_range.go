package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// formatRangeCollation formats the explicit COLLATION option of a range type,
// or returns "" when the range uses its subtype's default collation.
func formatRangeCollation(rt *RangeType) string {
	if rt == nil || rt.CollationName == "" {
		return ""
	}
	if rt.CollationSchema != "" {
		return migrationQualifiedName(rt.CollationSchema, rt.CollationName)
	}
	return quoteIdentAlways(rt.CollationName)
}

func formatCreateRange(cat *Catalog, schemaName, name string, rt *RangeType) string {
	subTypeName := cat.FormatType(rt.SubTypeOID, -1)
	qualName := migrationQualifiedName(schemaName, name)
	params := []string{"SUBTYPE = " + subTypeName}
	if rt.SubtypeOpclass != "" {
		params = append(params, "SUBTYPE_OPCLASS = "+formatCatalogIdentRef(rt.SubtypeOpclass))
	}
	if coll := formatRangeCollation(rt); coll != "" {
		params = append(params, "COLLATION = "+coll)
	}
	if rt.Canonical != "" {
		params = append(params, "CANONICAL = "+formatCatalogIdentRef(rt.Canonical))
	}
	if rt.SubtypeDiff != "" {
		params = append(params, "SUBTYPE_DIFF = "+formatCatalogIdentRef(rt.SubtypeDiff))
	}
	if mt := cat.TypeByOID(rt.MultirangeOID); mt != nil && mt.TypeName != "" {
		mrSchemaName := schemaName
		if ms := cat.schemas[mt.Namespace]; ms != nil && ms.Name != "" {
			mrSchemaName = ms.Name
		}
		if mrSchemaName != schemaName || mt.TypeName != makeMultirangeTypeName(name) {
			params = append(params, "MULTIRANGE_TYPE_NAME = "+migrationQualifiedName(mrSchemaName, mt.TypeName))
		}
	}
	return fmt.Sprintf("CREATE TYPE %s AS RANGE (%s)", qualName, strings.Join(params, ", "))
}

func buildCreateRangeOps(to *Catalog, schemaName, name string, rt *RangeType) []MigrationOp {
	qualName := migrationQualifiedName(schemaName, name)
	if rt.Canonical != "" {
		// Shell type lets canonical functions accept/return the range before full range definition.
		return []MigrationOp{
			{
				Type:          OpCreateType,
				SchemaName:    schemaName,
				ObjectName:    name,
				SQL:           fmt.Sprintf("CREATE TYPE %s", qualName),
				Transactional: true,
				Phase:         PhaseMain,
				Priority:      PriorityType,
			},
			{
				Type:          OpCreateType,
				SchemaName:    schemaName,
				ObjectName:    name,
				SQL:           formatCreateRange(to, schemaName, name, rt),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       't',
				ObjOID:        rt.OID,
				Priority:      PriorityTable,
			},
		}
	}
	return []MigrationOp{{
		Type:          OpCreateType,
		SchemaName:    schemaName,
		ObjectName:    name,
		SQL:           formatCreateRange(to, schemaName, name, rt),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       't',
		ObjOID:        rt.OID,
		Priority:      PriorityType,
	}}
}

// generateRangeDDL produces DDL operations for range type changes
// (CREATE TYPE ... AS RANGE, DROP TYPE, and replacements/warnings for modified ranges).
func generateRangeDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Ranges {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			ops = append(ops, buildCreateRangeOps(to, entry.SchemaName, entry.Name, entry.To)...)
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			qualName := migrationQualifiedName(entry.SchemaName, entry.Name)
			ops = append(ops, MigrationOp{
				Type:          OpDropType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP TYPE %s", qualName),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       't',
				ObjOID:        entry.From.OID,
				Priority:      PriorityType,
			})
		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			qualName := migrationQualifiedName(entry.SchemaName, entry.Name)
			fromSubType := from.FormatType(entry.From.SubTypeOID, -1)
			toSubType := to.FormatType(entry.To.SubTypeOID, -1)
			if fromSubType != toSubType {
				ops = append(ops, MigrationOp{
					Type:          OpAlterType,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("-- range type %s subtype changed from %s to %s; requires DROP + CREATE", qualName, fromSubType, toSubType),
					Warning:       fmt.Sprintf("range type %s cannot be ALTERed in place: subtype changed from %s to %s; DROP + CREATE would break dependent columns", qualName, fromSubType, toSubType),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       't',
					ObjOID:        entry.To.OID,
					Priority:      PriorityType,
				})
				continue
			}
			if !typeUsedByKeptColumn(from, to, nil, entry.From.OID) && !typeUsedByKeptColumn(from, to, nil, entry.From.MultirangeOID) {
				ops = append(ops,
					MigrationOp{
						Type:          OpDropType,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("DROP TYPE %s", qualName),
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       't',
						ObjOID:        entry.From.OID,
						Priority:      PriorityType,
					},
				)
				ops = append(ops, buildCreateRangeOps(to, entry.SchemaName, entry.Name, entry.To)...)
				continue
			}
			fromColl := formatRangeCollation(entry.From)
			if fromColl == "" {
				fromColl = "default"
			}
			toColl := formatRangeCollation(entry.To)
			if toColl == "" {
				toColl = "default"
			}
			changeDesc := fmt.Sprintf("collation changed from %s to %s", fromColl, toColl)
			if fromColl == toColl {
				changeDesc = "definition changed"
			}
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("-- range type %s %s; requires DROP + CREATE", qualName, changeDesc),
				Warning:       fmt.Sprintf("range type %s cannot be ALTERed in place: %s; DROP + CREATE would break dependent columns", qualName, changeDesc),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       't',
				ObjOID:        entry.To.OID,
				Priority:      PriorityType,
			})
		}
	}
	// Sort for determinism: by type then by schema+name.
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			return ops[i].Type < ops[j].Type
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})
	return ops
}
