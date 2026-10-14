package review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what certainly makes Bytebase's prior backup of a
// MariaDB change fail, which fails the task before any of the change runs,
// judged from the SQL and the synced schema alone. Bytebase backs MariaDB
// up with its MySQL code, which reads each statement with the MySQL
// parser; a statement only MariaDB accepts fails that parse, which this
// package, holding no MySQL parser, does not judge. The backup works like
// this (bytebase backend/runner/taskrun/database_migrate_executor.go,
// backupData, and backend/plugin/parser/mysql/backup.go,
// TransformDMLToSelect, registered for MariaDB too):
//
//   - It runs before the change, on the database the change runs in, with
//     the schema Bytebase last synced; a USE in the change does not move
//     it.
//   - It needs the bbdataarchive database, which it looks up before
//     anything else, and refuses a change longer than the size limit.
//   - It copies the rows each top-level UPDATE and DELETE would touch, one
//     table at a time: CREATE TABLE bbdataarchive.… LIKE the table, then
//     INSERT INTO it SELECT … FROM the statement's tables, as written, and
//     its WHERE, ORDER BY, and LIMIT, as written.
//   - It refuses UPDATE and DELETE on the same table, written the same
//     way once a name without a database takes the task's, and fails on a
//     table the synced schema does not hold.
//
// So the backup fails on a name the change makes, which does not exist
// yet when it runs: a table the statement reads, or a column the change
// adds, renames to, or declares on a table it creates, read in the WHERE
// clause, a join condition or USING list, or ORDER BY.
//
// A name the synced schema does not hold, and the change does not make, is
// not reported: the schema the review sees may predate an earlier spec of
// the same plan, and Bytebase syncs again before this spec's backup.
// Neither is a statement with a WITH clause, whose names may mean its own
// queries, a name of another database, a column inside a subquery or of a
// source whose columns the snapshot does not give, or anything after a
// USE in the change.

// checkPriorBackupSize reports a change longer than prior backup takes.
// It needs the SQL alone.
func checkPriorBackupSize(sql string, change review.Change, r *reporter) {
	if change.MaxBackupSize > 0 && len(sql) > change.MaxBackupSize {
		r.report(review.PriorBackup, -1, review.Range{}, fmt.Sprintf("the change is %d bytes, over the %d-byte limit of prior backup", len(sql), change.MaxBackupSize))
	}
}

// checkPriorBackupMixed reports a table the backup refuses: one an UPDATE
// and a DELETE of the change both touch, written the same way, a name
// without a database taken as one of the target's, db, as Bytebase fills
// it in. db is "" when the target's name is not known: then only names
// written alike are compared. Only a statement whose targets are certain
// is judged.
//
// bytebase: backend/plugin/parser/mysql/backup.go — collectSingleTables,
// generateSQL
func checkPriorBackupMixed(stmts []statement, db string, r *reporter) {
	kinds := make(map[string]string)
	reported := make(map[string]bool)
	for i := range stmts {
		s := &stmts[i]
		refs, kind := backupTargets(s)
		for _, ref := range refs {
			key := ref.Schema
			if key == "" {
				key = db
			}
			key += "." + ref.Name
			first, seen := kinds[key]
			if !seen {
				kinds[key] = kind
				continue
			}
			if first != kind && !reported[key] {
				reported[key] = true
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup cannot back up both UPDATE and DELETE on "+table(ref))
			}
		}
	}
}

// checkPriorBackupStatements reports, from the SQL alone, what makes
// Bytebase refuse the change: a statement with RETURNING, which the MySQL
// parser Bytebase reads every statement of a MariaDB change with rejects,
// and a multi-table DELETE whose target the backup cannot find in the
// statement's FROM list, which Bytebase matches by name or alias exactly,
// case included. A DELETE after a WITH clause, whose targets may be its
// queries, is not judged.
//
// bytebase: backend/plugin/parser/mysql/backup.go — prepareTransformation,
// extractTablesFromDelete
func checkPriorBackupStatements(stmts []statement, r *reporter) {
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.InsertStmt:
			if len(v.Returning) > 0 {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup reads the change as MySQL, which has no RETURNING clause")
			}
		case *ast.DeleteStmt:
			if len(v.Returning) > 0 {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup reads the change as MySQL, which has no RETURNING clause")
			}
		}
		v, ok := s.node.(*ast.DeleteStmt)
		if !ok || len(v.Using) == 0 || v.Loc.Start > s.loc.Start {
			continue
		}
		from := fromNames(v.Using)
		for _, e := range v.Tables {
			target, ok := e.(*ast.TableRef)
			if !ok {
				continue
			}
			name := target.Name
			if target.Alias != "" {
				name = target.Alias
			}
			if !from[name] {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup does not find DELETE target "+ident(name)+" in the FROM list, written the same way")
			}
		}
	}
}

