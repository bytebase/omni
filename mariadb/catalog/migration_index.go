package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// generateIndexDDL turns the index diff (TableDiffEntry.Indexes) into ALTER TABLE ADD and DROP
// index ops. An index cannot be altered in place, so a modified index is a DROP followed by an
// ADD. A PRIMARY KEY uses ADD/DROP PRIMARY KEY; every other index uses ADD ... KEY / DROP INDEX.
//
// A new table gets every non-PK index added after its CREATE, because formatCreateTable inlines
// only the PK and the AUTO_INCREMENT supporting key. Indexes that back a foreign key are left to
// the FK generator.
//
// Index ADDs run at priorityIndex, after the table CREATE and column ALTERs and before deferred
// FKs, so a later FK reuses the index instead of creating a duplicate. Index DROPs run in
// PhasePre so a drop-then-add of the same name never collides.
func generateIndexDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, addIndexOpsForNewTable(entry, n)...)
		case DiffModify:
			ops = append(ops, indexOpsForModifiedTable(entry, n)...)
		case DiffDrop:
			// DROP TABLE removes the indexes; nothing to emit.
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		if ops[i].Type != ops[j].Type {
			return indexOpRank(ops[i].Type) < indexOpRank(ops[j].Type)
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addIndexOpsForNewTable emits ADD ops for every user-managed non-PK index on a new table,
// using the same diff-able index set as diffIndexes.
func addIndexOpsForNewTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	// formatCreateTable renders these inline; adding again would be a duplicate key name.
	inlined := autoIncSupportingIndex(entry.To)
	inlinedPK := inlineImplicitPrimaryKey(entry.To)
	var ops []MigrationOp
	for _, idx := range orderedDiffableIndexes(entry.To) {
		if idx.Primary {
			continue // inline in CREATE TABLE
		}
		if inlined != nil && idx == inlined {
			continue // inline in CREATE TABLE (AUTO_INCREMENT supporting key)
		}
		if inlinedPK != nil && idx == inlinedPK {
			continue // inline in CREATE TABLE (implicit primary key)
		}
		ops = append(ops, addIndexOp(entry, table, idx, n))
	}
	return ops
}

// indexOpsForModifiedTable emits ADD / DROP / (DROP+ADD) ops from the per-index diff of a
// modified table.
func indexOpsForModifiedTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil || len(entry.Indexes) == 0 {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	// Release the retained IGNORED keys a replacement would otherwise promote to
	// the implicit primary key (4174): drop them before the replaced keys' drops
	// and re-add them, in target order, after the visible keys are back.
	for _, released := range ignoredKeysToRelease(entry) {
		drop := dropIndexOp(entry, table, released)
		drop.Priority = priorityIndexDrop - 2
		ops = append(ops, drop, addIndexOp(entry, table, released, n))
	}
	for _, ie := range entry.Indexes {
		switch ie.Action {
		case DiffAdd:
			if ie.To == nil {
				continue
			}
			ops = append(ops, addIndexOp(entry, table, ie.To, n))
		case DiffDrop:
			if ie.From == nil {
				continue
			}
			drop := dropIndexOp(entry, table, ie.From)
			// An IGNORED unique key drops before visible-key drops, so dropping a
			// visible key never leaves it as the sole implicit-PK candidate (4174).
			if !ie.From.Primary && ie.From.Unique && !ie.From.Visible {
				drop.Priority = priorityIndexDrop - 2
			}
			ops = append(ops, drop)
		case DiffModify:
			if ie.From == nil || ie.To == nil {
				continue
			}
			// A PRIMARY KEY change interacts with column changes; see primaryKeyModifyOps.
			if ie.To.Primary || ie.From.Primary {
				ops = append(ops, primaryKeyModifyOps(entry, table, ie.To, n)...)
				continue
			}
			ops = append(ops, dropIndexOp(entry, table, ie.From))
			ops = append(ops, addIndexOp(entry, table, ie.To, n))
		}
	}
	return ops
}

