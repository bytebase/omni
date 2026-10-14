package review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/oracle/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what certainly makes Bytebase's prior backup of an
// Oracle change fail, which fails the task before any of the change runs,
// judged from the SQL and the synced schema alone. The backup works like
// this (bytebase backend/runner/taskrun/database_migrate_executor.go,
// backupData, and backend/plugin/parser/plsql/backup.go,
// TransformDMLToSelect):
//
//   - It runs before the change, on the database the change runs in, with
//     the schema Bytebase last synced.
//   - It needs the BBDATAARCHIVE database, a schema, which it looks up
//     before anything else, and refuses a change longer than the size
//     limit.
//   - It backs up each top-level UPDATE and DELETE whose target is a
//     table, and skips one whose target is an inline view or a TABLE
//     collection, and MERGE, INSERT, and the DML of PL/SQL. A change with
//     nothing to back up fails the task: "no parse results".
//   - It groups the statements by table: the schema written, or else the
//     task's database, the name, and the database link. It refuses a group
//     with both UPDATE and DELETE, and a table referenced through a
//     database link.
//   - It copies a group with CREATE TABLE "BBDATAARCHIVE"."<prefix>_<table>_
//     <schema>" AS SELECT "<alias or table>".* FROM the statement's target
//     and WHERE clause, as written, and the FROM list of an UPDATE that has
//     a WHERE clause. A group of several statements selects
//     "<schema>"."<table>".* FROM "<schema>"."<table>" WHERE ROWID IN a
//     UNION of each statement's ROWIDs, so selected. It then runs COMMENT
//     ON TABLE on the copy, with the table's name inside the quoted
//     comment, unescaped.
//
// So the backup also fails on a table with a LONG or LONG RAW column,
// which CREATE TABLE AS SELECT rejects (ORA-00997); on a table whose name
// holds a single quote, which ends the comment's string; and on a name the
// change makes, which does not exist yet when the backup runs: a table the
// backup's query reads, or a column the change adds, renames to, or
// declares on a table it creates, read in the copied WHERE clause or a
// join condition of the copied FROM list.
//
// Oracle resolves an unqualified name through a public synonym when the
// schema holds no object of the name, and the snapshot holds no synonym:
// a table the change creates is judged missing only where the backup's
// query names it with its schema, as the outer query of a group of several
// statements does. A column the synced table lacks may be an invisible
// one, which Bytebase does not sync; then the change's own ADD of it
// fails (ORA-01430). A name the synced schema does not hold, and the
// change does not make, is not reported either: the schema the review sees
// may predate an earlier spec of the same plan. Neither is a column inside
// a subquery, nor anything from the change's first ALTER SESSION SET
// CURRENT_SCHEMA on, after which an unqualified name of the change means
// another schema's object. The copy's name, cut to Oracle's identifier
// length, may collide with another's on a version before 12.2, which the
// review does not know, and a nested table column makes CREATE TABLE AS
// SELECT fail (ORA-22913), which the column's type does not tell.

// checkPriorBackupSize reports a change longer than prior backup takes.
// It needs the SQL alone.
func checkPriorBackupSize(sql string, change review.Change, r *reporter) {
	if change.MaxBackupSize > 0 && len(sql) > change.MaxBackupSize {
		r.report(review.PriorBackup, -1, review.Range{}, fmt.Sprintf("the change is %d bytes, over the %d-byte limit of prior backup", len(sql), change.MaxBackupSize))
	}
}

// backupDML is an UPDATE or DELETE the backup copies the rows of.
type backupDML struct {
	pos   int
	stmt  *statement
	kind  string
	table *ast.ObjectName
	alias *ast.Alias
	where ast.ExprNode
	// from is the FROM list of an UPDATE, which the backup copies only
	// with a WHERE clause.
	from *ast.List
}

// backupDMLs returns the statements the backup copies the rows of: each
// top-level UPDATE and DELETE whose target is a table.
//
// bytebase: backend/plugin/parser/plsql/backup.go — prepareTransformation, extractOmniDML
func backupDMLs(stmts []statement) []backupDML {
	var out []backupDML
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.UpdateStmt:
			if v.Table != nil && v.Table.Name != "" {
				d := backupDML{pos: i, stmt: s, kind: "UPDATE", table: v.Table, alias: v.Alias, where: v.WhereClause}
				if v.WhereClause != nil {
					d.from = v.FromClause
				}
				out = append(out, d)
			}
		case *ast.DeleteStmt:
			if v.Table != nil && v.Table.Name != "" {
				out = append(out, backupDML{pos: i, stmt: s, kind: "DELETE", table: v.Table, alias: v.Alias, where: v.WhereClause})
			}
		}
	}
	return out
}

