package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return err
		}
		oldKey := toLower(pair.Old.Name)
		tbl := oldDB.Tables[oldKey]
		if tbl == nil {
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil {
			return errDupTable(pair.New.Name)
		}

		delete(oldDB.Tables, oldKey)
		tbl.Name = pair.New.Name
		tbl.Database = newDB
		newDB.Tables[newKey] = tbl
	}
	return nil
}
