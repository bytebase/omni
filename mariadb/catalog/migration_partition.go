package catalog

import (
	"fmt"
	"sort"
)

// generatePartitionDDL turns the partition signal (TableDiffEntry.PartitionChanged, plus the
// partitioning of added tables) into DDL. A changed spec is re-applied wholesale with ALTER TABLE
// ... PARTITION BY ..., or removed with REMOVE PARTITIONING when the target is not partitioned.
// Individual partition add/drop/reorganize is not attempted because the diff cannot tell which
// partitions changed.
//
// The PARTITION BY clause is rendered by showPartitioning over a copy whose per-partition engine
// is resolved to the table's engine (partitionSpecWithResolvedEngine, shared with the diff), so
// it is the clause the diff compares.
//
// PARTITION BY runs in PhaseMain at priorityPartition, after the table CREATE, column adds and
// modifies, and index adds, because it needs its columns to exist and every unique key to cover
// the partitioning columns. It runs before deferred FKs because a table cannot be partitioned
// while it is the child of a foreign key. REMOVE PARTITIONING runs in PhasePre at
// priorityRemovePartitioning, before any column DROP, because a column the live partition
// function references cannot be dropped.
//
// A repartition that drops a column of the old partition function is split into a PhasePre
// REMOVE PARTITIONING and a PhaseMain PARTITION BY (see partitionModifyOps).
func generatePartitionDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			// formatCreateTable does not render the partition clause.
			if op, ok := partitionByOp(entry, entry.To, n); ok {
				ops = append(ops, op)
			}
		case DiffModify:
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

// partitionModifyOps renders the partition change for a modified table: a PhasePre REMOVE
// PARTITIONING when the target is unpartitioned, otherwise a PhaseMain PARTITION BY, preceded by
// a REMOVE PARTITIONING when the same plan drops a column the old partition function references.
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
	if oldPartitionColumnDropped(entry) {
		return []MigrationOp{removePartitioningOp(entry, entry.To), op}
	}
	return []MigrationOp{op}
}

// oldPartitionColumnDropped reports whether dropping a column in this diff could be blocked by
// the old partition function, so the partitioning must be stripped first. COLUMNS and KEY
// partitioning list their columns, which are matched case-insensitively against the dropped set.
// EXPRESSION partitioning embeds columns in SQL text that cannot be extracted reliably, so any
// column drop counts; an unneeded REMOVE PARTITIONING is harmless because the PARTITION BY
// re-establishes the scheme.
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

// partitionByOp builds an ALTER TABLE ... PARTITION BY ... op that (re)defines tbl's
// partitioning. It returns ok=false when the table is not partitioned.
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

// removePartitioningOp builds an ALTER TABLE ... REMOVE PARTITIONING op in PhasePre at
// priorityRemovePartitioning, ahead of any column DROP.
func removePartitioningOp(entry *TableDiffEntry, tbl *Table) MigrationOp {
	return MigrationOp{
		Type:         OpAlterTable,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s REMOVE PARTITIONING", tableIdent(tbl)),
		Phase:        PhasePre,
		Priority:     priorityRemovePartitioning,
		sortName:     tableSortName(entry.Database, entry.Name),
	}
}

// priorityPartition orders a PARTITION BY op within PhaseMain after table, column, and index
// ops and before deferred FKs (PhasePost).
const priorityPartition = priorityForeignKey - 1

// priorityRemovePartitioning orders a REMOVE PARTITIONING op within PhasePre before every column
// and index/PK drop, because a column referenced by the live partition function cannot be
// dropped while the table is partitioned.
const priorityRemovePartitioning = priorityPrimaryKeyDrop - 1
