package catalog

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// generateColumnDDL produces one ALTER TABLE op per column change within modified tables:
// ADD COLUMN, DROP COLUMN, or MODIFY COLUMN with the full column definition. A rename is seen
// as DROP + ADD because columns are keyed by name. An ADD COLUMN carries no placement here;
// positionColumnAdds adds FIRST or AFTER once the final statement order is known.
func generateColumnDDL(_, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		// Added tables carry their columns inline in CREATE TABLE.
		if entry.Action != DiffModify || len(entry.Columns) == 0 {
			continue
		}
		// The original-case identifier; the lower-cased diff key can name the wrong object.
		table := tableIdent(entry.To)

		for _, col := range entry.Columns {
			switch col.Action {
			case DiffAdd:
				if col.To == nil {
					continue
				}
				ops = append(ops, addColumnOp(entry, table, col.To, n))

			case DiffDrop:
				if col.From == nil {
					continue
				}
				ops = append(ops, dropColumnOp(entry, table, col.Name, col.From))

			case DiffModify:
				if col.To == nil {
					continue
				}
				// A change of generated storage mode, or between plain and generated, is a
				// DROP + ADD rather than an in-place MODIFY.
				if generatedShapeChanged(col.From, col.To) {
					ops = append(ops, dropColumnOp(entry, table, col.Name, col.From))
					ops = append(ops, addColumnOp(entry, table, col.To, n))
					continue
				}
				clause := fmt.Sprintf("MODIFY COLUMN %s", formatColumnDefinition(entry.To, col.To, n))
				ops = append(ops, MigrationOp{
					Type:         OpModifyColumn,
					Database:     entry.Database,
					ObjectName:   entry.Name,
					ParentObject: entry.Name,
					SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
					Phase:        PhaseMain,
					Priority:     priorityColumn,
					sortName:     columnSortName(entry.Database, entry.Name, col.Name),
					alterClause:  clause,
					column:       toLower(col.Name),
				})
			}
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			return columnOpRank(ops[i].Type) < columnOpRank(ops[j].Type)
		}
		// Generated-column drops precede plain drops, then name.
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addColumnOp builds an ADD COLUMN op. A generated column is added after plain columns so
// any base column it references is already in place.
func addColumnOp(entry *TableDiffEntry, table string, to *Column, n *Normalizer) MigrationOp {
	clause := fmt.Sprintf("ADD COLUMN %s", formatColumnDefinition(entry.To, to, n))
	return MigrationOp{
		Type:         OpAddColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     columnAddPriority(to),
		sortName:     columnAddSortName(entry, to),
		alterClause:  clause,
		column:       toLower(to.Name),
	}
}

// columnAddPriority returns priorityColumn+1 for a generated column (add it after the plain
// columns its expression may reference) and priorityColumn for a plain column.
func columnAddPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn + 1
	}
	return priorityColumn
}

// dropColumnOp builds a DROP COLUMN op. A generated column is dropped before the plain
// columns its expression may reference.
func dropColumnOp(entry *TableDiffEntry, table, colName string, from *Column) MigrationOp {
	clause := fmt.Sprintf("DROP COLUMN %s", migrationQuoteIdent(colName))
	return MigrationOp{
		Type:         OpDropColumn,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhasePre,
		Priority:     columnDropPriority(from),
		sortName:     columnSortName(entry.Database, entry.Name, colName),
		alterClause:  clause,
		column:       toLower(colName),
	}
}

// columnDropPriority returns priorityColumn-1 for a generated column and priorityColumn for a
// plain column.
func columnDropPriority(c *Column) int {
	if c != nil && c.Generated != nil {
		return priorityColumn - 1
	}
	return priorityColumn
}

// generatedShapeChanged reports whether a column gained or lost the generated attribute or
// flipped VIRTUAL and STORED, which is rendered as DROP + ADD. A change to only the expression
// is a normal MODIFY.
func generatedShapeChanged(from, to *Column) bool {
	if from == nil || to == nil {
		return false
	}
	fromGen, toGen := from.Generated != nil, to.Generated != nil
	if fromGen != toGen {
		return true
	}
	if fromGen && toGen && from.Generated.Stored != to.Generated.Stored {
		return true
	}
	return false
}

// columnOpRank orders column op types so drops sort ahead of adds/modifies locally.
func columnOpRank(t MigrationOpType) int {
	if t == OpDropColumn {
		return 0
	}
	return 1
}

// columnSortName is the stable secondary sort key for a column op: lower-cased
// database.table.column.
func columnSortName(database, table, column string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(column)
}