// groupKey is the key Bytebase groups a statement's table by: the schema
// written, or else the task's database, db, and the name and database
// link. db is "" when the target's database is not known: then a name
// without a schema matches only another without one.
//
// bytebase: backend/plugin/parser/plsql/backup.go — omniDMLTableReference, generateSQL
func groupKey(t *ast.ObjectName, db string) string {
	schema := t.Schema
	if schema == "" {
		schema = db
	}
	return schema + "." + t.Name + "@" + t.DBLink
}

// checkPriorBackupStatements reports, from the SQL alone, what the backup
// refuses or cannot run: a change with no statement to back up, a table
// referenced through a database link, and a table whose name holds a
// single quote, which the comment Bytebase sets on the copy does not
// escape.
//
// bytebase: backend/plugin/parser/plsql/backup.go — prepareTransformation (errNoBackupableDML), generateSQLForTable
// bytebase: backend/runner/taskrun/database_migrate_executor.go — backupData (COMMENT ON TABLE)
func checkPriorBackupStatements(stmts []statement, r *reporter) {
	dmls := backupDMLs(stmts)
	if len(dmls) == 0 {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup copies the rows of top-level UPDATE and DELETE statements on a table, and the change has none, which fails the backup")
		return
	}
	quoted := make(map[string]bool)
	for _, d := range dmls {
		switch {
		case d.table.DBLink != "":
			r.report(review.PriorBackup, d.stmt.index, rangeOf(d.stmt.loc), "prior backup does not support table "+objectName(d.table)+", referenced through a database link")
		case strings.Contains(d.table.Name, "'") && !quoted[groupKey(d.table, "")]:
			quoted[groupKey(d.table, "")] = true
			r.report(review.PriorBackup, d.stmt.index, rangeOf(d.stmt.loc), "prior backup writes the name of table "+objectName(d.table)+" unescaped into a quoted comment, which its single quote breaks")
		}
	}
}

// checkPriorBackupTarget reports the target-dependent reasons the backup
// fails: a missing BBDATAARCHIVE database, which Bytebase looks up before
// anything else; a table both an UPDATE and a DELETE of the change touch;
// a table with a LONG or LONG RAW column; and a name the backup's query of
// a statement before the change's first ALTER SESSION SET CURRENT_SCHEMA
// reads that only the change makes.
//
// bytebase: backend/plugin/parser/plsql/backup.go — generateSQL, generateSQLForTable, writeBranchSelect
// bytebase: backend/plugin/advisor/oracle/rule_builtin_prior_backup_check.go — oracleFindLongColumns
func checkPriorBackupTarget(stmts []statement, backupExists bool, x schemaIndex, m changeMade, r *reporter) {
	if !backupExists {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup needs database "+backupDatabase+", which does not exist")
	}
	db := ""
	if x.schema != nil && x.schema.GetName() != "" {
		db = x.schema.GetName()
	}
	// The comment Bytebase sets on each copy holds the task's database
	// name unescaped, in quotes.
	if strings.Contains(db, "'") && len(backupDMLs(stmts)) > 0 {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup writes the database name "+escapeLines(db)+" unescaped into a quoted comment, which its single quote breaks")
	}
	type group struct {
		first    backupDML
		size     int
		reported bool
	}
	dmls := backupDMLs(stmts)
	groups := make(map[string]*group)
	var order []string
	for _, d := range dmls {
		if d.table.DBLink != "" {
			continue
		}
		key := groupKey(d.table, db)
		g, seen := groups[key]
		if !seen {
			groups[key] = &group{first: d, size: 1}
			order = append(order, key)
			continue
		}
		g.size++
		if g.first.kind != d.kind && !g.reported {
			g.reported = true
			r.report(review.PriorBackup, d.stmt.index, rangeOf(d.stmt.loc), "prior backup cannot back up both UPDATE and DELETE on "+objectName(d.table))
		}
	}
	if x.schema == nil {
		return
	}
	for _, key := range order {
		g := groups[key]
		k := tableKey{schema: x.db(), name: g.first.table.Name}
		if g.first.table.Schema != "" {
			k.schema = g.first.table.Schema
		}
		if columns := longColumns(x, k); len(columns) > 0 {
			r.report(review.PriorBackup, g.first.stmt.index, rangeOf(g.first.stmt.loc), "prior backup copies "+objectName(g.first.table)+" with CREATE TABLE AS SELECT, which rejects its LONG or LONG RAW column "+strings.Join(columns, ", ")+" (ORA-00997)")
		}
	}
	limit := firstSessionChange(stmts)
	for _, d := range dmls {
		if d.pos >= limit || d.table.DBLink != "" {
			continue
		}
		size := groups[groupKey(d.table, db)].size
		if msg := backupReadsMade(d, size, x, m); msg != "" {
			r.report(review.PriorBackup, d.stmt.index, rangeOf(d.stmt.loc), msg)
		}
	}
}

