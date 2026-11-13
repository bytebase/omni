package catalog

import (
	"fmt"
	"sort"
)

// SDL generate — table partitioning (ALTER TABLE ... PARTITION BY).
//
// Turns the coarse partition signal (TableDiffEntry.PartitionChanged, plus the partitioning
// state of added/dropped tables) into the DDL that re-defines a table's partitioning. There
// is no in-place "modify one partition clause", so a changed partition spec is re-applied
// wholesale with `ALTER TABLE ... <PARTITION BY ...>`: the engine drops the old scheme and
// repartitions in one statement. The three cases:
//
//   - to partitioned, from not (or different spec)  → ALTER TABLE ... PARTITION BY ...   (define/repartition)
//   - from partitioned, to not                       → no statement on OceanBase (see below)
//   - both partitioned, same canonical spec          → nothing (diff said PartitionChanged=false)
//
// Add/drop/reorganize-partition (which preserve data more cheaply) are not attempted: the
// bool diff cannot tell which partitions changed, and the wholesale re-emit is always
// correct for a declarative apply.
//
// The PARTITION BY clause is rendered by showPartitioning — the SAME deparser show.go uses for
// the stored form — over a copy whose per-partition engine is resolved to the table's engine
// (partitionSpecWithResolvedEngine, shared with the diff), so the apply readback diffs empty.
//
// OceanBase cannot take a table's partitioning away. It has no REMOVE PARTITIONING (1064),
// will not drop the last partition, and keeps a table repartitioned into one partition
// partitioned. So:
//
//   - A table that is to lose its partitioning gets the REMOVE PARTITIONING MySQL would run,
//     with a Warning: the plan has nothing else to offer, and the statement fails where the
//     change cannot be made.
//   - A repartition that also drops a column the old partitioning reads cannot strip the
//     partitioning first, as MySQL's plan does. The new PARTITION BY runs before the key
//     changes instead (priorityRepartitionBeforeKeys): the old primary key still holds the
//     column, and a key without it would be refused while the old partitioning reads it
//     (1503). The key is then replaced and the column dropped (columnsThatWait).
//   - Any other PARTITION BY runs in PhaseMain at priorityPartition, after the column and
//     index changes, because it needs its columns to exist and every unique key to cover
//     them.
//
// TestOracle_OceanBaseAlterTableClauses holds each of these to the engine.
func generatePartitionDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			// A freshly created table carries its partitioning via a follow-up ALTER (the CREATE
			// from generate-core renders only columns + PK, not the partition clause), unless
			// its partitioning is one only CREATE TABLE takes, which formatCreateTable writes.
			if entry.To != nil && createOnlyPartitioning(entry.To.Partitioning) {
				continue
			}
			if op, ok := partitionByOp(entry, entry.To, n); ok {
				ops = append(ops, op)
			}
		case DiffModify:
			// Only emit when the spec actually changed (PartitionChanged folds the canonical
			// comparison; re-checking it keeps this generator inert on a no-op table).
			if !entry.PartitionChanged {
				continue
			}
			ops = append(ops, partitionModifyOps(entry, n)...)
		case DiffDrop:
			// DROP TABLE removes the partitioning with the table; nothing to emit.
		}
	}

	sort.SliceStable(ops, func(i, j int) bool {
		return ops[i].sortName < ops[j].sortName
	})

	return ops
}

// partitionModifyOps renders the partition change for a modified table: the PARTITION BY of
// a table that stays or becomes partitioned, early when a column of the old partitioning is
// dropped (repartitionsBeforeKeys), and for a table that is to lose its partitioning the
// REMOVE PARTITIONING that OceanBase refuses, warned.
func partitionModifyOps(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	if entry.To.Partitioning == nil {
		return []MigrationOp{removePartitioningOp(entry, entry.To)}
	}
	op, ok := partitionByOp(entry, entry.To, n)
	if !ok {
		return nil
	}
	if createOnlyPartitioning(entry.To.Partitioning) {
		op.Warning = fmt.Sprintf("OceanBase takes SIZE (...) and PARTITIONS AUTO only in CREATE TABLE (error 1064 in ALTER TABLE): create table %s again with its partitioning", tableIdent(entry.To))
	}
	if repartitionsBeforeKeys(entry) {
		op.Priority = priorityRepartitionBeforeKeys
	}
	return []MigrationOp{op}
}

// repartitionsBeforeKeys reports whether a table's new PARTITION BY has to run before its
// key changes: when the table is partitioned on both sides and the plan drops a column the
// old partitioning reads. See the note on OceanBase at the top of this file.
func repartitionsBeforeKeys(entry *TableDiffEntry) bool {
	return entry.PartitionChanged && entry.From != nil && entry.To != nil &&
		entry.From.Partitioning != nil && entry.To.Partitioning != nil && oldPartitionColumnDropped(entry)
}

