package catalog

import (
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
)

// generateIndexDDL produces CREATE INDEX / DROP INDEX / ALTER INDEX operations from the diff.
func generateIndexDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops, renameOps, attachOps []MigrationOp
	attaching := make(map[uint32]bool) // to-catalog OIDs of indexes created to be attached

	for _, relEntry := range diff.Relations {
		switch relEntry.Action {
		case DiffAdd:
			// New table — generate CREATE INDEX for all standalone indexes, and
			// column renames for both standalone and constraint-backed indexes.
			if relEntry.To != nil {
				indexes := to.IndexesOf(relEntry.To.OID)
				for _, idx := range indexes {
					if idx.ParentIndexOID != 0 {
						continue // managed by the parent's index
					}
					if idx.ConstraintOID != 0 {
						renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
						continue // index creation managed by constraint DDL
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
					renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
				}
			}
		case DiffModify:
			if relEntry.From != nil && relEntry.To != nil && tableRebuilt(from, to, relEntry.From, relEntry.To) {
				for _, idx := range to.IndexesOf(relEntry.To.OID) {
					if idx.ParentIndexOID != 0 {
						continue
					}
					renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
				}
				continue
			}
			recreatedConIdxOIDs := make(map[uint32]bool)
			for _, ce := range relEntry.Constraints {
				if (ce.Action == DiffAdd || ce.Action == DiffModify) && ce.To != nil && ce.To.IndexOID != 0 {
					recreatedConIdxOIDs[ce.To.IndexOID] = true
					if idx := to.indexes[ce.To.IndexOID]; idx != nil && idx.ParentIndexOID == 0 {
						renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
					}
				}
			}
			for _, idxEntry := range relEntry.Indexes {
				switch idxEntry.Action {
				case DiffAdd:
					idx := idxEntry.To
					if idx.ConstraintOID != 0 {
						continue
					}
					if idx.ParentIndexOID != 0 {
						attachOps = append(attachOps, attachIndexOp(to, relEntry.SchemaName, idx))
						attaching[idx.OID] = true
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
					renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idx, defaultIndexColumnNames(to, idx), effectiveIndexColumnNames(to, idx))...)
				case DiffDrop:
					idx := idxEntry.From
					if idx.ConstraintOID != 0 {
						continue
					}
					sql := fmt.Sprintf("DROP INDEX %s", formatIndexTargetName(relEntry.SchemaName, idx))
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
						if idxFrom.ConstraintOID != 0 && idxTo.ConstraintOID != 0 && !recreatedConIdxOIDs[idxTo.OID] && !indexesStructurallyChanged(from, to, idxFrom, idxTo) {
							renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idxTo, effectiveIndexColumnNames(from, idxFrom), effectiveIndexColumnNames(to, idxTo))...)
						}
						continue
					}
					if idxTo.ParentIndexOID != 0 {
						attachOps = append(attachOps, attachIndexOp(to, relEntry.SchemaName, idxTo))
						if !indexesChanged(from, to, idxFrom, idxTo) {
							continue
						}
					}
					if !indexesStructurallyChanged(from, to, idxFrom, idxTo) {
						renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idxTo, effectiveIndexColumnNames(from, idxFrom), effectiveIndexColumnNames(to, idxTo))...)
						continue
					}
					// Indexes cannot be structurally ALTERed — DROP + CREATE.
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
					renameOps = append(renameOps, buildIndexColumnRenameOps(relEntry.SchemaName, idxTo, defaultIndexColumnNames(to, idxTo), effectiveIndexColumnNames(to, idxTo))...)
					// Dropping the index took its comment; an unchanged one makes
					// no comment diff, so restore it here.
					if txt, ok := to.GetComment('i', idxTo.OID, 0); ok && txt != "" {
						if fromTxt, _ := from.GetComment('i', idxFrom.OID, 0); fromTxt == txt {
							desc := relEntry.SchemaName + "." + idxTo.Name
							if csql := formatCommentSQL(to, 'i', desc, 0, txt); csql != "" {
								renameOps = append(renameOps, MigrationOp{
									Type:          OpComment,
									ObjectName:    desc,
									SQL:           csql,
									Transactional: true,
									Phase:         PhaseMain,
									ObjType:       'i',
									Priority:      PriorityMetadata,
								})
							}
						}
					}
				}
			}
		}
	}

	// Sort for determinism: drops before creates, then by schema + name.
	// A partition's index to be attached is created before the parent's
	// index, whose CREATE INDEX then adopts it instead of making one under a
	// name of its own.
	sort.SliceStable(ops, func(i, j int) bool {
		// Drops before creates.
		iDrop := ops[i].Type == OpDropIndex
		jDrop := ops[j].Type == OpDropIndex
		if iDrop != jDrop {
			return iDrop
		}
		if iChild, jChild := attaching[ops[i].ObjOID], attaching[ops[j].ObjOID]; iChild != jChild {
			return iChild
		}
		if ops[i].SchemaName != ops[j].SchemaName {
			return ops[i].SchemaName < ops[j].SchemaName
		}
		return ops[i].ObjectName < ops[j].ObjectName
	})

	// Attribute renames run after CREATE INDEX, and index attaches run last.
	ops = append(ops, renameOps...)
	return append(ops, attachOps...)
}

