package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — foreign keys.
//
// This is the MySQL analog of PG's foreign-key generator. It turns the FK diff
// (TableDiffEntry.ForeignKeys, populated by diff_foreign_key.go) into ALTER TABLE ... ADD/DROP
// FOREIGN KEY DDL. MySQL folds FK changes into ALTER TABLE; a FK cannot be altered in place, so
// a modified FK is a DROP followed by an ADD.
//
// Phase ordering — the crux of coordinating with the index node AND with the table/column
// generators. FK ADDs and FK DROPs sit in OPPOSITE phases:
//
//   - FK ADDs run in PhasePost (priorityForeignKey=99), AFTER the index node's PhaseMain index
//     ADDs (priorityIndex=30) and after table/column DDL. MySQL's ADD CONSTRAINT ... FOREIGN KEY
//     auto-creates a backing index ONLY when no covering index already exists; by deferring the
//     FK add to PhasePost, an index the user declared (added by the index node in PhaseMain) is
//     in place first and is REUSED — no duplicate index, no errno 1061/1553. It also guarantees
//     every referenced table/column exists before the FK is created (the reason PG defers FKs).
//   - FK DROPs run in PhasePre, BEFORE every dependent drop: a referenced table (errno 3730), a
//     referencing column (errno 1828), and the FK's own backing index (errno 1553). MySQL refuses
//     to drop any object a live FK still needs, so the FK constraint must be released FIRST. FK
//     constraint drops therefore use priorityForeignKeyConstraintDrop, ordered ahead of the
//     table-drop / index-drop / column-drop priorities the other generators use in PhasePre.
//
// Backing-index lifecycle on DROP — owned HERE, by contract with the index node:
//
//	When a FK is dropped, MySQL LEAVES the auto-created backing index behind (verified on both
//	live engines: after `ALTER TABLE c DROP FOREIGN KEY fk`, `KEY fk (pid)` remains). The index
//	node EXCLUDES FK-implicit backing indexes from its diff (diff_index.go fkImplicitIndexNames),
//	so it emits no DROP for that leftover — this node owns it. The leftover DROP INDEX runs in
//	PhasePre at priorityForeignKeyBackingIndexDrop: AFTER all FK constraint drops (so a backing
//	index SHARED by several FKs — MySQL keeps one index for FKs on the same columns — is not
//	dropped while another of those FKs still references it, errno 1553) and BEFORE column drops (a
//	DROP COLUMN auto-removes an index it empties, after which an explicit DROP INDEX fails errno
//	1091 — same hazard the index node guards with priorityIndexDrop). We drop the leftover index
//	only when (a) it was FK-implicit on the `from` side for `con`, (b) NO foreign key SURVIVING in
//	`to` still needs it (no `to` FK whose columns it covers), and (c) the target keeps no USER
//	index of that name — see leftoverBackingIndexToDrop.
//
// Rendering routes through the same canonical form show.go uses for the stored FK line, so the
// readback of the emitted DDL canonicalizes equal to `to` and apply-correctness holds.
//
// A per-plan seenBackingDrop set (keyed by the collision-free encodeKeyFields(db, table, index),
// so it never conflates indexes on different tables — and stays unambiguous even when an identifier
// contains a literal `.`) deduplicates the leftover-backing-index DROP INDEX: MySQL keeps ONE
// backing index for several FKs on the same leading columns (two unnamed FKs on `pid` share
// `KEY pid (pid)`), so when both such FKs are dropped, each would otherwise select the same index
// and emit a duplicate DROP INDEX — the second failing errno 1091. The set keeps exactly one drop
// per (database, table, index).
func generateForeignKeyDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp
	seenBackingDrop := map[string]bool{}

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addForeignKeyOpsForNewTable(entry)...)
		case DiffModify:
			ops = append(ops, foreignKeyOpsForModifiedTable(entry, seenBackingDrop)...)
		case DiffDrop:
			// The table itself is being dropped, but its FKs must be RELEASED first: a FK on this
			// table that references ANOTHER table being dropped in the same plan would otherwise
			// block that table's DROP (errno 3730 / 1451) — table drops all share priorityTable, so
			// the referenced parent can sort before this child. Emit DROP FOREIGN KEY for every FK
			// here in PhasePre (priorityForeignKeyConstraintDrop), ahead of all table drops. No
			// backing-index drop is needed — DROP TABLE removes the indexes.
			ops = append(ops, dropForeignKeyConstraintOpsForDroppedTable(entry)...)
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Priority != ops[j].Priority {
			return ops[i].Priority < ops[j].Priority
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// dropForeignKeyConstraintOpsForDroppedTable emits a DROP FOREIGN KEY for every FK on a table that
// is itself being DROPPED, in PhasePre ahead of all table drops. This releases the references so a
// referenced table dropped in the same plan is not blocked (errno 3730 / 1451). It deliberately
// emits NO backing-index drop: DROP TABLE removes the indexes, and the table is gone regardless.
func dropForeignKeyConstraintOpsForDroppedTable(entry *TableDiffEntry) []MigrationOp {
	if entry.From == nil {
		return nil
	}
	table := tableIdent(entry.From)
	var ops []MigrationOp
	for _, con := range orderedForeignKeys(entry.From) {
		ops = append(ops, MigrationOp{
			Type:         OpDropForeignKey,
			Database:     entry.Database,
			ObjectName:   entry.Name,
			ParentObject: entry.Name,
			SQL:          fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, migrationQuoteIdent(con.Name)),
			Phase:        PhasePre,
			Priority:     priorityForeignKeyConstraintDrop,
			sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
		})
	}
	return ops
}

