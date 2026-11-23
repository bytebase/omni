package catalog

import (
	"container/heap"
	"fmt"
	"maps"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

// LoadSDL parses a declarative MariaDB schema (SDL) and loads it into a new
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
// Order-tolerance works the way a dump loads:
//   - Foreign keys: validation is disabled for the whole load (foreign_key_checks
//     off), so a forward FK reference never errors; the constraint is still
//     recorded on the table. Because FK creation never requires the parent table
//     to exist, foreign keys impose no ordering constraint and are deliberately
//     left out of the dependency graph — including them would turn a valid
//     circular-FK schema into a false "dependency cycle".
//   - Views, indexes, triggers: these require their target table to already
//     exist (createIndex/createTrigger error otherwise; createView silently
//     yields an unresolved query). A dependency-respecting topological sort
//     places every dependent after the objects it references, so the resulting
//     catalog is identical regardless of input statement order.
//
// Session statements are positional: USE selects the current database and SET
// configures the session. The topological sort may freely reorder object
// statements, so rather than rely on the sorted position of USE, LoadSDL records
// the current database and session at each statement's original position and
// restores them before executing that statement. foreign_key_checks is likewise held off for
// the entire load regardless of any SET inside the SDL.
func LoadSDL(sql string) (*Catalog, error) {
	return LoadSDLWithVersion(sql, MariaDB118)
}

// LoadSDLWithVersion is LoadSDL with an explicit version recorded on the returned
// catalog's session, the version Diff and GenerateMigration compare against. It
// does not change how the SDL is parsed or stored. MariaDB118 is the only version,
// and LoadSDL uses it.
//
// The session's explicit_defaults_for_timestamp is seeded from the version's
// default; a `SET` inside the SDL runs after this and wins (exec.go).
func LoadSDLWithVersion(sql string, version Version) (*Catalog, error) {
	c := New()
	c.SetVersion(version)
	c.session.ExplicitDefaultsForTimestamp = NormalizerFor(version).ExplicitDefaultsForTimestamp

	list, err := mariadbparser.Parse(sql)
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
	for _, stmt := range stmts {
		if err := validateDDL(stmt); err != nil {
			return nil, err
		}
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
	c.loadingSDL = true
	defer func() { c.loadingSDL = false }()

	// Order statements so every dependent runs after the objects it references.
	deps := extractSDLDeps(stmts)
	order, err := topoSortSDL(stmts, deps)
	if err != nil {
		return c, err
	}

	sessions, final := sessionPerStatement(c, stmts, sql)
	for _, idx := range order {
		// Restore the database and session current at this statement's original
		// position, and keep FK checks off (a SET inside the SDL must not re-arm
		// validation).
		c.SetCurrentDatabase(dbAtIndex[idx])
		sessions[idx].restore(c)
		c.SetForeignKeyChecks(false)
		if err := c.processUtility(stmts[idx], sql); err != nil {
			return c, err
		}
	}
	c.SetCurrentDatabase(dbAtIndex[len(stmts)-1])
	final.restore(c)
	return c, nil
}

// sdlSession is the session a statement of an SDL text runs under.
type sdlSession struct {
	session       SessionState
	charsetClient string
	collationConn string
}

func (s sdlSession) restore(c *Catalog) {
	c.session, c.charsetClient, c.collationConn = s.session, s.charsetClient, s.collationConn
}

// sessionPerStatement runs the SDL's SET statements in their written order on
// a copy of c's session, and returns the session each statement starts under
// and the one the text ends with: a SET governs the objects written after it,
// wherever the sort puts them.
func sessionPerStatement(c *Catalog, stmts []nodes.Node, sql string) ([]sdlSession, sdlSession) {
	scratch := New()
	current := sdlSession{c.session, c.charsetClient, c.collationConn}
	sessions := make([]sdlSession, len(stmts))
	for i, stmt := range stmts {
		sessions[i] = current
		if set, ok := stmt.(*nodes.SetStmt); ok {
			current.restore(scratch)
			_ = scratch.processUtility(set, sql)
			current = sdlSession{scratch.session, scratch.charsetClient, scratch.collationConn}
		}
	}
	return sessions, current
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
		*nodes.CreateEventStmt,
		*nodes.CreateSequenceStmt:
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
		*nodes.DropEventStmt,
		*nodes.DropSequenceStmt:
		return fmt.Errorf("SDL does not allow DROP statements")
	case *nodes.AlterSequenceStmt:
		return fmt.Errorf("SDL does not allow ALTER SEQUENCE statements")
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
// In MariaDB the schema is the database, and table and view names share one
// namespace. An empty database means the current database at the point of
// reference.
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
// CREATE DATABASE does not change the current database.
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

	// The indexes declared apart from their table, so that CREATE TABLE ...
	// LIKE copies them too.
	indexesOf := make(map[objectKey][]int)
	for i, stmt := range stmts {
		if s, ok := stmt.(*nodes.CreateIndexStmt); ok && s.Table != nil {
			key := resolveTableRef(s.Table, dbAt[i])
			indexesOf[key] = append(indexesOf[key], i)
		}
	}

	var deps []sdlDep
	for i, stmt := range stmts {
		for _, ref := range sdlRefs(stmt, dbAt[i]) {
			if j, ok := keyToIdx[ref]; ok && j != i {
				deps = append(deps, sdlDep{from: i, to: j})
			}
		}
		if s, ok := stmt.(*nodes.CreateTableStmt); ok && s.Like != nil {
			for _, j := range indexesOf[resolveTableRef(s.Like, dbAt[i])] {
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
// subqueries and derived tables. An unqualified reference bound by a CTE in its
// own query scope or an enclosing one is excluded: a CTE is query-local and
// shadows a same-named schema object.
func selectTableRefs(sel *nodes.SelectStmt) []*nodes.TableRef {
	var refs []*nodes.TableRef
	collectSelectTableRefs(sel, nil, &refs)
	return refs
}

// collectSelectTableRefs adds to refs the references of sel not bound by a CTE
// in scope, the CTE names of the enclosing queries. sel's CTEs are visible in
// its body and in its subqueries; a CTE's own body sees the CTEs declared
// before it, or under WITH RECURSIVE all of them.
func collectSelectTableRefs(sel *nodes.SelectStmt, scope map[string]bool, refs *[]*nodes.TableRef) {
	if sel == nil {
		return
	}
	inner := maps.Clone(scope)
	if inner == nil {
		inner = make(map[string]bool)
	}
	// A query's own WITH binds in it and below it; a parenthesized set
	// operand's WITH binds in that operand alone, which its own call handles.
	for _, cte := range sel.CTEs {
		if cte != nil && cte.Name != "" {
			inner[toLower(cte.Name)] = true
		}
	}
	nodes.Inspect(sel, func(n nodes.Node) bool {
		switch n := n.(type) {
		case *nodes.SelectStmt:
			if n == sel {
				return true
			}
			collectSelectTableRefs(n, inner, refs)
			return false
		case *nodes.CommonTableExpr:
			cteScope := maps.Clone(scope)
			if cteScope == nil {
				cteScope = make(map[string]bool)
			}
			for _, cte := range sel.CTEs {
				if cte == n && !n.Recursive {
					break
				}
				cteScope[toLower(cte.Name)] = true
			}
			collectSelectTableRefs(n.Select, cteScope, refs)
			return false
		case *nodes.TableRef:
			// A schema-qualified reference always denotes a real object.
			if n.Schema != "" || !inner[toLower(n.Name)] {
				*refs = append(*refs, n)
			}
		}
		return true
	})
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
	case *nodes.CreateSequenceStmt:
		// Before tables: a default may read a sequence.
		return 1
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