func isTempIndex(schemaName string, idx *Index) bool {
	if isTempSchemaName(schemaName) {
		return true
	}
	if idx != nil && idx.Schema != nil {
		if isTempSchemaName(idx.Schema.Name) {
			return true
		}
		for _, r := range idx.Schema.Relations {
			if r != nil && r.OID == idx.RelOID && r.Persistence == 't' {
				return true
			}
		}
	}
	return false
}

func formatIndexTargetName(schemaName string, idx *Index) string {
	if idx == nil {
		return ""
	}
	if isTempIndex(schemaName, idx) {
		return quoteIdentifier(idx.Name)
	}
	return quoteIdentifier(schemaName) + "." + quoteIdentifier(idx.Name)
}

func buildIndexColumnRenameOps(schemaName string, idx *Index, fromNames, toNames []string) []MigrationOp {
	if len(fromNames) != len(toNames) || stringSliceEqual(fromNames, toNames) {
		return nil
	}
	cur := make([]string, len(fromNames))
	copy(cur, fromNames)
	used := make(map[string]bool, len(cur))
	for _, n := range cur {
		used[n] = true
	}
	targetSet := make(map[string]bool, len(toNames))
	for _, n := range toNames {
		targetSet[n] = true
	}
	var ops []MigrationOp
	qIdx := formatIndexTargetName(schemaName, idx)
	for {
		progressed := false
		for i := range cur {
			if cur[i] == toNames[i] || used[toNames[i]] {
				continue
			}
			ops = append(ops, MigrationOp{
				Type:          OpAlterTable,
				SchemaName:    schemaName,
				ObjectName:    idx.Name,
				SQL:           fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", qIdx, quoteIdentifier(cur[i]), quoteIdentifier(toNames[i])),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'i',
				ObjOID:        idx.OID,
				Priority:      PriorityIndex,
			})
			delete(used, cur[i])
			cur[i] = toNames[i]
			used[toNames[i]] = true
			progressed = true
		}
		if progressed {
			continue
		}
		// Remaining mismatches form rename cycles; move one participant aside to
		// an unused temporary name so the cycle unwinds cleanly.
		cycIdx := -1
		for i := range cur {
			if cur[i] != toNames[i] {
				cycIdx = i
				break
			}
		}
		if cycIdx < 0 {
			break
		}
		tmp := fmt.Sprintf("__omni_idx_col_%d", cycIdx)
		for k := 1; used[tmp] || targetSet[tmp]; k++ {
			tmp = fmt.Sprintf("__omni_idx_col_%d_%d", cycIdx, k)
		}
		ops = append(ops, MigrationOp{
			Type:          OpAlterTable,
			SchemaName:    schemaName,
			ObjectName:    idx.Name,
			SQL:           fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s", qIdx, quoteIdentifier(cur[cycIdx]), quoteIdentifier(tmp)),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'i',
			ObjOID:        idx.OID,
			Priority:      PriorityIndex,
		})
		delete(used, cur[cycIdx])
		cur[cycIdx] = tmp
		used[tmp] = true
	}
	return ops
}

