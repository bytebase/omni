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
//
// DDL in the change is ignored by the backup, so it is not a finding by
// itself. It matters when an UPDATE or DELETE touches a table that only
// the change creates: the backup reads that table before the change has
// created it.
//
// A name the synced schema cannot resolve, and the change does not
// create, is not reported. The schema the review sees may predate an
// earlier spec of the same plan, which may create the table; Bytebase
// syncs again before this spec's backup, which then finds it. A name no
// spec creates fails the UPDATE or DELETE itself, which is WalkThrough's
// to report.

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
	// created lists the relations the change creates, by name as written:
	// an unqualified one under an empty schema.
	// The backup runs before the whole change, so a relation the change
	// creates anywhere in it is not there yet.
	created := make(map[backupTable]bool)
	createdNames := make(map[string]bool)
	for i := range stmts {
		for _, t := range createdRelations(stmts[i].node) {
			created[t] = true
			createdNames[t.name] = true
		}
	}
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
				if createdNames[rv.Relname] {
					reportNotYetCreated(s, rv, r)
				}
				continue
			}
		} else if (rv.Catalogname == "" || rv.Catalogname == schema.GetName()) && (created[table] || created[backupTable{name: rv.Relname}]) && !hasBackupRelation(findSchema(schema, rv.Schemaname), rv.Relname) {
			// An unqualified creation may be in this schema; either way
			// the table is not there when the backup runs.
			reportNotYetCreated(s, rv, r)
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
