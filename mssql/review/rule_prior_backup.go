package review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what certainly makes Bytebase's prior backup of a
// SQL Server change fail, which fails the task before any of the change
// runs, judged from the SQL and the synced schema alone. The backup works
// like this (bytebase backend/runner/taskrun/database_migrate_executor.go,
// backupData, and backend/plugin/parser/tsql/backup.go):
//
//   - It runs before the change, on the database the change runs in, with
//     the schema Bytebase last synced; a USE in the change does not move
//     it.
//   - It needs the bbdataarchive database, which it looks up before
//     anything else, and refuses a change longer than the size limit.
//   - It copies the rows each top-level UPDATE and DELETE would touch, one
//     table at a time: SELECT * INTO bbdataarchive.dbo.… FROM a UNION of
//     one SELECT per statement, which copies, as written, the statement's
//     WITH clause, its TOP, its FROM list or its table, and its WHERE and
//     OPTION clauses.
//   - It refuses a table variable or a WITH query as the target, UPDATE
//     and DELETE on one table, WITH XMLNAMESPACES, statements on one table
//     with different WITH clauses, and WHERE CURRENT OF, a name without a
//     database or schema taking the task's database and dbo. A temporary
//     table is not backed up.
//
// So the backup also fails on TOP, which it writes after the select list,
// and on a name the change makes, which does not exist yet when it runs:
// a table the statement reads, or a column the change adds or declares on
// a table it creates, read in the WHERE clause or a join condition.
//
// A name the synced schema does not hold, and the change does not make, is
// not reported: the schema the review sees may predate an earlier spec of
// the same plan. Neither is a name of another database, an unqualified
// name more than one schema holds, which resolves by the session's
// default schema, a column inside a subquery, or anything from the
// change's first USE on.

// checkPriorBackupSize reports a change longer than prior backup takes.
// It needs the SQL alone.
func checkPriorBackupSize(sql string, change review.Change, r *reporter) {
	if change.MaxBackupSize > 0 && len(sql) > change.MaxBackupSize {
		r.report(review.PriorBackup, -1, review.Range{}, fmt.Sprintf("the change is %d bytes, over the %d-byte limit of prior backup", len(sql), change.MaxBackupSize))
	}
}

// dmlTarget is what the backup reads of an UPDATE or DELETE.
type dmlTarget struct {
	kind   string
	with   *ast.WithClause
	top    *ast.TopClause
	target ast.TableExpr
	from   *ast.List
	where  ast.ExprNode
	loc    ast.Loc
}

func dmlOf(n ast.Node) (dmlTarget, bool) {
	switch v := n.(type) {
	case *ast.UpdateStmt:
		return dmlTarget{"UPDATE", v.WithClause, v.Top, v.Relation, v.FromClause, v.WhereClause, v.Loc}, true
	case *ast.DeleteStmt:
		return dmlTarget{"DELETE", v.WithClause, v.Top, v.Relation, v.FromClause, v.WhereClause, v.Loc}, true
	}
	return dmlTarget{}, false
}

// backupTable returns the DML of a statement and the table its backup
// copies, or why the backup refuses the statement: a table variable or a
// WITH query as the target, WHERE CURRENT OF, or WITH XMLNAMESPACES. ref
// is nil, with no refusal, for a statement the backup skips, a temporary
// table's, or one whose target the rule does not know.
//
// bytebase: backend/plugin/parser/tsql/backup.go — prepareTransformation,
// generateSQLForTable
func backupTable(s *statement) (d dmlTarget, ref *ast.TableRef, refusal string) {
	d, ok := dmlOf(s.node)
	if !ok {
		return d, nil, ""
	}
	switch t := d.target.(type) {
	case *ast.TableVarRef:
		return d, nil, "prior backup does not support the table variable " + t.Name + " as the target"
	case *ast.TableRef:
		ref = backedUp(t, d.from)
	}
	switch {
	case ref == nil || strings.HasPrefix(ref.Object, "#"):
		return d, nil, ""
	case ref.Schema == "" && ref.Database == "" && withNames(d.with)[strings.ToLower(ref.Object)]:
		return d, nil, "prior backup does not support the WITH query " + ident(ref.Object) + " as the target"
	}
	if _, ok := d.where.(*ast.CurrentOfExpr); ok {
		return d, nil, "prior backup does not support WHERE CURRENT OF"
	}
	if d.with != nil && d.with.XmlNamespaces != nil && len(d.with.XmlNamespaces.Items) > 0 {
		return d, nil, "prior backup does not support WITH XMLNAMESPACES"
	}
	return d, ref, ""
}

