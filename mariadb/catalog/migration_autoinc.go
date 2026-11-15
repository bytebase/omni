package catalog

import (
	"slices"
	"sort"
	"strings"
)

// AUTO_INCREMENT / backing-key statement grouping.
//
// MariaDB checks at the end of every ALTER TABLE statement that the table has at most one
// AUTO_INCREMENT column and that this column is the first column of a key (errno 1075). The
// generator otherwise emits one clause per statement in a fixed order (index drops, column
// drops, column adds and modifies, index adds), so a change that is valid only when the
// AUTO_INCREMENT column and its backing key move together fails on the statement that exposes
// the unkeyed or doubled column. formatCreateTable already inlines the backing key for a new
// table (autoIncSupportingIndex); this pass merges the involved clauses of a modified table into
// one ALTER TABLE.
//
// The grouped cases are:
//   - ADD COLUMN or MODIFY COLUMN that gains AUTO_INCREMENT, grouped with the ADD of a key that
//     has the column first.
//   - Making one column AUTO_INCREMENT while another still holds the attribute, grouped with
//     the other column's de-AUTO_INCREMENT MODIFY.
//   - Dropping the last key covering a still-AUTO_INCREMENT column, grouped with the column's
//     DROP, its de-AUTO_INCREMENT MODIFY, or a replacement covering ADD.
//   - The combined `DROP PRIMARY KEY, ADD PRIMARY KEY` statement, which accepts extra clauses.
//   - A demoted old PK member's `MODIFY ... NULL`, which may share the statement that drops the
//     PK.
//
// The pass is rule-based rather than a general per-table clause merger, so every other statement
// keeps its own op, warnings, and sort position. When the target itself leaves an
// AUTO_INCREMENT column unkeyed, the ops stay ungrouped and the apply surfaces the server's
// error.
//
// Ops participate only through the grouping metadata their constructors set (alterClause,
// addsIndex, dropsIndex on MigrationOp); checks, views, options, and FK constraint adds are
// never regrouped.
//
// Not handled: a storage-engine conversion in the same plan. The table-option ALTER
// (ENGINE=...) runs before every column and index statement and carries no grouping metadata,
// so converting a MyISAM table whose auto column is covered only in a non-first key position to
// InnoDB fails on the ENGINE statement.

// primaryIndexKey is the dropsIndex marker for the PRIMARY KEY, which has no user-facing name.
// The name PRIMARY is reserved for it, so the lower-cased form cannot collide with a secondary
// index.
const primaryIndexKey = "primary"

// mergeAutoIncrementKeyOps rewrites the generated ops so that no statement boundary exposes an
// unkeyed or doubled AUTO_INCREMENT column. Ops merged into another statement are removed and
// the surviving statement carries their clauses and warnings. GenerateMigrationWithNormalizer
// calls it before sortMigrationOps.
func mergeAutoIncrementKeyOps(ops []MigrationOp, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	if diff == nil || len(ops) == 0 {
		return ops
	}
	m := newAIMerger(ops)
	for i := range diff.Tables {
		entry := &diff.Tables[i]
		if entry.Action != DiffModify || entry.From == nil || entry.To == nil {
			continue
		}
		m.mergeTable(entry, n)
	}
	return m.result()
}

// aiMerger tracks the clause lists of the ops while hazards are resolved. clauses[i] is nil
// until op i is touched; consumed[i] marks an op whose clauses moved into another statement,
// and mergedInto[i] records which one. byTable groups the op indices by owning table.
type aiMerger struct {
	ops        []MigrationOp
	clauses    [][]alterPart
	consumed   []bool
	touched    []bool
	mergedInto []int
	byTable    map[string][]int
}

