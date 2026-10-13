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
	table     *ast.TableRef
	name      string
}

// checkConstraints reports, from the SQL alone, what DisallowDropConstraint
// and RequirePrimaryKey can establish:
//
//   - DisallowDropConstraint: an ALTER TABLE ... DROP CONSTRAINT of a
//     primary key, unique, foreign key, or check constraint the change
//     defined by that name on the table earlier.
//   - RequirePrimaryKey: a CREATE TABLE without a primary key, and the drop
//     of a primary key the change defined, each when the change does not
//     settle the table later by giving it a key, dropping or renaming it,
//     or creating a table of its name. A temporary table, a table variable,
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
	for pos := range stmts {
		s := &stmts[pos]
		switch v := s.node.(type) {
		case *ast.CreateTableStmt:
			for _, c := range tableConstraints(v) {
				if c.Name != "" {
					made = append(made, defined{v.Name, c.Name, c.Type})
				}
			}
			if on[review.RequirePrimaryKey] && judgedTable(v.Name) && !createsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name) {
				r.report(review.RequirePrimaryKey, s.index, rangeOf(s.loc), "creates table "+object(v.Name)+" without a primary key")
			}
		case *ast.AlterTableStmt:
			if v.Actions == nil {
				continue
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
					if a.Constraint == nil || a.Constraint.Name == "" {
						continue
					}
					d, ok := find(v.Name, a.Constraint.Name)
					if !ok {
						open = append(open, droppedConstraint{s.index, pos, a.Loc, v.Name, a.Constraint.Name})
						continue
					}
					reportDrop(stmts, pos, s, a.Loc, v, d.kind, a.Constraint.Name, on, r)
				}
			}
		}
	}
	return open
}

// reportDrop reports the drop of a constraint of a known kind: under
// DisallowDropConstraint when the kind is one it covers, and under
// RequirePrimaryKey when it is the primary key and the change does not
// settle the table later.
func reportDrop(stmts []statement, pos int, s *statement, loc ast.Loc, v *ast.AlterTableStmt, kind ast.ConstraintType, name string, on map[review.Rule]bool, r *reporter) {
	label, ok := constraintKinds[kind]
	if !ok {
		return
	}
	if on[review.DisallowDropConstraint] {
		r.report(review.DisallowDropConstraint, s.index, rangeOf(loc), "drops "+label+" "+ident(name)+" of "+object(v.Name))
	}
	if on[review.RequirePrimaryKey] && kind == ast.ConstraintPrimaryKey && !addsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name) {
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
// give the table a primary key, drop it, rename it, or create a table of
// its name: then the rule does not follow it further.
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
			if v.ObjectType == ast.DropTable && v.Names != nil {
				for _, item := range v.Names.Items {
					if t, ok := item.(*ast.TableRef); ok && sameObject(t, ref) {
						return true
					}
				}
			}
		case *ast.RenameStmt:
			if sameObject(v.Name, ref) {
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
