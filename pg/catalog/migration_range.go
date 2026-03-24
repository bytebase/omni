package catalog

import (
	"fmt"
	"sort"
)

// generateRangeDDL produces DDL operations for range type changes
// (CREATE TYPE ... AS RANGE, DROP TYPE, and warnings for subtype changes).
func generateRangeDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Ranges {
		switch entry.Action {
		case DiffAdd:
			if entry.To == nil {
				continue
			}
			subTypeName := to.FormatType(entry.To.SubTypeOID, -1)
			qualName := migrationQualifiedName(entry.SchemaName, entry.Name)
			sql := fmt.Sprintf("CREATE TYPE %s AS RANGE (SUBTYPE = %s)", qualName, subTypeName)
			ops = append(ops, MigrationOp{
				Type:          OpCreateType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           sql,
				Transactional: true,
			})
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
			})
		case DiffModify:
			if entry.From == nil || entry.To == nil {
				continue
			}
			qualName := migrationQualifiedName(entry.SchemaName, entry.Name)
			fromSubType := from.FormatType(entry.From.SubTypeOID, -1)
			toSubType := to.FormatType(entry.To.SubTypeOID, -1)
			ops = append(ops, MigrationOp{
				Type:          OpAlterType,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("-- range type %s subtype changed from %s to %s; requires DROP + CREATE", qualName, fromSubType, toSubType),
				Warning:       fmt.Sprintf("range type %s cannot be ALTERed in place: subtype changed from %s to %s; DROP + CREATE would break dependent columns", qualName, fromSubType, toSubType),
				Transactional: true,
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
