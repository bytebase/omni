package catalog

import (
	"fmt"
	"slices"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
	mariadbparser "github.com/bytebase/omni/mariadb/parser"
)

// generateViewDDL turns the view diff (SchemaDiff.Views) into CREATE OR REPLACE VIEW and DROP
// VIEW IF EXISTS ops. An added view and a changed view emit the same CREATE OR REPLACE statement.
//
// CREATE runs in PhaseMain at priorityView, after table creation, column ALTERs, and routine
// CREATEs, so a view applies against the final tables and any routine its body calls. A view
// that references another view created in the same plan must be created after it; the ops share
// one priority, so viewCreateOrder prefixes the sortName with a topological depth. DROP runs in
// PhasePre, so a dropped and re-created view never collides, and needs no inter-view ordering.
//
// The body is emitted as stored (v.Definition); the diff's canonical comparison absorbs any
// database qualifier the server adds.
func generateViewDDL(_, _ *Catalog, diff *SchemaDiff, _ *Normalizer) []MigrationOp {
	if diff == nil || len(diff.Views) == 0 {
		return nil
	}

	// Keyed by lower-cased (database, name).
	depth := viewCreateOrder(diff.Views)

	var ops []MigrationOp
	for i := range diff.Views {
		entry := &diff.Views[i]
		switch entry.Action {
		case DiffAdd, DiffModify:
			if entry.To == nil {
				continue
			}
			ops = append(ops, createViewOp(entry, depth))
		case DiffDrop:
			if entry.From == nil {
				continue
			}
			ops = append(ops, dropViewOp(entry))
		}
	}

	return ops
}

// createViewOp builds a CREATE OR REPLACE VIEW op. The sortName carries the view's topological
// depth so a referenced view is created before its referrer.
func createViewOp(entry *ViewDiffEntry, depth map[tableKey]int) MigrationOp {
	v := entry.To
	d := depth[tableKey{db: toLower(entry.Database), name: toLower(entry.Name)}]
	return MigrationOp{
		Type:       OpCreateView,
		Database:   entry.Database,
		ObjectName: entry.Name,
		SQL:        formatCreateOrReplaceView(v),
		Phase:      PhaseMain,
		Priority:   priorityView,
		sortName:   viewCreateSortName(d, entry.Database, entry.Name),
	}
}

// dropViewOp builds a DROP VIEW IF EXISTS op in PhasePre.
func dropViewOp(entry *ViewDiffEntry) MigrationOp {
	return MigrationOp{
		Type:       OpDropView,
		Database:   entry.Database,
		ObjectName: entry.Name,
		SQL:        fmt.Sprintf("DROP VIEW IF EXISTS %s", viewIdent(entry.From, entry.Database, entry.Name)),
		Phase:      PhasePre,
		Priority:   priorityView,
		sortName:   viewSortName(entry.Database, entry.Name),
	}
}

// formatCreateOrReplaceView renders a view's CREATE OR REPLACE statement: ALGORITHM and SQL
// SECURITY preamble, database-qualified name, optional explicit column list, AS body, and
// optional WITH CHECK OPTION.
//
// DEFINER is not emitted. The diff ignores it, and `CREATE OR REPLACE DEFINER=<other> VIEW`
// needs elevated privileges that a least-privilege deployment lacks. Omitting it makes the
// definer the current user.
func formatCreateOrReplaceView(v *View) string {
	var b strings.Builder

	// CREATE VIEW looks a bare stored function call up in the session's
	// database (1305 otherwise), and the call has to stay bare: MariaDB stores
	// it as written. Measured on mariadb:11.8.8.
	if callsOwnFunctionsBare(v) {
		fmt.Fprintf(&b, "USE %s;\n", quoteIdent(v.databaseName()))
	}
	b.WriteString("CREATE OR REPLACE")

	algorithm := v.Algorithm
	if algorithm == "" {
		algorithm = "UNDEFINED"
	}
	fmt.Fprintf(&b, " ALGORITHM=%s", strings.ToUpper(algorithm))

	sqlSecurity := v.SqlSecurity
	if sqlSecurity == "" {
		sqlSecurity = "DEFINER"
	}
	fmt.Fprintf(&b, " SQL SECURITY %s", strings.ToUpper(sqlSecurity))

	fmt.Fprintf(&b, " VIEW %s", viewIdent(v, v.databaseName(), v.Name))

	if v.ExplicitColumns && len(v.Columns) > 0 {
		cols := make([]string, len(v.Columns))
		for i, c := range v.Columns {
			cols[i] = migrationQuoteIdent(c)
		}
		fmt.Fprintf(&b, " (%s)", strings.Join(cols, ","))
	}

	b.WriteString(" AS ")
	b.WriteString(qualifyViewRelations(v.Definition, v.databaseName()))

	if v.CheckOption != "" {
		fmt.Fprintf(&b, " WITH %s CHECK OPTION", strings.ToUpper(v.CheckOption))
	}

	return b.String()
}

// viewIdent returns the backtick-quoted database.view identifier using the view's original-case
// name and database, falling back to the diff entry's keys. A view with no database resolves to
// a bare name.
func viewIdent(v *View, database, name string) string {
	dbName := database
	viewName := name
	if v != nil {
		if v.Name != "" {
			viewName = v.Name
		}
		if dn := v.databaseName(); dn != "" {
			dbName = dn
		}
	}
	if dbName == "" {
		return migrationQuoteIdent(viewName)
	}
	return migrationQuoteIdent(dbName) + "." + migrationQuoteIdent(viewName)
}

// databaseName returns the view's database name, or "" if it has no database back-pointer.
func (v *View) databaseName() string {
	if v == nil || v.Database == nil {
		return ""
	}
	return v.Database.Name
}

