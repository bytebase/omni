package catalog

// MySQL declarative-schema (SDL) migration generator — tables + columns.
//
// This is the MySQL analog of the PostgreSQL reference at pg/catalog/migration.go. It
// turns a SchemaDiff (produced by diff-core) into ordered MySQL DDL that transforms the
// `from` (current/synced) catalog into the `to` (target/desired) catalog. The terminal
// output, MigrationPlan.SQL(), is what bytebase's mysqlDiffSDLMigration returns to the
// release executor.
//
// The defining correctness property is idempotence: when the diff is empty (target ==
// current), GenerateMigration emits an EMPTY plan and SQL() == "". A non-empty no-op plan
// is a bug — in ordering, rendering, or an upstream normalization gap — never a tolerable
// nit (correctness-protocol.md).
//
// Shape note (vs PG): GenerateMigration mirrors PG's dispatcher — one generate<obj>DDL per
// object kind, appended in dependency order, then sortMigrationOps. PG threads OIDs,
// schema namespaces, and a catalog dep graph through its topo sort. MySQL has none of
// that: identity is (database, table, column) by lower-cased name, and the table/column
// subset this node owns has no cross-object dependency to topo-sort (FK ordering is the
// FK breadth node's concern). So sortMigrationOps here is a deterministic phase + name
// sort, not a graph sort — which is all the table/column subset needs and is REQUIRED for
// idempotence/round-trip stability.
//
// Rendering principle: every column/type/option is rendered through normalize.go's
// canonical forms (CanonicalColumnType, ResolveColumnCharsetCollation, CanonicalNotNull,
// ...), NOT by re-stringifying the surface DDL. Because the emitted DDL IS the canonical
// stored form for the target version, apply-correctness holds by construction: apply the
// plan, read back SHOW CREATE, and the readback canonicalizes equal to `to`.
//
// Scope of THIS node (omni:generate-core): the MigrationPlan / MigrationOp types,
// GenerateMigration dispatcher, table DDL (migration_table.go), column DDL
// (migration_column.go), and sortMigrationOps. The breadth object kinds (indexes, FKs,
// constraints, checks, views, triggers, routines, events, partitions) are wired here against
// inert no-op generate stubs (each returns nil), mirroring PG's migration_<obj>.go layout;
// paired with the empty breadth SchemaDiff slices they contribute no ops, so a breadth node
// fills its slice and its hook with no change here.

import (
	"sort"
	"strings"
)

// MigrationPhase classifies when a DDL operation runs relative to others. For the
// table/column subset only two phases are used: drops run before adds/alters so a dropped
// table/column never collides with a re-created one of the same name. PhasePost is
// reserved for the FK breadth node (deferred FK creation), mirroring PG.
type MigrationPhase int

const (
	// PhasePre holds DROP operations (run first).
	PhasePre MigrationPhase = iota
	// PhaseMain holds CREATE + ALTER operations.
	PhaseMain
	// PhasePost holds deferred operations (FK constraints — breadth node).
	PhasePost
)

// MigrationOpType classifies a single MySQL DDL operation. MySQL uses MODIFY/CHANGE for
// column type/attribute changes (not PG's ALTER COLUMN ... TYPE), and folds index/PK/FK
// changes into ALTER TABLE rather than standalone statements — so the op-type set diverges
// from PG's. Only the table + column ops are emitted by this node; the rest are declared so
// breadth nodes and the drop-advice walker have a stable key set.
type MigrationOpType string

const (
	OpCreateTable     MigrationOpType = "CreateTable"
	OpDropTable       MigrationOpType = "DropTable"
	OpAlterTable      MigrationOpType = "AlterTable" // table-option change (ENGINE/CHARSET/...)
	OpAddColumn       MigrationOpType = "AddColumn"
	OpDropColumn      MigrationOpType = "DropColumn"
	OpModifyColumn    MigrationOpType = "ModifyColumn" // MySQL MODIFY COLUMN (in-place redefine)
	OpAddIndex        MigrationOpType = "AddIndex"     // breadth (index node)
	OpDropIndex       MigrationOpType = "DropIndex"    // breadth
	OpAddConstraint   MigrationOpType = "AddConstraint"
	OpDropConstraint  MigrationOpType = "DropConstraint"
	OpAddForeignKey   MigrationOpType = "AddForeignKey" // breadth (FK node)
	OpDropForeignKey  MigrationOpType = "DropForeignKey"
	OpAddCheck        MigrationOpType = "AddCheck" // breadth (check node)
	OpDropCheck       MigrationOpType = "DropCheck"
	OpCreateView      MigrationOpType = "CreateView" // breadth
	OpDropView        MigrationOpType = "DropView"
	OpCreateFunction  MigrationOpType = "CreateFunction" // breadth
	OpDropFunction    MigrationOpType = "DropFunction"
	OpCreateProcedure MigrationOpType = "CreateProcedure"
	OpDropProcedure   MigrationOpType = "DropProcedure"
	OpCreateTrigger   MigrationOpType = "CreateTrigger" // breadth
	OpDropTrigger     MigrationOpType = "DropTrigger"
	OpCreateEvent     MigrationOpType = "CreateEvent" // breadth
	OpDropEvent       MigrationOpType = "DropEvent"
)

