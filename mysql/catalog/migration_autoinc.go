package catalog

import (
	"sort"
	"strings"
)

// MySQL SDL generate — AUTO_INCREMENT / backing-key statement grouping.
//
// MySQL validates, at the END OF EVERY ALTER TABLE statement, that the table holds at most one
// AUTO_INCREMENT column and that this column is keyed — the FIRST column of some index (errno
// 1075, "there can be only one auto column and it must be defined as a key"). The plan generator
// otherwise emits ONE clause per statement in a fixed phase/priority order (index drops →
// column drops → column adds/modifies → index adds), so a change that is only valid when the
// AUTO_INCREMENT column and its backing key move TOGETHER fails on whichever ungrouped
// statement exposes the unkeyed (or doubled) auto column. formatCreateTable already inlines the
// backing key for a NEW table (autoIncSupportingIndex); this pass is the ALTER-path analog for
// MODIFIED tables: it detects the hazard statements and merges the involved clauses into one
// ALTER TABLE.
//
// Oracle evidence (probed live on 5.7.25 and 8.0.32, identical behavior; each grouped form this
// pass emits is proven by migration_autoinc_oracle_test.go):
//   - `ADD COLUMN c ... AUTO_INCREMENT` alone → errno 1075; with `, ADD <any key with c first>`
//     in the SAME statement → OK (UNIQUE, plain KEY, PRIMARY KEY; clause order is free).
//   - InnoDB and MEMORY require c to be the FIRST column of the key (a non-first position fails
//     even grouped); MyISAM alone also accepts a non-first position. A functional (expression)
//     first part does NOT satisfy the rule; DESC and INVISIBLE keys DO.
//   - `MODIFY COLUMN c ... AUTO_INCREMENT` (gaining the attribute) is subject to the same rule;
//     it works ungrouped only when a surviving key already covers c.
//   - Making one column AUTO_INCREMENT while another still holds the attribute fails ("only one
//     auto column"); the de-AUTO_INCREMENT MODIFY and the gaining clause must share a statement.
//   - Dropping the LAST key covering a still-AUTO_INCREMENT column → errno 1075 on the DROP
//     statement; grouped with the column's DROP, its de-AUTO_INCREMENT MODIFY, or a replacement
//     covering ADD (any of which restores validity at end of statement) → OK.
//   - The combined `DROP PRIMARY KEY, ADD PRIMARY KEY` statement accepts extra clauses (a new
//     AUTO_INCREMENT column's ADD, a stale covering key's DROP).
//   - A demoted old-PK member's `MODIFY ... NULL` may share the statement that drops the PK.
//
// The pass is deliberately RULE-BASED (grouping exactly the hazard cases) rather than a general
// per-table clause merger: every other statement keeps its own op, so op granularity, warnings,
// and the global sort order stay untouched for non-hazard plans. When a hazard has no in-plan
// resolution (the TARGET itself leaves an AUTO_INCREMENT column unkeyed — invalid in MySQL
// regardless of statement grouping), the ops are left ungrouped and the apply surfaces MySQL's
// own error, matching autoIncSupportingIndex's stance on unkeyable CREATEs.
//
// Ops participate only through the grouping metadata their constructors set (alterClause,
// addsIndex, dropsIndex on MigrationOp); everything else — checks, views, options, FK
// constraint adds — is invisible to the pass and never regrouped.
//
// FLAGGED LIMITATION (not handled): a storage-ENGINE conversion in the same plan. The
// table-option ALTER (ENGINE=...) runs at PhaseMain/priorityTable, before every column/index
// statement, and carries no grouping metadata — so converting a MyISAM table whose auto column
// is covered only in a NON-first key position to InnoDB fails errno 1075 on the ENGINE
// statement itself, even when the plan adds a first-position key later. Fixing it would mean
// grouping table-option clauses (a different op family) or hoisting index adds above
// priorityTable; both are out of this pass's contract. Such a plan also failed before the pass
// existed — nothing regresses.