// priorityForeignKeyConstraintDrop orders FK CONSTRAINT drops at the very front of PhasePre —
// ahead of table drops (priorityTable=10), index drops (priorityIndexDrop=15), and column drops
// (priorityColumn-1=19 / priorityColumn=20). A live FK blocks dropping any object it needs (a
// referenced table → errno 3730, a referencing column → errno 1828, its backing index → errno
// 1553), so every FK constraint must be released before any of those drops run. It is negative so
// it sorts before priorityTable without renumbering the shared constants.
const priorityForeignKeyConstraintDrop = -10

// priorityForeignKeyBackingIndexDrop orders the leftover-backing-index drop AFTER all FK
// constraint drops (so a backing index shared by multiple same-column FKs survives until the last
// of them is gone, errno 1553) and BEFORE column drops (a DROP COLUMN auto-removes the index,
// after which an explicit DROP INDEX fails errno 1091). It sits just below priorityIndexDrop so it
// shares the index node's "drop indexes before columns" guarantee.
const priorityForeignKeyBackingIndexDrop = priorityIndexDrop - 1

// addForeignKeyOpsForNewTable emits ADD CONSTRAINT ... FOREIGN KEY ops for every FK on a freshly
// created table. formatCreateTable (migration_table.go) renders only columns and the inline
// PRIMARY KEY — it does NOT inline FK constraints — so every FK on a new table is added here in
// PhasePost, after the table (and any referenced tables) exist.
func addForeignKeyOpsForNewTable(entry *TableDiffEntry) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, con := range orderedForeignKeys(entry.To) {
		ops = append(ops, addForeignKeyOp(entry, table, con))
	}
	return ops
}

// foreignKeyOpsForModifiedTable emits ADD / DROP / (DROP+ADD) ops from the per-FK diff of a
// modified table. A MODIFY is a DROP followed by an ADD (an FK cannot be altered in place): the
// DROP runs in PhasePre (priorityForeignKeyConstraintDrop) and the ADD in PhasePost
// (priorityForeignKey), so the old constraint is fully released before the new one is created —
// the name is free and any backing index the new FK needs is settled by then.
func foreignKeyOpsForModifiedTable(entry *TableDiffEntry, seenBackingDrop map[string]bool) []MigrationOp {
	if len(entry.ForeignKeys) == 0 {
		return nil
	}
	var ops []MigrationOp
	for _, fe := range entry.ForeignKeys {
		switch fe.Action {
		case DiffAdd:
			if fe.To == nil || entry.To == nil {
				continue
			}
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		case DiffDrop:
			if fe.From == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOps(entry, fe.From, seenBackingDrop)...)
		case DiffModify:
			if fe.From == nil || fe.To == nil || entry.To == nil {
				continue
			}
			// DROP the old FK (and its leftover backing index, if FK-implicit and no longer
			// needed), then ADD the new one. The DROP is in PhasePre, the ADD in PhasePost, so the
			// drop always precedes the add.
			ops = append(ops, dropForeignKeyOps(entry, fe.From, seenBackingDrop)...)
			ops = append(ops, addForeignKeyOp(entry, tableIdent(entry.To), fe.To))
		}
	}
	return ops
}

