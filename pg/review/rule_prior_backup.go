package review

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what makes Bytebase's prior backup of a PostgreSQL
// change fail, or copy other rows than the change touches, and only that.
// A failed backup fails the task before any of the change runs. The
// backup works like this (bytebase
// backend/runner/taskrun/database_migrate_executor.go, backupData, and
// backend/plugin/parser/pg/backup.go, TransformDMLToSelect):
//
//   - It runs before the change, on the schema Bytebase last synced.
//   - It refuses a change longer than the size limit.
//   - It copies the rows each top-level UPDATE and DELETE would touch,
//     one table at a time, into a table it creates in the bbdataarchive
//     schema of the same database, which it does not create. A change
//     with no top-level UPDATE or DELETE backs up nothing.
//   - It names the table each statement touches as written when the SQL
//     qualifies it, and otherwise looks it up in the synced schema along
//     the synced search path, or along the last SET search_path before
//     the statement; a name it cannot find fails the backup.
//   - It refuses a table touched by both UPDATE and DELETE, and a
//     statement whose WITH clause holds an INSERT, UPDATE, DELETE, or
//     MERGE, before it runs anything.
//   - It reads the rows with one query per table, which copies, as
//     written, the WITH clause of the table's first statement and each
//     statement's table, FROM or USING list, and WHERE clause, selects
//     through the first statement's alias, and runs with the change's
//     first SET role or search_path before it.
//
// DDL in the change is ignored by the backup, so it is not a finding by
// itself. It matters when the backup's query names what only the change
// makes, which does not exist yet when the query runs: the table, a
// relation of the FROM or USING list, a subquery, or the WITH clause, a
// column the change adds or renames to, or a routine named with its
// schema. The query also fails on a later statement that does not use
// the first one's alias, reads a WITH query named like the table in its
// place, reads another table, or none, when the search path it runs
// under finds the statement's table elsewhere than Bytebase grouped it,
// cannot take WHERE CURRENT OF, and leaves out a later statement's own
// WITH queries. The SET ROLE it runs first fails when the role is one the
// change creates. And the backup copies another table than a statement
// changes when the change's own statements find the table elsewhere along
// the search path than Bytebase does: after a RESET, SET ... TO DEFAULT, or
// set_config of search_path, under $user, the role or the session user
// the change sets, or where the change creates, drops, renames, or moves
// a relation of the name earlier on the path, renames a schema of it, or
// creates a temporary relation; or when the name holds another relation
// by the time the statement runs, which the change creates, renames, or
// moves there in place of the one the backup reads.
//
// A name the synced schema cannot resolve, and the change does not
// create, is not reported. The schema the review sees may predate an
// earlier spec of the same plan, which may create the table; Bytebase
// syncs again before this spec's backup, which then finds it. A name no
// spec creates fails the UPDATE or DELETE itself, which is WalkThrough's
// to report. Neither is what the backup's query may mean otherwise: an
// unqualified function or a type, which the system catalogs may hold, a
// routine call an overload of another signature may take, a column a
// subquery names that its own FROM list may hold, a column whose type the
// change alters, or a WITH query of the first statement standing for a
// later one's own WITH query or relation, which reads rows that may or
// may not be the statement's.

// backupSchema is where Bytebase writes a PostgreSQL prior backup.
const backupSchema = "bbdataarchive"

// checkPriorBackupStatements reports the top-level UPDATE and DELETE
// statements the backup refuses whatever the target: one whose WITH
// clause holds a query other than a SELECT, which copied into the backup
// would change data before the change does. It needs the SQL alone.
//
// bytebase: backend/plugin/parser/pg/backup.go — rejectDataModifyingCTE
func checkPriorBackupStatements(stmts []statement, r *reporter) {
	for i := range stmts {
		s := &stmts[i]
		switch s.node.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
		default:
			continue
		}
		names := dataModifyingCTEs(withOf(s.node))
		if len(names) == 0 {
			continue
		}
		quoted := make([]string, len(names))
		for j, name := range names {
			quoted[j] = ident(name)
		}
		r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), "prior backup refuses the data-modifying WITH query "+strings.Join(quoted, ", ")+" of this statement")
	}
}

// checkPriorBackupSize reports a change longer than prior backup takes.
// It needs the SQL alone.
func checkPriorBackupSize(sql string, change review.Change, r *reporter) {
	if change.MaxBackupSize > 0 && len(sql) > change.MaxBackupSize {
		r.report(review.PriorBackup, -1, review.Range{}, fmt.Sprintf("the change is %d bytes, over the %d-byte limit of prior backup", len(sql), change.MaxBackupSize))
	}
}

// backupTable is a table as the backup keys it: schema and name.
type backupTable struct {
	schema string
	name   string
}

// checkPriorBackupTarget reports the target-dependent reasons the backup
// fails or copies another table. Without a schema nothing is known about
// the target, and nothing is reported. sessionUser is the user the backup connects as.
func checkPriorBackupTarget(stmts []statement, schema *metadata.DatabaseSchemaMetadata, sessionUser string, r *reporter) {
	if schema == nil {
		return
	}
	path := backupSearchPath(schema.GetSearchPath())
	// kinds maps each table to the statement kind the backup saw first,
	// and mixed marks the tables already reported.
	kinds := make(map[backupTable]string)
	mixed := make(map[backupTable]bool)
	made := changeMade(stmts, schema, sessionUser)
	q := backupQueryOf(stmts, schema, sessionUser)
	// engine is the search path the change's own statements run under.
	engine := newEnginePath(schema, sessionUser)
	// undoable is the last rollback of the change: what a statement
	// before it changes, it may undo.
	undoable := -1
	for i := range stmts {
		if v, ok := stmts[i].node.(*ast.TransactionStmt); ok {
			switch v.Kind {
			case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE:
				undoable = i
			}
		}
	}
	// firsts is the first statement on each table, whose WITH clause and
	// alias the backup's query of the table carries.
	firsts := make(map[backupTable]ast.Node)
	backsUp := false
	for i := range stmts {
		s := &stmts[i]
		var rv *ast.RangeVar
		var kind string
		switch v := s.node.(type) {
		case *ast.UpdateStmt:
			rv, kind = v.Relation, "UPDATE"
		case *ast.DeleteStmt:
			rv, kind = v.Relation, "DELETE"
		case *ast.VariableSetStmt:
			if p, ok := backupSetSearchPath(v); ok {
				path = p
			}
			engine.set(v)
			continue
		default:
			engine.follow(s.node)
			continue
		}
		if rv == nil {
			continue
		}
		backsUp = true
		// The backup refuses the statement outright, which
		// checkPriorBackupStatements reports.
		if len(dataModifyingCTEs(withOf(s.node))) > 0 {
			continue
		}
		table := backupTable{schema: rv.Schemaname, name: rv.Relname}
		if rv.Schemaname == "" {
			table.schema = searchBackupRelation(schema, path, rv.Relname)
			if table.schema == "" {
				// The backup cannot find the table. Report it when the
				// change creates it: the table then exists only after the
				// backup has run, whatever the target's schema says later.
				if made.relationNames[rv.Relname] {
					reportNotYetCreated(s, rv, r)
				}
				continue
			}
		}
		if rv.Schemaname != "" && (rv.Catalogname == "" || rv.Catalogname == schema.GetName()) && !hasBackupRelation(findSchema(schema, rv.Schemaname), rv.Relname) && (made.relations[table] || made.relations[backupTable{name: rv.Relname}] || engine.holdsRelation(rv)) {
			// An unqualified creation may be in this schema, or a schema
			// rename gives the name; either way the table is not there
			// when the backup runs.
			reportNotYetCreated(s, rv, r)
		} else {
			first, seenFirst := firsts[table]
			if !seenFirst {
				first = s.node
				firsts[table] = first
			}
			if msg := selectionRefused(s.node, table, first, q, made); msg != "" {
				r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), msg)
			} else if at, h, ok := engine.resolve(rv); i > undoable && ok && !h.none() && h != (holding{synced: [2]string{table.schema, rv.Relname}}) {
				// Bytebase follows SET search_path but not RESET, DEFAULT,
				// set_config, $user, or what the change does to the
				// relations of a name: it backs up another table than the
				// statement changes.
				r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), "prior backup backs up "+relation(&ast.RangeVar{Schemaname: table.schema, Relname: rv.Relname})+" for this statement, which changes "+changedRelation(at, h, table))
			}
		}
		first, seen := kinds[table]
		if !seen {
			kinds[table] = kind
			continue
		}
		if first != kind && !mixed[table] {
			mixed[table] = true
			r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), "prior backup cannot back up both UPDATE and DELETE on "+relation(rv))
		}
	}
	if backsUp && findSchema(schema, backupSchema) == nil {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup needs schema "+backupSchema+", which does not exist")
	}
	if backsUp {
		if s, role := prependedRole(stmts); s != nil && madeRoles(stmts)[role] {
			r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup runs this SET ROLE before the change, which creates role "+ident(role))
		}
	}
}

// prependedRole returns the statement Bytebase runs before the backup's
// queries, the change's first SET role or search_path, when it is SET
// ROLE of a role, with the role.
func prependedRole(stmts []statement) (*statement, string) {
	for i := range stmts {
		v, ok := stmts[i].node.(*ast.VariableSetStmt)
		if !ok || !strings.EqualFold(v.Name, "role") && !strings.EqualFold(v.Name, "search_path") {
			continue
		}
		if !strings.EqualFold(v.Name, "role") || v.Kind != ast.VAR_SET_VALUE || v.Args == nil || len(v.Args.Items) != 1 {
			return nil, ""
		}
		role, ok := constString(v.Args.Items[0])
		if !ok || strings.EqualFold(role, "none") {
			return nil, ""
		}
		return &stmts[i], role
	}
	return nil, ""
}

// madeRoles returns the roles the change creates or gives a name, which
// do not exist before it: none it drops, alters, or renames away before
// it makes the name, which the role then may have had.
func madeRoles(stmts []statement) map[string]bool {
	roles := make(map[string]bool)
	before := make(map[string]bool)
	for i := range stmts {
		switch v := stmts[i].node.(type) {
		case *ast.CreateRoleStmt:
			if !before[v.Role] {
				roles[v.Role] = true
			}
		case *ast.RenameStmt:
			if v.RenameType == ast.OBJECT_ROLE {
				before[v.Subname] = true
				if !before[v.Newname] {
					roles[v.Newname] = true
				}
			}
		case *ast.DropRoleStmt:
			if v.Roles != nil {
				for _, item := range v.Roles.Items {
					if r, ok := item.(*ast.RoleSpec); ok {
						before[r.Rolename] = true
					}
				}
			}
		case *ast.AlterRoleStmt:
			if v.Role != nil {
				before[v.Role.Rolename] = true
			}
		}
	}
	return roles
}

func reportNotYetCreated(s *statement, rv *ast.RangeVar, r *reporter) {
	r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), "prior backup runs before the change, when "+relation(rv)+" does not exist yet")
}

// made is what the change makes, which the backup, run before all of
// it, cannot see: the relations it creates or gives a name, as written,
// an unqualified one under an empty schema; the columns it adds or gives
// a name, by the table as written; and the routine names it creates or
// gives a routine. A statement that may make nothing, IF NOT EXISTS or OR
// REPLACE, is left out: what it does depends on what already exists.
type made struct {
	relations     map[backupTable]bool
	relationNames map[string]bool
	columns       map[backupTable]map[string]bool
	functions     map[backupTable]bool
	// copies are the relations whose columns a relation the change makes
	// has, by that relation as written: the source of CREATE TABLE ...
	// LIKE, INHERITS, or PARTITION OF, and the relation a rename or a move
	// gives the name, each as the change's statement finds it: the synced
	// relation it holds, or one the change makes as written.
	copies map[backupTable][]*ast.RangeVar
}

