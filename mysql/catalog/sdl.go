package catalog

import (
	"container/heap"
	"fmt"

	nodes "github.com/bytebase/omni/mysql/ast"
	mysqlparser "github.com/bytebase/omni/mysql/parser"
)

// LoadSDL parses a declarative MySQL schema (SDL) and loads it into a new
// Catalog order-tolerantly. Unlike LoadSQL — which applies statements in textual
// order and therefore fails when an object references a not-yet-created object
// (a view over a later table, an index/trigger declared ahead of its table) —
// LoadSDL accepts a full target schema whose objects may reference each other
// out of declaration order.
//
// SDL describes a desired end state, so only the constructive subset of DDL is
// allowed: CREATE {DATABASE,TABLE,INDEX,VIEW,FUNCTION,PROCEDURE,TRIGGER,EVENT}
// plus the session statements USE and SET. DML (INSERT/UPDATE/DELETE/SELECT),
// destructive/mutating DDL (DROP, TRUNCATE, RENAME, ALTER) and the imperative
// CREATE TABLE ... AS SELECT are rejected with a clear error.
//
// Pipeline: parse → validate → extractDeps → topoSort → execute (with
// foreign_key_checks disabled and the per-statement current database restored).
//
// Order-tolerance is achieved the way MySQL itself loads a dump:
//   - Foreign keys: validation is disabled for the whole load (foreign_key_checks
//     off), so a forward FK reference never errors; the constraint is still
//     recorded on the table. Because FK creation never requires the parent table
//     to exist, foreign keys impose no ordering constraint and are deliberately
//     left out of the dependency graph — including them would turn a valid
//     circular-FK schema into a false "dependency cycle". This is the analog of
//     PostgreSQL deferring FK constraints to post-creation ALTER TABLE.
//   - Views, indexes, triggers: these require their target table to already
//     exist (createIndex/createTrigger error otherwise; createView silently
//     yields an unresolved query). A dependency-respecting topological sort
//     places every dependent after the objects it references, so the resulting
//     catalog is identical regardless of input statement order.
//
// Session statements are positional: USE selects the current database and SET
// configures the session. The topological sort may freely reorder object
// statements, so rather than rely on the sorted position of USE, LoadSDL records
// the current database at each statement's original position and restores it
// before executing that statement. foreign_key_checks is likewise held off for
// the entire load regardless of any SET inside the SDL.
func LoadSDL(sql string) (*Catalog, error) {
	c := New()

	list, err := mysqlparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) == 0 {
		return c, nil
	}
	stmts := list.Items

	// Validate every statement before mutating any catalog state.
	if err := validateSDL(stmts); err != nil {
		return nil, err
	}

	// Resolve the current database at each statement's ORIGINAL position so that,
	// after the sort reorders statements, every object is created under the
	// database that lexically governed it.
	dbAtIndex := currentDBPerStatement(stmts)

	// Disable FK validation for the whole load so forward references resolve, and
	// restore the normal session default on the way out. The returned catalog
	// then reflects a fully-validated schema for any later mutation.
	c.SetForeignKeyChecks(false)
	defer c.SetForeignKeyChecks(true)

	// Order statements so every dependent runs after the objects it references.
	deps := extractSDLDeps(stmts)
	order, err := topoSortSDL(stmts, deps)
	if err != nil {
		return c, err
	}

	for _, idx := range order {
		// Restore the database current at this statement's original position, and
		// keep FK checks off (a SET inside the SDL must not re-arm validation).
		c.SetCurrentDatabase(dbAtIndex[idx])
		c.SetForeignKeyChecks(false)
		if err := c.processUtility(stmts[idx]); err != nil {
			return c, err
		}
	}

	return c, nil
}

// ---------------------------------------------------------------------------
// Validation
// ---------------------------------------------------------------------------

