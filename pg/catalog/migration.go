package catalog

import (
	"fmt"
	"sort"
	"strings"
)

// MigrationPhase classifies when a DDL operation should be executed
// relative to other operations.
type MigrationPhase int

const (
	PhasePre  MigrationPhase = iota // DROP operations
	PhaseMain                       // CREATE + ALTER operations
	PhasePost                       // Deferred (FK constraints)
)

// Priority constants for tie-breaking within a phase during topological sort.
// Lower values are executed earlier.
const (
	PrioritySchema     = 0
	PriorityExtension  = 1
	PriorityType       = 2  // Enum/Domain/Range/Composite
	PrioritySequence   = 3
	PriorityFunction   = 4
	PriorityTable      = 5
	PriorityColumn     = 6  // uses parent table OID
	PriorityConstraint = 7  // non-FK
	PriorityView       = 8
	PriorityIndex      = 9
	PriorityTrigger    = 10
	PriorityPolicy     = 11
	PriorityMetadata   = 12 // Comment/Grant/Revoke
	PriorityFKDeferred = 99 // FK constraint (PhasePost)
)

// MigrationOpType classifies a single DDL operation.
type MigrationOpType string

const (
	OpCreateSchema    MigrationOpType = "CreateSchema"
	OpDropSchema      MigrationOpType = "DropSchema"
	OpAlterSchema     MigrationOpType = "AlterSchema"
	OpCreateTable     MigrationOpType = "CreateTable"
	OpDropTable       MigrationOpType = "DropTable"
	OpAddColumn       MigrationOpType = "AddColumn"
	OpDropColumn      MigrationOpType = "DropColumn"
	OpAlterColumn     MigrationOpType = "AlterColumn"
	OpAddConstraint   MigrationOpType = "AddConstraint"
	OpDropConstraint  MigrationOpType = "DropConstraint"
	OpCreateIndex     MigrationOpType = "CreateIndex"
	OpDropIndex       MigrationOpType = "DropIndex"
	OpCreateSequence  MigrationOpType = "CreateSequence"
	OpDropSequence    MigrationOpType = "DropSequence"
	OpAlterSequence   MigrationOpType = "AlterSequence"
	OpCreateFunction  MigrationOpType = "CreateFunction"
	OpDropFunction    MigrationOpType = "DropFunction"
	OpAlterFunction   MigrationOpType = "AlterFunction"
	OpCreateType      MigrationOpType = "CreateType"
	OpDropType        MigrationOpType = "DropType"
	OpAlterType       MigrationOpType = "AlterType"
	OpCreateTrigger   MigrationOpType = "CreateTrigger"
	OpDropTrigger     MigrationOpType = "DropTrigger"
	OpCreateView      MigrationOpType = "CreateView"
	OpDropView        MigrationOpType = "DropView"
	OpAlterView       MigrationOpType = "AlterView"
	OpCreateExtension MigrationOpType = "CreateExtension"
	OpDropExtension   MigrationOpType = "DropExtension"
	OpAlterExtension  MigrationOpType = "AlterExtension"
	OpCreatePolicy    MigrationOpType = "CreatePolicy"
	OpDropPolicy      MigrationOpType = "DropPolicy"
	OpAlterPolicy     MigrationOpType = "AlterPolicy"
	OpAlterTable      MigrationOpType = "AlterTable"
	OpComment         MigrationOpType = "Comment"
	OpGrant           MigrationOpType = "Grant"
	OpRevoke          MigrationOpType = "Revoke"
)

// MigrationOp represents a single DDL operation in a migration plan.
type MigrationOp struct {
	Type          MigrationOpType
	SchemaName    string
	ObjectName    string
	SQL           string
	Warning       string
	Transactional bool
	ParentObject  string // optional, for FK deferred creation

	// Metadata for dependency-driven ordering (populated but not yet used by GenerateMigration).
	Phase    MigrationPhase // PhasePre, PhaseMain, or PhasePost
	ObjType  byte           // 'r'=relation, 'f'=function, 'i'=index, 'c'=constraint, 't'=type, 'S'=sequence, 'n'=schema, 'T'=trigger, 'p'=policy, 'e'=extension
	ObjOID   uint32         // OID in source catalog (from for DROP, to for CREATE)
	Priority int            // tie-breaker for topo sort (lower = earlier)
}