// columnAddSortName is the sort key of an ADD COLUMN: the table's key, a NUL, and the column's
// zero-padded position in `to`. An identifier cannot contain U+0000, so a table's plain adds
// sort before its modifies, which may read them.
func columnAddSortName(entry *TableDiffEntry, c *Column) string {
	return columnSortName(entry.Database, entry.Name, fmt.Sprintf("\x00%04d", columnPosition(entry.To, c.Name)))
}

// positionColumnAdds appends FIRST or AFTER to every ADD COLUMN clause of the sorted plan so the
// column lands where `to` places it; MariaDB otherwise appends it. It runs after
// mergeAutoIncrementKeyOps, which can move an ADD COLUMN into another statement of its table, and
// walks the clauses in execution order, so an AFTER only names a column that exists by then: one
// that survives from `from`, or one an earlier clause added. A column that no existing column
// follows in `to` gets no clause.
func positionColumnAdds(ops []MigrationOp, diff *SchemaDiff) []MigrationOp {
	parts := make([][]alterPart, len(ops))
	added := make(map[string]map[string]bool)
	for i, op := range ops {
		parts[i] = op.alterParts
		if parts[i] == nil && op.Type == OpAddColumn {
			parts[i] = []alterPart{{clause: op.alterClause, addsColumn: op.column}}
		}
		for _, part := range parts[i] {
			if part.addsColumn == "" {
				continue
			}
			key := aiTableKey(op.Database, op.ObjectName)
			if added[key] == nil {
				added[key] = make(map[string]bool)
			}
			added[key][part.addsColumn] = true
		}
	}
	layouts := make(map[string]*columnLayout)
	for i := range diff.Tables {
		entry := &diff.Tables[i]
		key := aiTableKey(entry.Database, entry.Name)
		if added[key] != nil && entry.Action == DiffModify && entry.From != nil && entry.To != nil {
			layouts[key] = newColumnLayout(entry, added[key])
		}
	}
	for i := range ops {
		op := &ops[i]
		l := layouts[aiTableKey(op.Database, op.ObjectName)]
		joined := joinAlterParts(parts[i])
		if l == nil || parts[i] == nil || !strings.HasSuffix(op.SQL, joined) {
			continue
		}
		placed := make([]alterPart, len(parts[i]))
		for j, part := range parts[i] {
			placed[j] = part
			if part.addsColumn != "" {
				placed[j].clause += l.place(part.addsColumn)
			}
		}
		op.SQL = strings.TrimSuffix(op.SQL, joined) + joinAlterParts(placed)
	}
	return ops
}

// columnLayout tracks, for one modified table, which of the target's columns exist as the plan's
// ADD COLUMN clauses run.
type columnLayout struct {
	columns []*Column       // the target's columns in order, hidden ones excluded
	exists  map[string]bool // lower-cased name → the column exists at this point of the plan
}

// newColumnLayout returns the layout of a modified table whose plan adds the named columns. A
// column the plan drops and adds again does not exist until its ADD runs.
func newColumnLayout(entry *TableDiffEntry, added map[string]bool) *columnLayout {
	l := &columnLayout{exists: make(map[string]bool)}
	for _, c := range entry.To.Columns {
		if c == nil || c.Hidden != ColumnHiddenNone {
			continue
		}
		l.columns = append(l.columns, c)
		name := toLower(c.Name)
		l.exists[name] = !added[name] && entry.From.GetColumn(c.Name) != nil
	}
	return l
}

// place returns the placement clause for adding the named column and records that it exists.
func (l *columnLayout) place(name string) string {
	defer func() { l.exists[name] = true }()
	pos := slices.IndexFunc(l.columns, func(c *Column) bool { return toLower(c.Name) == name })
	if pos < 0 || !slices.ContainsFunc(l.columns[pos+1:], func(c *Column) bool { return l.exists[toLower(c.Name)] }) {
		return ""
	}
	for k := pos - 1; k >= 0; k-- {
		if l.exists[toLower(l.columns[k].Name)] {
			return " AFTER " + migrationQuoteIdent(l.columns[k].Name)
		}
	}
	return " FIRST"
}

// formatColumnDefinition is a column as SHOW CREATE TABLE writes it, its CHECK
// included: that is valid DDL for ADD and MODIFY, and the server stores from it
// exactly what the diff compares.
func formatColumnDefinition(tbl *Table, c *Column, n *Normalizer) string {
	return n.CanonicalColumn(tbl, c)
}
