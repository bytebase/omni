package catalog

import (
	"fmt"
	"maps"
	"slices"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// reorderTableChanges adjusts the ops of each modified table for what MariaDB
// refuses or loses otherwise:
//   - a generated column read by another one that is dropped too drops after it
//     (1054 otherwise);
//   - a key over a column the plan drops and adds again, to flip VIRTUAL and
//     STORED, is dropped and added around it, since the drop takes the column
//     out of the key, and so are the kept generated columns and checks that
//     read the column (1054, 4029);
//   - where every column of a table is dropped, the drops run after the adds,
//     since a table cannot lose its last column (1090);
//   - a column read by a generated column the plan keeps drops after that
//     column's new expression is in place (1054 otherwise);
//   - a foreign key whose only supporting key is dropped or replaced is dropped
//     first and added again at the end (1553 otherwise);
//   - so is a foreign key the plan keeps over a column whose type it changes,
//     on either side (1832, 1833), and one in another table whose referenced
//     key the plan drops or replaces;
//   - a generated column whose new expression reads a generated column the
//     plan adds is modified after the add (1054 otherwise);
//   - a key the plan keeps over an application-time period it replaces is
//     dropped before the period and added after it (4156 otherwise).
func reorderTableChanges(ops []MigrationOp, from, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil {
			continue
		}
		ops = orderGeneratedDrops(ops, e)
		ops = dropColumnsAfterGeneratedReads(ops, e)
		ops = modifyGeneratedAfterAddedReads(ops, e)
		ops = rebuildKeysOverReplacedPeriod(ops, e, n)
		ops = rebuildKeysOverFlippedColumns(ops, e, n)
		ops = rebuildReadersOfReplacedColumns(ops, e, n)
		ops = dropLastColumnsAfterAdds(ops, e)
		ops = releaseForeignKeysOfReplacedKeys(ops, e)
	}
	return releaseForeignKeysAcrossTables(ops, from, to, diff)
}

func opOfTable(op MigrationOp, e *TableDiffEntry) bool {
	return op.Database == e.Database && toLower(op.ParentObject) == toLower(e.Name)
}

func orderGeneratedDrops(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	dropped := make(map[string]*Column)
	for _, ce := range e.Columns {
		if ce.Action == DiffDrop && ce.From != nil && ce.From.Generated != nil {
			dropped[toLower(ce.Name)] = ce.From
		}
	}
	if len(dropped) < 2 {
		return ops
	}
	// rank counts the dropped generated columns that read a column, directly or
	// through another; a column drops after every one that reads it.
	var rank func(name string, seen map[string]bool) int
	rank = func(name string, seen map[string]bool) int {
		if seen[name] {
			return 0
		}
		seen[name] = true
		r := 0
		for other, col := range dropped {
			for _, ref := range col.Generated.Refs {
				if toLower(ref) == name {
					if d := rank(other, seen) + 1; d > r {
						r = d
					}
				}
			}
		}
		return r
	}
	for i := range ops {
		if ops[i].Type != OpDropColumn || !opOfTable(ops[i], e) {
			continue
		}
		if _, ok := dropped[ops[i].column]; ok {
			// The rank orders the drops within their priority, which stays
			// above the index drops (a drop takes the column out of its keys).
			ops[i].sortName = rankedColumnSortName(e, rank(ops[i].column, map[string]bool{}), ops[i].column)
		}
	}
	return ops
}

func dropColumnsAfterGeneratedReads(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	read := make(map[string]bool)
	removed := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action == DiffDrop {
			removed[toLower(ce.Name)] = true
		}
		if ce.Action == DiffModify && ce.From != nil && ce.From.Generated != nil {
			for _, ref := range ce.From.Generated.Refs {
				read[toLower(ref)] = true
			}
		}
	}
	if len(read) == 0 {
		return ops
	}
	// Only a column the plan removes for good: a replaced one drops ahead of
	// its own add (see replacedColumns).
	for i := range ops {
		if ops[i].Type == OpDropColumn && opOfTable(ops[i], e) && read[ops[i].column] && removed[ops[i].column] {
			// After every modify, modifyGeneratedAfterAddedReads's included.
			ops[i].Phase, ops[i].Priority = PhaseMain, priorityColumn+3
		}
	}
	return ops
}