// MigrationPlan holds an ordered list of DDL operations that transform
// one catalog state into another.
type MigrationPlan struct {
	Ops []MigrationOp
}

// SQL returns all operations joined with ";\n".
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

// Summary returns a human-readable summary grouping operations by type
// and counting create/drop/alter categories.
func (p *MigrationPlan) Summary() string {
	if len(p.Ops) == 0 {
		return "No changes"
	}

	counts := make(map[MigrationOpType]int)
	for _, op := range p.Ops {
		counts[op.Type]++
	}

	// Sort keys for determinism.
	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, string(k))
	}
	sort.Strings(keys)

	var creates, drops, alters int
	var lines []string
	for _, k := range keys {
		typ := MigrationOpType(k)
		n := counts[typ]
		lines = append(lines, fmt.Sprintf("  %s: %d", k, n))
		s := string(typ)
		switch {
		case strings.HasPrefix(s, "Create") || strings.HasPrefix(s, "Add"):
			creates += n
		case strings.HasPrefix(s, "Drop") || strings.HasPrefix(s, "Revoke"):
			drops += n
		default:
			alters += n
		}
	}

	header := fmt.Sprintf("%d operation(s): %d create, %d drop, %d alter",
		len(p.Ops), creates, drops, alters)
	return header + "\n" + strings.Join(lines, "\n")
}

// Filter returns a new MigrationPlan containing only operations for which
// fn returns true.
func (p *MigrationPlan) Filter(fn func(MigrationOp) bool) *MigrationPlan {
	var ops []MigrationOp
	for _, op := range p.Ops {
		if fn(op) {
			ops = append(ops, op)
		}
	}
	return &MigrationPlan{Ops: ops}
}

// HasWarnings returns true if any operation has a non-empty warning.
func (p *MigrationPlan) HasWarnings() bool {
	for _, op := range p.Ops {
		if op.Warning != "" {
			return true
		}
	}
	return false
}

// Warnings returns all operations that have a non-empty warning.
func (p *MigrationPlan) Warnings() []MigrationOp {
	var result []MigrationOp
	for _, op := range p.Ops {
		if op.Warning != "" {
			result = append(result, op)
		}
	}
	return result
}

// GenerateMigration produces a MigrationPlan that transforms the `from`
// catalog into the `to` catalog, using the precomputed diff.
func GenerateMigration(from, to *Catalog, diff *SchemaDiff) *MigrationPlan {
	var ops []MigrationOp

	// Phase 1: DROP (reverse dependency order)
	// Phase 2: CREATE (forward dependency order)
	// Phase 3: ALTER (modify existing)
	// Phase 4: Metadata

	ops = append(ops, generateSchemaDDL(from, to, diff)...)
	ops = append(ops, generateExtensionDDL(from, to, diff)...)
	ops = append(ops, generateEnumDDL(from, to, diff)...)
	ops = append(ops, generateDomainDDL(from, to, diff)...)
	ops = append(ops, generateRangeDDL(from, to, diff)...)
	ops = append(ops, generateSequenceDDL(from, to, diff)...)

	// Split function DDL: functions that are referenced by table CHECK constraints
	// must be created before the tables. Functions that reference tables (e.g.,
	// RETURNS SETOF table, or body references table) must come after tables.
	funcOps := generateFunctionDDL(from, to, diff)
	earlyFuncOps, lateFuncOps := splitFunctionOps(to, diff, funcOps)
	ops = append(ops, earlyFuncOps...)
	ops = append(ops, generateTableDDL(from, to, diff)...)
	ops = append(ops, generateColumnDDL(from, to, diff)...)
	ops = append(ops, generateConstraintDDL(from, to, diff)...)
	ops = append(ops, lateFuncOps...)
	ops = append(ops, generateViewDDL(from, to, diff)...)
	ops = append(ops, generateIndexDDL(from, to, diff)...)
	ops = append(ops, generatePartitionDDL(from, to, diff)...)
	ops = append(ops, generateTriggerDDL(from, to, diff)...)
	ops = append(ops, generatePolicyDDL(from, to, diff)...)
	ops = append(ops, generateCommentDDL(from, to, diff)...)
	ops = append(ops, generateGrantDDL(from, to, diff)...)

	// Post-processing: wrap ALTER COLUMN TYPE ops with DROP VIEW / CREATE VIEW
	// for dependent views. PG cannot ALTER a column type when a view depends on it.
	ops = wrapColumnTypeChangesWithViewOps(from, to, diff, ops)

	return &MigrationPlan{Ops: ops}
}

