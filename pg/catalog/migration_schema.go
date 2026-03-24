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
			ops = append(ops, MigrationOp{
				Type:          OpCreateSchema,
				ObjectName:    entry.Name,
				SQL:           sql,
				Transactional: true,
			})
		case DiffDrop:
			ops = append(ops, MigrationOp{
				Type:          OpDropSchema,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP SCHEMA %s CASCADE", quoteIdentifier(entry.Name)),
				Transactional: true,
			})
		case DiffModify:
			if entry.From != nil && entry.To != nil && entry.From.Owner != entry.To.Owner {
				ops = append(ops, MigrationOp{
					Type:          OpAlterSchema,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER SCHEMA %s OWNER TO %s", quoteIdentifier(entry.Name), quoteIdentifier(entry.To.Owner)),
					Transactional: true,
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