// validateSDL checks that every statement is allowed in SDL, returning the
// first disallowed statement as an error.
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
	// ---- Allowed constructive DDL ----
	case *nodes.CreateTableStmt:
		// CREATE TABLE ... AS SELECT is imperative: createTable cannot
		// faithfully materialize it as a schema snapshot and silently drops the
		// table. Reject it so the loader never reports success while omitting an
		// object.
		if s.Select != nil {
			return fmt.Errorf("SDL does not allow CREATE TABLE ... AS SELECT statements")
		}
		return nil
	case *nodes.CreateDatabaseStmt,
		*nodes.CreateIndexStmt,
		*nodes.CreateViewStmt,
		*nodes.CreateFunctionStmt, // function or procedure (IsProcedure)
		*nodes.CreateTriggerStmt,
		*nodes.CreateEventStmt:
		return nil

	// ---- Allowed session statements ----
	case *nodes.UseStmt, *nodes.SetStmt:
		return nil

	// ---- Rejected DML ----
	case *nodes.InsertStmt:
		return fmt.Errorf("SDL does not allow INSERT statements")
	case *nodes.UpdateStmt:
		return fmt.Errorf("SDL does not allow UPDATE statements")
	case *nodes.DeleteStmt:
		return fmt.Errorf("SDL does not allow DELETE statements")
	case *nodes.SelectStmt:
		return fmt.Errorf("SDL does not allow SELECT statements")

	// ---- Rejected destructive / mutating DDL ----
	case *nodes.DropTableStmt,
		*nodes.DropDatabaseStmt,
		*nodes.DropIndexStmt,
		*nodes.DropViewStmt,
		*nodes.DropRoutineStmt,
		*nodes.DropTriggerStmt,
		*nodes.DropEventStmt:
		return fmt.Errorf("SDL does not allow DROP statements")
	case *nodes.TruncateStmt:
		return fmt.Errorf("SDL does not allow TRUNCATE statements")
	case *nodes.RenameTableStmt:
		return fmt.Errorf("SDL does not allow RENAME statements")
	case *nodes.AlterTableStmt:
		return fmt.Errorf("SDL does not allow ALTER TABLE statements")
	case *nodes.AlterViewStmt:
		return fmt.Errorf("SDL does not allow ALTER VIEW statements")
	case *nodes.AlterDatabaseStmt:
		return fmt.Errorf("SDL does not allow ALTER DATABASE statements")
	case *nodes.AlterRoutineStmt:
		return fmt.Errorf("SDL does not allow ALTER routine statements")
	case *nodes.AlterEventStmt:
		return fmt.Errorf("SDL does not allow ALTER EVENT statements")

	default:
		return fmt.Errorf("SDL does not allow %T statements", stmt)
	}
}

// ---------------------------------------------------------------------------
// Name resolution
// ---------------------------------------------------------------------------

// objectKey identifies a schema object by (database, name), both lower-cased.
// MySQL has no schemas/namespaces or OIDs: "schema" == database, and table and
// view names share one namespace. An empty database means the current database
// at the point of reference.
type objectKey struct {
	db   string
	name string
}

// resolveTableRef returns the object key for a table/view reference, qualifying
// an unqualified name with the supplied current database.
func resolveTableRef(ref *nodes.TableRef, currentDB string) objectKey {
	if ref == nil {
		return objectKey{}
	}
	db := ref.Schema
	if db == "" {
		db = currentDB
	}
	return objectKey{db: toLower(db), name: toLower(ref.Name)}
}

// currentDBPerStatement returns, for each statement index, the database that is
// current at that point — i.e. the database selected by the most recent USE.
// CREATE DATABASE does not change the current database in MySQL.
func currentDBPerStatement(stmts []nodes.Node) []string {
	out := make([]string, len(stmts))
	current := ""
	for i, stmt := range stmts {
		if use, ok := stmt.(*nodes.UseStmt); ok {
			current = use.Database
		}
		out[i] = current
	}
	return out
}

// ---------------------------------------------------------------------------
// Dependency extraction
// ---------------------------------------------------------------------------

// sdlDep is a dependency edge: statement at index `from` must run after the
// statement at index `to`.
type sdlDep struct {
	from int
	to   int
}

