package catalog

import nodes "github.com/bytebase/omni/oceanbase/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	type plan struct {
		oldDB, newDB   *Database
		oldKey, newKey string
		newName        string
	}
	for _, pair := range stmt.Pairs {
		if err := rejectRemoteRef(pair.Old); err != nil {
			return err
		}
		if err := rejectRemoteRef(pair.New); err != nil {
			return err
		}
	}
	plans := make([]plan, 0, len(stmt.Pairs))

	// Simulate every pair in order against copies of the touched maps, so a
	// later pair sees the names an earlier one moved and a failed statement
	// leaves the catalog as it was.
	tables := map[*Database]map[string]*Table{}
	views := map[*Database]map[string]*View{}
	logs := map[*Database]map[string]bool{}
	known := map[*Database]map[string]bool{}
	tablesOf := func(db *Database) map[string]*Table {
		if tables[db] == nil {
			tables[db] = make(map[string]*Table, len(db.Tables))
			for k, v := range db.Tables {
				tables[db][k] = v
			}
		}
		return tables[db]
	}
	viewsOf := func(db *Database) map[string]*View {
		if views[db] == nil {
			views[db] = make(map[string]*View, len(db.Views))
			for k, v := range db.Views {
				views[db][k] = v
			}
		}
		return views[db]
	}
	logsOf := func(db *Database) map[string]bool {
		if logs[db] == nil {
			logs[db] = make(map[string]bool, len(db.MaterializedViewLogs))
			for k, v := range db.MaterializedViewLogs {
				logs[db][k] = v
			}
		}
		return logs[db]
	}
	knownOf := func(db *Database) map[string]bool {
		if known[db] == nil {
			known[db] = make(map[string]bool, len(db.MaterializedViewLogsKnown))
			for k, v := range db.MaterializedViewLogsKnown {
				known[db][k] = v
			}
		}
		return known[db]
	}

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
		oldKey := toLower(pair.Old.Name)
		fromTables, toTables := tablesOf(oldDB), tablesOf(newDB)
		fromViews, toViews := viewsOf(oldDB), viewsOf(newDB)
		if toTables[newKey] != nil || toViews[newKey] != nil {
			return errDupTable(pair.New.Name)
		}
		tbl := fromTables[oldKey]
		view := fromViews[oldKey]
		if tbl == nil {
			// The catalog holds a materialized view as the view its query is;
			// RENAME TABLE moves one like a base table.
			if view == nil || !view.Materialized {
				return &Error{Code: 1017, SQLState: "HY000", Message: "File not exist"}
			}
			delete(fromViews, oldKey)
			toViews[newKey] = view
		} else {
			delete(fromTables, oldKey)
			toTables[newKey] = tbl
		}
		fromLogs, toLogs := logsOf(oldDB), logsOf(newDB)
		fromKnown, toKnown := knownOf(oldDB), knownOf(newDB)
		if fromLogs[oldKey] {
			delete(fromLogs, oldKey)
			toLogs[newKey] = true
		}
		if fromKnown[oldKey] {
			delete(fromKnown, oldKey)
			toKnown[newKey] = true
		}
		plans = append(plans, plan{oldDB: oldDB, newDB: newDB, oldKey: oldKey, newKey: newKey, newName: pair.New.Name})
	}

	for _, p := range plans {
		if view := p.oldDB.Views[p.oldKey]; view != nil && view.Materialized {
			delete(p.oldDB.Views, p.oldKey)
			view.Name = p.newName
			view.Database = p.newDB
			p.newDB.Views[p.newKey] = view
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
