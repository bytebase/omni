package review

import (
	"context"
	"strings"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// scan follows the change statement by statement for the rules that ask
// what a statement does to an object of the target's schema:
// DisallowDropConstraint and RequirePrimaryKey. It simulates nothing. It
// checks each statement against the synced schema directly, and keeps
// track of what earlier statements may have changed: a name whose meaning
// a statement may have changed is not looked up again, and a statement
// whose effect the scan cannot bound (procedural code, a rollback, a
// cascading drop of anything but a relation) ends it. Every way the scan
// can be unsure leads to no finding.
type scan struct {
	on    map[review.Rule]bool
	r     *reporter
	index *schemaIndex // nil without a schema

	// path is the search path as last set, "$user" unexpanded, or nil
	// when unknown. user is the current user and session the session
	// user, "" when unknown.
	path    []string
	user    string
	session string

	// touched lists the relation names a statement created, dropped,
	// renamed, or moved, by (schema, name), with "" for a schema the
	// statement left to the search path: the synced schema no longer says
	// what they name. touchedName lists the same names alone.
	touched     map[[2]string]bool
	touchedName map[string]bool
	// unsettled lists the tables, by name, whose columns, constraints, or
	// inheritance a statement changed in a way the scan does not follow.
	unsettled map[string]bool
	// constraints lists the (table, constraint) names a statement added,
	// dropped, or renamed, and renamedKeys the names an index rename may
	// have given to or taken from a constraint, on whichever table.
	constraints map[[2]string]bool
	renamedKeys map[string]bool
	// dropped lists the (schema, table, constraint) a DROP CONSTRAINT of a
	// resolved table removed.
	dropped map[[3]string]bool
	// columns lists the (table, column) names a statement added, and
	// renamedFrom and renamedTo the old and new names of renamed columns.
	columns     map[[2]string]bool
	renamedFrom map[[2]string]bool
	renamedTo   map[[2]string]bool
	// schemas lists the schemas a statement created or dropped.
	schemas map[string]bool
	// referencedByChange and readByChange list the relations, by name, a
	// foreign key or a view the change created depends on, which the synced
	// schema does not record.
	referencedByChange map[string]bool
	readByChange       map[string]bool

	// pending are the tables RequirePrimaryKey reports at the end unless a
	// later statement settles them.
	pending []*pendingTable
	stopped bool
}

// scanTarget scans the change on one target. RequirePrimaryKey is about
// the end of the change, so it reports only when the scan reaches it.
func scanTarget(ctx context.Context, stmts []statement, on map[review.Rule]bool, target review.Target, r *reporter) error {
	s := &scan{
		on:          on,
		r:           r,
		user:        target.SessionUser,
		session:     target.SessionUser,
		touched:     make(map[[2]string]bool),
		touchedName: make(map[string]bool),
		unsettled:   make(map[string]bool),
		constraints: make(map[[2]string]bool),
		renamedKeys: make(map[string]bool),
		dropped:     make(map[[3]string]bool),
		columns:     make(map[[2]string]bool),
		renamedFrom: make(map[[2]string]bool),
		renamedTo:   make(map[[2]string]bool),
		schemas:     make(map[string]bool),

		referencedByChange: make(map[string]bool),
		readByChange:       make(map[string]bool),
	}
	if target.Schema != nil {
		s.index = newSchemaIndex(target.Schema)
		// A schema synced without its search path leaves it unknown.
		if path := parseSearchPath(target.Schema.GetSearchPath()); len(path) > 0 {
			s.path = path
		}
	}
	for i := range stmts {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.statement(&stmts[i])
		if s.stopped {
			return nil
		}
	}
	for _, p := range s.pending {
		if !p.settled {
			r.report(review.RequirePrimaryKey, p.statement, p.rng, p.message)
		}
	}
	return nil
}

// stop ends the scan: nothing after this statement is known, and no table
// is known to end the change without a key.
func (s *scan) stop() {
	s.stopped = true
}

func (s *scan) statement(st *statement) {
	switch v := st.node.(type) {
	case *ast.AlterTableStmt:
		s.alterTable(st, v)
	case *ast.CreateStmt:
		s.createTable(st, v)
		s.touch(v.Relation)
		if v.TableElts != nil {
			for _, item := range v.TableElts.Items {
				switch e := item.(type) {
				case *ast.ColumnDef:
					s.references(constraintsOf(e.Constraints)...)
				case *ast.Constraint:
					s.references(e)
				}
			}
		}
	case *ast.CreateForeignTableStmt:
		s.touch(v.Base.Relation)
	case *ast.CreateTableAsStmt:
		if v.Into != nil {
			s.touch(v.Into.Rel)
		}
		if v.Objtype == ast.OBJECT_MATVIEW {
			s.reads(v.Query)
		}
	case *ast.ViewStmt:
		s.touch(v.View)
		s.reads(v.Query)
	case *ast.CreateSeqStmt:
		s.touch(v.Sequence)
	case *ast.CompositeTypeStmt:
		s.touch(v.Typevar)
	case *ast.IndexStmt:
		// An unnamed index takes a generated name the server keeps clear of
		// every relation already in its schema.
		if v.Idxname != "" && v.Relation != nil {
			s.touchName(v.Relation.Schemaname, v.Idxname)
		}
	case *ast.SelectStmt:
		switch {
		case v.IntoClause != nil:
			s.touch(v.IntoClause.Rel)
		case selectsInto(v):
			s.stop()
		default:
			s.data(v)
		}
	case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
		s.data(v)
	case *ast.ExplainStmt:
		if !explainAnalyzes(v) {
			return
		}
		switch q := v.Query.(type) {
		case *ast.CreateTableAsStmt:
			if q.Into != nil {
				s.touch(q.Into.Rel)
			}
		case *ast.SelectStmt:
			if q.IntoClause != nil || selectsInto(q) {
				s.stop()
				return
			}
			s.data(q)
		default:
			s.data(q)
		}
	case *ast.RenameStmt:
		s.rename(v)
	case *ast.AlterObjectSchemaStmt:
		s.setSchema(v)
	case *ast.DropStmt:
		s.drop(v)
	case *ast.CreateSchemaStmt:
		if v.SchemaElts != nil && len(v.SchemaElts.Items) > 0 {
			s.stop()
			return
		}
		// A schema named after its owner leaves no name to record; it holds
		// only what the change creates in it.
		if v.Schemaname != "" {
			s.schemas[v.Schemaname] = true
		}
	case *ast.DoStmt, *ast.CallStmt:
		// Procedural code may run any DDL.
		s.stop()
	case *ast.DropOwnedStmt, *ast.ImportForeignSchemaStmt, *ast.AlterExtensionStmt, *ast.DiscardStmt, *ast.LoadStmt:
		// What these drop, create, or reset is not in the statement.
		s.stop()
	case *ast.CreateEventTrigStmt, *ast.AlterEventTrigStmt:
		// An event trigger runs code on every later DDL statement.
		s.stop()
	case *ast.TransactionStmt:
		switch v.Kind {
		case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE,
			ast.TRANS_STMT_COMMIT_PREPARED, ast.TRANS_STMT_ROLLBACK_PREPARED:
			s.stop()
		}
	case *ast.VariableSetStmt:
		s.setting(v)
	}
}

// references records the tables new foreign keys reference.
func (s *scan) references(constraints ...*ast.Constraint) {
	for _, c := range constraints {
		if c.Contype == ast.CONSTR_FOREIGN && c.Pktable != nil {
			s.referencedByChange[c.Pktable.Relname] = true
		}
	}
}

// reads records the relations a new view's query names.
func (s *scan) reads(query ast.Node) {
	ast.Inspect(query, func(n ast.Node) bool {
		if rv, ok := n.(*ast.RangeVar); ok {
			s.readByChange[rv.Relname] = true
		}
		return true
	})
}

// touch records a relation name a statement creates.
func (s *scan) touch(rv *ast.RangeVar) {
	if rv != nil {
		s.touchName(rv.Schemaname, rv.Relname)
	}
}

// touchName records a relation name a statement changed, in schema, or
// in whichever schema the search path finds when schema is "".
func (s *scan) touchName(schema, name string) {
	s.touched[[2]string{schema, name}] = true
	s.touchedName[name] = true
}

// isTouched reports whether a statement may have changed what the name
// means: a qualified name when a statement changed it in that schema or
// in an unknown one, an unqualified name when a statement changed it in
// any schema, which may come first on the search path.
func (s *scan) isTouched(rv *ast.RangeVar) bool {
	if rv.Schemaname == "" {
		return s.touchedName[rv.Relname]
	}
	return s.touched[[2]string{rv.Schemaname, rv.Relname}] || s.touched[[2]string{"", rv.Relname}]
}

// lookup resolves a relation name against the synced schema. ok is false
// when the name may no longer mean what the schema says, or names nothing
// the schema holds. An unqualified name needs a known search path, and a
// name beginning with pg_ may be a system catalog, which pg_catalog finds
// first.
func (s *scan) lookup(rv *ast.RangeVar) (schema string, kind relationKind, ok bool) {
	if s.index == nil || rv == nil || s.isTouched(rv) || s.schemas[rv.Schemaname] {
		return "", 0, false
	}
	if rv.Catalogname != "" && rv.Catalogname != s.index.database {
		return "", 0, false
	}
	if rv.Schemaname != "" {
		ns := s.index.schemas[rv.Schemaname]
		if ns == nil {
			return "", 0, false
		}
		kind, ok := ns.relations[rv.Relname]
		return rv.Schemaname, kind, ok
	}
	if strings.HasPrefix(rv.Relname, "pg_") {
		return "", 0, false
	}
	path, ok := s.searchPath()
	if !ok {
		return "", 0, false
	}
	for _, name := range path {
		if ns := s.index.schemas[name]; ns != nil {
			if kind, ok := ns.relations[rv.Relname]; ok {
				return name, kind, true
			}
			continue
		}
		// A system schema the snapshot leaves out may hold the name.
		if name == "information_schema" {
			return "", 0, false
		}
	}
	return "", 0, false
}

// searchPath expands "$user" in the current search path. Without a known
// current user, "$user" names no schema only when every schema of the
// database is listed by name.
func (s *scan) searchPath() ([]string, bool) {
	if s.path == nil {
		return nil, false
	}
	path := make([]string, 0, len(s.path))
	for _, name := range s.path {
		switch {
		case name != "$user":
			path = append(path, name)
		case s.user != "":
			path = append(path, s.user)
		case s.index.unlisted(s.path):
			return nil, false
		}
	}
	return path, true
}

// table resolves the table an ALTER TABLE acts on, with its constraints
// as the synced schema lists them, when the change has not altered it in
// a way the scan does not follow. A partition is not a table here: its
// keys come from its parent.
func (s *scan) table(rv *ast.RangeVar) (tableRef, bool) {
	schema, kind, ok := s.lookup(rv)
	if !ok || kind != kindTable || s.unsettled[rv.Relname] {
		return tableRef{}, false
	}
	return tableRef{schema, rv.Relname}, true
}

// setting follows SET and RESET. The scan follows SET search_path to a
// list of plain names and SET ROLE. A transaction-local setting reverts at
// commit, and RESET returns to the session's default, which the target
// does not record, so after those the path or the user is unknown.
func (s *scan) setting(v *ast.VariableSetStmt) {
	if v.Kind == ast.VAR_RESET_ALL {
		s.path = nil
		return
	}
	switch v.Name {
	case "search_path":
		if v.Kind != ast.VAR_SET_VALUE || v.IsLocal || !plainNames(v.Args) {
			s.path = nil
			return
		}
		s.path = nil
		for _, item := range v.Args.Items {
			name, _ := constString(item)
			s.path = append(s.path, name)
		}
	case "role":
		switch {
		case v.IsLocal:
			s.user = ""
		case v.Kind == ast.VAR_SET_DEFAULT || v.Kind == ast.VAR_RESET:
			s.user = s.session
		case v.Kind == ast.VAR_SET_VALUE && v.Args != nil && len(v.Args.Items) == 1:
			name, _ := constString(v.Args.Items[0])
			if name == "none" {
				name = s.session
			}
			s.user = name
		default:
			s.user = ""
		}
	case "session_authorization":
		s.user, s.session = "", ""
	}
}

// plainNames reports whether every SET search_path value is a name the
// setting reads back unchanged: lower-case letters, digits, underscores,
// and dollar signs. A value such as 'MySchema' or 'a, b' is read as a list
// of identifiers, which the scan does not redo.
func plainNames(args *ast.List) bool {
	if args == nil {
		return false
	}
	for _, item := range args.Items {
		name, ok := constString(item)
		if !ok || name == "" {
			return false
		}
		for i := 0; i < len(name); i++ {
			c := name[i]
			if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_' || c == '$') {
				return false
			}
		}
	}
	return true
}