// primaryIndexKey is the dropsIndex marker for the PRIMARY KEY, which has no user-facing name.
// MySQL reserves the name PRIMARY for it, so the lower-cased form cannot collide with a
// secondary index.
const primaryIndexKey = "primary"

// mergeAutoIncrementKeyOps rewrites the generated ops so that no statement boundary exposes an
// unkeyed or doubled AUTO_INCREMENT column (see the file comment). It returns the ops slice to
// use; ops merged into another statement are removed, the surviving statement carrying their
// clauses (and warnings). Called by GenerateMigrationWithNormalizer before sortMigrationOps.
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
// and mergedInto[i] records which statement (so later analysis knows WHEN the clauses now run).
// byTable groups the op indices by owning table once, so the per-table role collection does not
// rescan the whole plan for every AUTO_INCREMENT-bearing table.
type aiMerger struct {
	ops        []MigrationOp
	clauses    [][]string
	consumed   []bool
	touched    []bool
	mergedInto []int
	byTable    map[string][]int
}

func newAIMerger(ops []MigrationOp) *aiMerger {
	m := &aiMerger{
		ops:        ops,
		clauses:    make([][]string, len(ops)),
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

// aiTableKey is the collision-free per-table grouping key (encodeKeyFields length-prefixes the
// fields, so identifiers containing separator characters cannot alias another table).
func aiTableKey(database, table string) string {
	return encodeKeyFields("db", database, "table", table)
}

// clausesOf returns op i's current clause list, initializing it from alterClause.
func (m *aiMerger) clausesOf(i int) []string {
	if m.clauses[i] == nil {
		m.clauses[i] = []string{m.ops[i].alterClause}
	}
	return m.clauses[i]
}

// consume removes op i from the plan and returns its clauses for re-emission elsewhere.
func (m *aiMerger) consume(i int) []string {
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
	var moved []string
	for _, s := range srcs {
		moved = append(moved, m.consume(s)...)
		m.mergedInto[s] = target
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
	m.clauses[target] = append(m.clausesOf(target), clause)
	m.touched[target] = true
}

// result rebuilds the op slice: consumed ops disappear; touched ops re-render their SQL from
// the merged clause list (the original "ALTER TABLE <ident> " prefix is preserved byte-exactly
// by trimming the op's own alterClause suffix).
func (m *aiMerger) result() []MigrationOp {
	out := make([]MigrationOp, 0, len(m.ops))
	for i, op := range m.ops {
		if m.consumed[i] {
			continue
		}
		if m.touched[i] {
			prefix := strings.TrimSuffix(op.SQL, op.alterClause)
			op.SQL = prefix + strings.Join(m.clauses[i], ", ")
		}
		out = append(out, op)
	}
	return out
}

// tableOpRoles indexes one table's ops by their role in the hazard analysis. Only ops carrying
// grouping metadata (a non-empty alterClause that suffix-matches their SQL) are eligible.
type tableOpRoles struct {
	colAdd    map[string]int // lower column name → op index (ADD COLUMN)
	colModify map[string]int // lower column name → op index (MODIFY COLUMN)
	colDrop   map[string]int // lower column name → op index (DROP COLUMN)
	idxAdd    map[string]int // lower added-index name → op index (incl. combined PK change)
	idxDrop   map[string]int // lower dropped-index name → op index, PhasePre statements only
}

// collectRoles gathers the table's ops. Column ops are matched through the diff's column
// entries (recomputing the op's sortName — MySQL identifiers cannot contain '.', so the dotted
// key is unambiguous); index ops carry their identity in addsIndex/dropsIndex.
func (m *aiMerger) collectRoles(entry *TableDiffEntry) *tableOpRoles {
	r := &tableOpRoles{
		colAdd:    map[string]int{},
		colModify: map[string]int{},
		colDrop:   map[string]int{},
		idxAdd:    map[string]int{},
		idxDrop:   map[string]int{},
	}
	colSort := make(map[string]string, len(entry.Columns))
	for _, ce := range entry.Columns {
		colSort[columnSortName(entry.Database, entry.Name, ce.Name)] = toLower(ce.Name)
	}
	for _, i := range m.byTable[aiTableKey(entry.Database, entry.Name)] {
		op := m.ops[i]
		if m.consumed[i] {
			continue
		}
		if op.alterClause == "" || !strings.HasSuffix(op.SQL, op.alterClause) {
			continue
		}
		switch op.Type {
		case OpAddColumn, OpModifyColumn, OpDropColumn:
			name, ok := colSort[op.sortName]
			if !ok {
				continue
			}
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

// mergeTable runs the hazard analysis for one modified table: first the DROP direction (the
// last covering key leaving while the column is still AUTO_INCREMENT), then the ADD direction
// (the AUTO_INCREMENT attribute arriving without a surviving covering key, or migrating from
// another column).
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
// column is dropped by a PhasePre statement, so the LAST such drop would end its statement with
// an unkeyed auto column (errno 1075, oracle-probed). The covering drops are folded into
// whichever statement restores validity:
//
//	(a) the column's own DROP COLUMN (PhasePre — the drops arrive with the column removal);
//	(b) the column's de-AUTO_INCREMENT MODIFY (PhaseMain — the attribute leaves with the keys);
//	(c) a covering key ADD when the column stays AUTO_INCREMENT (PhaseMain — the coverage swaps
//	    atomically; the secondary-key analog of the combined PRIMARY KEY change).
//
// With no resolution the target itself is invalid (an unkeyed AUTO_INCREMENT column) and the
// ops stay ungrouped — MySQL reports the error.
func (m *aiMerger) mergeDropHazard(entry *TableDiffEntry, roles *tableOpRoles, fromAI, toAI string) {
	covering := coveringIndexes(entry.From, fromAI, entry.From.Engine)
	if len(covering) == 0 {
		return
	}
	var drops []int
	for _, idx := range covering {
		di, ok := roles.idxDrop[coveredIndexKey(idx)]
		if !ok {
			return // a covering key survives PhasePre — every drop statement stays valid
		}
		drops = append(drops, di)
	}
	sort.Ints(drops)

	if di, ok := roles.colDrop[fromAI]; ok {
		// (a) `DROP INDEX k, DROP COLUMN c` — one PhasePre statement at the column drop's slot.
		m.mergeInto(di, true, drops...)
		return
	}
	if mi, ok := roles.colModify[fromAI]; ok && !columnToIsAutoIncrement(entry, fromAI) {
		// (b) `MODIFY COLUMN c <no AUTO_INCREMENT>, DROP INDEX k` — one PhaseMain statement.
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
		// The combined PK statement already ends in its ADD clause (append the stale drops
		// after it — the oracle-proven shape); a standalone ADD takes the drops in front.
		prepend := m.ops[cand].dropsIndex == ""
		m.mergeInto(cand, prepend, drops...)
		// When the split-path PRIMARY KEY drop moved into this statement, the split's
		// standalone ADD PRIMARY KEY (same priority, name-sorted) could otherwise run BEFORE
		// it — while the old PK still exists (errno 1068). Pull the replacement PK add into
		// the same statement (oracle-proven single-statement form).
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

// pullDemotedPKMembers completes a merged statement that carries the PRIMARY KEY drop: an old
// PK member demoted to nullable in the same diff must not be MODIFYed while the PK still holds
// it (errno 1171), and the split-PK ordering that normally guarantees this no longer applies
// once the PK drop moved into a PhaseMain statement. Folding the demoting MODIFYs into the same
// statement satisfies both rules at once (oracle-probed on both versions).
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

// mergeAddHazard handles the ADD direction: the table's target AUTO_INCREMENT column has an
// ADD/MODIFY statement in the plan, and at that statement's boundary no surviving key covers it
// (a brand-new column is never covered; a modified one only by from-keys that no PhasePre
// statement dropped). The backing key's ADD is folded into the column statement, together with:
//   - the previous AUTO_INCREMENT column's de-attribute MODIFY (two live auto columns are
//     rejected even transiently, oracle-probed);
//   - the ADDs of any other new columns the backing key references (plain before generated —
//     a generated column cannot reference the auto column itself, so this order always works).
//
// When no diffed key covers the column, a FK-implicit backing index on the target table is
// synthesized as a last resort — the ALTER analog of formatCreateTable's last-resort inline;
// the deferred PhasePost FK then reuses it instead of auto-creating a duplicate.
func (m *aiMerger) mergeAddHazard(entry *TableDiffEntry, roles *tableOpRoles, fromAI, toAI string, n *Normalizer) {
	target, hasColOp := roles.colAdd[toAI]
	if !hasColOp {
		target, hasColOp = roles.colModify[toAI]
	}
	if !hasColOp || m.consumed[target] {
		return
	}

	// The de-AUTO_INCREMENT of a different previous auto column must join this statement even
	// when the new column is already covered (the "only one auto column" half of errno 1075).
	if fromAI != "" && fromAI != toAI {
		if mi, ok := roles.colModify[fromAI]; ok && !m.consumed[mi] && mi != target &&
			!columnToIsAutoIncrement(entry, fromAI) {
			m.mergeInto(target, true, mi)
		}
	}

	if m.addSideCovered(entry, roles, toAI) {
		return
	}

	cand := m.chooseCoveringAddOp(entry.To, roles, toAI)
	var candIdx *Index
	synthesized := ""
	if cand >= 0 {
		candIdx = m.ops[cand].addsIndex
	} else if idx := fkImplicitBackingToSynthesize(entry, toAI); idx != nil {
		candIdx = idx
		synthesized = "ADD " + formatIndexDefinition(idx, n)
	}
	if candIdx == nil {
		return // no in-plan backing — the target is invalid and MySQL reports it
	}

	m.pullReferencedColumnAdds(roles, target, candIdx, toAI, entry)
	if synthesized != "" {
		m.appendClause(target, synthesized)
	} else {
		m.mergeInto(target, false, cand)
	}
}

// addSideCovered reports whether a from-side key still covers the target auto column at its
// column statement: the column must not be re-created in this plan (its own PhasePre DROP
// strips it from every index) and the key must still exist when the PhaseMain column statement
// runs — either no statement drops it, or its drop was folded into a PhaseMain statement by the
// DROP-direction merges (a drop folded into a PhasePre statement, e.g. another column's DROP,
// still removes the key before PhaseMain).
func (m *aiMerger) addSideCovered(entry *TableDiffEntry, roles *tableOpRoles, toAI string) bool {
	if _, recreated := roles.colDrop[toAI]; recreated {
		return false
	}
	for _, idx := range coveringIndexes(entry.From, toAI, entry.To.Engine) {
		di, dropped := roles.idxDrop[coveredIndexKey(idx)]
		if !dropped {
			return true
		}
		if m.consumed[di] && m.mergedInto[di] >= 0 && m.ops[m.mergedInto[di]].Phase == PhaseMain {
			return true
		}
	}
	return false
}

// pullReferencedColumnAdds moves the ADD COLUMN statements of the backing key's OTHER new
// columns into the grouped statement (the key clause must be able to reference them), plain
// columns first, then generated ones, each name-ordered — matching the standalone emission
// order of the column generator. When a GENERATED key part is pulled, ALL of the table's plain
// column adds come along: the generated expression may reference a new plain column that is not
// itself a key part, and leaving that column's ADD in its own name-sorted statement could run
// it after this one (extra plain ADD clauses in the grouped statement are harmless,
// oracle-proven).
func (m *aiMerger) pullReferencedColumnAdds(roles *tableOpRoles, target int, idx *Index, toAI string, entry *TableDiffEntry) {
	var plain, generated []string
	for _, ic := range idx.Columns {
		if ic == nil || ic.Expr != "" {
			continue
		}
		name := toLower(ic.Name)
		if name == toAI {
			continue
		}
		if _, ok := roles.colAdd[name]; !ok {
			continue
		}
		if col := entry.To.GetColumn(ic.Name); col != nil && col.Generated != nil {
			generated = append(generated, name)
		} else {
			plain = append(plain, name)
		}
	}
	if len(generated) > 0 {
		plain = plain[:0]
		for name, ai := range roles.colAdd {
			if name == toAI || m.consumed[ai] || ai == target {
				continue
			}
			if col := entry.To.GetColumn(name); col != nil && col.Generated == nil {
				plain = append(plain, name)
			}
		}
	}
	sort.Strings(plain)
	sort.Strings(generated)
	for _, name := range append(plain, generated...) {
		ai := roles.colAdd[name]
		if !m.consumed[ai] && ai != target {
			m.mergeInto(target, false, ai)
		}
	}
}

// chooseCoveringAddOp picks the plan's best covering key ADD for the auto column: PRIMARY over
// UNIQUE over plain, first-column coverage over MyISAM's any-position coverage, index name as
// the deterministic tie-break. Returns -1 when no eligible ADD covers the column.
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
// coverage, then PRIMARY > UNIQUE > plain (names already iterate in sorted order, so the first
// candidate of a rank wins the tie).
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

// indexBacksAutoColumn reports whether the index satisfies errno 1075 for the named column
// under the given storage engine: the column is the index's FIRST plain column, or — for MyISAM
// only (oracle-probed; InnoDB and MEMORY reject it) — any plain column position.
func indexBacksAutoColumn(idx *Index, col, engine string) bool {
	return indexFirstColumnIs(idx, col) ||
		(engineAllowsNonFirstAutoKey(engine) && indexHasPlainColumn(idx, col))
}

// coveringIndexes returns the table's indexes that back the named AUTO_INCREMENT column
// (indexBacksAutoColumn) — the keys whose complete removal strands it (errno 1075).
func coveringIndexes(tbl *Table, col, engine string) []*Index {
	var out []*Index
	for _, idx := range tbl.Indexes {
		if idx == nil || isGeneratedInvisiblePrimaryKeyIndex(idx) {
			continue
		}
		if indexBacksAutoColumn(idx, col, engine) {
			out = append(out, idx)
		}
	}
	return out
}

// coveredIndexKey is the dropsIndex-space key of an index (the PRIMARY KEY has no user name).
func coveredIndexKey(idx *Index) string {
	if idx.Primary {
		return primaryIndexKey
	}
	return toLower(idx.Name)
}

// engineAllowsNonFirstAutoKey reports whether the storage engine accepts an AUTO_INCREMENT
// column in a NON-first position of a multi-column key (MyISAM's grouped-sequence feature).
// InnoDB and MEMORY require the first position (oracle-probed on 5.7.25 and 8.0.32).
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

// columnToIsAutoIncrement reports whether the diff leaves the named column AUTO_INCREMENT on
// the target side. A column without a diff entry is unchanged, so its target state is its
// current state.
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
// into the grouped column statement when NO diffed key covers the new auto column — the ALTER
// analog of formatCreateTable's last-resort inline. The index must cover the column, be
// FK-implicit on the target (the index differ excludes it, so no ADD op exists), and its name
// must be free on the from side (otherwise adding it would collide, errno 1061). The deferred
// PhasePost FK ADD then reuses the pre-created index instead of auto-creating a duplicate
// (oracle-probed). Candidates iterate in name order for determinism.
func fkImplicitBackingToSynthesize(entry *TableDiffEntry, col string) *Index {
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
	var best *Index
	for _, idx := range entry.To.Indexes {
		if idx == nil || idx.Primary || !implicit[toLower(idx.Name)] || fromNames[toLower(idx.Name)] {
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
