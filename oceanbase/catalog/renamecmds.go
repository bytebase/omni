package catalog

import nodes "github.com/bytebase/omni/oceanbase/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	type plan struct {
		oldDB, newDB     *Database
		oldKey, newKey   string
		newName          string
		oldName          string
		materializedView *View
	}
	plans := make([]plan, 0, len(stmt.Pairs))
	// Check every pair before moving any: a failed statement leaves the catalog
	// as it was.
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
		if newDB.Tables[newKey] != nil || newDB.Views[newKey] != nil {
			return errDupTable(pair.New.Name)
		}
		oldKey := toLower(pair.Old.Name)
		p := plan{oldDB: oldDB, newDB: newDB, oldKey: oldKey, newKey: newKey, newName: pair.New.Name, oldName: pair.Old.Name}
		if oldDB.Tables[oldKey] == nil {
			// The catalog holds a materialized view as the view its query is;
			// RENAME TABLE moves one like a base table.
			view := oldDB.Views[oldKey]
			if view == nil || !view.Materialized {
				return &Error{Code: 1017, SQLState: "HY000", Message: "File not exist"}
			}
			p.materializedView = view
		}
		plans = append(plans, p)
	}
	for _, p := range plans {
		if p.materializedView != nil {
			delete(p.oldDB.Views, p.oldKey)
			p.materializedView.Name = p.newName
			p.materializedView.Database = p.newDB
			p.newDB.Views[p.newKey] = p.materializedView
		} else {
			tbl := p.oldDB.Tables[p.oldKey]
			delete(p.oldDB.Tables, p.oldKey)
			tbl.Name = p.newName
			tbl.Database = p.newDB
			p.newDB.Tables[p.newKey] = tbl
		}
		if _, ok := p.oldDB.MaterializedViewLogs[p.oldKey]; ok {
			delete(p.oldDB.MaterializedViewLogs, p.oldKey)
			p.newDB.MaterializedViewLogs[p.newKey] = true
		}
		if p.oldDB.MaterializedViewLogsKnown[p.oldKey] {
			delete(p.oldDB.MaterializedViewLogsKnown, p.oldKey)
			p.newDB.MaterializedViewLogsKnown[p.newKey] = true
		}
	}
	return nil
}