// makes reports whether the change makes the column on a relation, by the
// relation as written: adds it, renames to it, declares it, or copies it
// from a relation that has it before the change, or that the change makes
// it on, following the copies. seen guards against a cycle.
func (m made) makes(t backupTable, column string, q backupQuery, seen map[backupTable]bool) bool {
	if m.columns[t][column] {
		return true
	}
	if seen[t] {
		return false
	}
	seen = withTable(seen, t)
	for _, from := range m.copies[t] {
		if c := sourceOf(from, nil, q).columns; c != nil && c[column] {
			return true
		}
		if m.makes(backupTable{from.Schemaname, from.Relname}, column, q, seen) {
			return true
		}
	}
	return false
}

// withTable returns seen with a table added.
func withTable(seen map[backupTable]bool, t backupTable) map[backupTable]bool {
	out := maps.Clone(seen)
	if out == nil {
		out = make(map[backupTable]bool)
	}
	out[t] = true
	return out
}

func changeMade(stmts []statement, schema *metadata.DatabaseSchemaMetadata, sessionUser string) made {
	m := made{
		relations:     make(map[backupTable]bool),
		relationNames: make(map[string]bool),
		columns:       make(map[backupTable]map[string]bool),
		functions:     make(map[backupTable]bool),
		copies:        make(map[backupTable][]*ast.RangeVar),
	}
	column := func(t backupTable, name string) {
		if m.columns[t] == nil {
			m.columns[t] = make(map[string]bool)
		}
		m.columns[t][name] = true
	}
	// function records a routine name, in the schema written, or "" for
	// whichever an unqualified name means.
	function := func(parts []string, name string) {
		t := backupTable{name: name}
		if len(parts) == 2 {
			t.schema = parts[0]
		}
		m.functions[t] = true
	}
	routine := func(t ast.ObjectType) bool {
		return t == ast.OBJECT_FUNCTION || t == ast.OBJECT_PROCEDURE || t == ast.OBJECT_ROUTINE
	}
	// engine finds the relations each statement names as it runs.
	engine := newEnginePath(schema, sessionUser)
	copied := func(t backupTable, rv *ast.RangeVar) {
		if source := engine.copySource(rv); source != nil {
			m.copies[t] = append(m.copies[t], source)
		}
	}
	for i := range stmts {
		if i > 0 {
			engine.run(stmts[i-1].node)
		}
		for _, t := range createdRelations(stmts[i].node) {
			m.relations[t] = true
			m.relationNames[t.name] = true
		}
		// A table the change creates has the columns it declares, which a
		// table of that name the backup reads may lack.
		for _, c := range createdColumns(stmts[i].node) {
			column(c.table, c.column)
		}
		for _, l := range createdLikes(stmts[i].node) {
			copied(l.table, l.source)
		}
		switch v := stmts[i].node.(type) {
		case *ast.RenameStmt:
			switch v.RenameType {
			case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
				if v.Relation != nil {
					copied(backupTable{v.Relation.Schemaname, v.Newname}, v.Relation)
				}
			}
		case *ast.AlterObjectSchemaStmt:
			switch v.ObjectType {
			case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
				if v.Relation != nil {
					copied(backupTable{v.Newschema, v.Relation.Relname}, v.Relation)
				}
			}
		}
		switch v := stmts[i].node.(type) {
		case *ast.AlterTableStmt:
			if v.Relation == nil || v.Cmds == nil {
				continue
			}
			for _, item := range v.Cmds.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok || ast.AlterTableType(cmd.Subtype) != ast.AT_AddColumn || cmd.Missing_ok {
					continue
				}
				if cd, ok := cmd.Def.(*ast.ColumnDef); ok {
					column(backupTable{v.Relation.Schemaname, v.Relation.Relname}, cd.Colname)
				}
			}
		case *ast.RenameStmt:
			switch {
			case v.RenameType == ast.OBJECT_COLUMN && v.Relation != nil:
				column(backupTable{v.Relation.Schemaname, v.Relation.Relname}, v.Newname)
			case routine(v.RenameType):
				if owa, ok := v.Object.(*ast.ObjectWithArgs); ok {
					function(nameParts(owa.Objname), v.Newname)
				}
			}
		case *ast.AlterObjectSchemaStmt:
			if owa, ok := v.Object.(*ast.ObjectWithArgs); ok && routine(v.ObjectType) {
				if parts := nameParts(owa.Objname); len(parts) > 0 {
					function([]string{v.Newschema, ""}, parts[len(parts)-1])
				}
			}
		case *ast.CreateFunctionStmt:
			if parts := nameParts(v.Funcname); len(parts) > 0 && len(parts) <= 2 && !v.IsOrReplace {
				function(parts, parts[len(parts)-1])
			}
		}
	}
	return m
}

// backupQuery is what the backup's queries see of the target: the synced
// schema, and the search path an unqualified name of them is found along,
// when it is known.
type backupQuery struct {
	schema *metadata.DatabaseSchemaMetadata
	path   []string
	known  bool
	// ctes are the visible WITH queries whose columns are known.
	ctes map[string]backupSource
}

// backupQueryOf returns what the backup's queries see. They run on a
// connection as the session user, after the change's first SET role or
// search_path, so along the search path that SET gives, or the synced one
// with $user the role or the session user.
//
// bytebase: backend/runner/taskrun/database_migrate_executor.go —
// getPrependStatements
func backupQueryOf(stmts []statement, schema *metadata.DatabaseSchemaMetadata, sessionUser string) backupQuery {
	q := backupQuery{schema: schema}
	user := sessionUser
	setting := schema.GetSearchPath()
	for i := range stmts {
		v, ok := stmts[i].node.(*ast.VariableSetStmt)
		if !ok || !strings.EqualFold(v.Name, "role") && !strings.EqualFold(v.Name, "search_path") {
			continue
		}
		if v.IsLocal {
			return q
		}
		// SET role TO DEFAULT and RESET ROLE return to the session user,
		// and the same of search_path, or SET search_path FROM CURRENT,
		// leave the synced path.
		if v.Kind == ast.VAR_SET_DEFAULT || v.Kind == ast.VAR_RESET || v.Kind == ast.VAR_SET_CURRENT {
			break
		}
		if v.Kind != ast.VAR_SET_VALUE || v.Args == nil {
			return q
		}
		var values []string
		for _, item := range v.Args.Items {
			s, ok := constString(item)
			if !ok {
				return q
			}
			values = append(values, s)
		}
		if strings.EqualFold(v.Name, "role") {
			if len(values) != 1 {
				return q
			}
			if !strings.EqualFold(values[0], "none") {
				user = values[0]
			}
		} else {
			q.path, q.known = values, true
			break
		}
		break
	}
	if !q.known {
		names := parseSearchPath(setting)
		if strings.TrimSpace(setting) == "" {
			names = []string{"$user", "public"}
		}
		q.path, q.known = names, true
	}
	for i, name := range q.path {
		if name == "$user" {
			if user == "" {
				return backupQuery{schema: schema}
			}
			q.path[i] = user
		}
	}
	return q
}

// enginePath follows the search path the change's statements run under,
// from the synced one, with $user the session user: SET search_path to
// names, RESET, DEFAULT, or FROM CURRENT, RESET ALL, SET ROLE, and SET
// SESSION AUTHORIZATION. It stops knowing it at anything else that sets
// one of them. It follows as well which relation each name holds as the
// statements run.
type enginePath struct {
	schema  *metadata.DatabaseSchemaMetadata
	initial []string
	path    []string
	// login is the user the change connects as, session the session
	// user, and user the current user, which $user means; each known
	// flag says whether the change has left it known.
	login        string
	session      string
	user         string
	pathKnown    bool
	sessionKnown bool
	userKnown    bool
	// triggers are the target's enabled event triggers.
	triggers []*metadata.EventTriggerMetadata
	// schemas are the schemas earlier statements created.
	schemas map[string]bool
	// holds maps a name earlier statements gave a relation or took one
	// from, by (schema, name), to the relation it then holds; temps are
	// the temporary relations they created, which the path finds first,
	// with the statement, and tempsClearedAt the last that dropped them
	// all.
	holds          map[[2]string]held
	temps          map[string]int
	tempsClearedAt int
	// partitioned are the partitioned tables they created.
	partitioned map[[2]string]bool
	// schemaFrom maps a name a schema rename gave to the synced schema
	// whose relations it holds, "" for none; goneSchemas lists the synced
	// schemas renamed away or dropped, whose relations their names no
	// longer hold.
	schemaFrom  map[string]string
	goneSchemas map[string]bool
	// lost lists the names, and lostSchemas the schemas, an earlier
	// statement changed in a way not followed, with the last such
	// statement; what a later one gives a name holds. A statement that may
	// drop relations it does not name leaves every relation before it
	// unsure: cascadeAt is the last, -1 for none. One that may run any
	// code leaves every name before it unknown: opaqueAt is the last, -1
	// for none. step is the statement followed.
	lost        map[string]int
	lostSchemas map[string]int
	cascadeAt   int
	opaqueAt    int
	step        int
}

// holding is the relation a name holds as the change runs: the synced
// relation of a name before the change, one the change creates, or, the
// zero holding, none.
type holding struct {
	synced [2]string
	made   bool
}

func (h holding) none() bool { return !h.made && h.synced == [2]string{} }

// held is a holding, with the statement that gave it.
type held struct {
	holding
	at int
}

func newEnginePath(schema *metadata.DatabaseSchemaMetadata, sessionUser string) *enginePath {
	initial := parseSearchPath(schema.GetSearchPath())
	if strings.TrimSpace(schema.GetSearchPath()) == "" {
		initial = []string{"$user", "public"}
	}
	return &enginePath{schema: schema, initial: initial, path: initial, login: sessionUser, session: sessionUser, user: sessionUser,
		pathKnown: true, sessionKnown: true, userKnown: true, triggers: enabledTriggers(schema),
		schemas: make(map[string]bool), holds: make(map[[2]string]held), temps: make(map[string]int), tempsClearedAt: -1,
		partitioned: make(map[[2]string]bool), schemaFrom: make(map[string]string), goneSchemas: make(map[string]bool),
		lost: make(map[string]int), lostSchemas: make(map[string]int), cascadeAt: -1, opaqueAt: -1}
}

// opaque notes a statement that may run any code: what every name holds,
// the path, and the users are not known after it.
func (e *enginePath) opaque() {
	e.opaqueAt = e.step
	e.pathKnown, e.sessionKnown, e.userKnown = false, false, false
}

// triggerMayFire reports whether an enabled event trigger of the target
// may fire on the statement, as the scan judges it, counting every DROP
// that may drop.
func triggerMayFire(triggers []*metadata.EventTriggerMetadata, n ast.Node) bool {
	return len(triggers) > 0 && (&scan{triggers: triggers}).triggerFires(n)
}

// userCalls reports whether a statement calls set_config of the role or
// the session user, or of a setting it does not name.
func userCalls(n ast.Node) bool {
	found := false
	ast.Inspect(n, func(x ast.Node) bool {
		fc, ok := x.(*ast.FuncCall)
		if found || !ok {
			return !found
		}
		parts := nameParts(fc.Funcname)
		builtin := len(parts) == 1 || len(parts) == 2 && parts[0] == "pg_catalog"
		if !builtin || parts[len(parts)-1] != "set_config" {
			return true
		}
		name, ok := constString(firstArg(fc))
		found = !ok || strings.EqualFold(name, "role") || strings.EqualFold(name, "session_authorization")
		return !found
	})
	return found
}