// ignoredKeysToRelease returns the retained IGNORED unique keys that other changes
// would, at some intermediate point, leave as the first eligible implicit primary
// key (4174) — whether a replaced or dropped key exposes them, or a nullable-to-NOT
// NULL column change makes them eligible before a new primary key is added. MariaDB
// enforces that invariant within a statement, so an atomic DROP+ADD does not help;
// these keys are instead dropped before the other changes and re-added afterwards,
// in target order. It covers any number of candidates and simultaneous drops. Keys
// the diff already changes, and keys the plan only adds, are left to their own
// operations.
func ignoredKeysToRelease(entry *TableDiffEntry) []*Index {
	if entry.From == nil || entry.To == nil {
		return nil
	}
	removed := make(map[string]bool)
	changed := make(map[string]bool)
	for _, ie := range entry.Indexes {
		if ie.From != nil {
			changed[toLower(ie.From.Name)] = true
			// An outright-dropped key leaves a window without it; a replaced
			// primary key does not, since its DROP+ADD is one statement.
			if ie.Action == DiffDrop || ie.Action == DiffModify && !ie.From.Primary {
				removed[toLower(ie.From.Name)] = true
			}
		}
		if ie.To != nil {
			changed[toLower(ie.To.Name)] = true
		}
	}
	// The intermediate state: a column is NOT NULL if it is in either schema, so a
	// nullability change in either direction that makes an IGNORED key eligible is
	// seen, with the source keys that survive. A primary key the plan adds or drops
	// — which would otherwise protect the ignored one — is absent.
	intermediate := *entry.To
	intermediate.Columns = worstCaseNullability(entry.From, entry.To)
	intermediate.Indexes = nil
	for _, idx := range entry.From.Indexes {
		if !removed[toLower(idx.Name)] {
			intermediate.Indexes = append(intermediate.Indexes, idx)
		}
	}
	if validateIgnoredPrimaryKey(&intermediate, false) == nil {
		return nil
	}
	var release []*Index
	for _, idx := range intermediate.Indexes {
		if idx.Primary || idx.Visible || !idx.Unique || changed[toLower(idx.Name)] {
			continue
		}
		if target := tableIndexByName(entry.To, idx.Name); target != nil {
			release = append(release, target)
		}
	}
	return release
}

// worstCaseNullability returns to's columns with each one marked NOT NULL if it is
// NOT NULL in either from or to, the state an IGNORED key is most likely eligible in
// at some point during the migration.
func worstCaseNullability(from, to *Table) []*Column {
	notNull := make(map[string]bool)
	for _, c := range from.Columns {
		if !c.Nullable {
			notNull[toLower(c.Name)] = true
		}
	}
	cols := make([]*Column, len(to.Columns))
	for i, c := range to.Columns {
		if c.Nullable && notNull[toLower(c.Name)] {
			clone := *c
			clone.Nullable = false
			cols[i] = &clone
			continue
		}
		cols[i] = c
	}
	return cols
}

// addIndexOp builds an ALTER TABLE ... ADD <index> op.
func addIndexOp(entry *TableDiffEntry, table string, idx *Index, n *Normalizer) MigrationOp {
	clause := "ADD " + formatIndexDefinition(idx, n)
	return MigrationOp{
		Type:         OpAddIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexAddSortName(entry, idx),
		alterClause:  clause,
		addsIndex:    idx,
	}
}

// priorityPrimaryKeyDrop orders a standalone DROP PRIMARY KEY within PhasePre before column drops
// and NOT NULL to NULL demotions. Dropping a PK member column removes the PK, so a later DROP
// PRIMARY KEY would fail, and a column still in the PK cannot be made nullable.
const priorityPrimaryKeyDrop = priorityIndexDrop - 1

// primaryKeyModifyOps renders a PRIMARY KEY change in one of two shapes:
//
//   - Combined: one PhaseMain statement `DROP PRIMARY KEY, ADD PRIMARY KEY (...)`. Required when
//     an AUTO_INCREMENT column is in the new PK, because a standalone DROP would leave it
//     unkeyed. Valid only when no old PK member is dropped or demoted in the same diff.
//   - Split: DROP PRIMARY KEY in PhasePre at priorityPrimaryKeyDrop, ADD PRIMARY KEY in PhaseMain
//     at priorityIndex, after the column adds and modifies so the new PK's columns exist and are
//     NOT NULL.
//
// mergeAutoIncrementKeyOps (migration_autoinc.go) later folds the split path's operations into
// one statement where an AUTO_INCREMENT column needs it.
func primaryKeyModifyOps(entry *TableDiffEntry, table string, to *Index, n *Normalizer) []MigrationOp {
	if canCombinePrimaryKeyModify(entry, to) {
		clause := "DROP PRIMARY KEY, ADD " + formatIndexDefinition(to, n)
		return []MigrationOp{{
			Type:         OpAddIndex,
			Database:     entry.Database,
			ObjectName:   entry.Name,
			ParentObject: entry.Name,
			SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
			Phase:        PhaseMain,
			Priority:     priorityIndex,
			sortName:     indexSortName(entry.Database, entry.Name, to.Name),
			alterClause:  clause,
			addsIndex:    to,
			dropsIndex:   primaryIndexKey,
		}}
	}
	dropPK := MigrationOp{
		Type:         OpDropIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY", table),
		Phase:        PhasePre,
		Priority:     priorityPrimaryKeyDrop,
		sortName:     indexSortName(entry.Database, entry.Name, to.Name),
		alterClause:  "DROP PRIMARY KEY",
		dropsIndex:   primaryIndexKey,
	}
	addPK := addIndexOp(entry, table, to, n)
	return []MigrationOp{dropPK, addPK}
}

// canCombinePrimaryKeyModify reports whether a PK change may be a single combined DROP+ADD
// statement: the new PK includes an AUTO_INCREMENT column, and no old PK member is dropped or
// demoted to nullable in this diff. Those column changes run before the combined PhaseMain
// statement while the old PK still holds the column, so its DROP PRIMARY KEY half would fail.
func canCombinePrimaryKeyModify(entry *TableDiffEntry, to *Index) bool {
	if oldPrimaryKeyMemberDroppedOrDemoted(entry) {
		return false
	}
	// A standalone DROP PRIMARY KEY would promote a retained IGNORED unique
	// key to the implicit primary key, which an IGNORED key cannot be (4174),
	// so the replacement is one statement.
	return newPrimaryKeyHasAutoIncrement(entry.To, to) || splitDropPromotesIgnored(entry)
}