// MigrationOp represents a single DDL operation in a migration plan. ObjectName is the
// table (or database-level object) name; Database is the qualifying scope (MySQL has no
// schema namespace). ParentObject names the owning table for sub-object ops (a column op's
// table) so the drop-advice walker can report "column x of table t". The ordering metadata
// (Phase, Priority, sortName) drives sortMigrationOps deterministically.
type MigrationOp struct {
	Type         MigrationOpType
	Database     string
	ObjectName   string
	SQL          string
	Warning      string
	ParentObject string // owning table for column/index/constraint ops; "" for table-level

	// Ordering metadata (used by sortMigrationOps).
	Phase    MigrationPhase
	Priority int    // tie-breaker within a phase (lower = earlier)
	sortName string // stable secondary key (table name, or table+column for column ops)

	// AUTO_INCREMENT-hazard grouping metadata (mergeAutoIncrementKeyOps,
	// migration_autoinc.go), set by the constructors of single-statement ALTER TABLE ops.
	// alterClause is the clause rendered after "ALTER TABLE <ident> " (empty for ops that
	// never participate in grouping); addsIndex/dropsIndex identify the index the op's own
	// statement adds / drops ("primary" stands for the PRIMARY KEY, which has no user name).
	alterClause string
	addsIndex   *Index
	dropsIndex  string
}

// Priority constants for tie-breaking within a phase. Lower runs earlier. For the
// table/column subset the only ordering that matters is: a table's own CREATE/DROP versus
// the column ALTERs on a surviving table. Breadth priorities are reserved so later nodes
// slot in without renumbering.
//
// Routines are direction-split (like priorityIndexDrop / priorityCheckDrop): their
// PhaseMain ops (CREATE/ALTER) run at priorityRoutine, BEFORE view creates, because
// CREATE VIEW eagerly validates the functions its body calls (Error 1305 when missing —
// live-verified on 8.0.32 and 5.7.25) while a routine body is lazily validated (a CREATE
// FUNCTION referencing a missing function/view/table is accepted on both versions), so
// hoisting every routine above every view is unconditionally safe and needs no per-edge
// dependency detection. Routine DROPs run at priorityRoutineDrop (migration_routine.go),
// AFTER view drops (priorityView) in PhasePre, so a dropped dependent view goes before
// the function it calls.
const (
	priorityTable      = 10
	priorityColumn     = 20
	priorityIndex      = 30
	priorityConstraint = 40
	priorityRoutine    = 45 // routine CREATE/ALTER (PhaseMain): before view creates
	priorityView       = 50
	priorityTrigger    = 70
	priorityEvent      = 80
	priorityForeignKey = 99 // FK deferred to PhasePost (breadth)
)

// MigrationPlan holds an ordered list of DDL operations that transform one catalog state
// into another. It is the MySQL analog of PG's MigrationPlan (pg/catalog/migration.go:104).
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL returns all operations joined with ";\n" (matching PG's MigrationPlan.SQL at
// pg/catalog/migration.go:109). An empty plan returns "" — the idempotence terminal: a
// no-op release emits no DDL.
func (p *MigrationPlan) SQL() string {
	if len(p.Ops) == 0 {
		return ""
	}
	parts := make([]string, len(p.Ops))
	for i, op := range p.Ops {
		parts[i] = op.SQL
	}
	return strings.Join(parts, ";\n")
}

// Filter returns a new MigrationPlan containing only operations for which fn returns true.
// bytebase uses it to strip ops it does not want to apply (PG strips the archive schema;
// pg/catalog/migration.go:163).
func (p *MigrationPlan) Filter(fn func(MigrationOp) bool) *MigrationPlan {
	var ops []MigrationOp
	for _, op := range p.Ops {
		if fn(op) {
			ops = append(ops, op)
		}
	}
	return &MigrationPlan{Ops: ops}
}

// HasWarnings reports whether any operation carries a non-empty warning.
func (p *MigrationPlan) HasWarnings() bool {
	for _, op := range p.Ops {
		if op.Warning != "" {
			return true
		}
	}
	return false
}

// Warnings returns the operations that carry a non-empty warning.
func (p *MigrationPlan) Warnings() []MigrationOp {
	var result []MigrationOp
	for _, op := range p.Ops {
		if op.Warning != "" {
			result = append(result, op)
		}
	}
	return result
}

// GenerateMigration produces a MigrationPlan that transforms the `from` catalog into the
// `to` catalog using the precomputed diff. It is the contract entry point
// (pg/catalog/migration.go:196) and uses the modern MySQL 8.0 canonical form by default,
// seeded — like Diff — with the `to` catalog's explicit_defaults_for_timestamp.
//
// Callers that know the target server version (bytebase's mysqlDiffSDLMigration, the oracle
// apply tests where 5.7-vs-8.0 stored forms diverge) should use
// GenerateMigrationWithNormalizer so the emitted DDL reproduces the version whose stored
// form the readback will canonicalize against — the same Diff/DiffWithNormalizer split
// diff-core established.
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	return GenerateMigrationWithNormalizer(from, to, diff, defaultNormalizer(to))
}