// follow notes what a statement other than SET does to the names the
// change's later statements find: set_config of search_path, the form
// pg_dump writes, sets the path, and any other call of it leaves the path
// unknown, as one of the role leaves the users; a relation it creates,
// renames, moves, or drops changes what the names hold. A statement an
// event trigger may fire on, code, extension DDL, and ROLLBACK, which may
// undo what came before, leave every name unknown.
//
// pg: src/backend/catalog/namespace.c — RangeVarGetCreationNamespace
func (e *enginePath) follow(n ast.Node) {
	defer func() { e.step++ }()
	if triggerMayFire(e.triggers, n) {
		defer e.opaque()
	}
	if calls := searchPathCalls(n); len(calls) > 0 {
		if path, ok := dumpSearchPath(n, calls); ok {
			e.path, e.pathKnown = append([]string{}, path...), true
		} else {
			e.pathKnown = false
		}
	}
	if userCalls(n) {
		e.sessionKnown, e.userKnown = false, false
	}
	if v, ok := n.(*ast.CreateSchemaStmt); ok {
		name := v.Schemaname
		if name == "" && v.Authrole != nil {
			name = v.Authrole.Rolename
		}
		e.schemas[name] = true
	}
	switch v := n.(type) {
	case *ast.DoStmt, *ast.CallStmt:
		// Code may set the path or change any relation.
		e.opaque()
		return
	case *ast.TransactionStmt:
		// A rollback undoes what the statements before it did as far as
		// the transaction they ran in reaches, which the change does not
		// say when Bytebase runs it in one.
		switch v.Kind {
		case ast.TRANS_STMT_ROLLBACK, ast.TRANS_STMT_ROLLBACK_TO, ast.TRANS_STMT_PREPARE, ast.TRANS_STMT_COMMIT_PREPARED, ast.TRANS_STMT_ROLLBACK_PREPARED:
			e.opaque()
		}
		return
	case *ast.DropOwnedStmt:
		e.cascadeAt = e.step
		return
	case *ast.CreateExtensionStmt, *ast.AlterExtensionStmt:
		// An extension's script, or an update, may create, drop, or move
		// relations it does not name; it keeps the path and the users.
		e.opaqueAt = e.step
		return
	case *ast.ImportForeignSchemaStmt:
		e.lostSchemas[v.LocalSchema] = e.step
		return
	case *ast.DiscardStmt:
		if v.Target == ast.DISCARD_ALL {
			e.session, e.user, e.path = e.login, e.login, e.initial
			e.pathKnown, e.sessionKnown, e.userKnown = true, true, true
		}
		if v.Target == ast.DISCARD_ALL || v.Target == ast.DISCARD_TEMP {
			clear(e.temps)
			e.tempsClearedAt = e.step
		}
		return
	case *ast.DropStmt:
		e.drop(v)
		return
	case *ast.RenameStmt:
		switch {
		case v.RenameType == ast.OBJECT_SCHEMA:
			e.renameSchema(v.Subname, v.Newname)
		case relationObject(v.RenameType) && v.Relation != nil:
			e.move(v.Relation, "", v.Newname)
		}
		return
	case *ast.AlterObjectSchemaStmt:
		switch {
		case v.ObjectType == ast.OBJECT_EXTENSION:
			// The extension's relations move, which it does not name.
			e.opaqueAt = e.step
		case relationObject(v.ObjectType) && v.Relation != nil:
			e.move(v.Relation, v.Newschema, v.Relation.Relname)
		}
		return
	}
	temps := tempRelations(n)
	for name := range temps {
		e.temps[name] = e.step
	}
	for _, t := range createdRelations(n) {
		if t.schema == "pg_temp" || t.schema == "" && temps[t.name] {
			continue
		}
		if t.schema == "" {
			t.schema = e.creationSchema()
		}
		if t.schema == "" {
			e.lost[t.name] = e.step
			continue
		}
		e.give([2]string{t.schema, t.name}, holding{made: true})
	}
	if v, ok := n.(*ast.CreateStmt); ok && v.Partspec != nil && v.Relation != nil {
		if at, ok := e.resolveCreation(v.Relation); ok {
			e.partitioned[at] = true
		}
	}
}

// resolveCreation returns the name a relation CREATE makes takes, false
// when that is not known.
func (e *enginePath) resolveCreation(rv *ast.RangeVar) ([2]string, bool) {
	if rv.Relpersistence == ast.RELPERSISTENCE_TEMP || rv.Schemaname == "pg_temp" {
		return [2]string{"pg_temp", rv.Relname}, true
	}
	schema := rv.Schemaname
	if schema == "" {
		schema = e.creationSchema()
	}
	return [2]string{schema, rv.Relname}, schema != ""
}

// give notes that a name holds a relation from this statement on.
func (e *enginePath) give(name [2]string, h holding) {
	e.holds[name] = held{h, e.step}
}

// move follows a rename, to name in the same schema when schema is "", or
// a move to schema: the new name holds what the old one did, which then
// holds none.
func (e *enginePath) move(rv *ast.RangeVar, schema, name string) {
	from, h, ok := e.resolve(rv)
	if !ok {
		e.lost[rv.Relname] = e.step
		e.lost[name] = e.step
		return
	}
	if h.none() {
		// The statement fails, or IF EXISTS does nothing.
		return
	}
	if from[0] == "pg_temp" {
		delete(e.temps, rv.Relname)
		if schema == "" {
			e.temps[name] = e.step
		}
		return
	}
	if schema == "" {
		schema = from[0]
	}
	e.give(from, holding{})
	e.give([2]string{schema, name}, h)
}

// drop follows a DROP: the relations it names, and those of a schema it
// names, hold none after it. One that may drop relations it does not name,
// by CASCADE, an extension's, or a partitioned table's partitions, leaves
// every relation before it unsure.
func (e *enginePath) drop(v *ast.DropStmt) {
	kind := ast.ObjectType(v.RemoveType)
	if v.Objects == nil {
		return
	}
	if ast.DropBehavior(v.Behavior) == ast.DROP_CASCADE || kind == ast.OBJECT_EXTENSION {
		e.cascadeAt = e.step
	}
	if kind == ast.OBJECT_SCHEMA {
		for _, obj := range v.Objects.Items {
			for _, name := range nameParts(listOf(obj)) {
				e.takeSchema(name)
				delete(e.schemas, name)
			}
		}
		return
	}
	if !relationObject(kind) {
		return
	}
	for _, obj := range v.Objects.Items {
		parts := nameParts(listOf(obj))
		if len(parts) == 0 || len(parts) > 3 {
			continue
		}
		rv := &ast.RangeVar{Relname: parts[len(parts)-1]}
		if len(parts) > 1 {
			rv.Schemaname = parts[len(parts)-2]
		}
		if len(parts) > 2 {
			rv.Catalogname = parts[0]
		}
		at, h, ok := e.resolve(rv)
		if !ok {
			e.lost[rv.Relname] = e.step
			continue
		}
		if h.none() {
			continue
		}
		if at[0] == "pg_temp" {
			delete(e.temps, rv.Relname)
			continue
		}
		if h.made && e.partitioned[at] || !h.made && partitionedTable(e.schema, h.synced) {
			e.cascadeAt = e.step
		}
		e.give(at, holding{})
	}
}

// takeSchema notes that a schema name holds no relations from this
// statement on, and returns the synced schema whose relations it held,
// "" for none, and the relations it held otherwise.
func (e *enginePath) takeSchema(name string) (string, map[string]held) {
	synced, ok := e.schemaFrom[name]
	if !ok && !e.goneSchemas[name] {
		synced = name
	}
	delete(e.schemaFrom, name)
	delete(e.lostSchemas, name)
	e.goneSchemas[name] = true
	held := make(map[string]held)
	for k, h := range e.holds {
		if k[0] == name {
			held[k[1]] = h
			delete(e.holds, k)
		}
	}
	return synced, held
}

// renameSchema follows ALTER SCHEMA ... RENAME TO: the new name holds the
// relations the old one did, which then holds none.
func (e *enginePath) renameSchema(from, to string) {
	lostAt, lost := e.lostSchemas[from]
	synced, held := e.takeSchema(from)
	if lost {
		defer func() { e.lostSchemas[to] = lostAt }()
	}
	e.takeSchema(to)
	e.schemaFrom[to] = synced
	for name, h := range held {
		e.holds[[2]string{to, name}] = h
	}
	for k := range e.partitioned {
		if k[0] == from {
			delete(e.partitioned, k)
			e.partitioned[[2]string{to, k[1]}] = true
		}
	}
	if e.schemas[from] {
		delete(e.schemas, from)
		e.schemas[to] = true
	}
}

// current returns what a name holds as the next statement runs, false
// when that is not known.
func (e *enginePath) current(name [2]string) (holding, bool) {
	since := -1
	if at, ok := e.lost[name[1]]; ok {
		since = at
	}
	if at, ok := e.lostSchemas[name[0]]; ok && at > since {
		since = at
	}
	if e.opaqueAt > since {
		since = e.opaqueAt
	}
	if h, ok := e.holds[name]; ok && h.at > since {
		if !h.none() && h.at < e.cascadeAt {
			return holding{}, false
		}
		return h.holding, true
	}
	if since >= 0 {
		return holding{}, false
	}
	synced, ok := e.schemaFrom[name[0]]
	if !ok && !e.goneSchemas[name[0]] {
		synced = name[0]
	}
	if synced == "" || !hasBackupRelation(findSchema(e.schema, synced), name[1]) {
		return holding{}, true
	}
	if e.cascadeAt >= 0 {
		return holding{}, false
	}
	return holding{synced: [2]string{synced, name[1]}}, true
}

// resolve returns the name a relation reference of the next statement
// finds and the relation it holds, the zero holding for none, false when
// that is not known: the name as written, a temporary relation the change
// creates, or the first on the path that holds a relation.
func (e *enginePath) resolve(rv *ast.RangeVar) ([2]string, holding, bool) {
	if rv.Catalogname != "" && rv.Catalogname != e.schema.GetName() {
		return [2]string{}, holding{}, false
	}
	temp := [2]string{"pg_temp", rv.Relname}
	if rv.Schemaname == "pg_temp" {
		if isTemp, known := e.temp(rv.Relname); known && isTemp {
			return temp, holding{made: true}, true
		}
		return [2]string{}, holding{}, false
	}
	if rv.Schemaname != "" {
		name := [2]string{rv.Schemaname, rv.Relname}
		h, ok := e.current(name)
		return name, h, ok
	}
	if !e.pathKnown || strings.HasPrefix(rv.Relname, "pg_") {
		return [2]string{}, holding{}, false
	}
	// The temporary schema comes first unless the path names it.
	if !slices.Contains(e.path, "pg_temp") {
		isTemp, known := e.temp(rv.Relname)
		if !known {
			return [2]string{}, holding{}, false
		}
		if isTemp {
			return temp, holding{made: true}, true
		}
	}
	for _, p := range e.path {
		if p == "$user" {
			if !e.userKnown || e.user == "" {
				return [2]string{}, holding{}, false
			}
			p = e.user
		}
		if p == "pg_temp" {
			isTemp, known := e.temp(rv.Relname)
			if !known {
				return [2]string{}, holding{}, false
			}
			if isTemp {
				return temp, holding{made: true}, true
			}
			continue
		}
		name := [2]string{p, rv.Relname}
		h, ok := e.current(name)
		if !ok {
			return [2]string{}, holding{}, false
		}
		if !h.none() {
			return name, h, true
		}
	}
	return [2]string{}, holding{}, true
}

// temp reports whether a temporary relation of the name exists, and
// whether that is known: not after code may have created or dropped one.
func (e *enginePath) temp(name string) (bool, bool) {
	if at, ok := e.temps[name]; ok {
		return true, at > e.opaqueAt
	}
	return false, e.opaqueAt < e.tempsClearedAt || e.opaqueAt < 0
}