// data follows a statement that changes data, not the schema. A call of
// set_config on search_path changes how later names resolve: the form
// pg_dump writes, a lone SELECT of it with constant arguments, sets the
// path, and any other call leaves it unknown.
func (s *scan) data(n ast.Node) {
	calls := searchPathCalls(n)
	if len(calls) == 0 {
		return
	}
	s.path = nil
	if path, ok := dumpSearchPath(n, calls); ok {
		// An empty path is a known path that finds nothing.
		s.path = append([]string{}, path...)
	}
}

// searchPathCalls lists the set_config calls that may set search_path:
// every call whose setting name is not a constant naming another setting.
func searchPathCalls(n ast.Node) []*ast.FuncCall {
	var calls []*ast.FuncCall
	ast.Inspect(n, func(n ast.Node) bool {
		fc, ok := n.(*ast.FuncCall)
		if !ok {
			return true
		}
		parts := nameParts(fc.Funcname)
		builtin := len(parts) == 1 || len(parts) == 2 && parts[0] == "pg_catalog"
		if !builtin || parts[len(parts)-1] != "set_config" {
			return true
		}
		if name, ok := constString(firstArg(fc)); ok && !strings.EqualFold(name, "search_path") {
			return true
		}
		calls = append(calls, fc)
		return true
	})
	return calls
}

