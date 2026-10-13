package review

import (
	"strings"

	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/review"
)

// constraintKinds names the constraint kinds DisallowDropConstraint
// covers.
var constraintKinds = map[ast.ConstraintType]string{
	ast.ConstraintPrimaryKey: "primary key",
	ast.ConstraintUnique:     "unique constraint",
	ast.ConstraintForeignKey: "foreign key",
	ast.ConstraintCheck:      "check constraint",
}

// droppedConstraint is an ALTER TABLE ... DROP CONSTRAINT the SQL alone
// cannot classify: what it drops depends on the constraint, which the rule
// looks up in the target's synced schema.
type droppedConstraint struct {
	// statement is the index in Result.Statements, and position the
	// statement's place in the parsed list.
	statement int
	position  int
	loc       ast.Loc
	alter     *ast.AlterTableStmt
	name      string
	// nested is a drop inside a routine body or a control-flow statement,
	// which RequirePrimaryKey does not judge.
	nested bool
}

// checkConstraints reports, from the SQL alone, what DisallowDropConstraint
// and RequirePrimaryKey can establish:
//
//   - DisallowDropConstraint: an ALTER TABLE ... DROP CONSTRAINT of a
//     primary key, unique, foreign key, or check constraint the change
//     defined by that name on the table earlier, at top level or inside a
//     routine body or a control-flow statement.
//   - RequirePrimaryKey: a CREATE TABLE without a primary key, at top
//     level or in a CREATE SCHEMA, and the drop of a primary key the
//     change defined, each when the change does not settle the table later
//     by giving it a key, dropping it or its database, renaming it, or
//     creating a table of its name. A temporary table, a table variable,
//     and a table SELECT INTO or CREATE TABLE AS SELECT makes are not
//     judged, and neither is a statement inside a routine body.
//
// A DROP CONSTRAINT of a name the change did not define is returned, for
// the target's synced schema to settle.
func checkConstraints(stmts []statement, on map[review.Rule]bool, r *reporter) []droppedConstraint {
	type defined struct {
		table *ast.TableRef
		name  string
		kind  ast.ConstraintType
	}
	var made []defined
	find := func(ref *ast.TableRef, name string) (defined, bool) {
		for i := len(made) - 1; i >= 0; i-- {
			if sameObject(made[i].table, ref) && sameName(made[i].name, name) {
				return made[i], true
			}
		}
		return defined{}, false
	}
	var open []droppedConstraint
	drop := func(pos int, s *statement, v *ast.AlterTableStmt, a *ast.AlterTableAction, nested bool) {
		if a.Constraint == nil || a.Constraint.Name == "" {
			return
		}
		d, ok := find(v.Name, a.Constraint.Name)
		if !ok {
			open = append(open, droppedConstraint{s.index, pos, a.Loc, v, a.Constraint.Name, nested})
			return
		}
		reportDrop(stmts, pos, s, a.Loc, v, d.kind, a.Constraint.Name, nested, on, r)
	}
	// visit reads a statement, or an element of a CREATE SCHEMA, which runs
	// as the statement does; loc is what a finding on it anchors.
	var visit func(pos int, s *statement, n ast.Node, loc ast.Loc)
	visit = func(pos int, s *statement, n ast.Node, loc ast.Loc) {
		switch v := n.(type) {
		case *ast.CreateTableStmt:
			for _, c := range tableConstraints(v) {
				if c.Name != "" {
					made = append(made, defined{v.Name, c.Name, c.Type})
				}
			}
			// A FileTable has a primary key in its fixed schema.
			if on[review.RequirePrimaryKey] && judgedTable(v.Name) && !v.IsFileTable && !createsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name) {
				r.report(review.RequirePrimaryKey, s.index, rangeOf(loc), "creates table "+object(v.Name)+" without a primary key")
			}
		case *ast.AlterTableStmt:
			if v.Actions == nil {
				return
			}
			for _, item := range v.Actions.Items {
				a, ok := item.(*ast.AlterTableAction)
				if !ok {
					continue
				}
				switch a.Type {
				case ast.ATAddConstraint:
					if a.Constraint != nil && a.Constraint.Name != "" {
						made = append(made, defined{v.Name, a.Constraint.Name, a.Constraint.Type})
					}
				case ast.ATDropConstraint:
					drop(pos, s, v, a, false)
				}
			}
		case *ast.UseStmt:
			// What the change defined before is of another database.
			made = nil
		case *ast.CreateSchemaStmt:
			if v.Elements != nil {
				for _, e := range v.Elements.Items {
					visit(pos, s, e, ast.NodeLoc(e))
				}
			}
		default:
			// What a routine body or a control-flow statement adds runs
			// only when it runs, so only its drops are read.
			ast.Inspect(n, func(n ast.Node) bool {
				v, ok := n.(*ast.AlterTableStmt)
				if !ok || v.Actions == nil {
					return true
				}
				for _, item := range v.Actions.Items {
					if a, ok := item.(*ast.AlterTableAction); ok && a.Type == ast.ATDropConstraint {
						drop(pos, s, v, a, true)
					}
				}
				return true
			})
		}
	}
	for pos := range stmts {
		visit(pos, &stmts[pos], stmts[pos].node, stmts[pos].loc)
	}
	return open
}