// run follows a statement other than UPDATE or DELETE.
func (e *enginePath) run(n ast.Node) {
	switch v := n.(type) {
	case *ast.VariableSetStmt:
		e.set(v)
	case *ast.UpdateStmt, *ast.DeleteStmt:
	default:
		e.follow(n)
	}
}

// copySource returns a relation the next statement copies the columns of,
// as the backup's query finds them: the synced relation the name holds,
// or, for one the change makes, the name as written; nil when that is not
// known, or the name holds none.
func (e *enginePath) copySource(rv *ast.RangeVar) *ast.RangeVar {
	_, h, ok := e.resolve(rv)
	switch {
	case !ok || h.none():
		return nil
	case h.made:
		return rv
	}
	return &ast.RangeVar{Schemaname: h.synced[0], Relname: h.synced[1]}
}

// holdsRelation reports whether a relation reference of the next statement
// certainly finds a relation.
func (e *enginePath) holdsRelation(rv *ast.RangeVar) bool {
	_, h, ok := e.resolve(rv)
	return ok && !h.none()
}

// changedRelation describes the relation a statement changes, at the name
// it finds, which is not table, the one the backup copies.
func changedRelation(at [2]string, h holding, table backupTable) string {
	name := relation(&ast.RangeVar{Schemaname: at[0], Relname: at[1]})
	switch {
	case at[0] == "pg_temp":
		return "the temporary table " + ident(at[1]) + " the change creates"
	case h.made && at == [2]string{table.schema, table.name}:
		return "the relation the change creates as " + name
	case !h.made && h.synced != at:
		return relation(&ast.RangeVar{Schemaname: h.synced[0], Relname: h.synced[1]}) + ", which the change renames or moves to " + name
	}
	return name
}

// relationObject reports whether DDL of the object type names a relation
// a DML statement may find.
func relationObject(t ast.ObjectType) bool {
	switch t {
	case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE:
		return true
	}
	return false
}

// partitionedTable reports whether the synced schema has the table
// partitioned.
func partitionedTable(schema *metadata.DatabaseSchemaMetadata, name [2]string) bool {
	for _, t := range findSchema(schema, name[0]).GetTables() {
		if t.GetName() == name[1] && len(t.GetPartitions()) > 0 {
			return true
		}
	}
	return false
}

// tempRelations returns the temporary relations a statement creates, IF
// NOT EXISTS or not.
func tempRelations(n ast.Node) map[string]bool {
	var rv *ast.RangeVar
	switch v := n.(type) {
	case *ast.ExplainStmt:
		if explainAnalyzes(v) {
			return tempRelations(v.Query)
		}
	case *ast.CreateStmt:
		rv = v.Relation
	case *ast.CreateSeqStmt:
		rv = v.Sequence
	case *ast.CreateTableAsStmt:
		if v.Into != nil {
			rv = v.Into.Rel
		}
	case *ast.SelectStmt:
		if into := intoOf(v); into != nil {
			rv = into.Rel
		}
	case *ast.ViewStmt:
		rv = v.View
	}
	if rv == nil || rv.Relpersistence != ast.RELPERSISTENCE_TEMP && rv.Schemaname != "pg_temp" {
		return nil
	}
	return map[string]bool{rv.Relname: true}
}

// creationSchema returns the schema an unqualified CREATE puts a relation
// in: the first of the path that exists, "" when that is not known.
func (e *enginePath) creationSchema() string {
	if !e.pathKnown {
		return ""
	}
	for _, p := range e.path {
		if p == "$user" {
			if !e.userKnown || e.user == "" {
				return ""
			}
			p = e.user
		}
		if _, lost := e.lostSchemas[p]; strings.HasPrefix(p, "pg_") || lost {
			return ""
		}
		_, renamed := e.schemaFrom[p]
		if findSchema(e.schema, p) != nil && !e.goneSchemas[p] || e.schemas[p] || renamed {
			return p
		}
	}
	return ""
}

// set follows a SET or RESET. RESET ALL leaves the role and the session
// user as they are.
//
// pg: src/backend/utils/misc/guc_tables.c — role, session_authorization
// (GUC_NO_RESET_ALL)
func (e *enginePath) set(v *ast.VariableSetStmt) {
	if v.Kind == ast.VAR_RESET_ALL {
		e.path, e.pathKnown = e.initial, true
		return
	}
	isPath, isRole := strings.EqualFold(v.Name, "search_path"), strings.EqualFold(v.Name, "role")
	isAuth := strings.EqualFold(v.Name, "session_authorization")
	if !isPath && !isRole && !isAuth {
		return
	}
	// unknown notes that the setting is not known.
	unknown := func() {
		switch {
		case isPath:
			e.pathKnown = false
		case isAuth:
			e.sessionKnown, e.userKnown = false, false
		default:
			e.userKnown = false
		}
	}
	if v.IsLocal {
		unknown()
		return
	}
	switch v.Kind {
	case ast.VAR_SET_DEFAULT, ast.VAR_RESET:
		switch {
		case isPath:
			e.path, e.pathKnown = e.initial, true
		case isAuth:
			e.session, e.user = e.login, e.login
			e.sessionKnown, e.userKnown = true, true
		default:
			e.user, e.userKnown = e.session, e.sessionKnown
		}
		return
	case ast.VAR_SET_CURRENT:
		return
	case ast.VAR_SET_VALUE:
	default:
		unknown()
		return
	}
	var values []string
	if v.Args != nil {
		for _, item := range v.Args.Items {
			s, ok := constString(item)
			if !ok {
				unknown()
				return
			}
			values = append(values, s)
		}
	}
	switch {
	case isPath && len(values) > 0:
		e.path, e.pathKnown = values, true
	case isRole && len(values) == 1 && strings.EqualFold(values[0], "none"):
		e.user, e.userKnown = e.session, e.sessionKnown
	case isRole && len(values) == 1:
		e.user, e.userKnown = values[0], true
	case isAuth && len(values) == 1:
		// The session user changes, and the current user with it.
		e.session, e.user = values[0], values[0]
		e.sessionKnown, e.userKnown = true, true
	default:
		unknown()
	}
}

// holders returns the schemas a relation reference of the backup's query
// may mean a synced relation of: the one written, or the first on the
// search path holding the name, or, with the path unknown, any holding
// it. A name only the system catalogs may hold has none.
func (q backupQuery) holders(rv *ast.RangeVar) []string {
	if rv.Schemaname != "" {
		if hasBackupRelation(findSchema(q.schema, rv.Schemaname), rv.Relname) {
			return []string{rv.Schemaname}
		}
		return nil
	}
	if q.known {
		for _, name := range q.path {
			if hasBackupRelation(findSchema(q.schema, name), rv.Relname) {
				return []string{name}
			}
		}
		return nil
	}
	var out []string
	for _, s := range q.schema.GetSchemas() {
		if hasBackupRelation(s, rv.Relname) {
			out = append(out, s.GetName())
		}
	}
	return out
}

// withOf returns the WITH clause of an UPDATE or DELETE.
func withOf(n ast.Node) *ast.WithClause {
	switch v := n.(type) {
	case *ast.UpdateStmt:
		return v.WithClause
	case *ast.DeleteStmt:
		return v.WithClause
	}
	return nil
}

// dmlParts returns what the backup copies of an UPDATE or DELETE: its
// table, and its FROM or USING list and WHERE clause.
func dmlParts(n ast.Node) (*ast.RangeVar, *ast.List, ast.Node) {
	switch v := n.(type) {
	case *ast.UpdateStmt:
		return v.Relation, v.FromClause, v.WhereClause
	case *ast.DeleteStmt:
		return v.Relation, v.UsingClause, v.WhereClause
	}
	return nil, nil, nil
}

// exposesFirst reports whether a later statement's table answers to the
// name the backup selects every arm through: the first statement's alias,
// which the later one gives too or, unaliased, is named; or, with none,
// the table's own name, which an alias would hide.
func exposesFirst(first, later *ast.RangeVar) bool {
	if a := aliasOf(first); a != "" {
		return aliasOf(later) == a || aliasOf(later) == "" && later.Relname == a
	}
	return aliasOf(later) == ""
}

// aliasOf returns the alias a statement gives its table, "" for none.
func aliasOf(rv *ast.RangeVar) string {
	if rv != nil && rv.Alias != nil {
		return rv.Alias.Aliasname
	}
	return ""
}

// selectionRefused reports why the query the backup runs for an UPDATE or
// DELETE fails, or does what the change does not ask, "" when the rule
// sees no reason. The backup copies, as written, the WITH clause of the
// table's first statement, the table and the FROM or USING list and WHERE
// clause of each statement, into one query per table, which selects
// through the first statement's alias, and runs it before the change:
//
//	CREATE TABLE bbdataarchive."…" AS [WITH …] SELECT a.* FROM t AS a[, …] [WHERE …] UNION …
//
// So the query fails on a later statement that does not use the first
// one's alias;
// a WITH query named like the table stands for it; a later statement's
// own WITH queries are left out; and it fails on a name the change makes,
// which does not exist yet: a relation, a column of one of the
// statement's relations, or a routine named with its schema. A name only
// the system catalogs may hold, an unqualified routine or a type, is not
// judged.
//
// bytebase: backend/plugin/parser/pg/backup.go — generateSQLForTable
func selectionRefused(n ast.Node, table backupTable, first ast.Node, q backupQuery, m made) string {
	target, from, where := dmlParts(n)
	if target == nil {
		return ""
	}
	isFirst := n == first
	with := withOf(first)
	firstTarget, _, _ := dmlParts(first)
	name := relation(&ast.RangeVar{Schemaname: table.schema, Relname: table.name})
	if !isFirst && !exposesFirst(firstTarget, target) {
		return "prior backup selects the rows of every statement on " + name + " through the first one's alias, which this statement does not use"
	}
	_, visible, q := withScope(with, nil, q, nil)
	if target.Schemaname == "" && visible[target.Relname] {
		return "prior backup copies the WITH clause of the first statement on " + name + ", whose query " + ident(target.Relname) + " then stands for the table"
	}
	if _, ok := where.(*ast.CurrentOfExpr); ok {
		return "prior backup copies WHERE CURRENT OF, which only UPDATE and DELETE take, into its query"
	}
	// The server finds the statement's table as written, along the path
	// the query runs under, which need not be the one Bytebase grouped it
	// by: the query then reads another table, or none.
	if target.Schemaname == "" && q.known {
		if h := q.holders(target); len(h) != 1 || h[0] != table.schema {
			found := "where it does not exist"
			if len(h) == 1 {
				found = "where it means " + relation(&ast.RangeVar{Schemaname: h[0], Relname: target.Relname})
			}
			return "prior backup reads " + ident(target.Relname) + " along the search path it runs under, " + found + ", not " + name
		}
	}
	// The parts the backup copies, beside the table.
	var parts []ast.Node
	if from != nil {
		parts = append(parts, from)
	}
	if where != nil {
		parts = append(parts, where)
	}
	// A later statement's own WITH queries are not in the backup's query.
	// A name of one that nothing else defines fails it. One the first
	// statement defines too, or a relation holds, reads other rows, which
	// may or may not be the statement's: that is not judged.
	if !isFirst {
		if own := withOf(n); own != nil && own.Ctes != nil {
			for _, item := range own.Ctes.Items {
				c, ok := item.(*ast.CommonTableExpr)
				if !ok || visible[c.Ctename] || q.holders(&ast.RangeVar{Relname: c.Ctename}) != nil {
					continue
				}
				named := false
				for _, part := range parts {
					relationsIn(part, visible, func(rv *ast.RangeVar) {
						named = named || rv.Schemaname == "" && rv.Relname == c.Ctename
					})
				}
				if named {
					return "prior backup copies only the WITH clause of the first statement on " + name + ", which does not define " + ident(c.Ctename)
				}
			}
		}
	}
	// A relation the change makes, which the query cannot find.
	var missing *ast.RangeVar
	check := func(rv *ast.RangeVar) {
		if missing != nil || rv == target || rv.Catalogname != "" && rv.Catalogname != q.schema.GetName() || systemSchema(rv.Schemaname) ||
			rv.Schemaname == "" && strings.HasPrefix(rv.Relname, "pg_") {
			return
		}
		if rv.Schemaname != "" {
			if m.relations[backupTable{rv.Schemaname, rv.Relname}] || m.relations[backupTable{name: rv.Relname}] {
				if q.holders(rv) == nil {
					missing = rv
				}
			}
			return
		}
		if m.relationNames[rv.Relname] && q.holders(rv) == nil {
			missing = rv
		}
	}
	if isFirst && with != nil {
		relationsIn(&ast.SelectStmt{WithClause: with}, nil, check)
	}
	for _, part := range parts {
		relationsIn(part, visible, check)
	}
	if missing != nil {
		return "prior backup runs before the change, when " + relation(missing) + " does not exist yet"
	}
	if msg := columnMissing(target, from, where, visible, q, m); msg != "" {
		return msg
	}
	if isFirst {
		if msg := withColumnMissing(with, q, m); msg != "" {
			return msg
		}
	}
	// A routine named with its schema, which the schema does not hold.
	var fn string
	inspect := func(node ast.Node) {
		ast.Inspect(node, func(x ast.Node) bool {
			call, ok := x.(*ast.FuncCall)
			if !ok || fn != "" {
				return fn == ""
			}
			if p := nameParts(call.Funcname); len(p) == 2 && !systemSchema(p[0]) && (m.functions[backupTable{p[0], p[1]}] || m.functions[backupTable{name: p[1]}]) && !hasFunction(findSchema(q.schema, p[0]), p[1]) {
				fn = ident(p[0]) + "." + ident(p[1])
			}
			return true
		})
	}
	if isFirst && with != nil {
		inspect(with)
	}
	for _, part := range parts {
		inspect(part)
	}
	if fn != "" {
		return "prior backup runs before the change, when function " + fn + " does not exist yet"
	}
	return ""
}