// dumpSearchPath reads SELECT pg_catalog.set_config('search_path', '...',
// false), the only statement in n, and returns the path it sets.
func dumpSearchPath(n ast.Node, calls []*ast.FuncCall) ([]string, bool) {
	sel, ok := n.(*ast.SelectStmt)
	if !ok || len(calls) != 1 || sel.FromClause != nil || sel.WhereClause != nil || sel.WithClause != nil ||
		sel.TargetList == nil || len(sel.TargetList.Items) != 1 {
		return nil, false
	}
	target, ok := sel.TargetList.Items[0].(*ast.ResTarget)
	if !ok || target.Val != calls[0] || calls[0].Args == nil || len(calls[0].Args.Items) != 3 {
		return nil, false
	}
	value, ok := constString(calls[0].Args.Items[1])
	if !ok {
		return nil, false
	}
	local, ok := calls[0].Args.Items[2].(*ast.A_Const)
	if !ok || local.Isnull {
		return nil, false
	}
	if b, ok := local.Val.(*ast.Boolean); !ok || b.Boolval {
		return nil, false
	}
	return parseSearchPath(value), true
}

func firstArg(fc *ast.FuncCall) ast.Node {
	if fc.Args == nil || len(fc.Args.Items) == 0 {
		return nil
	}
	return fc.Args.Items[0]
}

