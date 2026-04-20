package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

func (c *Catalog) dropTable(stmt *nodes.DropTableStmt) error {
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
		if db.Tables[key] == nil {
			if stmt.IfExists {
				continue
			}
			return errUnknownTable(db.Name, ref.Name)
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
			// Skip the table itself.
			if toLower(db.Name) == dbKey && toLower(tbl.Name) == tblKey {
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
		tbl.AutoIncrement = 0
	}
	return nil
}