// dataModifyingCTEs returns the names of a WITH clause's queries that are
// not a SELECT: an INSERT, UPDATE, DELETE, or MERGE.
func dataModifyingCTEs(w *ast.WithClause) []string {
	if w == nil || w.Ctes == nil {
		return nil
	}
	var names []string
	for _, item := range w.Ctes.Items {
		if c, ok := item.(*ast.CommonTableExpr); ok {
			if _, isSelect := c.Ctequery.(*ast.SelectStmt); !isSelect {
				names = append(names, c.Ctename)
			}
		}
	}
	return names
}

// systemSchema reports whether a schema is one the snapshot leaves out but
// the server has.
func systemSchema(name string) bool {
	return strings.HasPrefix(name, "pg_") || name == "information_schema"
}

// hasFunction reports whether a synced schema holds a routine of that name.
func hasFunction(s *metadata.SchemaMetadata, name string) bool {
	if s == nil {
		return false
	}
	for _, f := range s.GetFunctions() {
		if f.GetName() == name {
			return true
		}
	}
	for _, p := range s.GetProcedures() {
		if p.GetName() == name {
			return true
		}
	}
	return false
}

// backupSource is a relation the backup's query reads columns from: the
// statement's table or one of its FROM or USING list, with the synced
// relations it may mean and their columns; columns is nil when they are
// not known.
type backupSource struct {
	// name is how the query refers to the source, its alias or its name,
	// and rv the relation, nil for a join.
	name    string
	rv      *ast.RangeVar
	schemas []string
	columns map[string]bool
	// owners are the relations, as the change writes them, whose columns
	// the source shows: the relation itself, the partitioned table of a
	// partition, or the relations of a join.
	owners []backupTable
	label  string
	// overflow marks an alias's column list longer than every relation
	// the name may mean has columns before the change.
	overflow bool
	// ordered are the columns in order, as a join's alias list renames
	// them, nil when not known.
	ordered []string
}

// columnMissing reports a column the backup's query names that the change
// adds to, or renames to, one of the statement's relations, and that none
// of them has before the change. The WHERE clause and the join conditions
// are read, outside subqueries, whose names may mean their own relations.
// A name that may mean something else is not judged: a column of a
// source whose columns are not known, a whole-row reference, or a system
// column.
func columnMissing(target *ast.RangeVar, from *ast.List, where ast.Node, visible map[string]bool, q backupQuery, m made) string {
	var exprs []ast.Node
	if where != nil {
		exprs = append(exprs, where)
	}
	return columnsMissing(target, from, exprs, nil, exprs, visible, q, m)
}

// withColumnMissing reports a column a query of a copied WITH clause names
// that the change makes on one of the query's relations, which none of
// them has before the change: read in its select list, WHERE, GROUP BY, or
// HAVING clause, among its own FROM list and the WITH queries before it,
// in each branch of a set operation.
func withColumnMissing(w *ast.WithClause, q backupQuery, m made) string {
	msg, _, _ := withScope(w, nil, q, func(sel *ast.SelectStmt, visible map[string]bool, q backupQuery) string {
		return selectColumnMissing(sel, visible, q, m)
	})
	return msg
}

// withScope returns the WITH queries visible after a WITH clause, by name,
// and q with the columns of those whose columns are known, after judging
// each query with check, when given, among those before it.
func withScope(w *ast.WithClause, around map[string]bool, q backupQuery, check func(*ast.SelectStmt, map[string]bool, backupQuery) string) (string, map[string]bool, backupQuery) {
	visible := maps.Clone(around)
	if visible == nil {
		visible = make(map[string]bool)
	}
	q.ctes = maps.Clone(q.ctes)
	if q.ctes == nil {
		q.ctes = make(map[string]backupSource)
	}
	if w == nil || w.Ctes == nil {
		return "", visible, q
	}
	if w.Recursive {
		for _, item := range w.Ctes.Items {
			if c, ok := item.(*ast.CommonTableExpr); ok {
				visible[c.Ctename] = true
				delete(q.ctes, c.Ctename)
			}
		}
	}
	for _, item := range w.Ctes.Items {
		c, ok := item.(*ast.CommonTableExpr)
		if !ok {
			continue
		}
		sel, isSelect := c.Ctequery.(*ast.SelectStmt)
		if isSelect && check != nil {
			if msg := check(sel, visible, q); msg != "" {
				return msg, visible, q
			}
		}
		var shape backupSource
		if isSelect && !w.Recursive {
			shape = queryShape(sel, c.Aliascolnames, visible, q)
		}
		visible[c.Ctename] = true
		delete(q.ctes, c.Ctename)
		if shape.columns != nil {
			q.ctes[c.Ctename] = shape
		}
	}
	return "", visible, q
}

// selectColumnMissing is withColumnMissing for one query of the clause:
// its own WITH clause, then each branch of a set operation among its own
// FROM list.
func selectColumnMissing(sel *ast.SelectStmt, visible map[string]bool, q backupQuery, m made) string {
	if sel == nil {
		return ""
	}
	if sel.WithClause != nil {
		msg, inner, scoped := withScope(sel.WithClause, visible, q, func(sel *ast.SelectStmt, visible map[string]bool, q backupQuery) string {
			return selectColumnMissing(sel, visible, q, m)
		})
		if msg != "" {
			return msg
		}
		visible, q = inner, scoped
	}
	if sel.Op != ast.SETOP_NONE {
		if msg := selectColumnMissing(sel.Larg, visible, q, m); msg != "" {
			return msg
		}
		return selectColumnMissing(sel.Rarg, visible, q, m)
	}
	var exprs, qualified []ast.Node
	// WINDOW reads the input columns only.
	//
	// pg: src/backend/parser/parse_clause.c — transformWindowDefinitions
	for _, n := range []ast.Node{sel.TargetList, sel.WhereClause, sel.HavingClause, sel.WindowClause, sel.LimitCount, sel.LimitOffset} {
		if n != nil && !isNilList(n) {
			exprs = append(exprs, n)
		}
	}
	// An unqualified name of GROUP BY, ORDER BY, or DISTINCT ON may mean
	// an output of the query.
	for _, n := range []ast.Node{sel.GroupClause, sel.SortClause, sel.DistinctClause} {
		if n != nil && !isNilList(n) {
			qualified = append(qualified, n)
		}
	}
	var conditions []ast.Node
	for _, n := range []ast.Node{sel.WhereClause, sel.HavingClause} {
		if n != nil {
			conditions = append(conditions, n)
		}
	}
	return columnsMissing(nil, sel.FromClause, exprs, qualified, conditions, maps.Clone(visible), q, m)
}

// levelNames returns the names a query's own FROM list shows.
func levelNames(sel *ast.SelectStmt) map[string]bool {
	names := make(map[string]bool)
	if sel.FromClause != nil {
		for _, item := range sel.FromClause.Items {
			for name := range exposedNames(item) {
				names[name] = true
			}
		}
	}
	return names
}

// withNames returns bound with the names of each set added.
func withNames(bound map[string]bool, sets ...map[string]bool) map[string]bool {
	out := maps.Clone(bound)
	if out == nil {
		out = make(map[string]bool)
	}
	for _, set := range sets {
		for name := range set {
			out[name] = true
		}
	}
	return out
}

// exposedNames returns the names a FROM item shows: its alias, the name
// of a relation or a function without one, and an unaliased join's sides'.
func exposedNames(n ast.Node) map[string]bool {
	names := make(map[string]bool)
	// exposed records the alias, or the names, a FROM item shows.
	exposed := func(a *ast.Alias, otherwise ...string) {
		if a != nil && a.Aliasname != "" {
			names[a.Aliasname] = true
			return
		}
		for _, name := range otherwise {
			names[name] = true
		}
	}
	var add func(n ast.Node)
	add = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.RangeVar:
			exposed(v.Alias, v.Relname)
		case *ast.RangeTableSample:
			add(v.Relation)
		case *ast.RangeSubselect:
			exposed(v.Alias)
		case *ast.RangeFunction:
			var funcs []string
			ast.Inspect(v.Functions, func(y ast.Node) bool {
				if f, ok := y.(*ast.FuncCall); ok {
					if parts := nameParts(f.Funcname); len(parts) > 0 {
						funcs = append(funcs, parts[len(parts)-1])
					}
				}
				return true
			})
			exposed(v.Alias, funcs...)
		case *ast.RangeTableFunc:
			exposed(v.Alias, "xmltable")
		case *ast.JsonTable:
			exposed(v.Alias, "json_table")
		case *ast.JoinExpr:
			if v.Alias != nil && v.Alias.Aliasname != "" {
				exposed(v.Alias)
				return
			}
			add(v.Larg)
			add(v.Rarg)
		}
	}
	add(n)
	return names
}

// isNilList reports whether a node is a nil *ast.List.
func isNilList(n ast.Node) bool {
	l, ok := n.(*ast.List)
	return ok && l == nil
}

