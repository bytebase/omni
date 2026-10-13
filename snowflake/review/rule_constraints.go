package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/snowflake/ast"
)

// constraintKind is what a constraint is, as DisallowDropConstraint names
// it. unknownKind is a name the change touched in a way the rule does not
// follow.
type constraintKind int

const (
	unknownKind constraintKind = iota
	primaryKey
	uniqueKey
	foreignKey
	checkKey
)

func (k constraintKind) label() string {
	switch k {
	case primaryKey:
		return "primary key"
	case uniqueKey:
		return "unique constraint"
	case foreignKey:
		return "foreign key"
	case checkKey:
		return "check constraint"
	}
	return "constraint"
}

func kindOf(t ast.ConstraintType) constraintKind {
	switch t {
	case ast.ConstrPrimaryKey:
		return primaryKey
	case ast.ConstrUnique:
		return uniqueKey
	case ast.ConstrForeignKey:
		return foreignKey
	}
	return unknownKind
}

// droppedConstraint is an ALTER TABLE ... DROP CONSTRAINT the change alone
// does not classify: what it drops depends on the constraint, which the
// rule looks up in the target's synced schema.
type droppedConstraint struct {
	// statement is the index in Result.Statements.
	statement int
	loc       ast.Loc
	table     *ast.ObjectName
	name      ast.Ident
	// judgePK is whether RequirePrimaryKey judges the drop when it drops
	// the primary key: it is a top-level statement, its ALTER TABLE adds no
	// primary key back, and no later statement settles the table.
	judgePK bool
}

// reportDrop reports the drop of a constraint of a known kind, under
// DisallowDropConstraint, and under RequirePrimaryKey when it is the
// primary key and the drop is judged.
func reportDrop(d droppedConstraint, kind constraintKind, on map[review.Rule]bool, r *reporter) {
	if on[review.DisallowDropConstraint] {
		r.report(review.DisallowDropConstraint, d.statement, rangeOf(d.loc), "drops "+kind.label()+" "+ident(d.name)+" of "+object(d.table))
	}
	if on[review.RequirePrimaryKey] && kind == primaryKey && d.judgePK {
		r.report(review.RequirePrimaryKey, d.statement, rangeOf(d.loc), "removes the primary key of "+object(d.table)+", and the change adds none back")
	}
}

// constraintEvent is one thing the change did that a later DROP CONSTRAINT
// reads back.
type constraintEvent struct {
	// barrier is a statement after which a name may mean another object:
	// a USE, or a DROP or UNDROP of a schema or database.
	barrier bool
	// table is the table the event is about. With name empty, the table
	// was replaced: created, dropped, restored, renamed, or swapped.
	table *ast.ObjectName
	// name is the normalized constraint name, and kind what it now is.
	name string
	kind constraintKind
}

// constraintLog is what the change did to constraint names, in order.
type constraintLog struct {
	events []constraintEvent
}

func (c *constraintLog) barrier() {
	c.events = append(c.events, constraintEvent{barrier: true})
}

func (c *constraintLog) replace(table *ast.ObjectName) {
	if table != nil {
		c.events = append(c.events, constraintEvent{table: table})
	}
}

func (c *constraintLog) define(table *ast.ObjectName, name ast.Ident, kind constraintKind) {
	if table != nil && !name.IsEmpty() {
		c.events = append(c.events, constraintEvent{table: table, name: name.Normalize(), kind: kind})
	}
}

// lookup reads back what the change last did to the table's constraint of
// the name: decided is true when the change settles it, with kind
// unknownKind when the rule cannot tell. A definition counts only on a
// table written alike; a replacement of a table that may be this one, or a
// barrier, settles the name as unknown.
func (c *constraintLog) lookup(table *ast.ObjectName, name ast.Ident) (kind constraintKind, decided bool) {
	norm := name.Normalize()
	for i := len(c.events) - 1; i >= 0; i-- {
		e := c.events[i]
		switch {
		case e.barrier:
			return unknownKind, true
		case e.name == "":
			if sameObject(e.table, table) {
				return unknownKind, true
			}
		case e.name == norm && sameWritten(e.table, table):
			return e.kind, true
		}
	}
	return unknownKind, false
}

