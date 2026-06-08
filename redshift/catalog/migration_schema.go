package catalog

import (
	"fmt"
	"sort"
)

// generateSchemaDDL produces DDL operations for schema-level changes
// (CREATE, DROP, ALTER OWNER).
func generateSchemaDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Schemas {
		switch entry.Action {
		case DiffAdd:
			sql := fmt.Sprintf("CREATE SCHEMA %s", quoteIdentifier(entry.Name))
			if entry.To != nil && entry.To.Owner != "" {
				sql += fmt.Sprintf(" AUTHORIZATION %s", quoteIdentifier(entry.To.Owner))
			}
			var schemaOID uint32
			if entry.To != nil {
				schemaOID = entry.To.OID
			}
			ops = append(ops, MigrationOp{
				Type:          OpCreateSchema,
				ObjectName:    entry.Name,
				SQL:           sql,
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'n',
				ObjOID:        schemaOID,
				Priority:      PrioritySchema,
			})
		case DiffDrop:
			var schemaOID uint32
			if entry.From != nil {
				schemaOID = entry.From.OID
			}
			ops = append(ops, MigrationOp{
				Type:          OpDropSchema,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP SCHEMA %s CASCADE", quoteIdentifier(entry.Name)),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'n',
				ObjOID:        schemaOID,
				Priority:      PrioritySchema,
			})
		case DiffModify:
			if entry.From != nil && entry.To != nil && entry.From.Owner != entry.To.Owner {
				ops = append(ops, MigrationOp{
					Type:          OpAlterSchema,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SCHEMA %s OWNER TO %s", quoteIdentifier(entry.Name), quoteIdentifier(entry.To.Owner)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'n',
					ObjOID:        entry.To.OID,
					Priority:      PrioritySchema,
				})
			}
		}
	}
	// Sort for determinism: by type (Create < Drop < Alter), then by name.
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			return ops[i].Type < ops[j].Type
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})
	return ops
}