// addForeignKeyOp builds an ALTER TABLE ... ADD CONSTRAINT `name` FOREIGN KEY (...) REFERENCES
// ... op. It runs in PhasePost at priorityForeignKey so it lands after table/column/index DDL —
// every referenced table exists and any user-declared covering index is already present (so MySQL
// reuses it instead of auto-creating a duplicate backing index).
func addForeignKeyOp(entry *TableDiffEntry, table string, con *Constraint) MigrationOp {
	return MigrationOp{
		Type:         OpAddForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ADD %s", table, formatForeignKeyDefinition(con)),
		Phase:        PhasePost,
		Priority:     priorityForeignKey,
		sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
	}
}

// dropForeignKeyOps builds the ops to drop a foreign key: an ALTER TABLE ... DROP FOREIGN KEY
// `name` (PhasePre, priorityForeignKeyConstraintDrop — ahead of every dependent drop), plus —
// when the FK leaves behind an auto-created backing index nothing surviving needs — an ALTER
// TABLE ... DROP INDEX `name` for that leftover (PhasePre, priorityForeignKeyBackingIndexDrop —
// after ALL FK constraint drops, before column drops). The two priorities (not a sortName
// suffix) are what guarantee the constraint drop precedes the index drop globally, even across
// multiple FKs sharing one backing index.
//
// The leftover index is dropped ONLY when leftoverBackingIndexToDrop says so: it was FK-implicit
// on `from` for `con`, no FK surviving in `to` still needs it, and the target keeps no USER index
// of that name. seenBackingDrop deduplicates that DROP INDEX across FKs SHARING one backing index
// (two unnamed FKs on `pid` share `KEY pid (pid)`; without dedup each FK drop would emit
// `DROP INDEX pid` and the second fail errno 1091). The set is keyed by (db.table.index).
func dropForeignKeyOps(entry *TableDiffEntry, con *Constraint, seenBackingDrop map[string]bool) []MigrationOp {
	if entry.From == nil {
		return nil
	}
	table := tableIdent(entry.From)
	ops := []MigrationOp{{
		Type:         OpDropForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, migrationQuoteIdent(con.Name)),
		Phase:        PhasePre,
		Priority:     priorityForeignKeyConstraintDrop,
		sortName:     foreignKeySortName(entry.Database, entry.Name, con.Name),
	}}

	if idxName, ok := leftoverBackingIndexToDrop(entry, con); ok {
		// Dedup the leftover-backing-index DROP across FKs sharing one index using a COLLISION-FREE
		// key (encodeKeyFields length-prefixes each field), not the dotted db.table.index string:
		// the dotted form is ambiguous when an identifier itself contains a `.` (e.g. db `d`, table
		// `a`, index `b.c` and db `d`, table `a.b`, index `c` both fold to `d.a.b.c`), which would
		// suppress a legitimate second DROP. The op's sortName stays the dotted form so ordinary
		// identifiers keep their existing ordering.
		dedupKey := encodeKeyFields("db", entry.Database, "table", entry.Name, "index", idxName)
		if !seenBackingDrop[dedupKey] {
			seenBackingDrop[dedupKey] = true
			ops = append(ops, MigrationOp{
				// Op-type is OpDropForeignKey (NOT OpDropIndex) even though the SQL is DROP INDEX:
				// this index drop is part of the FK's lifecycle, owned by the FK node. Tagging it
				// OpDropForeignKey keeps it out of the index node's op-type space, which its contract
				// tests assert is empty on FK-only changes (no OpAddIndex/OpDropIndex). Nothing keys
				// off OpDropForeignKey expecting only literal DROP FOREIGN KEY statements — the only
				// consumer is MigrationPlan.SQL(), which concatenates op.SQL verbatim.
				Type:         OpDropForeignKey,
				Database:     entry.Database,
				ObjectName:   entry.Name,
				ParentObject: entry.Name,
				SQL:          fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, migrationQuoteIdent(idxName)),
				Phase:        PhasePre,
				Priority:     priorityForeignKeyBackingIndexDrop,
				sortName:     foreignKeySortName(entry.Database, entry.Name, idxName),
			})
		}
	}
	return ops
}

