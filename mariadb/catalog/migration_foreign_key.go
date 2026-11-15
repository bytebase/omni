package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateForeignKeyDDL turns the FK diff (TableDiffEntry.ForeignKeys) into ALTER TABLE ADD and
// DROP FOREIGN KEY ops. A FK cannot be altered in place, so a modified FK is a DROP followed by
// an ADD.
//
// FK ADDs run in PhasePost, after index ADDs and table/column DDL, so every referenced table
// exists and a user-declared covering index is reused instead of a duplicate backing index being
// auto-created. FK DROPs run in PhasePre at priorityForeignKeyConstraintDrop, ahead of every
// dependent drop, because a live FK blocks dropping a referenced table, a referencing column, or
// its backing index.
//
// Dropping a FK leaves its auto-created backing index behind. The index generator excludes
// FK-implicit indexes from its diff, so this generator drops that leftover (see
// leftoverBackingIndexToDrop). seenBackingDrop deduplicates the DROP INDEX when several FKs on
// the same columns share one backing index, keyed by encodeKeyFields so identifiers containing
// `.` cannot collide.
func generateForeignKeyDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp
	seenBackingDrop := map[string]bool{}

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addForeignKeyOpsForNewTable(entry, n)...)
		case DiffModify:
			ops = append(ops, foreignKeyOpsForModifiedTable(entry, seenBackingDrop, n)...)
		case DiffDrop:
			// A FK on a dropped table can block the drop of a table it references, and table
			// drops share one priority, so release the FKs ahead of all table drops.
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

// dropForeignKeyConstraintOpsForDroppedTable emits a DROP FOREIGN KEY for every FK on a dropped
// table, ahead of all table drops, so a referenced table dropped in the same plan is not
// blocked. It emits no backing-index drop because DROP TABLE removes the indexes.
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

// priorityForeignKeyConstraintDrop orders FK constraint drops at the front of PhasePre, ahead of
// table, index, and column drops. It is negative so it sorts before priorityTable.
const priorityForeignKeyConstraintDrop = -10

// priorityForeignKeyBackingIndexDrop orders the leftover-backing-index drop after all FK
// constraint drops, so an index shared by several FKs survives until the last is gone, and before
// column drops.
const priorityForeignKeyBackingIndexDrop = priorityIndexDrop - 1

// addForeignKeyOpsForNewTable emits ADD CONSTRAINT ... FOREIGN KEY ops for every FK on a new
// table, in PhasePost, because formatCreateTable does not inline FK constraints.
func addForeignKeyOpsForNewTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, con := range orderedForeignKeys(entry.To) {
		ops = append(ops, addForeignKeyOps(entry, table, con, n)...)
	}
	return ops
}

// foreignKeyOpsForModifiedTable emits ADD, DROP, or DROP+ADD ops from the per-FK diff of a
// modified table. The DROP runs in PhasePre and the ADD in PhasePost, so the old constraint is
// released before the new one is created.
func foreignKeyOpsForModifiedTable(entry *TableDiffEntry, seenBackingDrop map[string]bool, n *Normalizer) []MigrationOp {
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
			ops = append(ops, addForeignKeyOps(entry, tableIdent(entry.To), fe.To, n)...)
		case DiffDrop:
			if fe.From == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOps(entry, fe.From, seenBackingDrop)...)
		case DiffModify:
			if fe.From == nil || fe.To == nil || entry.To == nil {
				continue
			}
			ops = append(ops, dropForeignKeyOps(entry, fe.From, seenBackingDrop)...)
			ops = append(ops, addForeignKeyOps(entry, tableIdent(entry.To), fe.To, n)...)
		}
	}
	return ops
}

// addForeignKeyOps returns the FK's ADD, preceded by the ADD of its backing index when the
// server would otherwise name that index differently from `to` (fkBackingIndexOp).
func addForeignKeyOps(entry *TableDiffEntry, table string, con *Constraint, n *Normalizer) []MigrationOp {
	if op, ok := fkBackingIndexOp(entry, table, con, n); ok {
		return []MigrationOp{op, addForeignKeyOp(entry, table, con)}
	}
	return []MigrationOp{addForeignKeyOp(entry, table, con)}
}

