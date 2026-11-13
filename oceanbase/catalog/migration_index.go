package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// SDL generate — indexes (the full index/key surface).
//
// It turns the index diff (TableDiffEntry.Indexes, populated by diff_index.go) into ALTER
// TABLE ... ADD/DROP index DDL. An index is not altered in place, so a modified index is a
// DROP followed by an ADD, but for its visibility (ALTER INDEX) and for the primary key,
// which is replaced in one statement (primaryKeyModifyOp).
//
// Two table cases produce index ADDs:
//   - DiffModify: the per-index diff drives ADD / DROP / (DROP+ADD for MODIFY).
//   - DiffAdd: a freshly created table. formatCreateTable renders only columns and the inline
//     PRIMARY KEY (migration_table.go), so every other index of a new table is added here via
//     ALTER TABLE ... ADD. That holds for the key of an AUTO_INCREMENT column too: OceanBase
//     lets the column stand without one, where MySQL wants it in the CREATE TABLE.
//
// DiffDrop tables emit nothing: DROP TABLE removes their indexes wholesale.
//
// Rendering routes through the same canonical forms show.go uses for the stored index line, so
// the readback of the emitted DDL canonicalizes equal to `to` and apply-correctness holds. PK
// changes use ADD/DROP PRIMARY KEY; every other index uses ADD ... KEY / DROP INDEX.
//
// Ordering: index ADDs run in PhaseMain at priorityIndex (after the table CREATE and column
// ALTERs, before deferred FKs in PhasePost). Index DROPs run in PhasePre (before PhaseMain),
// so a drop-then-add of the same index name never collides. sortMigrationOps re-imposes the
// global phase/priority/name order.
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
		// Drops ahead of adds, as the phases have them, then by name. Ops that tie keep the
		// order they were made in: an ADD and the ALTER INDEX that hides what it added.
		if ri, rj := indexOpRank(ops[i].Type), indexOpRank(ops[j].Type); ri != rj {
			return ri < rj
		}
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// addIndexOpsForNewTable emits ADD ops for every non-PK index of a freshly created table (the
// PK is rendered inline by formatCreateTable).
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
		if isMultiValuedIndex(idx) {
			continue // inline in CREATE TABLE (formatCreateTable): OceanBase adds none later
		}
		ops = append(ops, addIndexOps(entry, table, idx, n)...)
	}
	return ops
}

// indexOpsForModifiedTable emits ADD / DROP / (DROP+ADD) ops from the per-index diff of a
// modified table, and a DROP and an ADD for an index the diff leaves alone that holds a
// column the plan drops and adds again (recreatedColumns): DROP COLUMN takes the column
// out of the index's key or its STORING, or the index with it.
func indexOpsForModifiedTable(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	table := tableIdent(entry.To)
	recreated := recreatedColumns(entry, n)
	inDiff := make(map[string]bool, len(entry.Indexes))
	var ops []MigrationOp
	for _, ie := range entry.Indexes {
		inDiff[toLower(ie.Name)] = true
		switch ie.Action {
		case DiffAdd:
			if ie.To == nil {
				continue
			}
			ops = append(ops, addIndexOps(entry, table, ie.To, n)...)
		case DiffDrop:
			if ie.From == nil {
				continue
			}
			ops = append(ops, dropIndexOp(entry, table, ie.From))
		case DiffModify:
			if ie.From == nil || ie.To == nil {
				continue
			}
			if ie.To.Primary || ie.From.Primary {
				ops = append(ops, primaryKeyModifyOp(entry, table, ie.To, n))
				continue
			}
			// An index that differs only in being visible keeps its data: ALTER INDEX says
			// so, where DROP and ADD would build it again.
			if onlyVisibilityDiffers(ie.From, ie.To, n) && !indexReadsAny(ie.From, recreated) {
				ops = append(ops, alterIndexVisibilityOp(entry, table, ie.To))
				continue
			}
			ops = append(ops, dropIndexOp(entry, table, ie.From))
			ops = append(ops, addIndexOps(entry, table, ie.To, n)...)
		}
	}
	if len(recreated) > 0 && entry.From != nil {
		kept := diffableIndexMap(entry.To)
		for _, idx := range orderedDiffableIndexes(entry.From) {
			to := kept[toLower(idx.Name)]
			if inDiff[toLower(idx.Name)] || to == nil || idx.Primary || !indexReadsAny(idx, recreated) {
				continue
			}
			ops = append(ops, dropIndexOp(entry, table, idx))
			ops = append(ops, addIndexOps(entry, table, to, n)...)
		}
	}
	return ops
}