// leftoverBackingIndexToDrop decides whether dropping `con` (a FK on the `from` table) leaves an
// auto-created backing index this node must drop, and returns its actual (original-case) name.
// It returns the index iff ALL of:
//   - the `from` table has an FK-implicit backing index for `con` (the index node's exact
//     detection: fkImplicitIndexNames over `from` + isAutoBackingIndexName + isPlainBackingIndexFor);
//   - NO foreign key SURVIVING in `to` still needs that index — i.e. no `to` FK whose columns the
//     index covers (left-prefix). MySQL keeps a SINGLE backing index for several FKs on the same
//     columns; dropping it while any of those FKs survives fails errno 1553. fkBackingIndexStillNeeded
//     enforces this; AND
//   - the `to` table keeps no USER index of that name (the index node owns a user-named index and
//     leaves it intentionally — diff_index.go userIndexNameSet).
//
// A whole-table DROP (DiffDrop) never reaches this helper: that path goes through
// dropForeignKeyConstraintOpsForDroppedTable, which releases the FK constraints but emits NO
// backing-index drop (DROP TABLE removes the indexes), so it never consults leftoverBackingIndexToDrop.
// This helper therefore always runs with a non-nil entry.To (the DiffModify path).
func leftoverBackingIndexToDrop(entry *TableDiffEntry, con *Constraint) (string, bool) {
	from := entry.From
	if from == nil {
		return "", false
	}
	implicit := fkImplicitIndexNames(from)
	// Find the backing-index name for THIS constraint: the FK-implicit index whose name matches
	// the auto-derived form for `con` (constraint name / first-column / collision suffix).
	var backing *Index
	for _, idx := range from.Indexes {
		if idx == nil {
			continue
		}
		if !implicit[toLower(idx.Name)] {
			continue
		}
		if isAutoBackingIndexName(idx.Name, con) && isPlainBackingIndexFor(idx, con.Columns) {
			backing = idx
			break
		}
	}
	if backing == nil {
		return "", false
	}
	// Keep it if a surviving FK in `to` still needs it (shared backing index, errno 1553).
	if fkBackingIndexStillNeeded(entry.To, backing) {
		return "", false
	}
	// Keep it if the target still carries a USER index of that name (index node's domain).
	if userIndexNameSet(entry.To)[toLower(backing.Name)] {
		return "", false
	}
	return backing.Name, true
}

// fkBackingIndexStillNeeded reports whether any FOREIGN KEY surviving in `to` would still need
// `idx` as a backing index — i.e. a `to` FK whose columns `idx` covers as a left-prefix. MySQL
// maintains ONE backing index for all FKs on the same leading columns, so dropping `idx` while
// such an FK survives fails errno 1553 ("needed in a foreign key constraint"). This mirrors the
// engine's own rule for whether an index is still required. A nil `to` (no surviving table) needs
// nothing.
func fkBackingIndexStillNeeded(to *Table, idx *Index) bool {
	if to == nil || idx == nil {
		return false
	}
	for _, con := range to.Constraints {
		if con == nil || con.Type != ConForeignKey || len(con.Columns) == 0 {
			continue
		}
		if indexCoversColumns(idx, con.Columns) {
			return true
		}
	}
	return false
}

// indexCoversColumns reports whether idx's leading columns match fkCols in order (a left-prefix
// cover), the condition under which MySQL treats idx as a usable backing index for an FK on
// fkCols. Prefix lengths / expressions / direction are ignored here — they do not change whether
// the index can back the FK for the purpose of the errno-1553 "still needed" check (this matches
// hasIndexCoveringColumns in tablecmds.go, the loader's own cover test).
func indexCoversColumns(idx *Index, fkCols []string) bool {
	if len(idx.Columns) < len(fkCols) {
		return false
	}
	for i, c := range fkCols {
		if !strings.EqualFold(idx.Columns[i].Name, c) {
			return false
		}
	}
	return true
}