// reportDrop reports the drop of a constraint of a known kind: under
// DisallowDropConstraint when the kind is one it covers, and under
// RequirePrimaryKey when it is the primary key, the drop is not nested,
// and the change does not settle the table later.
func reportDrop(stmts []statement, pos int, s *statement, loc ast.Loc, v *ast.AlterTableStmt, kind ast.ConstraintType, name string, nested bool, on map[review.Rule]bool, r *reporter) {
	label, ok := constraintKinds[kind]
	if !ok {
		return
	}
	if on[review.DisallowDropConstraint] {
		r.report(review.DisallowDropConstraint, s.index, rangeOf(loc), "drops "+label+" "+ident(name)+" of "+object(v.Name))
	}
	if on[review.RequirePrimaryKey] && !nested && kind == ast.ConstraintPrimaryKey && !addsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name) {
		r.report(review.RequirePrimaryKey, s.index, rangeOf(loc), "removes the primary key of "+object(v.Name)+", and the change adds none back")
	}
}

// judgedTable reports whether RequirePrimaryKey judges a table of the
// name: not a temporary table, whose name starts with #.
func judgedTable(ref *ast.TableRef) bool {
	return ref != nil && !strings.HasPrefix(ref.Object, "#")
}

// tableConstraints lists a CREATE TABLE's constraints, of the table and of
// its columns.
func tableConstraints(v *ast.CreateTableStmt) []*ast.ConstraintDef {
	var out []*ast.ConstraintDef
	add := func(l *ast.List) {
		if l == nil {
			return
		}
		for _, item := range l.Items {
			if c, ok := item.(*ast.ConstraintDef); ok {
				out = append(out, c)
			}
		}
	}
	add(v.Constraints)
	if v.Columns != nil {
		for _, item := range v.Columns.Items {
			if col, ok := item.(*ast.ColumnDef); ok {
				add(col.Constraints)
			}
		}
	}
	return out
}

// createsPrimaryKey reports whether a CREATE TABLE declares a primary key.
func createsPrimaryKey(v *ast.CreateTableStmt) bool {
	for _, c := range tableConstraints(v) {
		if c.Type == ast.ConstraintPrimaryKey {
			return true
		}
	}
	return false
}

// addsPrimaryKey reports whether an ALTER TABLE may give the table a
// primary key: ADD CONSTRAINT ... PRIMARY KEY, or a column added with
// PRIMARY KEY.
func addsPrimaryKey(v *ast.AlterTableStmt) bool {
	if v.Actions == nil {
		return false
	}
	for _, item := range v.Actions.Items {
		a, ok := item.(*ast.AlterTableAction)
		if !ok {
			continue
		}
		switch a.Type {
		case ast.ATAddConstraint:
			if a.Constraint != nil && a.Constraint.Type == ast.ConstraintPrimaryKey {
				return true
			}
		case ast.ATAddColumn:
			if a.Column != nil && a.Column.Constraints != nil {
				for _, c := range a.Column.Constraints.Items {
					if cd, ok := c.(*ast.ConstraintDef); ok && cd.Type == ast.ConstraintPrimaryKey {
						return true
					}
				}
			}
		}
	}
	return false
}

// settledAfter reports whether a statement after the given position may
// give the table a primary key, drop it or its database, rename it, or
// create a table of its name: then the rule does not follow it further.
func settledAfter(stmts []statement, pos int, ref *ast.TableRef) bool {
	for i := pos + 1; i < len(stmts); i++ {
		switch v := stmts[i].node.(type) {
		case *ast.AlterTableStmt:
			if sameObject(v.Name, ref) && addsPrimaryKey(v) {
				return true
			}
		case *ast.CreateTableStmt:
			if sameObject(v.Name, ref) {
				return true
			}
		case *ast.DropStmt:
			if v.Names == nil {
				continue
			}
			for _, item := range v.Names.Items {
				t, ok := item.(*ast.TableRef)
				switch {
				case !ok:
				case v.ObjectType == ast.DropTable && sameObject(t, ref):
					return true
				// A table without a database may be in the one dropped.
				case v.ObjectType == ast.DropDatabase && (ref.Database == "" || sameName(ref.Database, t.Object)):
					return true
				}
			}
		case *ast.RenameStmt:
			if v.ColumnName == "" && sameObject(v.Name, ref) {
				return true
			}
		case *ast.ExecStmt:
			if name, _, _, ok := spRename(v); ok {
				if t := parseObjectName(name); t != nil && sameObject(t, ref) {
					return true
				}
			}
		}
	}
	return false
}