// checkConstraints reports, from the SQL alone, what DisallowDropConstraint
// and RequirePrimaryKey can establish:
//
//   - DisallowDropConstraint: ALTER TABLE ... DROP PRIMARY KEY, DROP UNIQUE,
//     and DROP FOREIGN KEY, and DROP CONSTRAINT of a constraint the change
//     defined by that name on the table earlier, at top level or inside a
//     Snowflake Scripting block, a task's body, or an alert's action.
//   - RequirePrimaryKey: a CREATE TABLE without a primary key, and the drop
//     of a primary key, each a top-level statement the change does not
//     settle later by giving the table a key, dropping it or its schema or
//     database, renaming or swapping it, or creating a table of its name.
//     CREATE OR REPLACE TABLE makes a new table and is judged; a CTAS is
//     judged by what its column list declares, since only names and types
//     come from its query. Not judged: a temporary table, which lasts only
//     for its session (a transient one persists and is judged); CREATE TABLE
//     ... LIKE and CLONE, which copy the source's structure; IF NOT EXISTS,
//     which may make nothing; CREATE OR ALTER TABLE, which alters a table
//     that may already lack a key; and a statement inside a block, which
//     runs only as its control flow decides. Snowflake does not enforce a
//     primary key on a standard table, but the rule still asks for one.
//
// What the change defines is read only from a top-level statement; a
// nested one may not run, so its definitions only hide the name. A DROP
// CONSTRAINT of a name the change did not touch, on a table the change did
// not replace, before any USE or schema or database drop, is returned, for
// the target's synced schema to settle.
func checkConstraints(stmts []statement, on map[review.Rule]bool, r *reporter) []droppedConstraint {
	var c constraintLog
	var open []droppedConstraint
	for pos := range stmts {
		s := &stmts[pos]
		ast.Inspect(s.node, func(n ast.Node) bool {
			top := n == s.node
			switch v := n.(type) {
			case *ast.CreateTableStmt:
				if v.IfNotExists {
					return true
				}
				c.replace(v.Name)
				if top && v.Like == nil && v.Clone == nil {
					for _, d := range tableConstraints(v) {
						c.define(v.Name, d.name, d.kind)
					}
				}
				if top && on[review.RequirePrimaryKey] && judgedTable(v) && !createsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name) {
					r.report(review.RequirePrimaryKey, s.index, rangeOf(s.loc), "creates table "+object(v.Name)+" without a primary key")
				}
			case *ast.AlterTableStmt:
				open = append(open, c.alterTable(stmts, pos, s, v, top, on, r)...)
			case *ast.DropStmt:
				if v.Kind == ast.DropTable {
					c.replace(v.Name)
				}
			case *ast.UndropStmt:
				if v.Kind == ast.UndropTable {
					c.replace(v.Name)
				}
			case *ast.DropSchemaStmt, *ast.DropDatabaseStmt, *ast.UndropSchemaStmt, *ast.UndropDatabaseStmt:
				c.barrier()
			case *ast.UseStmt:
				if movesNames(v) {
					c.barrier()
				}
			}
			return true
		})
	}
	return open
}

// alterTable reads one ALTER TABLE, action by action: what it defines,
// renames, or replaces, and what it drops.
func (c *constraintLog) alterTable(stmts []statement, pos int, s *statement, v *ast.AlterTableStmt, top bool, on map[review.Rule]bool, r *reporter) []droppedConstraint {
	// A nested statement may not run: what it defines is not known.
	defined := func(k constraintKind) constraintKind {
		if top {
			return k
		}
		return unknownKind
	}
	var open []droppedConstraint
	for _, a := range v.Actions {
		if a == nil {
			continue
		}
		switch a.Kind {
		case ast.AlterTableAddConstraint:
			if a.Constraint != nil {
				c.define(v.Name, a.Constraint.Name, defined(kindOf(a.Constraint.Type)))
			}
		case ast.AlterTableAddColumn:
			for _, col := range a.Columns {
				if col != nil && col.InlineConstraint != nil {
					c.define(v.Name, col.InlineConstraint.Name, defined(kindOf(col.InlineConstraint.Type)))
				}
			}
		case ast.AlterTableRenameConstraint:
			kind, _ := c.lookup(v.Name, a.ConstraintName)
			c.define(v.Name, a.ConstraintName, unknownKind)
			c.define(v.Name, a.NewConstraintName, defined(kind))
		case ast.AlterTableRename, ast.AlterTableSwapWith:
			c.replace(v.Name)
			c.replace(a.NewName)
		case ast.AlterTableDropConstraint:
			judgePK := top && !addsPrimaryKey(v) && !settledAfter(stmts, pos, v.Name)
			switch {
			case a.IsPrimaryKey:
				if on[review.DisallowDropConstraint] {
					r.report(review.DisallowDropConstraint, s.index, rangeOf(a.Loc), "drops the primary key of "+object(v.Name))
				}
				if on[review.RequirePrimaryKey] && judgePK {
					r.report(review.RequirePrimaryKey, s.index, rangeOf(a.Loc), "removes the primary key of "+object(v.Name)+", and the change adds none back")
				}
			case a.DropUnique:
				if on[review.DisallowDropConstraint] {
					r.report(review.DisallowDropConstraint, s.index, rangeOf(a.Loc), "drops a unique constraint of "+object(v.Name))
				}
			case a.DropForeignKey:
				if on[review.DisallowDropConstraint] {
					r.report(review.DisallowDropConstraint, s.index, rangeOf(a.Loc), "drops a foreign key of "+object(v.Name))
				}
			case !a.ConstraintName.IsEmpty():
				d := droppedConstraint{statement: s.index, loc: a.Loc, table: v.Name, name: a.ConstraintName, judgePK: judgePK}
				kind, decided := c.lookup(v.Name, a.ConstraintName)
				switch {
				case !decided:
					open = append(open, d)
				case kind != unknownKind:
					reportDrop(d, kind, on, r)
				}
				c.define(v.Name, a.ConstraintName, unknownKind)
			}
		}
	}
	return open
}