func modifyGeneratedAfterAddedReads(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	added := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action == DiffAdd && ce.To != nil && ce.To.Generated != nil {
			added[toLower(ce.Name)] = true
		}
	}
	if len(added) == 0 {
		return ops
	}
	reads := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action != DiffModify || ce.To == nil || ce.To.Generated == nil {
			continue
		}
		for _, ref := range ce.To.Generated.Refs {
			if added[toLower(ref)] {
				reads[toLower(ce.Name)] = true
			}
		}
	}
	for i := range ops {
		if ops[i].Type == OpModifyColumn && opOfTable(ops[i], e) && reads[ops[i].column] {
			// After the generated adds, at priorityColumn+1 (columnAddPriority).
			ops[i].Priority = priorityColumn + 2
		}
	}
	return ops
}

func rebuildKeysOverReplacedPeriod(ops []MigrationOp, e *TableDiffEntry, n *Normalizer) []MigrationOp {
	from, to := e.From, e.To
	if from.AppPeriodName == "" || to.AppPeriodName == "" || !appPeriodChanged(from, to) {
		return ops
	}
	changed := make(map[string]bool)
	for _, ie := range e.Indexes {
		changed[toLower(ie.Name)] = true
	}
	table := tableIdent(to)
	for _, idx := range from.Indexes {
		if changed[toLower(idx.Name)] || !indexOverPeriod(idx, from.AppPeriodName) {
			continue
		}
		if target := tableIndexByName(to, idx.Name); target != nil {
			ops = append(ops, dropIndexOp(e, table, idx), addIndexOp(e, table, target, n))
		}
	}
	return ops
}

func indexOverPeriod(idx *Index, period string) bool {
	for _, ic := range idx.Columns {
		if ic.WithoutOverlaps && toLower(ic.Name) == toLower(period) {
			return true
		}
	}
	return false
}

// replacedColumns returns the columns the plan drops and adds again: those
// that flip between VIRTUAL, STORED and plain, and the generated columns, kept
// or modified, whose definition in from reads one, directly or through
// another, which MariaDB will not leave reading a dropped column (1054, 4029).
// Each maps to its depth, 0 for a flipped column and n for a reader n steps
// from one.
func replacedColumns(e *TableDiffEntry) map[string]int {
	gone := make(map[string]int)
	changed := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action != DiffModify {
			changed[toLower(ce.Name)] = true
		}
		if ce.Action == DiffModify && generatedShapeChanged(ce.From, ce.To) {
			gone[toLower(ce.Name)] = 0
		}
	}
	for depth := 1; len(gone) > 0; depth++ {
		added := false
		for _, col := range e.From.Columns {
			name := toLower(col.Name)
			if _, ok := gone[name]; ok || changed[name] || col.Generated == nil || e.To.GetColumn(col.Name) == nil {
				continue
			}
			for _, ref := range col.Generated.Refs {
				if d, ok := gone[toLower(ref)]; ok && d < depth {
					gone[name], added = depth, true
					break
				}
			}
		}
		if !added {
			break
		}
	}
	return gone
}

// rebuildReadersOfReplacedColumns drops the generated columns and kept table
// checks that read a replaced column ahead of it, and adds them again after
// it, deeper readers dropped first and added last. A modified reader is
// added with its new definition in place of its MODIFY.
func rebuildReadersOfReplacedColumns(ops []MigrationOp, e *TableDiffEntry, n *Normalizer) []MigrationOp {
	gone := replacedColumns(e)
	if len(gone) == 0 {
		return ops
	}
	ops = slices.DeleteFunc(ops, func(op MigrationOp) bool {
		d, ok := gone[op.column]
		return ok && d > 0 && op.Type == OpModifyColumn && opOfTable(op, e)
	})
	table := tableIdent(e.To)
	maxDepth := 0
	for _, d := range gone {
		maxDepth = max(maxDepth, d)
	}
	for _, col := range e.From.Columns {
		depth := gone[toLower(col.Name)]
		if depth == 0 {
			continue
		}
		drop, add := dropColumnOp(e, table, col.Name, col), addColumnOp(e, table, e.To.GetColumn(col.Name), n)
		drop.Priority, add.Priority = priorityColumn-2, priorityColumn+2
		drop.sortName = rankedColumnSortName(e, maxDepth-depth, col.Name)
		add.sortName = rankedColumnSortName(e, depth, fmt.Sprintf("%04d", columnPosition(e.To, col.Name)))
		ops = append(ops, drop, add)
	}
	changedChecks := make(map[string]bool)
	for _, ce := range e.Checks {
		changedChecks[toLower(ce.Name)] = true
	}
	for _, con := range orderedChecks(e.From) {
		if changedChecks[toLower(con.Name)] || !checkReadsAny(con, gone) {
			continue
		}
		for _, to := range orderedChecks(e.To) {
			if toLower(to.Name) == toLower(con.Name) {
				ops = append(ops, dropCheckOp(e, table, con), addCheckOp(e, table, to))
			}
		}
	}
	return ops
}

