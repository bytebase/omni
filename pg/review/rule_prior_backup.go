package review

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what makes Bytebase's prior backup of a PostgreSQL
// change fail, and only that: a failed backup fails the task before any
// of the change runs. The backup works like this (bytebase
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
// change creates.
//
// A name the synced schema cannot resolve, and the change does not
// create, is not reported. The schema the review sees may predate an
// earlier spec of the same plan, which may create the table; Bytebase
// syncs again before this spec's backup, which then finds it. A name no
// spec creates fails the UPDATE or DELETE itself, which is WalkThrough's
// to report. Neither is what the backup's query may mean otherwise: an
// unqualified function or a type, which the system catalogs may hold, a
// routine call an overload of another signature may take, a column a
// subquery names, a column whose type the change alters, or a WITH query
// of the first statement standing for a later one's own WITH query or
// relation, which reads rows that may or may not be the statement's.

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
// fails. Without a schema nothing is known about the target, and nothing
// is reported. sessionUser is the user the backup connects as.
func checkPriorBackupTarget(stmts []statement, schema *metadata.DatabaseSchemaMetadata, sessionUser string, r *reporter) {
	if schema == nil {
		return
	}
	path := backupSearchPath(schema.GetSearchPath())
	// kinds maps each table to the statement kind the backup saw first,
	// and mixed marks the tables already reported.
	kinds := make(map[backupTable]string)
	mixed := make(map[backupTable]bool)
	made := changeMade(stmts)
	q := backupQueryOf(stmts, schema, sessionUser)
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
			continue
		default:
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
		if rv.Schemaname != "" && (rv.Catalogname == "" || rv.Catalogname == schema.GetName()) && (made.relations[table] || made.relations[backupTable{name: rv.Relname}]) && !hasBackupRelation(findSchema(schema, rv.Schemaname), rv.Relname) {
			// An unqualified creation may be in this schema; either way
			// the table is not there when the backup runs.
			reportNotYetCreated(s, rv, r)
		} else {
			first, seenFirst := firsts[table]
			if !seenFirst {
				first = s.node
				firsts[table] = first
			}
			if msg := selectionRefused(s.node, table, first, q, made); msg != "" {
				r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), msg)
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
}

func changeMade(stmts []statement) made {
	m := made{
		relations:     make(map[backupTable]bool),
		relationNames: make(map[string]bool),
		columns:       make(map[backupTable]map[string]bool),
		functions:     make(map[backupTable]bool),
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
	for i := range stmts {
		for _, t := range createdRelations(stmts[i].node) {
			m.relations[t] = true
			m.relationNames[t.name] = true
		}
		// A table the change creates has the columns it declares, which a
		// table of that name the backup reads may lack.
		for _, c := range createdColumns(stmts[i].node) {
			column(c.table, c.column)
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
		if v.Kind != ast.VAR_SET_VALUE || v.IsLocal || v.Args == nil {
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
	if !isFirst && aliasOf(target) != aliasOf(firstTarget) {
		return "prior backup selects the rows of every statement on " + name + " through the first one's alias, which this statement does not use"
	}
	visible := make(map[string]bool)
	if with != nil && with.Ctes != nil {
		for _, item := range with.Ctes.Items {
			if c, ok := item.(*ast.CommonTableExpr); ok {
				visible[c.Ctename] = true
			}
		}
	}
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
}

// columnMissing reports a column the backup's query names that the change
// adds to, or renames to, one of the statement's relations, and that none
// of them has before the change. The WHERE clause and the join conditions
// are read, outside subqueries, whose names may mean their own relations.
// A name that may mean something else is not judged: a column of a
// source whose columns are not known, a whole-row reference, or a system
// column.
func columnMissing(target *ast.RangeVar, from *ast.List, where ast.Node, visible map[string]bool, q backupQuery, m made) string {
	sources := []backupSource{sourceOf(target, nil, q)}
	known := sources[0].columns != nil
	var quals []ast.Node
	// flatten adds the sources of a FROM item and returns them.
	var flatten func(n ast.Node) []backupSource
	flatten = func(n ast.Node) []backupSource {
		switch v := n.(type) {
		case *ast.RangeVar:
			src := sourceOf(v, visible, q)
			known = known && src.columns != nil
			sources = append(sources, src)
			return []backupSource{src}
		case *ast.RangeTableSample:
			return flatten(v.Relation)
		case *ast.JoinExpr:
			children := append(flatten(v.Larg), flatten(v.Rarg)...)
			if v.Quals != nil {
				quals = append(quals, v.Quals)
			}
			if v.IsNatural || v.UsingClause != nil {
				known = false
			}
			// An aliased join shows its relations' columns under its
			// alias, as named unless the alias renames them.
			if a := v.Alias; a != nil && a.Aliasname != "" {
				join := backupSource{name: a.Aliasname, label: ident(a.Aliasname)}
				columns := make(map[string]bool)
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
				}
				if a.Colnames == nil {
					join.columns = columns
				}
				sources = append(sources, join)
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
	if where != nil {
		quals = append(quals, where)
	}
	has := func(src backupSource, column string) bool {
		return src.columns == nil || src.columns[column]
	}
	var msg string
	judge := func(src backupSource, column string) {
		if has(src, column) {
			return
		}
		for _, o := range src.owners {
			if m.columns[o][column] {
				msg = "prior backup runs before the change, when column " + ident(column) + " of " + src.label + " does not exist yet"
				return
			}
		}
	}
	for _, qual := range quals {
		ast.Inspect(qual, func(x ast.Node) bool {
			if msg != "" {
				return false
			}
			switch v := x.(type) {
			case *ast.SubLink:
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
				for _, src := range sources {
					if src.name == fields[0] {
						judge(src, field.Str)
						return false
					}
				}
				return true
			case *ast.ColumnRef:
				fields, ok := columnFields(v)
				// database.schema.relation.column names the current database.
				if ok && len(fields) == 4 && fields[0] == q.schema.GetName() {
					fields = fields[1:]
				}
				if !ok || len(fields) > 3 || systemColumn(fields[len(fields)-1]) {
					return false
				}
				column := fields[len(fields)-1]
				switch len(fields) {
				case 3:
					// schema.relation.column names a source without an
					// alias in that schema.
					for _, src := range sources {
						if src.rv != nil && aliasOf(src.rv) == "" && src.rv.Relname == fields[1] && slices.Equal(src.schemas, []string{fields[0]}) {
							judge(src, column)
							break
						}
					}
				case 2:
					for _, src := range sources {
						if src.name == fields[0] {
							judge(src, column)
							break
						}
					}
				default:
					if !known {
						return false
					}
					for _, src := range sources {
						if src.name == column || has(src, column) {
							return false
						}
					}
					for _, src := range sources {
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
	if rv.Catalogname != "" && rv.Catalogname != q.schema.GetName() || rv.Schemaname == "" && visible[rv.Relname] {
		return src
	}
	src.owners = []backupTable{{name: rv.Relname}}
	src.schemas = q.holders(rv)
	if len(src.schemas) == 0 {
		return src
	}
	columns := make(map[string]bool)
	known := true
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
			for i, item := range rv.Alias.Colnames.Items {
				if alias, ok := item.(*ast.String); ok && i < len(names) {
					names[i] = alias.Str
				}
			}
		}
		for _, n := range names {
			columns[n] = true
		}
	}
	if known {
		src.columns = columns
	}
	return src
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
	case *ast.SelectStmt:
		if into := intoOf(v); into != nil && into.Rel != nil {
			out = intoColumns(into, v)
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
