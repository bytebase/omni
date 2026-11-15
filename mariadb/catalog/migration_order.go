package catalog

// reorderTableChanges adjusts the ops of each modified table for what MariaDB
// refuses or loses otherwise:
//   - a generated column read by another one that is dropped too drops after it
//     (1054 otherwise);
//   - a key over a column the plan drops and adds again, to flip VIRTUAL and
//     STORED, is dropped and added around it, since the drop takes the column
//     out of the key;
//   - where every column of a table is dropped, the drops run after the adds,
//     since a table cannot lose its last column (1090);
//   - a foreign key whose only supporting key is dropped or replaced is dropped
//     first and added again at the end (1553 otherwise).
func reorderTableChanges(ops []MigrationOp, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil {
			continue
		}
		ops = orderGeneratedDrops(ops, e)
		ops = rebuildKeysOverFlippedColumns(ops, e, n)
		ops = dropLastColumnsAfterAdds(ops, e)
		ops = releaseForeignKeysOfReplacedKeys(ops, e)
	}
	return ops
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