func checkReadsAny(con *Constraint, columns map[string]int) bool {
	expr, err := parseExpr(con.CheckExpr)
	if err != nil {
		return false
	}
	found := false
	nodes.Inspect(expr, func(node nodes.Node) bool {
		if ref, ok := node.(*nodes.ColumnRef); ok {
			if _, gone := columns[toLower(ref.Column)]; gone {
				found = true
			}
		}
		return !found
	})
	return found
}

func rebuildKeysOverFlippedColumns(ops []MigrationOp, e *TableDiffEntry, n *Normalizer) []MigrationOp {
	flipped := make(map[string]bool)
	for name := range replacedColumns(e) {
		flipped[name] = true
	}
	if len(flipped) == 0 {
		return ops
	}
	changed := make(map[string]bool)
	for _, ie := range e.Indexes {
		changed[toLower(ie.Name)] = true
	}
	table := tableIdent(e.To)
	for _, idx := range e.From.Indexes {
		if idx.Primary || changed[toLower(idx.Name)] || !indexReadsAny(idx, flipped) {
			continue
		}
		to := tableIndexByName(e.To, idx.Name)
		if to == nil {
			continue
		}
		ops = append(ops, dropIndexOp(e, table, idx), addIndexOp(e, table, to, n))
	}
	return ops
}

func indexReadsAny(idx *Index, columns map[string]bool) bool {
	for _, ic := range idx.Columns {
		if columns[toLower(ic.Name)] {
			return true
		}
	}
	return false
}

func dropLastColumnsAfterAdds(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	drops := 0
	for _, ce := range e.Columns {
		if ce.Action == DiffDrop {
			drops++
		}
	}
	if drops == 0 || drops < len(diffableColumns(e.From)) {
		return ops
	}
	for i := range ops {
		if ops[i].Type == OpDropColumn && opOfTable(ops[i], e) {
			ops[i].Phase, ops[i].Priority = PhaseMain, priorityColumn+1
		}
	}
	return ops
}

func releaseForeignKeysOfReplacedKeys(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	replaced := make(map[string]bool)
	for _, ie := range e.Indexes {
		if ie.Action == DiffDrop || ie.Action == DiffModify {
			replaced[toLower(ie.Name)] = true
		}
	}
	if len(replaced) == 0 {
		return ops
	}
	changedFK := make(map[string]bool)
	for _, fe := range e.ForeignKeys {
		changedFK[toLower(fe.Name)] = true
	}
	table := tableIdent(e.To)
	for _, con := range orderedForeignKeys(e.From) {
		if changedFK[toLower(con.Name)] || foreignKeyOf(e.To, con.Name) == nil {
			continue
		}
		needsReplaced, keptElsewhere := false, false
		for _, idx := range e.From.Indexes {
			if !indexCoversColumns(idx, con.Columns) {
				continue
			}
			if replaced[toLower(idx.Name)] {
				needsReplaced = true
			} else {
				keptElsewhere = true
			}
		}
		if !needsReplaced || keptElsewhere {
			continue
		}
		ops = append(ops, dropForeignKeyOps(e, con, map[string]bool{})...)
		ops = append(ops, addForeignKeyOp(e, table, foreignKeyOf(e.To, con.Name)))
	}
	return ops
}

func foreignKeyOf(t *Table, name string) *Constraint {
	for _, con := range t.Constraints {
		if con.Type == ConForeignKey && toLower(con.Name) == toLower(name) {
			return con
		}
	}
	return nil
}