// splitDropPromotesIgnored reports whether dropping the table's primary key on
// its own, before adding the replacement, would make a retained IGNORED key
// the implicit primary key (4174).
func splitDropPromotesIgnored(entry *TableDiffEntry) bool {
	if entry.To == nil {
		return false
	}
	withoutPK := *entry.To
	withoutPK.Indexes = nil
	for _, idx := range entry.To.Indexes {
		if !idx.Primary {
			withoutPK.Indexes = append(withoutPK.Indexes, idx)
		}
	}
	return validateIgnoredPrimaryKey(&withoutPK, false) != nil
}

// newPrimaryKeyHasAutoIncrement reports whether any column of the new PK is AUTO_INCREMENT.
func newPrimaryKeyHasAutoIncrement(tbl *Table, pk *Index) bool {
	if tbl == nil || pk == nil {
		return false
	}
	for _, ic := range pk.Columns {
		if col := tbl.GetColumn(ic.Name); col != nil && col.AutoIncrement {
			return true
		}
	}
	return false
}

// oldPrimaryKeyMemberDroppedOrDemoted reports whether any column of the old PRIMARY KEY is
// dropped or made nullable in this diff, which forces the split form. A column that leaves the
// PK but stays NOT NULL is fine for the combined form.
func oldPrimaryKeyMemberDroppedOrDemoted(entry *TableDiffEntry) bool {
	if entry.From == nil {
		return false
	}
	oldPKCols := make(map[string]bool)
	for _, idx := range entry.From.Indexes {
		if idx != nil && idx.Primary {
			for _, ic := range idx.Columns {
				oldPKCols[toLower(ic.Name)] = true
			}
			break
		}
	}
	if len(oldPKCols) == 0 {
		return false
	}
	for _, ce := range entry.Columns {
		if !oldPKCols[toLower(ce.Name)] {
			continue
		}
		if ce.Action == DiffDrop {
			return true
		}
		if ce.Action == DiffModify && ce.To != nil && ce.To.Nullable {
			return true
		}
	}
	return false
}

// priorityIndexDrop orders an index DROP within PhasePre after table drops and before any column
// drop. Dropping a column can remove an index it leaves empty, after which an explicit DROP
// INDEX would fail; plans apply one statement at a time, so the index drop goes first.
const priorityIndexDrop = priorityColumn - 5

// dropIndexOp builds an ALTER TABLE ... DROP op: DROP PRIMARY KEY for the PRIMARY KEY, DROP
// INDEX `name` otherwise.
func dropIndexOp(entry *TableDiffEntry, table string, idx *Index) MigrationOp {
	var clause string
	drops := toLower(idx.Name)
	if idx.Primary {
		clause = "DROP PRIMARY KEY"
		drops = primaryIndexKey
	} else {
		clause = "DROP INDEX " + migrationQuoteIdent(idx.Name)
	}
	return MigrationOp{
		Type:         OpDropIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhasePre,
		Priority:     priorityIndexDrop,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
		alterClause:  clause,
		dropsIndex:   drops,
	}
}

// formatIndexDefinition is a key as SHOW CREATE TABLE writes it.
func formatIndexDefinition(idx *Index, n *Normalizer) string {
	return n.CanonicalIndex(idx)
}

// orderedDiffableIndexes returns a new table's diff-able indexes ordered by lower-cased name.
// Foreign-key backing indexes are excluded; the FK generator adds them.
func orderedDiffableIndexes(t *Table) []*Index {
	m := diffableIndexMap(t, nil)
	out := make([]*Index, 0, len(m))
	for _, idx := range m {
		out = append(out, idx)
	}
	sort.Slice(out, func(i, j int) bool {
		return toLower(out[i].Name) < toLower(out[j].Name)
	})
	return out
}

// indexOpRank orders index op types so drops sort ahead of adds.
func indexOpRank(t MigrationOpType) int {
	if t == OpDropIndex {
		return 0
	}
	return 1
}

// indexAddSortName sorts the keys a table gains in the target's order: the first
// UNIQUE key over NOT NULL columns is the table's implicit primary key, which
// an IGNORED key cannot be (4174), so the order which one that is depends on
// has to hold.
func indexAddSortName(entry *TableDiffEntry, idx *Index) string {
	position := len(entry.To.Indexes)
	for i, other := range entry.To.Indexes {
		if other == idx || toLower(other.Name) == toLower(idx.Name) {
			position = i
			break
		}
	}
	return indexSortName(entry.Database, entry.Name, fmt.Sprintf("\x00%04d", position))
}

// indexSortName is the stable secondary sort key for an index op: lower-cased
// database.table.index.
func indexSortName(database, table, index string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(index)
}