// indexReadsAny reports whether idx holds one of the columns: as a key part, named or, for
// an expression part, found in its text, or as a column it stores. DROP COLUMN takes a
// stored column out of the index's STORING as it takes a key part out of the key.
func indexReadsAny(idx *Index, columns map[string]bool) bool {
	if len(columns) == 0 {
		return false
	}
	for _, name := range idx.Storing {
		if columns[toLower(name)] {
			return true
		}
	}
	for _, ic := range idx.Columns {
		if columns[toLower(ic.Name)] {
			return true
		}
		if ic.Expr != "" {
			for name := range columns {
				if strings.Contains(toLower(ic.Expr), name) {
					return true
				}
			}
		}
	}
	return false
}

// addIndexOps is addIndexOp, followed for an invisible index by the statement that hides
// it. OceanBase's ALTER TABLE ... ADD KEY does not take INVISIBLE, written plainly or as
// /*!80000 INVISIBLE */: the index comes out visible. CREATE TABLE, CREATE INDEX and ALTER
// INDEX do take it, so the index is added and then hidden.
// TestOracle_OceanBaseInvisibleIndex holds each of these to the engine.
func addIndexOps(entry *TableDiffEntry, table string, idx *Index, n *Normalizer) []MigrationOp {
	ops := []MigrationOp{addIndexOp(entry, table, idx, n)}
	if !idx.Visible && !idx.Primary && !idx.Vector {
		ops = append(ops, alterIndexVisibilityOp(entry, table, idx))
	}
	return ops
}

// alterIndexVisibilityOp builds the ALTER TABLE ... ALTER INDEX that makes idx as visible as
// it says. It sorts with the index's ADD and after it.
func alterIndexVisibilityOp(entry *TableDiffEntry, table string, idx *Index) MigrationOp {
	visibility := "VISIBLE"
	if !idx.Visible {
		visibility = "INVISIBLE"
	}
	return MigrationOp{
		Type:         OpAlterIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s ALTER INDEX %s %s", table, migrationQuoteIdent(idx.Name), visibility),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
	}
}

// onlyVisibilityDiffers reports whether two versions of an index are the same but for
// VISIBLE and INVISIBLE.
func onlyVisibilityDiffers(from, to *Index, n *Normalizer) bool {
	if from.Visible == to.Visible || from.Primary || to.Primary {
		return false
	}
	shown := *from
	shown.Visible = to.Visible
	return canonicalIndex(&shown, n) == canonicalIndex(to, n)
}

// addIndexOp builds an ALTER TABLE ... ADD <index> op. A PRIMARY KEY uses ADD PRIMARY KEY; every
// other kind uses ADD <UNIQUE|FULLTEXT|SPATIAL> KEY `name` (...). The op runs in PhaseMain at
// priorityIndex.
func addIndexOp(entry *TableDiffEntry, table string, idx *Index, n *Normalizer) MigrationOp {
	if idx.Vector {
		return addVectorIndexOp(entry, table, idx)
	}
	clause := "ADD " + formatIndexDefinition(idx, n)
	op := MigrationOp{
		Type:         OpAddIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", table, clause),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
	}
	if isMultiValuedIndex(idx) {
		op.Warning = fmt.Sprintf("OceanBase cannot add the multi-valued index %s to table %s (error 1235): it builds one only with its table",
			migrationQuoteIdent(idx.Name), table)
	}
	return op
}

// addVectorIndexOp builds the CREATE VECTOR INDEX a vector index needs, because OceanBase 4.3.5.6
// has no ALTER TABLE form of it: ADD VECTOR KEY is refused.
func addVectorIndexOp(entry *TableDiffEntry, table string, idx *Index) MigrationOp {
	return MigrationOp{
		Type:         OpAddIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          formatCreateVectorIndex(table, idx),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, idx.Name),
	}
}

// formatCreateVectorIndex renders a vector index as CREATE VECTOR INDEX, with its parameters
// where SHOW CREATE TABLE prints them, straight after the key parts.
func formatCreateVectorIndex(table string, idx *Index) string {
	var b strings.Builder
	fmt.Fprintf(&b, "CREATE VECTOR INDEX %s ON %s (", migrationQuoteIdent(idx.Name), table)
	cols := make([]string, 0, len(idx.Columns))
	for _, ic := range idx.Columns {
		cols = append(cols, migrationIndexColumn(ic))
	}
	b.WriteString(strings.Join(cols, ","))
	b.WriteString(")")
	if len(idx.VectorParams) > 0 {
		b.WriteString(" WITH (" + strings.Join(idx.VectorParams, ", ") + ")")
	}
	// The stored columns, so a changed set is applied (canonicalStoring).
	b.WriteString(formatStoring(idx))
	b.WriteString(formatIndexOptions(idx))
	return b.String()
}

