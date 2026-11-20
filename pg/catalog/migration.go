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
	PriorityRule         = 10
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
	OpCreateCollation    MigrationOpType = "CreateCollation"
	OpDropCollation      MigrationOpType = "DropCollation"
	OpAlterCollation     MigrationOpType = "AlterCollation"
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
	OpCreateRule         MigrationOpType = "CreateRule"
	OpDropRule           MigrationOpType = "DropRule"
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
	ops = append(ops, generateCollationDDL(from, to, diff)...)
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
	ops = append(ops, generateRuleDDL(from, to, diff)...)
	ops = append(ops, generateEventTriggerDDL(from, to, diff)...)
	ops = append(ops, generatePolicyDDL(from, to, diff)...)
	ops = append(ops, generateCommentDDL(from, to, diff)...)
	ops = append(ops, generateGrantDDL(from, to, diff)...)

	// Post-processing: detach/rebuild dependents of modified collations before
	// DROP COLLATION in PhasePre, and wrap ALTER COLUMN TYPE ops with DROP VIEW /
	// CREATE VIEW for dependent views.
	var extraViewTables, directViews map[uint32]bool
	ops, extraViewTables, directViews = wrapModifiedCollationDependents(from, to, diff, ops)
	ops = wrapColumnTypeChangesWithViewOps(from, to, diff, ops, extraViewTables, directViews)
	ops = reconcileRulesWithDroppedRelations(from, to, ops)
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
func wrapColumnTypeChangesWithViewOps(from, to *Catalog, diff *SchemaDiff, ops []MigrationOp, extraTableOIDs, directViewOIDs map[uint32]bool) []MigrationOp {
	// Find OIDs of tables that have column type changes or column drops.
	tableOIDs := make(map[uint32]bool)
	for oid := range extraTableOIDs {
		tableOIDs[oid] = true
	}
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
			if rel.From.RelKind == 'v' && rel.To.RelKind == 'v' &&
				(viewColumnsChanged(rel) || normalizeRelPersistence(rel.From.Persistence) != normalizeRelPersistence(rel.To.Persistence)) {
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
				if oldType != newType || columnCollationKey(from, col.From) != columnCollationKey(to, col.To) {
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
	if len(tableOIDs) == 0 && len(directViewOIDs) == 0 {
		return ops
	}

	// Use from.deps to find views/matviews that depend on these tables (transitively).
	type viewInfo struct {
		schema      string
		name        string
		oid         uint32
		relKind     byte
		persistence byte
	}
	seen := make(map[uint32]bool)
	var viewsToDrop []viewInfo

	queue := make([]uint32, 0, len(tableOIDs)+len(directViewOIDs))
	for oid := range tableOIDs {
		queue = append(queue, oid)
	}
	for oid := range directViewOIDs {
		rel := from.GetRelationByOID(oid)
		if rel == nil || (rel.RelKind != 'v' && rel.RelKind != 'm') || rel.Schema == nil || seen[oid] {
			continue
		}
		seen[oid] = true
		viewsToDrop = append(viewsToDrop, viewInfo{schema: rel.Schema.Name, name: rel.Name, oid: rel.OID, relKind: rel.RelKind, persistence: rel.Persistence})
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
			if rel == nil || (rel.RelKind != 'v' && rel.RelKind != 'm') {
				continue
			}
			seen[d.ObjOID] = true
			if rel.Schema == nil {
				continue
			}
			viewsToDrop = append(viewsToDrop, viewInfo{schema: rel.Schema.Name, name: rel.Name, oid: rel.OID, relKind: rel.RelKind, persistence: rel.Persistence})
			queue = append(queue, d.ObjOID)
		}
	}

	preDroppedRels := make(map[uint32]bool)
	preDroppedRules := make(map[uint32]bool)
	for _, op := range ops {
		if op.Phase == PhasePre {
			switch op.Type {
			case OpDropTable, OpDropView:
				if op.ObjOID != 0 {
					preDroppedRels[op.ObjOID] = true
				}
			case OpDropRule:
				if op.ObjOID != 0 {
					preDroppedRules[op.ObjOID] = true
				}
			}
		}
	}

	type ruleInfo struct {
		schema string
		rel    string
		rule   *Rule
	}
	seenRules := make(map[uint32]bool)
	var rulesToDrop []ruleInfo
	for _, d := range from.deps {
		if d.RefType != 'r' || d.ObjType != 'R' || (!tableOIDs[d.RefOID] && !seen[d.RefOID]) {
			continue
		}
		if seenRules[d.ObjOID] || preDroppedRules[d.ObjOID] {
			continue
		}
		rule := from.rules[d.ObjOID]
		if rule == nil || seen[rule.OwnerRelOID] || preDroppedRels[rule.OwnerRelOID] {
			continue
		}
		ownerRel := from.GetRelationByOID(rule.OwnerRelOID)
		if ownerRel == nil || ownerRel.Schema == nil {
			continue
		}
		seenRules[d.ObjOID] = true
		rulesToDrop = append(rulesToDrop, ruleInfo{schema: ownerRel.Schema.Name, rel: ownerRel.Name, rule: rule})
	}

	if len(viewsToDrop) == 0 && len(rulesToDrop) == 0 {
		return ops
	}

	// Sort views and rules for determinism.
	sort.Slice(viewsToDrop, func(i, j int) bool {
		if viewsToDrop[i].schema != viewsToDrop[j].schema {
			return viewsToDrop[i].schema < viewsToDrop[j].schema
		}
		return viewsToDrop[i].name < viewsToDrop[j].name
	})
	sort.Slice(rulesToDrop, func(i, j int) bool {
		if rulesToDrop[i].schema != rulesToDrop[j].schema {
			return rulesToDrop[i].schema < rulesToDrop[j].schema
		}
		if rulesToDrop[i].rel != rulesToDrop[j].rel {
			return rulesToDrop[i].rel < rulesToDrop[j].rel
		}
		return rulesToDrop[i].rule.Name < rulesToDrop[j].rule.Name
	})

	// Build drop and create ops. Resolve OIDs from `to` catalog for CREATE.
	var extraOps []MigrationOp
	dropViews := make(map[string]bool) // "schema.name" for dedup
	for _, v := range viewsToDrop {
		dropViews[v.schema+"."+v.name] = true
		qn := migrationQualifiedName(v.schema, v.name)
		if v.persistence == 't' && !isTempSchemaName(v.schema) {
			qn = quoteIdentAlways(v.name)
		}
		dropKind := "VIEW"
		if v.relKind == 'm' {
			dropKind = "MATERIALIZED VIEW"
		}

		// Use from-catalog OID for DROP (PhasePre uses from deps).
		extraOps = append(extraOps, MigrationOp{
			Type:          OpDropView,
			SchemaName:    v.schema,
			ObjectName:    v.name,
			SQL:           fmt.Sprintf("DROP %s IF EXISTS %s", dropKind, qn),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       'r',
			ObjOID:        v.oid,
			Priority:      PriorityView,
		})

		toRel := to.GetRelation(v.schema, v.name)
		if toRel != nil && (toRel.RelKind == 'v' || toRel.RelKind == 'm') {
			vEntry := RelationDiffEntry{SchemaName: v.schema, Name: v.name, To: toRel}
			if toRel.RelKind == 'm' {
				extraOps = append(extraOps, buildCreateMatViewOp(to, vEntry))
			} else {
				extraOps = append(extraOps, buildCreateViewOp(to, vEntry))
			}
		}
	}

	dropRules := make(map[string]bool) // "schema.rel.rule" for dedup
	for _, ri := range rulesToDrop {
		dropRules[ri.schema+"."+ri.rel+"."+ri.rule.Name] = true
		extraOps = append(extraOps, buildDropRuleOp(ri.schema, ri.rel, ri.rule))
		if toRel := to.GetRelation(ri.schema, ri.rel); toRel != nil && toRel.Rules != nil {
			if toRule := toRel.Rules[ri.rule.Name]; toRule != nil {
				extraOps = append(extraOps, buildCreateAndModeRuleOps(to, ri.schema, ri.rel, toRel, toRule)...)
			}
		}
	}

	// Filter out any existing view or rule ops for these objects (avoid duplicates).
	var filteredOps []MigrationOp
	for _, op := range ops {
		skip := false
		if op.Type == OpDropView || op.Type == OpCreateView || op.Type == OpAlterView {
			if dropViews[op.SchemaName+"."+op.ObjectName] {
				skip = true
			}
		}
		if (op.Type == OpDropRule || op.Type == OpCreateRule || (op.Type == OpAlterTable && op.ObjType == 'R')) &&
			dropRules[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] {
			skip = true
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
	// Lift type → owning base/range/relation (array, multirange, or row type).
	if objType == 't' {
		targetOID := objOID
		if bt := c.typeByOID[targetOID]; bt != nil && bt.Category == 'A' && bt.Elem != 0 {
			if idxs, ok := oidToIdx[depKey{'t', bt.Elem}]; ok {
				return idxs
			}
			targetOID = bt.Elem
		}
		for _, rt := range c.rangeTypes {
			if rt != nil && (rt.ArrayOID == objOID || rt.MultirangeOID == targetOID) {
				if idxs, ok := oidToIdx[depKey{'t', rt.OID}]; ok {
					return idxs
				}
			}
		}
		for _, rel := range c.relationByOID {
			if rel.RowTypeOID == targetOID {
				return oidToIdx[depKey{'r', rel.OID}]
			}
		}
	}
	// Lift trigger ('g' in deps) → trigger op ('T' in MigrationOp) or owning table.
	if objType == 'g' {
		if idxs, ok := oidToIdx[depKey{'T', objOID}]; ok {
			return idxs
		}
		if trig, ok := c.triggers[objOID]; ok {
			return oidToIdx[depKey{'r', trig.RelOID}]
		}
	}
	// Lift rule ('R' in deps) → owning relation if no dedicated rule op exists.
	if objType == 'R' {
		if rule, ok := c.rules[objOID]; ok {
			return oidToIdx[depKey{'r', rule.OwnerRelOID}]
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
	createdColl := make(map[string]bool)
	rebuilt := make(map[uint32]bool)
	for _, op := range ops {
		switch op.Type {
		case OpCreateType:
			created[op.SchemaName+"."+op.ObjectName] = true
		case OpCreateCollation:
			createdColl[op.SchemaName+"."+op.ObjectName] = true
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
		if op.Phase == PhasePre && op.Type == OpDropCollation && !createdColl[op.SchemaName+"."+op.ObjectName] && collUsedInPhaseMain(from, to, ops, rebuilt, op.ObjOID) {
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
					if !late[j] && (ops[j].Type == OpDropType || ops[j].Type == OpDropCollation) {
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
// a composite type in from has the type, an array of it, or a domain over it,
// and to still has that column. The relations in dropped keep none.
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
			if _, kept := toRel.colByName[col.Name]; kept && typeMatchesOrWraps(from, col.TypeOID, typ) {
				return true
			}
		}
	}
	return false
}

func typeMatchesOrWraps(from *Catalog, colTypeOID, targetOID uint32) bool {
	if colTypeOID == 0 || targetOID == 0 {
		return false
	}
	for cur := colTypeOID; cur != 0; {
		if cur == targetOID {
			return true
		}
		if t := from.typeByOID[cur]; t != nil && t.Category == 'A' && t.Elem != 0 {
			if t.Elem == targetOID {
				return true
			}
			cur = t.Elem
			continue
		}
		if dt := from.domainTypes[cur]; dt != nil && dt.BaseTypeOID != cur {
			cur = dt.BaseTypeOID
			continue
		}
		break
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

func generateCollationDDL(from, to *Catalog, diff *SchemaDiff) []MigrationOp {
	var ops []MigrationOp
	for _, entry := range diff.Collations {
		qn := migrationQualifiedName(entry.SchemaName, entry.Name)
		switch entry.Action {
		case DiffAdd:
			ops = append(ops, MigrationOp{
				Type:          OpCreateCollation,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("CREATE COLLATION %s %s", qn, entry.To.DefClause),
				Transactional: true,
				Phase:         PhaseMain,
				ObjType:       'L',
				ObjOID:        entry.To.OID,
				Priority:      PriorityType,
			})
			if entry.To.OwnerName != "" && entry.To.OwnerName != "postgres" {
				ops = append(ops, MigrationOp{
					Type:          OpAlterCollation,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER COLLATION %s OWNER TO %s", qn, quoteIdentifier(entry.To.OwnerName)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'L',
					ObjOID:        entry.To.OID,
					Priority:      PriorityType + 1,
				})
			}
		case DiffDrop:
			ops = append(ops, MigrationOp{
				Type:          OpDropCollation,
				SchemaName:    entry.SchemaName,
				ObjectName:    entry.Name,
				SQL:           fmt.Sprintf("DROP COLLATION %s", qn),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'L',
				ObjOID:        entry.From.OID,
				Priority:      PriorityType,
			})
		case DiffModify:
			if entry.From.DefClause != entry.To.DefClause {
				ops = append(ops,
					MigrationOp{
						Type:          OpDropCollation,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("DROP COLLATION %s", qn),
						Transactional: true,
						Phase:         PhasePre,
						ObjType:       'L',
						ObjOID:        entry.From.OID,
						Priority:      PriorityType,
					},
					MigrationOp{
						Type:          OpCreateCollation,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("CREATE COLLATION %s %s", qn, entry.To.DefClause),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'L',
						ObjOID:        entry.To.OID,
						Priority:      PriorityType,
					},
				)
				if entry.To.OwnerName != "" && (entry.To.OwnerName != "postgres" || entry.From.OwnerName != entry.To.OwnerName) {
					ops = append(ops, MigrationOp{
						Type:          OpAlterCollation,
						SchemaName:    entry.SchemaName,
						ObjectName:    entry.Name,
						SQL:           fmt.Sprintf("ALTER COLLATION %s OWNER TO %s", qn, quoteIdentifier(entry.To.OwnerName)),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'L',
						ObjOID:        entry.To.OID,
						Priority:      PriorityType + 1,
					})
				}
			} else if entry.From.OwnerName != entry.To.OwnerName && entry.To.OwnerName != "" {
				ops = append(ops, MigrationOp{
					Type:          OpAlterCollation,
					SchemaName:    entry.SchemaName,
					ObjectName:    entry.Name,
					SQL:           fmt.Sprintf("ALTER COLLATION %s OWNER TO %s", qn, quoteIdentifier(entry.To.OwnerName)),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'L',
					ObjOID:        entry.To.OID,
					Priority:      PriorityType,
				})
			}
		}
	}
	return ops
}

func collUsedInPhaseMain(from, to *Catalog, ops []MigrationOp, dropped map[uint32]bool, collOID uint32) bool {
	if collOID == 0 {
		return false
	}
	if collUsedByKeptColumn(from, to, dropped, collOID) {
		return true
	}
	preDroppedType := make(map[uint32]bool)
	preDroppedCon := make(map[uint32]bool)
	preDroppedConSQL := make(map[string]bool)
	preDroppedIdx := make(map[uint32]bool)
	preDroppedTrig := make(map[uint32]bool)
	preDroppedPol := make(map[uint32]bool)
	preDroppedFunc := make(map[uint32]bool)
	preDroppedCol := make(map[string]bool)
	for _, op := range ops {
		if op.Phase != PhasePre {
			continue
		}
		switch op.Type {
		case OpDropType:
			if op.ObjType == 't' {
				preDroppedType[op.ObjOID] = true
			}
		case OpDropConstraint:
			if op.ObjOID != 0 {
				preDroppedCon[op.ObjOID] = true
			}
			preDroppedConSQL[op.SQL] = true
			if op.ParentObject != "" {
				preDroppedConSQL[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
			}
		case OpDropIndex:
			preDroppedIdx[op.ObjOID] = true
		case OpDropTrigger:
			preDroppedTrig[op.ObjOID] = true
		case OpDropPolicy:
			preDroppedPol[op.ObjOID] = true
		case OpDropFunction:
			preDroppedFunc[op.ObjOID] = true
		case OpDropColumn:
			preDroppedCol[op.SQL] = true
		}
	}

	for _, d := range from.deps {
		if d.DepType == DepInternal || d.RefType != 'L' || d.RefOID != collOID {
			continue
		}
		switch d.ObjType {
		case 'r':
			if dropped[d.ObjOID] {
				continue
			}
			rel := from.GetRelationByOID(d.ObjOID)
			if rel == nil || rel.Schema == nil {
				continue
			}
			toRel := to.GetRelation(rel.Schema.Name, rel.Name)
			if toRel == nil || toRel.RelKind != rel.RelKind {
				continue
			}
			if d.ObjSubID > 0 {
				for _, col := range rel.Columns {
					if int32(col.AttNum) == d.ObjSubID {
						dropColSQL := fmt.Sprintf("ALTER TABLE %s.%s DROP COLUMN %s", quoteIdentifier(rel.Schema.Name), quoteIdentifier(rel.Name), quoteIdentifier(col.Name))
						if !preDroppedCol[dropColSQL] {
							return true
						}
						break
					}
				}
			} else if rel.RelKind == 'v' {
				return true
			}
		case 'f':
			if !preDroppedFunc[d.ObjOID] {
				return true
			}
		case 'p':
			if pol := from.policies[d.ObjOID]; pol != nil && !dropped[pol.RelOID] && !preDroppedPol[pol.OID] {
				return true
			}
		case 'c':
			if preDroppedCon[d.ObjOID] {
				continue
			}
			if con := from.constraints[d.ObjOID]; con != nil {
				if dropped[con.RelOID] {
					continue
				}
				if rel := from.GetRelationByOID(con.RelOID); rel != nil && rel.Schema != nil {
					key := rel.Schema.Name + "." + rel.Name + "." + con.Name
					if !preDroppedConSQL[key] {
						return true
					}
				}
			} else {
				for domOID, dt := range from.domainTypes {
					if preDroppedType[domOID] || dt == nil {
						continue
					}
					for _, dc := range dt.Constraints {
						if dc.OID == d.ObjOID {
							return true
						}
					}
				}
			}
		case 'i':
			if idx := from.indexes[d.ObjOID]; idx != nil && !dropped[idx.RelOID] && !preDroppedIdx[idx.OID] {
				if idx.ConstraintOID != 0 {
					if con := from.constraints[idx.ConstraintOID]; con != nil {
						if rel := from.GetRelationByOID(con.RelOID); rel != nil && rel.Schema != nil {
							key := rel.Schema.Name + "." + rel.Name + "." + con.Name
							if preDroppedCon[con.OID] || preDroppedConSQL[key] {
								continue
							}
						}
					}
				}
				return true
			}
		case 'g':
			if trig := from.triggers[d.ObjOID]; trig != nil && !dropped[trig.RelOID] && !preDroppedTrig[trig.OID] {
				return true
			}
		}
	}
	return false
}

func collUsedByKeptColumn(from, to *Catalog, dropped map[uint32]bool, collOID uint32) bool {
	if collOID == 0 {
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
			if _, kept := toRel.colByName[col.Name]; kept && col.UserCollationOID == collOID {
				return true
			}
		}
	}
	return false
}

func attnumSliceContains(xs []int16, target int16) bool {
	for _, x := range xs {
		if x == target {
			return true
		}
	}
	return false
}

func objDependsOnColumn(c *Catalog, objType byte, objOID, relOID uint32, attnum int16) bool {
	for _, d := range c.deps {
		if d.ObjType == objType && d.ObjOID == objOID && d.RefType == 'r' && d.RefOID == relOID && d.RefSubID == int32(attnum) {
			return true
		}
	}
	return false
}

func fkDependsOnConstraintOrIndex(c *Catalog, fk *Constraint, targetRelOID uint32, targetCon *Constraint, targetIdxOID uint32) bool {
	if fk == nil || fk.Type != ConstraintFK || fk.FRelOID != targetRelOID {
		return false
	}
	for _, d := range c.deps {
		if d.ObjType == 'c' && d.ObjOID == fk.OID {
			if targetCon != nil && ((d.RefType == 'c' && d.RefOID == targetCon.OID) || (targetCon.IndexOID != 0 && d.RefType == 'i' && d.RefOID == targetCon.IndexOID)) {
				return true
			}
			if targetIdxOID != 0 && d.RefType == 'i' && d.RefOID == targetIdxOID {
				return true
			}
		}
	}
	if targetCon != nil {
		for _, fc := range fk.FColumns {
			if attnumSliceContains(targetCon.Columns, fc) {
				return true
			}
		}
	}
	return false
}

// wrapModifiedCollationDependents detaches or rebuilds retained dependents of
// collations whose definition changed (DiffModify), as well as domains and
// ranges whose collation changed, so PhasePre drops succeed cleanly and
// PhaseMain restores the dependents against the new definition.
func wrapModifiedCollationDependents(from, to *Catalog, diff *SchemaDiff, ops []MigrationOp) ([]MigrationOp, map[uint32]bool, map[uint32]bool) {
	modifiedCollOIDs := make(map[uint32]bool)
	for _, entry := range diff.Collations {
		if entry.Action == DiffModify && entry.From != nil && entry.To != nil && entry.From.DefClause != entry.To.DefClause {
			modifiedCollOIDs[entry.From.OID] = true
		}
	}
	collChangedDomainOIDs := make(map[uint32]bool)
	for _, entry := range diff.Domains {
		if entry.Action == DiffModify && entry.From != nil && entry.To != nil {
			fromType := from.FormatType(entry.From.BaseTypeOID, entry.From.BaseTypMod)
			toType := to.FormatType(entry.To.BaseTypeOID, entry.To.BaseTypMod)
			if fromType == toType && domainCollationKey(entry.From, from) != domainCollationKey(entry.To, to) {
				collChangedDomainOIDs[entry.From.TypeOID] = true
			}
		}
	}
	rebuiltRangeOIDs := make(map[uint32]bool)
	for _, entry := range diff.Ranges {
		if entry.Action == DiffModify && entry.From != nil && entry.To != nil {
			fromSub := from.FormatType(entry.From.SubTypeOID, -1)
			toSub := to.FormatType(entry.To.SubTypeOID, -1)
			if fromSub == toSub && !typeUsedByKeptColumn(from, to, nil, entry.From.OID) && !typeUsedByKeptColumn(from, to, nil, entry.From.MultirangeOID) {
				rebuiltRangeOIDs[entry.From.OID] = true
			}
		}
	}
	for oid, rt := range from.rangeTypes {
		if rt != nil && modifiedCollOIDs[rt.UserCollationOID] {
			if sc := from.schemas[rt.Namespace]; sc != nil && resolveTypeOIDByName(to, sc.Name, rt.Name) != 0 {
				if !typeUsedByKeptColumn(from, to, nil, oid) && !typeUsedByKeptColumn(from, to, nil, rt.MultirangeOID) {
					rebuiltRangeOIDs[oid] = true
				}
			}
		}
	}
	if len(modifiedCollOIDs) == 0 && len(collChangedDomainOIDs) == 0 && len(rebuiltRangeOIDs) == 0 {
		return ops, nil, nil
	}

	preDroppedRel := make(map[uint32]bool)
	preDroppedIdx := make(map[uint32]bool)
	preDroppedCon := make(map[uint32]bool)
	preDroppedConKey := make(map[string]bool)
	preDroppedType := make(map[uint32]bool)
	preDroppedTrig := make(map[uint32]bool)
	preDroppedPol := make(map[uint32]bool)
	preDroppedRule := make(map[uint32]bool)
	preDroppedFunc := make(map[uint32]bool)
	mainCreatedIdx := make(map[string]bool)
	mainAddedCon := make(map[string]bool)
	mainCreatedTrig := make(map[string]bool)
	mainCreatedPol := make(map[string]bool)
	mainCreatedRule := make(map[string]bool)
	mainCreatedFunc := make(map[string]bool)
	mainComments := make(map[string]bool)

	for _, op := range ops {
		if op.Phase == PhasePre {
			switch op.Type {
			case OpDropTable, OpDropView:
				preDroppedRel[op.ObjOID] = true
			case OpDropType:
				if op.ObjType == 'r' {
					preDroppedRel[op.ObjOID] = true
				} else if op.ObjType == 't' && !collChangedDomainOIDs[op.ObjOID] && !rebuiltRangeOIDs[op.ObjOID] {
					preDroppedType[op.ObjOID] = true
				}
			case OpDropIndex:
				preDroppedIdx[op.ObjOID] = true
			case OpDropConstraint:
				if op.ObjOID != 0 {
					preDroppedCon[op.ObjOID] = true
				}
				if op.ParentObject != "" {
					preDroppedConKey[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
				}
				preDroppedConKey[op.SQL] = true
			case OpDropTrigger:
				preDroppedTrig[op.ObjOID] = true
			case OpDropPolicy:
				preDroppedPol[op.ObjOID] = true
			case OpDropRule:
				preDroppedRule[op.ObjOID] = true
			case OpDropFunction:
				preDroppedFunc[op.ObjOID] = true
			}
		} else {
			switch op.Type {
			case OpCreateIndex:
				mainCreatedIdx[op.SchemaName+"."+op.ObjectName] = true
			case OpAddConstraint:
				if op.ParentObject != "" {
					mainAddedCon[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
				}
			case OpCreateTrigger:
				mainCreatedTrig[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
			case OpCreatePolicy:
				if op.ParentObject != "" {
					mainCreatedPol[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
				}
				mainCreatedPol[op.SchemaName+"."+op.ObjectName] = true
			case OpCreateRule:
				mainCreatedRule[op.SchemaName+"."+op.ParentObject+"."+op.ObjectName] = true
			case OpCreateFunction, OpAlterFunction:
				mainCreatedFunc[op.SchemaName+"."+op.ObjectName] = true
			case OpComment:
				mainComments[op.SQL] = true
			}
		}
	}

	extraViewTables := make(map[uint32]bool)
	directViews := make(map[uint32]bool)
	var preExtra, mainExtra []MigrationOp

	// 1. Rebuild kept domains and ranges whose own COLLATE changed or references a modified collation.
	rebuiltDomains := make(map[uint32]*DomainType) // from domain TypeOID -> from DomainType
	rebuiltDomainNames := make(map[string]bool)
	rebuiltTypeRefOIDs := make(map[uint32]bool)
	rebuiltRangeTypeOIDs := make(map[uint32]bool)
	for oid := range rebuiltRangeOIDs {
		rebuiltRangeTypeOIDs[oid] = true
		if rt := from.rangeTypes[oid]; rt != nil {
			if rt.ArrayOID != 0 {
				rebuiltRangeTypeOIDs[rt.ArrayOID] = true
			}
			if rt.MultirangeOID != 0 {
				rebuiltRangeTypeOIDs[rt.MultirangeOID] = true
				if mrBT := from.typeByOID[rt.MultirangeOID]; mrBT != nil && mrBT.Array != 0 {
					rebuiltRangeTypeOIDs[mrBT.Array] = true
				}
			}
		}
	}
	for oid, dt := range from.domainTypes {
		if dt == nil || preDroppedType[oid] {
			continue
		}
		if modifiedCollOIDs[dt.UserCollationOID] || collChangedDomainOIDs[oid] {
			bt := from.typeByOID[oid]
			if bt == nil {
				continue
			}
			if sc := from.schemas[bt.Namespace]; sc != nil && resolveTypeOIDByName(to, sc.Name, bt.TypeName) != 0 {
				rebuiltDomains[oid] = dt
			}
		}
	}
	for changed := true; changed; {
		changed = false
		for oid, dt := range from.domainTypes {
			if dt == nil || preDroppedType[oid] || rebuiltDomains[oid] != nil {
				continue
			}
			baseRebuilt := rebuiltDomains[dt.BaseTypeOID] != nil || rebuiltRangeTypeOIDs[dt.BaseTypeOID]
			if !baseRebuilt {
				if baseBT := from.typeByOID[dt.BaseTypeOID]; baseBT != nil && baseBT.Category == 'A' && (rebuiltDomains[baseBT.Elem] != nil || rebuiltRangeTypeOIDs[baseBT.Elem]) {
					baseRebuilt = true
				}
			}
			if baseRebuilt {
				bt := from.typeByOID[oid]
				if bt == nil {
					continue
				}
				if sc := from.schemas[bt.Namespace]; sc != nil && resolveTypeOIDByName(to, sc.Name, bt.TypeName) != 0 {
					rebuiltDomains[oid] = dt
					changed = true
				}
			}
		}
	}
	var domOIDs []uint32
	for oid := range rebuiltDomains {
		domOIDs = append(domOIDs, oid)
	}
	sort.Slice(domOIDs, func(i, j int) bool { return domOIDs[i] < domOIDs[j] })
	for _, oid := range domOIDs {
		bt := from.typeByOID[oid]
		if bt == nil {
			continue
		}
		sc := from.schemas[bt.Namespace]
		if sc == nil {
			continue
		}
		qn := migrationQualifiedName(sc.Name, bt.TypeName)
		key := sc.Name + "." + bt.TypeName
		rebuiltDomainNames[key] = true
		rebuiltTypeRefOIDs[oid] = true
		if bt.Array != 0 {
			rebuiltTypeRefOIDs[bt.Array] = true
		}
		preDroppedType[oid] = true
		preExtra = append(preExtra, MigrationOp{
			Type:          OpDropType,
			SchemaName:    sc.Name,
			ObjectName:    bt.TypeName,
			SQL:           fmt.Sprintf("DROP DOMAIN %s", qn),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       't',
			ObjOID:        oid,
			Priority:      PriorityType,
		})
		if toOID := resolveTypeOIDByName(to, sc.Name, bt.TypeName); toOID != 0 {
			if toDT := to.domainTypes[toOID]; toDT != nil {
				mainExtra = append(mainExtra, MigrationOp{
					Type:          OpCreateType,
					SchemaName:    sc.Name,
					ObjectName:    bt.TypeName,
					SQL:           formatCreateDomain(to, sc.Name, bt.TypeName, toDT),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       't',
					ObjOID:        toOID,
					Priority:      PriorityType,
				})
			}
		}
	}

	rebuiltRangeNames := make(map[string]bool)
	var rangeOIDs []uint32
	for oid := range rebuiltRangeOIDs {
		rangeOIDs = append(rangeOIDs, oid)
	}
	sort.Slice(rangeOIDs, func(i, j int) bool { return rangeOIDs[i] < rangeOIDs[j] })
	for _, oid := range rangeOIDs {
		rt := from.rangeTypes[oid]
		if rt == nil {
			continue
		}
		sc := from.schemas[rt.Namespace]
		if sc == nil {
			continue
		}
		qn := migrationQualifiedName(sc.Name, rt.Name)
		key := sc.Name + "." + rt.Name
		rebuiltRangeNames[key] = true
		rebuiltTypeRefOIDs[oid] = true
		if rt.ArrayOID != 0 {
			rebuiltTypeRefOIDs[rt.ArrayOID] = true
		}
		if rt.MultirangeOID != 0 {
			rebuiltTypeRefOIDs[rt.MultirangeOID] = true
			if mrBT := from.typeByOID[rt.MultirangeOID]; mrBT != nil && mrBT.Array != 0 {
				rebuiltTypeRefOIDs[mrBT.Array] = true
			}
		}
		preDroppedType[oid] = true
		preExtra = append(preExtra, MigrationOp{
			Type:          OpDropType,
			SchemaName:    sc.Name,
			ObjectName:    rt.Name,
			SQL:           fmt.Sprintf("DROP TYPE %s", qn),
			Transactional: true,
			Phase:         PhasePre,
			ObjType:       't',
			ObjOID:        oid,
			Priority:      PriorityType,
		})
		if toOID := resolveTypeOIDByName(to, sc.Name, rt.Name); toOID != 0 {
			if toRT := to.rangeTypes[toOID]; toRT != nil {
				mainExtra = append(mainExtra, MigrationOp{
					Type:          OpCreateType,
					SchemaName:    sc.Name,
					ObjectName:    rt.Name,
					SQL:           formatCreateRange(to, sc.Name, rt.Name, toRT),
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       't',
					ObjOID:        toOID,
					Priority:      PriorityType,
				})
			}
		}
	}

	if len(rebuiltDomainNames) > 0 || len(rebuiltRangeNames) > 0 {
		filtered := ops[:0]
		for _, op := range ops {
			if (op.Type == OpAlterType || op.Type == OpCreateType || op.Type == OpDropType) && op.ObjType == 't' {
				key := op.SchemaName + "." + op.ObjectName
				if rebuiltDomainNames[key] || rebuiltRangeNames[key] {
					continue
				}
			}
			filtered = append(filtered, op)
		}
		ops = filtered
	}

	unwrappedDetachType := func(typeOID uint32) (string, bool) {
		isArr := false
		curOID := typeOID
		if bt := from.typeByOID[curOID]; bt != nil && bt.Category == 'A' && rebuiltDomains[bt.Elem] != nil {
			isArr = true
			curOID = bt.Elem
		}
		dt := rebuiltDomains[curOID]
		if dt == nil {
			return "", false
		}
		baseOID, baseMod := dt.BaseTypeOID, dt.BaseTypMod
		for rebuiltDomains[baseOID] != nil {
			inner := rebuiltDomains[baseOID]
			baseOID, baseMod = inner.BaseTypeOID, inner.BaseTypMod
		}
		if baseBT := from.typeByOID[baseOID]; baseBT != nil && baseBT.Category == 'A' && rebuiltDomains[baseBT.Elem] != nil {
			isArr = true
			inner := rebuiltDomains[baseBT.Elem]
			baseOID, baseMod = inner.BaseTypeOID, inner.BaseTypMod
			for rebuiltDomains[baseOID] != nil {
				inner = rebuiltDomains[baseOID]
				baseOID, baseMod = inner.BaseTypeOID, inner.BaseTypMod
			}
		}
		res := from.FormatType(baseOID, baseMod)
		if isArr {
			res += "[]"
		}
		return res, true
	}

	dropAndRestoreConstraint := func(con *Constraint) {
		if con == nil || con.ConParentID != 0 || con.Type == ConstraintTrigger || preDroppedRel[con.RelOID] {
			return
		}
		cRel := from.GetRelationByOID(con.RelOID)
		if cRel == nil || cRel.Schema == nil {
			return
		}
		cSc := cRel.Schema.Name
		conKey := cSc + "." + cRel.Name + "." + con.Name
		dropOp := buildDropConstraintOp(cSc, cRel.Name, con)
		if !preDroppedCon[con.OID] && !preDroppedConKey[conKey] && !preDroppedConKey[dropOp.SQL] {
			preDroppedCon[con.OID] = true
			preDroppedConKey[conKey] = true
			preDroppedConKey[dropOp.SQL] = true
			preExtra = append(preExtra, dropOp)
		}
		if !mainAddedCon[conKey] {
			if toCRel := to.GetRelation(cSc, cRel.Name); toCRel != nil {
				for _, toCon := range to.ConstraintsOf(toCRel.OID) {
					if toCon.Name == con.Name && toCon.ConParentID == 0 && toCon.Type != ConstraintTrigger {
						mainAddedCon[conKey] = true
						mainExtra = append(mainExtra, buildAddConstraintOp(to, cSc, cRel.Name, toCon))
						if txt, ok := to.GetComment('c', toCon.OID, 0); ok && txt != "" {
							desc := cSc + "." + cRel.Name + "." + toCon.Name
							if csql := formatCommentSQL(to, 'c', desc, 0, txt); csql != "" && !mainComments[csql] {
								mainComments[csql] = true
								mainExtra = append(mainExtra, MigrationOp{
									Type:          OpComment,
									ObjectName:    desc,
									SQL:           csql,
									Transactional: true,
									Phase:         PhaseMain,
									ObjType:       'c',
									Priority:      PriorityMetadata,
								})
							}
						}
						break
					}
				}
			}
		}
	}

	rebuildGenColDependents := func(fromRel *Relation, fromCol *Column, toRel *Relation, toCol *Column) {
		scName := fromRel.Schema.Name

		// Standalone indexes referencing the generated column (key, INCLUDE, expr, or WHERE).
		for _, idx := range from.IndexesOf(fromRel.OID) {
			if idx.ConstraintOID != 0 || idx.ParentIndexOID != 0 || idx.Schema == nil {
				continue
			}
			if !attnumSliceContains(idx.Columns, fromCol.AttNum) && !objDependsOnColumn(from, 'i', idx.OID, fromRel.OID, fromCol.AttNum) {
				continue
			}
			if idx.IsUnique {
				for _, otherCon := range from.constraints {
					if fkDependsOnConstraintOrIndex(from, otherCon, fromRel.OID, nil, idx.OID) {
						dropAndRestoreConstraint(otherCon)
					}
				}
			}
			idxSc := idx.Schema.Name
			if !preDroppedIdx[idx.OID] {
				preDroppedIdx[idx.OID] = true
				preExtra = append(preExtra, MigrationOp{
					Type:          OpDropIndex,
					SchemaName:    idxSc,
					ObjectName:    idx.Name,
					SQL:           fmt.Sprintf("DROP INDEX %s.%s", quoteIdentifier(idxSc), quoteIdentifier(idx.Name)),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'i',
					ObjOID:        idx.OID,
					Priority:      PriorityIndex,
				})
			}
			idxKey := idxSc + "." + idx.Name
			if !mainCreatedIdx[idxKey] {
				if toSc := to.GetSchema(idxSc); toSc != nil {
					if toIdx := toSc.Indexes[idx.Name]; toIdx != nil && toIdx.ConstraintOID == 0 && toIdx.ParentIndexOID == 0 {
						if toIdxRel := to.GetRelationByOID(toIdx.RelOID); toIdxRel != nil {
							mainCreatedIdx[idxKey] = true
							mainExtra = append(mainExtra, MigrationOp{
								Type:          OpCreateIndex,
								SchemaName:    idxSc,
								ObjectName:    toIdx.Name,
								SQL:           buildCreateIndexSQL(toIdx, toIdxRel, idxSc),
								Transactional: true,
								Phase:         PhaseMain,
								ObjType:       'i',
								ObjOID:        toIdx.OID,
								Priority:      PriorityIndex,
							})
							mainExtra = append(mainExtra, buildIndexColumnRenameOps(idxSc, toIdx, defaultIndexColumnNames(to, toIdx), effectiveIndexColumnNames(to, toIdx))...)
							if txt, ok := to.GetComment('i', toIdx.OID, 0); ok && txt != "" {
								desc := idxSc + "." + toIdx.Name
								if csql := formatCommentSQL(to, 'i', desc, 0, txt); csql != "" && !mainComments[csql] {
									mainComments[csql] = true
									mainExtra = append(mainExtra, MigrationOp{
										Type:          OpComment,
										ObjectName:    desc,
										SQL:           csql,
										Transactional: true,
										Phase:         PhaseMain,
										ObjType:       'i',
										Priority:      PriorityMetadata,
									})
								}
							}
						}
					}
				}
			}
		}

		// Constraints on fromRel referencing the generated column.
		for _, con := range from.ConstraintsOf(fromRel.OID) {
			if con.ConParentID != 0 || con.Type == ConstraintTrigger {
				continue
			}
			refsCol := attnumSliceContains(con.Columns, fromCol.AttNum) ||
				attnumSliceContains(con.IncludeColumns, fromCol.AttNum) ||
				attnumSliceContains(con.FKDelSetCols, fromCol.AttNum) ||
				(con.Type == ConstraintCheck && containsColumnRef(con.CheckExpr, fromCol.Name)) ||
				objDependsOnColumn(from, 'c', con.OID, fromRel.OID, fromCol.AttNum)
			if !refsCol {
				continue
			}
			if con.Type == ConstraintPK || con.Type == ConstraintUnique {
				for _, otherCon := range from.constraints {
					if fkDependsOnConstraintOrIndex(from, otherCon, fromRel.OID, con, con.IndexOID) {
						dropAndRestoreConstraint(otherCon)
					}
				}
			}
			dropAndRestoreConstraint(con)
		}

		// Inbound foreign keys on any table directly referencing (fromRel.OID, fromCol.AttNum).
		for _, otherCon := range from.constraints {
			if otherCon.Type == ConstraintFK && otherCon.ConParentID == 0 && otherCon.FRelOID == fromRel.OID && attnumSliceContains(otherCon.FColumns, fromCol.AttNum) {
				dropAndRestoreConstraint(otherCon)
			}
		}

		// Triggers on fromRel referencing the generated column.
		for _, trig := range from.TriggersOf(fromRel.OID) {
			if trig.ParentTriggerOID != 0 {
				continue
			}
			refsCol := attnumSliceContains(trig.Columns, fromCol.AttNum) ||
				(trig.WhenExpr != "" && containsColumnRef(trig.WhenExpr, fromCol.Name)) ||
				objDependsOnColumn(from, 'g', trig.OID, fromRel.OID, fromCol.AttNum)
			if !refsCol {
				continue
			}
			if !preDroppedTrig[trig.OID] {
				preDroppedTrig[trig.OID] = true
				preExtra = append(preExtra, buildDropTriggerOp(scName, fromRel.Name, trig.Name, trig.OID))
			}
			trigKey := scName + "." + fromRel.Name + "." + trig.Name
			if !mainCreatedTrig[trigKey] {
				for _, toTrig := range to.TriggersOf(toRel.OID) {
					if toTrig.Name == trig.Name && toTrig.ParentTriggerOID == 0 {
						mainCreatedTrig[trigKey] = true
						mainExtra = append(mainExtra, buildCreateTriggerOp(to, scName, fromRel.Name, toTrig))
						if txt, ok := to.GetComment('g', toTrig.OID, 0); ok && txt != "" {
							desc := scName + "." + fromRel.Name + "." + toTrig.Name
							if csql := formatCommentSQL(to, 'g', desc, 0, txt); csql != "" && !mainComments[csql] {
								mainComments[csql] = true
								mainExtra = append(mainExtra, MigrationOp{
									Type:          OpComment,
									ObjectName:    desc,
									SQL:           csql,
									Transactional: true,
									Phase:         PhaseMain,
									ObjType:       'g',
									Priority:      PriorityMetadata,
								})
							}
						}
						break
					}
				}
			}
		}

		// Policies on fromRel referencing the generated column.
		for _, pol := range from.policiesByRel[fromRel.OID] {
			refsCol := (pol.UsingExpr != "" && containsColumnRef(pol.UsingExpr, fromCol.Name)) ||
				(pol.CheckExpr != "" && containsColumnRef(pol.CheckExpr, fromCol.Name)) ||
				objDependsOnColumn(from, 'p', pol.OID, fromRel.OID, fromCol.AttNum)
			if !refsCol {
				continue
			}
			if !preDroppedPol[pol.OID] {
				preDroppedPol[pol.OID] = true
				preExtra = append(preExtra, buildDropPolicyOp(scName, fromRel.Name, pol.Name, pol.OID))
			}
			polKey := scName + "." + fromRel.Name + "." + pol.Name
			if !mainCreatedPol[polKey] && !mainCreatedPol[scName+"."+pol.Name] {
				for _, toPol := range to.policiesByRel[toRel.OID] {
					if toPol.Name == pol.Name {
						mainCreatedPol[polKey] = true
						mainExtra = append(mainExtra, buildCreatePolicyOp(scName, fromRel.Name, toPol))
						if txt, ok := to.GetComment('p', toPol.OID, 0); ok && txt != "" {
							desc := scName + "." + fromRel.Name + "." + toPol.Name
							if csql := formatCommentSQL(to, 'p', desc, 0, txt); csql != "" && !mainComments[csql] {
								mainComments[csql] = true
								mainExtra = append(mainExtra, MigrationOp{
									Type:          OpComment,
									ObjectName:    desc,
									SQL:           csql,
									Transactional: true,
									Phase:         PhaseMain,
									ObjType:       'p',
									Priority:      PriorityMetadata,
								})
							}
						}
						break
					}
				}
			}
		}

		// Rules referencing the generated column.
		for _, rule := range from.rules {
			if rule == nil || preDroppedRule[rule.OID] || preDroppedRel[rule.OwnerRelOID] {
				continue
			}
			if !objDependsOnColumn(from, 'R', rule.OID, fromRel.OID, fromCol.AttNum) {
				continue
			}
			rOwner := from.GetRelationByOID(rule.OwnerRelOID)
			if rOwner == nil || rOwner.Schema == nil {
				continue
			}
			rSc := rOwner.Schema.Name
			preDroppedRule[rule.OID] = true
			preExtra = append(preExtra, buildDropRuleOp(rSc, rOwner.Name, rule))
			ruleKey := rSc + "." + rOwner.Name + "." + rule.Name
			if !mainCreatedRule[ruleKey] {
				if toOwner := to.GetRelation(rSc, rOwner.Name); toOwner != nil && toOwner.Rules != nil {
					if toRule := toOwner.Rules[rule.Name]; toRule != nil {
						mainCreatedRule[ruleKey] = true
						mainExtra = append(mainExtra, buildCreateAndModeRuleOps(to, rSc, toOwner.Name, toOwner, toRule)...)
					}
				}
			}
		}

		// Unchanged column comment on toCol.
		if txt, ok := to.GetComment('r', toRel.OID, toCol.AttNum); ok && txt != "" {
			desc := scName + "." + fromRel.Name
			if csql := formatCommentSQL(to, 'r', desc, toCol.AttNum, txt); csql != "" && !mainComments[csql] {
				mainComments[csql] = true
				mainExtra = append(mainExtra, MigrationOp{
					Type:          OpComment,
					ObjectName:    desc,
					SQL:           csql,
					Transactional: true,
					Phase:         PhaseMain,
					ObjType:       'r',
					Priority:      PriorityMetadata,
				})
			}
		}
	}

	// 2. Detach or rebuild kept table columns and composite attributes.
	var relList []*Relation
	for _, r := range from.relationByOID {
		if r != nil && r.Schema != nil && (r.RelKind == 'r' || r.RelKind == 'p' || r.RelKind == 'c') && !preDroppedRel[r.OID] {
			relList = append(relList, r)
		}
	}
	sort.Slice(relList, func(i, j int) bool {
		if relList[i].Schema.Name != relList[j].Schema.Name {
			return relList[i].Schema.Name < relList[j].Schema.Name
		}
		return relList[i].Name < relList[j].Name
	})

	skipMainDropDefaultSQL := make(map[string]bool)
	rebuiltGenColPrefix := make(map[string]bool)
	flipCompDropAttrSQL := make(map[string]uint32)

	for _, fromRel := range relList {
		scName := fromRel.Schema.Name
		toRel := to.GetRelation(scName, fromRel.Name)
		if toRel == nil || toRel.RelKind != fromRel.RelKind {
			continue
		}
		tableName := fmt.Sprintf("%s.%s", quoteIdentifier(scName), quoteIdentifier(fromRel.Name))
		qn := migrationQualifiedName(scName, fromRel.Name)

		for _, fromCol := range fromRel.Columns {
			colUsesModColl := modifiedCollOIDs[fromCol.UserCollationOID]
			modCollDepCount := 0
			exprUsesRebuiltType := false
			for _, d := range from.deps {
				if d.ObjType == 'r' && d.ObjOID == fromRel.OID && d.ObjSubID == int32(fromCol.AttNum) {
					if d.RefType == 'L' && modifiedCollOIDs[d.RefOID] {
						modCollDepCount++
					}
					if d.RefType == 't' && rebuiltTypeRefOIDs[d.RefOID] && d.RefOID != fromCol.TypeOID {
						exprUsesRebuiltType = true
					}
				}
			}
			ownCount := 0
			if colUsesModColl {
				ownCount = 1
			}
			exprUsesModColl := modCollDepCount > ownCount || exprUsesRebuiltType
			rebuiltDomBase, colUsesRebuiltDom := unwrappedDetachType(fromCol.TypeOID)

			toIdx, kept := toRel.colByName[fromCol.Name]
			if !kept {
				if fromRel.RelKind == 'c' && (colUsesModColl || colUsesRebuiltDom) {
					sql := fmt.Sprintf("ALTER TYPE %s DROP ATTRIBUTE %s", qn, quoteIdentAlways(fromCol.Name))
					flipCompDropAttrSQL[sql] = fromRel.OID
				}
				continue
			}
			if !colUsesModColl && !exprUsesModColl && !colUsesRebuiltDom {
				continue
			}
			toCol := toRel.Columns[toIdx]

			if fromRel.RelKind == 'c' {
				detachType := from.FormatType(fromCol.TypeOID, fromCol.TypeMod)
				if colUsesRebuiltDom {
					detachType = rebuiltDomBase
				}
				preExtra = append(preExtra, MigrationOp{
					Type:          OpAlterType,
					SchemaName:    scName,
					ObjectName:    fromRel.Name,
					SQL:           fmt.Sprintf("ALTER TYPE %s ALTER ATTRIBUTE %s TYPE %s", qn, quoteIdentAlways(fromCol.Name), detachType),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        fromRel.OID,
					Priority:      PriorityColumn,
				})
				fromTypeStr := from.FormatType(fromCol.TypeOID, fromCol.TypeMod)
				toTypeStr := to.FormatType(toCol.TypeOID, toCol.TypeMod)
				if fromTypeStr == toTypeStr && columnCollationKey(from, fromCol) == columnCollationKey(to, toCol) {
					mainExtra = append(mainExtra, MigrationOp{
						Type:          OpAlterType,
						SchemaName:    scName,
						ObjectName:    fromRel.Name,
						SQL:           fmt.Sprintf("ALTER TYPE %s ALTER ATTRIBUTE %s TYPE %s%s", qn, quoteIdentAlways(toCol.Name), toTypeStr, formatColumnCollate(to, toCol)),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        toRel.OID,
						Priority:      PriorityColumn,
					})
				}
				continue
			}

			// Stored generated column: drop dependents and column in PhasePre, re-add in PhaseMain.
			if fromCol.Generated == 's' {
				extraViewTables[fromRel.OID] = true
				rebuildGenColDependents(fromRel, fromCol, toRel, toCol)
				colRef := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s", tableName, quoteIdentifier(fromCol.Name))
				rebuiltGenColPrefix[colRef] = true
				preExtra = append(preExtra, MigrationOp{
					Type:          OpDropColumn,
					SchemaName:    scName,
					ObjectName:    fromRel.Name,
					SQL:           fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", tableName, quoteIdentifier(fromCol.Name)),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        fromRel.OID,
					Priority:      PriorityColumn,
				})
				addOps := columnAddOps(to, tableName, scName, fromRel.Name, ColumnDiffEntry{Action: DiffAdd, Name: toCol.Name, To: toCol})
				for i := range addOps {
					addOps[i].Phase = PhaseMain
					addOps[i].ObjType = 'r'
					addOps[i].ObjOID = toRel.OID
					addOps[i].Priority = PriorityColumn
				}
				mainExtra = append(mainExtra, addOps...)
				continue
			}

			colRef := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s", tableName, quoteIdentifier(toCol.Name))
			fromTypeStr := from.FormatType(fromCol.TypeOID, fromCol.TypeMod)
			toTypeStr := to.FormatType(toCol.TypeOID, toCol.TypeMod)
			typeChanged := fromTypeStr != toTypeStr
			alterTypeAlreadyInMain := typeChanged || columnCollationKey(from, fromCol) != columnCollationKey(to, toCol)
			needDetachType := colUsesModColl || colUsesRebuiltDom
			needDropDefault := exprUsesModColl || (needDetachType && fromCol.HasDefault && fromCol.Default != "")

			if needDropDefault {
				dropDefSQL := fmt.Sprintf("%s DROP DEFAULT", colRef)
				skipMainDropDefaultSQL[dropDefSQL] = true
				preExtra = append(preExtra, MigrationOp{
					Type:          OpAlterColumn,
					SchemaName:    scName,
					ObjectName:    fromRel.Name,
					SQL:           dropDefSQL,
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        fromRel.OID,
					Priority:      PriorityColumn,
				})
			}

			if needDetachType {
				extraViewTables[fromRel.OID] = true
				if colUsesRebuiltDom {
					for _, con := range from.ConstraintsOf(fromRel.OID) {
						if con.Type == ConstraintCheck && containsColumnRef(con.CheckExpr, fromCol.Name) {
							dropAndRestoreConstraint(con)
						}
					}
				}
				detachType := fromTypeStr
				if colUsesRebuiltDom {
					detachType = rebuiltDomBase
				}
				preExtra = append(preExtra, MigrationOp{
					Type:          OpAlterColumn,
					SchemaName:    scName,
					ObjectName:    fromRel.Name,
					SQL:           fmt.Sprintf("%s TYPE %s", colRef, detachType),
					Transactional: true,
					Phase:         PhasePre,
					ObjType:       'r',
					ObjOID:        fromRel.OID,
					Priority:      PriorityColumn,
				})
				if !alterTypeAlreadyInMain {
					mainExtra = append(mainExtra, MigrationOp{
						Type:          OpAlterColumn,
						SchemaName:    scName,
						ObjectName:    fromRel.Name,
						SQL:           fmt.Sprintf("%s TYPE %s%s", colRef, toTypeStr, formatColumnCollate(to, toCol)),
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        toRel.OID,
						Priority:      PriorityColumn,
					})
					if stor := columnExplicitStorage(to, toCol); stor != "" && columnExplicitStorage(from, fromCol) == stor {
						mainExtra = append(mainExtra, MigrationOp{
							Type:          OpAlterColumn,
							SchemaName:    scName,
							ObjectName:    fromRel.Name,
							SQL:           fmt.Sprintf("%s SET STORAGE %s", colRef, stor),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        toRel.OID,
							Priority:      PriorityColumn,
						})
					}
					if comp := columnExplicitCompression(toCol); comp != "" && fromCol.Compression == toCol.Compression {
						mainExtra = append(mainExtra, MigrationOp{
							Type:          OpAlterColumn,
							SchemaName:    scName,
							ObjectName:    fromRel.Name,
							SQL:           fmt.Sprintf("%s SET COMPRESSION %s", colRef, comp),
							Transactional: true,
							Phase:         PhaseMain,
							ObjType:       'r',
							ObjOID:        toRel.OID,
							Priority:      PriorityColumn,
						})
					}
				}
			}

			if needDropDefault && toCol.HasDefault && toCol.Default != "" {
				setDefSQL := fmt.Sprintf("%s SET DEFAULT %s", colRef, toCol.Default)
				alreadySets := false
				for _, op := range ops {
					if op.Phase == PhaseMain && op.SQL == setDefSQL {
						alreadySets = true
						break
					}
				}
				if !alreadySets {
					mainExtra = append(mainExtra, MigrationOp{
						Type:          OpAlterColumn,
						SchemaName:    scName,
						ObjectName:    fromRel.Name,
						SQL:           setDefSQL,
						Transactional: true,
						Phase:         PhaseMain,
						ObjType:       'r',
						ObjOID:        toRel.OID,
						Priority:      PriorityColumn,
					})
				}
			}
		}
	}

	if len(skipMainDropDefaultSQL) > 0 || len(rebuiltGenColPrefix) > 0 || len(flipCompDropAttrSQL) > 0 {
		filtered := ops[:0]
		for _, op := range ops {
			if op.Phase == PhaseMain && skipMainDropDefaultSQL[op.SQL] {
				continue
			}
			if op.Phase == PhaseMain && op.Type == OpAlterColumn {
				skipGen := false
				for pfx := range rebuiltGenColPrefix {
					if strings.HasPrefix(op.SQL, pfx+" ") {
						skipGen = true
						break
					}
				}
				if skipGen {
					continue
				}
			}
			if fromOID, ok := flipCompDropAttrSQL[op.SQL]; ok && op.Phase == PhaseMain {
				op.Phase = PhasePre
				op.ObjOID = fromOID
				op.Priority = PriorityColumn
			}
			filtered = append(filtered, op)
		}
		ops = filtered
	}

	// 3. Rebuild kept constraints, indexes, triggers, policies, functions, and views
	// that directly reference a modified collation or rebuilt domain/range type in from.deps.
	for _, d := range from.deps {
		if d.DepType == DepInternal {
			continue
		}
		refMatch := (d.RefType == 'L' && modifiedCollOIDs[d.RefOID]) || (d.RefType == 't' && rebuiltTypeRefOIDs[d.RefOID])
		if !refMatch {
			continue
		}
		switch d.ObjType {
		case 'r':
			if d.ObjSubID == 0 {
				if rel := from.GetRelationByOID(d.ObjOID); rel != nil && (rel.RelKind == 'v' || rel.RelKind == 'm') && !preDroppedRel[rel.OID] {
					directViews[rel.OID] = true
				}
			}
		case 'c':
			if preDroppedCon[d.ObjOID] {
				continue
			}
			if con := from.constraints[d.ObjOID]; con != nil {
				if con.RelOID != 0 && !preDroppedRel[con.RelOID] && con.ConParentID == 0 {
					fromRel := from.GetRelationByOID(con.RelOID)
					if fromRel == nil || fromRel.Schema == nil {
						continue
					}
					scName := fromRel.Schema.Name
					preDroppedCon[con.OID] = true
					preExtra = append(preExtra, buildDropConstraintOp(scName, fromRel.Name, con))
					conKey := scName + "." + fromRel.Name + "." + con.Name
					if !mainAddedCon[conKey] {
						if toRel := to.GetRelation(scName, fromRel.Name); toRel != nil {
							for _, toCon := range to.ConstraintsOf(toRel.OID) {
								if toCon.Name == con.Name && toCon.Type != ConstraintTrigger {
									mainAddedCon[conKey] = true
									mainExtra = append(mainExtra, buildAddConstraintOp(to, scName, fromRel.Name, toCon))
									break
								}
							}
						}
					}
				}
			} else {
				for domOID, dt := range from.domainTypes {
					if preDroppedType[domOID] || dt == nil {
						continue
					}
					for _, fromDC := range dt.Constraints {
						if fromDC.OID != d.ObjOID {
							continue
						}
						bt := from.typeByOID[domOID]
						if bt == nil {
							break
						}
						sc := from.schemas[bt.Namespace]
						if sc == nil {
							break
						}
						qn := migrationQualifiedName(sc.Name, bt.TypeName)
						preDroppedCon[fromDC.OID] = true
						preExtra = append(preExtra, MigrationOp{
							Type:          OpDropConstraint,
							SchemaName:    sc.Name,
							ObjectName:    fromDC.Name,
							ParentObject:  bt.TypeName,
							SQL:           fmt.Sprintf("ALTER DOMAIN %s DROP CONSTRAINT %s", qn, quoteIdentAlways(fromDC.Name)),
							Transactional: true,
							Phase:         PhasePre,
							ObjType:       'c',
							ObjOID:        fromDC.OID,
							Priority:      PriorityConstraint,
						})
						if toOID := resolveTypeOIDByName(to, sc.Name, bt.TypeName); toOID != 0 {
							if toDT := to.domainTypes[toOID]; toDT != nil {
								for _, dc := range toDT.Constraints {
									if dc.Name == fromDC.Name {
										mainExtra = append(mainExtra, MigrationOp{
											Type:          OpAddConstraint,
											SchemaName:    sc.Name,
											ObjectName:    dc.Name,
											ParentObject:  bt.TypeName,
											SQL:           fmt.Sprintf("ALTER DOMAIN %s ADD CONSTRAINT %s CHECK (%s)", qn, quoteIdentAlways(dc.Name), dc.CheckExpr),
											Transactional: true,
											Phase:         PhaseMain,
											ObjType:       't',
											ObjOID:        toOID,
											Priority:      PriorityConstraint,
										})
										break
									}
								}
							}
						}
						break
					}
				}
			}
		case 'i':
			idx := from.indexes[d.ObjOID]
			if idx == nil || idx.Schema == nil || preDroppedRel[idx.RelOID] || idx.ParentIndexOID != 0 {
				continue
			}
			if idx.ConstraintOID != 0 {
				if con := from.constraints[idx.ConstraintOID]; con != nil && !preDroppedCon[con.OID] && con.ConParentID == 0 {
					fromRel := from.GetRelationByOID(con.RelOID)
					if fromRel != nil && fromRel.Schema != nil {
						scName := fromRel.Schema.Name
						preDroppedCon[con.OID] = true
						preExtra = append(preExtra, buildDropConstraintOp(scName, fromRel.Name, con))
						conKey := scName + "." + fromRel.Name + "." + con.Name
						if !mainAddedCon[conKey] {
							if toRel := to.GetRelation(scName, fromRel.Name); toRel != nil {
								for _, toCon := range to.ConstraintsOf(toRel.OID) {
									if toCon.Name == con.Name && toCon.Type != ConstraintTrigger {
										mainAddedCon[conKey] = true
										mainExtra = append(mainExtra, buildAddConstraintOp(to, scName, fromRel.Name, toCon))
										break
									}
								}
							}
						}
					}
				}
				continue
			}
			if preDroppedIdx[idx.OID] {
				continue
			}
			scName := idx.Schema.Name
			preDroppedIdx[idx.OID] = true
			preExtra = append(preExtra, MigrationOp{
				Type:          OpDropIndex,
				SchemaName:    scName,
				ObjectName:    idx.Name,
				SQL:           fmt.Sprintf("DROP INDEX %s.%s", quoteIdentifier(scName), quoteIdentifier(idx.Name)),
				Transactional: true,
				Phase:         PhasePre,
				ObjType:       'i',
				ObjOID:        idx.OID,
				Priority:      PriorityIndex,
			})
			idxKey := scName + "." + idx.Name
			if !mainCreatedIdx[idxKey] {
				if toSc := to.GetSchema(scName); toSc != nil {
					if toIdx := toSc.Indexes[idx.Name]; toIdx != nil && toIdx.ConstraintOID == 0 {
						if toRel := to.GetRelationByOID(toIdx.RelOID); toRel != nil {
							mainCreatedIdx[idxKey] = true
							mainExtra = append(mainExtra, MigrationOp{
								Type:          OpCreateIndex,
								SchemaName:    scName,
								ObjectName:    toIdx.Name,
								SQL:           buildCreateIndexSQL(toIdx, toRel, scName),
								Transactional: true,
								Phase:         PhaseMain,
								ObjType:       'i',
								ObjOID:        toIdx.OID,
								Priority:      PriorityIndex,
							})
							mainExtra = append(mainExtra, buildIndexColumnRenameOps(scName, toIdx, defaultIndexColumnNames(to, toIdx), effectiveIndexColumnNames(to, toIdx))...)
						}
					}
				}
			}
		case 'g':
			trig := from.triggers[d.ObjOID]
			if trig == nil || preDroppedTrig[trig.OID] || preDroppedRel[trig.RelOID] || trig.ParentTriggerOID != 0 {
				continue
			}
			fromRel := from.GetRelationByOID(trig.RelOID)
			if fromRel == nil || fromRel.Schema == nil {
				continue
			}
			scName := fromRel.Schema.Name
			preDroppedTrig[trig.OID] = true
			preExtra = append(preExtra, buildDropTriggerOp(scName, fromRel.Name, trig.Name, trig.OID))
			trigKey := scName + "." + fromRel.Name + "." + trig.Name
			if !mainCreatedTrig[trigKey] {
				if toRel := to.GetRelation(scName, fromRel.Name); toRel != nil {
					for _, toTrig := range to.triggersByRel[toRel.OID] {
						if toTrig.Name == trig.Name && toTrig.ParentTriggerOID == 0 {
							mainCreatedTrig[trigKey] = true
							mainExtra = append(mainExtra, buildCreateTriggerOp(to, scName, fromRel.Name, toTrig))
							break
						}
					}
				}
			}
		case 'p':
			pol := from.policies[d.ObjOID]
			if pol == nil || preDroppedPol[pol.OID] || preDroppedRel[pol.RelOID] {
				continue
			}
			fromRel := from.GetRelationByOID(pol.RelOID)
			if fromRel == nil || fromRel.Schema == nil {
				continue
			}
			scName := fromRel.Schema.Name
			preDroppedPol[pol.OID] = true
			preExtra = append(preExtra, buildDropPolicyOp(scName, fromRel.Name, pol.Name, pol.OID))
			polKey := scName + "." + pol.Name
			if !mainCreatedPol[polKey] {
				if toRel := to.GetRelation(scName, fromRel.Name); toRel != nil {
					for _, toPol := range to.policiesByRel[toRel.OID] {
						if toPol.Name == pol.Name {
							mainCreatedPol[polKey] = true
							mainExtra = append(mainExtra, buildCreatePolicyOp(scName, fromRel.Name, toPol))
							break
						}
					}
				}
			}
		case 'R':
			rule := from.rules[d.ObjOID]
			if rule == nil || preDroppedRule[rule.OID] || preDroppedRel[rule.OwnerRelOID] {
				continue
			}
			fromRel := from.GetRelationByOID(rule.OwnerRelOID)
			if fromRel == nil || fromRel.Schema == nil {
				continue
			}
			scName := fromRel.Schema.Name
			preDroppedRule[rule.OID] = true
			preExtra = append(preExtra, buildDropRuleOp(scName, fromRel.Name, rule))
			ruleKey := scName + "." + fromRel.Name + "." + rule.Name
			if !mainCreatedRule[ruleKey] {
				if toRel := to.GetRelation(scName, fromRel.Name); toRel != nil && toRel.Rules != nil {
					if toRule := toRel.Rules[rule.Name]; toRule != nil {
						mainCreatedRule[ruleKey] = true
						mainExtra = append(mainExtra, buildCreateAndModeRuleOps(to, scName, toRel.Name, toRel, toRule)...)
					}
				}
			}
		case 'f':
			up := from.userProcs[d.ObjOID]
			if up == nil || up.Schema == nil || preDroppedFunc[up.OID] {
				continue
			}
			scName := up.Schema.Name
			id := funcIdentity(from, up)
			preDroppedFunc[up.OID] = true
			preExtra = append(preExtra, buildDropFunctionOp(from, FunctionDiffEntry{Action: DiffDrop, SchemaName: scName, Identity: id, From: up}))
			fnKey := scName + "." + id
			if !mainCreatedFunc[fnKey] {
				for _, toUP := range to.userProcs {
					if toUP.Schema != nil && toUP.Schema.Name == scName && funcIdentity(to, toUP) == id {
						mainCreatedFunc[fnKey] = true
						mainExtra = append(mainExtra, buildCreateFunctionOp(to, FunctionDiffEntry{Action: DiffAdd, SchemaName: scName, Identity: id, To: toUP}))
						break
					}
				}
			} else {
				// Replace CREATE OR REPLACE in main with plain CREATE since we dropped in PhasePre.
				for i := range ops {
					if ops[i].Phase == PhaseMain && ops[i].Type == OpAlterFunction && ops[i].SchemaName == scName && ops[i].ObjectName == id {
						for _, toUP := range to.userProcs {
							if toUP.Schema != nil && toUP.Schema.Name == scName && funcIdentity(to, toUP) == id {
								ops[i] = buildCreateFunctionOp(to, FunctionDiffEntry{Action: DiffAdd, SchemaName: scName, Identity: id, To: toUP})
								break
							}
						}
					}
				}
			}
		}
	}

	ops = append(ops, preExtra...)
	ops = append(ops, mainExtra...)
	return ops, extraViewTables, directViews
}
