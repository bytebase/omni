package catalog

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// reorderTableChanges adjusts the ops of each modified table for what MariaDB
// refuses or loses otherwise:
//   - a column read by the expression or CHECK of a column the plan modifies
//     drops after the modify, where the plan removes it for good (1054);
//   - a generated column whose new expression reads a generated column the
//     plan adds is modified after the add (1054);
//   - a key the plan keeps over an application-time period it replaces is
//     dropped before the period and added after it (4156), and a kept
//     period's two columns change in one ALTER (4153);
//   - a column the plan replaces, to flip VIRTUAL, STORED or plain, takes the
//     keys over it and the generated columns and checks that read it out
//     first and back after (1054, 4029);
//   - every column drop runs after the dropped generated columns that read
//     it, and every add after the added columns its expression or CHECK reads
//     (1054, 4029);
//   - a table that loses every column it has keeps one until the end (1090);
//   - a foreign key whose only supporting key is dropped or replaced, on
//     either side, or that the plan keeps over a column whose type it
//     changes or that it replaces, is dropped first and added again at the
//     end (1553, 1828, 1832, 1833).
func reorderTableChanges(ops []MigrationOp, from, to *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	for i := range diff.Tables {
		e := &diff.Tables[i]
		if e.Action != DiffModify || e.From == nil || e.To == nil {
			continue
		}
		ops = dropColumnsAfterGeneratedReads(ops, e)
		ops = modifyGeneratedAfterAddedReads(ops, e)
		ops = rebuildKeysOverReplacedPeriod(ops, e, n)
		ops = modifyPeriodBoundsTogether(ops, e)
		ops = rebuildKeysOverFlippedColumns(ops, e, n)
		ops = rebuildReadersOfReplacedColumns(ops, e, n)
		ops = orderColumnDrops(ops, e)
		ops = orderColumnAdds(ops, e)
		ops = dropLastColumnsAfterAdds(ops, e)
		ops = hideColumnsLast(ops, e)
		ops = releaseForeignKeysOfReplacedKeys(ops, e)
	}
	return releaseForeignKeysAcrossTables(ops, from, to, diff)
}

func opOfTable(op MigrationOp, e *TableDiffEntry) bool {
	return op.Database == e.Database && toLower(op.ParentObject) == toLower(e.Name)
}

func dropColumnsAfterGeneratedReads(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	read := make(map[string]bool)
	removed := make(map[string]bool)
	sourceReads := columnReads(e.From)
	for _, ce := range e.Columns {
		if ce.Action == DiffDrop {
			removed[toLower(ce.Name)] = true
		}
		if ce.Action == DiffModify {
			for _, ref := range sourceReads[toLower(ce.Name)] {
				read[ref] = true
			}
		}
	}
	if len(read) == 0 {
		return ops
	}
	// Only a column the plan removes for good: a replaced one drops ahead of
	// its own add (see replacedColumns). A postponed column keeps the removed
	// columns it reads in turn; orderColumnDrops orders them.
	postponed := make(map[string]bool)
	var postpone func(name string)
	postpone = func(name string) {
		if !removed[name] || postponed[name] {
			return
		}
		postponed[name] = true
		for _, ref := range sourceReads[name] {
			postpone(ref)
		}
	}
	for name := range read {
		postpone(name)
	}
	for i := range ops {
		if ops[i].Type == OpDropColumn && opOfTable(ops[i], e) && postponed[ops[i].column] {
			// After every modify, modifyGeneratedAfterAddedReads's included.
			ops[i].Phase, ops[i].Priority = PhaseMain, priorityColumn+3
		}
	}
	return ops
}

