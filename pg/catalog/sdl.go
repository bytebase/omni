package catalog

import (
	"container/heap"
	"fmt"
	"sort"
	"strings"

	nodes "github.com/bytebase/omni/pg/ast"
	pgparser "github.com/bytebase/omni/pg/parser"
)

// LoadSDL parses a declarative schema definition (SDL) and loads it into a new
// Catalog. SDL only accepts CREATE/COMMENT/GRANT statements — DML and
// destructive DDL (DROP, ALTER TABLE ADD/DROP COLUMN, TRUNCATE, etc.) are
// rejected with a clear error.
//
// The pipeline: parse → validate → collectDeclared → extractDeps → topoSort → execute.
func LoadSDL(sql string) (*Catalog, error) {
	c := New()

	list, err := pgparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return c, nil
	}

	// Unwrap RawStmt wrappers and collect bare statements.
	stmts := make([]nodes.Node, 0, len(list.Items))
	for _, item := range list.Items {
		if raw, ok := item.(*nodes.RawStmt); ok {
			stmts = append(stmts, raw.Stmt)
		} else {
			stmts = append(stmts, item)
		}
	}

	// Validate all statements before execution.
	if err := validateSDL(stmts); err != nil {
		return nil, err
	}

	// Collect declared objects for closed-set matching.
	declared := collectDeclaredObjects(stmts)

	// Extract dependencies and topologically sort.
	deps := extractDeps(stmts, declared)
	ordered, deferred, err := topoSort(stmts, declared, deps)
	if err != nil {
		return nil, err
	}

	// Pre-create shell types (and their array types) for all composite types.
	// This resolves mutual references between composite types (e.g., type A
	// referencing type B[] and vice versa) by ensuring all type names exist
	// before any full definition.
	for _, stmt := range ordered {
		if ct, ok := stmt.(*nodes.CompositeTypeStmt); ok && ct.Typevar != nil {
			c.createCompositeShellType(ct.Typevar)
		}
	}

	// Execute in topologically sorted order.
	for _, stmt := range ordered {
		// Before executing a CompositeTypeStmt, remove its shell type so
		// DefineRelation does not see a conflicting type name.
		if ct, ok := stmt.(*nodes.CompositeTypeStmt); ok && ct.Typevar != nil {
			c.removeShellType(ct.Typevar)
		}
		if err := c.ProcessUtility(stmt); err != nil {
			return c, err
		}
	}

	// Execute deferred FK constraints as ALTER TABLE ADD CONSTRAINT.
	for _, fk := range deferred {
		if err := c.ProcessUtility(fk); err != nil {
			return c, err
		}
	}

	return c, nil
}

// ---------------------------------------------------------------------------
// Section 1.3: Declared Object Collection and Name Resolution
// ---------------------------------------------------------------------------

// collectDeclaredObjects scans all statements and returns the set of declared
// object names. Unqualified names get a "public." prefix. Functions include
// argument types in their identity: "public.myfunc(integer,text)".
func collectDeclaredObjects(stmts []nodes.Node) map[string]bool {
	declared := make(map[string]bool)
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *nodes.CreateStmt:
			if s.Relation != nil {
				declared[qualifiedRangeVar(s.Relation)] = true
			}
		case *nodes.ViewStmt:
			if s.View != nil {
				declared[qualifiedRangeVar(s.View)] = true
			}
		case *nodes.IndexStmt:
			// Indexes are not referenced by name in deps, skip.
		case *nodes.CreateSeqStmt:
			if s.Sequence != nil {
				declared[qualifiedRangeVar(s.Sequence)] = true
			}
		case *nodes.CreateSchemaStmt:
			if s.Schemaname != "" {
				declared[s.Schemaname] = true
			}
		case *nodes.CreateEnumStmt:
			declared[qualifiedNameFromList(s.TypeName)] = true
		case *nodes.CreateDomainStmt:
			declared[qualifiedNameFromList(s.Domainname)] = true
		case *nodes.CompositeTypeStmt:
			if s.Typevar != nil {
				declared[qualifiedRangeVar(s.Typevar)] = true
			}
		case *nodes.CreateRangeStmt:
			declared[qualifiedNameFromList(s.TypeName)] = true
		case *nodes.CreateFunctionStmt:
			declared[functionIdentity(s)] = true
		case *nodes.CreateExtensionStmt:
			// Extensions are identified by name only.
			if s.Extname != "" {
				declared["extension:"+s.Extname] = true
			}
		case *nodes.CreateTrigStmt:
			// Triggers are not referenced by name in deps.
		case *nodes.CreatePolicyStmt:
			// Policies are not referenced by name in deps.
		case *nodes.CreateTableAsStmt:
			if s.Into != nil && s.Into.Rel != nil {
				declared[qualifiedRangeVar(s.Into.Rel)] = true
			}
		case *nodes.CreateForeignTableStmt:
			if s.Base.Relation != nil {
				declared[qualifiedRangeVar(s.Base.Relation)] = true
			}
		}
	}
	return declared
}