// columnsMissing is columnMissing for a query of a table, nil for none,
// a FROM list, the expressions it reads, and the ones of which it judges
// only qualified names.
//
// conditions are the expressions of exprs a boolean is read from, where a
// bare name a source is named after cannot mean its whole row, which is
// no boolean.
func columnsMissing(target *ast.RangeVar, from *ast.List, exprs, qualifiedOnly, conditions []ast.Node, visible map[string]bool, q backupQuery, m made) string {
	var sources []backupSource
	known := true
	if target != nil {
		sources = []backupSource{sourceOf(target, nil, q)}
		known = sources[0].columns != nil
	}
	var quals []ast.Node
	has := func(src backupSource, column string) bool {
		return src.columns == nil || src.columns[column]
	}
	// makes reports whether the change makes the column on a relation the
	// source shows: adds it, renames to it, declares it, or gives the
	// relation a synced relation's columns, by LIKE, a rename, or a move.
	makes := func(src backupSource, column string) bool {
		for _, o := range src.owners {
			if m.makes(o, column, q, nil) {
				return true
			}
		}
		return false
	}
	var msg string
	judge := func(src backupSource, column string) {
		if !has(src, column) && makes(src, column) {
			msg = "prior backup runs before the change, when column " + ident(column) + " of " + src.label + " does not exist yet"
		}
	}
	// judgeSide judges a column one side of a join must have: none of its
	// sources has it before the change, which makes it on one of them.
	judgeSide := func(side []backupSource, column string) {
		for _, src := range side {
			if has(src, column) {
				return
			}
		}
		for _, src := range side {
			if makes(src, column) {
				judge(src, column)
				return
			}
		}
	}
	// every lists each source of the FROM list, a join's relations behind
	// its alias too, and scopes the sources a qualification inside an
	// aliased join reads, by its index in quals.
	every := slices.Clone(sources)
	scopes := make(map[int][]backupSource)
	// flatten adds the sources of a FROM item and returns them.
	var flatten func(n ast.Node) []backupSource
	flatten = func(n ast.Node) []backupSource {
		switch v := n.(type) {
		case *ast.RangeVar:
			src := sourceOf(v, visible, q)
			known = known && src.columns != nil
			sources = append(sources, src)
			every = append(every, src)
			return []backupSource{src}
		case *ast.RangeTableSample:
			return flatten(v.Relation)
		case *ast.RangeSubselect:
			// A LATERAL subquery may read the sources before it.
			if v.Lateral && v.Subquery != nil {
				quals = append(quals, &ast.SubLink{Subselect: v.Subquery})
			}
			src := subquerySource(v, visible, q)
			known = known && src.columns != nil
			sources = append(sources, src)
			every = append(every, src)
			return []backupSource{src}
		case *ast.RangeFunction:
			// The function's arguments may read the sources before it; what
			// it returns is known only from a column definition list.
			quals = append(quals, v.Functions)
			src := backupSource{}
			if v.Alias != nil {
				src.name, src.label = v.Alias.Aliasname, ident(v.Alias.Aliasname)
			}
			if !v.IsRowsfrom && !v.Ordinality && v.Coldeflist != nil && len(v.Coldeflist.Items) > 0 && v.Functions != nil && len(v.Functions.Items) == 1 {
				src.columns = make(map[string]bool)
				for _, item := range v.Coldeflist.Items {
					cd, ok := item.(*ast.ColumnDef)
					if !ok {
						src.columns = nil
						break
					}
					src.columns[cd.Colname] = true
				}
			}
			known = known && src.columns != nil
			sources = append(sources, src)
			every = append(every, src)
			return []backupSource{src}
		case *ast.RangeTableFunc:
			// XMLTABLE shows the columns it declares; its expressions may
			// read the sources before it.
			var names []string
			for _, n := range []ast.Node{v.Docexpr, v.Rowexpr} {
				if n != nil {
					quals = append(quals, n)
				}
			}
			if v.Columns != nil {
				for _, item := range v.Columns.Items {
					col, ok := item.(*ast.RangeTableFuncCol)
					if !ok {
						names = nil
						break
					}
					names = append(names, col.Colname)
					for _, n := range []ast.Node{col.Colexpr, col.Coldefexpr} {
						if n != nil {
							quals = append(quals, n)
						}
					}
				}
			}
			src := declaredSource(v.Alias, "xmltable", names)
			known = known && src.columns != nil
			sources = append(sources, src)
			every = append(every, src)
			return []backupSource{src}
		case *ast.JsonTable:
			// JSON_TABLE shows the columns it declares, nested ones too.
			if v.ContextItem != nil {
				quals = append(quals, v.ContextItem)
			}
			if v.Passing != nil {
				quals = append(quals, v.Passing)
			}
			if v.Columns != nil {
				quals = append(quals, v.Columns)
			}
			src := declaredSource(v.Alias, "json_table", jsonTableColumns(v.Columns))
			known = known && src.columns != nil
			sources = append(sources, src)
			every = append(every, src)
			return []backupSource{src}
		case *ast.JoinExpr:
			firstQual, firstSource := len(quals), len(sources)
			left, right := flatten(v.Larg), flatten(v.Rarg)
			children := append(slices.Clone(left), right...)
			if v.Quals != nil {
				quals = append(quals, v.Quals)
				conditions = append(conditions, v.Quals)
			}
			// USING names a column each side has.
			for _, column := range nameParts(v.UsingClause) {
				for _, side := range [][]backupSource{left, right} {
					if msg == "" {
						judgeSide(side, column)
					}
				}
			}
			// An aliased join shows its relations' columns under its
			// alias, as named unless the alias renames them.
			if a := v.Alias; a != nil && a.Aliasname != "" {
				join := backupSource{name: a.Aliasname, label: ident(a.Aliasname)}
				columns := make(map[string]bool)
				var ordered []string
				orderKnown := !v.IsNatural && v.UsingClause == nil
				for _, c := range children {
					join.owners = append(join.owners, c.owners...)
					if c.columns == nil {
						columns = nil
					}
					for name := range c.columns {
						if columns != nil {
							columns[name] = true
						}
					}
					if c.ordered == nil {
						orderKnown = false
					}
					ordered = append(ordered, c.ordered...)
				}
				// USING and NATURAL show the merged columns first, then each
				// side's others, known when each side is one source.
				//
				// pg: src/backend/parser/parse_clause.c — transformFromClauseItem
				if (v.IsNatural || v.UsingClause != nil) && len(left) == 1 && len(right) == 1 && left[0].ordered != nil && right[0].ordered != nil {
					merged := nameParts(v.UsingClause)
					if v.IsNatural {
						merged = nil
						for _, c := range left[0].ordered {
							if slices.Contains(right[0].ordered, c) {
								merged = append(merged, c)
							}
						}
					}
					ordered = slices.Clone(merged)
					for _, side := range [][]string{left[0].ordered, right[0].ordered} {
						for _, c := range side {
							if !slices.Contains(merged, c) {
								ordered = append(ordered, c)
							}
						}
					}
					orderKnown = true
				}
				switch {
				case a.Colnames == nil:
					join.columns = columns
					if orderKnown {
						join.ordered = ordered
					}
				case orderKnown && columns != nil:
					// The alias's list renames the join's columns in order.
					if len(a.Colnames.Items) > len(ordered) {
						join.overflow = true
					}
					names := slices.Clone(ordered)
					for i, item := range a.Colnames.Items {
						if alias, ok := item.(*ast.String); ok && i < len(names) {
							names[i] = alias.Str
						}
					}
					join.columns = make(map[string]bool)
					for _, n := range names {
						join.columns[n] = true
					}
					join.ordered = names
				}
				// The alias hides the join's relations: its conditions read
				// them, and the rest of the statement reads the join.
				inside := slices.Clone(sources[firstSource:])
				for i := firstQual; i < len(quals); i++ {
					if _, ok := scopes[i]; !ok {
						scopes[i] = inside
					}
				}
				sources = append(sources[:firstSource], join)
				every = append(every, join)
				return []backupSource{join}
			}
			return children
		default:
			known = false
			return nil
		}
	}
	if from != nil {
		for _, item := range from.Items {
			flatten(item)
		}
	}
	quals = append(quals, exprs...)
	// asBoolean are the bare names read as a boolean: a condition, or an
	// operand of AND, OR, or NOT in one.
	asBoolean := make(map[*ast.ColumnRef]bool)
	var mark func(n ast.Node)
	mark = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.ColumnRef:
			asBoolean[v] = true
		case *ast.BoolExpr:
			if v.Args != nil {
				for _, arg := range v.Args.Items {
					mark(arg)
				}
			}
		}
	}
	for _, c := range conditions {
		mark(c)
	}
	onlyQualified := make(map[int]bool)
	for _, n := range qualifiedOnly {
		onlyQualified[len(quals)] = true
		quals = append(quals, n)
	}
	// qualifiedFields returns the names of a column reference, with a
	// database.schema.relation.column one of the current database as
	// schema.relation.column; false for a system column or one of another
	// database.
	qualifiedFields := func(ref *ast.ColumnRef) ([]string, bool) {
		fields, ok := columnFields(ref)
		if ok && len(fields) == 4 && fields[0] == q.schema.GetName() {
			fields = fields[1:]
		}
		if !ok || len(fields) > 3 || systemColumn(fields[len(fields)-1]) {
			return nil, false
		}
		return fields, true
	}
	// scope is the sources the qualification being read sees.
	var scope []backupSource
	// named returns the source a qualified column reference names: by its
	// alias or name, or, schema.relation.column, a source without an alias
	// in that schema.
	named := func(fields []string) (backupSource, bool) {
		switch len(fields) {
		case 3:
			for _, src := range scope {
				if src.rv != nil && aliasOf(src.rv) == "" && src.rv.Relname == fields[1] && slices.Equal(src.schemas, []string{fields[0]}) {
					return src, true
				}
			}
		case 2:
			for _, src := range scope {
				if src.name == fields[0] {
					return src, true
				}
			}
		}
		return backupSource{}, false
	}
	// correlated judges the qualified names of a subquery that mean an
	// outer source: those no query level around them, its own included,
	// gives a relation of their own.
	var correlated func(n ast.Node, bound map[string]bool)
	correlated = func(n ast.Node, bound map[string]bool) {
		switch v := n.(type) {
		case *ast.SelectStmt:
			bound = withNames(bound, levelNames(v))
		case *ast.JoinExpr:
			// A join's own condition, and its relations, see the names
			// its sides show, which an alias hides from the rest.
			inner := withNames(bound, exposedNames(v.Larg), exposedNames(v.Rarg))
			for _, part := range []ast.Node{v.Larg, v.Rarg, v.Quals} {
				if part != nil {
					correlated(part, inner)
				}
			}
			return
		}
		ast.Inspect(n, func(y ast.Node) bool {
			if msg != "" {
				return false
			}
			if y == n {
				return true
			}
			switch w := y.(type) {
			case *ast.SelectStmt, *ast.JoinExpr:
				correlated(w, bound)
				return false
			case *ast.ColumnRef:
				fields, ok := qualifiedFields(w)
				if !ok || len(fields) < 2 || bound[fields[len(fields)-2]] {
					return false
				}
				if src, ok := named(fields); ok {
					judge(src, fields[len(fields)-1])
				}
				return false
			}
			return true
		})
	}
	// An alias's column list that names more columns than the relation has
	// before the change, which adds one.
	for _, src := range every {
		if msg == "" && src.overflow {
			for _, o := range src.owners {
				if len(m.columns[o]) > 0 || len(m.copies[o]) > 0 {
					what := src.label
					if src.rv == nil {
						what = "the join"
					}
					msg = "prior backup runs before the change, when " + what + " has fewer columns than alias " + ident(src.name) + " names"
					break
				}
			}
		}
	}
	for qi, qual := range quals {
		scope = sources
		if inside, ok := scopes[qi]; ok {
			scope = inside
		}
		ast.Inspect(qual, func(x ast.Node) bool {
			if msg != "" {
				return false
			}
			switch v := x.(type) {
			case *ast.SubLink:
				correlated(v.Subselect, nil)
				return false
			case *ast.A_Indirection:
				// (source).column selects a field of a source's row.
				ref, ok := v.Arg.(*ast.ColumnRef)
				if !ok || v.Indirection == nil || len(v.Indirection.Items) == 0 {
					return true
				}
				fields, ok := columnFields(ref)
				field, isField := v.Indirection.Items[0].(*ast.String)
				if !ok || len(fields) != 1 || !isField {
					return true
				}
				for _, src := range scope {
					if src.name == fields[0] {
						judge(src, field.Str)
						return false
					}
				}
				return true
			case *ast.ColumnRef:
				fields, ok := qualifiedFields(v)
				if !ok {
					return false
				}
				column := fields[len(fields)-1]
				switch len(fields) {
				case 3, 2:
					if src, ok := named(fields); ok {
						judge(src, column)
					}
				default:
					if !known || onlyQualified[qi] {
						return false
					}
					for _, src := range scope {
						if src.name == column && !asBoolean[v] || has(src, column) {
							return false
						}
					}
					for _, src := range scope {
						if judge(src, column); msg != "" {
							break
						}
					}
				}
				return false
			}
			return true
		})
	}
	return msg
}

