package review

import (
	"fmt"
	"strings"

	"github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/review"
)

// PriorBackup reports what certainly makes Bytebase's prior backup of a
// MySQL change fail, which fails the task before any of the change runs,
// judged from the SQL and the synced schema alone. The backup works like
// this (bytebase backend/runner/taskrun/database_migrate_executor.go,
// backupData, and backend/plugin/parser/mysql/backup.go,
// TransformDMLToSelect):
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
//     way, and fails on a table the synced schema does not hold.
//
// So the backup fails on a name the change makes, which does not exist
// yet when it runs: a table the statement reads, or a column the change
// adds, renames to, or declares on a table it creates, read in the WHERE
// clause, a join condition, or ORDER BY.
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

// checkPriorBackupStatements reports, from the SQL alone, a table the
// backup refuses: one an UPDATE and a DELETE of the change both touch,
// written the same way. Only a statement on one table is judged, whose
// target is certain.
//
// bytebase: backend/plugin/parser/mysql/backup.go — generateSQL
func checkPriorBackupStatements(stmts []statement, r *reporter) {
	kinds := make(map[string]string)
	reported := make(map[string]bool)
	for i := range stmts {
		s := &stmts[i]
		ref, kind := backupTarget(s.node)
		if ref == nil {
			continue
		}
		key := ref.Schema + "." + ref.Name
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

// backupTarget returns the table an UPDATE or DELETE of one table backs
// up, with the statement's kind, nil for another statement.
func backupTarget(n ast.Node) (*ast.TableRef, string) {
	switch v := n.(type) {
	case *ast.UpdateStmt:
		if len(v.Tables) == 1 {
			if ref, ok := v.Tables[0].(*ast.TableRef); ok {
				return ref, "UPDATE"
			}
		}
	case *ast.DeleteStmt:
		if len(v.Using) == 0 && len(v.Tables) == 1 {
			if ref, ok := v.Tables[0].(*ast.TableRef); ok {
				return ref, "DELETE"
			}
		}
	}
	return nil, ""
}

// changeMade is what the change makes, which the backup, run before all of
// it, cannot see: the tables and views it creates or renames to, and the
// columns it adds, renames to, or declares on a table it creates, by the
// table as written. A statement that may make nothing, IF NOT EXISTS or OR
// REPLACE, is left out.
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
	for i := range stmts {
		switch v := stmts[i].node.(type) {
		case *ast.CreateTableStmt:
			if v.IfNotExists || v.Table == nil {
				continue
			}
			m.relations = append(m.relations, v.Table)
			for _, col := range v.Columns {
				if col != nil {
					m.columns = append(m.columns, madeColumn{v.Table, col.Name})
				}
			}
		case *ast.CreateViewStmt:
			if !v.OrReplace && v.Name != nil {
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
						m.relations = append(m.relations, cmd.NewTable)
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
// anything else, and a name the backup's query reads that only the change
// makes.
func checkPriorBackupTarget(sql string, stmts []statement, t review.Target, x schemaIndex, m changeMade, r *reporter) {
	if !t.BackupDatabaseExists {
		r.report(review.PriorBackup, -1, review.Range{}, "prior backup needs database "+backupDatabase+", which does not exist")
	}
	if x.schema == nil || usesDatabase(stmts, len(stmts)) {
		return
	}
	for i := range stmts {
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
// before the change, read in the WHERE clause, a join condition, or ORDER
// BY.
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
	var flatten func(e ast.TableExpr)
	flatten = func(e ast.TableExpr) {
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
		case *ast.JoinClause:
			flatten(v.Left)
			flatten(v.Right)
			if on, ok := v.Condition.(*ast.OnCondition); ok && on.Expr != nil {
				quals = append(quals, on.Expr)
			}
		default:
			known = false
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
