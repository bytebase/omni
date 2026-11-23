package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
)

func (c *Catalog) dropTable(stmt *nodes.DropTableStmt) error {
	// A foreign key of a table the statement drops too does not hold a parent
	// back (v8.5.5: DROP TABLE p, c1, c2 runs where DROP TABLE p, c1 is 3730
	// for c2's key).
	dropping := map[string]bool{}
	for _, ref := range stmt.Tables {
		if db, err := c.resolveDatabase(ref.Schema); err == nil {
			dropping[toLower(db.Name)+"."+toLower(ref.Name)] = true
		}
	}
	c.droppingTables = dropping
	defer func() { c.droppingTables = nil }()
	for _, ref := range stmt.Tables {
		dbName := ref.Schema
		db, err := c.resolveDatabase(dbName)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		tbl := db.Tables[key]
		if tbl == nil {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		// A cached table cannot be dropped; NOCACHE it first (8242, read on
		// v8.5.5).
		if tbl.Cached {
			return errCachedTable("Drop Table")
		}
		// Check if any other table in any database has a FK referencing this table
		// (unless foreign_key_checks=0).
		if c.foreignKeyChecks {
			if err := c.checkFKReferences(db.Name, ref.Name); err != nil {
				return err
			}
		}
		delete(db.Tables, key)
	}
	return nil
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
		// does not count.
		if c.foreignKeyChecks {
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
