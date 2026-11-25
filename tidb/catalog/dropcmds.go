package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func (c *Catalog) dropTable(stmt *nodes.DropTableStmt) error {
	// The statement is judged whole before anything goes: v8.5.5 drops every
	// table it names that exists and answers 1051 for the first that does not
	// (DROP TABLE a, missing, b leaves neither a nor b), and a foreign key of
	// a table the statement drops too holds no parent back (DROP TABLE p, c1,
	// c2 runs where DROP TABLE p, c1 is 3730 for c2's key).
	type target struct {
		db  *Database
		key string
		ref string
		tbl *Table
	}
	var targets []target
	var missing error
	dropping := map[string]bool{}
	for _, ref := range stmt.Tables {
		db, err := c.resolveDatabase(ref.Schema)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		// DROP TABLE takes the local temporary table that shadows a table of
		// the name first, DROP TEMPORARY TABLE only a local temporary table,
		// and DROP GLOBAL TEMPORARY TABLE only a global one that no local
		// temporary table shadows: anything else of the name is 8007 there
		// (read on v8.5.5).
		var tbl *Table
		switch {
		case stmt.Temporary:
			tbl = db.TempTables[key]
		case stmt.Global:
			if db.TempTables[key] == nil && db.Tables[key] == nil && db.Views[key] == nil && db.Sequences[key] == nil {
				break
			}
			if tbl = db.Tables[key]; db.TempTables[key] != nil || tbl == nil || !tbl.GlobalTemporary {
				return errDropNotGlobalTemporary()
			}
		default:
			tbl = db.GetTable(key)
		}
		if tbl == nil {
			if !stmt.IfExists && missing == nil {
				missing = errUnknownTable(db.Name, ref.Name)
			}
			continue
		}
		// A cached table cannot be dropped; NOCACHE it first (8242, read on
		// v8.5.5).
		if tbl.Cached {
			return errCachedTable("Drop Table")
		}
		targets = append(targets, target{db, key, ref.Name, tbl})
		if !tbl.Temporary {
			dropping[toLower(db.Name)+"."+key] = true
		}
	}
	// Check if any other table in any database has a FK referencing this table
	// (unless foreign_key_checks=0).
	if c.foreignKeyChecks {
		c.droppingTables = dropping
		defer func() { c.droppingTables = nil }()
		for _, t := range targets {
			// No foreign key references a local temporary table.
			if t.tbl.Temporary {
				continue
			}
			if err := c.checkFKReferences(t.db.Name, t.ref); err != nil {
				return err
			}
		}
	}
	// DROP TEMPORARY TABLE is all or nothing: `DROP TEMPORARY TABLE t, nope`
	// is 1051 and leaves t (read on v8.5.5).
	if stmt.Temporary && missing != nil {
		return missing
	}
	for _, t := range targets {
		if t.tbl.Temporary {
			delete(t.db.TempTables, t.key)
		} else {
			delete(t.db.Tables, t.key)
		}
	}
	return missing
}

// checkFKReferences returns an error if any table in any database has a
// foreign key constraint that references the given table.
func (c *Catalog) checkFKReferences(dbName, tableName string) error {
	dbKey := toLower(dbName)
	tblKey := toLower(tableName)
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			// Skip the table itself, and the others the statement drops.
			if toLower(db.Name) == dbKey && toLower(tbl.Name) == tblKey {
				continue
			}
			if c.droppingTables[toLower(db.Name)+"."+toLower(tbl.Name)] {
				continue
			}
			for _, con := range tbl.Constraints {
				if con.Type != ConForeignKey {
					continue
				}
				refDB := con.RefDatabase
				if refDB == "" {
					refDB = db.Name
				}
				if toLower(refDB) == dbKey && toLower(con.RefTable) == tblKey {
					return errFKCannotDropParent(tableName, con.Name, tbl.Name)
				}
			}
		}
	}
	return nil
}

func (c *Catalog) truncateTable(stmt *nodes.TruncateStmt) error {
	for _, ref := range stmt.Tables {
		dbName := ref.Schema
		db, err := c.resolveDatabase(dbName)
		if err != nil {
			return err
		}
		tbl := db.GetTable(ref.Name)
		if tbl == nil {
			return errNoSuchTable(db.Name, ref.Name)
		}
		if tbl.Cached {
			return errCachedTable("Truncate Table")
		}
		// A table another table's foreign key references is not truncated
		// while foreign_key_checks is on (v8.5.5: TRUNCATE TABLE customers
		// under orders' fk_customer is 1701). A key of the table to itself
		// does not count, and no key references a local temporary table.
		if c.foreignKeyChecks && !tbl.Temporary {
			for _, ref := range c.referencingForeignKeys(tbl) {
				if ref.child == tbl {
					continue
				}
				return &Error{Code: 1701, SQLState: "42000", Message: fmt.Sprintf(
					"Cannot truncate a table referenced in a foreign key constraint (`%s`.`%s` CONSTRAINT `%s`)",
					ref.child.Database.Name, ref.child.Name, ref.con.Name)}
			}
		}
		tbl.AutoIncrement = 0
	}
	return nil
}