// extractSDLDeps walks each statement's structural references and produces edges
// to the declared objects whose prior existence the catalog requires. Only
// references that resolve to a table/view declared in the same SDL produce an
// edge; references to objects outside the SDL are ignored. Foreign keys are not
// included: with foreign_key_checks off, a FK never requires its parent to exist.
func extractSDLDeps(stmts []nodes.Node) []sdlDep {
	dbAt := currentDBPerStatement(stmts)

	// Map each declared table/view object key to its defining statement index.
	keyToIdx := make(map[objectKey]int, len(stmts))
	for i, stmt := range stmts {
		switch s := stmt.(type) {
		case *nodes.CreateTableStmt:
			if s.Table != nil {
				keyToIdx[resolveTableRef(s.Table, dbAt[i])] = i
			}
		case *nodes.CreateViewStmt:
			if s.Name != nil {
				keyToIdx[resolveTableRef(s.Name, dbAt[i])] = i
			}
		}
	}

	var deps []sdlDep
	for i, stmt := range stmts {
		for _, ref := range sdlRefs(stmt, dbAt[i]) {
			if j, ok := keyToIdx[ref]; ok && j != i {
				deps = append(deps, sdlDep{from: i, to: j})
			}
		}
	}
	return deps
}

// sdlRefs returns the object keys a statement structurally depends on among
// declared objects, qualifying unqualified names with currentDB.
func sdlRefs(stmt nodes.Node, currentDB string) []objectKey {
	var refs []objectKey
	add := func(ref *nodes.TableRef) {
		if ref != nil && ref.Name != "" {
			refs = append(refs, resolveTableRef(ref, currentDB))
		}
	}

	switch s := stmt.(type) {
	case *nodes.CreateTableStmt:
		// CREATE TABLE ... LIKE source table must exist first. (FK references are
		// intentionally excluded — see extractSDLDeps.)
		if s.Like != nil {
			add(s.Like)
		}

	case *nodes.CreateViewStmt:
		// A view depends on every table/view in its SELECT body, excluding names
		// bound by the query's own CTEs.
		for _, ref := range selectTableRefs(s.Select) {
			add(ref)
		}

	case *nodes.CreateIndexStmt:
		// An index requires its target table.
		add(s.Table)

	case *nodes.CreateTriggerStmt:
		// A trigger requires its target table.
		add(s.Table)
	}

	return refs
}

// selectTableRefs returns the base table/view references in a SELECT that denote
// real schema objects (for dependency edges), including those nested in joins,
// subqueries and derived tables. References whose unqualified name is bound by a
// CTE in the top-level query scope are excluded: a CTE is query-local and shadows
// a same-named schema object.
//
// CTE shadowing is applied only for top-level CTEs (the WITH clause on the view
// body and the set-operation branches that share its scope), not CTEs declared in
// nested subqueries — a nested CTE does not shadow an outer reference. Treating a
// nested CTE name as a top-level shadow could drop a real dependency edge; under-
// approximating the shadow set is safe because a stray edge to a non-existent
// object is simply ignored by extractSDLDeps.
func selectTableRefs(sel *nodes.SelectStmt) []*nodes.TableRef {
	if sel == nil {
		return nil
	}
	cteNames := make(map[string]bool)
	collectTopLevelCTENames(sel, cteNames)

	var refs []*nodes.TableRef
	nodes.Inspect(sel, func(n nodes.Node) bool {
		if tr, ok := n.(*nodes.TableRef); ok {
			// Only unqualified references can match a CTE name; a schema-qualified
			// reference always denotes a real object.
			if tr.Schema == "" && cteNames[toLower(tr.Name)] {
				return true
			}
			refs = append(refs, tr)
		}
		return true
	})
	return refs
}

