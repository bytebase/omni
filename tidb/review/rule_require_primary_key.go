package review

import (
	"github.com/bytebase/omni/review"
	"github.com/bytebase/omni/tidb/ast"
)

// checkRequirePrimaryKey reports two things, each only when the change
// ends that way: a table a CREATE TABLE made without a primary key, and a
// table whose primary key an ALTER TABLE drops without adding one back. A
// DROP PRIMARY KEY only runs on a table that has one, so neither needs the
// synced schema, and a table that had none before the change is never
// reported. A table made by CREATE TABLE ... LIKE or AS SELECT, whose
// columns and keys come from elsewhere, a CREATE TABLE IF NOT EXISTS,
// which may make nothing, and a temporary table are not judged, and
// neither is a statement inside a routine body, which runs only when
// called.
//
// Without simulating the change, the rule follows the table by name. A
// later statement that may give it a key, drop it or its database, rename
// it, or create a table of its name settles it, and it is not reported.
func checkRequirePrimaryKey(stmts []statement, r *reporter) {
	type pending struct {
		ref     *ast.TableRef
		index   int
		rng     review.Range
		message string
		settled bool
	}
	var list []*pending
	settle := func(ref *ast.TableRef) {
		for _, p := range list {
			if !p.settled && sameTable(p.ref, ref) {
				p.settled = true
			}
		}
	}
	for i := range stmts {
		s := &stmts[i]
		switch v := s.node.(type) {
		case *ast.CreateTableStmt:
			settle(v.Table)
			if v.Table == nil || v.Temporary || v.IfNotExists || v.Like != nil || v.Select != nil || createsPrimaryKey(v) {
				continue
			}
			list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "creates table " + table(v.Table) + " without a primary key"})
		case *ast.AlterTableStmt:
			if v.Table == nil {
				continue
			}
			adds, renames := false, false
			var dropped *ast.AlterTableCmd
			for _, cmd := range v.Commands {
				if cmd == nil {
					continue
				}
				switch {
				case addsPrimaryKey(cmd):
					adds = true
				case cmd.Type == ast.ATRenameTable:
					renames = true
				case dropped == nil && cmd.Type == ast.ATDropConstraint && sameName(cmd.Name, "PRIMARY"),
					dropped == nil && cmd.Type == ast.ATDropIndex && sameName(cmd.Name, "PRIMARY"):
					dropped = cmd
				}
			}
			if adds || renames {
				settle(v.Table)
			}
			if dropped != nil && !adds && !renames {
				list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(dropped.Loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
		case *ast.DropIndexStmt:
			if sameName(v.Name, "PRIMARY") {
				list = append(list, &pending{ref: v.Table, index: s.index, rng: rangeOf(s.loc), message: "removes the primary key of " + table(v.Table) + ", and the change adds none back"})
			}
		case *ast.DropTableStmt:
			for _, ref := range v.Tables {
				settle(ref)
			}
		case *ast.RenameTableStmt:
			for _, pair := range v.Pairs {
				if pair != nil {
					settle(pair.Old)
					settle(pair.New)
				}
			}
		case *ast.DropDatabaseStmt:
			// A table without a database may be in the one dropped.
			for _, p := range list {
				if p.ref != nil && (p.ref.Schema == "" || sameName(p.ref.Schema, v.Name)) {
					p.settled = true
				}
			}
		}
	}
	for _, p := range list {
		if !p.settled {
			r.report(review.RequirePrimaryKey, p.index, p.rng, p.message)
		}
	}
}

// createsPrimaryKey reports whether a CREATE TABLE declares a primary key,
// as a table constraint or on a column.
func createsPrimaryKey(v *ast.CreateTableStmt) bool {
	for _, c := range v.Constraints {
		if c != nil && c.Type == ast.ConstrPrimaryKey {
			return true
		}
	}
	for _, col := range v.Columns {
		if columnIsPrimaryKey(col) {
			return true
		}
	}
	return false
}

// addsPrimaryKey reports whether an ALTER TABLE subcommand may give the
// table a primary key: ADD PRIMARY KEY, or a column added, modified, or
// changed with PRIMARY KEY.
func addsPrimaryKey(cmd *ast.AlterTableCmd) bool {
	switch cmd.Type {
	case ast.ATAddConstraint, ast.ATAddIndex:
		return cmd.Constraint != nil && cmd.Constraint.Type == ast.ConstrPrimaryKey
	case ast.ATAddColumn, ast.ATModifyColumn, ast.ATChangeColumn:
		if columnIsPrimaryKey(cmd.Column) {
			return true
		}
		for _, col := range cmd.Columns {
			if columnIsPrimaryKey(col) {
				return true
			}
		}
	}
	return false
}

func columnIsPrimaryKey(col *ast.ColumnDef) bool {
	if col == nil {
		return false
	}
	for _, c := range col.Constraints {
		if c != nil && c.Type == ast.ColConstrPrimaryKey {
			return true
		}
	}
	return false
}