// orderedForeignKeys returns a table's FOREIGN KEY constraints in a deterministic order (by
// lower-cased constraint name) so a new table's ADD FOREIGN KEY sequence is stable across runs.
func orderedForeignKeys(t *Table) []*Constraint {
	if t == nil {
		return nil
	}
	var fks []*Constraint
	for _, con := range t.Constraints {
		if con != nil && con.Type == ConForeignKey {
			fks = append(fks, con)
		}
	}
	sort.Slice(fks, func(i, j int) bool {
		return toLower(fks[i].Name) < toLower(fks[j].Name)
	})
	return fks
}

// formatForeignKeyDefinition renders the constraint clause that follows ADD in an ALTER TABLE, in
// MySQL's canonical stored form:
//
//	CONSTRAINT `name` FOREIGN KEY (`c1`, `c2`) REFERENCES `db`.`tbl` (`r1`, `r2`)
//	    [ON DELETE <action>] [ON UPDATE <action>]
//
// It mirrors show.go's showConstraint (the stored FK line) so the readback of the ADD
// canonicalizes equal to `to`. The ON DELETE / ON UPDATE clauses are emitted ONLY for a
// non-default action (CASCADE / SET NULL / SET DEFAULT): RESTRICT and NO ACTION (and an absent
// clause) are the engine's default and are omitted, so the emitted DDL round-trips to the same
// stored form the canonical key folds to (see diff_foreign_key.go canonicalFKAction). Emitting an
// explicit RESTRICT/NO ACTION would also round-trip on equality, but omitting matches show.go and
// keeps the surface minimal.
func formatForeignKeyDefinition(con *Constraint) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CONSTRAINT %s FOREIGN KEY (", migrationQuoteIdent(con.Name))
	b.WriteString(quoteIdentList(con.Columns))
	b.WriteString(") REFERENCES ")
	if con.RefDatabase != "" {
		b.WriteString(migrationQuoteIdent(con.RefDatabase) + "." + migrationQuoteIdent(con.RefTable))
	} else {
		b.WriteString(migrationQuoteIdent(con.RefTable))
	}
	b.WriteString(" (")
	b.WriteString(quoteIdentList(con.RefColumns))
	b.WriteString(")")

	if con.OnDelete != "" && !isFKDefaultAction(con.OnDelete) {
		fmt.Fprintf(&b, " ON DELETE %s", strings.ToUpper(con.OnDelete))
	}
	if con.OnUpdate != "" && !isFKDefaultAction(con.OnUpdate) {
		fmt.Fprintf(&b, " ON UPDATE %s", strings.ToUpper(con.OnUpdate))
	}
	return b.String()
}

// isFKDefaultAction reports whether a referential action is a MySQL default that is omitted from
// the rendered FK clause. RESTRICT, NO ACTION, and an absent clause are all the "default" (no
// referential action) and are omitted — matching the canonical-equality collapse in
// canonicalFKAction, so the emitted DDL reads back to a form that diffs empty against `to` on both
// versions (8.0 omits NO ACTION but would echo RESTRICT; 5.7 omits RESTRICT but would echo NO
// ACTION — omitting BOTH from what we emit sidesteps the per-version echo entirely).
func isFKDefaultAction(action string) bool {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "", "RESTRICT", "NO ACTION":
		return true
	default:
		return false
	}
}

// quoteIdentList backtick-quotes each identifier and joins with ", " (the spacing show.go's FK
// line uses), preserving order.
func quoteIdentList(names []string) string {
	if len(names) == 0 {
		return ""
	}
	quoted := make([]string, len(names))
	for i, nm := range names {
		quoted[i] = migrationQuoteIdent(nm)
	}
	return strings.Join(quoted, ", ")
}

// foreignKeySortName is the stable secondary sort key for a FK op: lower-cased
// database.table.constraint.
func foreignKeySortName(database, table, constraint string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(constraint)
}