// fkBackingIndexOp returns an ADD for the FK-implicit index backing con in `to` when the FK's own
// ADD would not reproduce it. ADD CONSTRAINT `c` FOREIGN KEY names the index it creates `c`, so a
// backing index named after the first FK column must exist before the FK runs, which then reuses
// it. No ADD is needed when the index generator owns the name (a user index on the `from` side),
// when a `from` index of that name survives, or when CREATE TABLE already inlines it
// (autoIncSupportingIndex).
//
// The op is typed OpAddForeignKey because the index belongs to the FK's lifecycle; its grouping
// metadata lets the AUTO_INCREMENT pass use it as a covering key.
func fkBackingIndexOp(entry *TableDiffEntry, table string, con *Constraint, n *Normalizer) (MigrationOp, bool) {
	idx := fkImplicitBackingIndexes(entry.To)[con]
	if idx == nil || strings.EqualFold(idx.Name, con.Name) {
		return MigrationOp{}, false
	}
	if entry.Action == DiffAdd {
		if idx == autoIncSupportingIndex(entry.To) {
			return MigrationOp{}, false
		}
	} else {
		if userIndexNameSet(entry.From)[toLower(idx.Name)] {
			return MigrationOp{}, false
		}
		if fromIdx := tableIndexByName(entry.From, idx.Name); fromIdx != nil && fkBackingIndexStillNeeded(entry, fromIdx) {
			return MigrationOp{}, false
		}
	}
	clause := "ADD " + formatIndexDefinition(idx, n)
	return MigrationOp{
		Type:         OpAddForeignKey,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
		alterClause:  clause,
		addsIndex:    idx,
	}, true
}

// tableIndexByName returns t's index of the given name (case-insensitive), or nil.
func tableIndexByName(t *Table, name string) *Index {
	if t == nil {
		return nil
	}
	for _, idx := range t.Indexes {
		if idx != nil && strings.EqualFold(idx.Name, name) {
			return idx
		}
	}
	return nil
}

// addForeignKeyOp builds an ALTER TABLE ... ADD CONSTRAINT `name` FOREIGN KEY (...) REFERENCES
// ... op in PhasePost, after table, column, and index DDL.
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

// dropForeignKeyOps builds the ops to drop a foreign key: DROP FOREIGN KEY `name` at
// priorityForeignKeyConstraintDrop, plus a DROP INDEX for its leftover backing index at
// priorityForeignKeyBackingIndexDrop when leftoverBackingIndexToDrop says so. The two priorities
// keep every constraint drop ahead of the index drop, even across FKs sharing one backing index.
// seenBackingDrop deduplicates that DROP INDEX across such FKs.
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
		// encodeKeyFields is collision-free where a dotted db.table.index string is not.
		dedupKey := encodeKeyFields("db", entry.Database, "table", entry.Name, "index", idxName)
		if !seenBackingDrop[dedupKey] {
			seenBackingDrop[dedupKey] = true
			ops = append(ops, MigrationOp{
				// Typed OpDropForeignKey, not OpDropIndex, because the index belongs to the FK's
				// lifecycle; alterClause and dropsIndex still let the AUTO_INCREMENT pass see the
				// coverage it removes.
				Type:         OpDropForeignKey,
				Database:     entry.Database,
				ObjectName:   entry.Name,
				ParentObject: entry.Name,
				SQL:          fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, migrationQuoteIdent(idxName)),
				Phase:        PhasePre,
				Priority:     priorityForeignKeyBackingIndexDrop,
				sortName:     foreignKeySortName(entry.Database, entry.Name, idxName),
				alterClause:  "DROP INDEX " + migrationQuoteIdent(idxName),
				dropsIndex:   toLower(idxName),
			})
		}
	}
	return ops
}