// attachIndexOp attaches a partition's index to its parent's.
func attachIndexOp(to *Catalog, schemaName string, idx *Index) MigrationOp {
	parent := to.indexes[idx.ParentIndexOID]
	return MigrationOp{
		Type:          OpAlterTable,
		SchemaName:    schemaName,
		ObjectName:    idx.Name,
		SQL:           fmt.Sprintf("ALTER INDEX %s.%s ATTACH PARTITION %s.%s", quoteIdentifier(parent.Schema.Name), quoteIdentifier(parent.Name), quoteIdentifier(schemaName), quoteIdentifier(idx.Name)),
		Transactional: true,
		Phase:         PhaseMain,
		ObjType:       'i',
		ObjOID:        idx.OID,
		Priority:      PriorityIndex,
	}
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
	if (rel != nil && rel.Persistence == 't') || isTempSchemaName(schemaName) {
		b.WriteString(quoteIdentifier(rel.Name))
	} else {
		b.WriteString(quoteIdentifier(schemaName))
		b.WriteString(".")
		b.WriteString(quoteIdentifier(rel.Name))
	}

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
			if exprIdx < len(idx.Exprs) {
				b.WriteString(indexElemExpr(idx.Exprs[exprIdx]))
				exprIdx++
			}
		} else {
			colName := resolveAttNum(rel, idx.Columns[i])
			b.WriteString(quoteIdentifier(colName))
		}
		if i < len(idx.Collations) && len(idx.Collations[i]) > 0 {
			b.WriteString(" COLLATE ")
			for j, part := range idx.Collations[i] {
				if j > 0 {
					b.WriteString(".")
				}
				b.WriteString(quoteIdentifier(part))
			}
		}
		if i < len(idx.OpClasses) && idx.OpClasses[i] != "" {
			b.WriteString(" ")
			for j, part := range strings.Split(idx.OpClasses[i], ".") {
				if j > 0 {
					b.WriteString(".")
				}
				b.WriteString(quoteIdentifier(part))
			}
			if i < len(idx.OpClassOptions) && idx.OpClassOptions[i] != "" {
				b.WriteString(" (" + idx.OpClassOptions[i] + ")")
			}
		}
		// IndOption flags: bit 0 is DESC, bit 1 NULLS FIRST. DESC alone
		// sorts nulls first, so DESC without the bit needs NULLS LAST.
		if i < len(idx.IndOption) {
			desc, nullsFirst := idx.IndOption[i]&1 != 0, idx.IndOption[i]&2 != 0
			if desc {
				b.WriteString(" DESC")
			}
			if nullsFirst {
				b.WriteString(" NULLS FIRST")
			} else if desc {
				b.WriteString(" NULLS LAST")
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

// indexElemExpr writes an index expression as an index element: bare when
// it reads as a function call, parenthesized otherwise, as the grammar's
// index_elem requires.
//
// pg: src/backend/utils/adt/ruleutils.c — pg_get_indexdef_worker, looks_like_function
func indexElemExpr(expr string) string {
	if parsed, err := parseScalar("SELECT " + expr); err == nil {
		switch n := parsed.(type) {
		case *nodes.FuncCall:
			if n.Over == nil {
				return expr
			}
		case *nodes.CoalesceExpr, *nodes.MinMaxExpr, *nodes.SQLValueFunction:
			return expr
		}
	}
	return "(" + expr + ")"
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
