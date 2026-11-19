package catalog

import nodes "github.com/bytebase/omni/tidb/ast"

func (c *Catalog) renameTable(stmt *nodes.RenameTableStmt) error {
	// RENAME TABLE is one atomic statement: v8.5.5 answers 8242 for
	// `RENAME TABLE ordinary TO moved, cached TO other` and leaves both tables
	// under their own names, and 1146 for a pair naming no table the same way
	// (read on v8.5.5). So every pair is resolved and judged before any is
	// applied.
	op := "Rename Table"
	if len(stmt.Pairs) > 1 {
		// TiDB names the plural form for several pairs.
		op = "Rename Tables"
	}
	type rename struct {
		oldDB, newDB *Database
		oldKey       string
		tbl          *Table
		newName      string
	}
	var renames []rename
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
		if tbl.Cached {
			return errCachedTable(op)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			return err
		}
		newKey := toLower(pair.New.Name)
		if newDB.Tables[newKey] != nil || newDB.Sequences[newKey] != nil {
			return errDupTable(pair.New.Name)
		}
		renames = append(renames, rename{oldDB: oldDB, newDB: newDB, oldKey: oldKey, tbl: tbl, newName: pair.New.Name})
	}
	for _, r := range renames {
		delete(r.oldDB.Tables, r.oldKey)
		r.tbl.Name = r.newName
		r.tbl.Database = r.newDB
		r.newDB.Tables[toLower(r.newName)] = r.tbl
	}
	return nil
}