func newAIMerger(ops []MigrationOp) *aiMerger {
	m := &aiMerger{
		ops:        ops,
		clauses:    make([][]alterPart, len(ops)),
		consumed:   make([]bool, len(ops)),
		touched:    make([]bool, len(ops)),
		mergedInto: make([]int, len(ops)),
		byTable:    make(map[string][]int),
	}
	for i := range m.mergedInto {
		m.mergedInto[i] = -1
	}
	for i, op := range ops {
		key := aiTableKey(op.Database, op.ObjectName)
		m.byTable[key] = append(m.byTable[key], i)
	}
	return m
}

// aiTableKey is the collision-free per-table grouping key.
func aiTableKey(database, table string) string {
	return encodeKeyFields("db", database, "table", table)
}

// clausesOf returns op i's current clause list, initializing it from alterClause.
func (m *aiMerger) clausesOf(i int) []alterPart {
	if m.clauses[i] == nil && len(m.ops[i].alterParts) > 0 {
		m.clauses[i] = append([]alterPart(nil), m.ops[i].alterParts...)
	}
	if m.clauses[i] == nil {
		part := alterPart{clause: m.ops[i].alterClause}
		if m.ops[i].Type == OpAddColumn {
			part.addsColumn = m.ops[i].column
		}
		m.clauses[i] = []alterPart{part}
	}
	return m.clauses[i]
}

// consume removes op i from the plan and returns its clauses for re-emission elsewhere.
func (m *aiMerger) consume(i int) []alterPart {
	cl := m.clausesOf(i)
	m.consumed[i] = true
	return cl
}

// mergeInto moves the clauses of each src op into the target op, before (prepend) or after its
// current clauses. Warnings on consumed ops move to the target so none are lost.
func (m *aiMerger) mergeInto(target int, prepend bool, srcs ...int) {
	if len(srcs) == 0 {
		return
	}
	cur := m.clausesOf(target)
	var moved []alterPart
	for _, s := range srcs {
		moved = append(moved, m.consume(s)...)
		m.mergedInto[s] = target
		// The statement runs where the latest of its parts would, after what
		// each of them reads (orderColumnAdds).
		if src, dst := &m.ops[s], &m.ops[target]; src.Phase == dst.Phase && (src.Priority > dst.Priority || src.Priority == dst.Priority && src.sortName > dst.sortName) {
			dst.Priority, dst.sortName = src.Priority, src.sortName
		}
		if w := m.ops[s].Warning; w != "" && !strings.Contains(m.ops[target].Warning, w) {
			if m.ops[target].Warning != "" {
				m.ops[target].Warning += "; "
			}
			m.ops[target].Warning += w
		}
	}
	if prepend {
		m.clauses[target] = append(moved, cur...)
	} else {
		m.clauses[target] = append(cur, moved...)
	}
	m.touched[target] = true
}

// appendClause adds a synthesized clause (one not backed by an existing op) to the target.
func (m *aiMerger) appendClause(target int, clause string) {
	m.clauses[target] = append(m.clausesOf(target), alterPart{clause: clause})
	m.touched[target] = true
}

// result rebuilds the op slice: consumed ops disappear, and touched ops re-render their SQL
// from the merged clause list and keep it in alterParts.
func (m *aiMerger) result() []MigrationOp {
	out := make([]MigrationOp, 0, len(m.ops))
	for i, op := range m.ops {
		if m.consumed[i] {
			continue
		}
		if m.touched[i] {
			// An op already grouped (dropLastColumnsAfterAdds) renders its SQL
			// from alterParts, not alterClause, so strip that to keep the
			// prefix.
			rendered := op.alterClause
			if len(op.alterParts) > 0 {
				rendered = joinAlterParts(op.alterParts)
			}
			prefix := strings.TrimSuffix(op.SQL, rendered)
			op.SQL = prefix + joinAlterParts(m.clauses[i])
			op.alterParts = m.clauses[i]
		}
		out = append(out, op)
	}
	return out
}