// qualifiedRangeVar returns "schema.name" for a RangeVar, defaulting to "public".
func qualifiedRangeVar(rv *RangeVar) string {
	schema := rv.Schemaname
	if schema == "" {
		schema = "public"
	}
	return schema + "." + rv.Relname
}

// qualifiedNameFromList extracts a qualified name from a List of String nodes.
// For a single name, prepends "public.". For two names, uses "schema.name".
func qualifiedNameFromList(l *nodes.List) string {
	if l == nil || len(l.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(l.Items))
	for _, item := range l.Items {
		if s, ok := item.(*nodes.String); ok {
			parts = append(parts, s.Str)
		}
	}
	if len(parts) == 1 {
		return "public." + parts[0]
	}
	return strings.Join(parts, ".")
}

// functionIdentity returns a unique identity for a function including arg types:
// "public.myfunc(integer,text)".
func functionIdentity(s *nodes.CreateFunctionStmt) string {
	name := qualifiedNameFromList(s.Funcname)
	argTypes := make([]string, 0)
	if s.Parameters != nil {
		for _, item := range s.Parameters.Items {
			fp, ok := item.(*nodes.FunctionParameter)
			if !ok {
				continue
			}
			// Only IN/INOUT/VARIADIC params are part of the identity.
			if fp.Mode == nodes.FUNC_PARAM_OUT || fp.Mode == nodes.FUNC_PARAM_TABLE {
				continue
			}
			if fp.ArgType != nil {
				argTypes = append(argTypes, typeNameToString(fp.ArgType))
			}
		}
	}
	return name + "(" + strings.Join(argTypes, ",") + ")"
}

// typeNameToString extracts the last name from a TypeName's Names list.
// For built-in types like "pg_catalog.int4", returns "integer" (the common alias).
func typeNameToString(tn *nodes.TypeName) string {
	if tn == nil || tn.Names == nil || len(tn.Names.Items) == 0 {
		return ""
	}
	// Take the last name component.
	last := ""
	for _, item := range tn.Names.Items {
		if s, ok := item.(*nodes.String); ok {
			last = s.Str
		}
	}
	return last
}

// typeNameQualified returns the qualified type name from a TypeName.
// Returns "schema.name" or "public.name" for unqualified names.
func typeNameQualified(tn *nodes.TypeName) string {
	if tn == nil || tn.Names == nil || len(tn.Names.Items) == 0 {
		return ""
	}
	parts := make([]string, 0, len(tn.Names.Items))
	for _, item := range tn.Names.Items {
		if s, ok := item.(*nodes.String); ok {
			parts = append(parts, s.Str)
		}
	}
	if len(parts) == 1 {
		return "public." + parts[0]
	}
	// For "pg_catalog.int4" etc, return as-is.
	return strings.Join(parts, ".")
}

