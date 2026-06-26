package catalog

import (
	"fmt"
	"sort"
	"strings"
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
// Ordering has two cases, because a PARTITION BY and a REMOVE PARTITIONING have opposite
// dependencies on the column ops:
//
//   - PARTITION BY (define / repartition) runs in PhaseMain at priorityPartition — AFTER the
//     table CREATE (priorityTable), column adds/modifies (priorityColumn), and index adds
//     (priorityIndex), because PARTITION BY requires its referenced columns to exist and every
//     unique key (the PK included) to cover the partitioning columns. It runs BEFORE deferred
//     FKs (PhasePost): a table cannot be partitioned while it is the child of a foreign key.
//   - REMOVE PARTITIONING runs in PhasePre at priorityRemovePartitioning — BEFORE any column
//     DROP (PhasePre, priorityColumn-1/priorityColumn). MySQL rejects dropping a column that the
//     LIVE partition function/key still references (errno 1505/1054-class), so the table must be
//     un-partitioned first when the same plan also drops an old partitioning column.
//
// Repartitioning a table whose OLD partition function references a column DROPPED in the same plan
// is handled by splitting the change across phases: a PhasePre REMOVE PARTITIONING (so the column
// drop is no longer blocked by the live partition function) followed by the PhaseMain PARTITION BY
// that establishes the new scheme (see partitionModifyOps / oldPartitionColumnDropped). Otherwise a
// table emits at most one partition op, ordered by table name (sortName).
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

// partitionModifyOps renders the partition change for a modified table:
//   - to unpartitioned  → a single PhasePre REMOVE PARTITIONING.
//   - to (re)partitioned → a PhaseMain PARTITION BY. When the same plan ALSO drops a column the
//     OLD partition function references, a PhasePre REMOVE PARTITIONING is prepended so the table
//     is un-partitioned before that column drop runs (MySQL rejects dropping a column bound by the
//     live partition function); the PhaseMain PARTITION BY then re-establishes the new scheme.
func partitionModifyOps(entry *TableDiffEntry, n *Normalizer) []MigrationOp {
	if entry.To == nil {
		return nil
	}
	if entry.To.Partitioning == nil {
		// from was partitioned (PartitionChanged true + to unpartitioned), so strip it.
		return []MigrationOp{removePartitioningOp(entry, entry.To)}
	}
	op, ok := partitionByOp(entry, entry.To, n)
	if !ok {
		return nil
	}
	// Repartition: if a column referenced by the OLD partitioning is dropped in this plan, the
	// PhasePre DROP COLUMN would fail while the table is still partitioned by it — strip the old
	// partitioning first (PhasePre), then apply the new scheme (PhaseMain).
	if oldPartitionColumnDropped(entry) {
		return []MigrationOp{removePartitioningOp(entry, entry.To), op}
	}
	return []MigrationOp{op}
}

// oldPartitionColumnDropped reports whether any column the FROM partitioning references is being
// DROPPED in this diff. A dropped partition-key column requires the old partitioning to be removed
// first (a column bound by the live partition function cannot be dropped). Column references are
// taken from the structured KEY/COLUMNS lists and, for expression partitioning (RANGE/LIST/HASH
// over an expr), from an identifier scan of the expression — a superset that may also match a
// column merely named inside the expression, which only ever causes an extra, harmless PhasePre
// REMOVE PARTITIONING (the new PARTITION BY re-establishes the scheme regardless).
func oldPartitionColumnDropped(entry *TableDiffEntry) bool {
	if entry.From == nil || entry.From.Partitioning == nil {
		return false
	}
	dropped := droppedColumnSet(entry)
	if len(dropped) == 0 {
		return false
	}
	for _, col := range partitionReferencedColumns(entry.From.Partitioning) {
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

// partitionReferencedColumns returns the column names a partition spec depends on: the explicit
// KEY / RANGE COLUMNS / LIST COLUMNS column lists (and subpartition columns), plus identifiers
// scanned out of an expression-based partition/subpartition expression (RANGE/LIST/HASH over an
// expr). The expression scan is a conservative superset (it may include non-column identifiers
// such as function names), which is safe here: a false positive only adds a harmless REMOVE
// PARTITIONING before the repartition.
func partitionReferencedColumns(pi *PartitionInfo) []string {
	var cols []string
	cols = append(cols, pi.Columns...)
	cols = append(cols, pi.SubColumns...)
	cols = append(cols, identifiersInExpr(pi.Expr)...)
	cols = append(cols, identifiersInExpr(pi.SubExpr)...)
	return cols
}

// identifiersInExpr extracts identifier-like tokens from a partition expression, skipping
// backtick-quoted and string-literal content boundaries. It is used only to spot a dropped column
// referenced by an expression partition function; over-matching (e.g. catching a function name)
// is harmless (see partitionReferencedColumns).
func identifiersInExpr(expr string) []string {
	if expr == "" {
		return nil
	}
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() > 0 {
			out = append(out, b.String())
			b.Reset()
		}
	}
	for _, r := range expr {
		switch {
		case r == '_' || r == '$' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
		default:
			// Backticks, quotes, parentheses, operators, spaces all terminate an identifier run.
			flush()
		}
	}
	flush()
	return out
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

// removePartitioningOp builds an ALTER TABLE ... REMOVE PARTITIONING op. It runs in PhasePre at
// priorityRemovePartitioning, ahead of any column DROP, so a table partitioned by a column that is
// also dropped in the same plan is un-partitioned before the drop (which MySQL would otherwise
// reject — the column is still bound by the live partition function).
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

// priorityPartition orders a PARTITION BY op within PhaseMain after the table create
// (priorityTable), column alters (priorityColumn), and index adds (priorityIndex) — PARTITION BY
// needs the columns present and every unique key (PK) to cover the partitioning columns — and
// before deferred FKs (PhasePost), since a table cannot be partitioned while it is an FK child. It
// sits just above priorityForeignKey so partitioning is the last in-phase structural change before
// FKs.
const priorityPartition = priorityForeignKey - 1

// priorityRemovePartitioning orders a REMOVE PARTITIONING op within PhasePre BEFORE every column
// DROP. A column referenced by the live partition function cannot be dropped while the table is
// partitioned, so the stripping must run first. Column drops run at priorityColumn-1 (generated) /
// priorityColumn (plain) and index/PK drops at priorityIndexDrop/priorityPrimaryKeyDrop; this sits
// below all of them so REMOVE PARTITIONING is the very first PhasePre statement. (Dropping
// partitioning first is harmless when no partitioning column is dropped, so ordering it first
// unconditionally is safe.)
const priorityRemovePartitioning = priorityPrimaryKeyDrop - 1