// tableOpRoles indexes one table's ops by role. Only ops with a non-empty alterClause that
// suffix-matches their SQL are eligible.
type tableOpRoles struct {
	colAdd    map[string]int // lower column name → op index (ADD COLUMN)
	colModify map[string]int // lower column name → op index (MODIFY COLUMN)
	colDrop   map[string]int // lower column name → op index (DROP COLUMN)
	idxAdd    map[string]int // lower added-index name → op index (incl. combined PK change)
	idxDrop   map[string]int // lower dropped-index name → op index, PhasePre statements only
}

// collectRoles gathers the table's ops.
func (m *aiMerger) collectRoles(entry *TableDiffEntry) *tableOpRoles {
	r := &tableOpRoles{
		colAdd:    map[string]int{},
		colModify: map[string]int{},
		colDrop:   map[string]int{},
		idxAdd:    map[string]int{},
		idxDrop:   map[string]int{},
	}
	for _, i := range m.byTable[aiTableKey(entry.Database, entry.Name)] {
		op := m.ops[i]
		if m.consumed[i] {
			continue
		}
		// A statement grouping columns that read each other (rankColumnOps)
		// plays each of its columns' parts.
		if len(op.alterParts) > 0 {
			for _, part := range op.alterParts {
				if part.addsColumn != "" {
					r.colAdd[part.addsColumn] = i
				}
				if part.dropsColumn != "" {
					r.colDrop[part.dropsColumn] = i
				}
			}
			continue
		}
		if op.alterClause == "" || !strings.HasSuffix(op.SQL, op.alterClause) {
			continue
		}
		switch op.Type {
		case OpAddColumn, OpModifyColumn, OpDropColumn:
			name := op.column
			switch op.Type {
			case OpAddColumn:
				r.colAdd[name] = i
			case OpModifyColumn:
				r.colModify[name] = i
			default:
				r.colDrop[name] = i
			}
		default:
			if op.addsIndex != nil {
				r.idxAdd[toLower(op.addsIndex.Name)] = i
			}
			if op.dropsIndex != "" && op.addsIndex == nil && op.Phase == PhasePre {
				r.idxDrop[op.dropsIndex] = i
			}
		}
	}
	return r
}

// mergeTable runs the hazard analysis for one modified table: first the DROP direction (the last
// covering key leaving while the column is still AUTO_INCREMENT), then the ADD direction (the
// attribute arriving without a surviving covering key, or migrating from another column).
func (m *aiMerger) mergeTable(entry *TableDiffEntry, n *Normalizer) {
	fromAI := toLower(autoIncrementColumnName(entry.From))
	toAI := toLower(autoIncrementColumnName(entry.To))
	if fromAI == "" && toAI == "" {
		return
	}
	roles := m.collectRoles(entry)

	if fromAI != "" {
		m.mergeDropHazard(entry, roles, fromAI, toAI)
	}
	if toAI != "" {
		m.mergeAddHazard(entry, roles, fromAI, toAI, n)
	}
}