// splitFunctionOps separates function DDL into two groups:
// - early: functions that are needed by table CHECK/DEFAULT constraints (must be created before tables)
// - late: all other functions (may reference tables, so created after tables)
func splitFunctionOps(to *Catalog, diff *SchemaDiff, funcOps []MigrationOp) (early, late []MigrationOp) {
	// Collect function names referenced by CHECK constraints and DEFAULTs
	// in tables being added or modified.
	neededByTables := make(map[string]bool)
	for _, rel := range diff.Relations {
		if rel.Action != DiffAdd {
			continue
		}
		if rel.To == nil {
			continue
		}
		// Check constraints on the table.
		cons := to.ConstraintsOf(rel.To.OID)
		for _, c := range cons {
			if c.Type == ConstraintCheck && c.CheckExpr != "" {
				// Extract function references from CHECK expression.
				// Simple heuristic: look for qualified function calls.
				for _, fn := range diff.Functions {
					if fn.Action == DiffAdd && fn.To != nil {
						// Check if function name appears in the CHECK expression.
						funcName := fn.To.Name
						if strings.Contains(c.CheckExpr, funcName+"(") || strings.Contains(c.CheckExpr, funcName+" (") {
							neededByTables[fn.Identity] = true
						}
					}
				}
			}
		}
		// Also check column defaults.
		for _, col := range rel.To.Columns {
			if col.HasDefault && col.Default != "" {
				for _, fn := range diff.Functions {
					if fn.Action == DiffAdd && fn.To != nil {
						funcName := fn.To.Name
						if strings.Contains(col.Default, funcName+"(") || strings.Contains(col.Default, funcName+" (") {
							neededByTables[fn.Identity] = true
						}
					}
				}
			}
		}
	}

	for _, op := range funcOps {
		// Only CREATE operations can be early; DROP/ALTER operations must be late
		// because they may depend on tables being dropped first (e.g., triggers).
		if op.Type == OpCreateFunction && neededByTables[op.ObjectName] {
			early = append(early, op)
		} else {
			late = append(late, op)
		}
	}
	return early, late
}