// createCompositeShellType creates a shell type and its array type for a composite type.
// This allows mutual composite type references via arrays to be resolved.
func (c *Catalog) createCompositeShellType(rv *nodes.RangeVar) {
	schemaName := rv.Schemaname
	if schemaName == "" {
		schemaName = "public"
	}
	schema := c.schemaByName[schemaName]
	if schema == nil {
		return
	}
	name := rv.Relname
	key := typeKey{ns: schema.OID, name: name}
	if c.typeByName[key] != nil {
		return // already exists
	}

	// Create shell type.
	shellOID := c.oidGen.Next()
	shellType := &BuiltinType{
		OID:       shellOID,
		TypeName:  name,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'c',
		Category:  'C',
		IsDefined: false,
		Delim:     ',',
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	c.typeByOID[shellOID] = shellType
	c.typeByName[key] = shellType

	// Create array type (_name) so references like "name[]" resolve.
	arrayOID := c.oidGen.Next()
	arrayName := "_" + name
	arrayType := &BuiltinType{
		OID:       arrayOID,
		TypeName:  arrayName,
		Namespace: schema.OID,
		Len:       -1,
		ByVal:     false,
		Type:      'b',
		Category:  'A',
		IsDefined: false,
		Delim:     ',',
		Elem:      shellOID,
		Align:     'd',
		Storage:   'x',
		TypeMod:   -1,
	}
	shellType.Array = arrayOID
	c.typeByOID[arrayOID] = arrayType
	c.typeByName[typeKey{ns: schema.OID, name: arrayName}] = arrayType
}

// removeShellType removes an undefined shell type entry (and its array type)
// for a composite type so that DefineRelation does not see a conflict.
func (c *Catalog) removeShellType(rv *nodes.RangeVar) {
	schemaName := rv.Schemaname
	if schemaName == "" {
		schemaName = "public"
	}
	schema := c.schemaByName[schemaName]
	if schema == nil {
		return
	}
	key := typeKey{ns: schema.OID, name: rv.Relname}
	t := c.typeByName[key]
	if t == nil || t.IsDefined {
		return
	}
	// Remove the array type if it exists and is also a shell.
	if t.Array != 0 {
		arrayKey := typeKey{ns: schema.OID, name: "_" + rv.Relname}
		if at := c.typeByName[arrayKey]; at != nil && !at.IsDefined {
			delete(c.typeByName, arrayKey)
			delete(c.typeByOID, at.OID)
		}
	}
	delete(c.typeByName, key)
	delete(c.typeByOID, t.OID)
}

// rangeVarToNameList converts a RangeVar to a List of String nodes
// suitable for DefineStmt.Defnames.
func rangeVarToNameList(rv *nodes.RangeVar) *nodes.List {
	var items []nodes.Node
	if rv.Schemaname != "" {
		items = append(items, &nodes.String{Str: rv.Schemaname})
	}
	items = append(items, &nodes.String{Str: rv.Relname})
	return &nodes.List{Items: items}
}

// RangeVar is an alias used in this file for convenience.
type RangeVar = nodes.RangeVar

// ---------------------------------------------------------------------------
// Section 1.2: Basic Dependency Resolution
// ---------------------------------------------------------------------------

// dep represents a dependency: statement at index `from` depends on statement at index `to`.
type dep struct {
	from int  // index of the dependent statement
	to   int  // index of the dependency
	isFK bool // true if this dep originates from a FK constraint (deferred, not in within-layer sort)
}

// stmtPriority assigns a priority layer to each statement type for ordering.
// Lower priority = created first.
func stmtPriority(stmt nodes.Node) int {
	switch stmt.(type) {
	case *nodes.CreateSchemaStmt:
		return 0
	case *nodes.CreateExtensionStmt:
		return 1
	case *nodes.CreateEnumStmt, *nodes.CreateDomainStmt, *nodes.CompositeTypeStmt, *nodes.CreateRangeStmt:
		return 2
	case *nodes.CreateSeqStmt:
		return 3
	case *nodes.CreateStmt, *nodes.CreateForeignTableStmt:
		return 4
	case *nodes.CreateFunctionStmt:
		return 5
	case *nodes.ViewStmt, *nodes.CreateTableAsStmt:
		return 6
	case *nodes.IndexStmt:
		return 7
	case *nodes.CreateTrigStmt:
		return 7
	case *nodes.CreatePolicyStmt:
		return 7
	case *nodes.AlterSeqStmt:
		return 8
	case *nodes.AlterTableStmt:
		return 9
	case *nodes.CommentStmt:
		return 9
	case *nodes.GrantStmt:
		return 9
	case *nodes.AlterEnumStmt:
		return 9
	default:
		return 10
	}
}

// stmtName returns the declared object name for a statement, used for dep matching.
func stmtName(stmt nodes.Node) string {
	switch s := stmt.(type) {
	case *nodes.CreateStmt:
		if s.Relation != nil {
			return qualifiedRangeVar(s.Relation)
		}
	case *nodes.ViewStmt:
		if s.View != nil {
			return qualifiedRangeVar(s.View)
		}
	case *nodes.CreateSeqStmt:
		if s.Sequence != nil {
			return qualifiedRangeVar(s.Sequence)
		}
	case *nodes.CreateSchemaStmt:
		return s.Schemaname
	case *nodes.CreateEnumStmt:
		return qualifiedNameFromList(s.TypeName)
	case *nodes.CreateDomainStmt:
		return qualifiedNameFromList(s.Domainname)
	case *nodes.CompositeTypeStmt:
		if s.Typevar != nil {
			return qualifiedRangeVar(s.Typevar)
		}
	case *nodes.CreateRangeStmt:
		return qualifiedNameFromList(s.TypeName)
	case *nodes.CreateFunctionStmt:
		return functionIdentity(s)
	case *nodes.CreateExtensionStmt:
		return "extension:" + s.Extname
	case *nodes.CreateTableAsStmt:
		if s.Into != nil && s.Into.Rel != nil {
			return qualifiedRangeVar(s.Into.Rel)
		}
	case *nodes.CreateForeignTableStmt:
		if s.Base.Relation != nil {
			return qualifiedRangeVar(s.Base.Relation)
		}
	}
	return ""
}

// deferredFK represents a FK constraint that was extracted from a CreateStmt
// and must be applied after all tables are created.
type deferredFK struct {
	tableName *nodes.RangeVar
	cons      *nodes.Constraint
}

// extractDeps extracts dependencies from direct AST fields only (no expression walking).
func extractDeps(stmts []nodes.Node, declared map[string]bool) []dep {
	// Build name→index mapping.
	nameToIdx := make(map[string]int, len(stmts))
	for i, stmt := range stmts {
		name := stmtName(stmt)
		if name != "" {
			nameToIdx[name] = i
		}
	}

	var deps []dep
	for i, stmt := range stmts {
		refs, fkRefs := extractRefs(stmt, declared)
		for _, ref := range refs {
			if j, ok := nameToIdx[ref]; ok && j != i {
				deps = append(deps, dep{from: i, to: j})
			}
		}
		for _, ref := range fkRefs {
			if j, ok := nameToIdx[ref]; ok && j != i {
				deps = append(deps, dep{from: i, to: j, isFK: true})
			}
		}
	}
	return deps
}

// extractRefs returns the set of declared object names that a statement references.
// It returns two slices: regular refs and FK-only refs. FK refs are separated because
// FK constraints are deferred and should not contribute to within-layer cycle detection.
func extractRefs(stmt nodes.Node, declared map[string]bool) (refs []string, fkRefs []string) {
	addRef := func(name string) {
		if name != "" && declared[name] {
			refs = append(refs, name)
		}
	}
	addFKRef := func(name string) {
		if name != "" && declared[name] {
			fkRefs = append(fkRefs, name)
		}
	}

	// addExprDeps is a helper that collects expression-level deps and adds matching ones.
	addExprDeps := func(funcRefs, relRefs, typeRefs []string) {
		for _, fr := range funcRefs {
			// For function refs, try matching with "()" suffix for no-arg functions.
			addRef(fr + "()")
			// Also try matching the bare qualified name against declared set
			// in case the function identity matches with args.
			for d := range declared {
				if strings.HasPrefix(d, fr+"(") {
					addRef(d)
				}
			}
		}
		for _, rr := range relRefs {
			addRef(rr)
		}
		for _, tr := range typeRefs {
			addRef(tr)
		}
	}

	switch s := stmt.(type) {
	case *nodes.CreateStmt:
		// Column types.
		if s.TableElts != nil {
			for _, elt := range s.TableElts.Items {
				if col, ok := elt.(*nodes.ColumnDef); ok && col.TypeName != nil {
					addRef(typeNameQualified(col.TypeName))
				}
				// Inline FK constraint on column — FK refs go to fkRefs.
				if col, ok := elt.(*nodes.ColumnDef); ok && col.Constraints != nil {
					for _, c := range col.Constraints.Items {
						if cons, ok := c.(*nodes.Constraint); ok {
							if cons.Contype == nodes.CONSTR_FOREIGN && cons.Pktable != nil {
								addFKRef(qualifiedRangeVar(cons.Pktable))
							}
						}
					}
				}
				// Standalone FK constraint in TableElts (pgparser puts table-level
				// CONSTRAINT ... FOREIGN KEY nodes here).
				if cons, ok := elt.(*nodes.Constraint); ok {
					if cons.Contype == nodes.CONSTR_FOREIGN && cons.Pktable != nil {
						addFKRef(qualifiedRangeVar(cons.Pktable))
					}
				}
			}
		}
		// Table-level constraints in Constraints list: FK deps go to fkRefs.
		if s.Constraints != nil {
			for _, item := range s.Constraints.Items {
				if cons, ok := item.(*nodes.Constraint); ok {
					if cons.Contype == nodes.CONSTR_FOREIGN && cons.Pktable != nil {
						addFKRef(qualifiedRangeVar(cons.Pktable))
					}
				}
			}
		}
		// INHERITS.
		if s.InhRelations != nil {
			for _, item := range s.InhRelations.Items {
				if rv, ok := item.(*nodes.RangeVar); ok {
					addRef(qualifiedRangeVar(rv))
				}
			}
		}
		// Expression deps from defaults and CHECK constraints.
		fr, rr, tr := collectExprDepsFromCreateStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.ViewStmt:
		// Expression deps from the view query (replaces old extractSelectRefs).
		fr, rr, tr := collectExprDepsFromViewStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.IndexStmt:
		if s.Relation != nil {
			addRef(qualifiedRangeVar(s.Relation))
		}
		// Expression deps from index expressions and WHERE clause.
		fr, rr, tr := collectExprDepsFromIndexStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CreateTrigStmt:
		if s.Relation != nil {
			addRef(qualifiedRangeVar(s.Relation))
		}
		if s.Funcname != nil {
			// Build function identity — for triggers, function takes no args usually.
			funcName := qualifiedNameFromList(s.Funcname)
			// Try the no-arg version first.
			addRef(funcName + "()")
		}
		// Expression deps from WHEN clause.
		fr, rr, tr := collectExprDepsFromTriggerStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CreatePolicyStmt:
		if s.Table != nil {
			addRef(qualifiedRangeVar(s.Table))
		}
		// Expression deps from USING and WITH CHECK.
		fr, rr, tr := collectExprDepsFromPolicyStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CreateDomainStmt:
		// Expression deps from CHECK constraints and base type.
		fr, rr, tr := collectExprDepsFromDomainStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CreateFunctionStmt:
		// Structural deps from parameter types, return type, param defaults.
		fr, rr, tr := collectStructuralDepsFromFunctionStmt(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CompositeTypeStmt:
		// Structural deps from column types.
		fr, rr, tr := collectStructuralDepsFromCompositeType(s)
		addExprDeps(fr, rr, tr)

	case *nodes.CreateRangeStmt:
		// Structural deps from subtype.
		fr, rr, tr := collectStructuralDepsFromRangeType(s)
		addExprDeps(fr, rr, tr)

	case *nodes.AlterSeqStmt:
		// OWNED BY: look for "owned_by" in options.
		if s.Options != nil {
			for _, item := range s.Options.Items {
				if de, ok := item.(*nodes.DefElem); ok && de.Defname == "owned_by" {
					// The arg is a List of String nodes like ["tablename", "colname"].
					if l, ok := de.Arg.(*nodes.List); ok && len(l.Items) >= 2 {
						// Extract table name (all but last element).
						parts := make([]string, 0, len(l.Items)-1)
						for j := 0; j < len(l.Items)-1; j++ {
							if str, ok := l.Items[j].(*nodes.String); ok {
								parts = append(parts, str.Str)
							}
						}
						var tableName string
						if len(parts) == 1 {
							tableName = "public." + parts[0]
						} else {
							tableName = strings.Join(parts, ".")
						}
						addRef(tableName)
					}
				}
			}
		}

	case *nodes.AlterTableStmt:
		if s.Relation != nil {
			addRef(qualifiedRangeVar(s.Relation))
		}

	case *nodes.CommentStmt:
		// Comment targets are complex; for now just note that comments
		// go in the last priority layer so they execute after everything.

	case *nodes.GrantStmt:
		// Grants go in the last priority layer.
	}

	return refs, fkRefs
}

// extractSelectRefs extracts top-level RangeVar references from a SelectStmt's FromClause.
func extractSelectRefs(sel *nodes.SelectStmt, declared map[string]bool, refs *[]string) {
	if sel == nil {
		return
	}
	if sel.FromClause != nil {
		for _, item := range sel.FromClause.Items {
			if rv, ok := item.(*nodes.RangeVar); ok {
				name := qualifiedRangeVar(rv)
				if declared[name] {
					*refs = append(*refs, name)
				}
			}
		}
	}
}

// topoSort sorts statements using global topological sorting (Kahn's algorithm).
// Cross-layer dependencies (e.g., table CHECK → function) take precedence over
// default priority ordering. When no dependency constrains order, statements are
// ordered by (stmtPriority, original_index) for determinism.
// FK constraints are extracted from CreateStmts and deferred.
// Returns the ordered statements and the deferred FK ALTER TABLE statements.
func topoSort(stmts []nodes.Node, declared map[string]bool, deps []dep) ([]nodes.Node, []nodes.Node, error) {
	n := len(stmts)

	// Compute priority for each statement.
	priority := make([]int, n)
	for i, stmt := range stmts {
		priority[i] = stmtPriority(stmt)
	}

	// Build adjacency from ALL non-FK deps.
	// dep.from depends on dep.to → edge from dep.to to dep.from (dep.to must come first).
	adj := make([][]int, n)
	inDegree := make([]int, n)
	for _, d := range deps {
		if d.isFK {
			continue
		}
		adj[d.to] = append(adj[d.to], d.from)
		inDegree[d.from]++
	}

	// Kahn's algorithm with priority queue.
	// Tie-break by (stmtPriority, original_index) for determinism.
	// Use a heap for efficient min extraction.
	h := make(pqHeap, 0, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			h = append(h, pqEntry{idx: i, pri: priority[i]})
		}
	}
	heap.Init(&h)

	sorted := make([]int, 0, n)
	for h.Len() > 0 {
		e := heap.Pop(&h).(pqEntry)
		sorted = append(sorted, e.idx)
		for _, next := range adj[e.idx] {
			inDegree[next]--
			if inDegree[next] == 0 {
				heap.Push(&h, pqEntry{idx: next, pri: priority[next]})
			}
		}
	}

	// Handle cycles: if not all nodes were visited, check for composite type cycles.
	if len(sorted) != n {
		sortedSet := make(map[int]bool, len(sorted))
		for _, idx := range sorted {
			sortedSet[idx] = true
		}
		allComposite := true
		var remaining []int
		for i := 0; i < n; i++ {
			if !sortedSet[i] {
				remaining = append(remaining, i)
				if _, ok := stmts[i].(*nodes.CompositeTypeStmt); !ok {
					allComposite = false
				}
			}
		}
		if allComposite && len(remaining) > 0 {
			// Composite type cycles can be resolved via shell types.
			// Sort by (priority, index) for determinism.
			sort.Slice(remaining, func(a, b int) bool {
				pa, pb := priority[remaining[a]], priority[remaining[b]]
				if pa != pb {
					return pa < pb
				}
				return remaining[a] < remaining[b]
			})
			sorted = append(sorted, remaining...)
		} else {
			return nil, nil, fmt.Errorf("SDL dependency cycle detected")
		}
	}

	// Build ordered result, extracting FK constraints for deferral.
	var ordered []nodes.Node
	var deferredFKs []nodes.Node
	for _, idx := range sorted {
		stmt := stmts[idx]
		if cs, ok := stmt.(*nodes.CreateStmt); ok {
			stripped, fks := extractForeignKeys(cs)
			ordered = append(ordered, stripped)
			for _, fk := range fks {
				deferredFKs = append(deferredFKs, buildAlterTableAddConstraint(cs.Relation, fk))
			}
		} else {
			ordered = append(ordered, stmt)
		}
	}

	return ordered, deferredFKs, nil
}

// pqEntry is an entry in the priority queue for topological sorting.
type pqEntry struct {
	idx int
	pri int // stmtPriority
}

// pqHeap implements heap.Interface for Kahn's algorithm priority queue.
// Orders by (stmtPriority ASC, original index ASC).
type pqHeap []pqEntry

func (h pqHeap) Len() int      { return len(h) }
func (h pqHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h pqHeap) Less(i, j int) bool {
	if h[i].pri != h[j].pri {
		return h[i].pri < h[j].pri
	}
	return h[i].idx < h[j].idx
}
func (h *pqHeap) Push(x interface{}) {
	*h = append(*h, x.(pqEntry))
}
func (h *pqHeap) Pop() interface{} {
	old := *h
	n := len(old)
	e := old[n-1]
	*h = old[:n-1]
	return e
}

// extractForeignKeys removes FK constraints from a CreateStmt and returns
// the modified statement and the extracted FK constraints.
func extractForeignKeys(cs *nodes.CreateStmt) (*nodes.CreateStmt, []*nodes.Constraint) {
	var fks []*nodes.Constraint

	// Check table-level constraints in Constraints list.
	if cs.Constraints != nil {
		var kept []nodes.Node
		for _, item := range cs.Constraints.Items {
			if cons, ok := item.(*nodes.Constraint); ok && cons.Contype == nodes.CONSTR_FOREIGN {
				fks = append(fks, cons)
			} else {
				kept = append(kept, item)
			}
		}
		if len(fks) > 0 {
			// Clone the CreateStmt to avoid mutating the original.
			clone := *cs
			if len(kept) == 0 {
				clone.Constraints = nil
			} else {
				clone.Constraints = &nodes.List{Items: kept}
			}
			cs = &clone
		}
	}

	// Check standalone Constraint nodes in TableElts (pgparser puts table-level
	// constraints like CONSTRAINT ... FOREIGN KEY here, not in cs.Constraints).
	if cs.TableElts != nil {
		needCloneForTableElts := false
		for _, elt := range cs.TableElts.Items {
			if cons, ok := elt.(*nodes.Constraint); ok && cons.Contype == nodes.CONSTR_FOREIGN {
				needCloneForTableElts = true
				break
			}
		}
		if needCloneForTableElts {
			clone := *cs
			var newElts []nodes.Node
			for _, elt := range cs.TableElts.Items {
				if cons, ok := elt.(*nodes.Constraint); ok && cons.Contype == nodes.CONSTR_FOREIGN {
					fks = append(fks, cons)
				} else {
					newElts = append(newElts, elt)
				}
			}
			clone.TableElts = &nodes.List{Items: newElts}
			cs = &clone
		}
	}

	// Check inline column constraints.
	if cs.TableElts != nil {
		needClone := false
		for _, elt := range cs.TableElts.Items {
			if col, ok := elt.(*nodes.ColumnDef); ok && col.Constraints != nil {
				for _, c := range col.Constraints.Items {
					if cons, ok := c.(*nodes.Constraint); ok && cons.Contype == nodes.CONSTR_FOREIGN {
						needClone = true
						break
					}
				}
			}
			if needClone {
				break
			}
		}
		if needClone {
			clone := *cs
			newElts := make([]nodes.Node, len(cs.TableElts.Items))
			for i, elt := range cs.TableElts.Items {
				col, ok := elt.(*nodes.ColumnDef)
				if !ok || col.Constraints == nil {
					newElts[i] = elt
					continue
				}
				var colKept []nodes.Node
				for _, c := range col.Constraints.Items {
					if cons, ok := c.(*nodes.Constraint); ok && cons.Contype == nodes.CONSTR_FOREIGN {
						// Build a complete FK constraint with the column name in FkAttrs.
						fullFK := *cons
						if fullFK.FkAttrs == nil || len(fullFK.FkAttrs.Items) == 0 {
							fullFK.FkAttrs = &nodes.List{Items: []nodes.Node{&nodes.String{Str: col.Colname}}}
						}
						fks = append(fks, &fullFK)
					} else {
						colKept = append(colKept, c)
					}
				}
				colClone := *col
				if len(colKept) == 0 {
					colClone.Constraints = nil
				} else {
					colClone.Constraints = &nodes.List{Items: colKept}
				}
				newElts[i] = &colClone
			}
			clone.TableElts = &nodes.List{Items: newElts}
			cs = &clone
		}
	}

	return cs, fks
}

// buildAlterTableAddConstraint builds an ALTER TABLE ADD CONSTRAINT statement
// for a deferred FK constraint.
func buildAlterTableAddConstraint(rel *nodes.RangeVar, cons *nodes.Constraint) *nodes.AlterTableStmt {
	return &nodes.AlterTableStmt{
		Relation: rel,
		Cmds: &nodes.List{
			Items: []nodes.Node{
				&nodes.AlterTableCmd{
					Subtype: int(nodes.AT_AddConstraint),
					Def:     cons,
				},
			},
		},
		ObjType: int(nodes.OBJECT_TABLE),
	}
}

// validateSDL checks that every statement in the list is allowed in SDL.
// Returns the first disallowed statement as an error.
func validateSDL(stmts []nodes.Node) error {
	for _, stmt := range stmts {
		if err := validateSDLStmt(stmt); err != nil {
			return err
		}
	}
	return nil
}

// validateSDLStmt validates a single statement for SDL compliance.
func validateSDLStmt(stmt nodes.Node) error {
	switch s := stmt.(type) {
	// ---- Allowed DDL statements ----
	case *nodes.CreateStmt:
		return nil
	case *nodes.ViewStmt:
		return nil
	case *nodes.CreateFunctionStmt:
		return nil
	case *nodes.IndexStmt:
		return nil
	case *nodes.CreateSeqStmt:
		return nil
	case *nodes.CreateSchemaStmt:
		return nil
	case *nodes.CreateEnumStmt:
		return nil
	case *nodes.CreateDomainStmt:
		return nil
	case *nodes.CompositeTypeStmt:
		return nil
	case *nodes.CreateRangeStmt:
		return nil
	case *nodes.CreateExtensionStmt:
		return nil
	case *nodes.CreateTrigStmt:
		return nil
	case *nodes.CreatePolicyStmt:
		return nil
	case *nodes.CreateTableAsStmt:
		return nil
	case *nodes.CreateCastStmt:
		return nil
	case *nodes.CreateForeignTableStmt:
		return nil
	case *nodes.CommentStmt:
		return nil
	case *nodes.GrantStmt:
		return nil
	case *nodes.AlterSeqStmt:
		return nil
	case *nodes.AlterEnumStmt:
		return nil
	case *nodes.VariableSetStmt:
		return nil
	case *nodes.DefineStmt:
		return nil

	// ---- ALTER TABLE: only RLS commands allowed ----
	case *nodes.AlterTableStmt:
		return validateAlterTableSDL(s)

	// ---- Explicitly rejected DML ----
	case *nodes.InsertStmt:
		return fmt.Errorf("SDL does not allow INSERT statements")
	case *nodes.UpdateStmt:
		return fmt.Errorf("SDL does not allow UPDATE statements")
	case *nodes.DeleteStmt:
		return fmt.Errorf("SDL does not allow DELETE statements")
	case *nodes.SelectStmt:
		return fmt.Errorf("SDL does not allow SELECT statements")

	// ---- Explicitly rejected DDL ----
	case *nodes.DropStmt:
		return fmt.Errorf("SDL does not allow DROP statements")
	case *nodes.TruncateStmt:
		return fmt.Errorf("SDL does not allow TRUNCATE statements")
	case *nodes.DoStmt:
		return fmt.Errorf("SDL does not allow DO statements")

	// ---- Anything else ----
	default:
		return fmt.Errorf("SDL does not allow %T statements", stmt)
	}
}

// validateAlterTableSDL checks that an ALTER TABLE only contains RLS-related
// subcommands (ENABLE/DISABLE/FORCE/NO FORCE ROW LEVEL SECURITY).
func validateAlterTableSDL(s *nodes.AlterTableStmt) error {
	if s.Cmds == nil {
		return nil
	}
	for _, item := range s.Cmds.Items {
		cmd, ok := item.(*nodes.AlterTableCmd)
		if !ok {
			return fmt.Errorf("SDL does not allow ALTER TABLE with unknown command type")
		}
		subtype := nodes.AlterTableType(cmd.Subtype)
		switch subtype {
		case nodes.AT_EnableRowSecurity,
			nodes.AT_DisableRowSecurity,
			nodes.AT_ForceRowSecurity,
			nodes.AT_NoForceRowSecurity:
			// allowed
		default:
			return fmt.Errorf("SDL does not allow ALTER TABLE ADD/DROP COLUMN")
		}
	}
	return nil
}