// sourceOf returns a relation of the backup's query as a source of
// columns: the synced relations the name may mean, and the columns all of
// them have, or none known when the name may mean something whose
// columns the snapshot does not give: a WITH query, a materialized view,
// a sequence, or nothing the target holds.
func sourceOf(rv *ast.RangeVar, visible map[string]bool, q backupQuery) backupSource {
	src := backupSource{name: rv.Relname, rv: rv, label: relation(rv)}
	if a := aliasOf(rv); a != "" {
		src.name = a
	}
	if rv.Catalogname != "" && rv.Catalogname != q.schema.GetName() {
		return src
	}
	if rv.Schemaname == "" && visible[rv.Relname] {
		// A WITH query shows the columns of its query, when known.
		shape, ok := q.ctes[rv.Relname]
		if !ok || shape.columns == nil {
			return src
		}
		names := slices.Clone(shape.ordered)
		if rv.Alias != nil && rv.Alias.Colnames != nil {
			if len(rv.Alias.Colnames.Items) > len(names) {
				return src
			}
			for i, item := range rv.Alias.Colnames.Items {
				alias, ok := item.(*ast.String)
				if !ok {
					return src
				}
				names[i] = alias.Str
			}
		}
		src.columns = make(map[string]bool)
		for _, n := range names {
			src.columns[n] = true
		}
		src.ordered, src.owners = names, shape.owners
		return src
	}
	src.owners = []backupTable{{name: rv.Relname}}
	src.schemas = q.holders(rv)
	if len(src.schemas) == 0 {
		return src
	}
	columns := make(map[string]bool)
	known := true
	overflow, fits := false, false
	for _, s := range src.schemas {
		src.owners = append(src.owners, backupTable{s, rv.Relname})
		cols, parent, ok := relationColumns(findSchema(q.schema, s), rv.Relname)
		if parent != "" {
			src.owners = append(src.owners, backupTable{name: parent}, backupTable{s, parent})
		}
		if !ok {
			known = false
			continue
		}
		names := make([]string, len(cols))
		for i, c := range cols {
			names[i] = c.GetName()
		}
		// An alias's column list renames the columns in order.
		if rv.Alias != nil && rv.Alias.Colnames != nil {
			if len(rv.Alias.Colnames.Items) > len(names) {
				overflow = true
			} else {
				fits = true
			}
			for i, item := range rv.Alias.Colnames.Items {
				if alias, ok := item.(*ast.String); ok && i < len(names) {
					names[i] = alias.Str
				}
			}
		}
		for _, n := range names {
			columns[n] = true
		}
		if len(src.schemas) == 1 {
			src.ordered = names
		}
	}
	if known {
		src.columns = columns
		src.overflow = overflow && !fits
	} else {
		src.ordered = nil
	}
	return src
}

// declaredSource returns a table function of the FROM list as a source of
// the columns it declares, in order, which an alias's column list renames;
// named by its alias, or otherwise.
func declaredSource(alias *ast.Alias, otherwise string, names []string) backupSource {
	src := backupSource{name: otherwise, label: ident(otherwise)}
	if alias != nil && alias.Aliasname != "" {
		src.name, src.label = alias.Aliasname, ident(alias.Aliasname)
	}
	if len(names) == 0 {
		return src
	}
	names = slices.Clone(names)
	if alias != nil && alias.Colnames != nil {
		if len(alias.Colnames.Items) > len(names) {
			return src
		}
		for i, item := range alias.Colnames.Items {
			a, ok := item.(*ast.String)
			if !ok {
				return src
			}
			names[i] = a.Str
		}
	}
	src.columns = make(map[string]bool)
	for _, n := range names {
		src.columns[n] = true
	}
	src.ordered = names
	return src
}

// jsonTableColumns returns the names of the columns JSON_TABLE declares,
// of NESTED PATH ones too; nil when one is not a column.
func jsonTableColumns(columns *ast.List) []string {
	if columns == nil {
		return nil
	}
	var names []string
	for _, item := range columns.Items {
		col, ok := item.(*ast.JsonTableColumn)
		if !ok {
			return nil
		}
		if col.Coltype == ast.JTC_NESTED {
			nested := jsonTableColumns(col.Columns)
			if nested == nil {
				return nil
			}
			names = append(names, nested...)
			continue
		}
		if col.Name == "" {
			return nil
		}
		names = append(names, col.Name)
	}
	return names
}

// subquerySource returns a subquery or VALUES list of the FROM list as a
// source of the columns it shows, by its alias.
func subquerySource(v *ast.RangeSubselect, visible map[string]bool, q backupQuery) backupSource {
	if v.Alias == nil || v.Alias.Aliasname == "" {
		return backupSource{}
	}
	query, _ := v.Subquery.(*ast.SelectStmt)
	src := queryShape(query, v.Alias.Colnames, visible, q)
	src.name, src.label = v.Alias.Aliasname, ident(v.Alias.Aliasname)
	return src
}

// queryShape returns the columns a query shows, in order, which a column
// list renames: its outputs' aliases, the columns outputs read, the
// columns a star expands to, and columnN for VALUES, leftmost in a set
// operation. Its owners are the relations a star expands, whose columns
// the change makes the query show too. The columns are not known when an
// output has no name the statement gives, such as an expression, or a
// star expands a source whose columns are not known.
func queryShape(query *ast.SelectStmt, colnames *ast.List, visible map[string]bool, q backupQuery) backupSource {
	var src backupSource
	if query == nil {
		return src
	}
	// A WITH clause of the query's own stands for relations of its names.
	if query.WithClause != nil {
		_, visible, q = withScope(query.WithClause, visible, q, nil)
	}
	for query.Larg != nil {
		query = query.Larg
		if query.WithClause != nil {
			_, visible, q = withScope(query.WithClause, visible, q, nil)
		}
	}
	var names []string
	switch {
	case query.ValuesLists != nil && len(query.ValuesLists.Items) > 0:
		row, ok := query.ValuesLists.Items[0].(*ast.List)
		if !ok {
			return src
		}
		for i := range row.Items {
			names = append(names, fmt.Sprintf("column%d", i+1))
		}
	case query.TargetList != nil:
		for _, item := range query.TargetList.Items {
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				return src
			}
			if rt.Name != "" {
				names = append(names, rt.Name)
				continue
			}
			ref, ok := rt.Val.(*ast.ColumnRef)
			if !ok || ref.Fields == nil || len(ref.Fields.Items) == 0 {
				return backupSource{}
			}
			if _, star := ref.Fields.Items[len(ref.Fields.Items)-1].(*ast.A_Star); star {
				qualifier := nameParts(&ast.List{Items: ref.Fields.Items[:len(ref.Fields.Items)-1]})
				expanded, ok := expandStar(qualifier, len(ref.Fields.Items)-1, query.FromClause, visible, q)
				if !ok {
					return backupSource{}
				}
				names = append(names, expanded.ordered...)
				src.owners = append(src.owners, expanded.owners...)
				continue
			}
			fields, ok := columnFields(ref)
			if !ok {
				return backupSource{}
			}
			names = append(names, fields[len(fields)-1])
		}
	default:
		return src
	}
	if colnames != nil {
		if len(colnames.Items) > len(names) {
			return backupSource{}
		}
		for i, item := range colnames.Items {
			alias, ok := item.(*ast.String)
			if !ok {
				return backupSource{}
			}
			names[i] = alias.Str
		}
	}
	columns := make(map[string]bool)
	for _, n := range names {
		if n == "" {
			return backupSource{}
		}
		columns[n] = true
	}
	src.columns, src.ordered = columns, names
	return src
}

// expandStar returns the columns, in order, a star of a query's FROM list
// expands to, and the relations they come from: * every item's, which
// must each be a relation or a subquery, or relation.* the one item that
// shows the name, written as size parts. False when they are not known.
func expandStar(qualifier []string, size int, from *ast.List, visible map[string]bool, q backupQuery) (backupSource, bool) {
	var out backupSource
	if from == nil || len(qualifier) != size || size > 1 {
		return out, false
	}
	// item returns a FROM item as a source, false for one of another kind.
	item := func(n ast.Node) (backupSource, bool) {
		switch v := n.(type) {
		case *ast.RangeVar:
			return sourceOf(v, visible, q), true
		case *ast.RangeSubselect:
			return subquerySource(v, visible, q), true
		}
		return backupSource{}, false
	}
	add := func(src backupSource) bool {
		if src.ordered == nil {
			return false
		}
		out.ordered = append(out.ordered, src.ordered...)
		out.owners = append(out.owners, src.owners...)
		return true
	}
	if size == 0 {
		for _, n := range from.Items {
			src, ok := item(n)
			if !ok || !add(src) {
				return out, false
			}
		}
		return out, true
	}
	// find returns the item under the FROM list, or an unaliased join of
	// it, that shows the name.
	var found []backupSource
	var find func(n ast.Node) bool
	find = func(n ast.Node) bool {
		if j, ok := n.(*ast.JoinExpr); ok {
			if j.Alias != nil {
				return j.Alias.Aliasname != qualifier[0]
			}
			return find(j.Larg) && find(j.Rarg)
		}
		src, ok := item(n)
		if !ok {
			return false
		}
		if src.name == qualifier[0] {
			found = append(found, src)
		}
		return true
	}
	for _, n := range from.Items {
		if !find(n) {
			return out, false
		}
	}
	if len(found) != 1 || !add(found[0]) {
		return out, false
	}
	return out, true
}

// relationColumns returns the columns the snapshot gives a relation, in
// order: a table's, a partition's, which are its partitioned table's,
// named as parent, a view's, or a foreign table's; ok is false for
// another kind or a view the snapshot gives none.
func relationColumns(s *metadata.SchemaMetadata, name string) (columns []*metadata.ColumnMetadata, parent string, ok bool) {
	if s == nil {
		return nil, "", false
	}
	for _, t := range s.GetTables() {
		if t.GetName() == name {
			return t.GetColumns(), "", true
		}
	}
	for _, t := range s.GetTables() {
		if hasPartition(t.GetPartitions(), name) {
			return t.GetColumns(), t.GetName(), true
		}
	}
	for _, v := range s.GetViews() {
		if v.GetName() == name {
			return v.GetColumns(), "", len(v.GetColumns()) > 0
		}
	}
	for _, t := range s.GetExternalTables() {
		if t.GetName() == name {
			return t.GetColumns(), "", len(t.GetColumns()) > 0
		}
	}
	return nil, "", false
}

// columnFields returns the names of a column reference, false when it
// ends in *.
func columnFields(c *ast.ColumnRef) ([]string, bool) {
	if c.Fields == nil || len(c.Fields.Items) == 0 {
		return nil, false
	}
	var out []string
	for _, item := range c.Fields.Items {
		s, ok := item.(*ast.String)
		if !ok {
			return nil, false
		}
		out = append(out, s.Str)
	}
	return out, true
}