// checkPriorBackupStatements reports, from the SQL alone, what the backup
// refuses or cannot copy in a statement: what backupTable refuses, a
// FROM list that does not name the table copied, which Bytebase's query
// selects from it, and TOP.
func checkPriorBackupStatements(stmts []statement, r *reporter) {
	for i := range stmts {
		s := &stmts[i]
		d, ref, refusal := backupTable(s)
		switch {
		case refusal != "":
			r.report(review.PriorBackup, s.index, rangeOf(s.loc), refusal)
		case ref == nil:
		default:
			if !fromNames(d.from, ref) {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup selects "+object(ref)+".* from the statement's FROM list, which does not name it")
			}
			if d.top != nil {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup copies TOP after the select list of its query, which SQL Server rejects")
			}
		}
	}
}

// checkPriorBackupGroups reports the tables the backup refuses to group:
// one both an UPDATE and a DELETE touch, and one statements with different
// WITH clauses touch. The table is keyed as Bytebase keys it: the database
// and schema written, or the target's database, db, and dbo, and the name,
// as written. db is "" when the target's name is not known: then a name
// without a database matches only another without one.
//
// bytebase: backend/plugin/parser/tsql/backup.go — prepareTransformation,
// generateSQL
func checkPriorBackupGroups(sql string, stmts []statement, db string, r *reporter) {
	type group struct {
		kind     string
		with     string
		reported bool
	}
	groups := make(map[string]*group)
	for i := range stmts {
		s := &stmts[i]
		d, ref, refusal := backupTable(s)
		if ref == nil || refusal != "" {
			continue
		}
		report := func(msg string) {
			r.report(review.PriorBackup, s.index, rangeOf(s.loc), msg)
		}
		database := ref.Database
		if database == "" {
			database = db
		}
		schema := ref.Schema
		if schema == "" {
			schema = "dbo"
		}
		key := database + "." + schema + "." + ref.Object
		with := ""
		if d.with != nil {
			with = textOf(sql, d.with.Loc)
		}
		g, seen := groups[key]
		if !seen {
			groups[key] = &group{kind: d.kind, with: with}
			continue
		}
		switch {
		case g.reported:
		case g.kind != d.kind:
			g.reported = true
			report("prior backup cannot back up both UPDATE and DELETE on " + object(ref))
		case g.with != with:
			g.reported = true
			report("prior backup cannot back up statements on " + object(ref) + " with different WITH clauses")
		}
	}
}

// fromNames reports whether a FROM list may name the table the backup
// copies, which Bytebase's query selects from it: a table of its name, or
// an alias of it. A list that holds anything but tables and joins of them
// may, and so may no list.
func fromNames(from *ast.List, ref *ast.TableRef) bool {
	if from == nil || len(from.Items) == 0 {
		return true
	}
	names, plain := false, true
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.TableRef:
			names = names || sameName(v.Object, ref.Object) || sameName(v.Alias, ref.Object)
		case *ast.JoinClause:
			walk(v.Left)
			walk(v.Right)
		default:
			plain = false
		}
	}
	for _, item := range from.Items {
		walk(item)
	}
	return names || !plain
}

// backedUp returns the table the backup copies for a target: the table a
// FROM item gives the target's name as its alias, when the target names
// no schema or database, or else the target itself. Bytebase matches the
// alias as written.
func backedUp(target *ast.TableRef, from *ast.List) *ast.TableRef {
	if from == nil || target.Schema != "" || target.Database != "" {
		return target
	}
	var found *ast.TableRef
	var walk func(n ast.Node)
	walk = func(n ast.Node) {
		if found != nil {
			return
		}
		switch v := n.(type) {
		case *ast.TableRef:
			if v.Alias == target.Object {
				found = v
			}
		case *ast.AliasedTableRef:
			if t, ok := v.Table.(*ast.TableRef); ok && v.Alias == target.Object {
				found = t
			}
		case *ast.JoinClause:
			walk(v.Left)
			walk(v.Right)
		}
	}
	for _, item := range from.Items {
		walk(item)
	}
	if found == nil {
		return target
	}
	return found
}

