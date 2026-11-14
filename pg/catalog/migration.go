package catalog

import (
	"container/heap"
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
	// PhaseLateDrop holds the DROP operations that wait for a PhaseMain ALTER
	// to stop a column or attribute from using what they drop.
	PhaseLateDrop
)

// Priority constants for tie-breaking within a phase during topological sort.
// Lower values are executed earlier.
const (
	PrioritySchema       = 0
	PriorityExtension    = 1
	PriorityType         = 2 // Enum/Domain/Range/Composite
	PrioritySequence     = 3
	PriorityFunction     = 4
	PriorityTable        = 5
	PriorityColumn       = 6 // uses parent table OID
	PriorityConstraint   = 7 // non-FK
	PriorityView         = 8
	PriorityIndex        = 9
	PriorityTrigger      = 10
	PriorityEventTrigger = 10
	PriorityPolicy       = 11
	PriorityMetadata     = 12 // Comment/Grant/Revoke
	PriorityFKDeferred   = 99 // FK constraint (PhasePost)
)

// MigrationOpType classifies a single DDL operation.
type MigrationOpType string

const (
	OpCreateSchema       MigrationOpType = "CreateSchema"
	OpDropSchema         MigrationOpType = "DropSchema"
	OpAlterSchema        MigrationOpType = "AlterSchema"
	OpCreateTable        MigrationOpType = "CreateTable"
	OpDropTable          MigrationOpType = "DropTable"
	OpAddColumn          MigrationOpType = "AddColumn"
	OpDropColumn         MigrationOpType = "DropColumn"
	OpAlterColumn        MigrationOpType = "AlterColumn"
	OpAddConstraint      MigrationOpType = "AddConstraint"
	OpDropConstraint     MigrationOpType = "DropConstraint"
	OpCreateIndex        MigrationOpType = "CreateIndex"
	OpDropIndex          MigrationOpType = "DropIndex"
	OpCreateSequence     MigrationOpType = "CreateSequence"
	OpDropSequence       MigrationOpType = "DropSequence"
	OpAlterSequence      MigrationOpType = "AlterSequence"
	OpCreateFunction     MigrationOpType = "CreateFunction"
	OpDropFunction       MigrationOpType = "DropFunction"
	OpAlterFunction      MigrationOpType = "AlterFunction"
	OpCreateType         MigrationOpType = "CreateType"
	OpDropType           MigrationOpType = "DropType"
	OpAlterType          MigrationOpType = "AlterType"
	OpCreateTrigger      MigrationOpType = "CreateTrigger"
	OpDropTrigger        MigrationOpType = "DropTrigger"
	OpCreateEventTrigger MigrationOpType = "CreateEventTrigger"
	OpDropEventTrigger   MigrationOpType = "DropEventTrigger"
	OpAlterEventTrigger  MigrationOpType = "AlterEventTrigger"
	OpCreateView         MigrationOpType = "CreateView"
	OpDropView           MigrationOpType = "DropView"
	OpAlterView          MigrationOpType = "AlterView"
	OpCreateExtension    MigrationOpType = "CreateExtension"
	OpDropExtension      MigrationOpType = "DropExtension"
	OpAlterExtension     MigrationOpType = "AlterExtension"
	OpCreatePolicy       MigrationOpType = "CreatePolicy"
	OpDropPolicy         MigrationOpType = "DropPolicy"
	OpAlterPolicy        MigrationOpType = "AlterPolicy"
	OpAlterTable         MigrationOpType = "AlterTable"
	OpComment            MigrationOpType = "Comment"
	OpGrant              MigrationOpType = "Grant"
	OpRevoke             MigrationOpType = "Revoke"
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

	// Metadata for dependency-driven ordering (used by sortMigrationOps).
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
	ops = append(ops, generateCompositeDDL(from, to, diff)...)
	ops = append(ops, generateSequenceDDL(from, to, diff)...)

	// Function ordering is handled by the dependency-driven topological sort:
	// functions referenced by CHECK/DEFAULT deps are placed before tables,
	// functions with RETURNS SETOF deps are placed after tables.
	ops = append(ops, generateFunctionDDL(from, to, diff)...)
	ops = append(ops, generateTableDDL(from, to, diff)...)
	ops = append(ops, generateColumnDDL(from, to, diff)...)
	ops = append(ops, generateConstraintDDL(from, to, diff)...)
	ops = append(ops, generateViewDDL(from, to, diff)...)
	ops = append(ops, generateIndexDDL(from, to, diff)...)
	ops = append(ops, generatePartitionDDL(from, to, diff)...)
	ops = append(ops, generateTriggerDDL(from, to, diff)...)
	ops = append(ops, generateEventTriggerDDL(from, to, diff)...)
	ops = append(ops, generatePolicyDDL(from, to, diff)...)
	ops = append(ops, generateCommentDDL(from, to, diff)...)
	ops = append(ops, generateGrantDDL(from, to, diff)...)

	// Post-processing: wrap ALTER COLUMN TYPE ops with DROP VIEW / CREATE VIEW
	// for dependent views. PG cannot ALTER a column type when a view depends on it.
	ops = wrapColumnTypeChangesWithViewOps(from, to, diff, ops)
	ops = deferDropsOfTypesInUse(from, to, ops)

	// Dependency-driven ordering: sort ops by phase, then topologically within each phase.
	ops = sortMigrationOps(from, to, ops)

	return &MigrationPlan{Ops: ops}
}