// oldPartitionColumnDropped reports whether dropping a column in this diff could be blocked by the
// table's OLD (from) partition function, in which case the table must be repartitioned first.
// A column the live partition function reads is not dropped (errno 3855), so this errs
// toward yes:
//
//   - COLUMNS / KEY partitioning lists its dependency columns verbatim in Columns/SubColumns; an
//     exact (case-insensitive) match against the dropped set is reliable, including for names with
//     special characters.
//   - EXPRESSION partitioning (RANGE/LIST/HASH over Expr/SubExpr) embeds its columns inside SQL
//     text whose identifiers cannot be extracted reliably for every column name (quoting, special
//     characters). Rather than parse the expression, ANY column drop while the from is
//     expression-partitioned counts. Repartitioning early where it was not needed is harmless
//     (and an expression-partitioned table dropping any column is already an uncommon,
//     high-risk change).
func oldPartitionColumnDropped(entry *TableDiffEntry) bool {
	if entry.From == nil || entry.From.Partitioning == nil {
		return false
	}
	dropped := droppedColumnSet(entry)
	if len(dropped) == 0 {
		return false
	}
	pi := entry.From.Partitioning
	if pi.Expr != "" || pi.SubExpr != "" {
		// Expression partitioning: cannot reliably extract the referenced columns from SQL text;
		// any drop conservatively strips first (harmless if unnecessary).
		return true
	}
	for _, col := range append(append([]string{}, pi.Columns...), pi.SubColumns...) {
		if dropped[toLower(col)] {
			return true
		}
	}
	return false
}

// droppedColumnSet returns the lower-cased names of columns dropped in this table diff.
func droppedColumnSet(entry *TableDiffEntry) map[string]bool {
	var set map[string]bool
	for _, ce := range entry.Columns {
		if ce.Action == DiffDrop {
			if set == nil {
				set = make(map[string]bool)
			}
			set[toLower(ce.Name)] = true
		}
	}
	return set
}

// createOnlyPartitioning reports whether a partitioning has a clause OceanBase takes
// only in CREATE TABLE: a RANGE split size, SIZE ('size'), or automatic partitioning.
// ALTER TABLE's PARTITION BY refuses both (1064, simulate.tsv).
func createOnlyPartitioning(pi *PartitionInfo) bool {
	return pi != nil && (pi.SplitSize != "" || pi.AutoSize != "")
}

// partitionByOp builds an ALTER TABLE ... PARTITION BY ... op that (re)defines tbl's
// partitioning. Returns ok=false when the table is not partitioned (nothing to define).
func partitionByOp(entry *TableDiffEntry, tbl *Table, n *Normalizer) (MigrationOp, bool) {
	if tbl == nil || tbl.Partitioning == nil {
		return MigrationOp{}, false
	}
	clause := showPartitioning(partitionSpecWithResolvedEngine(tbl, n))
	return MigrationOp{
		Type:         OpAlterTable,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s %s", tableIdent(tbl), clause),
		Phase:        PhaseMain,
		Priority:     priorityPartition,
		sortName:     tableSortName(entry.Database, entry.Name),
	}, true
}

// removePartitioningOp builds the ALTER TABLE ... REMOVE PARTITIONING of a table that is to
// lose its partitioning. OceanBase answers it with a syntax error and has no other way to
// the same end, so the op carries a Warning. It runs first in PhasePre: nothing after it in
// the plan can be right once it has failed.
func removePartitioningOp(entry *TableDiffEntry, tbl *Table) MigrationOp {
	return MigrationOp{
		Type:         OpAlterTable,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s REMOVE PARTITIONING", tableIdent(tbl)),
		Warning:      fmt.Sprintf("OceanBase cannot remove the partitioning of table %s (error 1064): create the table again without it", tableIdent(tbl)),
		Phase:        PhasePre,
		Priority:     priorityRemovePartitioning,
		sortName:     tableSortName(entry.Database, entry.Name),
	}
}

// priorityPartition orders a PARTITION BY op within PhaseMain after the table create
// (priorityTable), column alters (priorityColumn), and index adds (priorityIndex) — PARTITION BY
// needs the columns present and every unique key (PK) to cover the partitioning columns — and
// before deferred FKs (PhasePost), since a table cannot be partitioned while it is an FK child. It
// sits just above priorityForeignKey so partitioning is the last in-phase structural change before
// FKs.
const priorityPartition = priorityForeignKey - 1

// priorityRepartitionBeforeKeys orders, within PhaseMain, the PARTITION BY of a repartition
// that lets a column go: after the column adds and modifies (priorityColumn, +1 for a
// generated column), so the new partitioning's columns exist, and before the index adds
// and the primary key's replacement (priorityIndex), which may leave out a column the old
// partitioning reads only once the table is partitioned by something else.
const priorityRepartitionBeforeKeys = priorityIndex - 1

// priorityRemovePartitioning orders a REMOVE PARTITIONING op first in PhasePre, ahead of the
// index and column drops (priorityIndexDrop, priorityColumn).
const priorityRemovePartitioning = priorityIndexDrop - 2