// longColumns returns the LONG and LONG RAW columns of the synced table,
// or updatable view, a key names.
func longColumns(x schemaIndex, k tableKey) []string {
	columns := x.table(k).GetColumns()
	if columns == nil {
		for _, v := range x.schemaOf(k).GetViews() {
			if v.GetName() == k.name {
				columns = v.GetColumns()
			}
		}
	}
	var out []string
	for _, c := range columns {
		switch strings.ToUpper(strings.TrimSpace(c.GetType())) {
		case "LONG", "LONG RAW":
			out = append(out, ident(c.GetName()))
		}
	}
	return out
}

// changeMade is what the change makes, which the backup, run before all of
// it, cannot see: the tables, views, and synonyms it creates or renames
// to, and the columns it adds, renames to, or declares on a table it
// creates, by the table. A statement that may make nothing, IF NOT EXISTS
// or OR REPLACE, is left out, and so is everything from the first ALTER
// SESSION SET CURRENT_SCHEMA on.
type changeMade struct {
	relations []tableKey
	columns   []madeColumn
}

type madeColumn struct {
	table  tableKey
	column string
}

// madeByChange collects what the change makes, each table keyed by the
// schema the statement names, "" for none: the review compares it with a
// reference once it knows the target's database.
func madeByChange(stmts []statement) changeMade {
	var m changeMade
	stmts = stmts[:firstSessionChange(stmts)]
	key := func(ref *ast.ObjectName) tableKey { return tableKey{schema: ref.Schema, name: ref.Name} }
	addTable := func(v *ast.CreateTableStmt) {
		if v.Name == nil || v.IfNotExists || v.Name.DBLink != "" {
			return
		}
		k := key(v.Name)
		m.relations = append(m.relations, k)
		if v.Columns != nil {
			for _, item := range v.Columns.Items {
				if col, ok := item.(*ast.ColumnDef); ok {
					m.columns = append(m.columns, madeColumn{k, col.Name})
				}
			}
		}
	}
	for i := range stmts {
		switch v := stmts[i].node.(type) {
		case *ast.CreateTableStmt:
			addTable(v)
		case *ast.CreateSchemaStmt:
			if v.Stmts == nil {
				continue
			}
			for _, item := range v.Stmts.Items {
				switch e := item.(type) {
				case *ast.CreateTableStmt:
					addTable(e)
				case *ast.CreateViewStmt:
					if e.Name != nil {
						m.relations = append(m.relations, key(e.Name))
					}
				}
			}
		case *ast.CreateViewStmt:
			if v.Name != nil && !v.OrReplace && !v.IfNotExists {
				m.relations = append(m.relations, key(v.Name))
			}
		case *ast.CreateJsonDualityViewStmt:
			if v.Name != nil && !v.OrReplace && !v.IfNotExists {
				m.relations = append(m.relations, key(v.Name))
			}
		case *ast.CreateSynonymStmt:
			if v.Name != nil && !v.OrReplace && !v.Public {
				m.relations = append(m.relations, key(v.Name))
			}
		case *ast.RenameStmt:
			if v.NewName != nil {
				m.relations = append(m.relations, key(v.NewName))
			}
		case *ast.AlterTableStmt:
			if v.Name == nil || v.Actions == nil {
				continue
			}
			k := key(v.Name)
			for _, item := range v.Actions.Items {
				cmd, ok := item.(*ast.AlterTableCmd)
				if !ok {
					continue
				}
				switch cmd.Action {
				case ast.AT_ADD_COLUMN:
					for _, col := range commandColumns(cmd) {
						m.columns = append(m.columns, madeColumn{k, col.Name})
					}
				case ast.AT_RENAME_COLUMN:
					m.columns = append(m.columns, madeColumn{k, cmd.NewName})
				case ast.AT_RENAME:
					if cmd.NewName != "" {
						m.relations = append(m.relations, tableKey{schema: k.schema, name: cmd.NewName})
					}
				}
			}
		}
	}
	return m
}

// sameTable reports whether a name the change makes and a reference of
// the backup's query name the same table: the same name, in the same
// schema, a schema left out meaning the target's database.
func sameTable(made tableKey, ref *ast.ObjectName, x schemaIndex) bool {
	own := func(schema string) bool { return schema == "" || schema == x.db() }
	return made.name == ref.Name && (made.schema == ref.Schema || own(made.schema) && own(ref.Schema))
}

func (m changeMade) madeRelation(ref *ast.ObjectName, x schemaIndex) bool {
	for _, made := range m.relations {
		if sameTable(made, ref, x) {
			return true
		}
	}
	return false
}

func (m changeMade) madeColumnOn(ref *ast.ObjectName, column string, x schemaIndex) bool {
	for _, made := range m.columns {
		if made.column == column && sameTable(made.table, ref, x) {
			return true
		}
	}
	return false
}

