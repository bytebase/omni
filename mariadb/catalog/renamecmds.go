package catalog

import nodes "github.com/bytebase/omni/mariadb/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return err
		}
		oldKey := toLower(pair.Old.Name)
		tbl := oldDB.Tables[oldKey]
		seq := oldDB.Sequences[oldKey]
		view := oldDB.Views[oldKey]
		if tbl == nil && seq == nil && view == nil {
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil || newDB.Sequences[newKey] != nil || newDB.Views[newKey] != nil {
			return errDupTable(pair.New.Name)
		}
		if view != nil {
			if oldDB != newDB {
				return errChangeSchemaNotAllowed(oldDB.Name, newDB.Name)
			}
			delete(oldDB.Views, oldKey)
			view.Name = pair.New.Name
			view.Database = newDB
			newDB.Views[newKey] = view
			continue
		}
		if seq != nil {
			delete(oldDB.Sequences, oldKey)
			seq.Name = pair.New.Name
			seq.Database = newDB
			newDB.Sequences[newKey] = seq
			continue
		}

		delete(oldDB.Tables, oldKey)
		tbl.Name = pair.New.Name
		tbl.Database = newDB
		newDB.Tables[newKey] = tbl
	}
	return nil
}