// mergeDropHazard handles the DROP direction: every from-side key covering the AUTO_INCREMENT
// column is dropped by a PhasePre statement, so the last drop would leave the column unkeyed. The
// covering drops are folded into whichever statement restores validity:
//
//	(a) the column's own DROP COLUMN (PhasePre);
//	(b) the column's de-AUTO_INCREMENT MODIFY (PhaseMain);
//	(c) a covering key ADD when the column stays AUTO_INCREMENT (PhaseMain).
//
// With no resolution the target itself is invalid and the ops stay ungrouped.
func (m *aiMerger) mergeDropHazard(entry *TableDiffEntry, roles *tableOpRoles, fromAI, toAI string) {
	covering := coveringIndexes(entry.From, fromAI, entry.From.Engine)
	if len(covering) == 0 {
		return
	}
	var drops []int
	for _, idx := range covering {
		di, ok := roles.idxDrop[coveredIndexKey(idx)]
		if !ok {
			return // a covering key survives PhasePre
		}
		drops = append(drops, di)
	}
	sort.Ints(drops)

	if di, ok := roles.colDrop[fromAI]; ok {
		// (a) `DROP INDEX k, DROP COLUMN c` in one PhasePre statement.
		m.mergeInto(di, true, drops...)
		return
	}
	if mi, ok := roles.colModify[fromAI]; ok && !columnToIsAutoIncrement(entry, fromAI) {
		// (b) `MODIFY COLUMN c <no AUTO_INCREMENT>, DROP INDEX k` in one PhaseMain statement.
		m.mergeInto(mi, false, drops...)
		m.pullDemotedPKMembers(entry, roles, mi, drops)
		return
	}
	if toAI == fromAI {
		// (c) the column stays AUTO_INCREMENT: fold the drops into a covering ADD.
		cand := m.chooseCoveringAddOp(entry.To, roles, fromAI)
		if cand < 0 {
			return
		}
		// The combined PK statement already ends in its ADD clause, so the drops follow it; a
		// standalone ADD takes them in front.
		prepend := m.ops[cand].dropsIndex == ""
		m.mergeInto(cand, prepend, drops...)
		// When the split path's PRIMARY KEY drop moved into this statement, its standalone ADD
		// PRIMARY KEY could otherwise run first, while the old PK still exists; pull it in too.
		if pkDropIn(m.ops, drops) {
			if pa, ok := roles.idxAdd[primaryIndexKey]; ok && pa != cand && !m.consumed[pa] {
				m.mergeInto(cand, false, pa)
			}
		}
		m.pullDemotedPKMembers(entry, roles, cand, drops)
	}
}

// pkDropIn reports whether any of the ops drops the PRIMARY KEY.
func pkDropIn(ops []MigrationOp, indices []int) bool {
	for _, i := range indices {
		if ops[i].dropsIndex == primaryIndexKey {
			return true
		}
	}
	return false
}

// pullDemotedPKMembers completes a merged statement that carries the PRIMARY KEY drop: an old PK
// member demoted to nullable in the same diff must not be modified while the PK still holds it,
// and the split-PK ordering no longer guarantees that once the PK drop moved into a PhaseMain
// statement. The demoting MODIFYs join the same statement.
func (m *aiMerger) pullDemotedPKMembers(entry *TableDiffEntry, roles *tableOpRoles, target int, drops []int) {
	if !pkDropIn(m.ops, drops) || m.ops[target].Phase != PhaseMain {
		return
	}
	oldPK := primaryKeyIndex(entry.From)
	if oldPK == nil {
		return
	}
	var demotes []int
	for _, ic := range oldPK.Columns {
		if ic == nil {
			continue
		}
		name := toLower(ic.Name)
		mi, ok := roles.colModify[name]
		if !ok || m.consumed[mi] || mi == target {
			continue
		}
		if ce := columnEntry(entry, name); ce != nil && ce.Action == DiffModify && ce.To != nil && ce.To.Nullable {
			demotes = append(demotes, mi)
		}
	}
	sort.Ints(demotes)
	m.mergeInto(target, false, demotes...)
}