// syncedKey keys a reference in the synced schema: the schema it names,
// or the target's database.
func syncedKey(ref *ast.ObjectName, x schemaIndex) tableKey {
	if ref.Schema != "" {
		return tableKey{schema: ref.Schema, name: ref.Name}
	}
	return tableKey{schema: x.db(), name: ref.Name}
}

// backupReadsMade returns why the backup's query of an UPDATE or DELETE
// fails on a name the change makes, "" for none: a table the query names
// with the database's schema, or the target of a group of several
// statements, which the outer query names so, that the synced schema does
// not hold; or a column the change makes on one of the statement's tables
// that the synced table lacks, read in the WHERE clause or a join
// condition.
func backupReadsMade(d backupDML, groupSize int, x schemaIndex, m changeMade) string {
	missing := func(ref *ast.ObjectName, qualified bool) bool {
		if ref == nil || ref.DBLink != "" || !qualified {
			return false
		}
		return x.lacksRelation(syncedKey(ref, x)) && m.madeRelation(ref, x)
	}
	if missing(d.table, d.table.Schema == x.db() || groupSize > 1) {
		return "prior backup runs before the change, when " + objectName(d.table) + " does not exist yet"
	}
	var reads []ast.Node
	if d.from != nil {
		reads = append(reads, d.from)
	}
	if d.where != nil {
		reads = append(reads, d.where)
	}
	var found *ast.ObjectName
	for _, n := range reads {
		ast.Inspect(n, func(n ast.Node) bool {
			if t, ok := n.(*ast.TableRef); ok && found == nil && t.Name != nil && missing(t.Name, t.Name.Schema == x.db()) {
				found = t.Name
			}
			return found == nil
		})
	}
	if found != nil {
		return "prior backup runs before the change, when " + objectName(found) + " does not exist yet"
	}
	return columnMissing(d, x, m)
}

// backupSource is a table the backup's query reads columns from, with the
// columns the synced schema gives it, nil when they are not known.
type backupSource struct {
	name    string
	ref     *ast.ObjectName
	columns map[string]bool
}

// columnMissing reports a column the backup's query reads that the change
// makes on one of the statement's tables, and that the synced table lacks.
// Only a name that can mean nothing else is judged: a column qualified by
// a table of the query, or an unqualified one when every table the query
// reads is one whose columns the snapshot gives and none has it.
// Subqueries are not read.
func columnMissing(d backupDML, x schemaIndex, m changeMade) string {
	var sources []backupSource
	known := true
	var quals []ast.Node
	if d.where != nil {
		quals = append(quals, d.where)
	}
	addRef := func(ref *ast.ObjectName, alias *ast.Alias) {
		src := backupSource{name: ref.Name, ref: ref}
		if alias != nil && alias.Name != "" {
			src.name = alias.Name
		}
		if t := x.table(syncedKey(ref, x)); t != nil && ref.DBLink == "" {
			src.columns = make(map[string]bool)
			for _, c := range t.GetColumns() {
				src.columns[c.GetName()] = true
			}
		} else {
			known = false
		}
		sources = append(sources, src)
	}
	addRef(d.table, d.alias)
	var flatten func(n ast.Node)
	flatten = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.TableRef:
			if v.Name != nil {
				addRef(v.Name, v.Alias)
				return
			}
			known = false
		case *ast.JoinClause:
			flatten(v.Left)
			flatten(v.Right)
			if v.On != nil {
				quals = append(quals, v.On)
			}
		default:
			known = false
		}
	}
	if d.from != nil {
		for _, item := range d.from.Items {
			flatten(item)
		}
	}
	var msg string
	judge := func(src backupSource, column string) {
		if src.columns == nil || src.columns[column] || src.ref == nil {
			return
		}
		if m.madeColumnOn(src.ref, column, x) {
			msg = "prior backup runs before the change, when column " + ident(column) + " of " + objectName(src.ref) + " does not exist yet"
		}
	}
	for _, q := range quals {
		ast.Inspect(q, func(n ast.Node) bool {
			if msg != "" {
				return false
			}
			switch v := n.(type) {
			case *ast.SubqueryExpr, *ast.SelectStmt, *ast.ExistsExpr:
				return false
			case *ast.ColumnRef:
				if v.Column == "" || v.Column == "*" {
					return false
				}
				if v.Table != "" {
					for _, src := range sources {
						if src.name == v.Table && (v.Schema == "" || src.ref.Schema == v.Schema || src.ref.Schema == "" && v.Schema == x.db()) {
							judge(src, v.Column)
							break
						}
					}
					return false
				}
				if !known {
					return false
				}
				for _, src := range sources {
					if src.columns[v.Column] {
						return false
					}
				}
				for _, src := range sources {
					if judge(src, v.Column); msg != "" {
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
