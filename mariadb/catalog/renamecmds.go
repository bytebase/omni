package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

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
			// RENAME TABLE IF EXISTS skips a missing table with a note
			// (mariadb:11.8.8).
			if stmt.IfExists {
				continue
			}
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			// A table moved to a schema that does not exist fails where InnoDB
			// renames its file: measured on mariadb:11.8.8, RENAME TABLE fp TO
			// nodb.fp answers 1025 (errno: 168), not 1049.
			if tbl != nil && pair.New.Schema != "" && c.GetDatabase(pair.New.Schema) == nil {
				return &Error{Code: 1025, SQLState: "HY000", Message: fmt.Sprintf(
					"Error on rename of './%s/%s' to './%s/%s' (errno: 168 \"Unknown (generic) error from engine\")",
					oldDB.Name, tbl.Name, pair.New.Schema, pair.New.Name)}
			}
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

		oldName := tbl.Name
		delete(oldDB.Tables, oldKey)
		tbl.Name = pair.New.Name
		tbl.Database = newDB
		newDB.Tables[newKey] = tbl
		// The foreign keys that reference the table follow it, and one that
		// moves to another schema keeps naming the parents it left
		// (mariadb:11.8.8, schemas@cross).
		c.retargetReferences(tbl, oldDB, oldName, nil, nil)
	}
	return nil
}
