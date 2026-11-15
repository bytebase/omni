package catalog

import (
	"maps"
	"slices"
)

// reorderTableChanges adjusts the ops of each modified table for what MariaDB
// refuses or loses otherwise:
//   - a generated column read by another one that is dropped too drops after it
//     (1054 otherwise);
//   - a key over a column the plan drops and adds again, to flip VIRTUAL and
//     STORED, is dropped and added around it, since the drop takes the column
//     out of the key;
//   - where every column of a table is dropped, the drops run after the adds,
//     since a table cannot lose its last column (1090);
//   - a column read by a generated column the plan keeps drops after that
//     column's new expression is in place (1054 otherwise);
//   - a foreign key whose only supporting key is dropped or replaced is dropped
//     first and added again at the end (1553 otherwise);
//   - so is a foreign key the plan keeps over a column whose type it changes,
//     on either side (1832, 1833).
func reorderTableChanges(ops []MigrationOp, from, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil {
			continue
		}
		ops = orderGeneratedDrops(ops, e)
		ops = dropColumnsAfterGeneratedReads(ops, e)
		ops = rebuildKeysOverFlippedColumns(ops, e, n)
		ops = dropLastColumnsAfterAdds(ops, e)
		ops = releaseForeignKeysOfReplacedKeys(ops, e)
	}
	return releaseForeignKeysOfRetypedColumns(ops, from, to, diff)
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
			ops[i].Priority = priorityColumn - 1 - len(dropped) + rank(ops[i].column, map[string]bool{})
		}
	}
	return ops
}

func dropColumnsAfterGeneratedReads(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	read := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action == DiffModify && ce.From != nil && ce.From.Generated != nil {
			for _, ref := range ce.From.Generated.Refs {
				read[toLower(ref)] = true
			}
		}
	}
	if len(read) == 0 {
		return ops
	}
	for i := range ops {
		if ops[i].Type == OpDropColumn && opOfTable(ops[i], e) && read[ops[i].column] {
			ops[i].Phase, ops[i].Priority = PhaseMain, priorityColumn+1
		}
	}
	return ops
}

func rebuildKeysOverFlippedColumns(ops []MigrationOp, e *TableDiffEntry, n *Normalizer) []MigrationOp {
	flipped := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action == DiffModify && generatedShapeChanged(ce.From, ce.To) {
			flipped[toLower(ce.Name)] = true
		}
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

func releaseForeignKeysOfRetypedColumns(ops []MigrationOp, from, to *Catalog, diff *SchemaDiff) []MigrationOp {
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
	if len(retyped) == 0 || from == nil || to == nil {
		return ops
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
			for _, con := range orderedForeignKeys(tbl) {
				toCon := foreignKeyOf(toTbl, con.Name)
				if toCon == nil || !foreignKeyOverRetyped(con, db.Name, tbl.Name, retyped, key) {
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
