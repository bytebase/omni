package catalog

// MariaDB declarative-schema (SDL) migration generator.
//
// GenerateMigration turns a SchemaDiff into ordered MariaDB DDL that transforms the `from`
// catalog into the `to` catalog. MigrationPlan.SQL() is the terminal output.
//
// When the diff is empty the plan is empty and SQL() == "". A non-empty no-op plan is a bug
// in ordering, rendering, or normalization.
//
// GenerateMigration calls one generate<obj>DDL per object kind, then sortMigrationOps.
// Identity is (database, table, column) by lower-cased name, and the sort is a deterministic
// phase + priority + name sort rather than a dependency graph.
//
// Columns, types, and options are rendered through the Normalizer's canonical forms rather
// than by re-stringifying the surface DDL, so the emitted DDL is what SHOW CREATE TABLE
// writes and the server stores exactly as the diff compares.

import (
	"sort"
	"strings"
)

// MigrationPhase classifies when a DDL operation runs relative to others. Drops run before
// adds/alters so a dropped object never collides with a re-created one of the same name.
type MigrationPhase int

const (
	// PhasePre holds DROP operations (run first).
	PhasePre MigrationPhase = iota
	// PhaseMain holds CREATE + ALTER operations.
	PhaseMain
	// PhasePost holds deferred operations such as foreign key constraints.
	PhasePost
)

// MigrationOpType classifies a single MariaDB DDL operation. The set is a stable key set for
// consumers such as drop-advice walkers.
type MigrationOpType string

const (
	OpCreateTable     MigrationOpType = "CreateTable"
	OpDropTable       MigrationOpType = "DropTable"
	OpAlterTable      MigrationOpType = "AlterTable" // table-option change
	OpAddColumn       MigrationOpType = "AddColumn"
	OpDropColumn      MigrationOpType = "DropColumn"
	OpModifyColumn    MigrationOpType = "ModifyColumn" // MODIFY COLUMN (in-place redefine)
	OpAddIndex        MigrationOpType = "AddIndex"
	OpDropIndex       MigrationOpType = "DropIndex"
	OpAddConstraint   MigrationOpType = "AddConstraint"
	OpDropConstraint  MigrationOpType = "DropConstraint"
	OpAddForeignKey   MigrationOpType = "AddForeignKey"
	OpDropForeignKey  MigrationOpType = "DropForeignKey"
	OpAddCheck        MigrationOpType = "AddCheck"
	OpDropCheck       MigrationOpType = "DropCheck"
	OpCreateView      MigrationOpType = "CreateView"
	OpDropView        MigrationOpType = "DropView"
	OpCreateFunction  MigrationOpType = "CreateFunction"
	OpDropFunction    MigrationOpType = "DropFunction"
	OpCreateProcedure MigrationOpType = "CreateProcedure"
	OpDropProcedure   MigrationOpType = "DropProcedure"
	OpCreateTrigger   MigrationOpType = "CreateTrigger"
	OpDropTrigger     MigrationOpType = "DropTrigger"
	OpCreateEvent     MigrationOpType = "CreateEvent"
	OpDropEvent       MigrationOpType = "DropEvent"
	OpCreateSequence  MigrationOpType = "CreateSequence"
	OpDropSequence    MigrationOpType = "DropSequence"
	OpAlterSequence   MigrationOpType = "AlterSequence"
	OpAlterVersioning MigrationOpType = "AlterVersioning" // system versioning and periods
)

// MigrationOp represents a single DDL operation in a migration plan. ObjectName is the
// table (or database-level object) name; Database is the qualifying scope. ParentObject names the owning table for sub-object ops (a column op's
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

	// column is the lower-cased column a column op adds, drops, or modifies. alterParts is the
	// clause list of a statement mergeAutoIncrementKeyOps assembled from several ops; nil when the
	// statement holds only alterClause. positionColumnAdds reads both.
	column     string
	alterParts []alterPart

	// changesVersionedColumns marks a versioning op that also drops or changes
	// columns of a versioned table, as a DROP SYSTEM VERSIONING that drops the ROW
	// START and ROW END columns with it, which needs keepVersionedHistory.
	changesVersionedColumns bool
}

// alterPart is one clause of a grouped ALTER TABLE statement. addsColumn is the lower-cased column
// an ADD COLUMN clause adds, empty for any other clause.
type alterPart struct {
	clause      string
	addsColumn  string
	dropsColumn string
}