// viewCreateOrder computes a topological depth for every created or modified view, so a view
// referencing another view in the same batch sorts after it. Depth 0 references no other batch
// view; otherwise depth is one more than the deepest batch view it references. Dependencies come
// from referencedRelations and are matched on the full (database, name) key.
//
// A view cannot reference itself or form a cycle, so the memoized DFS terminates; a visiting
// guard breaks any cycle at depth 0 on malformed input.
func viewCreateOrder(entries []ViewDiffEntry) map[tableKey]int {
	type viewInfo struct {
		key  tableKey
		refs map[tableKey]bool // relations referenced at FROM/JOIN positions
	}
	creates := make([]viewInfo, 0, len(entries))
	idxByKey := make(map[tableKey]int) // identity -> index into creates
	for i := range entries {
		e := &entries[i]
		if e.Action == DiffDrop || e.To == nil {
			continue
		}
		k := tableKey{db: toLower(e.Database), name: toLower(e.Name)}
		idxByKey[k] = len(creates)
		creates = append(creates, viewInfo{key: k, refs: referencedRelations(e.To.Definition, k.db)})
	}

	// Edge i -> j when view i references batch view j.
	refIdx := make([][]int, len(creates))
	for i, ci := range creates {
		for j, cj := range creates {
			if i == j {
				continue
			}
			if ci.refs[cj.key] {
				refIdx[i] = append(refIdx[i], j)
			}
		}
	}

	depth := make(map[tableKey]int, len(creates))
	memo := make([]int, len(creates))
	for i := range memo {
		memo[i] = -1
	}
	visiting := make([]bool, len(creates))
	var compute func(i int) int
	compute = func(i int) int {
		if memo[i] >= 0 {
			return memo[i]
		}
		if visiting[i] {
			// A cycle cannot occur in valid views; break it.
			return 0
		}
		visiting[i] = true
		best := 0
		for _, j := range refIdx[i] {
			if d := compute(j) + 1; d > best {
				best = d
			}
		}
		visiting[i] = false
		memo[i] = best
		return best
	}
	for i, ci := range creates {
		depth[ci.key] = compute(i)
	}
	return depth
}

// referencedRelations returns the lower-cased (database, name) keys a view body references as
// relations, with an unqualified reference resolved against ownDB. It walks the parsed AST, so
// CTE names, aliases, and string literals are excluded. If the body fails to parse or is not a
// SELECT, the empty set is returned and the view carries no ordering constraint.
func referencedRelations(body, ownDB string) map[tableKey]bool {
	out := make(map[tableKey]bool)
	list, err := mariadbparser.Parse(body)
	if err != nil || list == nil || len(list.Items) == 0 {
		return out
	}
	addRef := func(ref *nodes.TableRef) {
		if ref == nil || ref.Name == "" {
			return
		}
		db := ref.Schema
		if db == "" {
			db = ownDB
		}
		out[tableKey{db: toLower(db), name: toLower(ref.Name)}] = true
	}
	switch stmt := list.Items[0].(type) {
	case *nodes.SelectStmt:
		for _, ref := range selectTableRefs(stmt) {
			addRef(ref)
		}
	}
	return out
}

// viewSortName is the stable secondary sort key for a view-level op: lower-cased database.view.
func viewSortName(database, name string) string {
	return strings.ToLower(database) + "." + strings.ToLower(name)
}

// viewCreateSortName prefixes the view sort name with a zero-padded topological depth so a
// referenced view sorts before its referrer.
func viewCreateSortName(depth int, database, name string) string {
	return fmt.Sprintf("%05d.%s", depth, viewSortName(database, name))
}

// qualifyViewRelations qualifies with db the tables and views a stored view
// body names bare, CTEs aside. The catalog stores the body without its own
// database's qualifier, as the server does, while CREATE VIEW resolves a bare
// name against the session's database, not the view's. A body that does not
// parse is returned unchanged.
func qualifyViewRelations(body, db string) string {
	if db == "" || body == "" {
		return body
	}
	const prefix = "CREATE VIEW v AS "
	list, err := mariadbparser.Parse(prefix + body)
	if err != nil || len(list.Items) != 1 {
		return body
	}
	stmt, ok := list.Items[0].(*nodes.CreateViewStmt)
	if !ok {
		return body
	}
	var starts []int
	for _, ref := range selectTableRefs(stmt.Select) {
		start := ref.Loc.Start - len(prefix)
		if ref.Schema != "" || start < 0 || start >= len(body) {
			continue
		}
		if rest := body[start:]; strings.HasPrefix(rest, quoteIdent(ref.Name)) || hasPrefixFold(rest, ref.Name) {
			starts = append(starts, start)
		}
	}
	slices.Sort(starts)
	starts = slices.Compact(starts)
	var out strings.Builder
	prev := 0
	for _, start := range starts {
		out.WriteString(body[prev:start])
		out.WriteString(quoteIdent(db) + ".")
		prev = start
	}
	out.WriteString(body[prev:])
	return out.String()
}

// callsOwnFunctionsBare reports whether a view's body calls a stored function of
// its own database without qualifying it.
func callsOwnFunctionsBare(v *View) bool {
	if v.Database == nil || len(v.Database.Functions) == 0 {
		return false
	}
	const prefix = "CREATE VIEW v AS "
	list, err := mariadbparser.Parse(prefix + v.Definition)
	if err != nil || len(list.Items) != 1 {
		return false
	}
	found := false
	nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
		if call, ok := n.(*nodes.FuncCallExpr); ok && call.Schema == "" && v.Database.Functions[toLower(call.Name)] != nil {
			found = true
		}
		return !found
	})
	return found
}