func modifyGeneratedAfterAddedReads(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	// A column the plan replaces, dropping and adding it again, is added too.
	added := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action == DiffAdd || ce.Action == DiffModify && generatedShapeChanged(ce.From, ce.To) {
			added[toLower(ce.Name)] = true
		}
	}
	if len(added) == 0 {
		return ops
	}
	targetReads := columnReads(e.To)
	reads := make(map[string]bool)
	for _, ce := range e.Columns {
		if ce.Action != DiffModify || ce.To == nil {
			continue
		}
		for _, ref := range targetReads[toLower(ce.Name)] {
			if added[ref] {
				reads[toLower(ce.Name)] = true
			}
		}
	}
	for i := range ops {
		if ops[i].Type == OpModifyColumn && opOfTable(ops[i], e) && reads[ops[i].column] {
			// After every add, which runs at priorityColumn or, generated or
			// reading a generated one, priorityColumn+1.
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

// modifyPeriodBoundsTogether changes both columns of a kept application-time
// period in one ALTER: MariaDB holds the two to one type, so either change on
// its own is 4153.
func modifyPeriodBoundsTogether(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	from, to := e.From, e.To
	if from.AppPeriodName == "" || to.AppPeriodName == "" || appPeriodChanged(from, to) {
		return ops
	}
	start, end := -1, -1
	for i, op := range ops {
		if op.Type != OpModifyColumn || !opOfTable(op, e) || op.alterParts != nil {
			continue
		}
		switch op.column {
		case toLower(to.AppPeriodStartCol):
			start = i
		case toLower(to.AppPeriodEndCol):
			end = i
		}
	}
	if start < 0 || end < 0 {
		return ops
	}
	ops[start].alterParts = []alterPart{{clause: ops[start].alterClause}, {clause: ops[end].alterClause}}
	ops[start].SQL = "ALTER TABLE " + tableIdent(to) + " " + joinAlterParts(ops[start].alterParts)
	return slices.Delete(ops, end, end+1)
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
	// A replaced system-time period drops and adds both of its columns, the
	// one the target keeps as it is too.
	if systemPeriodReplaced(e) {
		for _, col := range e.From.Columns {
			if col.Generated != nil && col.Generated.RowBound != "" && e.To.GetColumn(col.Name) != nil {
				gone[toLower(col.Name)] = 0
			}
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
	for _, col := range e.From.Columns {
		depth := gone[toLower(col.Name)]
		if depth == 0 {
			continue
		}
		drop, add := dropColumnOp(e, table, col.Name, col), addColumnOp(e, table, e.To.GetColumn(col.Name), n)
		ops = append(ops, drop, add)
	}
	// A column whose own CHECK or expression DEFAULT reads a replaced column
	// is modified without them first, a column CHECK having no name to drop by
	// (1091), and back to its target definition after the adds.
	columnCheck := columnChecks(e.From)
	for _, col := range e.From.Columns {
		name := toLower(col.Name)
		if _, replaced := gone[name]; replaced {
			continue
		}
		con := columnCheck[name]
		checkReads := con != nil && checkReadsAny(con, gone)
		defaultReads := col.Default != nil && col.DefaultKind == ColumnDefaultExpression && exprReadsAny(*col.Default, gone)
		if !checkReads && !defaultReads {
			continue
		}
		target := e.To.GetColumn(col.Name)
		if target == nil {
			continue
		}
		context := *e.From
		context.Charset, context.Collation = "", ""
		released := *col
		if defaultReads {
			released.Default, released.DefaultKind = nil, ColumnDefaultNone
		}
		keep := map[string]*Constraint{}
		if con != nil && !checkReads {
			keep[name] = con
		}
		strip := "MODIFY COLUMN " + showColumnWithTable(&released, &context, keep, false)
		ops = append(ops, MigrationOp{
			Type: OpModifyColumn, Database: e.Database, ObjectName: e.Name, ParentObject: e.Name,
			SQL: "ALTER TABLE " + table + " " + strip, Phase: PhasePre, Priority: priorityCheckDrop,
			sortName: columnSortName(e.Database, e.Name, name), alterClause: strip,
		})
		if columnEntry(e, name) == nil {
			restore := "MODIFY COLUMN " + formatColumnDefinition(e.To, target, n)
			ops = append(ops, MigrationOp{
				Type: OpModifyColumn, Database: e.Database, ObjectName: e.Name, ParentObject: e.Name,
				SQL: "ALTER TABLE " + table + " " + restore, Phase: PhaseMain, Priority: priorityColumn + 3,
				sortName: columnSortName(e.Database, e.Name, name), alterClause: restore,
			})
		}
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
	return exprReadsAny(con.CheckExpr, columns)
}

func exprReadsAny(expr string, columns map[string]int) bool {
	for _, ref := range exprColumnRefs(expr) {
		if _, gone := columns[toLower(ref)]; gone {
			return true
		}
	}
	return false
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
		if changed[toLower(idx.Name)] || !indexReadsAny(idx, flipped) {
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

// orderColumnDrops runs the column drops of a table that precede its adds
// after every dropped column that reads them through its generated expression
// or column CHECK, directly or through another: MariaDB drops no column such a
// column still reads (1054, 4029). Columns that read each other drop in one
// ALTER. The order holds across permanent drops, replaced columns and their
// rebuilt readers, within one priority above the key drops.
func orderColumnDrops(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	in := func(op MigrationOp) bool {
		return op.Type == OpDropColumn && opOfTable(op, e)
	}
	reads := columnReads(e.From)
	// A column drops after the columns that read it.
	after := make(map[string][]string)
	for _, op := range ops {
		if in(op) {
			for _, ref := range reads[op.column] {
				after[ref] = append(after[ref], op.column)
			}
		}
	}
	return rankColumnOps(ops, e, in, after, func(op *MigrationOp, rank int) {
		if op.Phase == PhasePre {
			op.Priority = priorityColumn - 1
		}
		op.sortName = rankedColumnSortName(e, rank, op.column)
	})
}

// orderColumnAdds adds a table's columns after the added columns their
// generated expression or column CHECK reads (1054 otherwise): a column CHECK
// may read a column the target places after it. A column takes the latest
// priority among them; columns that read each other are added in one ALTER.
// Otherwise the adds keep their priority and the target's column order.
func orderColumnAdds(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	in := func(op MigrationOp) bool {
		return op.Type == OpAddColumn && op.Phase == PhaseMain && opOfTable(op, e)
	}
	priority := make(map[string]int)
	for _, op := range ops {
		if in(op) {
			priority[op.column] = op.Priority
		}
	}
	reads := columnReads(e.To)
	// A column takes the latest priority of the added columns it reads,
	// directly or through another.
	for changed := true; changed; {
		changed = false
		for name, p := range priority {
			for _, ref := range reads[name] {
				if q, ok := priority[ref]; ok && q > p {
					priority[name], p, changed = q, q, true
				}
			}
		}
	}
	return rankColumnOps(ops, e, in, reads, func(op *MigrationOp, rank int) {
		op.Priority = priority[op.column]
		op.sortName = rankedColumnSortName(e, rank, fmt.Sprintf("%04d", columnPosition(e.To, op.column)))
	})
}

// rankColumnOps ranks the ops in selects by the longest chain of their columns'
// edges among them (a column runs after every column it has an edge to),
// merges the columns that reach each other into the one statement of the
// cycle, and hands each remaining op its rank. A cycle's members share a rank.
func rankColumnOps(ops []MigrationOp, e *TableDiffEntry, selects func(MigrationOp) bool, edges map[string][]string, set func(op *MigrationOp, rank int)) []MigrationOp {
	member := make(map[string]bool)
	for _, op := range ops {
		if selects(op) && op.alterParts == nil {
			member[op.column] = true
		}
	}
	reaches := func(from, to string) bool {
		seen := map[string]bool{from: true}
		stack := []string{from}
		for len(stack) > 0 {
			name := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, next := range edges[name] {
				if next == to {
					return true
				}
				if member[next] && !seen[next] {
					seen[next] = true
					stack = append(stack, next)
				}
			}
		}
		return false
	}
	// comp maps each column to the first column, in the target's order, of
	// the cycle it is in, or to itself.
	names := slices.Sorted(maps.Keys(member))
	slices.SortStableFunc(names, func(a, b string) int { return columnPosition(e.To, a) - columnPosition(e.To, b) })
	comp := make(map[string]string)
	for _, a := range names {
		if _, ok := comp[a]; ok {
			continue
		}
		comp[a] = a
		for _, b := range names {
			if _, ok := comp[b]; !ok && reaches(a, b) && reaches(b, a) {
				comp[b] = a
			}
		}
	}
	ranks := make(map[string]int)
	var rank func(c string, seen map[string]bool) int
	rank = func(c string, seen map[string]bool) int {
		if r, ok := ranks[c]; ok {
			return r
		}
		if seen[c] {
			return 0
		}
		seen[c] = true
		r := 0
		for _, name := range names {
			if comp[name] != c {
				continue
			}
			for _, next := range edges[name] {
				if member[next] && comp[next] != c {
					r = max(r, rank(comp[next], seen)+1)
				}
			}
		}
		ranks[c] = r
		return r
	}
	for i := range ops {
		if selects(ops[i]) && ops[i].alterParts == nil {
			set(&ops[i], rank(comp[ops[i].column], map[string]bool{}))
		}
	}
	// One statement per cycle, at its first member, in the target's order.
	for _, c := range names {
		if comp[c] != c {
			continue
		}
		var cycle []int
		for _, name := range names {
			if comp[name] == c {
				for i := range ops {
					if selects(ops[i]) && ops[i].alterParts == nil && ops[i].column == name {
						cycle = append(cycle, i)
					}
				}
			}
		}
		if len(cycle) < 2 {
			continue
		}
		first := cycle[0]
		parts := make([]alterPart, 0, len(cycle))
		for _, i := range cycle {
			part := alterPart{clause: ops[i].alterClause}
			if ops[i].Type == OpAddColumn {
				part.addsColumn = ops[i].column
			} else {
				part.dropsColumn = ops[i].column
			}
			parts = append(parts, part)
			if ops[i].Priority > ops[first].Priority {
				ops[first].Priority = ops[i].Priority
			}
		}
		ops[first].alterParts = parts
		ops[first].SQL = "ALTER TABLE " + tableIdent(e.To) + " " + joinAlterParts(parts)
		drop := make(map[int]bool)
		for _, i := range cycle[1:] {
			drop[i] = true
		}
		kept := ops[:0]
		for i, op := range ops {
			if !drop[i] {
				kept = append(kept, op)
			}
		}
		ops = kept
	}
	return ops
}

// columnReads maps each column of t to the lower-cased columns its generated
// expression, column CHECK or expression DEFAULT reads, itself left out.
func columnReads(t *Table) map[string][]string {
	out := make(map[string][]string)
	checks := columnChecks(t)
	for _, col := range t.Columns {
		name := toLower(col.Name)
		var refs []string
		if col.Generated != nil {
			refs = append(refs, col.Generated.Refs...)
		}
		if con := checks[name]; con != nil {
			refs = append(refs, exprColumnRefs(con.CheckExpr)...)
		}
		if col.Default != nil && col.DefaultKind == ColumnDefaultExpression {
			refs = append(refs, exprColumnRefs(*col.Default)...)
		}
		for _, ref := range refs {
			if ref := toLower(ref); ref != name {
				out[name] = append(out[name], ref)
			}
		}
	}
	return out
}

func isAdded(ops []MigrationOp, e *TableDiffEntry, column string) bool {
	for _, op := range ops {
		if op.Type == OpAddColumn && opOfTable(op, e) && op.column == column {
			return true
		}
	}
	return false
}

// exprColumnRefs returns the columns an expression names.
func exprColumnRefs(text string) []string {
	expr, err := parseExpr(text)
	if err != nil {
		return nil
	}
	var refs []string
	nodes.Inspect(expr, func(node nodes.Node) bool {
		if ref, ok := node.(*nodes.ColumnRef); ok {
			refs = append(refs, ref.Column)
		}
		return true
	})
	return refs
}

// dropLastColumnsAfterAdds keeps a table from losing its last column (1090)
// where the plan drops every column it has, permanently or to add it again.
// The columns it only adds go first; where there are none, the first replaced
// column is dropped and added in one statement. A table that stays
// system-versioned likewise keeps a versioned column (4123): the versioned
// columns it only adds go first.
func dropLastColumnsAfterAdds(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	dropped := make(map[string]bool)
	for _, op := range ops {
		if op.Type == OpDropColumn && opOfTable(op, e) {
			dropped[op.column] = true
			for _, part := range op.alterParts {
				if part.dropsColumn != "" {
					dropped[part.dropsColumn] = true
				}
			}
		}
	}
	// So does DROP SYSTEM VERSIONING, its ROW START and ROW END columns.
	for _, col := range e.From.Columns {
		if droppedWithVersioning(e, col) {
			dropped[toLower(col.Name)] = true
		}
	}
	if len(dropped) == 0 {
		return ops
	}
	if len(dropped) < len(diffableColumns(e.From)) {
		if e.From.SystemVersioned && e.To.SystemVersioned && everyDropped(e, dropped, isVersionedColumn) {
			addFirst(ops, e, dropped, isVersionedColumn)
		}
		// A table keeps a visible column too (1113).
		if everyDropped(e, dropped, isVisibleColumn) {
			addFirst(ops, e, dropped, isVisibleColumn)
		}
		return ops
	}
	if addFirst(ops, e, dropped, func(*Column) bool { return true }) {
		return ops
	}
	// The replaced column the drop order takes last, which no other dropped
	// column reads, is dropped and added in one statement, together with the
	// columns grouped with it.
	last, add := -1, -1
	for d := range ops {
		if ops[d].Type != OpDropColumn || !opOfTable(ops[d], e) {
			continue
		}
		drops := partsOf(ops[d])
		for a := range ops {
			if ops[a].Type != OpAddColumn || !opOfTable(ops[a], e) || !sharesColumn(drops, partsOf(ops[a])) {
				continue
			}
			if last < 0 || ops[d].Phase > ops[last].Phase || ops[d].Phase == ops[last].Phase && (ops[d].Priority > ops[last].Priority || ops[d].Priority == ops[last].Priority && ops[d].sortName > ops[last].sortName) {
				last, add = d, a
			}
		}
	}
	if last < 0 {
		return ops
	}
	ops[add].alterParts = append(partsOf(ops[last]), partsOf(ops[add])...)
	ops[add].SQL = "ALTER TABLE " + tableIdent(e.To) + " " + joinAlterParts(ops[add].alterParts)
	return slices.Delete(ops, last, last+1)
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
		// A column the plan drops and adds again is held by its foreign keys
		// as much as one whose type changes (1828).
		if e.From != nil && e.To != nil {
			for name := range replacedColumns(e) {
				retyped[key(e.Database, e.Name, name)] = true
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
	return columnSortName(e.Database, e.Name, fmt.Sprintf("\x01%04d\x00%s", rank, toLower(key)))
}

// addFirst moves ahead of the drops the adds of the columns pick takes that the
// plan only adds, unless a column reads one added after it, a replaced one
// among them, directly or through another. It reports whether it moved any.
func addFirst(ops []MigrationOp, e *TableDiffEntry, dropped map[string]bool, pick func(*Column) bool) bool {
	reads := columnReads(e.To)
	early := make(map[string]bool)
	for i := range ops {
		if ops[i].Type == OpAddColumn && opOfTable(ops[i], e) && !dropped[ops[i].column] {
			if col := e.To.GetColumn(ops[i].column); col != nil && pick(col) {
				early[ops[i].column] = true
			}
		}
	}
	for changed := true; changed; {
		changed = false
		for name := range early {
			for _, ref := range reads[name] {
				if !early[ref] && ref != name && isAdded(ops, e, ref) {
					delete(early, name)
					changed = true
					break
				}
			}
		}
	}
	for i := range ops {
		if ops[i].Type == OpAddColumn && opOfTable(ops[i], e) && early[ops[i].column] {
			// Ahead of the drops, in the order orderColumnAdds gave them:
			// their sort names start where rankedColumnSortName's do not.
			ops[i].Phase, ops[i].Priority = PhasePre, priorityColumn-1
			ops[i].sortName = strings.Replace(ops[i].sortName, "\x01", "\x00", 1)
		}
	}
	return len(early) > 0
}

// isVersionedColumn reports whether a column of a system-versioned table keeps
// history: one not marked WITHOUT SYSTEM VERSIONING, other than ROW START and
// ROW END.
func isVersionedColumn(col *Column) bool {
	return col.SystemVersioning != "WITHOUT" && (col.Generated == nil || col.Generated.RowBound == "")
}

func isVisibleColumn(col *Column) bool { return !col.Invisible }

// everyDropped reports whether the plan drops every column of the source that
// is picks, there being one.
func everyDropped(e *TableDiffEntry, dropped map[string]bool, is func(*Column) bool) bool {
	found := false
	for _, col := range e.From.Columns {
		if is(col) {
			if !dropped[toLower(col.Name)] {
				return false
			}
			found = true
		}
	}
	return found
}

// partsOf returns the clauses of a column op's statement.
func partsOf(op MigrationOp) []alterPart {
	if op.alterParts != nil {
		return append([]alterPart(nil), op.alterParts...)
	}
	part := alterPart{clause: op.alterClause}
	if op.Type == OpAddColumn {
		part.addsColumn = op.column
	} else {
		part.dropsColumn = op.column
	}
	return []alterPart{part}
}

// sharesColumn reports whether a column the drops drop is one the adds add.
func sharesColumn(drops, adds []alterPart) bool {
	for _, d := range drops {
		for _, a := range adds {
			if d.dropsColumn != "" && d.dropsColumn == a.addsColumn {
				return true
			}
		}
	}
	return false
}

// hideColumnsLast runs a MODIFY that makes a column INVISIBLE after the table's
// other column changes, the one that makes another visible among them: a
// table keeps a visible column (1113).
func hideColumnsLast(ops []MigrationOp, e *TableDiffEntry) []MigrationOp {
	for i := range ops {
		if ops[i].Type != OpModifyColumn || !opOfTable(ops[i], e) || ops[i].column == "" {
			continue
		}
		from, to := e.From.GetColumn(ops[i].column), e.To.GetColumn(ops[i].column)
		if from != nil && to != nil && !from.Invisible && to.Invisible && ops[i].Phase == PhaseMain {
			ops[i].Priority = max(ops[i].Priority, priorityColumn+3)
		}
	}
	return ops
}