// fromNames returns the names a FROM list gives its tables: the alias, or
// the table's name.
func fromNames(exprs []ast.TableExpr) map[string]bool {
	names := make(map[string]bool)
	var collect func(e ast.TableExpr)
	collect = func(e ast.TableExpr) {
		switch t := e.(type) {
		case *ast.TableRef:
			if t.Alias != "" {
				names[t.Alias] = true
			} else {
				names[t.Name] = true
			}
		case *ast.JoinClause:
			collect(t.Left)
			collect(t.Right)
		}
	}
	for _, e := range exprs {
		collect(e)
	}
	return names
}

// backupTargets returns the tables an UPDATE or DELETE backs up, with the
// statement's kind: the table of a single-table UPDATE or DELETE, and
// each target of a multi-table DELETE, as the table its name means in the
// FROM list, matched as Bytebase matches it. A multi-table UPDATE, whose
// targets depend on its SET columns, and a multi-table DELETE after a WITH
// clause, whose targets may be its queries, give none.
//
// bytebase: backend/plugin/parser/mysql/backup.go — extractTablesFromDelete,
// collectSingleTables
func backupTargets(s *statement) ([]*ast.TableRef, string) {
	switch v := s.node.(type) {
	case *ast.UpdateStmt:
		if len(v.Tables) == 1 {
			if ref, ok := v.Tables[0].(*ast.TableRef); ok {
				return []*ast.TableRef{ref}, "UPDATE"
			}
		}
	case *ast.DeleteStmt:
		if len(v.Using) == 0 {
			if len(v.Tables) == 1 {
				if ref, ok := v.Tables[0].(*ast.TableRef); ok {
					return []*ast.TableRef{ref}, "DELETE"
				}
			}
			return nil, ""
		}
		if v.Loc.Start > s.loc.Start {
			return nil, ""
		}
		from := make(map[string]*ast.TableRef)
		var collect func(e ast.TableExpr)
		collect = func(e ast.TableExpr) {
			switch t := e.(type) {
			case *ast.TableRef:
				key := t.Name
				if t.Alias != "" {
					key = t.Alias
				}
				from[key] = t
			case *ast.JoinClause:
				collect(t.Left)
				collect(t.Right)
			}
		}
		for _, e := range v.Using {
			collect(e)
		}
		var refs []*ast.TableRef
		for _, e := range v.Tables {
			target, ok := e.(*ast.TableRef)
			if !ok {
				continue
			}
			key := target.Name
			if target.Alias != "" {
				key = target.Alias
			}
			if ref, ok := from[key]; ok {
				refs = append(refs, ref)
			}
		}
		return refs, "DELETE"
	}
	return nil, ""
}

// changeMade is what the change makes, which the backup, run before all of
// it, cannot see: the tables and views it creates or renames to, and the
// columns it adds, renames to, or declares on a table it creates, by the
// table as written. A relation counts only where the synced schema lacks
// it, so IF NOT EXISTS and OR REPLACE make theirs; the columns of a CREATE
// TABLE IF NOT EXISTS, which may keep the table it finds, do not count.
// What the change makes after a USE is left out.
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
	for i := range stmts[:beforeUse(stmts)] {
		switch v := stmts[i].node.(type) {
		case *ast.CreateTableStmt:
			if v.Table == nil {
				continue
			}
			// IF NOT EXISTS makes the table when the synced schema lacks
			// it, which the reader checks, but may keep another one's
			// columns.
			m.relations = append(m.relations, v.Table)
			if v.IfNotExists {
				continue
			}
			for _, col := range v.Columns {
				if col != nil {
					m.columns = append(m.columns, madeColumn{v.Table, col.Name})
				}
			}
		case *ast.CreateViewStmt:
			if v.Name != nil {
				m.relations = append(m.relations, v.Name)
			}
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair != nil && pair.New != nil {
					m.relations = append(m.relations, pair.New)
				}
			}
		case *ast.AlterTableStmt:
			if v.Table == nil {
				continue
			}
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				switch cmd.Type {
				case ast.ATRenameTable:
					if cmd.NewTable != nil {
						// A new name without a database stays in the table's.
						to := *cmd.NewTable
						if to.Schema == "" {
							to.Schema = v.Table.Schema
						}
						m.relations = append(m.relations, &to)
					} else if cmd.NewName != "" {
						m.relations = append(m.relations, &ast.TableRef{Schema: v.Table.Schema, Name: cmd.NewName})
					}
				case ast.ATAddColumn:
					if cmd.Column != nil {
						m.columns = append(m.columns, madeColumn{v.Table, cmd.Column.Name})
					}
					for _, col := range cmd.Columns {
						if col != nil {
							m.columns = append(m.columns, madeColumn{v.Table, col.Name})
						}
					}
				case ast.ATChangeColumn:
					if cmd.Column != nil && !sameName(cmd.Name, cmd.Column.Name) {
						m.columns = append(m.columns, madeColumn{v.Table, cmd.Column.Name})
					}
				case ast.ATRenameColumn:
					m.columns = append(m.columns, madeColumn{v.Table, cmd.NewName})
				}
			}
		}
	}
	return m
}