// constString returns the value of a string constant.
func constString(n ast.Node) (string, bool) {
	c, ok := n.(*ast.A_Const)
	if !ok || c.Isnull {
		return "", false
	}
	str, ok := c.Val.(*ast.String)
	if !ok {
		return "", false
	}
	return str.Str, true
}

// selectsInto reports whether a branch of a set operation carries INTO.
func selectsInto(s *ast.SelectStmt) bool {
	if s == nil {
		return false
	}
	return s.IntoClause != nil || selectsInto(s.Larg) || selectsInto(s.Rarg)
}

// isTemp reports whether a relation is created in the temporary schema.
func isTemp(rv *ast.RangeVar) bool {
	return rv != nil && (rv.Relpersistence == 't' || rv.Schemaname == "pg_temp" || strings.HasPrefix(rv.Schemaname, "pg_temp_"))
}

// isRelationKind reports whether an object type lives in the relation
// namespace.
func isRelationKind(t ast.ObjectType) bool {
	switch t {
	case ast.OBJECT_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_SEQUENCE, ast.OBJECT_INDEX:
		return true
	}
	return false
}

// rename follows RENAME. A relation renamed is touched under both names,
// and a pending table it may be keeps the new name too. An index rename
// renames the constraint the index backs, if any.
func (s *scan) rename(v *ast.RenameStmt) {
	switch {
	case isRelationKind(v.RenameType) && v.Relation != nil:
		s.touchName(v.Relation.Schemaname, v.Relation.Relname)
		s.touchName(v.Relation.Schemaname, v.Newname)
		for _, p := range s.matching(v.Relation) {
			p.names[v.Newname] = true
		}
		if v.RenameType == ast.OBJECT_INDEX {
			s.renamedKeys[v.Relation.Relname] = true
			s.renamedKeys[v.Newname] = true
		}
	case v.RenameType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			schema := ""
			if len(parts) >= 2 {
				schema = parts[len(parts)-2]
			}
			s.touchName(schema, parts[len(parts)-1])
			s.touchName(schema, v.Newname)
		}
	case v.RenameType == ast.OBJECT_TABCONSTRAINT && v.Relation != nil:
		s.constraints[[2]string{v.Relation.Relname, v.Subname}] = true
		s.constraints[[2]string{v.Relation.Relname, v.Newname}] = true
	case v.RenameType == ast.OBJECT_COLUMN && v.Relation != nil:
		s.renamedFrom[[2]string{v.Relation.Relname, v.Subname}] = true
		s.renamedTo[[2]string{v.Relation.Relname, v.Newname}] = true
	case v.RenameType == ast.OBJECT_SCHEMA:
		s.stop()
	}
}

