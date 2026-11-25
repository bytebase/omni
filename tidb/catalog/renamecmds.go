package catalog

import (
	"fmt"

	nodes "github.com/bytebase/omni/tidb/ast"
)

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
		seq          *Sequence
		newName      string
	}
	// tables, views and sequences are the namespaces the pairs see, one
	// shallow copy per database they touch, so a refused statement leaves the
	// catalog's own maps alone.
	tables := map[*Database]map[string]*Table{}
	views := map[*Database]map[string]*View{}
	sequences := map[*Database]map[string]*Sequence{}
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
	seqNS := func(db *Database) map[string]*Sequence {
		if m, ok := sequences[db]; ok {
			return m
		}
		m := make(map[string]*Sequence, len(db.Sequences))
		for k, v := range db.Sequences {
			m[k] = v
		}
		sequences[db] = m
		return m
	}
	var renames []rename
	for _, pair := range stmt.Pairs {
		oldDB, err := c.resolveDatabase(pair.Old.Schema)
		if err != nil {
			return err
		}
		oldKey := toLower(pair.Old.Name)
		if oldDB.TempTables[oldKey] != nil {
			return errLocalTemporaryDDL("RENAME TABLE")
		}
		tbl := tableNS(oldDB)[oldKey]
		var view *View
		var seq *Sequence
		if tbl == nil {
			view = viewNS(oldDB)[oldKey]
		}
		if tbl == nil && view == nil {
			seq = seqNS(oldDB)[oldKey]
		}
		if tbl == nil && view == nil && seq == nil {
			return errNoSuchTable(oldDB.Name, pair.Old.Name)
		}
		if tbl != nil && tbl.Cached {
			return errCachedTable(op)
		}

		newDB, err := c.resolveDatabase(pair.New.Schema)
		if err != nil {
			// A table moved to a schema that does not exist fails where the
			// rename is applied (v8.5.5: RENAME TABLE accounts TO
			// simulation_nope.accounts is 1025, not 1049).
			if tbl != nil && pair.New.Schema != "" && c.GetDatabase(pair.New.Schema) == nil {
				return &Error{Code: 1025, SQLState: "HY000", Message: fmt.Sprintf(
					"Error on rename of '%s.%s' to '%s.%s' (errno: 168 - Database `%s` doesn't exist)",
					oldDB.Name, pair.Old.Name, pair.New.Schema, pair.New.Name, pair.New.Schema)}
			}
			return err
		}
		if view != nil && newDB != oldDB {
			return errRenameViewAcrossDB(oldDB.Name, newDB.Name)
		}
		newKey := toLower(pair.New.Name)
		if tableNS(newDB)[newKey] != nil || viewNS(newDB)[newKey] != nil || seqNS(newDB)[newKey] != nil {
			return errDupTableQualified(newDB.Name, pair.New.Name)
		}
		switch {
		case tbl != nil:
			delete(tableNS(oldDB), oldKey)
			tableNS(newDB)[newKey] = tbl
		case view != nil:
			delete(viewNS(oldDB), oldKey)
			viewNS(newDB)[newKey] = view
		default:
			delete(seqNS(oldDB), oldKey)
			seqNS(newDB)[newKey] = seq
		}
		renames = append(renames, rename{oldDB: oldDB, newDB: newDB, oldKey: oldKey, tbl: tbl, view: view, seq: seq, newName: pair.New.Name})
	}
	for db, m := range tables {
		db.Tables = m
	}
	for db, m := range views {
		db.Views = m
	}
	for db, m := range sequences {
		db.Sequences = m
	}
	for _, r := range renames {
		switch {
		case r.tbl != nil:
			oldName := r.tbl.Name
			r.tbl.Name = r.newName
			r.tbl.Database = r.newDB
			// The foreign keys that reference the table follow it, pair by
			// pair, and one that moves to another schema keeps naming the
			// parents it left (fk_retarget.go).
			c.retargetReferences(r.tbl, r.oldDB, oldName, nil, nil)
		case r.view != nil:
			r.view.Name = r.newName
			r.view.Database = r.newDB
		case r.seq != nil:
			r.seq.Name = r.newName
			r.seq.Database = r.newDB
		}
	}
	return nil
}
