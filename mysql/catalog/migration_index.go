package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MySQL SDL generate — indexes (the full index/key surface).
//
// This is the MySQL analog of PG's generateIndexDDL (pg/catalog/migration_index.go). It turns
// the index diff (TableDiffEntry.Indexes, populated by diff_index.go) into ALTER TABLE ...
// ADD/DROP index DDL. MySQL folds index changes into ALTER TABLE rather than standalone
// CREATE/DROP INDEX statements, and an index cannot be altered in place, so a modified index is
// a DROP followed by an ADD.
//
// Two table cases produce index ADDs:
//   - DiffModify: the per-index diff drives ADD / DROP / (DROP+ADD for MODIFY).
//   - DiffAdd: a freshly created table. generate-core's formatCreateTable renders only columns
//     and the inline PRIMARY KEY (migration_table.go) — it does NOT inline secondary/UNIQUE/
//     FULLTEXT/SPATIAL indexes — so every non-PK, non-FK-implicit index on a new table is added
//     here via ALTER TABLE ... ADD. The PK is already inline (skip it); FK-implicit-backing
//     indexes are owned by the FK node (skip them, same exclusion as the diff).
//
// DiffDrop tables emit nothing: DROP TABLE removes their indexes wholesale.
//
// Rendering routes through the same canonical forms show.go uses for the stored index line, so
// the readback of the emitted DDL canonicalizes equal to `to` and apply-correctness holds. PK
// changes use ADD/DROP PRIMARY KEY; every other index uses ADD ... KEY / DROP INDEX.
//
// Ordering: index ADDs run in PhaseMain at priorityIndex (after the table CREATE and column
// ALTERs, before deferred FKs in PhasePost — so a FK added in the same plan reuses an index this
// node created rather than auto-creating a duplicate). Index DROPs run in PhasePre (before
// PhaseMain), so a drop-then-add of the same index name never collides. sortMigrationOps
// re-imposes the global phase/priority/name order.
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
			// Drops ahead of adds locally (also enforced by phase) for readability.
			return indexOpRank(ops[i].Type) < indexOpRank(ops[j].Type)
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addIndexOpsForNewTable emits ADD ops for every user-managed non-PK index on a freshly created
// table (the PK is rendered inline by formatCreateTable; FK-implicit-backing indexes are the FK
// node's). It uses the diff-able index set so the GIPK and FK-implicit indexes are excluded
// exactly as in diffIndexes.
func addIndexOpsForNewTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	var ops []MigrationOp
	for _, idx := range orderedDiffableIndexes(entry.To) {
		if idx.Primary {
			continue // inline in CREATE TABLE
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
			ops = append(ops, dropIndexOp(entry, table, ie.From))
		case DiffModify:
			if ie.From == nil || ie.To == nil {
				continue
			}
			// The PRIMARY KEY interacts with column changes, so it needs special ordering — see
			// primaryKeyModifyOps. A secondary index is independent: DROP then ADD (the drop in
			// PhasePre always precedes the PhaseMain add, freeing the name for re-creation).
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

// addIndexOp builds an ALTER TABLE ... ADD <index> op. A PRIMARY KEY uses ADD PRIMARY KEY; every
// other kind uses ADD <UNIQUE|FULLTEXT|SPATIAL> KEY `name` (...). The op runs in PhaseMain at
// priorityIndex.
func addIndexOp(entry *TableDiffEntry, table string, idx *Index, n *Normalizer) MigrationOp {
	return MigrationOp{
		Type:         OpAddIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ADD %s", table, formatIndexDefinition(idx, n)),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
	}
}

// priorityPrimaryKeyDrop orders a standalone DROP PRIMARY KEY within PhasePre BEFORE column drops
// AND column NOT NULL→NULL demotions. When a PK change also drops a column that was a member of
// the OLD PK, or demotes a removed PK member to nullable, the PK must be gone first: a PhasePre
// DROP COLUMN of a PK member auto-removes the PK (then a later DROP PRIMARY KEY fails errno 1091),
// and MySQL rejects making a still-PK column nullable (errno 1171). Dropping the PK first frees
// both. It sits just below the index-drop priority so the PK is the very first thing released.
const priorityPrimaryKeyDrop = priorityIndexDrop - 1

// primaryKeyModifyOps renders a PRIMARY KEY change, choosing between two correct shapes because a
// PK change interleaves with the column generator's ops (which this node does not control):
//
//   - COMBINED, one statement `ALTER TABLE ... DROP PRIMARY KEY, ADD PRIMARY KEY (...)` in
//     PhaseMain. This is required when an AUTO_INCREMENT column is a member of the new PK: a
//     standalone DROP PRIMARY KEY would leave that column momentarily unkeyed, which MySQL rejects
//     (errno 1075). The combined statement has no unkeyed window. It is only valid, though, when no
//     OLD-PK-member column is being dropped in the same diff — a PhasePre DROP COLUMN of a PK
//     member auto-removes the PK before this PhaseMain statement runs, so its DROP PRIMARY KEY half
//     would fail (errno 1091).
//   - SPLIT: a standalone DROP PRIMARY KEY in PhasePre (at priorityPrimaryKeyDrop, before column
//     drops and NOT NULL demotions) plus an ADD PRIMARY KEY in PhaseMain (at priorityIndex, after
//     the column generator's adds/modifies at priorityColumn, so the new PK's columns are already
//     present and NOT NULL). This is correct whenever a PK member is dropped or demoted, and for
//     all non-AUTO_INCREMENT PK changes.
//
// The two cases are disjoint in practice: an AUTO_INCREMENT column must stay keyed, so it is never
// the dropped PK member; the combined form is chosen only when AUTO_INCREMENT keys the new PK and
// no PK-member column drop forces the split. `to` always carries the PRIMARY index for a PK
// change (the diff keys it by the PRIMARY name), so DROP PRIMARY KEY is unconditional.
func primaryKeyModifyOps(entry *TableDiffEntry, table string, to *Index, n *Normalizer) []MigrationOp {
	if canCombinePrimaryKeyModify(entry, to) {
		return []MigrationOp{{
			Type:         OpAddIndex,
			Database:     entry.Database,
			ObjectName:   entry.Name,
			ParentObject: entry.Name,
			SQL:          fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY, ADD %s", table, formatIndexDefinition(to, n)),
			Phase:        PhaseMain,
			Priority:     priorityIndex,
			sortName:     indexSortName(entry.Database, entry.Name, to.Name),
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
	}
	addPK := addIndexOp(entry, table, to, n)
	return []MigrationOp{dropPK, addPK}
}

// canCombinePrimaryKeyModify reports whether a PK change may be emitted as a single combined
// DROP+ADD statement: the new PK must include an AUTO_INCREMENT column (the reason the combined
// form is needed — to avoid an unkeyed AUTO_INCREMENT window, errno 1075) AND no column that was a
// member of the OLD PK may be dropped in this diff (such a drop runs in PhasePre and auto-removes
// the PK before the combined PhaseMain statement, breaking its DROP PRIMARY KEY half). Otherwise
// the split form (DROP PK early, ADD PK late) is used.
func canCombinePrimaryKeyModify(entry *TableDiffEntry, to *Index) bool {
	if !newPrimaryKeyHasAutoIncrement(entry.To, to) {
		return false
	}
	return !oldPrimaryKeyMemberDropped(entry)
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

// oldPrimaryKeyMemberDropped reports whether any column being dropped in this diff was a member of
// the old PRIMARY KEY — the case that forces the split form (the column drop auto-removes the PK).
func oldPrimaryKeyMemberDropped(entry *TableDiffEntry) bool {
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
		if ce.Action == DiffDrop && oldPKCols[toLower(ce.Name)] {
			return true
		}
	}
	return false
}

// priorityIndexDrop orders an index DROP within PhasePre BEFORE any column DROP. MySQL's
// ALTER TABLE ... DROP COLUMN auto-removes an index that the dropped column leaves empty (a
// single-column index, or the last remaining column of a composite one), after which an
// explicit DROP INDEX for that index fails with errno 1091 ("Can't DROP 'k'; check that
// column/key exists"). Because the migration plan is applied one statement at a time, the index
// DROP must precede the column DROP. Column drops run at priorityColumn-1 (generated) /
// priorityColumn (plain) within PhasePre, and table drops at priorityTable=10; this sits below
// the column-drop range and above the table-drop, so the order is:
// table-drop (10) < index-drop (15) < generated-col-drop (19) < plain-col-drop (20).
// Dropping an index whose columns are NOT being dropped is harmless to sequence first (it is
// independent), so always ordering index drops ahead of column drops is safe and simple.
const priorityIndexDrop = priorityColumn - 5

// dropIndexOp builds an ALTER TABLE ... DROP op. A PRIMARY KEY uses DROP PRIMARY KEY (it has no
// user name); every other index uses DROP INDEX `name`. The op runs in PhasePre at
// priorityIndexDrop so all index drops precede any add AND precede column drops (see
// priorityIndexDrop), and so a dropped index name is free for re-creation in the same plan.
func dropIndexOp(entry *TableDiffEntry, table string, idx *Index) MigrationOp {
	var clause string
	if idx.Primary {
		clause = "DROP PRIMARY KEY"
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
	}
}

// formatIndexDefinition renders the index clause that follows ADD in an ALTER TABLE, in MySQL's
// canonical stored form: PRIMARY KEY (...) | UNIQUE KEY `name` (...) | FULLTEXT KEY `name` (...)
// | SPATIAL KEY `name` (...) | KEY `name` (...), with the column parts, an optional USING type,
// KEY_BLOCK_SIZE, COMMENT, and INVISIBLE marker.
//
// It mirrors show.go's showIndex (the stored index line) so the readback of the ADD
// canonicalizes equal to `to`. The one addition over showIndex is KEY_BLOCK_SIZE: show.go never
// emits it (SHOW CREATE drops the index-level value on 8.0), but the ADD must carry it so the
// value is actually applied — on 5.7 the readback then echoes it (matches), and on 8.0 the
// readback drops it but canonicalIndexKeyBlockSize also ignores it on 8.0, so the round-trip is
// empty either way.
func formatIndexDefinition(idx *Index, _ *Normalizer) string {
	var b strings.Builder

	switch {
	case idx.Primary:
		b.WriteString("PRIMARY KEY (")
	case idx.Unique:
		fmt.Fprintf(&b, "UNIQUE KEY %s (", migrationQuoteIdent(idx.Name))
	case idx.Fulltext:
		fmt.Fprintf(&b, "FULLTEXT KEY %s (", migrationQuoteIdent(idx.Name))
	case idx.Spatial:
		fmt.Fprintf(&b, "SPATIAL KEY %s (", migrationQuoteIdent(idx.Name))
	default:
		fmt.Fprintf(&b, "KEY %s (", migrationQuoteIdent(idx.Name))
	}

	cols := make([]string, 0, len(idx.Columns))
	for _, ic := range idx.Columns {
		cols = append(cols, migrationIndexColumn(ic))
	}
	b.WriteString(strings.Join(cols, ","))
	b.WriteString(")")

	// USING clause: only the user-selectable BTREE/HASH, never for PRIMARY/FULLTEXT/SPATIAL
	// (matching show.go), so a redundant FULLTEXT/SPATIAL IndexType echo is not emitted.
	if !idx.Primary && !idx.Fulltext && !idx.Spatial && idx.IndexType != "" {
		fmt.Fprintf(&b, " USING %s", strings.ToUpper(idx.IndexType))
	}

	// KEY_BLOCK_SIZE: emit when set so the value is applied. (show.go omits it from the stored
	// line; we emit it on the ADD — see the function doc for why this round-trips on both
	// versions.)
	if idx.KeyBlockSize > 0 {
		fmt.Fprintf(&b, " KEY_BLOCK_SIZE=%d", idx.KeyBlockSize)
	}

	// COMMENT.
	if idx.Comment != "" {
		fmt.Fprintf(&b, " COMMENT %s", quoteStringLiteral(idx.Comment))
	}

	// INVISIBLE (8.0). The /*!80000 ... */ version guard makes it a no-op on 5.7, matching how
	// show.go renders an invisible index, so the same DDL is valid on both engines.
	if !idx.Visible {
		b.WriteString(" /*!80000 INVISIBLE */")
	}

	return b.String()
}

// orderedDiffableIndexes returns a table's diff-able indexes (GIPK + FK-implicit excluded, the
// same set diffIndexes uses) in a deterministic order: by lower-cased name. Determinism keeps a
// new table's ADD INDEX sequence stable across runs. This is a single-table (new-table) context
// with no other side, so FK-implicit indexes are excluded unconditionally (nil other-side set) —
// the FK node adds the backing index for a newly created table's foreign keys.
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

// indexOpRank orders index op types so drops sort ahead of adds locally (phase ordering also
// enforces this globally).
func indexOpRank(t MigrationOpType) int {
	if t == OpDropIndex {
		return 0
	}
	return 1
}

// indexSortName is the stable secondary sort key for an index op: lower-cased
// database.table.index.
func indexSortName(database, table, index string) string {
	return strings.ToLower(database) + "." + strings.ToLower(table) + "." + strings.ToLower(index)
}