// systemColumn reports whether a name is one of the system columns every
// table has.
func systemColumn(name string) bool {
	switch name {
	case "ctid", "xmin", "xmax", "cmin", "cmax", "tableoid":
		return true
	}
	return false
}

func statementLoc(n ast.Node) ast.Loc {
	switch v := n.(type) {
	case *ast.UpdateStmt:
		return v.Loc
	case *ast.DeleteStmt:
		return v.Loc
	}
	return ast.NoLoc()
}

// createdRelations returns the relations a statement creates, renames, or
// moves to a new name, as written, the tables and views of CREATE SCHEMA
// in the new schema.
func createdRelations(n ast.Node) []backupTable {
	if v, ok := n.(*ast.CreateSchemaStmt); ok {
		name := v.Schemaname
		if name == "" && v.Authrole != nil {
			// CREATE SCHEMA AUTHORIZATION role names the schema after the role.
			name = v.Authrole.Rolename
		}
		if name == "" || v.SchemaElts == nil {
			return nil
		}
		var out []backupTable
		for _, elt := range v.SchemaElts.Items {
			for _, t := range createdRelations(elt) {
				if t.schema == "" {
					t.schema = name
				}
				out = append(out, t)
			}
		}
		return out
	}
	// EXPLAIN ANALYZE runs the statement it explains.
	if v, ok := n.(*ast.ExplainStmt); ok {
		if explainAnalyzes(v) {
			return createdRelations(v.Query)
		}
		return nil
	}
	if t, ok := createdRelation(n); ok {
		return []backupTable{t}
	}
	return nil
}

// createdRelation returns the relation a statement creates, renames, or
// moves to a new name, as written. A statement that may create nothing,
// IF NOT EXISTS or OR REPLACE, is left out: what it does depends on what
// already exists.
func createdRelation(n ast.Node) (backupTable, bool) {
	var rv *ast.RangeVar
	switch v := n.(type) {
	case *ast.CreateStmt:
		if !v.IfNotExists {
			rv = v.Relation
		}
	case *ast.CreateForeignTableStmt:
		if !v.Base.IfNotExists {
			rv = v.Base.Relation
		}
	case *ast.CreateSeqStmt:
		if !v.IfNotExists {
			rv = v.Sequence
		}
	case *ast.CreateTableAsStmt:
		if !v.IfNotExists && v.Into != nil {
			rv = v.Into.Rel
		}
	case *ast.SelectStmt:
		// INTO may sit on a branch of a set operation.
		if into := intoOf(v); into != nil {
			rv = into.Rel
		}
	case *ast.ViewStmt:
		if !v.Replace {
			rv = v.View
		}
	case *ast.RenameStmt:
		switch v.RenameType {
		case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE:
			if v.Relation != nil {
				return backupTable{schema: v.Relation.Schemaname, name: v.Newname}, true
			}
		}
	case *ast.AlterObjectSchemaStmt:
		switch v.ObjectType {
		case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW, ast.OBJECT_SEQUENCE:
			if v.Relation != nil {
				return backupTable{schema: v.Newschema, name: v.Relation.Relname}, true
			}
		}
	}
	if rv == nil {
		return backupTable{}, false
	}
	return backupTable{schema: rv.Schemaname, name: rv.Relname}, true
}

// createdColumn is a column a table the change creates declares.
type createdColumn struct {
	table  backupTable
	column string
}

// createdColumns returns the columns CREATE TABLE declares for the table
// it creates, the tables of CREATE SCHEMA in the new schema.
func createdColumns(n ast.Node) []createdColumn {
	var out []createdColumn
	switch v := n.(type) {
	case *ast.ExplainStmt:
		// EXPLAIN ANALYZE runs the statement it explains.
		if explainAnalyzes(v) {
			return createdColumns(v.Query)
		}
	case *ast.CreateSchemaStmt:
		name := v.Schemaname
		if name == "" && v.Authrole != nil {
			name = v.Authrole.Rolename
		}
		if name == "" || v.SchemaElts == nil {
			return nil
		}
		for _, elt := range v.SchemaElts.Items {
			for _, c := range createdColumns(elt) {
				if c.table.schema == "" {
					c.table.schema = name
				}
				out = append(out, c)
			}
		}
	case *ast.CreateStmt:
		if v.IfNotExists || v.Relation == nil || v.TableElts == nil {
			return nil
		}
		for _, item := range v.TableElts.Items {
			if cd, ok := item.(*ast.ColumnDef); ok {
				out = append(out, createdColumn{backupTable{v.Relation.Schemaname, v.Relation.Relname}, cd.Colname})
			}
		}
	case *ast.CreateTableAsStmt:
		if v.IfNotExists || v.Into == nil || v.Into.Rel == nil {
			return nil
		}
		query, _ := v.Query.(*ast.SelectStmt)
		out = intoColumns(v.Into, query)
	case *ast.ViewStmt:
		if v.Replace || v.View == nil {
			return nil
		}
		query, _ := v.Query.(*ast.SelectStmt)
		out = intoColumns(&ast.IntoClause{Rel: v.View, ColNames: v.Aliases}, query)
	case *ast.SelectStmt:
		if into := intoOf(v); into != nil && into.Rel != nil {
			out = intoColumns(into, v)
		}
	}
	return out
}

// createdLike is a relation CREATE TABLE ... LIKE, INHERITS, or PARTITION
// OF copies the columns of.
type createdLike struct {
	table  backupTable
	source *ast.RangeVar
}

// createdLikes returns the relations the tables a statement creates copy
// the columns of, by LIKE, INHERITS, or PARTITION OF, the tables of CREATE
// SCHEMA in the new schema.
func createdLikes(n ast.Node) []createdLike {
	var out []createdLike
	switch v := n.(type) {
	case *ast.ExplainStmt:
		if explainAnalyzes(v) {
			return createdLikes(v.Query)
		}
	case *ast.CreateSchemaStmt:
		name := v.Schemaname
		if name == "" && v.Authrole != nil {
			name = v.Authrole.Rolename
		}
		if name == "" || v.SchemaElts == nil {
			return nil
		}
		for _, elt := range v.SchemaElts.Items {
			for _, l := range createdLikes(elt) {
				if l.table.schema == "" {
					l.table.schema = name
				}
				out = append(out, l)
			}
		}
	case *ast.CreateStmt:
		if v.IfNotExists || v.Relation == nil {
			return nil
		}
		t := backupTable{v.Relation.Schemaname, v.Relation.Relname}
		if v.TableElts != nil {
			for _, item := range v.TableElts.Items {
				if l, ok := item.(*ast.TableLikeClause); ok && l.Relation != nil {
					out = append(out, createdLike{t, l.Relation})
				}
			}
		}
		// PARTITION OF names its parent here as well.
		if v.InhRelations != nil {
			for _, item := range v.InhRelations.Items {
				if parent, ok := item.(*ast.RangeVar); ok {
					out = append(out, createdLike{t, parent})
				}
			}
		}
	}
	return out
}

// intoColumns returns the columns CREATE TABLE AS or SELECT INTO gives the
// table it creates, as far as the statement names them: the column list
// written, then the names of the query's outputs, an alias or the column
// an output reads, leftmost in a set operation.
func intoColumns(into *ast.IntoClause, query *ast.SelectStmt) []createdColumn {
	table := backupTable{into.Rel.Schemaname, into.Rel.Relname}
	var names []string
	if into.ColNames != nil {
		names = nameParts(into.ColNames)
	}
	for query != nil && query.Larg != nil {
		query = query.Larg
	}
	if query != nil && query.TargetList != nil {
		for i, item := range query.TargetList.Items {
			if i < len(names) {
				continue
			}
			rt, ok := item.(*ast.ResTarget)
			if !ok {
				continue
			}
			name := rt.Name
			if ref, ok := rt.Val.(*ast.ColumnRef); ok && name == "" {
				if fields, ok := columnFields(ref); ok {
					name = fields[len(fields)-1]
				}
			}
			if name != "" {
				names = append(names, name)
			}
		}
	}
	out := make([]createdColumn, len(names))
	for i, n := range names {
		out[i] = createdColumn{table, n}
	}
	return out
}

// intoOf returns the INTO clause of a SELECT or of a branch of its set
// operation.
func intoOf(s *ast.SelectStmt) *ast.IntoClause {
	if s == nil {
		return nil
	}
	if s.IntoClause != nil {
		return s.IntoClause
	}
	if into := intoOf(s.Larg); into != nil {
		return into
	}
	return intoOf(s.Rarg)
}

// backupSearchPath is the synced search path as the backup reads it: the
// listed schemas without "$user" and the system schemas. An empty path
// is searched as public.
//
// bytebase: backend/store/model/pg_search_path.go — ResolvePGSearchPath
func backupSearchPath(setting string) []string {
	var path []string
	for _, name := range parseSearchPath(setting) {
		if name == "" || name == "$user" || isSystemSchema(name) {
			continue
		}
		path = append(path, name)
	}
	return path
}

func isSystemSchema(name string) bool {
	for _, s := range []string{"pg_catalog", "information_schema", "pg_toast", "pg_temp_1", "pg_temp_2", "pg_global"} {
		if strings.EqualFold(name, s) {
			return true
		}
	}
	return false
}

// backupSetSearchPath reads SET search_path as the backup does: every
// string value, as written, replaces the path. A SET with no string value
// changes nothing.
//
// bytebase: backend/plugin/parser/pg/backup.go — prepareTransformation
func backupSetSearchPath(v *ast.VariableSetStmt) ([]string, bool) {
	if !strings.EqualFold(v.Name, "search_path") || v.Args == nil {
		return nil, false
	}
	var path []string
	for _, item := range v.Args.Items {
		if s, ok := constString(item); ok {
			path = append(path, s)
		}
	}
	return path, len(path) > 0
}

// searchBackupRelation returns the first schema on the path holding a
// relation of that name, or "".
//
// bytebase: backend/store/model/pg_searcher.go — SearchRelation
func searchBackupRelation(schema *metadata.DatabaseSchemaMetadata, path []string, name string) string {
	if len(path) == 0 {
		path = []string{"public"}
	}
	for _, p := range path {
		if hasBackupRelation(findSchema(schema, p), name) {
			return p
		}
	}
	return ""
}

func findSchema(schema *metadata.DatabaseSchemaMetadata, name string) *metadata.SchemaMetadata {
	for _, s := range schema.GetSchemas() {
		if s.GetName() == name {
			return s
		}
	}
	return nil
}

// hasBackupRelation reports whether the schema holds a relation of that
// name the backup's lookup finds: a table or any of its partitions, a
// view, a materialized view, an external table, or a sequence.
func hasBackupRelation(s *metadata.SchemaMetadata, name string) bool {
	if s == nil {
		return false
	}
	for _, t := range s.GetTables() {
		if t.GetName() == name || hasPartition(t.GetPartitions(), name) {
			return true
		}
	}
	for _, v := range s.GetViews() {
		if v.GetName() == name {
			return true
		}
	}
	for _, v := range s.GetMaterializedViews() {
		if v.GetName() == name {
			return true
		}
	}
	for _, t := range s.GetExternalTables() {
		if t.GetName() == name {
			return true
		}
	}
	for _, q := range s.GetSequences() {
		if q.GetName() == name {
			return true
		}
	}
	return false
}

func hasPartition(partitions []*metadata.TablePartitionMetadata, name string) bool {
	for _, p := range partitions {
		if p.GetName() == name || hasPartition(p.GetSubpartitions(), name) {
			return true
		}
	}
	return false
}
