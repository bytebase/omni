package catalog

import nodes "github.com/bytebase/omni/mysql/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	type renameUndo struct {
		oldDB    *Database
		newDB    *Database
		oldKey   string
		newKey   string
		view     *View
		oldName  string
		tbl      *Table
		snapshot Table
	}
	var applied []renameUndo
	rollback := func(err error) error {
		for i := len(applied) - 1; i >= 0; i-- {
			u := applied[i]
			if u.view != nil {
				delete(u.newDB.Views, u.newKey)
				u.view.Name = u.oldName
				u.view.Database = u.oldDB
				u.oldDB.Views[u.oldKey] = u.view
				continue
			}
			delete(u.newDB.Tables, u.newKey)
			*u.tbl = u.snapshot
			u.oldDB.Tables[u.oldKey] = u.tbl
		}
		return err
	}

	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return rollback(err)
		}
		oldKey := toLower(pair.Old.Name)
		tbl := oldDB.Tables[oldKey]
		view := oldDB.Views[oldKey]
		if tbl == nil && view == nil {
			return rollback(errNoSuchTable(oldDB.Name, pair.Old.Name))
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return rollback(err)
		}
		if view != nil && newDB != oldDB {
			return rollback(errForbidSchemaChange(oldDB.Name, newDB.Name))
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil || newDB.Views[newKey] != nil {
			return rollback(errDupTable(pair.New.Name))
		}

		if view != nil {
			applied = append(applied, renameUndo{
				oldDB:   oldDB,
				newDB:   newDB,
				oldKey:  oldKey,
				newKey:  newKey,
				view:    view,
				oldName: view.Name,
			})
			delete(oldDB.Views, oldKey)
			view.Name = pair.New.Name
			view.Database = newDB
			newDB.Views[newKey] = view
			continue
		}

		applied = append(applied, renameUndo{
			oldDB:    oldDB,
			newDB:    newDB,
			oldKey:   oldKey,
			newKey:   newKey,
			tbl:      tbl,
			snapshot: cloneTable(tbl),
		})
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
