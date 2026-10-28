package review

import (
	"slices"

	"github.com/bytebase/omni/mariadb/ast"
	"github.com/bytebase/omni/review"
)

// droppedIndex is a DROP INDEX the SQL alone cannot classify: whether it
// drops a unique key depends on the index, which the rule looks up in the
// target's synced schema.
type droppedIndex struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	loc       ast.Loc
	table     *ast.TableRef
	name      string
}

// checkDisallowDropConstraint reports the drop of a primary key, foreign
// key, unique key, or check constraint: ALTER TABLE ... DROP PRIMARY KEY,
// DROP FOREIGN KEY, and DROP CONSTRAINT, which drops a foreign key, unique
// key, or check constraint of the name, and DROP INDEX of the
// primary key, at top level or inside a routine body. A DROP INDEX of
// another name drops a unique key when the change made a unique key of
// that name on the table earlier, without IF NOT EXISTS. One of a name the
// change has not touched before is returned, for the target's synced
// schema to settle;
// one the change created, renamed, or dropped before is not judged, and
// neither is one on a table the change creates, without IF NOT EXISTS, or
// renames to, whose indexes are not the synced table's. A table named in
// another case, or with a database where the earlier statement named
// none or another, than the earlier statement wrote it is not judged:
// where table names are case-sensitive, or the database differs, it is
// another table. What a routine,
// trigger, or event body does runs only when it runs, so the body's own
// drops are reported, except one of an index the body itself touched or
// on a table it made, and the body touches no index the rest of the
// change reads.
func checkDisallowDropConstraint(sql string, stmts []statement, r *reporter) []droppedIndex {
	type index struct {
		table  *ast.TableRef
		name   string
		unique bool
		// unsure is an index a routine, trigger, or event body touched:
		// which of its statements ran before a later one depends on the
		// body's control flow, so a drop of it is not judged.
		unsure bool
	}
	var touched []index
	find := func(ref *ast.TableRef, name string) (index, bool) {
		for i := len(touched) - 1; i >= 0; i-- {
			if sameTable(touched[i].table, ref) && sameName(touched[i].name, name) {
				return touched[i], true
			}
		}
		return index{}, false
	}
	// inBody is whether the statement read is a routine, trigger, or
	// event definition, and foreign whether that body resolves an
	// unqualified table in a database the definition names, which may not
	// be the target's.
	inBody, foreign := false, false
	touch := func(ref *ast.TableRef, name string, unique bool) {
		if name != "" {
			touched = append(touched, index{ref, name, unique, inBody})
		}
	}
	// forget drops what the change touched on a table of the name, which
	// is gone with the table.
	forget := func(ref *ast.TableRef) {
		kept := touched[:0]
		for _, i := range touched {
			if !sameTable(i.table, ref) {
				kept = append(kept, i)
			}
		}
		touched = kept
	}
	var replaced []*ast.TableRef
	replace := func(ref *ast.TableRef) {
		replaced = append(replaced, ref)
		forget(ref)
	}
	isReplaced := func(ref *ast.TableRef) bool {
		for _, t := range replaced {
			if sameTable(t, ref) {
				return true
			}
		}
		return false
	}
	var open []droppedIndex
	droppedName := func(pos int, s *statement, loc ast.Loc, ref *ast.TableRef, name string) {
		if foreign && ref.Schema == "" {
			return
		}
		i, ok := find(ref, name)
		switch {
		case ok && !sameCase(i.table, ref):
		case ok && i.unsure:
		case ok && i.unique:
			r.report(review.DisallowDropConstraint, s.index, rangeOf(loc), "drops unique key "+ident(name)+" of "+table(ref))
		case !ok && !isReplaced(ref):
			open = append(open, droppedIndex{s.index, pos, loc, ref, name})
		}
		touch(ref, name, false)
	}
	for pos := range stmts {
		s := &stmts[pos]
		// A routine, trigger, or event definition changes nothing until
		// its body runs: what the body does only keeps its own later
		// drops from being judged, then is forgotten.
		inBody, foreign = false, false
		var keptTouched []index
		var keptReplaced []*ast.TableRef
		switch v := s.node.(type) {
		case *ast.CreateFunctionStmt:
			foreign = v.Name != nil && v.Name.Schema != ""
		case *ast.CreateTriggerStmt:
			foreign = v.Table != nil && v.Table.Schema != ""
		}
		switch s.node.(type) {
		case *ast.CreateFunctionStmt, *ast.CreateTriggerStmt, *ast.CreateEventStmt, *ast.AlterEventStmt:
			inBody = true
			keptTouched, keptReplaced = slices.Clone(touched), slices.Clone(replaced)
		}
		ast.Inspect(s.node, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.UseStmt:
				// What the change touched before on a table named without a
				// database is of another database; one named with its
				// database still is what it was.
				touched = slices.DeleteFunc(touched, func(i index) bool { return i.table.Schema == "" })
				replaced = slices.DeleteFunc(replaced, func(t *ast.TableRef) bool { return t.Schema == "" })
			case *ast.CreateTableStmt:
				// IF NOT EXISTS may leave the table of the name as it was.
				if v.IfNotExists {
					break
				}
				replace(v.Table)
				for _, c := range v.Constraints {
					if c != nil && (c.Type == ast.ConstrUnique || c.Type == ast.ConstrPrimaryKey) {
						touch(v.Table, c.Name, c.Type == ast.ConstrUnique)
					}
				}
			case *ast.DropTableStmt:
				// DROP TEMPORARY drops only a temporary table, which may
				// have hidden the table of the name: that table is as it was.
				for _, ref := range v.Tables {
					if v.Temporary {
						forget(ref)
					} else {
						replace(ref)
					}
				}
			case *ast.CreateIndexStmt:
				// IF NOT EXISTS may leave an index of the name as it was.
				if !v.IfNotExists {
					touch(v.Table, v.IndexName, v.Unique)
				}
			case *ast.RenameTableStmt:
				for _, pair := range v.Pairs {
					if pair != nil {
						replace(pair.New)
					}
				}
			case *ast.DropIndexStmt:
				loc := v.Loc
				if rng := anchor(s, n, v.Loc); rng.End > rng.Start {
					loc = ast.Loc{Start: rng.Start, End: rng.End}
				}
				if sameName(v.Name, "PRIMARY") {
					r.report(review.DisallowDropConstraint, s.index, rangeOf(loc), "drops the primary key of "+table(v.Table))
				} else {
					droppedName(pos, s, loc, v.Table, v.Name)
				}
			case *ast.AlterTableStmt:
				for _, cmd := range v.Commands {
					if cmd == nil {
						continue
					}
					switch cmd.Type {
					case ast.ATAddConstraint, ast.ATAddIndex:
						// A check constraint makes no index; its name is
						// another namespace's.
						if cmd.Constraint != nil && cmd.Constraint.Type != ast.ConstrCheck {
							touch(v.Table, keyName(cmd.Constraint), cmd.Constraint.Type == ast.ConstrUnique)
						}
					case ast.ATRenameIndex:
						// The new name keeps what the change knows of the old.
						old, known := find(v.Table, cmd.Name)
						touch(v.Table, cmd.Name, false)
						touch(v.Table, cmd.NewName, known && sameCase(old.table, v.Table) && old.unique)
					case ast.ATRenameTable:
						// A new name without a database stays in the table's.
						to := &ast.TableRef{Name: cmd.NewName}
						if cmd.NewTable != nil {
							copied := *cmd.NewTable
							to = &copied
						}
						if to.Schema == "" {
							to.Schema = v.Table.Schema
						}
						replace(to)
					}
					var msg string
					switch dropKindOf(sql, cmd) {
					case dropPrimaryKey:
						msg = "drops the primary key of " + table(v.Table)
					case dropForeignKey:
						msg = "drops foreign key " + ident(cmd.Name) + " of " + table(v.Table)
					case dropConstraint:
						msg = "drops constraint " + ident(cmd.Name) + " of " + table(v.Table)
					case dropIndex:
						droppedName(pos, s, cmd.Loc, v.Table, cmd.Name)
					}
					if msg != "" {
						r.report(review.DisallowDropConstraint, s.index, rangeOf(cmd.Loc), msg)
					}
				}
			}
			return true
		})
		if inBody {
			touched, replaced = keptTouched, keptReplaced
		}
	}
	return open
}

// checkDroppedIndex reports a DROP INDEX of a unique key of the target's
// synced table. The table is found as MariaDB resolves it, and not at all
// when a USE before the statement may have moved the current database.
func checkDroppedIndex(stmts []statement, d droppedIndex, x schemaIndex, r *reporter) {
	t := x.table(d.table, usesDatabase(stmts, d.position))
	if t == nil {
		return
	}
	if unique, found := uniqueIndex(t, d.name); found && unique {
		r.report(review.DisallowDropConstraint, d.statement, rangeOf(d.loc), "drops unique key "+ident(d.name)+" of "+table(d.table))
	}
}
