package review

import (
	"fmt"
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
//   - It refuses a table touched by both UPDATE and DELETE.
//   - It reads the rows with one query per table, which copies, as
//     written, the WITH clause of the table's first statement and each
//     statement's FROM or USING list and WHERE clause, and runs it with
//     the change's first SET role or search_path before it.
//
// DDL in the change is ignored by the backup, so it is not a finding by
// itself. It matters when the backup's query names what only the change
// makes, which does not exist yet when the query runs: the table, a
// relation of the FROM or USING list, a subquery, or the WITH clause, a
// column the change adds or renames to, or a function named with its
// schema. The query also runs a data-modifying WITH query of a table's
// first statement, before the change and without backing up its rows,
// and leaves out a later statement's own WITH queries.
//
// A name the synced schema cannot resolve, and the change does not
// create, is not reported. The schema the review sees may predate an
// earlier spec of the same plan, which may create the table; Bytebase
// syncs again before this spec's backup, which then finds it. A name no
// spec creates fails the UPDATE or DELETE itself, which is WalkThrough's
// to report. Neither is what the backup's query may mean otherwise: an
// unqualified function or a type, which the system catalogs may hold, a
// column a subquery names, or a column whose type the change alters.

// backupSchema is where Bytebase writes a PostgreSQL prior backup.
const backupSchema = "bbdataarchive"

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
// is reported.
func checkPriorBackupTarget(stmts []statement, schema *metadata.DatabaseSchemaMetadata, r *reporter) {
	if schema == nil {
		return
	}
	path := backupSearchPath(schema.GetSearchPath())
	// kinds maps each table to the statement kind the backup saw first,
	// and mixed marks the tables already reported.
	kinds := make(map[backupTable]string)
	mixed := make(map[backupTable]bool)
	made := changeMade(stmts)
	// firstWith is the WITH clause of the first statement on each table,
	// the one the backup's query of the table carries.
	firstWith := make(map[backupTable]*ast.WithClause)
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
			with, seenWith := firstWith[table]
			if !seenWith {
				with = withOf(s.node)
				firstWith[table] = with
			}
			if msg := selectionRefused(s.node, table, with, !seenWith, schema, made); msg != "" {
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
}

func reportNotYetCreated(s *statement, rv *ast.RangeVar, r *reporter) {
	r.report(review.PriorBackup, s.index, rangeOf(statementLoc(s.node)), "prior backup runs before the change, when "+relation(rv)+" does not exist yet")
}

// made is what the change makes, which the backup, run before all of
// it, cannot see: the relations it creates or gives a name, as written,
// an unqualified one under an empty schema; the columns it adds or gives
// a name, by the name of the table as written; and the functions it
// creates, by name. A statement that may make nothing, IF NOT EXISTS or
// OR REPLACE, is left out: what it does depends on what already exists.
type made struct {
	relations     map[backupTable]bool
	relationNames map[string]bool
	columns       map[string]map[string]bool
	functions     map[string]bool
}

func changeMade(stmts []statement) made {
	m := made{
		relations:     make(map[backupTable]bool),
		relationNames: make(map[string]bool),
		columns:       make(map[string]map[string]bool),
		functions:     make(map[string]bool),
	}
	column := func(table, name string) {
		if m.columns[table] == nil {
			m.columns[table] = make(map[string]bool)
		}
		m.columns[table][name] = true
	}
	for i := range stmts {
		for _, t := range createdRelations(stmts[i].node) {
			m.relations[t] = true
			m.relationNames[t.name] = true
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
					column(v.Relation.Relname, cd.Colname)
				}
			}
		case *ast.RenameStmt:
			if v.RenameType == ast.OBJECT_COLUMN && v.Relation != nil {
				column(v.Relation.Relname, v.Newname)
			}
		case *ast.CreateFunctionStmt:
			if parts := nameParts(v.Funcname); len(parts) > 0 && !v.IsOrReplace {
				m.functions[parts[len(parts)-1]] = true
			}
		}
	}
	return m
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

// selectionRefused reports why the query the backup runs for an UPDATE or
// DELETE fails, or does what the change does not ask, "" when the scan
// sees no reason. The backup copies the statement's FROM or USING list
// and WHERE clause as written into one query per table, after the WITH
// clause of the table's first statement, and runs it before the change:
//
//	CREATE TABLE bbdataarchive."…" AS [WITH …] SELECT t.* FROM t[, …] [WHERE …]
//
// So it runs a data-modifying WITH query of that first statement, and it
// fails on a name the change makes, which does not exist yet: a relation,
// a column of one of the statement's tables, or a function named with
// its schema, or on a query of a later statement naming a WITH query of
// its own. A name only the system catalogs may hold, an unqualified
// function or a type, is not judged.
//
// bytebase: backend/plugin/parser/pg/backup.go — generateSQLForTable
func selectionRefused(n ast.Node, table backupTable, with *ast.WithClause, first bool, schema *metadata.DatabaseSchemaMetadata, m made) string {
	var target *ast.RangeVar
	var from *ast.List
	var where ast.Node
	switch v := n.(type) {
	case *ast.UpdateStmt:
		target, from, where = v.Relation, v.FromClause, v.WhereClause
	case *ast.DeleteStmt:
		target, from, where = v.Relation, v.UsingClause, v.WhereClause
	default:
		return ""
	}
	if first && with != nil && modifiesData(with) {
		return "prior backup runs the data-modifying WITH query of this statement before the change, and backs up none of the rows it changes"
	}
	// The parts the backup copies, and the WITH queries they can name.
	var parts []ast.Node
	if from != nil {
		parts = append(parts, from)
	}
	if where != nil {
		parts = append(parts, where)
	}
	visible := make(map[string]bool)
	if with != nil && with.Ctes != nil {
		for _, item := range with.Ctes.Items {
			if c, ok := item.(*ast.CommonTableExpr); ok {
				visible[c.Ctename] = true
			}
		}
	}
	anywhere := func(name string) bool {
		for _, sc := range schema.GetSchemas() {
			if hasBackupRelation(sc, name) {
				return true
			}
		}
		return false
	}
	// A later statement's own WITH queries are not in the backup's query.
	if !first {
		if own := withOf(n); own != nil && own.Ctes != nil {
			for _, item := range own.Ctes.Items {
				c, ok := item.(*ast.CommonTableExpr)
				if !ok || visible[c.Ctename] || anywhere(c.Ctename) {
					continue
				}
				named := false
				for _, part := range parts {
					relationsIn(part, visible, func(rv *ast.RangeVar) {
						named = named || rv.Schemaname == "" && rv.Relname == c.Ctename
					})
				}
				if named {
					return "prior backup copies only the WITH clause of the first statement on " + relation(&ast.RangeVar{Schemaname: table.schema, Relname: table.name}) + ", which does not define " + ident(c.Ctename)
				}
			}
		}
	}
	// A relation the change makes, which no schema of the target holds.
	var missing *ast.RangeVar
	check := func(rv *ast.RangeVar) {
		if missing != nil || rv == target || rv.Catalogname != "" && rv.Catalogname != schema.GetName() || systemSchema(rv.Schemaname) {
			return
		}
		if rv.Schemaname != "" {
			if (m.relations[backupTable{rv.Schemaname, rv.Relname}] || m.relations[backupTable{name: rv.Relname}]) && !hasBackupRelation(findSchema(schema, rv.Schemaname), rv.Relname) {
				missing = rv
			}
			return
		}
		if m.relationNames[rv.Relname] && !strings.HasPrefix(rv.Relname, "pg_") && !anywhere(rv.Relname) {
			missing = rv
		}
	}
	if first && with != nil {
		relationsIn(&ast.SelectStmt{WithClause: with}, nil, check)
	}
	for _, part := range parts {
		relationsIn(part, visible, check)
	}
	if missing != nil {
		return "prior backup runs before the change, when " + relation(missing) + " does not exist yet"
	}
	if msg := columnMissing(target, from, where, visible, schema, m); msg != "" {
		return msg
	}
	// A function named with its schema, which the schema does not hold.
	var fn string
	inspect := func(node ast.Node) {
		ast.Inspect(node, func(x ast.Node) bool {
			call, ok := x.(*ast.FuncCall)
			if !ok || fn != "" {
				return fn == ""
			}
			if p := nameParts(call.Funcname); len(p) == 2 && !systemSchema(p[0]) && m.functions[p[1]] && !hasFunction(findSchema(schema, p[0]), p[1]) {
				fn = ident(p[0]) + "." + ident(p[1])
			}
			return true
		})
	}
	if first && with != nil {
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

// modifiesData reports whether a WITH clause holds an INSERT, UPDATE,
// DELETE, or MERGE, which runs wherever the clause is copied.
func modifiesData(w *ast.WithClause) bool {
	if w.Ctes == nil {
		return false
	}
	for _, item := range w.Ctes.Items {
		if c, ok := item.(*ast.CommonTableExpr); ok {
			switch c.Ctequery.(type) {
			case *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.MergeStmt:
				return true
			}
		}
	}
	return false
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
// tables its name may mean.
type backupSource struct {
	rv     *ast.RangeVar
	tables []*metadata.TableMetadata
}

// columnMissing reports a column the backup's query names that the change
// adds to, or renames to, one of the statement's tables, and that none of
// them has before the change. The WHERE clause and the join conditions
// are read, outside subqueries, whose names may mean their own tables. A
// name that may mean something else is not judged: a column of a source
// that is not a synced table, a whole-row reference, or a system column.
func columnMissing(target *ast.RangeVar, from *ast.List, where ast.Node, visible map[string]bool, schema *metadata.DatabaseSchemaMetadata, m made) string {
	// The server resolves the statement's table as written when the backup
	// runs, which may be along another path than the one Bytebase looked
	// it up with: an unqualified name may mean a table of it in any schema.
	targetTables, known := sourceTables(target, nil, schema)
	sources := []backupSource{{rv: target, tables: targetTables}}
	var quals []ast.Node
	var flatten func(n ast.Node)
	flatten = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.RangeVar:
			tables, ok := sourceTables(v, visible, schema)
			if !ok {
				known = false
			}
			sources = append(sources, backupSource{rv: v, tables: tables})
		case *ast.JoinExpr:
			flatten(v.Larg)
			flatten(v.Rarg)
			if v.Quals != nil {
				quals = append(quals, v.Quals)
			}
			if v.IsNatural || v.UsingClause != nil {
				known = false
			}
		default:
			known = false
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
	// name is how the query refers to a source: its alias, or its name.
	name := func(src backupSource) string {
		if src.rv.Alias != nil && src.rv.Alias.Aliasname != "" {
			return src.rv.Alias.Aliasname
		}
		return src.rv.Relname
	}
	has := func(src backupSource, column string) bool {
		if src.tables == nil {
			return true
		}
		for _, t := range src.tables {
			for _, c := range t.GetColumns() {
				if c.GetName() == column {
					return true
				}
			}
		}
		return false
	}
	var msg string
	judge := func(src backupSource, column string) {
		if !has(src, column) && m.columns[src.rv.Relname][column] {
			msg = "prior backup runs before the change, when column " + ident(column) + " of " + relation(src.rv) + " does not exist yet"
		}
	}
	for _, q := range quals {
		ast.Inspect(q, func(x ast.Node) bool {
			if msg != "" {
				return false
			}
			switch v := x.(type) {
			case *ast.SubLink:
				return false
			case *ast.ColumnRef:
				fields, ok := columnFields(v)
				if !ok || len(fields) > 2 || systemColumn(fields[len(fields)-1]) {
					return false
				}
				if len(fields) == 2 {
					for _, src := range sources {
						if name(src) == fields[0] {
							judge(src, fields[1])
							break
						}
					}
					return false
				}
				column := fields[0]
				if !known {
					return false
				}
				for _, src := range sources {
					if name(src) == column || has(src, column) {
						return false
					}
				}
				for _, src := range sources {
					judge(src, column)
					if msg != "" {
						break
					}
				}
				return false
			}
			return true
		})
	}
	return msg
}

// sourceTables returns the synced tables a FROM or USING relation may
// mean: the one in the schema written, or one of that name in any schema.
// ok is false when the name may mean something else: a WITH query, a
// relation of another kind, or none the target holds.
func sourceTables(rv *ast.RangeVar, visible map[string]bool, schema *metadata.DatabaseSchemaMetadata) ([]*metadata.TableMetadata, bool) {
	if rv.Catalogname != "" || rv.Schemaname == "" && visible[rv.Relname] {
		return nil, false
	}
	var schemas []*metadata.SchemaMetadata
	if rv.Schemaname != "" {
		schemas = []*metadata.SchemaMetadata{findSchema(schema, rv.Schemaname)}
	} else {
		schemas = schema.GetSchemas()
	}
	var tables []*metadata.TableMetadata
	for _, sc := range schemas {
		if sc == nil || !hasBackupRelation(sc, rv.Relname) {
			continue
		}
		t := findTable(sc, rv.Relname)
		if t == nil {
			return nil, false
		}
		tables = append(tables, t)
	}
	return tables, len(tables) > 0
}

func findTable(s *metadata.SchemaMetadata, name string) *metadata.TableMetadata {
	if s == nil {
		return nil
	}
	for _, t := range s.GetTables() {
		if t.GetName() == name {
			return t
		}
	}
	return nil
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
		case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
			if v.Relation != nil {
				return backupTable{schema: v.Relation.Schemaname, name: v.Newname}, true
			}
		}
	case *ast.AlterObjectSchemaStmt:
		switch v.ObjectType {
		case ast.OBJECT_TABLE, ast.OBJECT_FOREIGN_TABLE, ast.OBJECT_VIEW, ast.OBJECT_MATVIEW:
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
