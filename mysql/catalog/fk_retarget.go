package catalog

import (
	"slices"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// retargetReferences points the foreign keys that reference tbl at what a
// rename made of it: tbl was oldName in oldDB, and columns maps each column a
// statement renamed, by lowered old name, to its new name. InnoDB keeps a
// foreign key on its parent through renames, so after RENAME TABLE p TO q,
// ALTER TABLE p RENAME TO o.q (another schema) or ALTER TABLE p RENAME COLUMN
// id TO pid, SHOW CREATE TABLE of the child references `q`, `o`.`q` or
// (`pid`) (mysql:8.0.46). A table that moves to another schema also keeps
// referencing the parents it left behind, so its own references to tables of
// oldDB name oldDB. It returns what undoes the changes.
func (c *Catalog) retargetReferences(tbl *Table, oldDB *Database, oldName string, columns map[string]string) (undo func()) {
	// A change is saved by the table and the position of the constraint, not
	// the constraint: a rollback may first put back a table's snapshot, whose
	// constraints are copies, and it is those the undo must reach.
	type saved struct {
		child   *Table
		at      int
		db, tbl string
		cols    []string
	}
	var changed []saved
	moved := tbl.Database != oldDB
	for _, db := range c.databases {
		for _, child := range db.Tables {
			// The database a reference without a schema names is its table's,
			// which for tbl itself is the one it was in.
			owner := child.Database
			if child == tbl {
				owner = oldDB
			}
			for at, con := range child.Constraints {
				if con.Type != ConForeignKey {
					continue
				}
				save := func(con *Constraint) {
					changed = append(changed, saved{child, at, con.RefDatabase, con.RefTable, slices.Clone(con.RefColumns)})
				}
				refDB := con.RefDatabase
				if refDB == "" && owner != nil {
					refDB = owner.Name
				}
				if toLower(refDB) != toLower(oldDB.Name) {
					continue
				}
				if toLower(con.RefTable) != toLower(oldName) {
					// A table that moved keeps the parents of oldDB it names.
					if child == tbl && moved && con.RefDatabase == "" {
						save(con)
						con.RefDatabase = oldDB.Name
					}
					continue
				}
				save(con)
				con.RefTable = tbl.Name
				switch {
				case child.Database == tbl.Database && (moved || child == tbl):
					con.RefDatabase = ""
				case moved:
					con.RefDatabase = tbl.Database.Name
				}
				for i, col := range con.RefColumns {
					if renamed, ok := columns[toLower(col)]; ok {
						con.RefColumns[i] = renamed
					}
				}
			}
		}
	}
	return func() {
		for i := len(changed) - 1; i >= 0; i-- {
			s := changed[i]
			if s.at < len(s.child.Constraints) {
				con := s.child.Constraints[s.at]
				con.RefDatabase, con.RefTable, con.RefColumns = s.db, s.tbl, s.cols
			}
		}
	}
}

// fkRef is one foreign key, con, of the table child.
type fkRef struct {
	child *Table
	con   *Constraint
}

// referencingForeignKeys returns the foreign keys of every table, parent
// itself included, that reference parent.
func (c *Catalog) referencingForeignKeys(parent *Table) []fkRef {
	if parent == nil || parent.Database == nil {
		return nil
	}
	return c.referencingForeignKeysOf(parent.Database.Name, parent.Name)
}

// referencingForeignKeysOf returns the foreign keys that reference the table
// dbName.name: a table an ALTER TABLE renames is still referenced by its old
// name until the statement ends (retargetReferences).
func (c *Catalog) referencingForeignKeysOf(dbName, name string) []fkRef {
	var refs []fkRef
	for _, db := range c.databases {
		for _, child := range db.Tables {
			for _, con := range child.Constraints {
				if con.Type != ConForeignKey {
					continue
				}
				refDB := con.RefDatabase
				if refDB == "" && child.Database != nil {
					refDB = child.Database.Name
				}
				if toLower(refDB) == toLower(dbName) && toLower(con.RefTable) == toLower(name) {
					refs = append(refs, fkRef{child, con})
				}
			}
		}
	}
	return refs
}

// alterColumnRenames maps each column an ALTER TABLE renames, by lowered old
// name, to its last new name: RENAME COLUMN and CHANGE, followed through a
// later command that renames the new name again.
func alterColumnRenames(commands []*nodes.AlterTableCmd) map[string]string {
	var renames map[string]string
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}
		from, to := cmd.Name, ""
		switch cmd.Type {
		case nodes.ATRenameColumn:
			to = cmd.NewName
		case nodes.ATChangeColumn:
			if cmd.Column != nil {
				to = cmd.Column.Name
			}
		}
		if to == "" || from == to {
			continue
		}
		if renames == nil {
			renames = make(map[string]string)
		}
		followed := false
		for old, now := range renames {
			if toLower(now) == toLower(from) {
				renames[old] = to
				followed = true
			}
		}
		if !followed {
			renames[toLower(from)] = to
		}
	}
	return renames
}
