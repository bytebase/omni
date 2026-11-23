package catalog

import (
	"fmt"
	"strings"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

func (c *Catalog) dropTable(stmt *nodes.DropTableStmt) error {
	// A table another table's foreign key references is not dropped, and the
	// statement goes on to the next: measured on mariadb:11.8.8 over
	// features@fk, DROP TABLE customers, orders, items answers 1451 and leaves
	// customers and orders, and DROP TABLE items, customers leaves customers
	// and orders. Each table is judged once the ones before it went.
	var refused error
	for _, ref := range stmt.Tables {
		dbName := ref.Schema
		db, err := c.resolveDatabase(dbName)
		if err != nil {
			if stmt.IfExists {
				continue
			}
			return err
		}
		key := toLower(ref.Name)
		// DROP TABLE drops a sequence too. Measured on mariadb:11.8.8.
		// DROP TEMPORARY TABLE drops only a temporary one.
		if seq := db.Sequences[key]; seq != nil && db.Tables[key] == nil && (!stmt.Temporary || seq.Temporary) {
			delete(db.Sequences, key)
			continue
		}
		if db.Tables[key] == nil {
			if stmt.IfExists {
				continue
			}
			if !stmt.Temporary && db.Views[key] != nil {
				return errIsAView(db.Name, ref.Name)
			}
			return errUnknownTable(db.Name, ref.Name)
		}
		// Check if any other table in any database has a FK referencing this table
		// (unless foreign_key_checks=0).
		if c.foreignKeyChecks {
			if child, _ := c.findFKReference(db.Name, ref.Name); child != nil {
				if refused == nil {
					refused = errRowIsReferenced()
				}
				continue
			}
		}
		delete(db.Tables, key)
	}
	return refused
}

// findFKReference returns a table other than the given one, in any
// database, with a foreign key that references the given table, and that
// foreign key; nil when there is none. A foreign key of the table to
// itself does not count: dropping the table drops it too.
func (c *Catalog) findFKReference(dbName, tableName string) (*Table, *Constraint) {
	dbKey := toLower(dbName)
	tblKey := toLower(tableName)
	for _, db := range c.databases {
		for _, tbl := range db.Tables {
			// Skip the table itself.
			if toLower(db.Name) == dbKey && toLower(tbl.Name) == tblKey {
				continue
			}
			for _, con := range tbl.Constraints {
				if con.Type != ConForeignKey {
					continue
				}
				refDB := con.RefDatabase
				if refDB == "" {
					refDB = db.Name
				}
				if toLower(refDB) == dbKey && toLower(con.RefTable) == tblKey {
					return tbl, con
				}
			}
		}
	}
	return nil, nil
}

func (c *Catalog) truncateTable(stmt *nodes.TruncateStmt) error {
	for _, ref := range stmt.Tables {
		dbName := ref.Schema
		db, err := c.resolveDatabase(dbName)
		if err != nil {
			return err
		}
		tbl := db.GetTable(ref.Name)
		if tbl == nil {
			if db.Sequences[toLower(ref.Name)] != nil {
				return &Error{Code: 1031, SQLState: "HY000",
					Message: fmt.Sprintf("Storage engine SEQUENCE of the table `%s`.`%s` doesn't have this option", db.Name, ref.Name)}
			}
			return errNoSuchTable(db.Name, ref.Name)
		}
		if tbl.SystemVersioned {
			return errTruncateSystemVersioned(ref.Name)
		}
		// A table another table's foreign key references is not truncated while
		// foreign_key_checks is on (mariadb:11.8.8: TRUNCATE TABLE customers
		// under orders' fk_customer is 1701). A foreign key to itself does not
		// count.
		if c.foreignKeyChecks {
			if child, con := c.findFKReference(db.Name, tbl.Name); child != nil {
				return errTruncateReferenced(child, con, db.Name, tbl.Name)
			}
		}
		tbl.AutoIncrement = 0
		// TRUNCATE gives the table a tablespace of its own again, after an
		// 1814 warning where DISCARD TABLESPACE took it away: measured on
		// mariadb:11.8.8, `DISCARD TABLESPACE`, `TRUNCATE TABLE a` and
		// `IMPORT TABLESPACE` answers 1813 (the tablespace is there again), not
		// 1030.
		tbl.TablespaceDiscarded = false
	}
	return nil
}

// errTruncateReferenced is 1701, TRUNCATE of a table con of child references,
// with the foreign key written as MariaDB writes it.
func errTruncateReferenced(child *Table, con *Constraint, dbName, table string) error {
	quote := func(names []string) string {
		out := make([]string, len(names))
		for i, n := range names {
			out[i] = backtickIDQ(n)
		}
		return strings.Join(out, ",")
	}
	return &Error{Code: 1701, SQLState: "42000", Message: fmt.Sprintf(
		"Cannot truncate a table referenced in a foreign key constraint (%s.%s, CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s (%s))",
		backtickIDQ(child.Database.Name), backtickIDQ(child.Name), backtickIDQ(con.Name), quote(con.Columns),
		backtickIDQ(dbName), backtickIDQ(table), quote(con.RefColumns))}
}
