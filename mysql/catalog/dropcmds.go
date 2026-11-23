package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/mysql/ast"
)

func (c *Catalog) dropTable(stmt *nodes.DropTableStmt) error {
	// DROP TABLE is atomic: every table is found before any is dropped, and a
	// table the statement drops does not hold back the parents it references,
	// so DROP TABLE parent, child drops both (mysql:8.0.46).
	type target struct {
		db  *Database
		key string
		tbl *Table
	}
	var targets []target
	dropping := map[*Table]bool{}
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
		targets = append(targets, target{db, key, tbl})
		dropping[tbl] = true
	}
	// Check if any other table in any database has a FK referencing this table
	// (unless foreign_key_checks=0).
	if c.foreignKeyChecks {
		for _, t := range targets {
			if childTable, conName, ok := c.findReferencingFKExcept(t.db.Name, t.tbl.Name, dropping); ok {
				return errFKCannotDropParent(t.tbl.Name, conName, childTable)
			}
		}
	}
	for _, t := range targets {
		delete(t.db.Tables, t.key)
	}
	return nil
}

// checkFKReferences returns an error if any table in any database has a
func (c *Catalog) findReferencingFK(dbName, tableName string) (childTable, conName string, found bool) {
	return c.findReferencingFKExcept(dbName, tableName, nil)
}

// findReferencingFKExcept is findReferencingFK over the tables except the
// ones skip holds.
func (c *Catalog) findReferencingFKExcept(dbName, tableName string, skip map[*Table]bool) (childTable, conName string, found bool) {
	dbKey := toLower(dbName)
	tblKey := toLower(tableName)
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			// Skip the table itself.
			if toLower(db.Name) == dbKey && toLower(tbl.Name) == tblKey || skip[tbl] {
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
					return tbl.Name, con.Name, true
				}
			}
		}
	}
	return "", "", false
}

// foreign key constraint that references the given table.
func (c *Catalog) checkFKReferences(dbName, tableName string) error {
	if childTable, conName, ok := c.findReferencingFK(dbName, tableName); ok {
		return errFKCannotDropParent(tableName, conName, childTable)
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
		if tableHasDiscardedTablespace(tbl) {
			return errTablespaceDiscarded(db.Name, tbl.Name)
		}
		// A table another table's foreign key references is not truncated
		// while foreign_key_checks is on (1701); its own references to itself
		// do not count (mysql:8.0.46).
		if c.foreignKeyChecks {
			for _, ref := range c.referencingForeignKeys(tbl) {
				if ref.child == tbl {
					continue
				}
				return &Error{Code: 1701, SQLState: "42000",
					Message: fmt.Sprintf("Cannot truncate a table referenced in a foreign key constraint (`%s`.`%s`, CONSTRAINT `%s`)", ref.child.Database.Name, ref.child.Name, ref.con.Name)}
			}
		}
		tbl.AutoIncrement = 0
	}
	return nil
}