// mergeAddHazard handles the ADD direction: the target AUTO_INCREMENT column has an ADD or
// MODIFY statement in the plan, and at its end no surviving key covers it. The backing key's ADD
// is folded into the column statement, together with:
//   - the previous AUTO_INCREMENT column's de-attribute MODIFY;
//   - the ADDs of any other new columns the backing key references.
//
// When no diffed key covers the column, a FK-implicit backing index on the target table is
// synthesized as a last resort; the deferred FK then reuses it instead of auto-creating a
// duplicate.
func (m *aiMerger) mergeAddHazard(entry *TableDiffEntry, roles *tableOpRoles, fromAI, toAI string, n *Normalizer) {
	target, hasColOp := roles.colAdd[toAI]
	if !hasColOp {
		target, hasColOp = roles.colModify[toAI]
	}
	if !hasColOp || m.consumed[target] {
		return
	}

	// The previous auto column's de-AUTO_INCREMENT must join this statement even when the new
	// column is already covered.
	if fromAI != "" && fromAI != toAI {
		if mi, ok := roles.colModify[fromAI]; ok && !m.consumed[mi] && mi != target &&
			!columnToIsAutoIncrement(entry, fromAI) {
			m.mergeInto(target, true, mi)
		}
	}

	// An add that runs ahead of the drops (dropLastColumnsAfterAdds) takes
	// the previous auto column's drop with it: a table holds one auto column.
	if fromAI != "" && fromAI != toAI && m.ops[target].Phase == PhasePre {
		if di, ok := roles.colDrop[fromAI]; ok && !m.consumed[di] && di != target {
			m.mergeInto(target, true, di)
		}
	}

	if m.addSideCovered(entry, roles, toAI, target) {
		return
	}

	cand := m.chooseCoveringAddOp(entry.To, roles, toAI)
	var candIdx *Index
	synthesized := ""
	if cand >= 0 {
		candIdx = m.ops[cand].addsIndex
	} else if idx := m.fkImplicitBackingToSynthesize(entry, roles, toAI, target); idx != nil {
		candIdx = idx
		synthesized = "ADD " + formatIndexDefinition(idx, n)
	}
	if candIdx == nil {
		return // no in-plan backing; the target is invalid and the server reports it
	}

	m.pullReferencedColumnOps(roles, target, candIdx, toAI, entry)
	if synthesized != "" {
		m.appendClause(target, synthesized)
	} else {
		m.mergeInto(target, false, cand)
	}
}

// addSideCovered reports whether a from-side key still covers the target auto column at the end
// of the target statement. The column must not be re-created in this plan, and a dropped key
// counts only when its DROP finally executes in a statement strictly after the target, resolved
// through the merge chain.
func (m *aiMerger) addSideCovered(entry *TableDiffEntry, roles *tableOpRoles, toAI string, target int) bool {
	if _, recreated := roles.colDrop[toAI]; recreated {
		return false
	}
	for _, idx := range coveringIndexes(entry.From, toAI, entry.To.Engine) {
		di, dropped := roles.idxDrop[coveredIndexKey(idx)]
		if !dropped {
			return true
		}
		final := m.resolveMergedTarget(di)
		if final != target && m.opExecutesBefore(target, final) {
			return true
		}
	}
	return false
}

// resolveMergedTarget follows the mergedInto chain to the statement that finally carries op i's
// clauses.
func (m *aiMerger) resolveMergedTarget(i int) int {
	for m.consumed[i] && m.mergedInto[i] >= 0 && m.mergedInto[i] != i {
		i = m.mergedInto[i]
	}
	return i
}

// opExecutesBefore reports whether op i's statement runs before op j's under sortMigrationOps'
// order: phase, priority, sortName, then original position.
func (m *aiMerger) opExecutesBefore(i, j int) bool {
	a, b := m.ops[i], m.ops[j]
	if a.Phase != b.Phase {
		return a.Phase < b.Phase
	}
	if a.Priority != b.Priority {
		return a.Priority < b.Priority
	}
	if a.sortName != b.sortName {
		return a.sortName < b.sortName
	}
	return i < j
}