// checkPriorBackupTarget reports the target-dependent reasons the backup
// fails: a missing bbdataarchive database, which Bytebase looks up before
// anything else, and a name the backup's query of a statement before the
// change's first USE reads that only the change makes.
func checkPriorBackupTarget(sql string, stmts []statement, t review.Target, x schemaIndex, m changeMade, r *reporter) {
	if !t.BackupDatabaseExists {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup needs database "+backupDatabase+", which does not exist")
	}
	if x.schema == nil {
		return
	}
	// The backup copies a table with CREATE TABLE ... LIKE, which refuses
	// a view. A name without a database means the task's, whatever a USE
	// does.
	for i := range stmts {
		s := &stmts[i]
		refs, _ := backupTargets(s)
		for _, ref := range refs {
			if x.inDatabase(ref) && x.isView(ref.Name) {
				r.report(review.PriorBackup, s.index, rangeOf(s.loc), "prior backup copies tables with CREATE TABLE ... LIKE, which refuses view "+table(ref))
			}
		}
	}
	for i := range stmts[:beforeUse(stmts)] {
		s := &stmts[i]
		var loc ast.Loc
		switch v := s.node.(type) {
		case *ast.UpdateStmt:
			loc = v.Loc
		case *ast.DeleteStmt:
			loc = v.Loc
		default:
			continue
		}
		// A WITH clause before the statement, which the parser does not
		// keep, may name queries the statement reads.
		if loc.Start > s.loc.Start {
			continue
		}
		if msg := backupReadsMade(s.node, x, m); msg != "" {
			r.report(review.PriorBackup, s.index, rangeOf(s.loc), msg)
		}
	}
}

// madeRelation reports whether the change makes a table or view the
// reference may mean in the target's database.
func (m changeMade) madeRelation(x schemaIndex, ref *ast.TableRef) bool {
	for _, made := range m.relations {
		if x.sameRelation(made, ref) {
			return true
		}
	}
	return false
}

// madeColumnOn reports whether the change makes the column on the table
// the reference means.
func (m changeMade) madeColumnOn(x schemaIndex, ref *ast.TableRef, column string) bool {
	for _, made := range m.columns {
		if strings.EqualFold(made.column, column) && x.sameRelation(made.table, ref) {
			return true
		}
	}
	return false
}

// sameRelation reports whether two references of the change mean the same
// table of the target's database: a reference without a database means
// the target's.
func (x schemaIndex) sameRelation(a, b *ast.TableRef) bool {
	if a == nil || b == nil || !x.same(a.Name, b.Name) {
		return false
	}
	db := func(ref *ast.TableRef) string {
		if ref.Schema == "" {
			return x.schema.GetName()
		}
		return ref.Schema
	}
	return x.same(db(a), db(b))
}

// inDatabase reports whether a reference means a table of the target's
// database.
func (x schemaIndex) inDatabase(ref *ast.TableRef) bool {
	return ref.Schema == "" || x.same(ref.Schema, x.schema.GetName())
}