// leftoverBackingIndexToDrop decides whether dropping `con` (a FK on the `from` table) leaves an
// auto-created backing index this generator must drop, and returns its original-case name. It
// returns the index only when the `from` table has an FK-implicit backing index for `con`, no FK
// surviving in `to` still needs it, and `to` keeps no user index of that name. It runs only on
// the DiffModify path, so entry.To is non-nil.
func leftoverBackingIndexToDrop(entry *TableDiffEntry, con *Constraint) (string, bool) {
	from := entry.From
	if from == nil {
		return "", false
	}
	implicit := fkImplicitIndexNames(from)
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
	// A surviving FK still needs it.
	if fkBackingIndexStillNeeded(entry, backing) {
		return "", false
	}
	// The index generator owns a user index of that name.
	if userIndexNameSet(entry.To)[toLower(backing.Name)] {
		return "", false
	}
	return backing.Name, true
}

// fkBackingIndexStillNeeded reports whether any FOREIGN KEY retained in `entry.To` has columns
// `idx` covers as a left-prefix, so dropping `idx` would break it. A foreign key the plan only
// adds is excluded: it creates its own backing index, so it never keeps a differently named one
// alive. A nil target needs nothing.
func fkBackingIndexStillNeeded(entry *TableDiffEntry, idx *Index) bool {
	to := entry.To
	if to == nil || idx == nil {
		return false
	}
	retained := foreignKeyNameSet(entry.From)
	backing := fkImplicitBackingIndexes(to)
	for _, con := range to.Constraints {
		if con == nil || con.Type != ConForeignKey || len(con.Columns) == 0 {
			continue
		}
		if !retained[toLower(con.Name)] {
			continue
		}
		// A key of the name that the target makes over other columns, as a
		// foreign key that loses a column does, replaces this one.
		if own := backing[con]; own != nil && toLower(own.Name) == toLower(idx.Name) && !sameKeyColumns(own, idx) {
			continue
		}
		if indexCoversColumns(idx, con.Columns) {
			return true
		}
	}
	return false
}

// foreignKeyNameSet returns the lower-cased names of a table's foreign keys. A nil table has none.
func foreignKeyNameSet(t *Table) map[string]bool {
	set := make(map[string]bool)
	if t == nil {
		return set
	}
	for _, con := range t.Constraints {
		if con != nil && con.Type == ConForeignKey {
			set[toLower(con.Name)] = true
		}
	}
	return set
}

func sameKeyColumns(a, b *Index) bool {
	if len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if toLower(a.Columns[i].Name) != toLower(b.Columns[i].Name) {
			return false
		}
	}
	return true
}

// indexCoversColumns reports whether idx's leading columns match fkCols in order. Prefix
// lengths, expressions, and direction are ignored.
func indexCoversColumns(idx *Index, fkCols []string) bool {
	if len(idx.Columns) < len(fkCols) {
		return false
	}
	for i, c := range fkCols {
		// A prefix key part does not back a foreign key (1553), so it never
		// supplies surviving coverage.
		if !strings.EqualFold(idx.Columns[i].Name, c) || idx.Columns[i].Length > 0 {
			return false
		}
	}
	return true
}

// orderedForeignKeys returns a table's FOREIGN KEY constraints ordered by lower-cased name.
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

// formatForeignKeyDefinition renders the constraint clause that follows ADD in an ALTER TABLE:
//
//	CONSTRAINT `name` FOREIGN KEY (`c1`, `c2`) REFERENCES `db`.`tbl` (`r1`, `r2`)
//	    [ON DELETE <action>] [ON UPDATE <action>]
//
// An action is omitted when absent or the default; an explicit RESTRICT is emitted.
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

	if con.OnDelete != "" && !isFKDefault(con.OnDelete) {
		fmt.Fprintf(&b, " ON DELETE %s", strings.ToUpper(con.OnDelete))
	}
	if con.OnUpdate != "" && !isFKDefault(con.OnUpdate) {
		fmt.Fprintf(&b, " ON UPDATE %s", strings.ToUpper(con.OnUpdate))
	}
	return b.String()
}

// quoteIdentList backtick-quotes each identifier and joins with ", ", preserving order.
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