// namedConstraint is a constraint a CREATE TABLE names.
type namedConstraint struct {
	name ast.Ident
	kind constraintKind
}

// tableConstraints lists a CREATE TABLE's named constraints, of the table
// and of its columns.
func tableConstraints(v *ast.CreateTableStmt) []namedConstraint {
	var out []namedConstraint
	for _, col := range v.Columns {
		if col != nil && col.InlineConstraint != nil && !col.InlineConstraint.Name.IsEmpty() {
			out = append(out, namedConstraint{col.InlineConstraint.Name, kindOf(col.InlineConstraint.Type)})
		}
	}
	for _, tc := range v.Constraints {
		if tc != nil && !tc.Name.IsEmpty() {
			out = append(out, namedConstraint{tc.Name, kindOf(tc.Type)})
		}
	}
	return out
}

// judgedTable reports whether RequirePrimaryKey judges a CREATE TABLE: not
// a temporary table (TEMPORARY or one of its synonyms, VOLATILE among
// them), not IF NOT EXISTS, not CREATE OR ALTER, and not LIKE or CLONE.
func judgedTable(v *ast.CreateTableStmt) bool {
	return v.Name != nil && !v.Temporary && !v.Volatile && !v.IfNotExists && !v.OrAlter && v.Like == nil && v.Clone == nil
}

// createsPrimaryKey reports whether a CREATE TABLE declares a primary key,
// as a table constraint or on a column.
func createsPrimaryKey(v *ast.CreateTableStmt) bool {
	for _, tc := range v.Constraints {
		if tc != nil && tc.Type == ast.ConstrPrimaryKey {
			return true
		}
	}
	for _, col := range v.Columns {
		if col != nil && col.InlineConstraint != nil && col.InlineConstraint.Type == ast.ConstrPrimaryKey {
			return true
		}
	}
	return false
}

// addsPrimaryKey reports whether an ALTER TABLE may give the table a
// primary key: ADD PRIMARY KEY, or a column added with PRIMARY KEY.
func addsPrimaryKey(v *ast.AlterTableStmt) bool {
	for _, a := range v.Actions {
		if a == nil {
			continue
		}
		switch a.Kind {
		case ast.AlterTableAddConstraint:
			if a.Constraint != nil && a.Constraint.Type == ast.ConstrPrimaryKey {
				return true
			}
		case ast.AlterTableAddColumn:
			for _, col := range a.Columns {
				if col != nil && col.InlineConstraint != nil && col.InlineConstraint.Type == ast.ConstrPrimaryKey {
					return true
				}
			}
		}
	}
	return false
}

// settledAfter reports whether a statement after the given position, at
// top level or nested, may give the table a primary key, drop it or its
// schema or database, rename or swap it, or create a table of its name:
// then the rule does not follow it further.
func settledAfter(stmts []statement, pos int, ref *ast.ObjectName) bool {
	settled := false
	for i := pos + 1; i < len(stmts) && !settled; i++ {
		ast.Inspect(stmts[i].node, func(n ast.Node) bool {
			switch v := n.(type) {
			case *ast.CreateTableStmt:
				settled = settled || sameObject(v.Name, ref)
			case *ast.AlterTableStmt:
				if sameObject(v.Name, ref) && addsPrimaryKey(v) {
					settled = true
				}
				for _, a := range v.Actions {
					if a != nil && (a.Kind == ast.AlterTableRename || a.Kind == ast.AlterTableSwapWith) &&
						(sameObject(v.Name, ref) || sameObject(a.NewName, ref)) {
						settled = true
					}
				}
			case *ast.DropStmt:
				settled = settled || v.Kind == ast.DropTable && sameObject(v.Name, ref)
			// A table without a schema or database may be in the one dropped.
			case *ast.DropSchemaStmt:
				settled = settled || v.Name != nil && (ref.Schema.IsEmpty() || sameName(ref.Schema, v.Name.Name))
			case *ast.DropDatabaseStmt:
				settled = settled || v.Name != nil && (ref.Database.IsEmpty() || sameName(ref.Database, v.Name.Name))
			}
			return !settled
		})
	}
	return settled
}

// movesNames reports whether a USE moves what a name means: USE DATABASE,
// USE SCHEMA, and USE with no object kind, which names a database or a
// schema. USE ROLE, USE WAREHOUSE, and USE SECONDARY ROLES do not.
//
// sql: https://docs.snowflake.com/en/sql-reference/sql/use-database
func movesNames(v *ast.UseStmt) bool {
	switch v.Kind {
	case ast.UseDefault, ast.UseDatabase, ast.UseSchema:
		return true
	}
	return false
}