func releaseForeignKeysAcrossTables(ops []MigrationOp, from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	retyped := make(map[string]bool)
	entries := make(map[string]*TableDiffEntry)
	key := func(db, table, column string) string {
		return encodeKeyFields("db", toLower(db), "table", toLower(table), "column", toLower(column))
	}
	for i := range diff.Tables {
		e := &diff.Tables[i]
		entries[key(e.Database, e.Name, "")] = e
		if e.Action != DiffModify {
			continue
		}
		for _, ce := range e.Columns {
			if ce.Action == DiffModify && ce.From != nil && ce.To != nil && columnRetyped(ce.From, ce.To) {
				retyped[key(e.Database, e.Name, ce.Name)] = true
			}
		}
	}
	if from == nil || to == nil {
		return ops
	}
	// parentKeyReplaced reports whether the plan drops or replaces every key of
	// the referenced table that the foreign key can use.
	parentKeyReplaced := func(con *Constraint, db string) bool {
		refDB := con.RefDatabase
		if refDB == "" {
			refDB = db
		}
		e := entries[key(refDB, con.RefTable, "")]
		if e == nil || e.Action != DiffModify || e.From == nil {
			return false
		}
		replaced := make(map[string]bool)
		for _, ie := range e.Indexes {
			if ie.Action == DiffDrop || ie.Action == DiffModify {
				replaced[toLower(ie.Name)] = true
			}
		}
		needsReplaced := false
		for _, idx := range e.From.Indexes {
			if !indexCoversColumns(idx, con.RefColumns) {
				continue
			}
			if !replaced[toLower(idx.Name)] {
				return false
			}
			needsReplaced = true
		}
		return needsReplaced
	}
	planned := make(map[string]bool)
	for _, op := range ops {
		if op.Type == OpDropForeignKey {
			planned[op.SQL] = true
		}
	}
	for _, db := range from.Databases() {
		toDB := to.GetDatabase(db.Name)
		if toDB == nil {
			continue
		}
		for _, name := range slices.Sorted(maps.Keys(db.Tables)) {
			tbl, toTbl := db.Tables[name], toDB.Tables[name]
			if toTbl == nil {
				continue
			}
			e := entries[key(db.Name, tbl.Name, "")]
			if e == nil {
				e = &TableDiffEntry{Action: DiffModify, Database: db.Name, Name: tbl.Name, From: tbl, To: toTbl}
			}
			changedFK := make(map[string]bool)
			for _, fe := range e.ForeignKeys {
				changedFK[toLower(fe.Name)] = true
			}
			for _, con := range orderedForeignKeys(tbl) {
				toCon := foreignKeyOf(toTbl, con.Name)
				if toCon == nil || changedFK[toLower(con.Name)] {
					continue
				}
				if !foreignKeyOverRetyped(con, db.Name, tbl.Name, retyped, key) && !parentKeyReplaced(con, db.Name) {
					continue
				}
				drop := dropForeignKeyOps(e, con, map[string]bool{})
				if planned[drop[0].SQL] {
					continue
				}
				planned[drop[0].SQL] = true
				ops = append(ops, drop...)
				ops = append(ops, addForeignKeyOp(e, tableIdent(toTbl), toCon))
			}
		}
	}
	return ops
}

// columnRetyped reports whether a column's type, character set or collation
// changes: what InnoDB refuses to change under a foreign key, while a comment,
// a default or NULL it takes. Measured on mariadb:11.8.8.
func columnRetyped(from, to *Column) bool {
	return toLower(from.ColumnType) != toLower(to.ColumnType) ||
		toLower(from.Charset) != toLower(to.Charset) || toLower(from.Collation) != toLower(to.Collation)
}

func foreignKeyOverRetyped(con *Constraint, db, table string, retyped map[string]bool, key func(db, table, column string) string) bool {
	for _, col := range con.Columns {
		if retyped[key(db, table, col)] {
			return true
		}
	}
	refDB := con.RefDatabase
	if refDB == "" {
		refDB = db
	}
	for _, col := range con.RefColumns {
		if retyped[key(refDB, con.RefTable, col)] {
			return true
		}
	}
	return false
}

// rankedColumnSortName orders a table's column ops of one priority by rank, then
// by key.
func rankedColumnSortName(e *TableDiffEntry, rank int, key string) string {
	return columnSortName(e.Database, e.Name, fmt.Sprintf("\x00%04d\x00%s", rank, toLower(key)))
}