// collectTopLevelCTENames adds the CTE names declared in the top-level scope of a
// query into names. The top-level scope spans the query's own WITH clause plus the
// set-operation branches and parenthesized wrapper that share that scope; it does
// not descend into FROM-clause subqueries or CTE bodies (those open their own
// scopes).
func collectTopLevelCTENames(sel *nodes.SelectStmt, names map[string]bool) {
	if sel == nil {
		return
	}
	for _, cte := range sel.CTEs {
		if cte != nil && cte.Name != "" {
			names[toLower(cte.Name)] = true
		}
	}
	// UNION/INTERSECT/EXCEPT branches and a parenthesized wrapper share the
	// enclosing WITH scope.
	collectTopLevelCTENames(sel.Left, names)
	collectTopLevelCTENames(sel.Right, names)
	collectTopLevelCTENames(sel.ParenSource, names)
}

// ---------------------------------------------------------------------------
// Topological ordering
// ---------------------------------------------------------------------------

// sdlPriority assigns an ordering layer to each statement type; lower runs
// first. Session statements (SET, USE) and CREATE DATABASE share the lowest
// layer so the database exists and the session is configured before any object
// is created. Their relative order is preserved by the original-index tie-break,
// and object placement into the correct database is handled separately via the
// per-statement current-database restore in LoadSDL.
func sdlPriority(stmt nodes.Node) int {
	switch stmt.(type) {
	case *nodes.SetStmt, *nodes.UseStmt, *nodes.CreateDatabaseStmt:
		return 0
	case *nodes.CreateTableStmt:
		return 1
	case *nodes.CreateFunctionStmt:
		// Routines before views/triggers in case a view or trigger body calls a
		// routine (best-effort; routine bodies are not analyzed at load).
		return 2
	case *nodes.CreateViewStmt:
		return 3
	case *nodes.CreateIndexStmt:
		return 4
	case *nodes.CreateTriggerStmt:
		return 5
	case *nodes.CreateEventStmt:
		return 6
	default:
		return 7
	}
}

// topoSortSDL returns statement indices ordered so every dependent runs after
// its dependencies (Kahn's algorithm), breaking ties by (priority, original
// index) for determinism. A genuine dependency cycle (e.g. two views referencing
// each other) is reported as an error.
func topoSortSDL(stmts []nodes.Node, deps []sdlDep) ([]int, error) {
	n := len(stmts)

	priority := make([]int, n)
	for i, stmt := range stmts {
		priority[i] = sdlPriority(stmt)
	}

	// dep.from depends on dep.to → edge to.from must precede from. De-duplicate
	// edges so parallel references don't inflate the in-degree.
	adj := make([][]int, n)
	inDegree := make([]int, n)
	seen := make(map[sdlDep]bool, len(deps))
	for _, d := range deps {
		if seen[d] {
			continue
		}
		seen[d] = true
		adj[d.to] = append(adj[d.to], d.from)
		inDegree[d.from]++
	}

	h := make(sdlHeap, 0, n)
	for i := 0; i < n; i++ {
		if inDegree[i] == 0 {
			h = append(h, sdlHeapEntry{idx: i, pri: priority[i]})
		}
	}
	heap.Init(&h)

	order := make([]int, 0, n)
	for h.Len() > 0 {
		e := heap.Pop(&h).(sdlHeapEntry)
		order = append(order, e.idx)
		for _, next := range adj[e.idx] {
			inDegree[next]--
			if inDegree[next] == 0 {
				heap.Push(&h, sdlHeapEntry{idx: next, pri: priority[next]})
			}
		}
	}

	if len(order) != n {
		return nil, fmt.Errorf("SDL dependency cycle detected")
	}
	return order, nil
}

// sdlHeapEntry is a priority-queue entry for Kahn's algorithm: order by
// (priority ASC, original index ASC).
type sdlHeapEntry struct {
	idx int
	pri int
}

type sdlHeap []sdlHeapEntry

func (h sdlHeap) Len() int      { return len(h) }
func (h sdlHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h sdlHeap) Less(i, j int) bool {
	if h[i].pri != h[j].pri {
		return h[i].pri < h[j].pri
	}
	return h[i].idx < h[j].idx
}
func (h *sdlHeap) Push(x any) {
	*h = append(*h, x.(sdlHeapEntry))
}
func (h *sdlHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	*h = old[:n-1]
	return e
}