// setSchema follows SET SCHEMA. A pending table it may be can then be in
// any schema.
func (s *scan) setSchema(v *ast.AlterObjectSchemaStmt) {
	switch {
	case isRelationKind(v.ObjectType) && v.Relation != nil:
		s.touch(v.Relation)
		s.touchName(v.Newschema, v.Relation.Relname)
		for _, p := range s.matching(v.Relation) {
			p.schema = ""
		}
	case v.ObjectType == ast.OBJECT_TYPE:
		if parts := nameParts(listOf(v.Object)); len(parts) > 0 {
			s.touchName("", parts[len(parts)-1])
		}
	case v.ObjectType == ast.OBJECT_EXTENSION:
		s.stop()
	}
}

// drop follows DROP. A cascading drop of a relation also takes the
// relations that go with it (its partitions, indexes, owned sequences,
// and the views that read it), the foreign keys that reference it, and
// any inheritance child, which the synced schema does not record; a
// cascading drop of anything else may take columns of any table, which
// ends the scan. DROP EXTENSION drops the extension's objects.
func (s *scan) drop(v *ast.DropStmt) {
	kind := ast.ObjectType(v.RemoveType)
	cascade := v.Behavior == int(ast.DROP_CASCADE)
	if kind == ast.OBJECT_EXTENSION || cascade && !isRelationKind(kind) {
		s.stop()
		return
	}
	if kind == ast.OBJECT_SCHEMA && v.Objects != nil {
		// Without CASCADE the schema must be empty; its name may then be
		// created again.
		for _, obj := range v.Objects.Items {
			if parts := nameParts(listOf(obj)); len(parts) == 1 {
				s.schemas[parts[0]] = true
			}
		}
		return
	}
	if !isRelationKind(kind) && kind != ast.OBJECT_TYPE || v.Objects == nil {
		return
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) >= 2 {
			rv.Schemaname = parts[len(parts)-2]
		}
		s.touch(rv)
		if kind == ast.OBJECT_INDEX {
			s.renamedKeys[rv.Relname] = true
		}
		for _, p := range s.matching(rv) {
			p.settled = true
		}
		if cascade {
			s.dropDependents(rv)
		}
	}
	if cascade {
		// An inheritance child of a dropped table goes with it.
		for _, p := range s.pending {
			if p.inherits || p.existed {
				p.settled = true
			}
		}
	}
}

// dropDependents touches what a cascading drop of a relation takes with
// it, as far as the synced schema records it, in every schema the name
// may be in.
func (s *scan) dropDependents(rv *ast.RangeVar) {
	s.dropReferences(rv.Relname)
	if s.index == nil {
		return
	}
	for name, ns := range s.index.schemas {
		if rv.Schemaname != "" && rv.Schemaname != name {
			continue
		}
		for _, d := range ns.dependents[rv.Relname] {
			s.touchName(name, d)
		}
		for _, view := range s.index.readers[tableRef{name, rv.Relname}] {
			s.touchName(name, view)
		}
	}
}

// dropReferences records that the foreign keys referencing a table, by
// name in any schema, may be gone.
func (s *scan) dropReferences(table string) {
	if s.index == nil {
		return
	}
	for _, fk := range s.index.foreignKeys {
		if fk.referenced.table == table {
			s.constraints[[2]string{fk.owner.table, fk.name}] = true
		}
	}
}

// constraintKnown reports whether a table's constraint of that name is
// still the one the synced schema lists.
func (s *scan) constraintKnown(table, name string) bool {
	return !s.constraints[[2]string{table, name}] && !s.renamedKeys[name]
}
