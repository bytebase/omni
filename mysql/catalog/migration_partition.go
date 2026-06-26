package catalog

import (
	"fmt"
	"sort"
)

// MySQL SDL generate — table partitioning (ALTER TABLE ... PARTITION BY / REMOVE PARTITIONING).
//
// Turns the coarse partition signal (TableDiffEntry.PartitionChanged, plus the partitioning
// state of added/dropped tables) into the DDL that re-defines a table's partitioning. MySQL
// folds partitioning into ALTER TABLE; there is no in-place "modify one partition clause", so a
// changed partition spec is re-applied wholesale with `ALTER TABLE ... <PARTITION BY ...>` — the
// engine drops the old scheme and repartitions in one statement. The three cases:
//
//   - to partitioned, from not (or different spec)  → ALTER TABLE ... PARTITION BY ...   (define/repartition)
//   - from partitioned, to not                       → ALTER TABLE ... REMOVE PARTITIONING (strip)
//   - both partitioned, same canonical spec          → nothing (diff said PartitionChanged=false)
//
// A coarse re-emit of the full PARTITION BY clause is apply-correct: `ALTER TABLE t PARTITION
// BY ...` is accepted on both 5.7 and 8.0 and yields exactly the target scheme, so the readback
// canonicalizes equal to `to`. Add/drop/reorganize-partition (which preserve data more cheaply)
// are not attempted — the bool diff cannot tell which partitions changed, and the wholesale
// re-emit is always correct for a declarative apply.
//
// The PARTITION BY clause is rendered by showPartitioning — the SAME deparser show.go uses for
// the stored form — over a copy whose per-partition engine is resolved to the table's engine
// (partitionSpecWithResolvedEngine, shared with the diff). That guarantees the emitted clause is
// byte-for-byte what the engine stores, so the apply readback diffs empty. showPartitioning wraps
// the clause in a /*!50100 ... */ (or /*!50500 ... */ for RANGE/LIST COLUMNS) version comment;
// MySQL executes that comment as an ALTER suffix on every supported version (verified on live
// 5.7.25 + 8.0.32), so the same DDL is valid on both engines.
//
// Ordering: partition ops run in PhaseMain at priorityPartition — AFTER the table CREATE
// (priorityTable), the column adds/modifies (priorityColumn), and the index adds (priorityIndex),
// because PARTITION BY requires its referenced columns to exist and every unique key (the PK
// included) to cover the partitioning columns; the PK/indexes must therefore already be in place.
// They run BEFORE deferred FKs (PhasePost): a table cannot be partitioned while it is the child
// of a foreign key, so partitioning must be established before FKs are added. A table either
// gains/changes partitioning or loses it in a single plan, never both, so at most one partition
// op per table — ops are ordered by table name (sortName) alone.
//
// implemented by omni:partitions breadth node
func generatePartitionDDL(_, _ *Catalog, diff *SchemaDiff, n *Normalizer) []MigrationOp {
	var ops []MigrationOp

	for i := range diff.Tables {
		entry := &diff.Tables[i]
		switch entry.Action {
		case DiffAdd:
			// A freshly created table carries its partitioning via a follow-up ALTER (the CREATE
			// from generate-core renders only columns + PK, not the partition clause).
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

// partitionModifyOps renders the partition change for a modified table: REMOVE PARTITIONING when
// the target is unpartitioned, otherwise a wholesale PARTITION BY re-definition.
func partitionModifyOps(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	if entry.To.Partitioning == nil {
		// from was partitioned (PartitionChanged true + to unpartitioned), so strip it.
		return []MigrationOp{removePartitioningOp(entry, entry.To)}
	}
	if op, ok := partitionByOp(entry, entry.To, n); ok {
		return []MigrationOp{op}
	}
	return nil
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

// removePartitioningOp builds an ALTER TABLE ... REMOVE PARTITIONING op.
func removePartitioningOp(entry *TableDiffEntry, tbl *Table) MigrationOp {
	return MigrationOp{
		Type:         OpAlterTable,
		Database:     entry.Database,
		ObjectName:   entry.Name,
		ParentObject: entry.Name,
		SQL:          fmt.Sprintf("ALTER TABLE %s REMOVE PARTITIONING", tableIdent(tbl)),
		Phase:        PhaseMain,
		Priority:     priorityPartition,
		sortName:     tableSortName(entry.Database, entry.Name),
	}
}

// priorityPartition orders partition ops within PhaseMain after the table create (priorityTable),
// column alters (priorityColumn), and index adds (priorityIndex) — PARTITION BY needs the columns
// present and every unique key (PK) to cover the partitioning columns — and before deferred FKs
// (PhasePost), since a table cannot be partitioned while it is an FK child. It sits just above
// priorityForeignKey so partitioning is the last in-phase structural change before FKs.
const priorityPartition = priorityForeignKey - 1