// pullReferencedColumnOps moves the column statements the backing key depends on into the
// grouped statement, so the key clause sees those columns in their final shape:
//   - ADD COLUMN of every other new plain key-part column, then generated ones, each in target
//     column order;
//   - MODIFY COLUMN of every existing plain key-part column changed in this plan, because a
//     prefix part would otherwise be validated against the old type;
//   - when the key has generated or functional parts, all of the table's plain column adds,
//     because the expression may reference a new column that is not itself a key part, and for
//     functional parts the generated adds as well.
func (m *aiMerger) pullReferencedColumnOps(roles *tableOpRoles, target int, idx *Index, toAI string, entry *TableDiffEntry) {
	var plain, generated, modifies []string
	pullAllAdds, pullGeneratedAdds := false, false
	for _, ic := range idx.Columns {
		if ic == nil {
			continue
		}
		if ic.Expr != "" {
			pullAllAdds, pullGeneratedAdds = true, true
			continue
		}
		name := toLower(ic.Name)
		if name == toAI {
			continue
		}
		if _, ok := roles.colAdd[name]; ok {
			if col := entry.To.GetColumn(ic.Name); col != nil && col.Generated != nil {
				generated = append(generated, name)
				pullAllAdds = true
			} else {
				plain = append(plain, name)
			}
			continue
		}
		if _, ok := roles.colModify[name]; ok {
			modifies = append(modifies, name)
		}
	}
	if pullAllAdds {
		plain = plain[:0]
		for name, ai := range roles.colAdd {
			if name == toAI || m.consumed[ai] || ai == target {
				continue
			}
			col := entry.To.GetColumn(name)
			if col == nil {
				continue
			}
			if col.Generated == nil {
				plain = append(plain, name)
			} else if pullGeneratedAdds && !slices.Contains(generated, name) {
				generated = append(generated, name)
			}
		}
	}
	byPosition := func(names []string) {
		slices.SortFunc(names, func(a, b string) int {
			return columnPosition(entry.To, a) - columnPosition(entry.To, b)
		})
	}
	byPosition(modifies)
	byPosition(plain)
	byPosition(generated)
	for _, name := range modifies {
		mi := roles.colModify[name]
		if !m.consumed[mi] && mi != target {
			m.mergeInto(target, false, mi)
		}
	}
	for _, name := range append(plain, generated...) {
		ai := roles.colAdd[name]
		if !m.consumed[ai] && ai != target {
			m.mergeInto(target, false, ai)
		}
	}
}

// columnPosition returns the index of the named column in t.Columns, or -1.
func columnPosition(t *Table, name string) int {
	return slices.IndexFunc(t.Columns, func(c *Column) bool { return strings.EqualFold(c.Name, name) })
}

// chooseCoveringAddOp picks the plan's best covering key ADD for the auto column: PRIMARY over
// UNIQUE over plain, first-column coverage over any-position coverage, then index name. It
// returns -1 when no eligible ADD covers the column.
func (m *aiMerger) chooseCoveringAddOp(to *Table, roles *tableOpRoles, col string) int {
	names := make([]string, 0, len(roles.idxAdd))
	for name := range roles.idxAdd {
		names = append(names, name)
	}
	sort.Strings(names)
	best := -1
	var bestIdx *Index
	bestLeftmost := false
	for _, name := range names {
		oi := roles.idxAdd[name]
		if m.consumed[oi] {
			continue
		}
		idx := m.ops[oi].addsIndex
		if !indexBacksAutoColumn(idx, col, to.Engine) {
			continue
		}
		leftmost := indexFirstColumnIs(idx, col)
		if best < 0 || coveringAddPreferred(idx, leftmost, bestIdx, bestLeftmost) {
			best, bestIdx, bestLeftmost = oi, idx, leftmost
		}
	}
	return best
}

// coveringAddPreferred ranks candidate backing keys: first-column coverage beats any-position
// coverage, then PRIMARY > UNIQUE > plain. Names iterate in sorted order, so the first
// candidate of a rank wins.
func coveringAddPreferred(cand *Index, candLeftmost bool, cur *Index, curLeftmost bool) bool {
	if candLeftmost != curLeftmost {
		return candLeftmost
	}
	if cand.Primary != cur.Primary {
		return cand.Primary
	}
	if cand.Unique != cur.Unique {
		return cand.Unique
	}
	return false
}