// wrapColumnTypeChangesWithViewOps detects ALTER COLUMN TYPE operations and
// injects synthetic DROP VIEW + CREATE VIEW ops for any views that depend on
// the modified table (via catalog deps). PG cannot alter a column type if a
// view references that table.
//
// The ops are appended (not manually positioned); sortMigrationOps handles
// ordering via the dependency graph.
func wrapColumnTypeChangesWithViewOps(from, to *Catalog, diff *SchemaDiff, ops []MigrationOp) []MigrationOp {
	// Find OIDs of tables that have column type changes or column drops.
	tableOIDs := make(map[uint32]bool)
	for _, rel := range diff.Relations {
		if rel.Action != DiffModify {
			continue
		}
		needsViewWrap := false

		// RelKind change (e.g., table → partitioned) or inheritance change
		// trigger DROP TABLE CASCADE which destroys dependent views.
		if rel.From != nil && rel.To != nil {
			if rel.From.RelKind != rel.To.RelKind {
				needsViewWrap = true
			}
			if !inhParentsEqual(from, to, rel.From.InhParents, rel.To.InhParents) {
				needsViewWrap = true
			}
		}

		if !needsViewWrap {
			for _, col := range rel.Columns {
				// Column drop: PG cannot drop a column when a view depends on the table.
				if col.Action == DiffDrop {
					needsViewWrap = true
					break
				}
				// Column type change.
				if col.Action != DiffModify || col.From == nil || col.To == nil {
					continue
				}
				// Column type or collation change: both are ALTER COLUMN TYPE.
				oldType := from.FormatType(col.From.TypeOID, col.From.TypeMod)
				newType := to.FormatType(col.To.TypeOID, col.To.TypeMod)
				if oldType != newType || col.From.CollationName != col.To.CollationName {
					needsViewWrap = true
					break
				}
			}
		}
		if needsViewWrap {
			// Look up the table OID from the `from` catalog (where deps are recorded).
			r := from.GetRelation(rel.SchemaName, rel.Name)
			if r != nil {
				tableOIDs[r.OID] = true
			}
		}
	}
	if len(tableOIDs) == 0 {
		return ops
	}

	// Use from.deps to find views that depend on these tables (transitively).
	// A view depends on a table if there is a dep entry where RefType='r',
	// RefOID=tableOID, ObjType='r', and the dependent is a view (RelKind='v').
	type viewInfo struct {
		schema string
		name   string
		oid    uint32
	}
	seen := make(map[uint32]bool)
	var viewsToDrop []viewInfo

	// Collect views transitively: a view may depend on another view that
	// depends on the table. We do a BFS over the dep graph.
	queue := make([]uint32, 0, len(tableOIDs))
	for oid := range tableOIDs {
		queue = append(queue, oid)
	}
	for len(queue) > 0 {
		refOID := queue[0]
		queue = queue[1:]
		for _, d := range from.deps {
			if d.RefType != 'r' || d.RefOID != refOID || d.ObjType != 'r' {
				continue
			}
			if seen[d.ObjOID] {
				continue
			}
			rel := from.GetRelationByOID(d.ObjOID)
			if rel == nil || rel.RelKind != 'v' {
				continue
			}
			seen[d.ObjOID] = true
			if rel.Schema == nil {
				continue
			}
			viewsToDrop = append(viewsToDrop, viewInfo{schema: rel.Schema.Name, name: rel.Name, oid: rel.OID})
			// Enqueue this view's OID for transitive dependency discovery.
			queue = append(queue, d.ObjOID)
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

	// Build drop and create ops. Resolve OIDs from `to` catalog for CREATE.
	var extraOps []MigrationOp
	dropViews := make(map[string]bool) // "schema.name" for dedup
	for _, v := range viewsToDrop {
		dropViews[v.schema+"."+v.name] = true
		qn := migrationQualifiedName(v.schema, v.name)

		// Use from-catalog OID for DROP (PhasePre uses from deps).
		extraOps = append(extraOps, MigrationOp{
			Type:          OpDropView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			SQL:           fmt.Sprintf("DROP VIEW IF EXISTS %s", qn),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        v.oid,
			Priority:      PriorityView,
		})

		// Use to-catalog OID for CREATE (PhaseMain uses to deps).
		var toOID uint32
		toRel := to.GetRelation(v.schema, v.name)
		if toRel != nil {
			toOID = toRel.OID
		}
		def, _ := to.GetViewDefinition(v.schema, v.name)
		extraOps = append(extraOps, MigrationOp{
			Type:          OpCreateView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			SQL:           fmt.Sprintf("CREATE VIEW %s AS %s", qn, strings.TrimRight(def, " \t\n\r;")),
			Transactional: true,
			Phase:         PhaseMain,
			ObjType:       'r',
			ObjOID:        toOID,
			Priority:      PriorityView,
		})
	}

	// Filter out any existing view ops for these views (avoid duplicates).
	var filteredOps []MigrationOp
	for _, op := range ops {
		skip := false
		if op.Type == OpDropView || op.Type == OpCreateView || op.Type == OpAlterView {
			if dropViews[op.SchemaName+"."+op.ObjectName] {
				skip = true
			}
		}
		if !skip {
			filteredOps = append(filteredOps, op)
		}
	}

	// Append extra ops — sortMigrationOps handles ordering.
	return append(filteredOps, extraOps...)
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

// ---------------------------------------------------------------------------
// Dependency-driven migration ordering (Step 2)
// ---------------------------------------------------------------------------

// depKey uniquely identifies a catalog object for OID-to-op mapping.
type depKey struct {
	objType byte
	objOID  uint32
}

// migPQEntry is an entry in the priority queue for migration topological sort.
type migPQEntry struct {
	idx int // index into the ops slice
	pri int // Priority from MigrationOp
	ord int // original index for stable tie-break
}

// migPQHeap implements heap.Interface for Kahn's algorithm priority queue.
// Orders by (Priority ASC, original index ASC).
type migPQHeap []migPQEntry

func (h migPQHeap) Len() int      { return len(h) }
func (h migPQHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h migPQHeap) Less(i, j int) bool {
	if h[i].pri != h[j].pri {
		return h[i].pri < h[j].pri
	}
	return h[i].ord < h[j].ord
}
func (h *migPQHeap) Push(x interface{}) {
	*h = append(*h, x.(migPQEntry))
}
func (h *migPQHeap) Pop() interface{} {
	old := *h
	n := len(old)
	e := old[n-1]
	*h = old[:n-1]
	return e
}

// liftDepToOp maps a DepEntry's object-side (ObjType/ObjOID) to the migration
// op indices that own it. Constraints and indexes are lifted to their parent
// table's op. Type references to a relation's row type are lifted to the
// relation's op.
func liftDepToOp(c *Catalog, objType byte, objOID uint32, oidToIdx map[depKey][]int) []int {
	// Direct match first.
	if idxs, ok := oidToIdx[depKey{objType, objOID}]; ok {
		return idxs
	}
	// Lift constraint → owning table.
	if objType == 'c' {
		if con, ok := c.constraints[objOID]; ok {
			return oidToIdx[depKey{'r', con.RelOID}]
		}
		// Lift domain constraint → owning domain type.
		for _, dt := range c.domainTypes {
			for _, dc := range dt.Constraints {
				if dc.OID == objOID {
					return oidToIdx[depKey{'t', dt.TypeOID}]
				}
			}
		}
	}
	// Lift index → owning table.
	if objType == 'i' {
		for _, idxList := range c.indexesByRel {
			for _, idx := range idxList {
				if idx.OID == objOID {
					return oidToIdx[depKey{'r', idx.RelOID}]
				}
			}
		}
	}
	// Lift type → owning relation (row type of a table/view/composite).
	if objType == 't' {
		for _, rel := range c.relationByOID {
			if rel.RowTypeOID == objOID {
				return oidToIdx[depKey{'r', rel.OID}]
			}
		}
	}
	return nil
}

// topoSortOps sorts migration ops using Kahn's algorithm with priority
// tie-break. It uses catalog deps to build the adjacency graph.
// For forward (CREATE) sorting (reverse=false): if dep says A depends on B, B comes before A.
// For reverse (DROP) sorting (reverse=true): if dep says A depends on B, A comes before B
// (drop dependents first). Priority tie-break is negated so higher-priority objects
// (views=8, triggers=10) are dropped before lower-priority ones (tables=5).
//
// Returns (sorted, deferred): sorted ops in dependency order, and any ops that
// were deferred to PhasePost to break cycles (e.g., CHECK constraints).
func topoSortOps(c *Catalog, ops []MigrationOp, reverse bool) (sorted []MigrationOp, deferred []MigrationOp) {
	if len(ops) == 0 {
		return nil, nil
	}
	n := len(ops)

	// Build OID → op index map (multiple ops can share an OID).
	oidToIdx := make(map[depKey][]int)
	for i, op := range ops {
		if op.ObjOID != 0 {
			key := depKey{op.ObjType, op.ObjOID}
			oidToIdx[key] = append(oidToIdx[key], i)
		}
	}

	// Build adjacency from catalog deps.
	adj := make([][]int, n)
	inDegree := make([]int, n)
	edgeSeen := make(map[[2]int]bool) // avoid duplicate edges

	for _, d := range c.deps {
		if d.DepType == DepInternal {
			continue
		}

		// Skip FK constraint deps — FK ops are already in PhasePost and their
		// deps should not create edges in the PhaseMain graph.
		if d.ObjType == 'c' {
			if con, ok := c.constraints[d.ObjOID]; ok && con.Type == ConstraintFK {
				continue
			}
		}

		// Find ops for the dependent object (lift constraint/index → table).
		depIdxs := liftDepToOp(c, d.ObjType, d.ObjOID, oidToIdx)
		// Find ops for the referenced object (lift constraint/index → table).
		refIdxs := liftDepToOp(c, d.RefType, d.RefOID, oidToIdx)

		if len(depIdxs) == 0 || len(refIdxs) == 0 {
			continue
		}

		// Forward: referenced (refIdx) must come BEFORE dependent (depIdx).
		// Reverse: dependent (depIdx) must come BEFORE referenced (refIdx).
		for _, refIdx := range refIdxs {
			for _, depIdx := range depIdxs {
				if depIdx == refIdx {
					continue // self-reference
				}
				// Skip edges between ops that share the same owning object
				// (same ObjType + ObjOID). Multiple ops on the same table
				// (e.g., AddColumn + ALTER TABLE ENABLE RLS) would otherwise
				// form false bidirectional edges via deps that reference the
				// table from both sides (e.g., a constraint's DepAuto on its
				// own table lifted to two distinct op indices).
				if ops[depIdx].ObjType == ops[refIdx].ObjType && ops[depIdx].ObjOID == ops[refIdx].ObjOID {
					continue
				}
				var from, to int
				if reverse {
					from, to = depIdx, refIdx
				} else {
					from, to = refIdx, depIdx
				}
				edge := [2]int{from, to}
				if edgeSeen[edge] {
					continue
				}
				edgeSeen[edge] = true
				adj[from] = append(adj[from], to)
				inDegree[to]++
			}
		}
	}

	// Kahn's algorithm with min-heap (priority, original index).
	// For reverse sort, negate priority so higher-priority objects come first
	// (e.g., views=8 before tables=5 → -8 < -5).
	priOf := func(i int) int {
		if reverse {
			return -ops[i].Priority
		}
		return ops[i].Priority
	}

	h := make(migPQHeap, 0, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			h = append(h, migPQEntry{idx: i, pri: priOf(i), ord: i})
		}
	}
	heap.Init(&h)

	sortedIdxs := make([]int, 0, n)
	for h.Len() > 0 {
		e := heap.Pop(&h).(migPQEntry)
		sortedIdxs = append(sortedIdxs, e.idx)
		for _, next := range adj[e.idx] {
			inDegree[next]--
			if inDegree[next] == 0 {
				heap.Push(&h, migPQEntry{idx: next, pri: priOf(next), ord: next})
			}
		}
	}

	// Convert sorted indices to ops.
	sorted = make([]MigrationOp, 0, len(sortedIdxs))
	for _, idx := range sortedIdxs {
		sorted = append(sorted, ops[idx])
	}

	// Cycle detection: if Kahn's didn't consume all ops, there is a cycle.
	if len(sortedIdxs) < n {
		sortedSet := make(map[int]bool, len(sortedIdxs))
		for _, idx := range sortedIdxs {
			sortedSet[idx] = true
		}
		var remaining []int
		for i := 0; i < n; i++ {
			if !sortedSet[i] {
				remaining = append(remaining, i)
			}
		}
		// Try to break cycle by deferring CHECK constraints to PhasePost.
		var deferredIdxs []int
		var stillRemaining []int
		for _, idx := range remaining {
			op := ops[idx]
			if op.Type == OpAddConstraint && op.ObjType == 'c' && op.Phase == PhaseMain {
				deferredIdxs = append(deferredIdxs, idx)
			} else {
				stillRemaining = append(stillRemaining, idx)
			}
		}
		if len(deferredIdxs) > 0 && len(stillRemaining) == 0 {
			// All remaining ops are deferrable CHECK constraints.
			for _, idx := range deferredIdxs {
				op := ops[idx]
				op.Phase = PhasePost
				op.Priority = PriorityFKDeferred
				deferred = append(deferred, op)
			}
		} else {
			// Unresolvable cycle: append remaining ops sorted by priority with a
			// warning so that the migration still produces output (best-effort).
			sort.Slice(remaining, func(a, b int) bool {
				pa, pb := ops[remaining[a]].Priority, ops[remaining[b]].Priority
				if pa != pb {
					return pa < pb
				}
				return remaining[a] < remaining[b]
			})
			for _, idx := range remaining {
				op := ops[idx]
				if op.Warning == "" {
					op.Warning = "unresolvable dependency cycle detected"
				}
				sorted = append(sorted, op)
			}
		}
	}

	return sorted, deferred
}

// deferDropsOfTypesInUse moves to PhaseLateDrop each type drop whose type a
// column or composite attribute uses in from and still has in to: PhaseMain
// changes its type, and only then can the type go. The drops of what such a
// type uses, and of the schemas holding them, wait with it. A drop that the
// plan pairs with a create of the same name keeps its place, since the create
// needs the name free; so does one whose users the plan drops and rebuilds.
func deferDropsOfTypesInUse(from, to *Catalog, ops []MigrationOp) []MigrationOp {
	created := make(map[string]bool)
	rebuilt := make(map[uint32]bool)
	for _, op := range ops {
		switch op.Type {
		case OpCreateType:
			created[op.SchemaName+"."+op.ObjectName] = true
		case OpDropTable, OpDropType, OpDropView:
			if op.Phase == PhasePre && op.ObjType == 'r' {
				rebuilt[op.ObjOID] = true
			}
		}
	}
	// typeOID is the type a drop removes: the type itself, or a composite's
	// row type.
	typeOID := func(op MigrationOp) uint32 {
		if op.ObjType == 'r' {
			if rel := from.GetRelationByOID(op.ObjOID); rel != nil {
				return rel.RowTypeOID
			}
			return 0
		}
		return op.ObjOID
	}
	late := make(map[int]bool)
	for i, op := range ops {
		if op.Phase == PhasePre && op.Type == OpDropType && !created[op.SchemaName+"."+op.ObjectName] && typeUsedByKeptColumn(from, to, rebuilt, typeOID(op)) {
			late[i] = true
		}
	}
	if len(late) == 0 {
		return ops
	}

	// Close over what the late drops depend on, lifting each dependency to
	// the drop that removes its object, as topoSortOps does.
	oidToIdx := make(map[depKey][]int)
	for i, op := range ops {
		if op.Phase == PhasePre && op.ObjOID != 0 {
			oidToIdx[depKey{op.ObjType, op.ObjOID}] = append(oidToIdx[depKey{op.ObjType, op.ObjOID}], i)
		}
	}
	for changed := true; changed; {
		changed = false
		for _, d := range from.deps {
			if d.DepType == DepInternal {
				continue
			}
			for _, i := range liftDepToOp(from, d.ObjType, d.ObjOID, oidToIdx) {
				if !late[i] {
					continue
				}
				for _, j := range liftDepToOp(from, d.RefType, d.RefOID, oidToIdx) {
					if !late[j] && ops[j].Type == OpDropType {
						late[j], changed = true, true
					}
				}
			}
		}
	}
	schemas := make(map[string]bool)
	for i := range late {
		schemas[ops[i].SchemaName] = true
	}
	for i, op := range ops {
		if op.Phase == PhasePre && op.Type == OpDropSchema && schemas[op.ObjectName] {
			late[i] = true
		}
	}
	for i := range late {
		ops[i].Phase = PhaseLateDrop
	}
	return ops
}

// typeUsedByKeptColumn reports whether a column of a table or an attribute of
// a composite type in from has the type, or an array of it, and to still has
// that column. The relations in dropped keep none.
func typeUsedByKeptColumn(from, to *Catalog, dropped map[uint32]bool, typ uint32) bool {
	if typ == 0 {
		return false
	}
	for _, rel := range from.relationByOID {
		if rel.RelKind != 'r' && rel.RelKind != 'p' && rel.RelKind != 'c' || rel.Schema == nil || dropped[rel.OID] {
			continue
		}
		toRel := to.GetRelation(rel.Schema.Name, rel.Name)
		if toRel == nil || toRel.RelKind != rel.RelKind {
			continue
		}
		for _, col := range rel.Columns {
			used := col.TypeOID == typ
			if t := from.typeByOID[col.TypeOID]; t != nil && t.Elem == typ && t.Category == 'A' {
				used = true
			}
			if _, kept := toRel.colByName[col.Name]; used && kept {
				return true
			}
		}
	}
	return false
}

// sortMigrationOps separates ops into phases and sorts each phase.
// PhasePre: reverse topological sort using from catalog deps (drop dependents first).
// PhaseMain: topologically sorted using to catalog deps (forward ordering).
// PhasePost: sorted by name for determinism.
//
// If the PhaseMain sort detects a cycle, it may reclassify some ops (e.g.,
// CHECK constraints) to PhasePost to break the cycle. Those ops are then
// collected into postOps for the final assembly.
func sortMigrationOps(from, to *Catalog, ops []MigrationOp) []MigrationOp {
	var preOps, mainOps, postOps, lateDropOps []MigrationOp
	for _, op := range ops {
		switch op.Phase {
		case PhasePre:
			preOps = append(preOps, op)
		case PhaseMain:
			mainOps = append(mainOps, op)
		case PhasePost:
			postOps = append(postOps, op)
		case PhaseLateDrop:
			lateDropOps = append(lateDropOps, op)
		}
	}

	// PhasePre: reverse topological sort using from catalog deps.
	// Dependents are dropped before the objects they depend on.
	preOps, _ = topoSortOps(from, preOps, true)

	// PhaseMain: topological sort using to catalog deps.
	// topoSortOps returns deferred ops (cycle-breaking) separately.
	mainSorted, mainDeferred := topoSortOps(to, mainOps, false)
	postOps = append(postOps, mainDeferred...)

	// PhasePost: sort by name for determinism.
	sort.SliceStable(postOps, func(i, j int) bool {
		return postOps[i].ObjectName < postOps[j].ObjectName
	})

	// PhaseLateDrop: the same reverse order as PhasePre.
	lateDropOps, _ = topoSortOps(from, lateDropOps, true)

	result := make([]MigrationOp, 0, len(ops))
	result = append(result, preOps...)
	result = append(result, mainSorted...)
	result = append(result, postOps...)
	result = append(result, lateDropOps...)
	return result
}