// backupReadsMade returns why the backup's query of an UPDATE or DELETE
// fails on a name the change makes, "" for none: a table of the statement,
// or of a subquery of its WHERE clause, the synced schema does not hold;
// or a column the change makes that none of the statement's tables has
// before the change, read in the WHERE clause, a join condition or USING
// list, or ORDER BY.
func backupReadsMade(n ast.Node, x schemaIndex, m changeMade) string {
	var from []ast.TableExpr
	var quals []ast.Node
	switch v := n.(type) {
	case *ast.UpdateStmt:
		from = v.Tables
		if v.Where != nil {
			quals = append(quals, v.Where)
		}
		for _, o := range v.OrderBy {
			if o != nil {
				quals = append(quals, o)
			}
		}
	case *ast.DeleteStmt:
		from = v.Tables
		if len(v.Using) > 0 {
			from = v.Using
		}
		if v.Where != nil {
			quals = append(quals, v.Where)
		}
		for _, o := range v.OrderBy {
			if o != nil {
				quals = append(quals, o)
			}
		}
	}
	// A table the query reads that does not exist yet.
	var missing *ast.TableRef
	check := func(ref *ast.TableRef) {
		if missing == nil && x.inDatabase(ref) && m.madeRelation(x, ref) && !x.hasRelation(ref.Name) {
			missing = ref
		}
	}
	for _, e := range from {
		ast.Inspect(e, func(n ast.Node) bool {
			if ref, ok := n.(*ast.TableRef); ok {
				check(ref)
			}
			return true
		})
	}
	for _, q := range quals {
		ast.Inspect(q, func(n ast.Node) bool {
			if ref, ok := n.(*ast.TableRef); ok {
				check(ref)
			}
			return true
		})
	}
	if missing != nil {
		return "prior backup runs before the change, when " + table(missing) + " does not exist yet"
	}
	return columnMissing(from, quals, x, m)
}

// backupSource is a table the backup's query reads columns from, with
// the columns the synced schema gives it, nil when they are not known.
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
// Subqueries are not read, and neither is a derived table's column.
func columnMissing(from []ast.TableExpr, quals []ast.Node, x schemaIndex, m changeMade) string {
	var sources []backupSource
	known := true
	// using is the columns a JOIN ... USING names, each with the sources
	// of one side, every one of which may hold it.
	type usingColumn struct {
		column string
		side   []backupSource
	}
	var using []usingColumn
	// flatten adds the sources of a table expression and returns them.
	var flatten func(e ast.TableExpr) []backupSource
	flatten = func(e ast.TableExpr) []backupSource {
		switch v := e.(type) {
		case *ast.TableRef:
			src := backupSource{name: v.Name, ref: v}
			if v.Alias != "" {
				src.name = v.Alias
			}
			if t := x.table(v, false); t != nil {
				src.columns = make(map[string]bool)
				for _, c := range t.GetColumns() {
					src.columns[strings.ToLower(c.GetName())] = true
				}
			} else {
				known = false
			}
			sources = append(sources, src)
			return []backupSource{src}
		case *ast.JoinClause:
			left, right := flatten(v.Left), flatten(v.Right)
			switch c := v.Condition.(type) {
			case *ast.OnCondition:
				if c.Expr != nil {
					quals = append(quals, c.Expr)
				}
			case *ast.UsingCondition:
				for _, column := range c.Columns {
					using = append(using, usingColumn{column, left}, usingColumn{column, right})
				}
			}
			return append(left, right...)
		default:
			// A source whose columns are not known, which may hold any.
			known = false
			return []backupSource{{}}
		}
	}
	for _, e := range from {
		flatten(e)
	}
	var msg string
	judge := func(src backupSource, column string) {
		if src.columns == nil || src.columns[strings.ToLower(column)] || src.ref == nil {
			return
		}
		if m.madeColumnOn(x, src.ref, column) {
			msg = "prior backup runs before the change, when column " + ident(column) + " of " + table(src.ref) + " does not exist yet"
		}
	}
	// USING names a column each side has: one no source of a side may
	// hold is missing there.
	for _, u := range using {
		has := false
		for _, src := range u.side {
			has = has || src.columns == nil || src.columns[strings.ToLower(u.column)]
		}
		for _, src := range u.side {
			if !has && msg == "" {
				judge(src, u.column)
			}
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
				if v.Star || v.Column == "" {
					return false
				}
				if v.Table != "" {
					for _, src := range sources {
						if !x.same(src.name, v.Table) {
							continue
						}
						if v.Schema != "" && (src.ref == nil || src.ref.Alias != "" || !x.same(dbOf(x, src.ref), v.Schema)) {
							continue
						}
						judge(src, v.Column)
						break
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

// dbOf returns the database a reference names, the target's when it names
// none.
func dbOf(x schemaIndex, ref *ast.TableRef) string {
	if ref.Schema == "" {
		return x.schema.GetName()
	}
	return ref.Schema
}