// joinAlterParts renders a clause list as it follows "ALTER TABLE <ident> ".
func joinAlterParts(parts []alterPart) string {
	clauses := make([]string, len(parts))
	for i, p := range parts {
		clauses[i] = p.clause
	}
	return strings.Join(clauses, ", ")
}

// Priority constants for tie-breaking within a phase. Lower runs earlier.
//
// Routine CREATE/ALTER runs at priorityRoutine, before view creates, so a view body can call
// a routine created in the same plan. Routine DROPs run at priorityRoutineDrop
// (migration_routine.go), after view drops in PhasePre, so a dropped dependent view goes
// before the routine it calls.
const (
	priorityTable      = 10
	priorityColumn     = 20
	priorityIndex      = 30
	priorityConstraint = 40
	priorityRoutine    = 45 // routine CREATE/ALTER (PhaseMain): before view creates
	priorityView       = 50
	priorityTrigger    = 70
	priorityEvent      = 80
	priorityForeignKey = 99 // FK deferred to PhasePost
)

// MigrationPlan holds an ordered list of DDL operations that transform one catalog state
// into another.
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL returns all operations joined with ";\n". An empty plan returns "".
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
// Callers use it to strip ops they do not want to apply.
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
// `to` catalog using the precomputed diff and the default Normalizer for `to`. Callers that
// fix the normalizer explicitly use GenerateMigrationWithNormalizer.
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	return GenerateMigrationWithNormalizer(from, to, diff, defaultNormalizer(to))
}

// GenerateMigrationWithNormalizer generates the plan with the Normalizer that fixes the
// stored form the rendered DDL must reproduce. A nil normalizer falls back to the default
// for `to`.
func GenerateMigrationWithNormalizer(from, to *Catalog, diff *SchemaDiff, n *Normalizer) *MigrationPlan {
	if diff == nil {
		return &MigrationPlan{}
	}
	if n == nil {
		n = defaultNormalizer(to)
	}

	var ops []MigrationOp

	// Append order does not matter; sortMigrationOps imposes the final order.
	ops = append(ops, generateTableDDL(from, to, diff, n)...)
	ops = append(ops, generateColumnDDL(from, to, diff, n)...)

	ops = append(ops, generateIndexDDL(from, to, diff, n)...)
	ops = append(ops, generateConstraintDDL(from, to, diff, n)...)
	ops = append(ops, generateForeignKeyDDL(from, to, diff, n)...)
	ops = append(ops, generateCheckDDL(from, to, diff, n)...)
	ops = append(ops, generateViewDDL(from, to, diff, n)...)
	ops = append(ops, generateRoutineDDL(from, to, diff, n)...)
	ops = append(ops, generateTriggerDDL(from, to, diff, n)...)
	ops = append(ops, generateEventDDL(from, to, diff, n)...)
	ops = append(ops, generatePartitionDDL(from, to, diff, n)...)
	ops = append(ops, generateSequenceDDL(diff)...)
	ops = append(ops, generateVersioningDDL(diff, n)...)

	// An AUTO_INCREMENT column and its backing key must land in one statement, because
	// MariaDB checks the key rule at the end of every ALTER; see migration_autoinc.go.
	ops = dropViewsBeforeSequences(ops, diff)
	ops = reorderTableChanges(ops, from, to, diff, n)
	ops = mergeAutoIncrementKeyOps(ops, diff, n)
	ops = reorderAddsAfterMerge(ops, diff)
	ops = foldEngineIntoAutoKey(ops, diff, n)

	ops = sortMigrationOps(ops)
	ops = positionColumnAdds(ops, diff)
	ops = restatePageChecksum(ops, diff)
	ops = keepVersionedHistory(ops, diff)
	ops = explicitTimestampDefaults(ops, diff)

	return &MigrationPlan{Ops: ops}
}

// sortMigrationOps imposes a deterministic total order on the ops, so the same diff always
// produces byte-identical DDL. The order is:
//
//  1. Phase: PhasePre (drops) before PhaseMain (creates/alters) before PhasePost (deferred FKs).
//  2. Priority within a phase, so a table create/drop precedes column alters on it.
//  3. sortName lexicographic (database.table[.column]).
//  4. Original index (sort.SliceStable), so multiple ALTERs on one column keep their emitted
//     sequence.
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