// withNames returns the names of a WITH clause's queries, lower-cased.
func withNames(w *ast.WithClause) map[string]bool {
	names := make(map[string]bool)
	if w == nil || w.CTEs == nil {
		return names
	}
	for _, item := range w.CTEs.Items {
		if c, ok := item.(*ast.CommonTableExpr); ok {
			names[strings.ToLower(c.Name)] = true
		}
	}
	return names
}

func textOf(sql string, loc ast.Loc) string {
	if loc.Start < 0 || loc.End > len(sql) || loc.End <= loc.Start {
		return ""
	}
	return sql[loc.Start:loc.End]
}

// changeMade is what the change makes, which the backup, run before all of
// it, cannot see: the tables, views, and synonyms it creates, SELECT INTO,
// CREATE SCHEMA, and BEGIN ... END among them, and the columns it adds,
// renames to with sp_rename, or declares on a table it creates, by the
// table as written. A statement that may make nothing, IF NOT EXISTS, is
// left out, and so is what the change makes after a USE.
type changeMade struct {
	relations []*ast.TableRef
	columns   []madeColumn
}

type madeColumn struct {
	table  *ast.TableRef
	column string
}

func madeByChange(stmts []statement) changeMade {
	var m changeMade
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.CreateTableStmt:
			if v.IfNotExists || v.Name == nil {
				return
			}
			m.relations = append(m.relations, v.Name)
			if v.Columns != nil {
				for _, item := range v.Columns.Items {
					if col, ok := item.(*ast.ColumnDef); ok {
						m.columns = append(m.columns, madeColumn{v.Name, col.Name})
					}
				}
			}
		case *ast.CreateViewStmt:
			if v.Name != nil && !v.OrAlter {
				m.relations = append(m.relations, v.Name)
			}
		case *ast.SelectStmt:
			if v.IntoTable != nil {
				m.relations = append(m.relations, v.IntoTable)
			}
		case *ast.CreateSynonymStmt:
			if v.Name != nil {
				m.relations = append(m.relations, v.Name)
			}
		case *ast.CreateSchemaStmt:
			// The elements of CREATE SCHEMA run as the statement does.
			if v.Elements != nil {
				for _, e := range v.Elements.Items {
					visit(e)
				}
			}
		case *ast.ExecStmt:
			// A column sp_rename renames to, in the target's database.
			if name, newName, kind, ok := spRename(v); ok && strings.EqualFold(kind, "COLUMN") && v.Name.Database == "" && v.Name.Server == "" {
				if ref, _ := splitColumn(name); ref != nil {
					m.columns = append(m.columns, madeColumn{ref, newName})
				}
			}
		case *ast.BeginEndStmt:
			// So do the statements of a BEGIN ... END block.
			if v.Stmts != nil {
				for _, e := range v.Stmts.Items {
					visit(e)
				}
			}
		case *ast.AlterTableStmt:
			if v.Actions == nil {
				return
			}
			for _, item := range v.Actions.Items {
				if a, ok := item.(*ast.AlterTableAction); ok && a.Type == ast.ATAddColumn && a.Column != nil {
					m.columns = append(m.columns, madeColumn{v.Name, a.Column.Name})
				}
			}
		}
	}
	for i := range stmts[:beforeUse(stmts)] {
		visit(stmts[i].node)
	}
	return m
}

func (m changeMade) madeRelation(ref *ast.TableRef) bool {
	for _, made := range m.relations {
		if sameObject(made, ref) {
			return true
		}
	}
	return false
}

func (m changeMade) madeColumnOn(ref *ast.TableRef, column string) bool {
	for _, made := range m.columns {
		if sameName(made.column, column) && sameObject(made.table, ref) {
			return true
		}
	}
	return false
}

// checkPriorBackupTarget reports the target-dependent reasons the backup
// fails: a missing bbdataarchive database, which Bytebase looks up before
// anything else, and a name the backup's query of a statement before the
// change's first USE reads that only the change makes.
func checkPriorBackupTarget(stmts []statement, backupExists bool, x schemaIndex, m changeMade, r *reporter) {
	if !backupExists {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup needs database "+backupDatabase+", which does not exist")
	}
	if x.schema == nil {
		return
	}
	for i := range stmts[:beforeUse(stmts)] {
		s := &stmts[i]
		d, ok := dmlOf(s.node)
		if !ok || d.with != nil {
			continue
		}
		if msg := backupReadsMade(d, x, m); msg != "" {
			r.report(review.PriorBackup, s.index, rangeOf(s.loc), msg)
		}
	}
}