// primaryKeyModifyOp renders a changed primary key as one statement, ALTER TABLE ... DROP
// PRIMARY KEY, ADD PRIMARY KEY (...), in PhaseMain at priorityIndex: after the column adds
// and modifies, so the new key's columns are there and NOT NULL.
//
// One statement, because OceanBase refuses the DROP alone where it takes the pair: on a
// table a foreign key references, and on one whose foreign key is over the key's own
// columns (4714 for both). It is also the one pair of clauses OceanBase takes together
// where a third makes it refuse (1235, Multiple complex DDLs about primary in single
// stmt), so the statement holds nothing else.
//
// What the old key forbids while it stands waits for this statement: dropping one of its
// columns and making one nullable (columnsThatWait). A key the plan drops outright is a
// DROP PRIMARY KEY in PhasePre (dropIndexOp), ahead of the column drops.
// TestOracle_OceanBaseAlterTableClauses holds each of these to the engine.
func primaryKeyModifyOp(entry *TableDiffEntry, table string, to *Index, n *Normalizer) MigrationOp {
	return MigrationOp{
		Type:         OpAddIndex,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s DROP PRIMARY KEY, ADD %s", table, formatIndexDefinition(to, n)),
		Phase:        PhaseMain,
		Priority:     priorityIndex,
		sortName:     indexSortName(entry.Database, entry.Name, to.Name),
	}
}

// priorityIndexDrop orders an index DROP within PhasePre BEFORE any column DROP. A DROP COLUMN
// removes an index that the dropped column leaves empty (a single-column index, or the last
// remaining column of a composite one), after which an explicit DROP INDEX for that index
// fails with errno 1091 ("Can't DROP 'k'; check that column/key exists"). Because the
// migration plan is applied one statement at a time, the index DROP must precede the column
// DROP. Column drops run at priorityColumn-1 (generated) / priorityColumn (plain) within
// PhasePre, and table drops at priorityTable=10; this sits below the column-drop range and
// above the table-drop, so the order is:
// table-drop (10) < index-drop (15) < generated-col-drop (19) < plain-col-drop (20).
// A dropped primary key is among the index drops, which is what lets its columns be
// dropped after it: OceanBase refuses to drop a column the key holds (1235).
const priorityIndexDrop = priorityColumn - 5

// dropIndexOp builds an ALTER TABLE ... DROP op. A PRIMARY KEY uses DROP PRIMARY KEY (it has no
// user name); every other index uses DROP INDEX `name`. The op runs in PhasePre at
// priorityIndexDrop so all index drops precede any add AND precede column drops (see
// priorityIndexDrop), and so a dropped index name is free for re-creation in the same plan.
func dropIndexOp(entry *TableDiffEntry, table string, idx *Index) MigrationOp {
	clause := "DROP INDEX " + migrationQuoteIdent(idx.Name)
	if idx.Primary {
		clause = "DROP PRIMARY KEY"
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
// emits it (SHOW CREATE drops the index-level value on 8.0), but the ADD carries what the
// target states — the readback drops it and canonicalIndex ignores it, so the round-trip is
// empty.
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

	// A fulltext index's parser, so a changed one is applied (canonicalFulltextParser).
	b.WriteString(formatFulltextParser(idx))
	// The stored columns, so a changed set is applied (canonicalStoring).
	b.WriteString(formatStoring(idx))
	b.WriteString(formatIndexOptions(idx))
	return b.String()
}

// formatStoring renders an index's STORING clause, naming each column as the table does.
func formatStoring(idx *Index) string {
	if len(idx.Storing) == 0 {
		return ""
	}
	names := make([]string, 0, len(idx.Storing))
	for _, name := range idx.Storing {
		if idx.Table != nil {
			if col := idx.Table.GetColumn(name); col != nil {
				name = col.Name
			}
		}
		names = append(names, migrationQuoteIdent(name))
	}
	return " STORING (" + strings.Join(names, ",") + ")"
}

// formatIndexOptions renders what follows an index's key parts and any parameters of its kind.
func formatIndexOptions(idx *Index) string {
	var b strings.Builder

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

	// Index scope, last, as OceanBase prints it. Without this a LOCAL to GLOBAL change
	// recreates the index at the default scope, so the change never lands and the next
	// diff reports it again. The diff folds an unstated LOCAL away (canonicalIndexScope),
	// so what reaches here is a scope someone wrote.
	if idx.Scope != "" {
		b.WriteString(" " + strings.ToUpper(idx.Scope))
	}

	return b.String()
}

// orderedDiffableIndexes returns a table's diff-able indexes (the same set diffIndexes uses) in
// a deterministic order: by lower-cased name. Determinism keeps a new table's ADD INDEX sequence
// stable across runs.
func orderedDiffableIndexes(t *Table) []*Index {
	m := diffableIndexMap(t)
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
