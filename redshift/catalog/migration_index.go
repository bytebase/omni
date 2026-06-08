package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateIndexDDL produces CREATE INDEX / DROP INDEX operations from the diff.
func generateIndexDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp

	for _, relEntry := range diff.Relations {
		switch relEntry.Action {
		case DiffAdd:
			// New table — generate CREATE INDEX for all standalone indexes.
			if relEntry.To != nil {
				indexes := to.IndexesOf(relEntry.To.OID)
				for _, idx := range indexes {
					if idx.ConstraintOID != 0 {
						continue // managed by constraint DDL
					}
					sql := buildCreateIndexSQL(idx, relEntry.To, relEntry.SchemaName)
					ops = append(ops, MigrationOp{
						Type:          OpCreateIndex,
						SchemaName:    relEntry.SchemaName,
						ObjectName:    idx.Name,
						SQL:           sql,
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'i',
						ObjOID:        idx.OID,
						Priority:      PriorityIndex,
					})
				}
			}
		case DiffModify:
			for _, idxEntry := range relEntry.Indexes {
				switch idxEntry.Action {
				case DiffAdd:
					idx := idxEntry.To
					if idx.ConstraintOID != 0 {
						continue
					}
					rel := relEntry.To
					sql := buildCreateIndexSQL(idx, rel, relEntry.SchemaName)
					ops = append(ops, MigrationOp{
						Type:          OpCreateIndex,
						SchemaName:    relEntry.SchemaName,
						ObjectName:    idx.Name,
						SQL:           sql,
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'i',
						ObjOID:        idx.OID,
						Priority:      PriorityIndex,
					})
				case DiffDrop:
					idx := idxEntry.From
					if idx.ConstraintOID != 0 {
						continue
					}
					sql := fmt.Sprintf("DROP INDEX %s.%s",
						quoteIdentifier(relEntry.SchemaName),
						quoteIdentifier(idx.Name))
					ops = append(ops, MigrationOp{
						Type:          OpDropIndex,
						SchemaName:    relEntry.SchemaName,
						ObjectName:    idx.Name,
						SQL:           sql,
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
					// Indexes cannot be ALTERed — DROP + CREATE.
					dropSQL := fmt.Sprintf("DROP INDEX %s.%s",
						quoteIdentifier(relEntry.SchemaName),
						quoteIdentifier(idxFrom.Name))
					ops = append(ops, MigrationOp{
						Type:          OpDropIndex,
						SchemaName:    relEntry.SchemaName,
						ObjectName:    idxFrom.Name,
						SQL:           dropSQL,
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'i',
						ObjOID:        idxFrom.OID,
						Priority:      PriorityIndex,
					})
					createSQL := buildCreateIndexSQL(idxTo, relEntry.To, relEntry.SchemaName)
					ops = append(ops, MigrationOp{
						Type:          OpCreateIndex,
						SchemaName:    relEntry.SchemaName,
						ObjectName:    idxTo.Name,
						SQL:           createSQL,
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'i',
						ObjOID:        idxTo.OID,
						Priority:      PriorityIndex,
					})
				}
			}
		}
	}

	// Sort for determinism: drops before creates, then by schema + name.
	sort.SliceStable(ops, func(i, j int) bool {
		// Drops before creates.
		iDrop := ops[i].Type == OpDropIndex
		jDrop := ops[j].Type == OpDropIndex
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

// buildCreateIndexSQL builds a CREATE [UNIQUE] INDEX statement.
func buildCreateIndexSQL(idx *Index, rel *Relation, schemaName string) string {
	var b strings.Builder
	b.WriteString("CREATE ")
	if idx.IsUnique {
		b.WriteString("UNIQUE ")
	}
	b.WriteString("INDEX ")
	b.WriteString(quoteIdentifier(idx.Name))
	b.WriteString(" ON ")
	b.WriteString(quoteIdentifier(schemaName))
	b.WriteString(".")
	b.WriteString(quoteIdentifier(rel.Name))

	// Access method — omit default btree.
	if idx.AccessMethod != "" && idx.AccessMethod != "btree" {
		b.WriteString(" USING ")
		b.WriteString(idx.AccessMethod)
	}

	// Key columns.
	nKey := idx.NKeyColumns
	if nKey == 0 {
		nKey = len(idx.Columns)
	}

	b.WriteString(" (")
	exprIdx := 0
	for i := 0; i < nKey; i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		if idx.Columns[i] == 0 {
			// Expression column.
			if exprIdx < len(idx.Exprs) {
				b.WriteString(idx.Exprs[exprIdx])
				exprIdx++
			}
		} else {
			colName := resolveAttNum(rel, idx.Columns[i])
			b.WriteString(quoteIdentifier(colName))
		}
		// IndOption flags.
		if i < len(idx.IndOption) {
			opt := idx.IndOption[i]
			if opt&1 != 0 {
				b.WriteString(" DESC")
			}
			if opt&2 != 0 {
				b.WriteString(" NULLS FIRST")
			}
		}
	}
	b.WriteString(")")

	// INCLUDE columns (after NKeyColumns).
	if nKey < len(idx.Columns) {
		b.WriteString(" INCLUDE (")
		first := true
		for i := nKey; i < len(idx.Columns); i++ {
			if !first {
				b.WriteString(", ")
			}
			first = false
			colName := resolveAttNum(rel, idx.Columns[i])
			b.WriteString(quoteIdentifier(colName))
		}
		b.WriteString(")")
	}

	// WHERE clause.
	if idx.WhereClause != "" {
		b.WriteString(" WHERE ")
		b.WriteString(idx.WhereClause)
	}

	return b.String()
}

// resolveAttNum finds a column name by its attnum.
func resolveAttNum(rel *Relation, attnum int16) string {
	for _, col := range rel.Columns {
		if col.AttNum == attnum {
			return col.Name
		}
	}
	return fmt.Sprintf("col%d", attnum)
}