// backupReadsMade returns why the backup's query of an UPDATE or DELETE
// fails on a name the change makes, "" for none: a table of the statement,
// or of a subquery of its WHERE clause, no synced schema it may mean
// holds; or a column the change makes that none of the statement's tables
// has before the change, read in the WHERE clause or a join condition.
func backupReadsMade(d dmlTarget, x schemaIndex, m changeMade) string {
	var from []ast.Node
	if d.from != nil && len(d.from.Items) > 0 {
		for _, item := range d.from.Items {
			from = append(from, item)
		}
	} else if d.target != nil {
		from = append(from, d.target)
	}
	var missing *ast.TableRef
	check := func(n ast.Node) bool {
		if ref, ok := n.(*ast.TableRef); ok && missing == nil && x.inDatabase(ref) && !strings.HasPrefix(ref.Object, "#") && m.madeRelation(ref) && !x.holdsRelation(ref) {
			missing = ref
		}
		return true
	}
	for _, n := range from {
		ast.Inspect(n, check)
	}
	if d.where != nil {
		ast.Inspect(d.where, check)
	}
	if missing != nil {
		return "prior backup runs before the change, when " + object(missing) + " does not exist yet"
	}
	return columnMissing(from, d.where, x, m)
}

// backupSource is a table the backup's query reads columns from, with the
// columns the synced schema gives it, nil when they are not known.
type backupSource struct {
	name    string
	ref     *ast.TableRef
	columns map[string]bool
}

// columnMissing reports a column the backup's query names that the change
// makes on one of the statement's tables, and that none of them has before
// the change. Only a name that can mean nothing else is judged: a column
// qualified by a table of the statement, or an unqualified one when every
// source is a table whose columns the snapshot gives and none has it.
// Subqueries are not read.
func columnMissing(from []ast.Node, where ast.ExprNode, x schemaIndex, m changeMade) string {
	var sources []backupSource
	known := true
	var quals []ast.Node
	if where != nil {
		quals = append(quals, where)
	}
	addRef := func(ref *ast.TableRef, alias string) {
		src := backupSource{name: ref.Object, ref: ref}
		if alias != "" {
			src.name = alias
		}
		if t := x.table(ref); t != nil {
			src.columns = make(map[string]bool)
			for _, c := range t.GetColumns() {
				src.columns[strings.ToLower(c.GetName())] = true
			}
		} else {
			known = false
		}
		sources = append(sources, src)
	}
	var flatten func(n ast.Node)
	flatten = func(n ast.Node) {
		switch v := n.(type) {
		case *ast.TableRef:
			addRef(v, v.Alias)
		case *ast.AliasedTableRef:
			if t, ok := v.Table.(*ast.TableRef); ok && v.Columns == nil {
				addRef(t, v.Alias)
				return
			}
			known = false
			sources = append(sources, backupSource{name: v.Alias})
		case *ast.JoinClause:
			flatten(v.Left)
			flatten(v.Right)
			if v.Condition != nil {
				quals = append(quals, v.Condition)
			}
		default:
			known = false
		}
	}
	for _, n := range from {
		flatten(n)
	}
	var msg string
	judge := func(src backupSource, column string) {
		if src.columns == nil || src.columns[strings.ToLower(column)] || src.ref == nil {
			return
		}
		if m.madeColumnOn(src.ref, column) {
			msg = "prior backup runs before the change, when column " + ident(column) + " of " + object(src.ref) + " does not exist yet"
		}
	}
	for _, q := range quals {
		ast.Inspect(q, func(n ast.Node) bool {
			if msg != "" {
				return false
			}
			switch v := n.(type) {
			case *ast.SubqueryExpr, *ast.SelectStmt:
				return false
			case *ast.ColumnRef:
				if v.Column == "" || v.Server != "" || v.Database != "" {
					return false
				}
				if v.Table != "" {
					for _, src := range sources {
						if sameName(src.name, v.Table) && (v.Schema == "" || src.ref != nil && sameName(src.ref.Schema, v.Schema)) {
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
					if src.columns[strings.ToLower(v.Column)] {
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
