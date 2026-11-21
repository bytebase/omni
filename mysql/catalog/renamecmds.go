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
		view := oldDB.Views[oldKey]
		if tbl == nil && view == nil {
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		if view != nil && newDB != oldDB {
			return errForbidSchemaChange(oldDB.Name, newDB.Name)
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil || newDB.Views[newKey] != nil {
			return errDupTable(pair.New.Name)
		}

		if view != nil {
			delete(oldDB.Views, oldKey)
			view.Name = pair.New.Name
			view.Database = newDB
			newDB.Views[newKey] = view
			continue
		}

		oldDBName := oldDB.Name
		oldTableName := tbl.Name
		delete(oldDB.Tables, oldKey)
		tbl.Name = pair.New.Name
		tbl.Database = newDB
		newDB.Tables[newKey] = tbl
		stripTableQualifiersFromTableExprs(tbl, oldDBName, oldTableName, true, nil)
	}
	return nil
}
