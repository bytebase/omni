package catalog

import nodes "github.com/bytebase/omni/oceanbase/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return err
		}
		// OceanBase checks the new name (1050) before the old one, and answers 1017 for
		// a table that is not there (TestOracle_OceanBaseNameOrder).
		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil {
			return errDupTable(pair.New.Name)
		}
		oldKey := toLower(pair.Old.Name)
		tbl := oldDB.Tables[oldKey]
		if tbl == nil {
			return &Error{Code: 1017, SQLState: "HY000", Message: "File not exist"}
		}

		delete(oldDB.Tables, oldKey)
		tbl.Name = pair.New.Name
		tbl.Database = newDB
		newDB.Tables[newKey] = tbl
	}
	return nil
}