// GenerateMigrationWithNormalizer is the version-aware workhorse: the Normalizer fixes the
// target MySQL version (and explicit_defaults_for_timestamp) whose stored form the rendered
// DDL must reproduce, so apply-correctness holds on that engine. A nil normalizer falls
// back to the default for `to`. GenerateMigration delegates here.
//
// The dispatcher mirrors PG: one generate<obj>DDL per kind, appended in dependency order,
// then sortMigrationOps. Only the table generator (which descends into columns) is
// implemented in generate-core; the breadth hooks are present but inert until their node
// populates the matching SchemaDiff slice.
func GenerateMigrationWithNormalizer(from, to *Catalog, diff *SchemaDiff, n *Normalizer) *MigrationPlan {
	if diff == nil {
		return &MigrationPlan{}
	}
	if n == nil {
		n = defaultNormalizer(to)
	}

	var ops []MigrationOp

	// Table + column DDL (this node). Order of appends does not matter — sortMigrationOps
	// imposes the deterministic phase/name ordering.
	ops = append(ops, generateTableDDL(from, to, diff, n)...)
	ops = append(ops, generateColumnDDL(from, to, diff, n)...)

	// Breadth object kinds (mirroring PG's per-object generators and diff-core's empty slices).
	// Each generator is an inert no-op stub today (returns nil) — the matching SchemaDiff slice
	// is empty, so these contribute no ops and the plan stays empty on a no-op release. A breadth
	// node fills its slice and its generator with no change here.
	ops = append(ops, generateIndexDDL(from, to, diff, n)...)
	ops = append(ops, generateConstraintDDL(from, to, diff, n)...)
	ops = append(ops, generateForeignKeyDDL(from, to, diff, n)...)
	ops = append(ops, generateCheckDDL(from, to, diff, n)...)
	ops = append(ops, generateViewDDL(from, to, diff, n)...)
	ops = append(ops, generateRoutineDDL(from, to, diff, n)...)
	ops = append(ops, generateTriggerDDL(from, to, diff, n)...)
	ops = append(ops, generateEventDDL(from, to, diff, n)...)
	ops = append(ops, generatePartitionDDL(from, to, diff, n)...)

	// Group the clauses that MySQL requires to land in ONE statement: an AUTO_INCREMENT
	// column and its backing key must never be separated by a statement boundary (errno 1075
	// checks every ALTER's end state) — see migration_autoinc.go.
	ops = mergeAutoIncrementKeyOps(ops, diff, n)

	ops = sortMigrationOps(ops)

	return &MigrationPlan{Ops: ops}
}

// sortMigrationOps imposes a deterministic total order on the ops. Determinism is REQUIRED:
// idempotence/round-trip stability depends on the same diff always producing byte-identical
// DDL. The order is:
//
//  1. Phase: PhasePre (drops) before PhaseMain (creates/alters) before PhasePost (deferred FKs).
//     Within drops, a child (column) drop precedes its parent (table) drop is NOT needed —
//     DROP TABLE removes the whole table — but dropping a column on a SURVIVING table and
//     dropping a whole table are independent; phase ordering keeps all drops ahead of all
//     adds so a "drop table t / create table t" rename-in-place never reverses.
//  2. Priority within a phase: a table create/drop (priorityTable) before column alters
//     (priorityColumn) so CREATE TABLE precedes any ALTER ... ADD COLUMN referencing it
//     (defensive; for the current scope a freshly-added table carries its columns inline).
//  3. sortName lexicographic (database.table[.column]) for a stable tie-break.
//  4. Original index as the final stable tie-break (sort.SliceStable) so multiple ALTERs on
//     one column keep their emitted sequence (e.g. a MODIFY chain).
//
// This is the MySQL analog of PG's sortMigrationOps (pg/catalog/migration.go:699), reduced
// to a phase+name sort because the table/column subset has no cross-object dep graph.
func sortMigrationOps(ops []MigrationOp) []MigrationOp {
	sorted := make([]MigrationOp, len(ops))
	copy(sorted, ops)
	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i], sorted[j]
		if a.Phase != b.Phase {
			return a.Phase < b.Phase
		}
		if a.Priority != b.Priority {
			return a.Priority < b.Priority
		}
		if a.sortName != b.sortName {
			return a.sortName < b.sortName
		}
		return false // stable: preserve original order for equal keys
	})
	return sorted
}