// indexBacksAutoColumn reports whether the index satisfies the AUTO_INCREMENT key rule for the
// named column under the given storage engine: the column is the index's first plain column, or,
// for MyISAM, any plain column.
func indexBacksAutoColumn(idx *Index, col, engine string) bool {
	return indexFirstColumnIs(idx, col) ||
		(engineAllowsNonFirstAutoKey(engine) && indexHasPlainColumn(idx, col))
}

// coveringIndexes returns the table's indexes that back the named AUTO_INCREMENT column; removing
// all of them leaves it unkeyed.
func coveringIndexes(tbl *Table, col, engine string) []*Index {
	var out []*Index
	for _, idx := range tbl.Indexes {
		if idx == nil {
			continue
		}
		if indexBacksAutoColumn(idx, col, engine) {
			out = append(out, idx)
		}
	}
	return out
}

// coveredIndexKey is the dropsIndex key of an index.
func coveredIndexKey(idx *Index) string {
	if idx.Primary {
		return primaryIndexKey
	}
	return toLower(idx.Name)
}

// engineAllowsNonFirstAutoKey reports whether the storage engine accepts an AUTO_INCREMENT
// column in a non-first position of a multi-column key. Only MyISAM does.
func engineAllowsNonFirstAutoKey(engine string) bool {
	return strings.EqualFold(strings.TrimSpace(engine), "MyISAM")
}

// indexHasPlainColumn reports whether any plain (non-expression) key part is the named column.
func indexHasPlainColumn(idx *Index, col string) bool {
	for _, ic := range idx.Columns {
		if ic != nil && ic.Expr == "" && strings.EqualFold(ic.Name, col) {
			return true
		}
	}
	return false
}

// columnToIsAutoIncrement reports whether the diff leaves the named column AUTO_INCREMENT on the
// target side; a column without a diff entry is unchanged.
func columnToIsAutoIncrement(entry *TableDiffEntry, col string) bool {
	if ce := columnEntry(entry, col); ce != nil {
		return ce.To != nil && ce.To.AutoIncrement
	}
	c := entry.To.GetColumn(col)
	return c != nil && c.AutoIncrement
}

// columnEntry returns the diff's column entry for the lower-cased name, or nil.
func columnEntry(entry *TableDiffEntry, col string) *ColumnDiffEntry {
	for i := range entry.Columns {
		if toLower(entry.Columns[i].Name) == col {
			return &entry.Columns[i]
		}
	}
	return nil
}

// fkImplicitBackingToSynthesize returns the target table's FK-implicit backing index to inline
// into the grouped column statement when no diffed key covers the new auto column. The index
// must cover the column, be FK-implicit on the target (so the index differ emits no ADD), and
// have a free name when the grouped statement runs: absent on the from side, or freed by an
// in-plan DROP that executes at or before the target. Candidates iterate in name order.
func (m *aiMerger) fkImplicitBackingToSynthesize(entry *TableDiffEntry, roles *tableOpRoles, col string, target int) *Index {
	implicit := fkImplicitIndexNames(entry.To)
	if len(implicit) == 0 {
		return nil
	}
	fromNames := make(map[string]bool, len(entry.From.Indexes))
	for _, idx := range entry.From.Indexes {
		if idx != nil {
			fromNames[toLower(idx.Name)] = true
		}
	}
	nameFree := func(name string) bool {
		if !fromNames[name] {
			return true
		}
		di, dropped := roles.idxDrop[name]
		if !dropped {
			return false
		}
		final := m.resolveMergedTarget(di)
		return final == target || m.opExecutesBefore(final, target)
	}
	var best *Index
	for _, idx := range entry.To.Indexes {
		if idx == nil || idx.Primary || !implicit[toLower(idx.Name)] || !nameFree(toLower(idx.Name)) {
			continue
		}
		if !indexBacksAutoColumn(idx, col, entry.To.Engine) {
			continue
		}
		if best == nil || toLower(idx.Name) < toLower(best.Name) {
			best = idx
		}
	}
	return best
}
