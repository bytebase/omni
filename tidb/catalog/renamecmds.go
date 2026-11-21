package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	// RENAME TABLE is one atomic statement: v8.5.5 answers 8242 for
	// `RENAME TABLE ordinary TO moved, cached TO other` and leaves both tables
	// under their own names, 1146 for a pair naming no table and 1050 for a
	// destination another pair already took (`RENAME TABLE a TO x, b TO x`
	// names x twice). The pairs are resolved one after another against the
	// namespace the pairs before them leave -- `RENAME TABLE a TO a2, a2 TO a`
	// is accepted, because a2 holds a by the time the second pair reads it --
	// and nothing is committed until every pair holds.
	op := "Rename Table"
	if len(stmt.Pairs) > 1 {
		// TiDB names the plural form for several pairs.
		op = "Rename Tables"
	}
	type rename struct {
		oldDB, newDB *Database
		oldKey       string
		tbl          *Table
		view         *View
		newName      string
	}
	// tables and views are the namespaces the pairs see, one shallow copy per
	// database they touch, so a refused statement leaves the catalog's own maps
	// alone.
	tables := map[*Database]map[string]*Table{}
	views := map[*Database]map[string]*View{}
	tableNS := func(db *Database) map[string]*Table {
		if m, ok := tables[db]; ok {
			return m
		}
		m := make(map[string]*Table, len(db.Tables))
		for k, v := range db.Tables {
			m[k] = v
		}
		tables[db] = m
		return m
	}
	viewNS := func(db *Database) map[string]*View {
		if m, ok := views[db]; ok {
			return m
		}
		m := make(map[string]*View, len(db.Views))
		for k, v := range db.Views {
			m[k] = v
		}
		views[db] = m
		return m
	}
	var renames []rename
	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return err
		}
		oldKey := toLower(pair.Old.Name)
		tbl := tableNS(oldDB)[oldKey]
		var view *View
		if tbl == nil {
			view = viewNS(oldDB)[oldKey]
		}
		if tbl == nil && view == nil {
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}
		if tbl != nil && tbl.Cached {
			return errCachedTable(op)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		if view != nil && newDB != oldDB {
			return errRenameViewAcrossDB(oldDB.Name, newDB.Name)
		}
		newKey := toLower(pair.New.Name)
		if tableNS(newDB)[newKey] != nil || viewNS(newDB)[newKey] != nil || newDB.Sequences[newKey] != nil {
			return errDupTableQualified(newDB.Name, pair.New.Name)
		}
		if tbl != nil {
			delete(tableNS(oldDB), oldKey)
			tableNS(newDB)[newKey] = tbl
		} else {
			delete(viewNS(oldDB), oldKey)
			viewNS(newDB)[newKey] = view
		}
		renames = append(renames, rename{oldDB: oldDB, newDB: newDB, oldKey: oldKey, tbl: tbl, view: view, newName: pair.New.Name})
	}
	for db, m := range tables {
		db.Tables = m
	}
	for db, m := range views {
		db.Views = m
	}
	for _, r := range renames {
		if r.tbl != nil {
			r.tbl.Name = r.newName
			r.tbl.Database = r.newDB
		} else if r.view != nil {
			r.view.Name = r.newName
			r.view.Database = r.newDB
		}
	}
	return nil
}