// wrapColumnTypeChangesWithViewOps detects ALTER COLUMN TYPE operations and
// inserts DROP VIEW before and CREATE VIEW after them for any views that
// depend on the modified table. PG cannot alter a column type if a view
// references that table.
func wrapColumnTypeChangesWithViewOps(from, to *Catalog, diff *SchemaDiff, ops []MigrationOp) []MigrationOp {
	// Find all tables that have column type changes.
	tablesWithTypeChange := make(map[string]bool) // "schema.table"
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		for _, col := range rel.Columns {
			if col.Action != DiffModify || col.From == nil || col.To == nil {
				continue
			}
			oldType := from.FormatType(col.From.TypeOID, col.From.TypeMod)
			newType := to.FormatType(col.To.TypeOID, col.To.TypeMod)
			if oldType != newType {
				tablesWithTypeChange[rel.SchemaName+"."+rel.Name] = true
			}
		}
	}
	if len(tablesWithTypeChange) == 0 {
		return ops
	}

	// Find views that depend on these tables.
	type viewInfo struct {
		schema string
		name   string
	}
	var viewsToDrop []viewInfo
	seen := make(map[string]bool)
	for tableKey := range tablesWithTypeChange {
		parts := strings.SplitN(tableKey, ".", 2)
		if len(parts) != 2 {
			continue
		}
		tableName := parts[1]
		for _, s := range to.UserSchemas() {
			for _, rel := range s.Relations {
				if rel.RelKind != 'v' {
					continue
				}
				viewKey := s.Name + "." + rel.Name
				if seen[viewKey] {
					continue
				}
				def, err := to.GetViewDefinition(s.Name, rel.Name)
				if err != nil {
					continue
				}
				if strings.Contains(def, quoteIdentAlways(tableName)) || strings.Contains(def, tableName) {
					viewsToDrop = append(viewsToDrop, viewInfo{schema: s.Name, name: rel.Name})
					seen[viewKey] = true
				}
			}
		}
	}
	if len(viewsToDrop) == 0 {
		return ops
	}

	// Sort views for determinism.
	sort.Slice(viewsToDrop, func(i, j int) bool {
		if viewsToDrop[i].schema != viewsToDrop[j].schema {
			return viewsToDrop[i].schema < viewsToDrop[j].schema
		}
		return viewsToDrop[i].name < viewsToDrop[j].name
	})

	// Build drop and create ops.
	var dropOps, createOps []MigrationOp
	for _, v := range viewsToDrop {
		qn := migrationQualifiedName(v.schema, v.name)
		// Resolve view OID for metadata.
		var viewOID uint32
		rel := to.GetRelation(v.schema, v.name)
		if rel != nil {
			viewOID = rel.OID
		}
		dropOps = append(dropOps, MigrationOp{
			Type:          OpDropView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			SQL:           fmt.Sprintf("DROP VIEW IF EXISTS %s", qn),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        viewOID,
			Priority:      PriorityView,
		})

		def, _ := to.GetViewDefinition(v.schema, v.name)
		createOps = append(createOps, MigrationOp{
			Type:          OpCreateView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			SQL:           fmt.Sprintf("CREATE VIEW %s AS %s", qn, strings.TrimRight(def, " \t\n\r;")),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'r',
			ObjOID:        viewOID,
			Priority:      PriorityView,
		})
	}

	// Filter out any existing view ops for these views (avoid duplicates).
	var filteredOps []MigrationOp
	for _, op := range ops {
		skip := false
		if op.Type == OpDropView || op.Type == OpCreateView || op.Type == OpAlterView {
			for _, v := range viewsToDrop {
				if op.SchemaName == v.schema && op.ObjectName == v.name {
					skip = true
					break
				}
			}
		}
		if !skip {
			filteredOps = append(filteredOps, op)
		}
	}

	// Insert: dropOps at the beginning, createOps at the end.
	var result []MigrationOp
	result = append(result, dropOps...)
	result = append(result, filteredOps...)
	result = append(result, createOps...)
	return result
}

// resolveTypeOIDByName looks up a type OID by schema + name from a catalog.
// Returns 0 if not found.
func resolveTypeOIDByName(c *Catalog, schemaName, typeName string) uint32 {
	if c == nil {
		return 0
	}
	schema := c.schemaByName[schemaName]
	if schema == nil {
		return 0
	}
	oid, _, err := c.ResolveType(TypeName{Schema: schemaName, Name: typeName, TypeMod: -1})
	if err != nil {
		return 0
	}
	return oid
}

// phaseForOpType returns the migration phase for a given op type.
func phaseForOpType(opType MigrationOpType) MigrationPhase {
	s := string(opType)
	if strings.HasPrefix(s, "Drop") {
		return PhasePre
	}
	return PhaseMain
}

